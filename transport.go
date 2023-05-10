package ziti_caddy

import (
	"context"
	"fmt"
	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp/reverseproxy"
	"github.com/openziti/sdk-golang/ziti"
	"net"
	"net/http"
	"time"
)

var zitiContexts = ziti.NewSdkCollection()

func init() {
	zitiContexts.ConfigTypes = append(zitiContexts.ConfigTypes, ziti.InterceptV1, ziti.ClientConfigV1)

	caddy.RegisterModule(ZitiTransport{})
	caddy.RegisterNetwork("ziti", newZitiListener)
}

type ZitiTransport struct {
	Identity   string `json:"identity,omitempty"`
	Service    string `json:"service,omitempty"`
	Terminator string `json:"terminator,omitempty"`

	ztx   ziti.Context
	caddy caddy.Context
	proxy *reverseproxy.HTTPTransport
}

func (z *ZitiTransport) Cleanup() error {
	var err error
	if z.ztx != nil {
		zitiContexts.Remove(z.ztx)
		z.ztx.Close()
		z.ztx = nil
	}

	if z.proxy != nil {
		err = z.proxy.Cleanup()
		z.proxy = nil
	}
	return err
}

func (z *ZitiTransport) RoundTrip(request *http.Request) (*http.Response, error) {
	return z.proxy.RoundTrip(request)
}

func (z *ZitiTransport) UnmarshalCaddyfile(d *caddyfile.Dispenser) error {

	if !d.Next() {
		return fmt.Errorf("missing ziti transport configuration")
	}

	for d.NextBlock(0) {
		switch d.Val() {
		case "identity":
			d.Args(&z.Identity)
		case "service":
			d.Args(&z.Service)
		case "terminator":
			d.Args(&z.Terminator)
		default:
			return fmt.Errorf("unsupported configuration option `%s`", d.Val())
		}
	}
	return nil
}

// CaddyModule returns the Caddy module information.
func (ZitiTransport) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID: "http.reverse_proxy.transport.ziti",
		New: func() caddy.Module {
			return new(ZitiTransport)
		},
	}
}

func (z *ZitiTransport) Provision(ctx caddy.Context) error {
	log := ctx.Logger(z)
	log.Info("ZitiTransport is loading")

	z.caddy = ctx

	// default proxy roundTripper
	ctr := new(reverseproxy.HTTPTransport)
	err := ctr.Provision(z.caddy)
	if err != nil {
		return err
	}

	opts := &ziti.Options{
		RefreshInterval: 30 * time.Second,
	}
	z.ztx, err = zitiContexts.NewContextFromFileWithOpts(z.Identity, opts)
	if err != nil {
		return err
	}

	err = z.ztx.Authenticate()
	if err != nil {
		return err
	}

	if z.Service != "" {
		options := &ziti.DialOptions{
			ConnectTimeout: time.Duration(ctr.DialTimeout),
			Identity:       z.Terminator,
		}
		ctr.Transport.DialContext = func(ctx context.Context, network, addr string) (net.Conn, error) {
			return z.ztx.DialWithOptions(z.Service, options)
		}
	} else {
		ctr.Transport.DialContext = func(ctx context.Context, network, addr string) (net.Conn, error) {
			return z.ztx.DialAddr(network, addr)
		}
	}

	z.proxy = ctr

	return err
}

var (
	_ caddyfile.Unmarshaler = (*ZitiTransport)(nil)
	_ http.RoundTripper     = (*ZitiTransport)(nil)
	_ caddy.CleanerUpper    = (*ZitiTransport)(nil)
)
