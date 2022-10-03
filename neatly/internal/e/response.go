package e

import "github.com/gin-gonic/gin"

type errorResponse struct {
	Message string `json:"message"`
}

func NewErrorResponse(ctx *gin.Context, code int, msg string) {
	ctx.AbortWithStatusJSON(code, errorResponse{msg})
}
