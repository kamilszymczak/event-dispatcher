package track

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/kamilszymczak/event-dispatcher/request"
	"github.com/kamilszymczak/event-dispatcher/requestSource"
)

type Track[T requestSource.Payload] struct {
	refreshRate time.Duration
	requests []request.Requestable[T]
}

func New[T requestSource.Payload]() (*Track[T], error) {
	track := &Track[T]{
		refreshRate: time.Minute * 1,
	}

	return track, nil
}

func (t *Track[T]) RefreshRate(rate time.Duration) {
	t.refreshRate = rate
}

func (t *Track[T]) AddRequest(reqs ...request.Requestable[T]) {
	t.requests = append(t.requests, reqs...)
}

func (t *Track[T]) Listen() (<-chan request.Requestable[T]) {
	out := make(chan request.Requestable[T])

	for i, request := range t.requests {
		go fetchData(out, i, request)
	}

	return out
}

func fetchData[T requestSource.Payload](channel chan<- request.Requestable[T], i int, request request.Requestable[T]) {
	res, err := http.Get(request.GetUrl())

	if i < 1 {
		fmt.Println("sleeping 2 sec " , i)
		time.Sleep(2 * time.Second)
		fmt.Println("done sleeping")
	}

	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)

	var output T
	json.Unmarshal(body, &output)

	request.SetData(output)
	channel <- request
}
