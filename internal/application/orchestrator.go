package application

import (
	"fmt"
	"strconv"

	"github.com/hidnt/calc_go_yandex_2/pkg/calculation"
)

type Expression struct {
	ID      string               `json:"id"`
	Status  string               `json:"status"`
	Result  float64              `json:"result"`
	Actions []calculation.Action `json:"-"`
}

type Task struct {
	ID             string  `json:"id"`
	Arg1           float64 `json:"arg1"`
	Arg2           float64 `json:"arg2"`
	Operation      string  `json:"operation"`
	Operation_time int     `json:"operation_time"`
}

type Orchestrator struct {
	Expr                []Expression `json:"expressions"`
	ComputingAmount     int          `json:"-"`
	TimeAddition        int          `json:"-"`
	TimeSubstraction    int          `json:"-"`
	TimeMultiplications int          `json:"-"`
	TimeDivisions       int          `json:"-"`
}

func (o *Orchestrator) SetTime(cfg Config) {
	o.ComputingAmount = cfg.ComputingAmount
	o.TimeAddition = cfg.TimeAddition
	o.TimeSubstraction = cfg.TimeSubstraction
	o.TimeMultiplications = cfg.TimeMultiplications
	o.TimeDivisions = cfg.TimeDivisions
}

func NewOrchestrator() *Orchestrator {
	o := Orchestrator{
		Expr: []Expression{},
	}
	o.SetTime(*ConfigFromEnv())
	return &o
}

func (o *Orchestrator) ParceActionToTask(indexExpr int, indexTask int) (Task, error) {
	var t int = 0
	var arg1, arg2 float64

	switch o.Expr[indexExpr].Actions[indexTask].Operation {
	case "+":
		t = o.TimeAddition
	case "-":
		t = o.TimeSubstraction
	case "*":
		t = o.TimeMultiplications
	case "/":
		t = o.TimeDivisions
	}

	if o.Expr[indexExpr].Actions[indexTask].IdDepends[0] == "-1" {
		arg1 = o.Expr[indexExpr].Actions[indexTask].Arg1
	} else if d, _ := strconv.Atoi(o.Expr[indexExpr].Actions[indexTask].IdDepends[0]); o.Expr[indexExpr].Actions[d].Completed {
		arg1 = o.Expr[indexExpr].Actions[d].Result
	} else {
		return Task{}, fmt.Errorf("expression not available")
	}

	if o.Expr[indexExpr].Actions[indexTask].IdDepends[1] == "-1" {
		arg2 = o.Expr[indexExpr].Actions[indexTask].Arg2
	} else if d, _ := strconv.Atoi(o.Expr[indexExpr].Actions[indexTask].IdDepends[1]); o.Expr[indexExpr].Actions[d].Completed {
		arg2 = o.Expr[indexExpr].Actions[d].Result
	} else {
		return Task{}, fmt.Errorf("expression not available")
	}

	return Task{ID: o.Expr[indexExpr].Actions[indexTask].ID,
		Arg1:           arg1,
		Arg2:           arg2,
		Operation:      o.Expr[indexExpr].Actions[indexTask].Operation,
		Operation_time: t,
	}, nil
}
