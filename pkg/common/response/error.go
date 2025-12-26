package response

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

func HandelError(ctx *gin.Context, err error) {
	status := -1
	var serviceError *Error
	ok := errors.As(err, &serviceError)
	if !ok {
		err = &Error{
			Cause:   err,
			Message: err.Error(),
			Class:   EUnknown,
		}
	} else {
		err = serviceError
	}

	var causeErr *Error
	errors.As(err, &causeErr)
	cause := causeErr.Cause

	if err.(*Error).Cause == nil {
		errors.As(err, &cause)
	}

	switch {
	case IsBadArg(err):
		status = http.StatusBadRequest
	case IsAccess(err):
		status = http.StatusForbidden
	case IsNotFound(err):
		status = http.StatusNotFound
	case IsConflict(err):
		status = http.StatusConflict
	case IsValidation(err):
		status = http.StatusUnprocessableEntity
	case IsUnauthorized(err):
		status = http.StatusUnauthorized
	default:
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, BaseResponse{
			Payload: nil,
			Meta: ErrResponse{
				Message: err.Error(),
				Causes:  cause,
				Code:    serviceError.ErrCode,
			},
		})
	}

	if status != -1 {
		ctx.AbortWithStatusJSON(status,
			BaseResponse{
				Payload: nil,
				Meta: ErrResponse{
					Message: err.Error(),
					Causes:  cause,
					Code:    serviceError.ErrCode,
				},
			},
		)
	}

	// nolinter
	_ = ctx.AbortWithError(status, err.(*Error))
}
