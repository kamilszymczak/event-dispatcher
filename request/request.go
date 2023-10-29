package request

import (
	"time"

	"github.com/kamilszymczak/event-dispatcher/requestSource"
)

type Requestable[T requestSource.Payload] interface {
	GetUrl() string
	GetData() T
	SetData(T) T
}

type RefreshRater interface {
	SetRefreshRate(rate time.Duration)
	GetRefreshRate() time.Duration
	HasRefreshRate() bool
}

type RequestableRefreshRater[T requestSource.Payload] interface {
	Requestable[T]
	RefreshRater
}

type Request[T requestSource.Payload] struct {
	Url         string
	Data        T
	RefreshRate *time.Duration
}

func (r *Request[T]) GetUrl() string {
	return r.Url
}

func (r *Request[T]) GetData() T {
	return r.Data
}

func (r *Request[T]) SetData(data T) T {
	r.Data = data
	return data
}

func (r *Request[T]) GetRefreshRate() time.Duration {
	return *r.RefreshRate
}

func (r *Request[T]) SetRefreshRate(rate time.Duration) {
	r.RefreshRate = &rate
}

func (r *Request[T]) HasRefreshRate() bool {
	return r.RefreshRate != nil
}

func New[T requestSource.Payload](url string) RequestableRefreshRater[T] {
	request := &Request[T]{
		Url:         url,
		RefreshRate: nil,
	}
	return request
}
