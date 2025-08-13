package repository

import (
	"sync"

	"github.com/oatsmoke/20250813/internal/model"
)

type TaskRepository struct {
	data map[int64]*model.Task
	mu   sync.RWMutex
}

func NewTaskRepository() *TaskRepository {
	return &TaskRepository{
		data: make(map[int64]*model.Task),
	}
}

func (r *TaskRepository) Set(key int64, value *model.Task) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.data[key] = value
}

func (r *TaskRepository) Get(key int64) (*model.Task, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	value, ok := r.data[key]

	return value, ok
}

func (r *TaskRepository) GetAll() map[int64]*model.Task {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return r.data
}
