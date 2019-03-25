package tx

import (
	"airpush/auction/bid"
	"context"
	"sync"
	"time"
)

const GLOBAL_TIMEOUT = time.Duration(100) * time.Millisecond

// settings setter
type TransactionSetOption func(*Transaction)

// global tx timeout
func SetTimeout(ttl int) TransactionSetOption {
	return func(t *Transaction) {
		t.timeout = time.Duration(ttl) * time.Millisecond
	}
}

// root struct
type Transaction struct {
	timeout time.Duration
}

// new module
func New(opts ...TransactionSetOption) (proto *Transaction) {

	proto = &Transaction{
		timeout: GLOBAL_TIMEOUT,
	}

	// set custom settings
	for _, opt := range opts {
		opt(proto)
	}

	return
}

// execute transaction
func (tx *Transaction) Exec(bids []*bid.Bid) {

	done := make(chan bool)

	// tx timeout
	ctx, cancel := context.WithTimeout(context.Background(), tx.timeout)
	defer cancel()

	// parent rutine for async requests
	go func(bids []*bid.Bid) {

		var wg sync.WaitGroup
		wg.Add(len(bids))

		for _, b := range bids {

			// rutine for single async request
			go func(b *bid.Bid) {
				defer wg.Add(1)
				b.WithContext(ctx).Do()
			}(b)
		}

		wg.Wait()

	}(bids)

	select {
	case <-done:
		return
	case <-ctx.Done():
		return
	}
}