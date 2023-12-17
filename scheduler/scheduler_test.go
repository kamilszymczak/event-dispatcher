package scheduler

import (
	"reflect"
	"testing"
	"time"

	"github.com/mixer/clock"
)

func TestSchedulerWithIntegerParameters(t *testing.T) {
	testCases := []struct {
		name 		string
		interval 	time.Duration
		repeat		int
		want		int
		instanceVar	int
		fn 			any
		fnParams	[]int
	}{
		{
			name: "scheduler run twice",
			interval: 1 * time.Second,
			repeat: 2,
			want: 2,
			instanceVar: 0,
			fn: func(count *int){
				*count++
			},
		},
		{
			name: "scheduler run once function with two parameters",
			interval: 1 * time.Second,
			repeat: 1,
			want: 4,
			instanceVar: 0,
			fn: func(result *int, a, b int) {
				*result = a+b
			},
			fnParams: []int{1, 3},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			fc := clock.NewMockClock()
			ticker := fc.NewTicker(tc.interval)

			var t1 *Job
			if len(tc.fnParams) > 0 {
				arguments := []interface{}{&tc.instanceVar}
				for _, x := range tc.fnParams {
					arguments = append(arguments, x)
				}
				t1 = Every(ticker).Repeat(tc.repeat).Do(tc.fn, arguments...)
			} else {
				t1 = Every(ticker).Repeat(tc.repeat).Do(tc.fn, &tc.instanceVar)
			}

			fc.AddTime(tc.interval)
			t1.Wait()

			if tc.instanceVar != tc.want {
				t.Errorf("Expected %d runs, got %d", tc.want, tc.instanceVar)
			}
		})
	}
}

func TestSchedulerWithStringParameters(t *testing.T) {
	testCases := []struct {
		name 		string
		interval 	time.Duration
		repeat		int
		want		string
		instanceVar	string
		fn 			any
		fnParams	[]string
	}{
		{
			name: "scheduler run once function with variadic parameters",
			interval: 500 * time.Millisecond,
			repeat: 2,
			want: "HelloWorldHelloWorld",
			instanceVar: "",
			fn: func (result *string, args ...string) {
				for _, x := range args {
					*result = (*result) + x
				}
			},
			fnParams: []string{"Hello", "World"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			fc := clock.NewMockClock()
			ticker := fc.NewTicker(tc.interval)

			var t1 *Job
			if len(tc.fnParams) > 0 {
				arguments := []interface{}{&tc.instanceVar}
				for _, x := range tc.fnParams {
					arguments = append(arguments, x)
				}
				t1 = Every(ticker).Repeat(tc.repeat).Do(tc.fn, arguments...)
			} else {
				t1 = Every(ticker).Repeat(tc.repeat).Do(tc.fn, &tc.instanceVar)
			}

			fc.AddTime(tc.interval)
			t1.Wait()

			if tc.instanceVar != tc.want {
				t.Errorf("Expected %s runs, got %s", tc.want, tc.instanceVar)
			}
		})
	}
}

func TestValidFuncArguments(t *testing.T) {
	foo := func(str string) {
		return 
	}

	want := true
	got := validArguments(reflect.ValueOf(foo).Type(), "hello")

	if got != want {
		t.Errorf("Expected %t got %t", want, got)
	} 
}

func TestInvalidFuncArguments(t *testing.T) {
	foo := func(str string) {
		return 
	}

	want := false
	got := validArguments(reflect.ValueOf(foo).Type(), "hello", "world")

	if got != want {
		t.Errorf("Expected %t got %t", want, got)
	} 
}

func TestValidArgumentsVariadicFunc(t *testing.T) {
	foo := func(args ...any) bool {
		return len(args) > 0
	}

	want := true
	got := validArguments(reflect.ValueOf(foo).Type(), 1, 2, 3)

	if got != want {
		t.Errorf("Expected %t got %t", want, got)
	} 
}

func TestValidArgumentsNoArgsVariadicFunc(t *testing.T) {
	foo := func(args ...any) bool {
		return len(args) > 0
	}

	want := true
	got := validArguments(reflect.ValueOf(foo).Type())

	if got != want {
		t.Errorf("Expected %t got %t", want, got)
	} 
}

func TestValidArgumentsValidArgAndEmptyVariadicFunc(t *testing.T) {
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

func TestValidArgumentsValidVariadicFuncArguments(t *testing.T) {
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