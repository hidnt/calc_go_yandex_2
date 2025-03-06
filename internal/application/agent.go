package application

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type RequestTaskData struct {
	ID  string  `json:"ID"`
	Msg string  `json:"msg,omitempty"`
	Res float64 `json:"result"`
}

func agent() {
	for {
		task, err := getTask()
		if err == nil {
			res, err2 := solveTask(task)
			data := RequestTaskData{
				ID:  task.ID,
				Res: res,
			}
			if err2 != nil {
				data.Msg = err2.Error()
			}
			sendResult(data)
		} else {
			time.Sleep(time.Second)
		}
	}
}

func getTask() (*Task, error) {
	res, err := http.Get("http://localhost:8080/internal/task")
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("no available tasks")
	}

	var data Task

	if err = json.NewDecoder(res.Body).Decode(&data); err != nil {
		return nil, err
	}
	return &data, nil
}

func solveTask(task *Task) (float64, error) {
	timer := time.NewTimer(time.Duration(task.Operation_time) * time.Millisecond)
	<-timer.C
	switch task.Operation {
	case "+":
		return task.Arg1 + task.Arg2, nil
	case "-":
		return task.Arg1 - task.Arg2, nil
	case "*":
		return task.Arg1 * task.Arg2, nil
	case "/":
		if task.Arg2 == 0 {
			return 0, fmt.Errorf("division by zero")
		}
		return task.Arg1 / task.Arg2, nil
	default:
		return 0, fmt.Errorf("unknown operation")
	}
}

func sendResult(data RequestTaskData) error {
	jsonData, _ := json.Marshal(data)
	res, err := http.Post("http://localhost:8080/internal/task", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to send")
	}
	res.Body.Close()
	return nil
}
