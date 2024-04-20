package scheduler

import (
	"context"
	"fmt"
	"log"
	"reflect"

	"github.com/mixer/clock"
)

type Scheduler interface {
	Every(interval clock.Ticker)
	Repeat(repeats int)
	Do(fn any, args ...any) *Job
	Wait()
}

type Job struct {
	Ctx context.Context
	Cancel context.CancelFunc
	interval clock.Ticker
	jobFunc func()
	repeats int
	running bool
	// quit chan struct{}
}

func newJob(interval clock.Ticker) *Job {
	context, cancel := context.WithCancel(context.Background())
	j := &Job{
		Ctx: context,
		Cancel: cancel,
		interval: interval,
		jobFunc: nil,
		repeats: -1,
	}
	return j
}

func Every(interval clock.Ticker) *Job {
	j := newJob(interval)
	return j
}

func (j *Job) Repeat(repeats int) *Job {
	j.repeats = repeats
	return j
}

func allArgumentsOfSameType(args ...any) bool {
	m := map[reflect.Type]bool{}
	for _ , x := range args {
		m[reflect.TypeOf(x)] = true
	}
	return len(m) == 1	
}

func validArguments(fArgs reflect.Type, args ...any) bool {
	// fmt.Printf("function args: %v given args: %v Variadic: %v\n", fArgs.NumIn(), args, fArgs.IsVariadic())
	if fArgs.IsVariadic() && fArgs.In(0).Elem().String() == "interface {}"{
		if len(args) > 1 && reflect.TypeOf(args[0]) != reflect.TypeOf(args[1]) {
			return false
		}
		return true
	}
	if fArgs.IsVariadic() && fArgs.In(0) != reflect.TypeOf(args[0]) {
		return false
	}
	if fArgs.NumIn() > 0 && fArgs.In(0) != reflect.TypeOf(args[0]) {
		return false
	} 
	if fArgs.IsVariadic() && len(args) > 0 && fArgs.In(0) == reflect.TypeOf(args[0]) {
		return true
	}
	if !fArgs.IsVariadic() && fArgs.NumIn() == len(args) {
		return true
	}
	if fArgs.IsVariadic() && fArgs.NumIn() == 1 {
		return true
	}
	if fArgs.IsVariadic() && len(args) >= fArgs.NumIn()-1 {
		return true
	}
	return false
}

// Runs function in a goroutine and returns the job.
// Optional args parameter to provide paramters to provided function
func (j *Job) Do(fn any, args ...any) *Job {
	f := reflect.ValueOf(fn)

	if !validArguments(f.Type(), args...) {
		fmt.Println("Invalid number of arguments")
		return nil
	}

	in := make([]reflect.Value, len(args))
	for i, param := range args {
		in[i] = reflect.ValueOf(param)
	}

	go func() {
		defer j.interval.Stop()
		defer func(){j.running = false}()
		i := 0
		j.running = true

		L: for{
			i++
			f.Call(in)

			if(i >= j.repeats){
				j.Cancel()
				return
			}

			select {
				case <- j.interval.Chan():
					log.Print("received from ticker ch")
					continue
				case <- j.Ctx.Done():
					break L
			}
		}
	}()
	return j
}

func (j *Job) Stop() {
	j.Cancel()
}

// Blocking until job finishes
func (j *Job) Wait() {
	<-j.Ctx.Done()
}

// Checks if job running indefinitely 
func (j *Job) Indefinite() bool {
	return j.repeats == -1
}
