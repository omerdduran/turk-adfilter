package pipeline

import (
	"context"
	"errors"
	"net"
	"testing"
	"time"
)

type mockLookuper struct {
	ips []net.IP
	err error
}

func (m mockLookuper) lookupIP(context.Context, string) ([]net.IP, error) {
	return m.ips, m.err
}

func ip(s string) net.IP { return net.ParseIP(s) }

func TestQuorumActive(t *testing.T) {
	q := &QuorumResolver{
		resolvers: []lookuper{
			mockLookuper{ips: []net.IP{ip("104.21.1.1")}},
			mockLookuper{ips: []net.IP{ip("172.67.1.1")}},
		},
		need: 2,
	}
	r := q.Resolve(context.Background(), "x.com")
	if !r.Active || r.Class != ClassReal {
		t.Fatalf("iki real IP → active/real bekleniyordu: %+v", r)
	}
}

func TestQuorumNotMet(t *testing.T) {
	q := &QuorumResolver{
		resolvers: []lookuper{
			mockLookuper{ips: []net.IP{ip("104.21.1.1")}},
			mockLookuper{err: errors.New("nxdomain")},
		},
		need: 2,
	}
	if r := q.Resolve(context.Background(), "x.com"); r.Active {
		t.Errorf("quorum(2) sağlanmadı → inactive bekleniyordu: %+v", r)
	}
}

func TestSinkholeInactive(t *testing.T) {
	q := &QuorumResolver{
		resolvers: []lookuper{
			mockLookuper{ips: []net.IP{ip("0.0.0.0")}},
			mockLookuper{ips: []net.IP{ip("0.0.0.0")}},
		},
		need: 2,
	}
	r := q.Resolve(context.Background(), "x.com")
	if r.Active || r.Class != ClassSinkhole {
		t.Errorf("0.0.0.0 → inactive/sinkhole bekleniyordu: %+v", r)
	}
}

func TestClassifyIP(t *testing.T) {
	if classifyIP(ip("0.0.0.0")) != ClassSinkhole {
		t.Error("0.0.0.0 sinkhole olmalı")
	}
	if classifyIP(ip("195.175.254.2")) != ClassBlock {
		t.Error("BTK IP block olmalı")
	}
	if classifyIP(ip("104.21.1.1")) != ClassReal {
		t.Error("normal IP real olmalı")
	}
	// SSRF koruması: dahili IP'ler asla real sayılmaz.
	for _, s := range []string{"192.168.1.1", "10.0.0.5", "172.16.0.1", "169.254.169.254", "100.64.0.1"} {
		if classifyIP(ip(s)) != ClassInternal {
			t.Errorf("%s dahili (internal) sayılmalı — SSRF koruması", s)
		}
	}
}

func TestNewQuorumMajority(t *testing.T) {
	// 3 resolver → 2-of-3 majority.
	q := NewQuorumResolver([]string{"1.1.1.1:53", "1.0.0.1:53", "8.8.8.8:53"}, time.Second)
	if q.need != 2 {
		t.Errorf("3 resolver için need=2 (majority) olmalı, %d alındı", q.need)
	}
	// 1 resolver → need 1.
	if NewQuorumResolver([]string{"1.1.1.1:53"}, time.Second).need != 1 {
		t.Error("tek resolver için need=1 olmalı")
	}
}

func TestQuorumInternalNotActive(t *testing.T) {
	// Aday dahili IP'ye çözülüyorsa aktif SAYILMAZ (SSRF).
	q := &QuorumResolver{
		resolvers: []lookuper{
			mockLookuper{ips: []net.IP{ip("169.254.169.254")}}, // bulut metadata
			mockLookuper{ips: []net.IP{ip("169.254.169.254")}},
		},
		need: 2,
	}
	if r := q.Resolve(context.Background(), "evil.com"); r.Active || r.Class != ClassInternal {
		t.Errorf("dahili IP → inactive/internal bekleniyordu: %+v", r)
	}
}

func TestGuardDial(t *testing.T) {
	// Probe dial-guard: dahili/loopback hedefleri reddeder.
	for _, addr := range []string{"127.0.0.1:443", "192.168.1.1:443", "169.254.169.254:80", "10.0.0.1:443"} {
		if guardDial("tcp", addr, nil) == nil {
			t.Errorf("guardDial(%q) reddetmeliydi", addr)
		}
	}
	if guardDial("tcp", "104.21.1.1:443", nil) != nil {
		t.Error("guardDial genel IP'yi kabul etmeli")
	}
}
