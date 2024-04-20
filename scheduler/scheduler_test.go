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

func TestSchedulerCancelJob(t *testing.T) {
	const (
		interval time.Duration = 200 * time.Millisecond
		repeat int = 3
		want string = "HelloHello"
	)

	fn := func (result *string, args ...string) {
		for _, x := range args {
			*result = (*result) + x
		}
	}
	fnParams := []string{"Hello"}
	instanceVar := ""

	ticker := clock.C.NewTicker(interval)

	arguments := []interface{}{&instanceVar}
	for _, x := range fnParams {
		arguments = append(arguments, x)
	}
	job := Every(ticker).Repeat(repeat).Do(fn, arguments...)

	// cancel job after task expected to have executed twice
	time.Sleep(250 * time.Millisecond)
	job.Stop()
	job.Wait()

	if instanceVar != want {
		t.Errorf("Expected %s runs, got %s", want, instanceVar)
	}
}

func TestSchedulerFunctionArgumentsValidation(t *testing.T) {
	testCases := []struct {
		name 		string
		want		bool
		instanceVar	string
		fn 			any
		fnParams	[]interface{}
	}{
		{
			name: "Valid function arguments",
			want: true,
			fn: func(str string) { 
			},
			fnParams: []interface{}{"hello"},
		},
		{
			name: "Valid function arguments for variadic function",
			want: true,
			fn: func(args ...any) bool {
				return len(args) > 0
			},
			fnParams: []interface{}{"hello", "world"},
		},
		{
			name: "Invalid function arguments for variadic function providing different parameters to function",
			want: false,
			fn: func(args ...any) bool {
				return len(args) > 0
			},
			fnParams: []interface{}{"hello", 1},
		},
		{
			name: "Invalid function arguments for function taking in a single value while providing a variadic parameter",
			want: false,
			fn: func(str string) { 
			},
			fnParams: []interface{}{"hello", "world"},
		},
		{
			name: "Valid arguments for variadic function when providing an empty slice",
			want: true,
			fn: func(args ...any) bool {
				return len(args) > 0
			},
			fnParams: []interface{}{},
		},
		{
			name: "Valid arguments for function with normal parameter and variadic parameter",
			want: true,
			fn: func(first *int, args ...int) {
			},
			fnParams: []interface{}{new(int), 1, 2},
		},
		{
			name: "Valid arguments for function with normal parameter and empty variadic parameter",
			want: true,
			fn: func(first *int, args ...int) {
			},
			fnParams: []interface{}{new(int)},
		},
		{
			name: "Invalid arguments for function taking a pointer parameter and variadic parameter. 1st parameter not provided. Type mismatch",
			want: false,
			fn: func(first *int, args ...int) {
			},
			fnParams: []interface{}{1, 2},
		},
		{
			name: "Invalid arguments for function taking an int variadic parameter while providing a string",
			want: false,
			fn: func(args ...int) {
			},
			fnParams: []interface{}{"a", "b"},
		},
		{
			name: "Invalid arguments for function taking a pointer and variadic parameter when not providing the pointer",
			want: false,
			fn: func(first *int, args ...int) {
			},
			fnParams: []interface{}{"a", "b"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := validArguments(reflect.ValueOf(tc.fn).Type(), tc.fnParams...)

			if got != tc.want {
				t.Errorf("Expected %t got %t", tc.want, got)
			}
		})
	}
}

// func TestSchedulerThrowsError(t *testing.T) {
// 	const (
// 		interval time.Duration = 1 * time.Second
// 		expectedRunCount int = 2
// 	)
// 	runCount := 0
// 	fc := clock.NewMockClock()
// 	ticker := fc.NewTicker(interval)

// 	// Given a function to execute
// 	increment := func(count *int){
// 		*count++
// 	}

// 	t1 := Every(ticker).Repeat(expectedRunCount).Do(increment, runCount)
// 	fc.AddTime(interval)
// 	t1.Wait()

// 	if runCount != expectedRunCount {
// 		t.Errorf("Expected %d runs, got %d", expectedRunCount, runCount)
// 	}
// }