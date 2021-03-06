package server

import (
	"net"

	"github.com/mholt/caddy"
	"github.com/valyala/fasthttp"
)

// make sure FastServer implement GracefulServer
var _ caddy.GracefulServer = new(FastServer)

func NewFastServer(cfg ServerConfig) *FastServer {
	srv := &FastServer{
		Addr:   cfg.Addr,
		Server: cfg.makeServer(),
	}
	return srv
}

type FastServer struct {
	*fasthttp.Server
	Addr string
}

func (s *FastServer) Listen() (net.Listener, error) {
	ln, err := net.Listen("tcp4", s.Addr)
	if err != nil {
		return nil, err
	}
	if s.TCPKeepalive {
		if tcpln, ok := ln.(*net.TCPListener); ok {
			return &tcpKeepaliveListener{
				TCPListener:     tcpln,
				keepalivePeriod: s.TCPKeepalivePeriod,
			}, nil
		}
	}
	return ln, nil
}

func (s *FastServer) WrapListener(ln net.Listener) net.Listener {
	return ln
}

func (s *FastServer) Serve(ln net.Listener) error {
	return s.Server.Serve(ln)
}

func (s *FastServer) ListenPacket() (net.PacketConn, error) {
	return nil, nil
}

func (s *FastServer) ServePacket(net.PacketConn) error {
	return nil
}

func (s *FastServer) Address() string {
	return s.Addr
}

func (s *FastServer) Stop() error {
	return s.Server.Shutdown()
}
