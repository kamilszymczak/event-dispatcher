package scheduler

import (
	"reflect"
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
	
	timingsHelper(t, timings[:], interval, thresholdBetweenRuns)

	if runCount != expectedRunCount {
		t.Errorf("Expected %d runs, got %d", expectedRunCount, runCount)
	}
}

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

func TestValidArgumentsValid(t *testing.T) {
	foo := func(str string) {
		return 
	}

	want := true
	got := validArguments(reflect.ValueOf(foo).Type(), "hello")

	if got != want {
		t.Errorf("Expected %t got %t", want, got)
	} 
}

func TestValidArgumentsInvalid(t *testing.T) {
	foo := func(str string) {
		return 
	}

	want := false
	got := validArguments(reflect.ValueOf(foo).Type(), "hello", "world")

	if got != want {
		t.Errorf("Expected %t got %t", want, got)
	} 
}

func TestValidArgumentsValidOnlyVariadic(t *testing.T) {
	foo := func(args ...any) bool {
		return len(args) > 0
	}

	want := true
	got := validArguments(reflect.ValueOf(foo).Type(), 1, 2, 3)

	if got != want {
		t.Errorf("Expected %t got %t", want, got)
	} 
}

func TestValidArgumentsValidNoArgsVariadic(t *testing.T) {
	foo := func(args ...any) bool {
		return len(args) > 0
	}

	want := true
	got := validArguments(reflect.ValueOf(foo).Type())

	if got != want {
		t.Errorf("Expected %t got %t", want, got)
	} 
}

func TestValidArgumentsValidArgAndEmptyVariadic(t *testing.T) {
	foo := func(first *int, args ...int) {
		var sum int = *first
		for _, a := range args {
			sum = sum + a
		}
		*first = sum
	}

	want := true
	got := validArguments(reflect.ValueOf(foo).Type(), new(int))

	if got != want {
		t.Errorf("Expected %t got %t", want, got)
	} 
}

func TestValidArgumentsValidArgAndVariadicArg(t *testing.T) {
	foo := func(first *int, args ...int) {
		var sum int = *first
		for _, a := range args {
			sum = sum + a
		}
		*first = sum
	}

	want := true
	got := validArguments(reflect.ValueOf(foo).Type(), new(int), 1, 2)

	if got != want {
		t.Errorf("Expected %t got %t", want, got)
	} 
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