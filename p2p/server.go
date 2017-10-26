package p2p

import (
	"fmt"
	"net"
)

type ServerContext struct {
	willShutdown chan bool
	port         int
	handler      HandleConnection
}

func NewServer(port int, handler HandleConnection) *ServerContext {
	return &ServerContext{
		willShutdown: make(chan bool),
		port:         port,
		handler:      handler,
	}
}

func (ctx *ServerContext) Start() error {
	addr, err := net.ResolveTCPAddr("tcp4", fmt.Sprintf(":%d", ctx.port))
	if err != nil {
		return err
	}

	listener, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return err
	}

	defer listener.Close()

	incoming := make(chan net.Conn)
	go acceptConnections(listener, incoming)

	for {
		select {
		case <-ctx.willShutdown:
			break
		case conn := <-incoming:
			go ctx.handler(conn)
		}
	}

	return nil
}

func acceptConnections(listener *net.TCPListener, incoming chan net.Conn) {
	for {
		if conn, err := listener.Accept(); err != nil {
			panic(err)
		} else {
			incoming <- conn
		}
	}
}

func (ctx *ServerContext) Shutdown() {
	ctx.willShutdown <- true
}
