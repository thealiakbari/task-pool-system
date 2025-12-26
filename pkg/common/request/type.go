package request

import (
	"context"
	"time"

	"github.com/thealiakbari/task-pool-system/pkg/common/validation"
)

type SortType = string

const (
	SortTypeASC  SortType = "ASC"  // ASC
	SortTypeDESC SortType = "DESC" /* DESC */
)

// swagger:model SortSpec
type SortSpec struct {
	SortType *SortType `json:"sortType" form:"sortType" default:"DESC" required:"false" enums:"ASC,DESC"`
}

// Pagination used to RAW DTO, in transport layer. it will be transformed to `Portion` for
// calling service methods
type Pagination struct {
	// Starts from 1
	Page     int `json:"page" form:"page" required:"false" minimum:"1" default:"1"`
	PageSize int `json:"pageSize" form:"pageSize"  required:"false" default:"12"`
}

func (p Pagination) Validate(ctx context.Context) error {
	return validation.Validate(ctx, p)
}

// Portion transformed version of `Pagination`, for service and repository layers
type Portion struct {
	Offset int
	Limit  int
}

type NumberRange struct {
	From *int64
	To   *int64
}

type DateRange struct {
	From *time.Time
	To   *time.Time
}
