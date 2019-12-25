// +build unit_test

package fanout

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/caddyserver/caddy"
	"github.com/coredns/coredns/plugin/pkg/dnstest"
	"github.com/coredns/coredns/plugin/test"
	"github.com/miekg/dns"
)

func TestFanoutCanReturnUnsuccessRespnse(t *testing.T) {
	s := newServer(func(w dns.ResponseWriter, r *dns.Msg) {
		msg := nxdomainMsg()
		msg.SetRcode(r, msg.Rcode)
		w.WriteMsg(msg)
	})
	f := NewFanout()
	c := createFanoutClient(s.Addr)
	f.addClient(c)
	defer f.Close()
	req := new(dns.Msg)
	req.SetQuestion("example1.", dns.TypeA)
	writer := &cachedDNSWriter{ResponseWriter: new(test.ResponseWriter)}
	f.ServeDNS(context.TODO(), writer, req)
	if len(writer.answers) != 1 {
		fmt.Println(len(writer.answers))
		t.FailNow()
	}
	if writer.answers[0].MsgHdr.Rcode != dns.RcodeNameError {
		t.Error("fanout plugin returns first negative answer if other answers on request are negative")
	}
}
func TestFanoutTwoServersNotSuccessResponse(t *testing.T) {
	rcode := 1
	rcodeMutex := sync.Mutex{}
	s1 := newServer(func(w dns.ResponseWriter, r *dns.Msg) {
		if r.Question[0].Name == "example1." {
			msg := nxdomainMsg()
			rcodeMutex.Lock()
			msg.SetRcode(r, rcode)
			rcode++
			rcode %= dns.RcodeNotZone
			rcodeMutex.Unlock()
			w.WriteMsg(msg)
			//let another server answer
			<-time.After(time.Millisecond * 100)
		}
	})
	s2 := newServer(func(w dns.ResponseWriter, r *dns.Msg) {
		if r.Question[0].Name == "example1." {
			msg := dns.Msg{
				Answer: []dns.RR{makeRecordA("example1. 3600	IN	A 10.0.0.1")},
			}
			msg.SetReply(r)
			w.WriteMsg(&msg)
			//let another server answer
			<-time.After(time.Millisecond * 100)
		}
	})
	defer s1.close()
	defer s2.close()
	c1 := createFanoutClient(s1.Addr)
	c2 := createFanoutClient(s2.Addr)
	f := NewFanout()
	f.addClient(c1)
	f.addClient(c2)
	defer f.Close()
	writer := &cachedDNSWriter{ResponseWriter: new(test.ResponseWriter)}
	for i := 0; i < 10; i++ {
		req := new(dns.Msg)
		req.SetQuestion("example1.", dns.TypeA)
		f.ServeDNS(context.TODO(), writer, req)
	}
	for _, m := range writer.answers {
		if m.MsgHdr.Rcode != dns.RcodeSuccess {
			t.Error("fanout should return only positive answers")
		}
	}
}

func TestFanoutTwoServers(t *testing.T) {
	const expected = 1
	var mutex sync.Mutex
	answerCount1 := 0
	answerCount2 := 0
	s1 := newServer(func(w dns.ResponseWriter, r *dns.Msg) {
		if r.Question[0].Name == "example1." {
			msg := dns.Msg{
				Answer: []dns.RR{makeRecordA("example1 3600	IN	A 10.0.0.1")},
			}
			mutex.Lock()
			answerCount1++
			mutex.Unlock()
			msg.SetReply(r)
			w.WriteMsg(&msg)
		}
	})
	s2 := newServer(func(w dns.ResponseWriter, r *dns.Msg) {
		if r.Question[0].Name == "example2." {
			msg := dns.Msg{
				Answer: []dns.RR{makeRecordA("example2. 3600	IN	A 10.0.0.1")},
			}
			mutex.Lock()
			answerCount2++
			mutex.Unlock()
			msg.SetReply(r)
			w.WriteMsg(&msg)
		}
	})
	defer s1.close()
	defer s2.close()

	c1 := createFanoutClient(s1.Addr)
	c2 := createFanoutClient(s2.Addr)
	f := NewFanout()
	f.addClient(c1)
	f.addClient(c2)
	defer f.Close()

	req := new(dns.Msg)
	req.SetQuestion("example1.", dns.TypeA)
	f.ServeDNS(context.TODO(), &test.ResponseWriter{}, req)
	<-time.After(time.Second)
	req = new(dns.Msg)
	req.SetQuestion("example2.", dns.TypeA)
	f.ServeDNS(context.TODO(), &test.ResponseWriter{}, req)
	mutex.Lock()
	defer mutex.Unlock()
	if answerCount2 != expected || answerCount1 != expected {
		t.Errorf("Expected number of health checks to be %d, got s1: %d, s2: %d", expected, answerCount1, answerCount2)
	}
}

func TestFanout(t *testing.T) {
	s := newServer(func(w dns.ResponseWriter, r *dns.Msg) {
		ret := new(dns.Msg)
		ret.SetReply(r)
		ret.Answer = append(ret.Answer, test.A("example.org. IN A 127.0.0.1"))
		w.WriteMsg(ret)
	})
	defer s.close()
	c := caddy.NewTestController("dns", "fanout "+s.Addr)
	f, err := parseFanout(c)
	if err != nil {
		t.Errorf("Failed to create fanout: %s", err)
	}
	f.OnStartup()
	defer f.OnShutdown()

	m := new(dns.Msg)
	m.SetQuestion("example.org.", dns.TypeA)
	rec := dnstest.NewRecorder(&test.ResponseWriter{})

	if _, err := f.ServeDNS(context.TODO(), rec, m); err != nil {
		t.Fatal("Expected to receive reply, but didn't")
	}
	if x := rec.Msg.Answer[0].Header().Name; x != "example.org." {
		t.Errorf("Expected %s, got %s", "example.org.", x)
	}
}
