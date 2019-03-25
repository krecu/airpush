package bid

import (
	"airpush/client"
	"context"
	"sync"
	"time"
)

type Response struct {
	Time time.Duration
	Data interface{}
}

type Bid struct {
	mu sync.Mutex

	res Response
	client *client.Client
	ctx context.Context
	err []string
}

func (b *Bid) WithContext(ctx context.Context) *Bid {
	b.ctx = ctx
	return b
}

func (b *Bid) Do() {
	b.client.Do()
	return
}

func (b *Bid) AddErr(err string) {
	defer b.mu.Unlock()
	b.mu.Lock()

	b.err = append(b.err, err)
}