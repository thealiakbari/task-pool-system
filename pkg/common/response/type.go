package response

// ErrSwaggerResponse swagger generator //nolinter
type ErrSwaggerResponse struct {
	Payload map[string]interface{} `json:"payload"`
	Meta    struct {
		Code    int64  `json:"code"`
		Message string `json:"message"`
		Causes  []any  `json:"causes"`
	} `json:"meta"`
}

// ErrValidationSwaggerResponse swagger generator //nolinter
type ErrValidationSwaggerResponse struct {
	Payload map[string]interface{} `json:"payload"`
	Meta    struct {
		Code    int64  `json:"code"`
		Message string `json:"message"`
		Causes  []struct {
			Field   string `json:"field"`
			Message string `json:"message"`
		} `json:"causes"`
	} `json:"meta"`
}

type ErrResponse struct {
	Code    int64  `json:"code"`
	Message string `json:"message"`
	Causes  any    `json:"causes"`
}

type BaseResponse struct {
	Payload any         `json:"payload"`
	Meta    ErrResponse `json:"meta"`
}

type PaginationInfo struct {
	PageSize   int64 `json:"pageSize" form:"pageSize"`
	Page       int64 `json:"page" form:"page"`
	TotalItems int64 `json:"totalItems"`
}

type DefaultSort struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type ListResponse struct {
	Pagination  PaginationInfo `json:"pagination"`
	DefaultSort *DefaultSort   `json:"defaultSort,omitempty"`
	Items       any            `json:"items"`
}

func PaginationListResponse(items any, count int64, pageSize int64, page int64) ListResponse {
	if page == 0 {
		page = 1
	}
	return ListResponse{
		Pagination: PaginationInfo{
			PageSize:   pageSize,
			TotalItems: count,
			Page:       page,
		},
		Items: items,
	}
}

func PaginationAndSortListResponse(items any, count int64, pageSize int64, page int64, sortBy, sortType string) ListResponse {
	if page == 0 {
		page = 1
	}
	response := ListResponse{
		Pagination: PaginationInfo{
			PageSize:   pageSize,
			TotalItems: count,
			Page:       page,
		},
		Items: items,
	}

	if sortBy != "" && sortType != "" {
		response.DefaultSort = &DefaultSort{
			Key:   sortBy,
			Value: sortType,
		}
	}

	return response
}
