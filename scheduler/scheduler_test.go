package scheduler

import (
	"log"
	"reflect"
	"testing"
	"time"

	"github.com/mixer/clock"
)

func TestSchedulerRunTwiceEverySecond(t *testing.T) {
	const (
		interval time.Duration = 1 * time.Second
		expectedRunCount int = 2
	)
	runCount := 0
	fc := clock.NewMockClock()
	ticker := fc.NewTicker(interval)

	// Given a function to execute
	increment := func(){
		runCount++
	}

	t1 := Every(ticker).Repeat(expectedRunCount).Do(increment)
	fc.AddTime(interval)
	t1.Wait()

	if runCount != expectedRunCount {
		t.Errorf("Expected %d runs, got %d", expectedRunCount, runCount)
	}
}

func TestSchedulerRunOnceFuncWithParams(t *testing.T) {
	const (
		interval time.Duration = 1 * time.Second
		want int = 4
	)
	got := new(int)
	fc := clock.NewMockClock()
	ticker := fc.NewTicker(interval)

	addTwo := func (result *int, a, b int) {
		*result = a+b
	}

	t1 := Every(ticker).Repeat(1).Do(addTwo, got, 1, 3)
	fc.AddTime(interval)
	t1.Wait()

	if *got != want {
		t.Errorf("Expected %d got %d", want, *got)
	} 
}

func TestSchedulerRunThreeTimesFuncWithParams(t *testing.T) {
	const (
		interval time.Duration = 1 * time.Second
		want int = 3
	)
	got := new(int)
	// fc := clockwork.NewFakeClock()
	// ticker := fc.NewTicker(interval)
	fc := clock.NewMockClock()
	ticker := fc.NewTicker(interval)

	increment := func(result *int) {
		log.Print("function executed")
		(*result)++
	}

	t1 := Every(ticker).Repeat(3).Do(increment, got)
	fc.AddTime(2 * time.Second)
	t1.Wait()

	if *got != want {
		t.Errorf("Expected %d got %d", want, *got)
	} 
}

func TestSchedulerWithVariadicFunc(t *testing.T) {
	const (
		interval time.Duration = 500 * time.Millisecond
		want string = "HelloWorldHelloWorld"
	)
	got := new(string)

	fc := clock.NewMockClock()
	ticker := fc.NewTicker(interval)

	append := func (result *string, args ...string) {
		for _, x := range args {
			*result = (*result) + x
		}
	}

	t1 := Every(ticker).Repeat(2).Do(append, got, "Hello", "World")
	fc.AddTime(1 * time.Second)
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