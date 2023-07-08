package scheduler

import (
	// "reflect"
	"fmt"
	"time"
)

type Scheduler interface {
	Schedule(f func(), interval time.Duration) *Job
	Every(interval time.Duration)
	Repeat(repeats int)
}

type Job struct {
	interval time.Duration
	jobFunc func()
	repeats int
	running bool
	quit chan bool
}

func newJob(interval time.Duration) *Job {
	j := &Job{
		interval: interval,
		jobFunc: nil,
		repeats: -1,
		quit: make(chan bool, 1),
	}
	return j
}

func Every(interval time.Duration) *Job {
	j := newJob(interval)
	return j
}

func (j *Job) Repeat(repeats int) *Job {
	j.repeats = repeats
	return j
}

func (j *Job) Do(f func()) *Job {
	ticker := time.NewTicker(j.interval)

	go func() {
		defer ticker.Stop()
		i := 0
		j.running = true

		L: for{
			i++
			f()

			select {
				case <- ticker.C:
					if(i >= j.repeats){
						close(j.quit)
						break L
					}
					continue

				case <- j.quit:
					break L
			}
		}
		j.running = false
	}()
	return j
}

func (j *Job) Stop() {
	close(j.quit)
}

// Blocking until job finishes
func (j *Job) Wait() {
	<-j.quit
}

func (j *Job) 


// func scheduleWithParams(f func(), params []any) (bool, error) {
// 	fun := reflect.TypeOf(f)
// 	paramsCount := fun.NumIn()
// 	if len(params) != paramsCount {
// 		return _, ErrParamsNotAdapted
// 	}

// }