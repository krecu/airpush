package dsp

import "airpush/client"

type Dsp struct {
	name string
	client *client.Client
}

func New(name string, client *client.Client) *Dsp {
	return &Dsp{
		client: client,
		name: name,
	}
}

func (dsp *Dsp) GetClient() *client.Client {
	return dsp.client
}

func (dsp *Dsp) GetName() string {
	return dsp.name
}