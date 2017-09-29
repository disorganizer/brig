package client

import (
	"context"
	"net"

	"github.com/disorganizer/brig/brigd/capnp"
	"zombiezen.com/go/capnproto2/rpc"
)

type Client struct {
	ctx  context.Context
	conn *rpc.Conn

	api capnp.API
}

func Connect(ctx context.Context) (*Client, error) {
	conn, err := net.Dial("tcp", "localhost:6666")
	if err != nil {
		return nil, err
	}

	transport := rpc.StreamTransport(conn)
	clientConn := rpc.NewConn(transport)
	api := capnp.API{Client: clientConn.Bootstrap(ctx)}

	return &Client{
		ctx:  context.Background(),
		conn: clientConn,
		api:  api,
	}, nil
}

func (cl *Client) Close() error {
	return cl.conn.Close()
}
