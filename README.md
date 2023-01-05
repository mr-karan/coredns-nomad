# nomad

## Name

*nomad* - DNS interface to Nomad native service discovery.

## Description

The nomad plugin serves DNS records for services registered with Nomad. Nomad 1.3+ comes with [support for discovering services](https://www.hashicorp.com/blog/nomad-service-discovery) with an in-built service catalogue that is available via the HTTP API. This plugin extends the HTTP API and provides a DNS interface for querying the service catalogue.

The query can be looked up with the format `service.namespace.nomad`. The plugin currently handles A, AAAA and SRV records. Refer to [#Usage Example](#usage-example) for more details.

## Compilation

NOTE: For now, this is a standalone package until [https://github.com/coredns/coredns/pull/5833](https://github.com/coredns/coredns/pull/5833) gets merged in upstream.

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

## Usage Example

### A record

```
dig redis.default.nomad @127.0.0.1 -p 1053    

; <<>> DiG 9.18.1-1ubuntu1.2-Ubuntu <<>> redis.default.nomad @127.0.0.1 -p 1053
;; global options: +cmd
;; Got answer:
;; ->>HEADER<<- opcode: QUERY, status: NOERROR, id: 54986
;; flags: qr aa rd; QUERY: 1, ANSWER: 3, AUTHORITY: 0, ADDITIONAL: 1
;; WARNING: recursion requested but not available

;; OPT PSEUDOSECTION:
; EDNS: version: 0, flags:; udp: 1232
; COOKIE: bdc9237f49a1f744 (echoed)
;; QUESTION SECTION:
;redis.default.nomad.		IN	A

;; ANSWER SECTION:
redis.default.nomad.	10	IN	A	192.168.29.76
redis.default.nomad.	10	IN	A	192.168.29.76
redis.default.nomad.	10	IN	A	192.168.29.76

;; Query time: 4 msec
;; SERVER: 127.0.0.1#1053(127.0.0.1) (UDP)
;; WHEN: Thu Jan 05 12:12:25 IST 2023
;; MSG SIZE  rcvd: 165
```

### SRV Record

Since an A record doesn't contain the port number, SRV record can be used to query the port number of a service.

```
dig redis.default.nomad @127.0.0.1 -p 1053 SRV

; <<>> DiG 9.18.1-1ubuntu1.2-Ubuntu <<>> redis.default.nomad @127.0.0.1 -p 1053 SRV
;; global options: +cmd
;; Got answer:
;; ->>HEADER<<- opcode: QUERY, status: NOERROR, id: 49945
;; flags: qr aa rd; QUERY: 1, ANSWER: 3, AUTHORITY: 0, ADDITIONAL: 4
;; WARNING: recursion requested but not available

;; OPT PSEUDOSECTION:
; EDNS: version: 0, flags:; udp: 1232
; COOKIE: 14572535f3ba6648 (echoed)
;; QUESTION SECTION:
;redis.default.nomad.		IN	SRV

;; ANSWER SECTION:
redis.default.nomad.	8	IN	SRV	10 10 25395 redis.default.nomad.
redis.default.nomad.	8	IN	SRV	10 10 20888 redis.default.nomad.
redis.default.nomad.	8	IN	SRV	10 10 26292 redis.default.nomad.

;; ADDITIONAL SECTION:
redis.default.nomad.	8	IN	A	192.168.29.76
redis.default.nomad.	8	IN	A	192.168.29.76
redis.default.nomad.	8	IN	A	192.168.29.76

;; Query time: 0 msec
;; SERVER: 127.0.0.1#1053(127.0.0.1) (UDP)
;; WHEN: Thu Jan 05 12:12:20 IST 2023
;; MSG SIZE  rcvd: 339
```

### plugin.cfg

This plugin is intended to appear twoard the end of the plugin list, usually near the `proxy` plugin declaration.
