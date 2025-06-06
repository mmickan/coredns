package kubernetes

import (
	"testing"

	"github.com/coredns/coredns/request"

	"github.com/miekg/dns"
)

func TestParseRequest(t *testing.T) {
	tests := []struct {
		query        string
		expected     string // output from r.String()
		multicluster bool
	}{
		// valid SRV request
		{"_http._tcp.webs.mynamespace.svc.inter.webs.tests.", "http.tcp...webs.mynamespace.svc", false},
		// A request of endpoint
		{"1-2-3-4.webs.mynamespace.svc.inter.webs.tests.", "..1-2-3-4..webs.mynamespace.svc", false},
		// bare zone
		{"inter.webs.tests.", "......", false},
		// bare svc type
		{"svc.inter.webs.tests.", "......", false},
		// bare pod type
		{"pod.inter.webs.tests.", "......", false},
		// SRV request with empty segments
		{"..webs.mynamespace.svc.inter.webs.tests.", "....webs.mynamespace.svc", false},
		// A multicluster request with a clusterid
		{"1-2-3-4.cluster1.webs.mynamespace.svc.inter.webs.tests.", "..1-2-3-4.cluster1.webs.mynamespace.svc", true},
	}
	for i, tc := range tests {
		m := new(dns.Msg)
		m.SetQuestion(tc.query, dns.TypeA)
		state := request.Request{Zone: zone, Req: m}

		r, e := parseRequest(state.Name(), state.Zone, tc.multicluster)
		if e != nil {
			t.Errorf("Test %d, expected no error, got '%v'.", i, e)
		}
		rs := r.String()
		if rs != tc.expected {
			t.Errorf("Test %d, expected (stringified) recordRequest: %s, got %s", i, tc.expected, rs)
		}
	}
}

func TestParseInvalidRequest(t *testing.T) {
	invalid := []string{
		"webs.mynamespace.pood.inter.webs.test.",                 // Request must be for pod or svc subdomain.
		"too.long.for.what.I.am.trying.to.pod.inter.webs.tests.", // Too long.
	}

	for i, query := range invalid {
		m := new(dns.Msg)
		m.SetQuestion(query, dns.TypeA)
		state := request.Request{Zone: zone, Req: m}

		if _, e := parseRequest(state.Name(), state.Zone, false); e == nil {
			t.Errorf("Test %d: expected error from %s, got none", i, query)
		}
	}
}

const zone = "inter.webs.tests."
