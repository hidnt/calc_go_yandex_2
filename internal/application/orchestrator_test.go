package application

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"
	"time"
)

type testApplication struct {
	name        string
	exprs       []string
	exceptedRes []Expression
}

type CalcRequest struct {
	Expression string `json:"expression"`
}

type CalcResponse struct {
	ID string `json:"id"`
}

func TestOrchestratorAgent(t *testing.T) {
	testCases := testApplication{
		name: "expressions",
		exceptedRes: []Expression{
			{
				ID:     "0",
				Status: "completed",
				Result: 891,
			},
			{
				ID:     "1",
				Status: "division by zero",
				Result: 0,
			},
		},
		exprs: []string{"1012+123-24*10-4", "1/0"},
	}

	app := New()
	go app.RunServer()

	for _, expr := range testCases.exprs {
		timer := time.NewTimer(3 * time.Second)
		<-timer.C
		requestBody, _ := json.Marshal(CalcRequest{Expression: expr})
		http.Post("http://localhost:8080/api/v1/calculate", "application/json", bytes.NewBuffer(requestBody))
	}

	count := 0
	for next := true; next; next = (count != 0) {
		count = 0
		for _, ex := range orchestrator.Expr {
			if ex.Status == "under consideration" {
				count++
			}
		}
	}

	ans, _ := http.Get("http://localhost:8080/api/v1/expressions")
	var data Orchestrator
	err := json.NewDecoder(ans.Body).Decode(&data)

	if err != nil {
		t.Fatalf("data can not be decode")
	}

	for i, _ := range testCases.exceptedRes {
		if testCases.exceptedRes[i].ID != data.Expr[i].ID {
			t.Fatalf("%s should be equal %s", data.Expr[i].ID, testCases.exceptedRes[i].ID)
		}
		if testCases.exceptedRes[i].Status != data.Expr[i].Status {
			t.Fatalf("%s should be equal %s", data.Expr[i].Status, testCases.exceptedRes[i].Status)
		}
		if testCases.exceptedRes[i].Result != data.Expr[i].Result {
			t.Fatalf("%f should be equal %f", data.Expr[i].Result, testCases.exceptedRes[i].Result)
		}
	}
}
