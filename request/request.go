package request

import (
	"time"
)

type RequestUrler interface {
	Url() string
}

type RefreshRater interface {
	SetRefreshRate(rate time.Duration)
	RefreshRate() time.Duration
	HasRefreshRate() bool
}

type RequestableRefreshRater interface {
	RequestUrler
	RefreshRater
	Request() *Request
}

type Request struct {
	url         string
	refreshRate *time.Duration
}

func New(url string) *Request {
	request := &Request{
		url:         url,
		refreshRate: nil,
	}
	return request
}

func (r *Request) Url() string {
	return r.url
}

func (r *Request) RefreshRate() time.Duration {
	return *r.refreshRate
}

func (r *Request) SetRefreshRate(rate time.Duration) {
	r.refreshRate = &rate
}

func (r *Request) HasRefreshRate() bool {
	return r.refreshRate != nil
}

func (r *Request) Request() *Request {
	return r
}
