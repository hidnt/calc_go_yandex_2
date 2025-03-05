package calculation

import (
	"fmt"
	"strconv"
	"strings"
)

type Action struct {
	ID           string
	Arg1         float64
	Arg2         float64
	Result       float64
	Operation    string
	IdDepends    []string
	Completed    bool
	NowCalculate bool
}

func Calc(expression string, curID int) ([]Action, error) {
	var parts []string
	var curPart string
	for _, char := range expression {
		if char == ' ' {
			continue
		}
		if char == '-' && (curPart == "" || curPart == "(") {
			curPart += string(char)
			continue
		}
		if strings.ContainsAny(string(char), "0123456789.") {
			curPart += string(char)
			continue
		}
		if curPart != "" {
			parts = append(parts, curPart)
			curPart = ""
		}
		if strings.ContainsAny(string(char), "+*/()-") {
			parts = append(parts, string(char))
			continue
		}
		return nil, ErrUnknownOp
	}

	if curPart != "" {
		parts = append(parts, curPart)
	}

	var nums []string
	var operators []string
	var actions []Action
	curAction := 0

	priority := map[string]int{
		"+": 1,
		"-": 1,
		"*": 2,
		"/": 2,
		"(": 0,
	}

	calculate := func() error {
		if len(operators) == 0 {
			return ErrNotEnoughtOp
		}
		if len(nums) < 2 && len(operators) >= 1 {
			err := ErrNotEnoughtNums
			return err
		}

		depends := []string{}
		var left, right float64

		if num, err := strconv.ParseFloat(nums[len(nums)-2], 64); err == nil {
			left = num
			depends = append(depends, "-1")
		} else {
			left = 0
			depends = append(depends, nums[len(nums)-2][1:])
		}
		if num, err := strconv.ParseFloat(nums[len(nums)-1], 64); err == nil {
			right = num
			depends = append(depends, "-1")
		} else {
			right = 0
			depends = append(depends, nums[len(nums)-1][1:])
		}
		nums = nums[:len(nums)-2]

		operator := operators[len(operators)-1]
		operators = operators[:len(operators)-1]

		actions = append(actions, Action{
			ID:        fmt.Sprintf("%d_%d", curID, curAction),
			Arg1:      left,
			Arg2:      right,
			Result:    0,
			Operation: operator,
			IdDepends: depends,
		})

		nums = append(nums, fmt.Sprintf("d%d", curAction))
		curAction++

		return nil
	}

	for _, part := range parts {
		if _, err := strconv.ParseFloat(part, 64); err == nil {
			nums = append(nums, part)
			continue
		}
		if part == "(" {
			operators = append(operators, part)
			continue
		}
		if part == ")" {
			for len(operators) > 0 && operators[len(operators)-1] != "(" {
				if err := calculate(); err != nil {
					return nil, err
				}
			}
			if len(operators) == 0 {
				return nil, ErrIncorrectPriorOp
			}
			operators = operators[:len(operators)-1]
		} else {
			for len(operators) > 0 && priority[operators[len(operators)-1]] >= priority[part] {
				if err := calculate(); err != nil {
					return nil, err
				}
			}
			operators = append(operators, part)
		}
	}

	for len(operators) > 0 {
		if err := calculate(); err != nil {
			return nil, err
		}
	}

	if len(nums) != 1 {
		return nil, ErrCalc
	}

	return actions, nil
}
