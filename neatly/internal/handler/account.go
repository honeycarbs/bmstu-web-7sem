package handler

import (
	"errors"
	"github.com/gin-gonic/gin"
	"neatly/internal/model/account"
	"neatly/pkg/e"
	"net/http"
)

// @Summary Register
// @Tags register
// @Description create account
// @ID create-account
// @Accept  json
// @Produce  json
// @Param dto body account.RegisterAccountDTO true "account info"
// @Success 201 {string} string 1
// @Failure 500 {object} e.ErrorResponse
// @Failure default {object} e.ErrorResponse
// @Router /api/v1/accounts/register [post]
func (h *Handler) register(ctx *gin.Context) {
	var (
		err error
		dto account.RegisterAccountDTO
		u   account.Account
	)

	if err = ctx.BindJSON(&dto); err != nil {
		h.logger.Error(err)
		e.NewErrorResponse(ctx, http.StatusBadRequest, err)
		return
	}

	u, err = h.mappers.Account.MapRegisterAccountDTO(dto)
	if err != nil {
		h.logger.Error(err)
		e.NewErrorResponse(ctx, http.StatusInternalServerError, err)
		return
	}

	err = h.services.Authorisation.CreateAccount(&u)
	if err != nil {
		e.NewErrorResponse(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.Writer.WriteHeader(http.StatusCreated)
}

// @Summary Login
// @Tags account
// @Description login
// @ID login
// @Accept  json
// @Produce  json
// @Param dto body account.LoginAccountDTO true "credentials"
// @Success 200 {object} account.WithTokenDTO
// @Failure 500 {object} e.ErrorResponse
// @Failure default {object} e.ErrorResponse
// @Router /api/v1/accounts/login [post]
func (h *Handler) login(ctx *gin.Context) {
	var loginDto account.LoginAccountDTO

	if err := ctx.BindJSON(&loginDto); err != nil {
		h.logger.Error(err)
		e.NewErrorResponse(ctx, http.StatusBadRequest, err)
		return
	}

	a := h.mappers.Account.MapLogInAccountDTO(loginDto)
	//a = h.services.Authorisation.

	token, err := h.services.Authorisation.GenerateJWT(&a)
	if err != nil {
		if errors.Is(err, &e.PasswordDoesNotMatchErr{}) {
			e.NewErrorResponse(ctx, http.StatusUnauthorized, err)
		}
		return
	}
	loginWithTokenDto := h.mappers.Account.MapAccountWithTokenDTO(token, a)

	ctx.JSON(http.StatusOK, loginWithTokenDto)
}
