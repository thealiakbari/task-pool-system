package dto

import (
	"context"

	"github.com/thealiakbari/task-pool-system/pkg/common/validation"
)

type CreateTaskRequest struct {
	Title       string `json:"title" validate:"required"`
	Description string `json:"description" validate:"required"`
}

func (c CreateTaskRequest) Validate(ctx context.Context) error {
	return validation.Validate(ctx, c)
}
