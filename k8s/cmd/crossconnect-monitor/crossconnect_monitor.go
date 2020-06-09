package main

import (
	"context"
	"flag"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"

	"github.com/networkservicemesh/networkservicemesh/controlplane/api/crossconnect"
	metricspkg "github.com/networkservicemesh/networkservicemesh/controlplane/pkg/metrics"
	"github.com/networkservicemesh/networkservicemesh/k8s/pkg/networkservice/clientset/versioned"
	"github.com/networkservicemesh/networkservicemesh/k8s/pkg/networkservice/namespace"
	"github.com/networkservicemesh/networkservicemesh/pkg/tools"
)

var closing = false
var managers = map[string]string{}

func monitorCrossConnects(address string, continuousMonitor bool, cache map[string]string) {
	var err error
	logrus.Infof("Starting CrossConnections Monitor on %s", address)
	conn, err := tools.DialTCP(address)
	if err != nil {
		logrus.Errorf("failure to communicate with the socket %s with error: %+v", address, err)
		return
	}
	defer conn.Close()
	forwarderClient := crossconnect.NewMonitorCrossConnectClient(conn)

	// Looping indefinitely or until grpc returns an error indicating the other end closed connection.
	stream, err := forwarderClient.MonitorCrossConnects(context.Background(), &empty.Empty{})

	if err != nil {
		logrus.Warningf("Error: %+v.", err)
		return
	}

	if _, ok := cache["prometheus"]; !ok {
		cache["prometheus"] = "off"
	}

	t := proto.TextMarshaler{}
	for {
		event, err := stream.Recv()
		if err != nil {
			logrus.Errorf("Error: %+v.", err)
			return
		}
		data := fmt.Sprintf("\u001b[31m*** %s\n\u001b[0m", event.Type)
		data += fmt.Sprintf("\u001b[31m*** %s\n\u001b[0m", address)
		for _, cc := range event.CrossConnects {
			if cc != nil {
				data += fmt.Sprintf("\u001b[32m%s\n\u001b[0m", t.Text(cc))
			}
		}
		println(data)

		if event.Type == crossconnect.CrossConnectEventType_UPDATE && cache["prometheus"] == "on" {
			trackMetrics(event)
		}
		if !continuousMonitor {
			logrus.Infof("Monitoring of server: %s. is complete...", address)
			delete(managers, address)
			return
		}
	}
}

func servePrometheus(metricsIdentifiers map[string]string, metricsData *crossconnect.Metrics) {
	srcPod := metricsIdentifiers[metricspkg.SrcPodKey]
	srcNamespace := metricsIdentifiers[metricspkg.SrcNamespaceKey]
	dstPod := metricsIdentifiers[metricspkg.DstPodKey]
	dstNamespace := metricsIdentifiers[metricspkg.DstNamespaceKey]
	if srcPod == "" || srcNamespace == "" || dstPod == "" || dstNamespace == "" {
		logrus.Infof("failed to serve Prometheus - connection data is missing; src and dst pod names and namespaces should not be nill %s, %s, %s, %s",
			srcPod, srcNamespace, dstPod, dstNamespace)
	} else {
		prometheusCtx := metricspkg.BuildPrometheusMetricsContext(srcPod, srcNamespace, dstPod, dstNamespace)
		prometheusMetrics := metricspkg.BuildPrometheusMetrics()
		metricspkg.CollectAllMetrics(prometheusCtx, prometheusMetrics, metricsData)
	}
}

func getCrossConnectMetricsMap(event *crossconnect.CrossConnectEvent) map[string]map[string]string {
	crossConnectMetricsMap := map[string]map[string]string{}
	for _, cc := range event.CrossConnects {
		ccID, metricsIdentifiers, err := metricspkg.GetMetricsIdentifiers(cc)

		if err != nil {
			logrus.Warningf("failed to get metric identifier: %v", err)
		}
		crossConnectMetricsMap[ccID] = metricsIdentifiers
	}
	return crossConnectMetricsMap
}

func trackMetrics(event *crossconnect.CrossConnectEvent) {
	crossConnectMetricsMap := getCrossConnectMetricsMap(event)
	if len(crossConnectMetricsMap) == 0 {
		logrus.Infof("failed to get cross-connect metrics map from event: %v", event)
	}

	// event Metrics contain single key-value of type
	// SRC/DTS + cross connect Id and metrics map
	for metricName, metrics := range event.Metrics {
		// Specifying cross connect by `crossConnectID`, parsed from `metricName`.
		// `communicationSide` can be 'SRC' or 'DST' in order to apply metrics
		// data to the cross connection source or destination respectively.
		ccID, communicationSide, err := parseMetricName(metricName)
		if err != nil {
			logrus.Warningf("error parsing metric name: %v", errors.Errorf("failed to get metrics: %v", err))
		}

		switch communicationSide {
		case "SRC":
			servePrometheus(crossConnectMetricsMap[ccID], metrics)
			continue
		case "DST":
			originalMetricsIdentifiers := map[string]string{}
			for k, v := range crossConnectMetricsMap[ccID] {
				originalMetricsIdentifiers[k] = v
			}

			// In cross connect source/destination represent client/endpoint.
			// When metrics are attached to event, they represent specific
			// cross connect by Id and client or endpoint by "SRC" or "DST" respectively.
			// This might be confusing and not clear when someone is observing metrics
			// for specific connection, as people would consider src/dst as src/dst of a
			// connection traffic. This is why we want to track the metrics from the user
			// perspective and when they are from "DST" to apply that as traffic source.
			crossConnectMetricsMap[ccID][metricspkg.SrcPodKey] = originalMetricsIdentifiers[metricspkg.DstPodKey]
			crossConnectMetricsMap[ccID][metricspkg.SrcNamespaceKey] = originalMetricsIdentifiers[metricspkg.DstNamespaceKey]
			crossConnectMetricsMap[ccID][metricspkg.DstPodKey] = originalMetricsIdentifiers[metricspkg.SrcPodKey]
			crossConnectMetricsMap[ccID][metricspkg.DstNamespaceKey] = originalMetricsIdentifiers[metricspkg.SrcNamespaceKey]
			servePrometheus(crossConnectMetricsMap[ccID], metrics)
		default:
			logrus.Warningf("error parsing metric: %v", errors.Errorf("error: communication side should be 'SRC' or 'DST', but got %s", communicationSide))
		}
	}
	logrus.Infof("no metrics available in event: %v", event)
}

func parseMetricName(metricName string) (string, string, error) {
	metricNameSlice := strings.Split(metricName, "-")
	if len(metricNameSlice) != 2 {
		return "", "", errors.Errorf("cannot parse metric to get key and crossconnect id. Inaproprite metric name received. Should be of type SRC-id or DST-id, but got: %s ", metricName)
	}
	if metricNameSlice[0] != "SRC" && metricNameSlice[0] != "DST" {
		return "", "", errors.Errorf("metric key should be SRC or DST, but got: %s", metricNameSlice[0])
	}
	// Returning crossConnect id and SRC/DST
	return metricNameSlice[1], metricNameSlice[0], nil
}

func lookForNSMServers() {
	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	// check if CRD is installed
	config, err := rest.InClusterConfig()
	if err != nil {
		logrus.Println("Unable to get in cluster config, attempting to fall back to kubeconfig", err)
		config, err = clientcmd.BuildConfigFromFlags("", *kubeconfig)
		if err != nil {
			logrus.Fatalln("Unable to build config", err)
		}
	}

	// Initialize clientset
	nsmClientSet, err := versioned.NewForConfig(config)
	if err != nil {
		logrus.Fatalln("Unable to initialize nsmd-k8s", err)
	}

	cache := map[string]string{}
	prom, err := tools.ReadEnvBool(metricspkg.PrometheusEnv, metricspkg.PrometheusDefault)
	if err != nil {
		logrus.Warningf("failed to serve prometheus: env PROMETHEUS=%t, err: %v", prom, err)
	} else if prom {
		cache["prometheus"] = "on"
	}

	nsmNamespace := namespace.GetNamespace()
	for !closing {
		result, err := nsmClientSet.NetworkserviceV1alpha1().NetworkServiceManagers(nsmNamespace).List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			logrus.Fatalln("Unable to find NSMs", err)
		}
		for i := range result.Items {
			mgr := &result.Items[i]
			if _, ok := managers[mgr.Spec.URL]; !ok {
				logrus.Printf("Adding manager: %s at %s", mgr.Name, mgr.Spec.URL)
				managers[mgr.Spec.URL] = "true"
				go monitorCrossConnects(mgr.Spec.URL, true, cache)
			}
		}
		time.Sleep(time.Second)
	}
}
