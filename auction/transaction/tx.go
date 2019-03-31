package transaction

import (
	"airpush/auction/bid"
	"context"
	"fmt"
	"sync"
	"time"
)

const GLOBAL_TIMEOUT = time.Duration(100) * time.Millisecond

// settings setter
type TransactionOption func(*Transaction)

// global tx timeout
func SetTimeout(duration time.Duration) TransactionOption {
	return func(t *Transaction) {
		t.timeout = duration
	}
}

// SetBids
func SetBids(bids []*bid.Bid) TransactionOption {
	return func(t *Transaction) {
		t.bids = bids
	}
}

// root struct
type Transaction struct {
	bids []*bid.Bid
	timeout time.Duration
}

// new module
func New(opts ...TransactionOption) (proto *Transaction) {

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
func (tx *Transaction) Do() (err error) {

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
				defer wg.Done()
				b.Do()
			}(b)
		}

		wg.Wait()

		// signal complete requests
		done <- true
	}(tx.bids)

	select {
	case <-done:
		return
	case <-ctx.Done():
		err = fmt.Errorf("tx timeout")
		return
	}
}