package handler

import (
	"github.com/gin-gonic/gin"
	"neatly/internal/e"
	"neatly/internal/model/auth"
	"net/http"
)

// @Summary Register
// @Tags register
// @Description create account
// @ID create-account
// @Accept  json
// @Produce  json
// @Param dto body auth.RegisterAccountDTO true "account info"
// @Success 201
// @Failure 500 {object} e.ErrorResponse
// @Failure default {object} e.ErrorResponse
// @Router /api/v1/auth/register [post]
func (h *Handler) register(ctx *gin.Context) {
	var (
		err error
		dto auth.RegisterAccountDTO
		u   auth.Account
	)

	if err = ctx.BindJSON(&dto); err != nil {
		h.logger.Error(err)
		e.NewErrorResponse(ctx, http.StatusBadRequest, err)
		return
	}

	u, err = h.mappers.Auth.MapRegisterAccountDTO(dto)
	if err != nil {
		h.logger.Error(err)
		e.NewErrorResponse(ctx, http.StatusInternalServerError, err)
		return
	}

	err = h.services.Authorisation.CreateUser(&u)
	if err != nil {
		e.NewErrorResponse(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.Writer.WriteHeader(http.StatusCreated)
}

// @Summary Login
// @Tags auth
// @Description login
// @ID login
// @Accept  json
// @Produce  json
// @Param dto body auth.LoginAccountDTO true "credentials"
// @Success 200 {object} auth.JwtDTO
// @Failure 500 {object} e.ErrorResponse
// @Failure default {object} e.ErrorResponse
// @Router /api/v1/auth/login [post]
func (h *Handler) login(ctx *gin.Context) {
	var dto auth.LoginAccountDTO

	if err := ctx.BindJSON(&dto); err != nil {
		h.logger.Error(err)
		e.NewErrorResponse(ctx, http.StatusBadRequest, err)
		return
	}

	u := h.mappers.Auth.MapLogInAccountDTO(dto)

	token, err := h.services.Authorisation.GenerateJWT(&u)
	if err != nil {
		e.NewErrorResponse(ctx, http.StatusInternalServerError, err)
		return
	}
	jdto := h.mappers.Auth.MapJwtDTO(token)

	ctx.JSON(http.StatusOK, jdto)
}
