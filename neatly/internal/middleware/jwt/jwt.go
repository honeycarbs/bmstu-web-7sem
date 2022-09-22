package jwt

import (
	"context"
	"encoding/json"
	"github.com/cristalhq/jwt/v3"
	"neatly/internal/session"
	"neatly/pkg/logging"
	"net/http"
	"strings"
	"time"
)

type UserClaims struct {
	jwt.RegisteredClaims
	Email string `json:"email"`
}

func Middleware(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger := logging.GetLogger()
		authHeader := strings.Split(r.Header.Get("Authorization"), "Bearer ")
		if len(authHeader) != 2 {
			logger.Error("Malformed token")
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("malformed token"))
			return
		}

		logger.Debug("Create JWT verifier")
		jwtToken := authHeader[1]
		key := []byte(session.GetConfig().JWT.Secret)
		verifier, err := jwt.NewVerifierHS(jwt.HS256, key)
		if err != nil {
			unauthorized(w, err)
			return
		}

		logger.Debug("Parse and verify token")
		token, err := jwt.ParseAndVerifyString(jwtToken, verifier)
		if err != nil {
			unauthorized(w, err)
			return
		}

		var uc UserClaims
		err = json.Unmarshal(token.RawClaims(), &uc)
		if err != nil {
			unauthorized(w, err)
			return
		}
		if valid := uc.IsValidAt(time.Now()); !valid {
			logger.Error("Token has been expired.")
			unauthorized(w, err)
			return
		}

		ctx := context.WithValue(r.Context(), "user_id", uc.ID)
		h(w, r.WithContext(ctx))
	}
}

func unauthorized(w http.ResponseWriter, err error) {
	logging.GetLogger().Error(err)
	w.WriteHeader(http.StatusUnauthorized)
	w.Write([]byte("unauthorized"))
}
