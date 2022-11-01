package account

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"neatly/internal/handlers/middleware"
	"neatly/internal/mapper"
	"neatly/internal/model"
	"neatly/internal/model/dto"
	"neatly/internal/service"
	"neatly/pkg/e"
	"neatly/pkg/logging"
	"net/http"
	"strconv"
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
	service service.AccountServiceImpl
	mapper  mapper.AccountMapper
}

func NewHandler(logger logging.Logger, service service.AccountServiceImpl, mapper mapper.AccountMapper) *Handler {
	return &Handler{logger: logger, service: service, mapper: mapper}
}

func (h *Handler) Register(router *gin.Engine) {
	groupName := fmt.Sprintf("%v/v%v%v", apiURLGroup, apiVersion, accountsURLGroup)

	h.logger.Tracef("Register route: %v", groupName)

	auth := router.Group(groupName)
	{
		auth.POST(registerURL, h.register)
		auth.POST(loginURL, h.login)
	}

	accounts := router.Group(groupName, middleware.Authenticate)
	{
		accounts.GET("/:id", h.getAccount)
	}
}

// @Summary Register
// @Tags account
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
		in  dto.RegisterAccountDTO
		a   model.Account
	)

	if err = ctx.BindJSON(&in); err != nil {
		h.logger.Error(err)
		e.NewErrorResponse(ctx, http.StatusBadRequest, err)
		return
	}

	a, err = h.mapper.MapRegisterAccountDTO(in)
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

	ctx.JSON(http.StatusCreated, fmt.Sprintf(
		"%v/v%v%v/%v",
		apiURLGroup,
		apiVersion,
		accountsURLGroup,
		a.ID,
	))
}

// @Summary getAccount
// @Security ApiKeyAuth
// @Tags account
// @Description get account
// @ID get-account
// @Accept  json
// @Produce  json
// @Param id   path  string  true  "id"
// @Success 200 {object} account.GetAccountDTO
// @Failure 500 {object} e.ErrorResponse
// @Failure default {object} e.ErrorResponse
// @Router /api/v1/accounts/:id [get]
func (h *Handler) getAccount(ctx *gin.Context) {
	fromTokenID, err := middleware.GetUserID(ctx)
	if err != nil {
		e.NewErrorResponse(ctx, http.StatusInternalServerError, err)
		return
	}
	h.logger.Info(fromTokenID)

	fromUrlID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		h.logger.Info("error while getting id from request")
		e.NewErrorResponse(ctx, http.StatusInternalServerError, err)
		return
	}
	h.logger.Info(fromUrlID)

	if fromUrlID != fromTokenID {
		e.NewErrorResponse(ctx, http.StatusForbidden, err)
		return
	}

	a, err := h.service.GetOne(fromUrlID)
	if err != nil {
		e.NewErrorResponse(ctx, http.StatusInternalServerError, err)
		return
	}

	out := h.mapper.MapAccountDTO(a)
	ctx.JSON(http.StatusOK, out)
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
	//h.logger.Info("Got login request")
	var loginDto dto.LoginAccountDTO

	if err := ctx.BindJSON(&loginDto); err != nil {
		h.logger.Error(err)
		e.NewErrorResponse(ctx, http.StatusBadRequest, err)
		return
	}

	a := h.mapper.MapLogInAccountDTO(loginDto)

	token, err := h.service.GenerateJWT(&a)
	if err != nil {
		e.NewErrorResponse(ctx, http.StatusUnauthorized, err)
		return
	}
	loginWithTokenDto := h.mapper.MapAccountWithTokenDTO(token, a)

	ctx.JSON(http.StatusOK, loginWithTokenDto)
}
