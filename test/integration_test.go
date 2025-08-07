package test

import (
	"encoding/json"
	"github.com/thealiakbari/scheduler/internal/application"
	"github.com/thealiakbari/scheduler/internal/domain"
	"github.com/thealiakbari/scheduler/internal/infrastructure/persistence"
	"github.com/thealiakbari/scheduler/internal/queue"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestTaskCreationAndProcessing(t *testing.T) {
	storage, _ := persistence.NewStorage()
	pq := queue.NewPriorityQueue()
	svc := application.NewTaskService(storage, pq)

	handler := http.HandlerFunc(svc.HandleCreateTask)
	payload := `{"priority":0, "payload": {"data": "test"}}`
	req := httptest.NewRequest("POST", "/tasks", strings.NewReader(payload))
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}
	var resp application.TaskResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatal(err)
	}
	if resp.Status != domain.StatusPending {
		t.Errorf("expected status pending, got %s", resp.Status)
	}

	if err := svc.ProcessNextTask(); err != nil {
		t.Fatal(err)
	}

	status, err := svc.GetTaskStatus(resp.ID)
	if err != nil {
		t.Fatal(err)
	}
	if status != domain.StatusCompleted {
		t.Errorf("expected status completed, got %s", status)
	}
}
