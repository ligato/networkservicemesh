package prefixcollector

import (
	"context"
	"os"
	"strings"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/client-go/kubernetes"
	"k8s.io/kubernetes/cmd/kubeadm/app/apis/kubeadm/v1beta2"

	"github.com/networkservicemesh/networkservicemesh/sdk/prefix_pool"
)

const (
	// ExcludedPrefixesEnv is the name of the env variable to define excluded prefixes
	ExcludedPrefixesEnv = "EXCLUDED_PREFIXES"
)

func getExcludedPrefixesChan(clientset *kubernetes.Clientset) (<-chan prefix_pool.PrefixPool, error) {
	prefixes := getExcludedPrefixesFromEnv()

	// trying to get excludePrefixes from kubeadm-config, if it exists
	if configMapPrefixes, err := getExcludedPrefixesFromConfigMap(clientset); err == nil {
		poolCh := make(chan prefix_pool.PrefixPool, 1)
		pool, err := prefix_pool.NewPrefixPool(append(prefixes, configMapPrefixes...)...)
		if err != nil {
			return nil, err
		}
		poolCh <- pool
		return poolCh, nil
	}

	// seems like we don't have kubeadm-config in cluster, starting monitor client
	return monitorSubnets(clientset, prefixes...), nil
}

func getExcludedPrefixesFromEnv() []string {
	excludedPrefixesEnv, ok := os.LookupEnv(ExcludedPrefixesEnv)
	if !ok {
		return []string{}
	}
	logrus.Infof("Getting excludedPrefixes from ENV: %v", excludedPrefixesEnv)
	return strings.Split(excludedPrefixesEnv, ",")
}

func getExcludedPrefixesFromConfigMap(clientset *kubernetes.Clientset) ([]string, error) {
	kubeadmConfig, err := clientset.CoreV1().
		ConfigMaps("kube-system").
		Get(context.TODO(), "kubeadm-config", metav1.GetOptions{})
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	clusterConfiguration := &v1beta2.ClusterConfiguration{}
	err = yaml.NewYAMLOrJSONDecoder(strings.NewReader(kubeadmConfig.Data["ClusterConfiguration"]), 4096).
		Decode(clusterConfiguration)
	if err != nil {
		return nil, err
	}

	podSubnet := clusterConfiguration.Networking.PodSubnet
	serviceSubnet := clusterConfiguration.Networking.ServiceSubnet

	if podSubnet == "" {
		return nil, errors.New("ClusterConfiguration.Networking.PodSubnet is empty")
	}
	if serviceSubnet == "" {
		return nil, errors.New("ClusterConfiguration.Networking.ServiceSubnet is empty")
	}

	return []string{
		podSubnet,
		serviceSubnet,
	}, nil
}

func monitorSubnets(clientset *kubernetes.Clientset, additionalPrefixes ...string) <-chan prefix_pool.PrefixPool {
	logrus.Infof("Start monitoring prefixes to exclude")
	poolCh := make(chan prefix_pool.PrefixPool, 1)

	go func() {
		for {
			errCh := make(chan error)
			go monitorReservedSubnets(poolCh, errCh, clientset, additionalPrefixes)
			err := <-errCh
			logrus.Error(err)
		}
	}()

	return poolCh
}

func monitorReservedSubnets(poolCh chan prefix_pool.PrefixPool, errCh chan<- error, clientset *kubernetes.Clientset, additionalPrefixes []string) {
	pw, err := WatchPodCIDR(clientset)
	if err != nil {
		errCh <- err
		return
	}
	defer pw.Stop()

	sw, err := WatchServiceIpAddr(clientset)
	if err != nil {
		errCh <- err
		return
	}
	defer sw.Stop()

	var podSubnet, serviceSubnet string
	for {
		select {
		case subnet, ok := <-pw.ResultChan():
			if !ok {
				return
			}
			podSubnet = subnet.String()
		case subnet, ok := <-sw.ResultChan():
			if !ok {
				return
			}
			serviceSubnet = subnet.String()
		}
		sendPrefixPool(poolCh, podSubnet, serviceSubnet, additionalPrefixes)
	}
}

func sendPrefixPool(poolCh chan prefix_pool.PrefixPool, podSubnet, serviceSubnet string, additionalPrefixes []string) {
	pool, err := getPrefixPool(podSubnet, serviceSubnet, additionalPrefixes)
	if err != nil {
		logrus.Errorf("Failed to create a prefix pool: %v", err)
		return
	}
	select {
	case <-poolCh:
	default:
	}
	poolCh <- pool
}

func getPrefixPool(podSubnet, serviceSubnet string, additionalPrefixes []string) (prefix_pool.PrefixPool, error) {
	prefixes := additionalPrefixes
	if len(podSubnet) > 0 {
		prefixes = append(prefixes, podSubnet)
	}
	if len(serviceSubnet) > 0 {
		prefixes = append(prefixes, serviceSubnet)
	}

	pool, err := prefix_pool.NewPrefixPool(prefixes...)
	if err != nil {
		return nil, err
	}

	return pool, nil
}
