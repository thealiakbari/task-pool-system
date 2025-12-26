package pg

import (
	"context"

	"github.com/thealiakbari/task-pool-system/internal/domain/task/entity"
	"github.com/thealiakbari/task-pool-system/internal/ports/outbound/task"
	"github.com/thealiakbari/task-pool-system/pkg/common/db"
)

type TaskConfig struct {
	db db.DBWrapper
}

func NewTaskRepository(db db.DBWrapper) task.TaskRepository {
	return TaskConfig{
		db: db,
	}
}

func (u TaskConfig) Create(ctx context.Context, in entity.Task) (res entity.Task, err error) {
	err = db.GormConnection(ctx, u.db.DB).Save(&in).Error
	if err != nil {
		return entity.Task{}, err
	}

	return in, nil
}

func (u TaskConfig) Update(ctx context.Context, in entity.Task) (err error) {
	err = db.GormConnection(ctx, u.db.DB).Save(&in).Error
	if err != nil {
		return err
	}

	return nil
}

func (u TaskConfig) FindByIdOrEmpty(ctx context.Context, id string) (res entity.Task, err error) {
	err = db.GormConnection(ctx, u.db.DB).Model(&res).Order("created_at desc").Find(&res, "id = ?", id).Limit(1).Error
	if err != nil {
		return entity.Task{}, err
	}

	return res, nil
}

func (u TaskConfig) FindByIds(ctx context.Context, ids []string) (res []entity.Task, err error) {
	err = db.GormConnection(ctx, u.db.DB).Model(&res).Find(&res, "id IN (?)", ids).Error
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (u TaskConfig) Purge(ctx context.Context, id string) (err error) {
	err = db.GormConnection(ctx, u.db.DB).Unscoped().Delete(&entity.Task{}, "id = ?", id).Error
	if err != nil {
		return err
	}

	return nil
}

func (u TaskConfig) Delete(ctx context.Context, id string) (err error) {
	err = db.GormConnection(ctx, u.db.DB).Delete(&entity.Task{}, "id = ?", id).Error
	if err != nil {
		return err
	}

	return nil
}

func (u TaskConfig) FilterFind(ctx context.Context, query []any, order string, limit int, offset int) (res []entity.Task, err error) {
	err = db.GormConnection(ctx, u.db.DB).Model(&res).Order(order).
		Limit(limit).
		Offset(offset).
		Find(&res, query...).Error
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (u TaskConfig) FilterCount(ctx context.Context, query []any) (res int64, err error) {
	countQuery := db.GormConnection(ctx, u.db.DB).Model(&entity.Task{})
	if len(query) > 1 {
		countQuery = countQuery.Where(query[0], query[1:]...)
	} else if len(query) == 1 {
		countQuery = countQuery.Where(query[0])
	}

	err = countQuery.Count(&res).Error
	if err != nil {
		return 0, err
	}

	return res, nil
}
