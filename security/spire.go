package security

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"github.com/sirupsen/logrus"
	"github.com/spiffe/spire/api/workload"
	proto "github.com/spiffe/spire/proto/spire/api/workload"
	"net"
	"time"
)

type spireCertObtainer struct {
	stopCh            chan struct{}
	errorCh           chan error
	workloadAPIClient workload.X509Client
}

func NewSpireCertObtainer(agentAddress string, timeout time.Duration) CertificateObtainer {
	return &spireCertObtainer{
		stopCh:            make(chan struct{}),
		errorCh:           make(chan error),
		workloadAPIClient: newWorkloadAPIClient(agentAddress, timeout),
	}
}

func newWorkloadAPIClient(agentAddress string, timeout time.Duration) workload.X509Client {
	addr := &net.UnixAddr{
		Net:  "unix",
		Name: agentAddress,
	}
	config := &workload.X509ClientConfig{
		Addr:    addr,
		Timeout: timeout,
	}
	return workload.NewX509Client(config)
}

func (s *spireCertObtainer) ObtainCertificates() <-chan *RetrievedCerts {
	certCh := make(chan *RetrievedCerts)

	go func() {
		if err := s.workloadAPIClient.Start(); err != nil {
			logrus.Error(err.Error())
			s.errorCh <- err
			close(certCh)
			return
		}
	}()
	defer s.workloadAPIClient.Stop()

	go func() {
		defer close(certCh)

		updateCh := s.workloadAPIClient.UpdateChan()
		for {
			select {
			case svidResponse := <-updateCh:
				logrus.Infof("Received new SVID: %v", svidResponse.Svids[0].SpiffeId)
				if c, err := readCertificates(svidResponse); err == nil {
					certCh <- c
				} else {
					s.errorCh <- err
					return
				}
			case <-s.stopCh:
				return
			}
		}
	}()

	return certCh
}

func (s *spireCertObtainer) Stop() {
	close(s.stopCh)
}

func (s *spireCertObtainer) Error() error {
	return <-s.errorCh
}

func readCertificates(svidResponse *proto.X509SVIDResponse) (*RetrievedCerts, error) {
	svid := svidResponse.Svids[0]
	keyPair, err := tls.X509KeyPair(svid.GetX509Svid(), svid.GetX509SvidKey())
	if err != nil {
		return nil, err
	}

	caPool := x509.NewCertPool()
	if ok := caPool.AppendCertsFromPEM(svid.GetBundle()); !ok {
		return nil, errors.New("failed to append ca cert to pool")
	}

	return &RetrievedCerts{
		TLSCert:  &keyPair,
		CABundle: caPool,
	}, nil
}