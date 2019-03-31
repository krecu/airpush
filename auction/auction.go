package auction

import (
	"airpush/auction/bid"
	"airpush/auction/dsp"
	"airpush/auction/transaction"
	"fmt"
	"sort"
	"time"
)

// settings setter
type AuctionOption func(*Auction)

// dsps list
func SetDsp(dsp []*dsp.Dsp) AuctionOption {
	return func(a *Auction) {
		a.dsp = dsp
	}
}

// global tx timeout
func SetTimeout(duration time.Duration) AuctionOption {
	return func(a *Auction) {
		a.timeout = duration
	}
}

type Auction struct {
	timeout time.Duration
	dsp []*dsp.Dsp
}

func New(opts ...AuctionOption) (proto *Auction) {
	proto = &Auction{}

	// set custom settings
	for _, opt := range opts {
		opt(proto)
	}

	return
}

func (a *Auction) Do() (win *bid.Bid, err error) {

	// formed bids
	var rBids, nBids []*bid.Bid

	for _, d := range a.dsp {
		rBids = append(rBids, bid.New(bid.SetDsp(d)))
	}

	err = transaction.New(transaction.SetTimeout(a.timeout), transaction.SetBids(rBids)).Do()
	if err != nil {
		return
	}

	// filter good bids
	for _, b := range rBids {
		if len(b.GetErr()) == 0 {
			nBids = append(nBids, b)
		}
	}

	if len(nBids) > 0 {
		// sort
		sort.Sort(bid.OrderBids(nBids))
		win = nBids[0]
	} else {
		err = fmt.Errorf("empty auction")
	}

	return
}