package transform

import (
	"github.com/google/uuid"
	"github.com/thealiakbari/task-pool-system/internal/application/task/domain/dto"
	"github.com/thealiakbari/task-pool-system/internal/domain/task/entity"
	"time"
)

func CreateTaskRequestToEntity(in dto.CreateTaskRequest) entity.Task {
	out := entity.Task{
		Title:       in.Title,
		Description: in.Description,
		Status:      entity.StatusPending,
		Duration:    time.Duration(1+time.Now().Unix()%5) * time.Second,
	}

	return out
}

func UpdateTaskRequestToEntity(in dto.UpdateTaskRequest, id string) (out entity.Task, err error) {
	out = entity.Task{
		Title:       in.Title,
		Description: in.Description,
	}

	idUUID, err := uuid.Parse(id)
	if err != nil {
		return out, err
	}

	out.Id = idUUID
	return out, nil
}

func TaskEntityToTaskDto(in entity.Task) dto.Task {
	return dto.Task{
		Id:          in.Id,
		Title:       in.Title,
		Description: in.Description,
		Status:      in.Status,
		Duration:    in.Duration,
		CreatedAt:   in.CreatedAt,
		UpdatedAt:   in.UpdatedAt,
	}
}

func TasksEntityToTasksDto(in []entity.Task) []dto.Task {
	items := make([]dto.Task, 0, len(in))
	for _, v := range in {
		items = append(items, TaskEntityToTaskDto(v))
	}

	return items
}
