package calculation

import (
	"slices"
	"testing"
)

type testCase struct {
	name        string
	expression  string
	id          int
	exceptedRes []Action
	wantError   bool
}

func TestCalc(t *testing.T) {
	testCases := []testCase{
		{
			name:       "correct expression",
			expression: "1012+123-24*10-4",
			id:         0,
			exceptedRes: []Action{
				{
					ID:           "0_0",
					Arg1:         1012,
					Arg2:         123,
					Result:       0,
					Operation:    "+",
					IdDepends:    []string{"-1", "-1"},
					Completed:    false,
					NowCalculate: false,
				},
				{
					ID:           "0_1",
					Arg1:         24,
					Arg2:         10,
					Result:       0,
					Operation:    "*",
					IdDepends:    []string{"-1", "-1"},
					Completed:    false,
					NowCalculate: false,
				},
				{
					ID:           "0_2",
					Arg1:         0,
					Arg2:         0,
					Result:       0,
					Operation:    "-",
					IdDepends:    []string{"0", "1"},
					Completed:    false,
					NowCalculate: false,
				},
				{
					ID:           "0_3",
					Arg1:         0,
					Arg2:         4,
					Result:       0,
					Operation:    "-",
					IdDepends:    []string{"2", "-1"},
					Completed:    false,
					NowCalculate: false,
				},
			},
			wantError: false,
		},
		{
			name:       "correct expression",
			expression: "1+1",
			id:         0,
			exceptedRes: []Action{
				{
					ID:           "0_0",
					Arg1:         1,
					Arg2:         1,
					Result:       0,
					Operation:    "+",
					IdDepends:    []string{"-1", "-1"},
					Completed:    false,
					NowCalculate: false,
				},
			},
			wantError: false,
		},
		{
			name:        "incorrect expression",
			expression:  "1238)",
			exceptedRes: []Action{},
			wantError:   true,
		},
		{
			name:        "incorrect expression 2",
			expression:  "124+2-",
			exceptedRes: []Action{},
			wantError:   true,
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			ans, err := Calc(testCase.expression, testCase.id)
			if testCase.wantError {
				if err == nil {
					t.Fatalf("Excepted an err")
				}
			} else {
				if err != nil {
					t.Fatalf("Successful case is %s, but returns error: %s", testCase.expression, err.Error())
				}
				isEqual := true
				for i := 0; i < len(ans); i++ {
					if testCase.exceptedRes[i].ID != ans[i].ID ||
						testCase.exceptedRes[i].Arg1 != ans[i].Arg1 ||
						testCase.exceptedRes[i].Arg2 != ans[i].Arg2 ||
						testCase.exceptedRes[i].Result != ans[i].Result ||
						testCase.exceptedRes[i].Operation != ans[i].Operation ||
						testCase.exceptedRes[i].Completed != ans[i].Completed ||
						testCase.exceptedRes[i].NowCalculate != ans[i].NowCalculate ||
						!slices.Equal(testCase.exceptedRes[i].IdDepends, ans[i].IdDepends) {
						isEqual = false
						break
					}
				}
				if !isEqual {
					t.Fatalf("%v should be equal %v", ans, testCase.exceptedRes)
				}
			}
		})
	}
}
