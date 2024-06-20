package poller

import (
	"errors"
	"fmt"
	"io"
	"log"
	"log/slog"
	"net/http"
	"sync"
	"time"

	"github.com/jonboulle/clockwork"
	"github.com/kamilszymczak/event-dispatcher/config"
	"github.com/kamilszymczak/event-dispatcher/response"
	"github.com/kamilszymczak/event-dispatcher/scheduler"
)

type Event struct {
	Response   response.ResponseAccessor
	Observable *Observable
}

type Observable struct {
	Address  string
	interval *time.Duration
}

type dispatchFunc func(response.ResponseAccessor) bool

type poller struct {
	apiUrl       string
	interval     time.Duration
	observables  []Observable
	eventChan    chan Event
	errorHandler func(err error)
	jobsRunning  sync.WaitGroup
	responseType response.Responser
	jobs         map[string]*scheduler.Job
	dispatchFunc dispatchFunc
	shuffle 	 *delay
	stopObservableAfterDispatch	bool
}

func New(url string, responseObject response.Responser) *poller {
	rate := config.GetConfig().FetchRate()

	return &poller{
		apiUrl:       url,
		interval:     time.Duration(rate) * time.Millisecond,
		observables:  make([]Observable, 0),
		eventChan:    make(chan Event),
		responseType: responseObject,
		jobs:		  make(map[string]*scheduler.Job),
		shuffle:	  newDelay(),
	}
}

func (p *poller) Listen() <-chan Event {
	slog.Info("Poller listener started", "poller endpoint", p.apiUrl)
	p.jobsRunning.Add(len(p.observables))

	for _, obs := range p.observables {
		slog.Info("Scheduling job.", "observable", obs)
		scheduleJob(p, obs)
	}

	go p.waitForJobsToComplete()
	return p.eventChan
}

func (t *poller) executeJob(observable Observable) {
	defer t.jobsRunning.Done()

	ticker := clockwork.NewRealClock().NewTicker(*observable.interval)
	job := scheduler.Every(ticker).Do(t.poolData, observable)
	t.jobs[observable.Address] = job
	job.Wait()
}

func (t *poller) fetchData(observable Observable) ([]byte, error) {
	slog.Info("Pooling data started.", "observable", observable.Address)

	resp, err := http.Get(fmt.Sprintf("%s%s", t.apiUrl, observable.Address))
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		slog.Warn("Fetching data unsuccessful.", "observable address", observable.Address, "response status", resp.StatusCode)
		t.jobs[observable.Address].Cancel()
		return nil, fmt.Errorf("fetching data unsuccessful, response status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	return body, nil
}

func parseData(p *poller, body []byte) (response.ResponseAccessor, error) {
	typedResponse, err := p.responseType.Unmarshal(body)

	if err != nil {
		return nil, err
	}

	v, ok := typedResponse.(response.ResponseAccessor)
	if !ok {
		return nil, errors.New("event response does not implement the Responser interface")
	}

	return v, nil
}

func buildEvent(observable Observable, response response.ResponseAccessor) Event {
	return Event{
		Response:   response,
		Observable: &observable,
	}
}

func (t *poller) poolData(observable Observable) {
	data, err := t.fetchData(observable)
	if err != nil {
		log.Print(err.Error())
		return
	}

	parsedResponse, _ := parseData(t, data)
	event := buildEvent(observable, parsedResponse)
	slog.Info("Response parsed and event built.", "observable address", event.Observable.Address,"event response", event.Response)

	if t.dispatchFunc == nil {
		t.eventChan <- event
		return
	}

	if t.dispatchFunc(event.Response) {
		t.eventChan <- event

		if t.stopObservableAfterDispatch {
			slog.Info("Observable's response has been dispatched, cancelling it's job.", "observable", observable.Address)
			t.jobs[observable.Address].Cancel()
		}
	}
}

func (p *poller) GetEventChannel() <-chan Event {
	return p.eventChan
}

func (p *poller) AddObservable(obs ...Observable) {
	for i, o := range obs {
		if o.interval == nil {
			obs[i].interval = &p.interval
		}
	}

	p.observables = append(p.observables, obs...)
}

func (p *poller) waitForJobsToComplete() {
	defer close(p.eventChan)
	p.jobsRunning.Wait()
	slog.Info("All jobs finished running, closing listener channel.")
}

func (p *poller) stopAll() {
	for _, j := range p.jobs {
		j.Stop()
	}
}

func (p *poller) SetDispatchFunc(fn dispatchFunc) {
	p.dispatchFunc = fn
}

func (p *poller) StopObservableAfterDispatched(toggle bool) {
	p.stopObservableAfterDispatch = toggle
}

type delay struct {
	current	time.Duration
	toggle	bool
	gap		time.Duration
}

// Shuffle is used to specify if Observables should have different start times
// instead of all starting at the same time, therefore if same interval will hit
// endpoint at the same time
func (p *poller) Shuffle(toggle bool) {
	p.shuffle.toggle = toggle
}

func newDelay() *delay {
	return &delay{
		current: 0,
		toggle: false,
		gap: time.Duration(config.GetConfig().Request.DelayGap) * time.Millisecond,
	}
}

func scheduleJob(p *poller, observable Observable) {
	if !p.shuffle.toggle {
		go p.executeJob(observable)
		return
	}

	delayJob(p, observable)
}

func delayJob(p *poller, observable Observable) {
	job := func() {
		p.executeJob(observable)
	}

	p.shuffle.delayFunction(job)	
}

func (d *delay) delayFunction(fn func()) {
	time.AfterFunc(d.current, fn)
	d.current += d.gap
}