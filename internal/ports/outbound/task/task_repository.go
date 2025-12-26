package task

import (
	"context"

	"github.com/thealiakbari/task-pool-system/internal/domain/task/entity"
)

type TaskRepository interface {
	Create(ctx context.Context, in entity.Task) (res entity.Task, err error)
	Update(ctx context.Context, in entity.Task) (err error)
	FindByIds(ctx context.Context, ids []string) (res []entity.Task, err error)
	FindByIdOrEmpty(ctx context.Context, id string) (res entity.Task, err error)
	Purge(ctx context.Context, id string) (err error)
	Delete(ctx context.Context, id string) (err error)
	FilterFind(ctx context.Context, query []any, order string, limit int, offset int) (res []entity.Task, err error)
	FilterCount(ctx context.Context, query []any) (res int64, err error)
}
