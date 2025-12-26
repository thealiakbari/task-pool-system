package task

import (
	"context"
	"errors"

	"github.com/thealiakbari/task-pool-system/internal/domain/task/entity"
	taskInterface "github.com/thealiakbari/task-pool-system/internal/ports/inbound/task"
	"github.com/thealiakbari/task-pool-system/internal/ports/outbound/task"
	"github.com/thealiakbari/task-pool-system/pkg/common/logger"
)

type TaskConfig struct {
	Logger   logger.Logger
	TaskRepo task.TaskRepository
}

type taskService struct {
	TaskConfig
}

func NewTaskService(config TaskConfig) taskInterface.TaskService {
	u := taskService{config}
	u.Logger = config.Logger.ForService(u)
	return u
}

func (u taskService) Create(ctx context.Context, req entity.Task) (res entity.Task, err error) {
	if err = req.Validate(ctx); err != nil {
		u.Logger.Warnf(ctx, "validation error:%v", err)
		return entity.Task{}, err
	}

	taskEntity, err := u.TaskRepo.Create(ctx, req)
	if err != nil {
		u.Logger.Errorf(ctx, "Cannot create task item: %v", err)
		return entity.Task{}, err
	}

	return taskEntity, nil
}

func (u taskService) Update(ctx context.Context, req entity.Task) (res entity.Task, err error) {
	if err = req.Validate(ctx); err != nil {
		u.Logger.Warnf(ctx, "validation error:%v", err)
		return entity.Task{}, err
	}

	err = u.TaskRepo.Update(ctx, req)
	if err != nil {
		return entity.Task{}, err
	}

	return req, nil
}

func (u taskService) GetByIdOrEmpty(ctx context.Context, id string) (res entity.Task, err error) {
	if id == "" {
		err = errors.New("id must not be empty")
		return entity.Task{}, err
	}

	taskEntity, err := u.TaskRepo.FindByIdOrEmpty(ctx, id)
	if err != nil {
		return entity.Task{}, err
	}

	return taskEntity, nil
}

func (u taskService) Purge(ctx context.Context, id string) (err error) {
	if id == "" {
		err := errors.New("id must not be empty")
		return err
	}

	err = u.TaskRepo.Purge(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

func (u taskService) Delete(ctx context.Context, id string) (err error) {
	if id == "" {
		err := errors.New("id must not be empty")
		return err
	}

	err = u.TaskRepo.Delete(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

func (u taskService) GetByPhoneNumberOrEmpty(ctx context.Context, phoneNumber string) (res entity.Task, err error) {
	if phoneNumber == "" {
		err := errors.New("phoneNumber must not be empty")
		return entity.Task{}, err
	}

	taskEntity, err := u.TaskRepo.FindByPhoneNumberOrEmpty(ctx, phoneNumber)
	if err != nil {
		return entity.Task{}, err
	}

	return taskEntity, nil
}
