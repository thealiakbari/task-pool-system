package utiles

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/thealiakbari/task-pool-system/pkg/common/request"
)

const ConvertTimeLayout = "2006-01-02 15:04:05"

func ConvertUnixIntoDate(unix int64, separator string) string {
	t := time.Unix(unix, 0).UTC() // UTC returns t with the location set to UTC.
	return t.Format(fmt.Sprintf("2006%s01%s02", separator, separator))
}

func ConvertStringToDateTime(dateTime string) time.Time {
	date, err := time.Parse(ConvertTimeLayout, dateTime)
	if err != nil {
		panic(err)
	}

	return date
}

func ConvertStringIntoFloat64(str string) float64 {
	floatNum, err := strconv.ParseFloat(str, 64)
	if err != nil {
		return 0.0
	}
	return floatNum
}

func ConvertStringDateIntoUnix(date string, separator string) (int64, error) {
	layout := fmt.Sprintf("2006%s01%s02", separator, separator)
	t, err := time.Parse(layout, date)
	if err != nil {
		return 0, err
	}
	return t.Unix(), nil
}

func ConvertStringToInt64(str string) int64 {
	i, _ := strconv.ParseInt(str, 10, 64)
	return i
}

// func PaginationConvertor(page string, offset string) response.Pagination {
// 	finalOffset := int((ConvertStringToInt64(page) - 1) * ConvertStringToInt64(offset))
// 	return request.Pagination{PageSize: finalOffset, Page: int(ConvertStringToInt64(offset))}
// }

func PaginationNormalizer(pagination request.Pagination, ctx context.Context) (request.Pagination, error) {
	err := pagination.Validate(ctx)
	if err != nil {
		return request.Pagination{}, err
	}
	// Handling the page and default of 0 was move to `PaginationToPortion`
	// What is passed to pagination, will remain to reiterative to caller as response
	if pagination.PageSize <= 0 {
		pagination.PageSize = 12
	}
	return pagination, nil
}

func PaginationNormalizerFromParams(page, pageSize string, ctx context.Context) (request.Pagination, error) {
	pag := request.Pagination{
		Page:     int(ConvertStringToInt64(page)),
		PageSize: int(ConvertStringToInt64(pageSize)),
	}
	// Handling the page and default of 0 was move to `PaginationToPortion`
	// What is passed to pagination, will remain to reiterative to caller as response
	return PaginationNormalizer(pag, ctx)
}

func PaginationToPortion(pagination request.Pagination) request.Portion {
	if pagination.Page >= 1 {
		pagination.Page -= 1
	} else {
		pagination.Page = 0
	}

	return request.Portion{
		Limit:  pagination.PageSize,
		Offset: pagination.Page * pagination.PageSize,
	}
}

func ConvertToPointerBool(bool bool) *bool {
	return &bool
}

func ConvertToUUID(strs []string) ([]uuid.UUID, error) {
	out := make([]uuid.UUID, len(strs))

	for i, str := range strs {
		val, err := uuid.Parse(str)
		if err != nil {
			return nil, fmt.Errorf("error parsing UUID from string '%s': %v", str, err)
		}
		out[i] = val
	}

	return out, nil
}

func Ptr[T any](in T) *T {
	return &in
}
