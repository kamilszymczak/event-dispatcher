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

	// When scheduler runs the function twice every second
	t1 := Every(interval).Repeat(expectedRunCount).Do(increment)
	t1.Wait()
	
	// Then run count equals 2 and duration between two runs not more than run interval + threshold  
	timingsHelper(t, timings[:], interval, thresholdBetweenRuns)

	if runCount != expectedRunCount {
		t.Errorf("Expected %d runs, got %d", expectedRunCount, runCount)
	}
}

func timingsHelper(t *testing.T, arr []time.Time, interval time.Duration, threshold time.Duration) {
	for i, x := range arr[:len(arr) - 1] {
		runTime := arr[i+1].Sub(x)
		t.Logf("Timings difference %f", runTime.Seconds())

		if runTime >= interval + threshold {
			t.Logf("Expected duration between runs under %f, got %f", (interval + threshold).Seconds(), runTime.Seconds())
		}
	}
}