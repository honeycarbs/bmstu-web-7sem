package middleware

import (
	"errors"
	"github.com/gin-gonic/gin"
	"neatly/pkg/e"
	"neatly/pkg/jwt"
	"neatly/pkg/logging"
	"net/http"
	"strings"
)

const (
	authorizationHeader = "Authorization"
	userCtx             = "user_id"
)

func Authenticate(ctx *gin.Context) {
	header := ctx.GetHeader(authorizationHeader)
	if header == "" {
		e.NewErrorResponse(ctx, http.StatusUnauthorized, errors.New("unauthorized"))
		return
	}
	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 {
		e.NewErrorResponse(ctx, http.StatusUnauthorized, errors.New("malformed token"))
		return
	}
	userID, err := jwt.GetIdFromToken(headerParts[1])
	if err != nil {
		e.NewErrorResponse(ctx, http.StatusUnauthorized, err)
	}

	logging.GetLogger().Info("authorized")

	ctx.Set(userCtx, userID)
}

func GetUserID(ctx *gin.Context) (int, error) {
	id, ok := ctx.Get(userCtx)
	if !ok {
		return 0, errors.New("can't get authorization parameters")
	}

	idNum, ok := id.(int)
	if !ok {
		return 0, errors.New("can't get authorization params")
	}

	return idNum, nil
}
