package middleware

import (
	"errors"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"neatly/pkg/e"
	"neatly/pkg/jwt"
	"neatly/pkg/logging"
	"net/http"
	"strings"
	"time"
)

const (
	authorizationHeader = "Authorization"
	userCtx             = "user_id"
)

func CorsMiddleware(router *gin.Engine) {
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:5173", "http://localhost:8084"}
	config.AllowHeaders = []string{"Content-Type"}
	config.ExposeHeaders = []string{"Content-Length"}
	config.AllowCredentials = true

	config.AllowOriginFunc = func(origin string) bool {
		return origin == "http://localhost:5173"
	}

	config.MaxAge = 12 * time.Hour

	router.Use(cors.New(config))
}

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
