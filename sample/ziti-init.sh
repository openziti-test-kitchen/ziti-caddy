
# create ziti service and policies
ziti edge create service caddy-service
ziti edge create service-edge-router-policy caddy-serp --service-roles '@caddy-service' --edge-router-roles '#all'

# create caddy-host identity and policies
ziti edge create identity caddy-host -o ./caddy-host.jwt && ziti edge enroll ./caddy-host.jwt -o ./caddy-host.json
ziti edge create edge-router-policy caddy-host-erp --identity-roles '@caddy-host' --edge-router-roles '#all'
ziti edge create service-policy caddy-hosting Bind --service-roles '@caddy-service' --identity-roles '@caddy-host'

# create caddy-client identity and policies
ziti edge create identity caddy-client -o ./caddy-client.jwt && ziti edge enroll ./caddy-client.jwt -o ./caddy-client.json
ziti edge create edge-router-policy caddy-client-erp --identity-roles '@caddy-client' --edge-router-roles '#all'
ziti edge create service-policy caddy-dialing Dial --service-roles '@caddy-service' --identity-roles '@caddy-client'

rm -f ./caddy-client.jwt ./caddy-host.jwt