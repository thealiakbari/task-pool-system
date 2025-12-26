package dto

import (
	"github.com/thealiakbari/task-pool-system/internal/domain/task/entity"
	"time"

	"github.com/google/uuid"
)

type Task struct {
	Id          uuid.UUID     `json:"id"`
	Title       string        `json:"title"`
	Description string        `json:"description"`
	Status      entity.Status `json:"status"`
	Duration    time.Duration `json:"duration"`
	CreatedAt   time.Time     `json:"createdAt"`
	UpdatedAt   time.Time     `json:"updatedAt"`
}
