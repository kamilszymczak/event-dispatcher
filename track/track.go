package track

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/kamilszymczak/event-dispatcher/request"
	"github.com/kamilszymczak/event-dispatcher/requestSource"
	"github.com/kamilszymczak/event-dispatcher/scheduler"
)

type Track[T requestSource.Payload] struct {
	refreshRate time.Duration
	requests    []request.RequestableRefreshRater[T]
	computeFunc func(int, int) bool
	jobsRunning sync.WaitGroup
}

func New[T requestSource.Payload]() (*Track[T], error) {
	track := &Track[T]{
		refreshRate: 2 * time.Second,
	}

	return track, nil
}

func (t *Track[T]) JobsRunning() *sync.WaitGroup {
	return &t.jobsRunning
}

func (t *Track[T]) RefreshRate(rate time.Duration) {
	t.refreshRate = rate
}

func (t *Track[T]) AddRequest(reqs ...request.RequestableRefreshRater[T]) {
	for _, req := range reqs {
		if !req.HasRefreshRate() {
			req.SetRefreshRate(t.refreshRate)
		}
	}
	t.requests = append(t.requests, reqs...)
}

func (t *Track[T]) Listen() <-chan request.Requestable[T] {
	out := make(chan request.Requestable[T])
	t.jobsRunning.Add(len(t.requests))

	for _, request := range t.requests {
		go t.executeJob(out, request)
	}

	go t.waitForJobsToComplete(out)

	return out
}

func (t *Track[T]) executeJob(ch chan<- request.Requestable[T], req request.RequestableRefreshRater[T]) {
	defer t.jobsRunning.Done()

	job := scheduler.Every(req.GetRefreshRate()).Repeat(3).Do(fetchData[T], ch, req)
	job.Wait()
}

type DataFetcher[T requestSource.Payload] interface {
	fetchData(channel chan<- request.Requestable[T], request request.Requestable[T])
}

type DataFetch struct{}

func fetchData[T requestSource.Payload](channel chan<- request.Requestable[T], request request.Requestable[T]) {
	res, err := http.Get(request.GetUrl())

	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)

	var output T
	if err := json.Unmarshal(body, &output); err != nil {
		log.Fatal(err)
	}

	request.SetData(output)
	channel <- request
}

func (t *Track[T]) waitForJobsToComplete(ch chan<- request.Requestable[T]) {
	defer close(ch)
	t.jobsRunning.Wait()
}

func (t *Track[T]) SetComputeFunc(fn func(current int, new int) bool) {
	t.computeFunc = fn
}
