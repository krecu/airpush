package auction

type Auction struct {

}

func New() (proto *Auction) {
	proto = &Auction{}
	// init clients for persisten connection
	return
}

func (a *Auction) Do() {
	//build bids
	//run tx
	//return best bid
}