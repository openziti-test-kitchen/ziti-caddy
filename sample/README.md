Complete ziti-caddy Sample
-----

This folder contains files that would allow you to run a 
complete end-to-end zitified Caddyserver sample.

# Setup Prerequisites 
* Golang - pick a method approprate for your platform
* OpenZiti CLI - install or download a release from [Github](https://github.com/openziti/ziti/releases/latest)
      or build your own
```shell
$ go install github.com/openziti/ziti/ziti@latest
```
* OpenZiti network - any of the quickstart methods documented [here](https://openziti.io/docs/learn/quickstarts/)
> **_NOTE_**: for development purposes you can create a transient OpenZiti network with my new
> favorite:
> ```
> ziti egde quicktart
> ```

# Setup
Run [`ziti-init.sh`](./ziti-init.sh) script. It creates the following:
- `caddy-service` - OpenZiti service
- `caddy-host` - identity to host the service
- `caddy-client` - identity to access the service
along with all necessary policies

# Run Caddy server on overlay network
In this first exercise we are going to run zitified Caddyserver in the dark mode:
- no open/listening ports
- service is only available on the OpenZiti overlay network

[Caddyfile.server](Caddyfile.server) is configured to use `caddy-host.json` identity 
and to bind to `caddy-service`. 
It serves up the content of the file system to client on the overlay network.

```shell
$ cd sample
$ go run ../cmd/ziti-caddy run --config Caddyfile.server
2023/09/29 13:54:45.004	INFO	using provided configuration	{"config_file": "Caddyfile.server", "config_adapter": ""}
2023/09/29 13:54:45.005	WARN	admin	admin endpoint disabled
2023/09/29 13:54:45.005	WARN	http	server is listening only on the HTTP port, so no automatic HTTPS will be applied to this server	{"server_name": "srv0", "http_port": 80}
2023/09/29 13:54:45.005	INFO	tls.cache.maintenance	started background certificate maintenance	{"cache": "0xc000a2ce00"}
2023/09/29 13:54:45.006	INFO	tls	cleaning storage unit	{"description": "FileStorage:/home/eugene/.local/share/caddy"}
2023/09/29 13:54:45.007	INFO	tls	finished cleaning storage units
2023/09/29 13:54:45.060	INFO	http.log	server running	{"name": "srv0", "protocols": ["h1", "h2", "h3"]}
2023/09/29 13:54:45.061	INFO	autosaved config (load with --resume flag)	{"file": "/home/eugene/.config/caddy/autosave.json"}
2023/09/29 13:54:45.061	INFO	serving initial configuration
INFO[0000] new service session                           session token=f3709f65-edf9-43ba-b8b0-aec3d4ebb410

```

After the server is started you can check that it opened no listening ports. 
You'll need another zitified application or tunneler to access the service. 
Luckily for us, we can use a sample ziti-embedded app from [OpenZiti Golang SDK repo](https://github.com/openziti/sdk-golang)

```shell
$ go run github.com/openziti/sdk-golang/example/curlz@latest http://caddy-service caddy-client.json
<!DOCTYPE html>
<html>
	<head>
		<title>/</title>
		<meta charset="utf-8">
		<meta name="viewport" content="width=device-width, initial-scale=1.0">
<style>
* { padding: 0; margin: 0; }

body {
	font-family: sans-serif;
	text-rendering: optimizespeed;
...
```

# Run Caddy reverse_proxy to a sevice on overlay network
In this second exercise we are using a zitified Caddyserver to proxy to a service on the OpenZiti network.
The service is provided by the Caddyservice instance is [step 1](#run-caddy-server-on-overlay-network)
In this case Caddyserver is [configured](./Caddyfile.proxy) to accept HTTP requests `localhost:8080` 
and proxy them with `reverse_proxy` module that uses ziti transport connecting to a service.

```shell
$ cd sample
$ go run ../cmd/ziti-caddy run --config Caddyfile.proxy
2023/09/29 13:55:09.269	INFO	using provided configuration	{"config_file": "Caddyfile.proxy", "config_adapter": ""}
2023/09/29 13:55:09.270	WARN	admin	admin endpoint disabled
2023/09/29 13:55:09.270	WARN	http	server is listening only on the HTTP port, so no automatic HTTPS will be applied to this server	{"server_name": "srv0", "http_port": 8080}
2023/09/29 13:55:09.270	INFO	tls.cache.maintenance	started background certificate maintenance	{"cache": "0xc0003b5880"}
2023/09/29 13:55:09.270	INFO	http.reverse_proxy.transport.ziti	ZitiTransport is loading
```

And no we can use a web browser or other tools to get the response:
```shell
$ curl -s localhost:8080
<!DOCTYPE html>
<html>
	<head>
		<title>/</title>
		<meta charset="utf-8">
		<meta name="viewport" content="width=device-width, initial-scale=1.0">
<style>
* { padding: 0; margin: 0; }

body {
	font-family: sans-serif;
	text-rendering: optimizespeed;
	background-color: #ffffff;
...
```

# OpenZiti module flexibility
[Combined configuration](Caddyfile.combined) file merges the two above exercises into a single process: `localhost:8080` is 
proxied over OpenZiti overlay back into the same Caddyserver process to a Ziti listener serving files.

It is not very useful on its own but shows the flexibility of OpenZiti Caddy module (and OpenZiti SDK):
multiple identities and services can be used at the same time within the same process.

