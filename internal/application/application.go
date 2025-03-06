package application

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/hidnt/calc_go_yandex_2/pkg/calculation"
)

type Config struct {
	Addr                string
	ComputingAmount     int
	TimeAddition        int
	TimeSubstraction    int
	TimeMultiplications int
	TimeDivisions       int
}

type Application struct {
	config *Config
	orch   *Orchestrator
}

func New() *Application {
	return &Application{
		config: ConfigFromEnv(),
		orch:   orchestrator,
	}
}

type RequestCalc struct {
	Expression string `json:"expression"`
}

type ResponseCalc struct {
	ID string `json:"id"`
}

type ResponseExprID struct {
	Expression Expression `json:"expression"`
}

var curId = 0
var orchestrator *Orchestrator = NewOrchestrator()

func ConfigFromEnv() *Config {
	config := new(Config)
	config.Addr = os.Getenv("PORT")
	if config.Addr == "" {
		config.Addr = "8080"
	}
	var err error
	config.ComputingAmount, err = strconv.Atoi(os.Getenv("COMPUTING_POWER"))
	if err != nil || config.ComputingAmount <= 0 {
		config.ComputingAmount = 1
	}
	config.TimeAddition, err = strconv.Atoi(os.Getenv("TIME_ADDITION_MS"))
	if err != nil || config.TimeAddition < 0 {
		config.TimeAddition = 1000
	}
	config.TimeSubstraction, err = strconv.Atoi(os.Getenv("TIME_SUBTRACTION_MS"))
	if err != nil || config.TimeSubstraction < 0 {
		config.TimeSubstraction = 1000
	}
	config.TimeMultiplications, err = strconv.Atoi(os.Getenv("TIME_MULTIPLICATIONS_MS"))
	if err != nil || config.TimeMultiplications < 0 {
		config.TimeMultiplications = 1000
	}
	config.TimeDivisions, err = strconv.Atoi(os.Getenv("TIME_DIVISIONS_MS"))
	if err != nil || config.TimeDivisions < 0 {
		config.TimeDivisions = 1000
	}
	return config
}

func CalcHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	request := new(RequestCalc)
	defer r.Body.Close()
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ResponseCalc{ID: fmt.Sprintf("%d", curId)})
		return
	}
	request.Expression = strings.TrimSpace(request.Expression)
	actions, err := calculation.Calc(request.Expression, curId)
	expr := Expression{
		ID:      fmt.Sprintf("%d", curId),
		Result:  0,
		Actions: actions,
	}
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(ResponseCalc{ID: fmt.Sprintf("%d", curId)})
		expr.Status = fmt.Sprintf("%v", err)
	} else {
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(ResponseCalc{ID: fmt.Sprintf("%d", curId)})
		expr.Status = "under consideration"
	}
	curId++
	orchestrator.Expr = append(orchestrator.Expr, expr)
}

func GetExprHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(orchestrator)
}

func GetExprHandlerID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id := r.URL.Path[len("/api/v1/expressions/"):]
	if id == "" {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(orchestrator.Expr)
		return
	}
	n, err := strconv.Atoi(id[1:])
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(orchestrator.Expr)
		return
	} else if (len(orchestrator.Expr) > n) && (n >= 0) {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(ResponseExprID{orchestrator.Expr[n]})
	} else {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(orchestrator.Expr)
	}
}

func TaskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		GetTaskHandler(w, r)
	case "POST":
		CompleteTaskHandler(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func GetTaskHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	exprAvailable := false
	taskBreak := false
	for i := 0; i < len(orchestrator.Expr); i++ {
		if orchestrator.Expr[i].Status == "under consideration" {
			exprAvailable = true
		}
	}
	if !exprAvailable {
		w.WriteHeader(http.StatusNotFound)
	} else {
		for exprID := 0; exprID < len(orchestrator.Expr); exprID++ {
			for actionID := 0; actionID < len(orchestrator.Expr[exprID].Actions); actionID++ {
				if orchestrator.Expr[exprID].Actions[actionID].Completed || orchestrator.Expr[exprID].Actions[actionID].NowCalculate {
					continue
				}
				task, err := orchestrator.ParceActionToTask(exprID, actionID)
				if err != nil {
					continue
				}
				json.NewEncoder(w).Encode(task)
				taskBreak = true
				orchestrator.Expr[exprID].Actions[actionID].NowCalculate = true
				break
			}
			if taskBreak {
				break
			}
		}
	}
}

func CompleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	request := new(RequestTaskData)
	defer r.Body.Close()
	err := json.NewDecoder(r.Body).Decode(&request)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	ids := strings.Split(request.ID, "_")
	exprID, _ := strconv.Atoi(ids[0])
	actionID, _ := strconv.Atoi(ids[1])

	log.Printf("Expression %d, task %d was completed", exprID, actionID)

	if request.Msg != "" {
		w.WriteHeader(http.StatusUnprocessableEntity)
		orchestrator.Expr[exprID].Status = "division by zero"
	} else {
		orchestrator.Expr[exprID].Actions[actionID].Completed = true
		orchestrator.Expr[exprID].Actions[actionID].Result = request.Res
	}

	completed := 0
	for _, a := range orchestrator.Expr[exprID].Actions {
		if a.Completed {
			completed++
		}
	}

	if completed == len(orchestrator.Expr[exprID].Actions) {
		orchestrator.Expr[exprID].Status = "completed"
		orchestrator.Expr[exprID].Result = request.Res
	}
}

func (a *Application) RunServer() error {
	for i := 0; i < a.config.ComputingAmount; i++ {
		go agent()
	}

	http.Handle("/api/v1/", http.StripPrefix("/api/v1", http.FileServer(http.Dir("./static."))))

	http.HandleFunc("/api/v1/calculate", CalcHandler)
	http.HandleFunc("/api/v1/expressions", GetExprHandler)
	http.HandleFunc("/api/v1/expressions/", GetExprHandlerID)
	http.HandleFunc("/internal/task", TaskHandler)

	log.Printf("Запущен сервер на :%s", a.config.Addr)
	return http.ListenAndServe(":"+a.config.Addr, nil)
}
