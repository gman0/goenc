package p2p

import (
	"fmt"
	"net"
)

type ClientContext struct {
	addr    string
	port    int
	handler HandleConnection
	conn    net.Conn
}

type ClientConfig struct {
	Address string
	Port    int
	Handler HandleConnection
}

func NewClient(addr string, port int, handler HandleConnection) *ClientContext {
	return &ClientContext{
		addr:    addr,
		port:    port,
		handler: handler,
	}
}

func (ctx *ClientContext) Start() error {
	service := fmt.Sprintf("%s:%d", ctx.addr, ctx.port)

	serverAddr, err := net.ResolveTCPAddr("tcp", service)
	if err != nil {
		return err
	}

	conn, err := net.DialTCP("tcp", nil, serverAddr)
	if err != nil {
		return err
	}

	ctx.handler(conn)

	return nil
}

func (ctx *ClientContext) Shutdown() error {
	if ctx.conn != nil {
		return ctx.conn.Close()
	}

	return nil
}
