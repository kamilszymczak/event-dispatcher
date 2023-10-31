package track

import (
	"encoding/json"
	"github.com/kamilszymczak/event-dispatcher/scheduler"
	"io"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/kamilszymczak/event-dispatcher/request"
	"github.com/kamilszymczak/event-dispatcher/response"
)

type Tracker[T response.Payload] interface {
	AddRequest(reqs ...request.RequestableRefreshRater)
	JobsRunning() *sync.WaitGroup
	RefreshRate(rate time.Duration)
	Listen() <-chan response.Responser
}

type Track[T response.Payload] struct {
	responses   response.Response[T]
	refreshRate time.Duration
	requests    []request.RequestableRefreshRater
	computeFunc func(int, int) bool
	jobsRunning sync.WaitGroup
}

func New[T response.Payload]() (*Track[T], error) {
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

func (t *Track[T]) AddRequest(reqs ...request.RequestableRefreshRater) {
	for _, req := range reqs {
		if !req.HasRefreshRate() {
			req.SetRefreshRate(t.refreshRate)
		}
	}
	t.requests = append(t.requests, reqs...)
}

func (t *Track[T]) Listen() <-chan response.Responser {
	out := make(chan response.Responser)
	t.jobsRunning.Add(len(t.requests))

	for _, req := range t.requests {
		go t.executeJob(out, req)
	}

	go t.waitForJobsToComplete(out)
	return out
}

func (t *Track[T]) executeJob(ch chan<- response.Responser, req request.RequestableRefreshRater) {
	defer t.jobsRunning.Done()

	job := scheduler.Every(req.GetRefreshRate()).Repeat(3).Do(fetchData[T], ch, req)
	job.Wait()
}

//type DataFetcher[T responseSource.Payload] interface {
//	fetchData(channel chan<- request.RequestUrler, request request.RequestUrler)
//}
//
//type DataFetch struct{}

func fetchData[T response.Payload](channel chan<- response.Responser, request request.RequestableRefreshRater) {
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

	resp := response.New[T](*request.GetRequest(), output)

	channel <- resp.GetData()
}

func (t *Track[T]) waitForJobsToComplete(ch chan<- response.Responser) {
	defer close(ch)
	t.jobsRunning.Wait()
}

func (t *Track[T]) SetComputeFunc(fn func(current int, new int) bool) {
	t.computeFunc = fn
}
