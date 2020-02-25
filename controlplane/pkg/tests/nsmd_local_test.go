package tests

import (
	"context"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/networkservicemesh/api/pkg/api/networkservice/mechanisms/common"

	"github.com/golang/protobuf/ptypes/empty"
	. "github.com/onsi/gomega"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	"github.com/networkservicemesh/api/pkg/api/networkservice"
)

type nseWithOptions struct {
	netns             string
	srcIp             string
	dstIp             string
	need_ip_neighbors bool
	connection        *networkservice.Connection
}

func (impl *nseWithOptions) Request(ctx context.Context, in *networkservice.NetworkServiceRequest, opts ...grpc.CallOption) (*networkservice.Connection, error) {
	var mechanism *networkservice.Mechanism

	if in.Connection.Labels != nil {
		if val, ok := in.Connection.Labels["nse_sleep"]; ok {
			delay, err := strconv.Atoi(val)
			if err == nil {
				logrus.Infof("Delaying NSE init: %v", delay)
				<-time.After(time.Duration(delay) * time.Second)
			}
		}
	}
	mechanism = &networkservice.Mechanism{
		Type: in.MechanismPreferences[0].Type,
		Parameters: map[string]string{
			common.NetNSInodeKey: impl.netns,
			// TODO: Fix this terrible hack using xid for getting a unique interface name
			common.InterfaceNameKey: "nsm" + in.GetConnection().GetId(),
		},
	}

	conn := &networkservice.Connection{
		Id:             in.GetConnection().GetId(),
		NetworkService: in.GetConnection().GetNetworkService(),
		Mechanism:      mechanism,
		Context: &networkservice.ConnectionContext{
			IpContext: &networkservice.IPContext{
				SrcIpAddr: impl.srcIp,
				DstIpAddr: impl.dstIp,
			},
		},
	}

	if impl.need_ip_neighbors {
		conn.GetContext().GetIpContext().IpNeighbors = []*networkservice.IpNeighbor{
			{
				Ip:              "127.0.0.1",
				HardwareAddress: "ff-ee-ff-ee-ff",
			},
		}
	}
	impl.connection = conn
	return conn, nil
}

func (nseWithOptions) Close(ctx context.Context, in *networkservice.Connection, opts ...grpc.CallOption) (*empty.Empty, error) {
	return nil, nil
}

// Below only tests

func TestNSMDRequestClientConnectionRequest(t *testing.T) {
	g := NewWithT(t)

	storage := NewSharedStorage()
	srv := NewNSMDFullServer(Master, storage)
	defer srv.Stop()
	srv.AddFakeForwarder("test_data_plane", "tcp:some_addr")

	srv.TestModel.AddEndpoint(context.Background(), srv.RegisterFakeEndpoint("golden_network", "test", Master))

	nsmClient, conn := srv.requestNSMConnection("nsm")
	defer conn.Close()

	request := CreateRequest()

	nsmResponse, err := nsmClient.Request(context.Background(), request)
	g.Expect(err).To(BeNil())
	g.Expect(nsmResponse.GetNetworkService()).To(Equal("golden_network"))
	logrus.Print("End of test")
}

func TestNSENoSrc(t *testing.T) {
	g := NewWithT(t)

	storage := NewSharedStorage()
	srv := NewNSMDFullServer(Master, storage)
	defer srv.Stop()

	srv.serviceRegistry.localTestNSE = &nseWithOptions{
		netns: "12",
		//srcIp: "169083138/30",
		dstIp: "10.20.1.2/30",
	}
	srv.AddFakeForwarder("test_data_plane", "tcp:some_addr")

	srv.TestModel.AddEndpoint(context.Background(), srv.RegisterFakeEndpoint("golden_network", "test", Master))

	nsmClient, conn := srv.requestNSMConnection("nsm")
	defer conn.Close()

	request := CreateRequest()

	nsmResponse, err := nsmClient.Request(context.Background(), request)
	println(err.Error())
	g.Expect(strings.Contains(err.Error(), "failure Validating NSE Connection: ConnectionContext.SrcIp is required cannot be empty/nil")).To(Equal(true))
	g.Expect(nsmResponse).To(BeNil())
}

func TestNSEIPNeghtbours(t *testing.T) {
	g := NewWithT(t)

	storage := NewSharedStorage()
	srv := NewNSMDFullServer(Master, storage)
	defer srv.Stop()
	srv.serviceRegistry.localTestNSE = &nseWithOptions{
		netns:             "12",
		srcIp:             "10.20.1.1/30",
		dstIp:             "10.20.1.2/30",
		need_ip_neighbors: true,
	}

	srv.AddFakeForwarder("test_data_plane", "tcp:some_addr")
	srv.TestModel.AddEndpoint(context.Background(), srv.RegisterFakeEndpoint("golden_network", "test", Master))

	nsmClient, conn := srv.requestNSMConnection("nsm")
	defer conn.Close()

	request := CreateRequest()

	nsmResponse, err := nsmClient.Request(context.Background(), request)
	g.Expect(err).To(BeNil())
	g.Expect(nsmResponse.GetNetworkService()).To(Equal("golden_network"))
	logrus.Print("End of test")

	originl, ok := srv.serviceRegistry.localTestNSE.(*nseWithOptions)
	g.Expect(ok).To(Equal(true))

	g.Expect(len(originl.connection.GetContext().GetIpContext().GetIpNeighbors())).To(Equal(1))
	g.Expect(originl.connection.GetContext().GetIpContext().GetIpNeighbors()[0].Ip).To(Equal("127.0.0.1"))
	g.Expect(originl.connection.GetContext().GetIpContext().GetIpNeighbors()[0].HardwareAddress).To(Equal("ff-ee-ff-ee-ff"))
}

func TestSlowNSE(t *testing.T) {
	g := NewWithT(t)

	storage := NewSharedStorage()
	srv := NewNSMDFullServer(Master, storage)
	defer srv.Stop()

	srv.serviceRegistry.localTestNSE = &nseWithOptions{
		netns: "12",
		srcIp: "169083138/30",
		dstIp: "169083137/30",
	}
	srv.AddFakeForwarder("test_data_plane", "tcp:some_addr")

	srv.TestModel.AddEndpoint(context.Background(), srv.RegisterFakeEndpoint("golden_network", "test", Master))

	nsmClient, conn := srv.requestNSMConnection("nsm")
	defer conn.Close()

	request := CreateRequest()

	request.Connection.Labels = map[string]string{}
	request.Connection.Labels["nse_sleep"] = "1"

	ctx, canceOp := context.WithTimeout(context.Background(), 30*time.Millisecond)
	defer canceOp()
	nsmResponse, err := nsmClient.Request(ctx, request)
	println(err.Error())
	g.Expect(strings.Contains(err.Error(), "rpc error: code = DeadlineExceeded desc = context deadline exceeded")).To(Equal(true))
	g.Expect(nsmResponse).To(BeNil())
}

func TestSlowDP(t *testing.T) {
	g := NewWithT(t)

	storage := NewSharedStorage()
	srv := NewNSMDFullServer(Master, storage)
	defer srv.Stop()

	srv.serviceRegistry.localTestNSE = &nseWithOptions{
		netns: "12",
		srcIp: "10.20.1.1/30",
		dstIp: "10.20.1.2/30",
	}
	srv.AddFakeForwarder("test_data_plane", "tcp:some_addr")

	srv.TestModel.AddEndpoint(context.Background(), srv.RegisterFakeEndpoint("golden_network", "test", Master))

	nsmClient, conn := srv.requestNSMConnection("nsm")
	defer conn.Close()

	request := CreateRequest()

	request.Connection.Labels = map[string]string{}
	request.Connection.Labels["forwarder_sleep"] = "1"

	ctx, cancelOp := context.WithTimeout(context.Background(), 30*time.Millisecond)
	defer cancelOp()
	nsmResponse, err := nsmClient.Request(ctx, request)
	println(err.Error())
	g.Expect(strings.Contains(err.Error(), "rpc error: code = DeadlineExceeded desc = context deadline exceeded")).To(Equal(true))
	g.Expect(nsmResponse).To(BeNil())
}
