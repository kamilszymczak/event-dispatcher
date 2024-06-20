package poller

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/jonboulle/clockwork"
	"github.com/kamilszymczak/event-dispatcher/config"
	"github.com/kamilszymczak/event-dispatcher/response"
	"github.com/kamilszymczak/event-dispatcher/scheduler"
)

type Event struct {
	Response   []byte
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
	jobs         []*scheduler.Job
	dispatchFunc dispatchFunc
	shuffle 	 *delay
}

func New(url string, responseObject response.Responser) *poller {
	rate := config.GetConfig().FetchRate()

	return &poller{
		apiUrl:       url,
		interval:     time.Duration(rate) * time.Millisecond,
		observables:  make([]Observable, 0),
		eventChan:    make(chan Event),
		responseType: responseObject,
		shuffle:	  newDelay(),
	}
}

func (p *poller) Listen() <-chan Event {
	log.Print("Poller listener started")
	p.jobsRunning.Add(len(p.observables))

	for _, obs := range p.observables {
		scheduleJob(p, obs)
	}

	go p.waitForJobsToComplete()
	return p.eventChan
}

func (t *poller) executeJob(observable Observable) {
	defer t.jobsRunning.Done()

	ticker := clockwork.NewRealClock().NewTicker(*observable.interval)
	job := scheduler.Every(ticker).Do(t.poolData, observable)
	t.jobs = append(t.jobs, job)
	job.Wait()
}

func (t *poller) fetchData(observable Observable) ([]byte, error) {
	log.Printf("Pooling data for %s", observable.Address)

	resp, err := http.Get(fmt.Sprintf("%s%s", t.apiUrl, observable.Address))
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("fetching data unsuccessful, response status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	return body, nil
}

func buildEvent(observable Observable, body []byte) Event {
	return Event{
		Response:   body,
		Observable: &observable,
	}
}

func (t *poller) poolData(observable Observable) {
	data, err := t.fetchData(observable)
	if err != nil {
		log.Print(err.Error())
		return
	}

	event := buildEvent(observable, data)

	if t.dispatchFunc == nil {
		t.eventChan <- event
		return
	}

	response, _ := t.HandleEvent(event)
	if t.dispatchFunc(response) {
		t.eventChan <- event
	}
}

func (t *poller) HandleEvent(event Event) (response.ResponseAccessor, error) {
	typedResponse, err := t.responseType.Unmarshal(event.Response)

	if err != nil {
		return nil, err
	}

	v, ok := typedResponse.(response.ResponseAccessor)
	if !ok {
		return nil, errors.New("event response does not implement the Responser interface")
	}

	return v, nil
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
}

func (p *poller) stopAll() {
	for _, j := range p.jobs {
		j.Stop()
	}
}

func (p *poller) SetDispatchFunc(fn dispatchFunc) {
	p.dispatchFunc = fn
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