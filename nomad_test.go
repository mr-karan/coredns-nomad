package nomad

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/coredns/coredns/plugin/pkg/dnstest"
	"github.com/coredns/coredns/plugin/test"
	nomad "github.com/hashicorp/nomad/api"
	"github.com/miekg/dns"
)

func TestNomad(t *testing.T) {
	var demoNomad = Nomad{
		Next: test.ErrorHandler(),
		ttl:  uint32(defaultTTL),
	}

	var cases = []test.Case{
		{
			Qname: "example.default.nomad.",
			Qtype: dns.TypeA,
			Rcode: dns.RcodeSuccess,
			Answer: []dns.RR{
				test.A("example.default.nomad.	30	IN	A	1.2.3.4"),
			},
		},
		{
			Qname: "fakeipv6.default.nomad.",
			Qtype: dns.TypeAAAA,
			Rcode: dns.RcodeSuccess,
			Answer: []dns.RR{
				test.AAAA("fakeipv6.default.nomad.	30	IN	AAAA	1:2:3::4"),
			},
		},
		{
			Qname: "multi.default.nomad.",
			Qtype: dns.TypeA,
			Rcode: dns.RcodeSuccess,
			Answer: []dns.RR{
				test.A("multi.default.nomad.	30	IN	A	1.2.3.4"),
				test.A("multi.default.nomad.	30	IN	A	1.2.3.5"),
				test.A("multi.default.nomad.	30	IN	A	1.2.3.6"),
			},
		},
		{
			Qname:  "nonexistent.default.nomad.",
			Qtype:  dns.TypeA,
			Rcode:  dns.RcodeNameError,
			Answer: []dns.RR{},
		},
		{
			Qname:  "example.default.nomad.",
			Qtype:  dns.TypeSRV,
			Rcode:  dns.RcodeSuccess,
			Answer: []dns.RR{test.SRV("example.default.nomad.	30	IN	SRV	10 10 23202 example.default.nomad.")},
			Extra:  []dns.RR{test.A("example.default.nomad.  30       IN      A       1.2.3.4")},
		},
	}

	ctx := context.Background()

	// Setup a fake Nomad server.
	nomadServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		default:
			t.Errorf("Not implemented: %v", r.URL.Path)
			return
		case "/v1/service/example":
			w.Write([]byte(`[{"Address":"1.2.3.4","Namespace":"default","Port":23202,"ServiceName":"example"}]`))
		case "/v1/service/fakeipv6":
			w.Write([]byte(`[{"Address":"1:2:3::4","Namespace":"default","Port":8000,"ServiceName":"fakeipv6"}]`))
		case "/v1/service/multi":
			w.Write([]byte(`[{"Address":"1.2.3.4","Namespace":"default","Port":25395,"ServiceName":"multi"},{"Address":"1.2.3.5","Namespace":"default","Port":20888,"ServiceName":"multi"},{"Address":"1.2.3.6","Namespace":"default","Port":26292,"ServiceName":"multi"}]`))
		case "/v1/service/nonexistent":
			w.Write([]byte(`[]`))
		}
	}))
	defer nomadServer.Close()

	// Configure the plugin to use the fake Nomad server.
	cfg := nomad.DefaultConfig()
	cfg.Address = nomadServer.URL
	client, err := nomad.NewClient(cfg)
	if err != nil {
		t.Errorf("Failed to create Nomad client: %v", err)
		return
	}
	demoNomad.client = client

	runTests(ctx, t, &demoNomad, cases)
}

func runTests(ctx context.Context, t *testing.T, n *Nomad, cases []test.Case) {
	for i, tc := range cases {
		r := tc.Msg()
		w := dnstest.NewRecorder(&test.ResponseWriter{})

		_, err := n.ServeDNS(ctx, w, r)
		if err != tc.Error {
			t.Errorf("Test %d: %v", i, err)
			return
		}

		if w.Msg == nil {
			t.Errorf("Test %d: nil message", i)
		}
		if err := test.SortAndCheck(w.Msg, tc); err != nil {
			t.Errorf("Test %d: %v", i, err)
		}
	}
}
