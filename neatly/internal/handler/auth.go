package handler

import (
	"github.com/gin-gonic/gin"
	"neatly/internal/e"
	"neatly/internal/model/auth"
	"net/http"
)

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

	ctx.JSON(http.StatusOK, u)
}

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
