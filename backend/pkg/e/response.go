package e

import (
	"errors"
	"github.com/gin-gonic/gin"
)

var (
	ClientNoteError    = errors.New("note does not exist or does not belong to user")
	ClientTagError     = errors.New("tag does not exist or does not belong to user")
	ClientAccountError = errors.New("username already exists")
	InternalDBError    = errors.New("database error occurred")
)

func NewErrorResponse(ctx *gin.Context, status int, err error) {
	er := ErrorResponse{
		Code:    status,
		Message: err.Error(),
	}

	ctx.AbortWithStatusJSON(status, er)
}

type ErrorResponse struct {
	Code    int    `json:"code" example:"400"`
	Message string `json:"message" example:"status bad request"`
}
