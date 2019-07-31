package fanout

import (
	"context"
	"time"

	"github.com/coredns/coredns/plugin/dnstap"
	"github.com/coredns/coredns/plugin/dnstap/msg"
	"github.com/coredns/coredns/request"

	tap "github.com/dnstap/golang-dnstap"
	"github.com/miekg/dns"
)

func checkErr(err error) {
	if err == nil {
		return
	}
	log.Error(err)
}

func toDnstap(ctx context.Context, host string, req request.Request, reply *dns.Msg, start time.Time) error {
	tapper := dnstap.TapperFromContext(ctx)
	if tapper == nil {
		return nil
	}
	// Query
	b := msg.New().Time(start).HostPort(host)
	b.SocketProto = tap.SocketProtocol_TCP

	if tapper.Pack() {
		b.Msg(req.Req)
	}
	m, err := b.ToOutsideQuery(tap.Message_FORWARDER_QUERY)
	if err != nil {
		return err
	}
	tapper.TapMessage(m)

	if reply != nil {
		if tapper.Pack() {
			b.Msg(reply)
		}
		m, err := b.Time(time.Now()).ToOutsideResponse(tap.Message_FORWARDER_RESPONSE)
		if err != nil {
			return err
		}
		tapper.TapMessage(m)
	}

	return nil
}
