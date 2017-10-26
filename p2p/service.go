package p2p

import (
	"net"
)

type HandleConnection func(conn net.Conn)

type Service struct {
	server  *ServerContext
	clients []*ClientContext
}

func New(port int, serverConn HandleConnection) *Service {
	srv := &Service{
		server: NewServer(port, serverConn),
	}

	return srv
}

func (srv *Service) Start() error {
	return srv.server.Start()
}

func (srv *Service) Shutdown() {
	for _, cl := range srv.clients {
		cl.Shutdown()
	}
	srv.clients = nil

	srv.server.Shutdown()
}

func (srv *Service) AddPeer(conf *ClientConfig) error {
	peer := NewClient(conf.Address, conf.Port, conf.Handler)
	if err := peer.Start(); err != nil {
		return err
	}

	srv.clients = append(srv.clients, peer)

	return nil
}
