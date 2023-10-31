package request

import (
	"time"
)

type RequestUrler interface {
	GetUrl() string
}

type RefreshRater interface {
	SetRefreshRate(rate time.Duration)
	GetRefreshRate() time.Duration
	HasRefreshRate() bool
}

type RequestableRefreshRater interface {
	RequestUrler
	RefreshRater
	GetRequest() *Request
}

type Request struct {
	Url         string
	RefreshRate *time.Duration
}

func New(url string) *Request {
	request := &Request{
		Url:         url,
		RefreshRate: nil,
	}
	return request
}

func (r *Request) GetUrl() string {
	return r.Url
}

func (r *Request) GetRefreshRate() time.Duration {
	return *r.RefreshRate
}

func (r *Request) SetRefreshRate(rate time.Duration) {
	r.RefreshRate = &rate
}

func (r *Request) HasRefreshRate() bool {
	return r.RefreshRate != nil
}

func (r *Request) GetRequest() *Request {
	return r
}
