package ziti_caddy

import (
	"context"
	"fmt"
	"github.com/openziti/sdk-golang/ziti"
	"github.com/openziti/sdk-golang/ziti/edge"
	"net"
	"strings"
)

type zitiListener struct {
	ztx        ziti.Context
	service    string
	terminator string
	l          edge.Listener
}

func (z *zitiListener) String() string {
	return fmt.Sprintf("%s/%s", z.Network(), z.Address())
}

func (z *zitiListener) Network() string {
	return "ziti"
}

func (z *zitiListener) Address() string {
	if z.terminator == "" {
		return z.service
	} else {
		return z.service + "/" + z.terminator
	}
}

func (z *zitiListener) Accept() (net.Conn, error) {
	return z.l.Accept()
}

func (z *zitiListener) Close() error {
	_ = z.l.Close()
	z.ztx.Close()
	return nil
}

func (z *zitiListener) Addr() net.Addr {
	return z
}

func newZitiListener(
	ctx context.Context, network, addr, portRange string, portOffset uint,
	cfg net.ListenConfig) (any, error) {
	// address always ends with ":port", so just strip it
	// the rest should be "<service>[/<terminator>]@<identity>
	addr = strings.Split(addr, ":")[0]

	s := strings.Split(addr, "@")
	srv, identity := s[0], s[1]

	s = strings.Split(srv, "/")
	service := s[0]
	terminator := ""
	if len(s) > 1 {
		terminator = s[1]
	}

	ztx, err := ziti.NewContextFromFile(identity)
	if err != nil {
		return nil, err
	}

	opts := &ziti.ListenOptions{
		Identity: terminator,
	}

	conn, err := ztx.ListenWithOptions(service, opts)
	if err != nil {
		ztx.Close()
		return nil, err
	}

	l := &zitiListener{
		ztx: ztx,
		l:   conn,
	}
	return l, nil
}

var (
	_ net.Addr     = (*zitiListener)(nil)
	_ net.Listener = (*zitiListener)(nil)
)
