package application

import (
	"testing"
)

type testCase struct {
	name        string
	task        Task
	exceptedRes float64
	wantError   bool
}

func TestCalc(t *testing.T) {
	testCases := []testCase{
		{
			name: "correct 1",
			task: Task{
				ID:             "0",
				Arg1:           10,
				Arg2:           10,
				Operation:      "*",
				Operation_time: 0,
			},
			exceptedRes: 100,
			wantError:   false,
		},
		{
			name: "correct 2",
			task: Task{
				ID:             "0",
				Arg1:           10,
				Arg2:           10,
				Operation:      "-",
				Operation_time: 0,
			},
			exceptedRes: 0,
			wantError:   false,
		},
		{
			name: "correct 1",
			task: Task{
				ID:             "0",
				Arg1:           10,
				Arg2:           10,
				Operation:      "/",
				Operation_time: 0,
			},
			exceptedRes: 1,
			wantError:   false,
		},
		{
			name: "div by zero",
			task: Task{
				ID:             "0",
				Arg1:           10,
				Arg2:           0,
				Operation:      "/",
				Operation_time: 0,
			},
			exceptedRes: 0,
			wantError:   true,
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			ans, err := solveTask(&testCase.task)
			if testCase.wantError {
				if err == nil {
					t.Fatalf("Excepted an err")
				}
			} else {
				if err != nil {
					t.Fatalf("Successful case is %f, but returns error: %s", testCase.exceptedRes, err.Error())
				}
				if ans != testCase.exceptedRes {
					t.Fatalf("%v should be equal %v", ans, testCase.exceptedRes)
				}
			}
		})
	}
}
