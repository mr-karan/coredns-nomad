package nomad

import (
	"bytes"
	"context"
	golog "log"
	"strings"
	"testing"

	"github.com/coredns/coredns/plugin/pkg/dnstest"
	"github.com/coredns/coredns/plugin/test"
	"github.com/miekg/dns"
)

func TestNomad(t *testing.T) {
	// Create a new Example Plugin. Use the test.ErrorHandler as the next plugin.
	n := Nomad{Next: test.ErrorHandler()}

	// Setup a new output buffer that is *not* standard output, so we can check if
	// example is really being printed.
	b := &bytes.Buffer{}
	golog.SetOutput(b)

	ctx := context.TODO()
	r := new(dns.Msg)
	r.SetQuestion("jobby-groupy-tasky.nomad", dns.TypeA)
	// Create a new Recorder that captures the result, this isn't actually used in this test
	// as it just serves as something that implements the dns.ResponseWriter interface.
	rec := dnstest.NewRecorder(&test.ResponseWriter{})

	// Call our plugin directly, and check the result.
	n.ServeDNS(ctx, rec, r)
	if a := b.String(); !strings.Contains(a, "[INFO] plugin/nomad: nomad") {
		t.Errorf("Failed to print '%s', got %s", "[INFO] plugin/nomad: nomad", a)
	}
}

func TestServeDNS(t *testing.T) {
	n := Nomad{Next: test.ErrorHandler()}

	var cases = []test.Case{
		{
			Qname: "jobA.groupA.taskA.nomad.",
			Qtype: dns.TypeA,
			Rcode: dns.RcodeSuccess,
			Answer: []dns.RR{
				test.A("jobA.groupA.taskA.nomad.	5	IN	A	1.2.3.4"),
			},
		},
		// {
		// 	Qname: "jobA.groupA.taskA.nomad.",
		// 	Qtype: dns.TypeAAAA,
		// 	Rcode: dns.RcodeSuccess,
		// 	Answer: []dns.RR{
		// 		test.AAAA("jobA.groupA.taskA.nomad.	5	IN	AAAA	1:2:3::4"),
		// 	},
		// },
		// {
		// 	Qname: "jobB.groupB.taskB.nomad.", Qtype: dns.TypeA,
		// 	Rcode: dns.RcodeSuccess,
		// 	Answer: []dns.RR{
		// 		test.A("jobB.groupB.taskB.nomad.	5	IN	A	1.2.3.5"),
		// 		test.A("jobB.groupB.taskB.nomad.	5	IN	A	1.2.3.6"),
		// 	},
		// },
		// {
		// 	Qname: "example.", Qtype: dns.TypeA,
		// 	Rcode: dns.RcodeSuccess,
		// 	Ns:    []dns.RR{k.soa()},
		// },
		// {
		// 	Qname: "nonexistent-node.example.", Qtype: dns.TypeA,
		// 	Rcode: dns.RcodeNameError,
		// 	Ns:    []dns.RR{k.soa()},
		// },
	}

	// n.client = fake.NewSimpleClientset()
	ctx := context.Background()

	// TODO: SET UP NOMAD OR FAKE NOMAD SERVER

	runTests(t, ctx, &n, cases)
}

func runTests(t *testing.T, ctx context.Context, n *Nomad, cases []test.Case) {
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
