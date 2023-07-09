package scheduler

import (
	"testing"
	"time"
)

func TestSchedulerRunTwiceEverySecond(t *testing.T) {
	const (
		interval time.Duration = 1 * time.Second
		expectedRunCount int = 2
		thresholdBetweenRuns time.Duration = 20 * time.Millisecond
	)
	runCount := 0
	var timings [expectedRunCount]time.Time

	// Given a function to execute
	increment := func(){
		timings[runCount] = time.Now()
		runCount++
	}

	t1 := Every(interval).Repeat(expectedRunCount).Do(increment)
	t1.Wait()
	
	// Then run count equals 2 and duration between two runs not more than run interval + threshold  
	timingsHelper(t, timings[:], interval, thresholdBetweenRuns)

	if runCount != expectedRunCount {
		t.Errorf("Expected %d runs, got %d", expectedRunCount, runCount)
	}
}

// func TestSchedulerIndefiniteRepeats(t *testing.T) {
// 	const (
// 		interval time.Duration = 1 * time.Second
// 		thresholdBetweenRuns time.Duration = 20 * time.Millisecond
// 	)

// 	runCount := 0
// 	var timings []time.Time

// 	increment := func(){
// 		timings[runCount] = time.Now()
// 		runCount++
// 	}

// 	t1 := Every(interval).Do(increment)

// }

func TestSchedulerRunOnceFuncWithParams(t *testing.T) {
	got := new(int)
	want := 4
	addTwo := func (result *int, a, b int) {
		*result = a+b
	}

	t1 := Every(time.Second).Repeat(1).Do(addTwo, got, 1, 3)
	t1.Wait()

	if *got != want {
		t.Errorf("Expected %d got %d", want, *got)
	} 
}

func TestSchedulerRunThreeTimesFuncWithParams(t *testing.T) {
	got := new(int)
	want := 3
	increment := func(result *int) {
		(*result)++
	}

	t1 := Every(500 * time.Millisecond).Repeat(3).Do(increment, got)
	t1.Wait()

	if *got != want {
		t.Errorf("Expected %d got %d", want, *got)
	} 
}

func TestSchedulerWithVariadicFunc(t *testing.T) {
	got := new(string)
	want := "HelloWorldHelloWorld"
	append := func (result *string, args ...string) {
		for _, x := range args {
			*result = (*result) + x
		}
	}

	t1 := Every(500 * time.Millisecond).Repeat(2).Do(append, got, "Hello", "World")
	t1.Wait()

	if *got != want {
		t.Errorf("Expected %s got %s", want, *got)
	} 
}

func TestSchedulerWithInvalidArgs(t *testing.T) {

}


func timingsHelper(t *testing.T, arr []time.Time, interval time.Duration, threshold time.Duration) {
	for i, x := range arr[:len(arr) - 1] {
		runTime := arr[i+1].Sub(x)
		t.Logf("Timings difference %f", runTime.Seconds())

		if runTime >= interval + threshold {
			t.Errorf("Expected duration between runs under %f, got %f", (interval + threshold).Seconds(), runTime.Seconds())
		}
	}
}