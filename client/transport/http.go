package transport

import (
	"context"
	"io/ioutil"
	"net/http"
)

// settings setter
type HttpTransportOption func(*BaseHttpTransport)

// addr
func SetAddr(addr string) HttpTransportOption {
	return func(t *BaseHttpTransport) {
		t.addr = addr
	}
}

// context
func WithContext(ctx context.Context) HttpTransportOption {
	return func(t *BaseHttpTransport) {
		t.ctx = ctx
	}
}

// BaseHttpTransport
type BaseHttpTransport struct {
	client *http.Client
	addr string
	ctx context.Context
}

// NewHttpTransport
func NewHttpTransport(opts ...HttpTransportOption) (proto *BaseHttpTransport) {

	proto = &BaseHttpTransport{
		ctx: context.Background(),
	}

	// set custom transport params
	for _, opt := range opts {
		opt(proto)
	}

	proto.client = &http.Client{
		// @todo - control in parent context
		//Timeout: time.Duration(200) * time.Millisecond,
		//Transport: &http.Transport{
		//	DialContext: (&net.Dialer{
		//		Timeout:   time.Duration(200) * time.Millisecond,
		//		KeepAlive: time.Duration(200) * time.Millisecond,
		//		DualStack: true,
		//	}).DialContext,
		//	MaxIdleConns:          100,
		//	IdleConnTimeout:       time.Duration(200) * time.Millisecond,
		//	TLSHandshakeTimeout:   time.Duration(200) * time.Millisecond,
		//	ExpectContinueTimeout: time.Duration(200) * time.Millisecond,
		//},
	}

	return
}

// transport.Do
func (t *BaseHttpTransport) Do(ctx context.Context) ([]byte, error) {

	// init request
	req, err := http.NewRequest(http.MethodGet, t.addr, nil)
	if err != nil {
		return nil, err
	}

	// inherit parent context
	req = req.WithContext(ctx)
	res, err := t.client.Do(req)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	defer func() {
		_ = res.Body.Close()
	}()

	return body, nil
}