package account

import (
	"errors"
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
	// TODO: logout endpoint
	groupName := fmt.Sprintf("%v/v%v%v", apiURLGroup, apiVersion, accountsURLGroup)

	h.logger.Tracef("Register route: %v", groupName)

	auth := router.Group(groupName)
	{
		auth.POST(registerURL, h.RegisterAccount)
		auth.POST(loginURL, h.Login)
	}

	accounts := router.Group(groupName, middleware.Authenticate)
	{
		accounts.GET("/:id", h.GetAccount)
	}
}

// RegisterAccount creates account
// @Summary RegisterAccount
// @Tags account
// @Description create account
// @ID create-account
// @Accept  json
// @Produce  json
// @Param dto body dto.RegisterAccountDTO true "account info"
// @Success 201 {string} string 1
// @Failure 500 {object} e.ErrorResponse
// @Failure 409 {object} e.ErrorResponse
// @Failure default {object} e.ErrorResponse
// @Router /api/v1/accounts/register [post]
func (h *Handler) RegisterAccount(ctx *gin.Context) {
	var (
		err error
		in  dto.RegisterAccountDTO
		a   model.Account
	)
	h.logger.Info("Got registration request")

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

	h.logger.Infof("CreateAccountDTO mapped: %v", a)

	err = h.service.CreateAccount(&a)
	if err != nil {
		if errors.Is(err, e.ClientAccountError) {
			e.NewErrorResponse(ctx, http.StatusConflict, err)
		} else {
			e.NewErrorResponse(ctx, http.StatusInternalServerError, err)
		}
		return
	}

	h.logger.Infof("Inserted into database successfully: account id is %v", a.ID)

	ctx.JSON(http.StatusCreated, fmt.Sprintf(
		"%v/v%v%v/%v",
		apiURLGroup,
		apiVersion,
		accountsURLGroup,
		a.ID,
	))
}

// GetAccount
// @Summary GetAccount
// @Security ApiKeyAuth
// @Tags account
// @Description get account
// @ID get-account
// @Accept  json
// @Produce  json
// @Param id   path  string  true  "id"
// @Success 200 {object} dto.GetAccountDTO
// @Failure 500 {object} e.ErrorResponse
// @Failure 403 {object} e.ErrorResponse
// @Failure default {object} e.ErrorResponse
// @Router /api/v1/accounts/{id} [get]
func (h *Handler) GetAccount(ctx *gin.Context) {
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
		e.NewErrorResponse(ctx, http.StatusForbidden, errors.New("permission denied"))
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

// Login
// @Summary Login
// @Tags account
// @Description login
// @ID login
// @Accept  json
// @Produce  json
// @Param dto body dto.LoginAccountDTO true "credentials"
// @Success 200 {object} dto.WithTokenDTO
// @Failure 500 {object} e.ErrorResponse
// @Failure default {object} e.ErrorResponse
// @Router /api/v1/accounts/login [post]
func (h *Handler) Login(ctx *gin.Context) {
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
	ctx.SetCookie("token", token, 36000, "/", "localhost", false, true)

	loginWithTokenDto := h.mapper.MapAccountWithTokenDTO(token, a)

	ctx.JSON(http.StatusOK, loginWithTokenDto)
}
