package handler

import (
	"errors"
	"github.com/gin-gonic/gin"
	"neatly/internal/e"
	"neatly/pkg/jwt"
	"net/http"
	"strings"
)

const (
	authorizationHeader = "Authorization"
	userCtx             = "user_id"
)

func (h *Handler) userIdentity(ctx *gin.Context) {
	header := ctx.GetHeader(authorizationHeader)
	if header == "" {
		e.NewErrorResponse(ctx, http.StatusUnauthorized, "unauthorized")
		return
	}
	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 {
		e.NewErrorResponse(ctx, http.StatusUnauthorized, "malformed token")
		return
	}
	userID, err := jwt.GetIdFromToken(headerParts[1])
	if err != nil {
		e.NewErrorResponse(ctx, http.StatusUnauthorized, err.Error())
	}

	ctx.Set(userCtx, userID)
}

func (h *Handler) getUserID(ctx *gin.Context) (int, error) {
	id, ok := ctx.Get(userCtx)
	if !ok {
		e.NewErrorResponse(ctx, http.StatusInternalServerError, "error umm")
		return 0, errors.New("error umm")
	}

	idNum, ok := id.(int)
	if !ok {
		e.NewErrorResponse(ctx, http.StatusInternalServerError, "error umm")
		return 0, errors.New("error umm")
	}

	return idNum, nil
}
