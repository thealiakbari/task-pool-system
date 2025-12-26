package entity

import (
	"context"
	"time"

	"github.com/thealiakbari/task-pool-system/pkg/common/db"
)

type Status string

const (
	StatusPending   Status = "PENDING"
	StatusRunning   Status = "RUNNING"
	StatusCompleted Status = "COMPLETED"
	StatusFailed    Status = "FAILED"
)

type Task struct {
	db.UniversalModel
	Title       string `gorm:"column:title;type:varchar(255);not null" validate:"required"`
	Description string `gorm:"column:description;type:text;not null" validate:"required"`
	Status      Status
	Duration    time.Duration
}

func (u Task) Validate(ctx context.Context) error {
	return nil
}
