package scheduler

import (
	"reflect"
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

// Runs function in a goroutine and returns the job.
// Optional args parameter to provide paramters to provided function
func (j *Job) Do(fn any, args ...any) *Job {
	f := reflect.ValueOf(fn)

	// TODO improve validation when working with variadic function
	// if len(args) != f.Type().NumIn() {
	// 	fmt.Println("Invalid number of arguments")
	// 	return nil
	// }

	in := make([]reflect.Value, len(args))
	for i, param := range args {
		in[i] = reflect.ValueOf(param)
	}

	ticker := time.NewTicker(j.interval)

	go func() {
		defer ticker.Stop()
		i := 0
		j.running = true

		L: for{
			i++
			f.Call(in)

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

// func (j *Job) Do(f func()) *Job {
// 	ticker := time.NewTicker(j.interval)

// 	go func() {
// 		defer ticker.Stop()
// 		i := 0
// 		j.running = true

// 		L: for{
// 			i++
// 			f()

// 			select {
// 				case <- ticker.C:
// 					if(i >= j.repeats){
// 						close(j.quit)
// 						break L
// 					}
// 					continue

// 				case <- j.quit:
// 					break L
// 			}
// 		}
// 		j.running = false
// 	}()
// 	return j
// }

func (j *Job) Stop() {
	close(j.quit)
}

// Blocking until job finishes
func (j *Job) Wait() {
	<-j.quit
}

// Checks if job running indefinitely 
func (j *Job) Indefinite() bool {
	return j.repeats == -1
}


// func scheduleWithParams(f func(), params []any) (bool, error) {
// 	fun := reflect.TypeOf(f)
// 	paramsCount := fun.NumIn()
// 	if len(params) != paramsCount {
// 		return _, ErrParamsNotAdapted
// 	}

// }