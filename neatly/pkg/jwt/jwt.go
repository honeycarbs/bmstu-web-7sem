package jwt

import (
	"encoding/json"
	"errors"
	"github.com/cristalhq/jwt/v3"
	"neatly/internal/model/account"
	"neatly/internal/session"
	"time"
)

const (
	tokenTTL = 12 * time.Hour
)

type UserClaims struct {
	jwt.RegisteredClaims
	UserID int
}

func GenerateAccessToken(u account.Account) (string, error) {
	key := []byte(session.GetConfig().JWT.Secret)

	signer, err := jwt.NewSignerHS(jwt.HS256, key)
	if err != nil {
		return "", err
	}
	builder := jwt.NewBuilder(signer)
	claims := UserClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenTTL)),
		},
		UserID: u.ID,
	}

	token, err := builder.Build(claims)
	if err != nil {
		return "", err
	}

	return token.String(), nil
}

func GetIdFromToken(token string) (int, error) {
	key := []byte(session.GetConfig().JWT.Secret)
	verifier, err := jwt.NewVerifierHS(jwt.HS256, key)
	if err != nil {
		return 0, err
	}

	tok, err := jwt.ParseAndVerifyString(token, verifier)
	if err != nil {
		return 0, err
	}

	var uc UserClaims
	err = json.Unmarshal(tok.RawClaims(), &uc)
	if err != nil {
		return 0, err
	}
	if valid := uc.IsValidAt(time.Now()); !valid {
		return 0, errors.New("token has been expired")
	}

	return uc.UserID, nil
}
