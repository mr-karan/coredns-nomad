# nomad

## Name

*nomad* - DNS interface to Nomad native service discovery.

## Description

The nomad plugin serves DNS records for services registered with Nomad. Nomad 1.3+ comes with [support for discovering services](https://www.hashicorp.com/blog/nomad-service-discovery) with an in-built service catalogue that is available via the HTTP API. This plugin extends the HTTP API and provides a DNS interface for querying the service catalogue.

## Compilation

NOTE: For now, this is a standalone package until this PR gets merged in upstream.

It will require you to use `go get` or as a dependency on [plugin.cfg](https://github.com/coredns/coredns/blob/master/plugin.cfg).

The [manual](https://coredns.io/manual/toc/#what-is-coredns) will have more information about how to configure and extend the server with external plugins.

A simple way to consume this plugin, is by adding the following on [plugin.cfg](https://github.com/coredns/coredns/blob/master/plugin.cfg), and recompile it as [detailed on coredns.io](https://coredns.io/2017/07/25/compile-time-enabling-or-disabling-plugins/#build-with-compile-time-configuration-file).

~~~
nomad:github.com/mr-karan/coredns-nomad
~~~

You can put it after the `kubernetes` plugin in the above list.

After this you can compile coredns by running:

``` sh
make
```

## Syntax

~~~ txt
nomad {
    address URL
    token TOKEN
    ttl DURATION
}
~~~

* `address` The address where a Nomad agent (server) is available. **URL** defaults to `http://127.0.0.1:4646`.

* `token` The SecretID of an ACL token to use to authenticate API requests with if the Nomad cluster has ACL enabled. **TOKEN** defaults to `""`.

* `ttl` allows you to set a custom TTL for responses. **DURATION** defaults to `30 seconds`. The minimum TTL allowed is `0` seconds, and the maximum is capped at `3600` seconds. Setting TTL to 0 will prevent records from being cached. The unit for the value is seconds.

## Metrics

If monitoring is enabled (via the *prometheus* directive) the following metric is exported:

* `coredns_nomad_success_requests_total{namespace,server}` - Counter of DNS requests handled successfully.
* `coredns_nomad_failed_requests_total{namespace,server}` - Counter of DNS requests failed.

The `server` label indicated which server handled the request. `namespace` indicates the namespace of the service in the query.

## Ready

This plugin reports readiness to the ready plugin. It will be ready only when it has successfully connected to the Nomad server. It queries the [`/v1/agent/self`](https://developer.hashicorp.com/nomad/api-docs/agent#query-self) endpoint to check if it is ready.

## Examples

Enable nomad with and resolve all services with `.nomad` as the suffix. `cache` plugin is used to cache the responses for 30 seconds. This avoids a lookup to the Nomad server for every request.

```
nomad:1053 {
    nomad {
	  	address http://127.0.0.1:4646
    }
    cache 30
}
```

You can see the [Corefile.example](./Corefile.example) for a full Corefile example.

## Authentication

`nomad` plugin uses a default Nomad configuration to create an API client. Options like the HTTP address and the token can be specified in Corefile. However, Nomad Go SDK can also additionally read these environment variables.

- `NOMAD_TOKEN`
- `NOMAD_ADDR`
- `NOMAD_REGION`
- `NOMAD_NAMESPACE`
- `NOMAD_HTTP_AUTH`
- `NOMAD_CACERT`
- `NOMAD_CAPATH`
- `NOMAD_CLIENT_CERT`
- `NOMAD_CLIENT_KEY`
- `NOMAD_TLS_SERVER_NAME`
- `NOMAD_SKIP_VERIFY`

You can read about them in detail [here](https://www.nomadproject.io/docs/runtime/environment).

### plugin.cfg

This plugin is intended to appear twoard the end of the plugin list, usually near the `proxy` plugin declaration.
