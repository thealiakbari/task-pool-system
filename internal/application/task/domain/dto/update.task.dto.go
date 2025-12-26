package dto

import (
	"context"

	"github.com/thealiakbari/task-pool-system/pkg/common/validation"
)

type UpdateTaskRequest struct {
	Title       string `json:"title" validate:"required"`
	Description string `json:"description" validate:"required"`
}

func (u UpdateTaskRequest) Validate(ctx context.Context) error {
	return validation.Validate(ctx, u)
}
