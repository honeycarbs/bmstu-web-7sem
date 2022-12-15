//go:build e2e
// +build e2e

package e2e_test

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
	"neatly/internal/handlers/account"
	"neatly/internal/handlers/middleware"
	"neatly/internal/mapper"
	"neatly/internal/model"
	"neatly/internal/repository"
	"neatly/internal/service"
	"neatly/pkg/logging"
	"neatly/pkg/testdocker"
	"net/http/httptest"
	"testing"
)

var (
	registerAccountRequestBody = `{"name": "TestUser", "username": "testuser", "email": "test@user.com", "password": "testusertestuser"}`
	testAccount                = model.Account{Name: "TestUser", Username: "testuser", Email: "test@user.com", Password: "testusertestuser"}
)

func TestE2E_AccountCreateAndAuthorize(t *testing.T) {
	// === ARRANGE ===
	gin.SetMode(gin.TestMode)
	router := gin.New()

	middleware.CorsMiddleware(router)

	logging.Init()
	logger := logging.GetLogger()

	client, err := testdocker.GetTestResource("../etc/migrations")
	if err != nil {
		t.Fatal(err)
	}

	repo := repository.NewAccountRepositoryImpl(client, logger)
	serv := service.NewAccountServiceImpl(repo, logger)
	mppr := mapper.NewAccountMapper(logger)

	handler := account.NewHandler(logger, *serv, *mppr)
	handler.Register(router)

	expectedUserID := 1

	// === ACT ===
	t.Run("E2E: register + auth suite", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest(
			"POST", "/api/v1/accounts/register",
			bytes.NewBufferString(registerAccountRequestBody),
		)

		router.ServeHTTP(w, r)

		// === ASSERT ===
		assert.Equal(t, w.Code, 201)
		assert.Equal(t, w.Body.String(), fmt.Sprintf(`"/api/v1/accounts/%v"`, expectedUserID))

		a, err := serv.GetOne(expectedUserID)
		if err != nil {
			t.Fatalf("[service] Can't get test account: %s", err)
		}

		canAuth := repo.AuthorizeAccount(&testAccount)
		if canAuth != nil {
			t.Fatalf("[repository] Can't authorize test account: %s", err)
		}
		err = testAccount.CheckPassword(testAccount.Password)
		assert.Equal(t, err, nil)

		assert.Equal(t, a.ID, expectedUserID)
	})
}
