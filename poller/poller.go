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
	Response []byte
	Observable *Observable
}

type Observable struct {
	Address		string
	interval	*time.Duration	
}

type poller struct {
	apiUrl   		string
	interval 		time.Duration
	observables		[]Observable
	eventChan		chan Event
	errorHandler	func(err error) 
	jobsRunning 	sync.WaitGroup
	responseType	response.Responser
}

func New(url string, responseObject response.Responser) *poller {
	rate := config.GetConfig().FetchRate()

	return &poller{
		apiUrl:			url,
		interval: 		time.Duration(rate) * time.Millisecond,
		observables:	make([]Observable, 0),
		eventChan: 		make(chan Event),
		responseType: 	responseObject,	
	}
}

func (p *poller) Listen() <-chan Event {
	log.Print("Poller listener started")
	p.jobsRunning.Add(len(p.observables))

	for _, obs := range p.observables {
		go p.executeJob(obs)
	}

	go p.waitForJobsToComplete()
	return p.eventChan
}

func (t *poller) executeJob(observable Observable) {
	defer t.jobsRunning.Done()

	ticker := clockwork.NewRealClock().NewTicker(*observable.interval)
	job := scheduler.Every(ticker).Repeat(3).Do(t.poolData, observable)
	job.Wait()
}

func (t *poller) poolData(observable Observable) {
	log.Print("pooling data for " + observable.Address)

	//TODO - Make a httpClient struct for below to be able to mock with dummy api server
	res, err := http.Get(fmt.Sprintf("%s%s", t.apiUrl, observable.Address))
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	e := Event{
		Response: body,
		Observable: &observable,
	}

	t.eventChan <- e
}

func (t *poller) HandleEvent(event Event) (response.ResponseAccessor, error) {
	typedResponse := t.responseType.Unmarshal(event.Response)

	v, ok := typedResponse.(response.ResponseAccessor);
	if !ok {
		return nil, errors.New("Event response does not implement the Responser interface")
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