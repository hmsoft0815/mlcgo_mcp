package db

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/mlcmcp/task-manager/internal/models"
)

type TaskStore struct {
	tasks    map[string]*models.Task
	path     string
	planMode bool
	mu       sync.RWMutex
}

func NewTaskStore(path string) (*TaskStore, error) {
	s := &TaskStore{
		tasks: make(map[string]*models.Task),
		path:  path,
	}
	if err := s.load(); err != nil {
		return nil, err
	}
	return s, nil
}

func (s *TaskStore) load() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	data, err := os.ReadFile(s.path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	if len(data) == 0 {
		return nil
	}

	return json.Unmarshal(data, &s.tasks)
}

func (s *TaskStore) save() error {
	data, err := json.MarshalIndent(s.tasks, "", "  ")
	if err != nil {
		return err
	}

	// Ensure directory exists
	dir := filepath.Dir(s.path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	return os.WriteFile(s.path, data, 0644)
}

func (s *TaskStore) Create(subject, description, activeForm string, metadata map[string]interface{}) *models.Task {
	s.mu.Lock()
	defer s.mu.Unlock()

	id := fmt.Sprintf("%d", len(s.tasks)+1)
	now := time.Now()
	task := &models.Task{
		ID:          id,
		Subject:     subject,
		Description: description,
		ActiveForm:  activeForm,
		Status:      models.StatusPending,
		Metadata:    metadata,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	s.tasks[id] = task
	s.save()
	return task
}

func (s *TaskStore) Get(id string) (*models.Task, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	t, ok := s.tasks[id]
	return t, ok
}

func (s *TaskStore) List() []*models.Task {
	s.mu.RLock()
	defer s.mu.RUnlock()

	list := make([]*models.Task, 0, len(s.tasks))
	for _, t := range s.tasks {
		if t.Status != models.StatusDeleted {
			list = append(list, t)
		}
	}
	return list
}

func (s *TaskStore) Update(id string, update func(*models.Task)) (*models.Task, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	t, ok := s.tasks[id]
	if !ok {
		return nil, fmt.Errorf("task %s not found", id)
	}

	update(t)
	t.UpdatedAt = time.Now()
	s.save()
	return t, nil
}

func (s *TaskStore) SetPlanMode(active bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.planMode = active
}

func (s *TaskStore) IsPlanMode() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.planMode
}
