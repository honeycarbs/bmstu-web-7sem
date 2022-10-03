package handler

import (
	"github.com/gin-gonic/gin"
	"neatly/internal/e"
	"neatly/internal/model/user"
	"net/http"
)

func (h *Handler) register(ctx *gin.Context) {
	var (
		err error
		dto user.RegisterUserDTO
	)

	if err = ctx.BindJSON(&dto); err != nil {
		h.logger.Error(err)
		e.NewErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	u, err := h.userMapper.MapUserRegisterDTO(dto)
	if err != nil {
		h.logger.Error(err)
		e.NewErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	err = h.services.Authorisation.CreateUser(&u)
	if err != nil {
		e.NewErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, map[string]interface{}{
		"id": u.ID,
	})
}

func (h *Handler) login(ctx *gin.Context) {
	var dto user.LoginUserDTO

	if err := ctx.BindJSON(&dto); err != nil {
		h.logger.Error(err)
		e.NewErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	u := h.userMapper.MapUserLogInUserDTO(dto)

	token, err := h.services.Authorisation.GenerateJWT(&u)
	if err != nil {
		e.NewErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, map[string]interface{}{
		"token": token,
	})

}
