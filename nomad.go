package nomad

import (
	"context"
	"fmt"

	"github.com/coredns/coredns/plugin"
	clog "github.com/coredns/coredns/plugin/pkg/log"
	"github.com/coredns/coredns/request"
	nomad "github.com/hashicorp/nomad/api"

	"github.com/miekg/dns"
)

const pluginName = "nomad"

var log = clog.NewWithPlugin(pluginName)

// TODO: Figure out what to do with Zones
// Nomad is a plugin that serves records for Nomad services
type Nomad struct {
	Zones       []string
	NomadClient *nomad.Client
	Next        plugin.Handler
}

// ServeDNS implements the plugin.Handler interface. This method gets called when example is used
// in a Server.
func (n Nomad) ServeDNS(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {
	state := request.Request{W: w, Req: r}
	qname := state.Name()

	fmt.Println("qname:")
	fmt.Println(qname)

	// TODO A: GET A STUBBED OUT REPLY WORKING WITH FAKE DATA PASSING 1ST TEST

	// TODO B: GET NOMAD RUNNING WITH THE A JOB TO MATCH

	// TODO C: MAKE THIS MAKE A REQUEST TO NOMAD FOR THE SD INFO

	// TODO D: PASS THE SD INFO INTO THE RESPONSE & UPDATE THE TEST

	m := new(dns.Msg)
	m.SetReply(r)

	w.WriteMsg(m)
	return dns.RcodeSuccess, nil
}

// Name implements the Handler interface.
func (n Nomad) Name() string { return "nomad" }
