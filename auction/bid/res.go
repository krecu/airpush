package bid

type RtbResponse struct {
	Dsp string `json:"dsp"`
	Build string `json:"time_req"`
	Bid BidResponse
}

type BidResponse struct {
	Wait string `json:"time_wait"`
	Cpm float64 `json:"cpm"`
}