package account

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"neatly/internal/mapper"
	"neatly/internal/model/account"
	"neatly/internal/service"
	"neatly/pkg/e"
	"neatly/pkg/logging"
	"net/http"
)

const (
	accountsURLGroup = "/accounts"
	registerURL      = "/register"
	loginURL         = "/login"
	apiURLGroup      = "/api"
	apiVersion       = "1"
)

type Handler struct {
	logger  logging.Logger
	service service.Account
	mapper  mapper.Account
}

func NewHandler(logger logging.Logger, service service.Account, mapper mapper.Account) *Handler {
	return &Handler{logger: logger, service: service, mapper: mapper}
}

func (h *Handler) Register(router *gin.Engine) {
	groupName := fmt.Sprintf("%v/v%v%v", apiURLGroup, apiVersion, accountsURLGroup)

	h.logger.Tracef("Register route: %v", groupName)

	group := router.Group(groupName)
	{
		group.POST(registerURL, h.register)
		group.POST(loginURL, h.login)
	}
}

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
		a   account.Account
	)

	if err = ctx.BindJSON(&dto); err != nil {
		h.logger.Error(err)
		e.NewErrorResponse(ctx, http.StatusBadRequest, err)
		return
	}

	a, err = h.mapper.MapRegisterAccountDTO(dto)
	if err != nil {
		h.logger.Error(err)
		e.NewErrorResponse(ctx, http.StatusInternalServerError, err)
		return
	}

	err = h.service.CreateAccount(&a)
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

	a := h.mapper.MapLogInAccountDTO(loginDto)

	token, err := h.service.GenerateJWT(&a)
	if err != nil {
		if errors.Is(err, &account.PasswordDoesNotMatchErr{}) {
			e.NewErrorResponse(ctx, http.StatusUnauthorized, err)
		}
		return
	}
	loginWithTokenDto := h.mapper.MapAccountWithTokenDTO(token, a)

	ctx.JSON(http.StatusOK, loginWithTokenDto)
}
