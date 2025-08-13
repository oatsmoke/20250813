package service

import (
	"context"
	"fmt"
	"time"

	"github.com/oatsmoke/20250813/internal/model"
)

type taskRepository interface {
	Set(key int64, value *model.Task)
	Get(key int64) (*model.Task, bool)
	GetAll() map[int64]*model.Task
}

type logger interface {
	Print(msg string)
}

const (
	starting = "starting"
	running  = "running"
	stopped  = "stopped"
)

type TaskService struct {
	srvCtx         context.Context
	taskRepository taskRepository
	loggerService  logger
}

func NewTaskService(ctx context.Context, taskRepository taskRepository, loggerService logger) *TaskService {
	return &TaskService{
		srvCtx:         ctx,
		taskRepository: taskRepository,
		loggerService:  loggerService,
	}
}

func (s *TaskService) Start(duration int64) string {
	key := time.Now().UnixNano()
	value := &model.Task{
		Status:  starting,
		Seconds: duration,
	}
	s.taskRepository.Set(key, value)

	msg := fmt.Sprintf("[%d] = %s", key, value.Status)
	s.loggerService.Print(msg)

	go s.run(s.srvCtx, key)

	return msg
}

func (s *TaskService) Get(key int64) *model.Task {
	value, _ := s.taskRepository.Get(key)

	return value
}

func (s *TaskService) GetAll(status string) map[int64]*model.Task {
	tasks := s.taskRepository.GetAll()

	if status != "" {
		filteredTasks := make(map[int64]*model.Task)
		for key, value := range tasks {
			if value.Status == status {
				filteredTasks[key] = value
			}
		}

		return filteredTasks
	}

	return tasks
}

func (s *TaskService) run(ctx context.Context, key int64) {
	timer := time.NewTicker(time.Second)
	defer timer.Stop()

	for {
		select {

		case <-ctx.Done():
			s.stop(key)
			return

		case <-timer.C:
			value, ok := s.taskRepository.Get(key)
			if !ok || value.Seconds < 1 {
				s.stop(key)
				return
			}

			v := &model.Task{
				Status:  running,
				Seconds: value.Seconds - 1,
			}
			s.taskRepository.Set(key, v)

			msg := fmt.Sprintf("[%d] = %d", key, v.Seconds)
			s.loggerService.Print(msg)
		}
	}
}

func (s *TaskService) stop(key int64) {
	v := &model.Task{
		Status:  stopped,
		Seconds: 0,
	}
	s.taskRepository.Set(key, v)

	msg := fmt.Sprintf("[%d] = %s", key, v.Status)
	s.loggerService.Print(msg)
}
