package nomad

import (
	"strings"

	"github.com/coredns/caddy"
	"github.com/coredns/coredns/core/dnsserver"
	"github.com/coredns/coredns/plugin"
	nomad "github.com/hashicorp/nomad/api"
)

// init registers this plugin.
func init() { plugin.Register(pluginName, setup) }

func setup(c *caddy.Controller) error {
	n := Nomad{}
	err := parse(c, n)

	if err != nil {
		return plugin.Error("nomad", err)
	}

	dnsserver.GetConfig(c).AddPlugin(func(next plugin.Handler) plugin.Handler {
		n.Next = next
		return n
	})

	return nil
}

func parse(c *caddy.Controller, n Nomad) error {
	nomadConfig := nomad.DefaultConfig()
	nomadConfig.TLSConfig.Insecure = false

	for c.Next() {
		for c.NextBlock() {
			selector := strings.ToLower(c.Val())

			switch selector {
			case "address":
				nomadConfig.Address = c.RemainingArgs()[0]
			case "token":
				nomadConfig.SecretID = c.RemainingArgs()[0]
			case "tls-insecure":
				nomadConfig.TLSConfig.Insecure = true
			default:
				return c.Errf("unknown property '%s'", selector)
			}
		}
	}

	nomadClient, err := nomad.NewClient(nomadConfig)
	if err != nil {
		return plugin.Error("nomad", err)
	}
	n.NomadClient = nomadClient

	return nil
}
