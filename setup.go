package nomad

import (
	"strconv"
	"strings"

	"github.com/coredns/caddy"
	"github.com/coredns/coredns/core/dnsserver"
	"github.com/coredns/coredns/plugin"
	nomad "github.com/hashicorp/nomad/api"
)

// init registers this plugin.
func init() { plugin.Register(pluginName, setup) }

// setup is the function that gets called when the config parser see the token "nomad". Setup is responsible
// for parsing any extra options the nomad plugin may have. The first token this function sees is "nomad".
func setup(c *caddy.Controller) error {
	n := &Nomad{
		ttl: uint32(defaultTTL),
	}
	if err := parse(c, n); err != nil {
		return plugin.Error("nomad", err)
	}

	// Do a ping check to check if the Nomad server is reachable.
	_, err := n.client.Agent().Self()
	if err != nil {
		return plugin.Error("nomad", err)
	}
	// Mark the plugin as ready to use.
	// https://github.com/coredns/coredns/blob/master/plugin.md#readiness
	n.Ready()

	dnsserver.GetConfig(c).AddPlugin(func(next plugin.Handler) plugin.Handler {
		n.Next = next
		return n
	})

	return nil
}

func parse(c *caddy.Controller, n *Nomad) error {
	cfg := nomad.DefaultConfig()

	for c.Next() {
		for c.NextBlock() {
			selector := strings.ToLower(c.Val())

			switch selector {
			case "address":
				cfg.Address = c.RemainingArgs()[0]
			case "token":
				cfg.SecretID = c.RemainingArgs()[0]
			case "ttl":
				t, err := strconv.Atoi(c.RemainingArgs()[0])
				if err != nil {
					return c.Err("error parsing ttl: " + err.Error())
				}
				if t < 0 || t > 3600 {
					return c.Errf("ttl must be in range [0, 3600]: %d", t)
				}
				n.ttl = uint32(t)
			default:
				return c.Errf("unknown property '%s'", selector)
			}
		}
	}

	// Create a new Nomad client.
	nomadClient, err := nomad.NewClient(cfg)
	if err != nil {
		return plugin.Error("nomad", err)
	}
	n.client = nomadClient

	return nil
}
