package bid

import (
	"airpush/auction/dsp"
	"sync"
	"time"
)

// settings setter
type BidOption func(*Bid)

// SetDsp
func SetDsp(d *dsp.Dsp) BidOption {
	return func(t *Bid) {
		t.dsp = d
	}
}

// Bid
type Bid struct {
	mu sync.Mutex
	dsp *dsp.Dsp
	res *RtbResponse
	err []string
}

// Bid.New()
func New(opts ...BidOption) (proto *Bid) {

	proto = &Bid{}

	// set custom settings
	for _, opt := range opts {
		opt(proto)
	}
	return
}

// execute bid request
func (b *Bid) Do() {

	startTime := time.Now()
	b.res = &RtbResponse{
		Dsp: b.dsp.GetName(),
	}

	// calc request time
	defer func() {
		b.res.Build = time.Since(startTime).String()
	}()

	buf, err := b.dsp.GetClient().Do()
	if err != nil {
		b.err = append(b.err, err.Error())
		return
	}

	res := new(BidResponse)
	err = res.UnmarshalJSON(buf)
	if err != nil {
		b.err = append(b.err, err.Error())
		return
	}

	b.res.Bid = *res
}

// get bid response
func (b *Bid) GetRes() *RtbResponse {
	defer b.mu.Unlock()
	b.mu.Lock()

	return b.res
}

// get bid errors
func (b *Bid) GetErr() []string {
	defer b.mu.Unlock()
	b.mu.Lock()

	return b.err
}

// order bids by cpm
type OrderBids []*Bid

func (a OrderBids) Len() int      { return len(a) }
func (a OrderBids) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a OrderBids) Less(i, j int) bool {
	return a[i].GetRes().Bid.Cpm > a[j].GetRes().Bid.Cpm
}
