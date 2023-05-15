Zitified Caddy Server
----

This project [zitifies](https://docs.openziti.io/docs/reference/glossary#zitification-zitified-zitify)
[Caddy Server](https://caddyserver.com/). It allows caddy server to be either 
a native Ziti Service provider or a front-end proxy to a ziti service.

This project is implemented as Caddy service module. It provides two components:
- ziti listener: allows to bind to a ziti service
- ziti transport: allows to use ziti service as a reverse proxy backend

Configuration of these features is driven by [Caddyfile](https://caddyserver.com/docs/caddyfile-tutorial)

## Build Zitified Caddy server
```
$ git clone https://github.com/openziti-test-kitchen/ziti-caddy
$ cd ziti-caddy
$ go build ./...
```
After the build is done `ziti-caddy` executable should be present in the directory. It is a drop-in replacement for
the standard Caddy executable.

## Running Zitified Caddy 

#### Service Side
To configure caddy to accept connections over ziti network you need to provide a ziti binding address
in the server configuration:
```
	# ziti address format: ziti:/<service_name>[/<terminator>]@<ziti_identity_file>
	bind ziti/caddy-http@caddy-server.json
```

To try it you'll need to create a [ziti service](https://docs.openziti.io/docs/learn/core-concepts/services/overview) 
and [ziti identity](https://docs.openziti.io/docs/learn/core-concepts/identities/overview) that has permission to `Bind`
to the service.
[Caddyfile.server](Caddyfile.server) is provided as an example. It runs a file server that is only accessible over the ziti network.
```
$ export ZITI_IDENTITY=caddy-server.json
$ export ZITI_SERVICE=caddy-http
$ ./ziti-caddy run --config Caddyfile.server
```

#### Proxy side
To configure caddy to proxy request over ziti network to a ziti service your need to configure `reverseproxy.transport`
```
reverse_proxy <intercept address> {
	# intercept address, will be used if service id not set in the 'ziti' block
	transport ziti {
		# required, ziti identity file
		identity caddy-proxy.json

		# optional, intercept address will used if service is not specified
		service caddy-http
		terminator 
	}
}
```

For this exercise you need to create a [ziti service](https://docs.openziti.io/docs/learn/core-concepts/services/overview)
and [ziti identity](https://docs.openziti.io/docs/learn/core-concepts/identities/overview) that has permission to `Dial` the service.
[Caddyfile.proxy](Caddyfile.proxy) is provided as an example. It runs a reverse proxy that is forwards requests to the ziti service.
```
$ export ZITI_IDENTITY=caddy-proxy.json
$ export ZITI_SERVICE=caddy-http
$ ./ziti-caddy run --config Caddyfile.proxy
```