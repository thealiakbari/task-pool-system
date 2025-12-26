package dto

import (
	"context"

	"github.com/thealiakbari/task-pool-system/pkg/common/request"
)

type GetTaskRequest struct {
	Ids   []string `form:"ids"`
	Title []string `json:"titles"`

	request.Pagination `json:"-"`
}

func (g GetTaskRequest) Validate(ctx context.Context) error {
	return nil
}
