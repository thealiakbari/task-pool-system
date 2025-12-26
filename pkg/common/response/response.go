package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreatedResponse(ctx *gin.Context, body any) {
	ctx.JSON(http.StatusCreated, BaseResponse{
		Payload: body,
		Meta: ErrResponse{
			Causes: []any{},
		},
	})
}

func OKResponse(ctx *gin.Context, body any) {
	ctx.JSON(http.StatusOK, BaseResponse{
		Payload: body,
		Meta: ErrResponse{
			Causes: []any{},
		},
	})
}

func NoContentResponse(ctx *gin.Context) {
	ctx.AbortWithStatus(http.StatusNoContent)
}

func NotFoundResponse(ctx *gin.Context) {
	ctx.AbortWithStatus(http.StatusNotFound)
}
