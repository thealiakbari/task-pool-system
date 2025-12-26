package task

import (
	"context"

	"github.com/thealiakbari/task-pool-system/internal/domain/task/entity"
)

type TaskService interface {
	Create(ctx context.Context, entity entity.Task) (res entity.Task, err error)
	Update(ctx context.Context, entity entity.Task) (res entity.Task, err error)
	GetByIdOrEmpty(ctx context.Context, id string) (res entity.Task, err error)
	Delete(ctx context.Context, id string) (err error)
	Purge(ctx context.Context, id string) (err error)
	GetByPhoneNumberOrEmpty(ctx context.Context, phoneNumber string) (res entity.Task, err error)
}
