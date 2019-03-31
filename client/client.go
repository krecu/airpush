package client

import (
	"airpush/client/transport"
	"context"
	"fmt"
	"time"
)

// support connection type
const CONN_TYPE_HTTP  = "http"
const CONN_TYPE_GRPC  = "grpc"

// settings setter
type ClientOption func(*Client)

// SetAddr
func SetAddr(addr string) ClientOption {
	return func(t *Client) {
		t.addr = addr
	}
}

// WithTimeout
func WithTimeout(duration time.Duration) ClientOption {
	return func(t *Client) {
		t.timeout = duration
	}
}

// SetConnectionType
func SetConnectionType(ctype string) ClientOption {
	return func(t *Client) {
		t.cType = ctype
	}
}

// Client struct
// param: cType - connection type
// param: addr - endpoint address
type Client struct {
	cType string
	addr string
	timeout time.Duration
	transport transport.Transport
}

// construct client
func New(opts ...ClientOption) (proto *Client, err error) {

	proto = new(Client)

	// set custom transport params
	for _, opt := range opts {
		opt(proto)
	}

	// init transport type
	switch proto.cType {
	case CONN_TYPE_HTTP:
		proto.transport = transport.NewHttpTransport(transport.SetAddr(proto.addr))
	case CONN_TYPE_GRPC:
		err = fmt.Errorf("grpc transprt no implement")
	}


	return
}

// execute request
func (c *Client) Do() (buf []byte, err error){

	done := make(chan bool)

	// timeout
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	go func() {
		buf, err = c.transport.Do(ctx)
		done <- true
	}()

	select {
	case <-done:
		return
	case <-ctx.Done():
		err = fmt.Errorf("requet timeout")
		return
	}
}