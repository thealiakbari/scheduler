package application

import (
	"encoding/json"
	"github.com/thealiakbari/scheduler/internal/domain"
	"io/ioutil"
	"net/http"
)

type TaskRequest struct {
	Priority int             `json:"priority"`
	Payload  json.RawMessage `json:"payload"`
}

type TaskResponse struct {
	ID     string        `json:"id"`
	Status domain.Status `json:"status"`
}

func (s *TaskService) HandleCreateTask(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var req TaskRequest
	if err := json.Unmarshal(body, &req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	pri := domain.Priority(req.Priority)
	task, err := s.CreateTask(pri, req.Payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	res := TaskResponse{
		ID:     task.ID,
		Status: task.Status,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

func (s *TaskService) HandleGetTaskStatus(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "missing task id", http.StatusBadRequest)
		return
	}
	status, err := s.GetTaskStatus(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": string(status)})
}
