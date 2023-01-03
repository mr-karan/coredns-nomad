package nomad

import (
	"testing"

	"github.com/coredns/caddy"
)

// TestSetup tests the various things that should be parsed by setup.
// Make sure you also test for parse errors.
func TestSetup(t *testing.T) {
	c := caddy.NewTestController("dns", `nomad {
		address http://127.0.0.1:4646
		token 4649b287-1213-6080-8b77-f115f5b4e8e0
		tls-insecure
}`)
	if err := setup(c); err != nil {
		t.Fatalf("Expected no errors, but got: %v", err)
	}

	c = caddy.NewTestController("dns", `nomad {
		address http://127.0.0.1:4646
		token 4649b287-1213-6080-8b77-f115f5b4e8e0
		tls-insecure
		foo bar
}`)
	if err := setup(c); err == nil {
		t.Fatalf("Expected errors, but got: %v", err)
	}
}
