package server_test

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"jwt-enabled-api/internal/app"
	"jwt-enabled-api/internal/auth"
	"jwt-enabled-api/internal/config"
	"jwt-enabled-api/internal/hasher"
	"jwt-enabled-api/internal/repo"
	"jwt-enabled-api/internal/server"

	"github.com/gavv/httpexpect/v2"
	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

func TestUser(t *testing.T) {
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("failed creating new zap logger: %s", err)
	}

	signingKey, err := auth.ECPrivateKeyFromBase64(
		[]byte(config.Config.Auth.ECDSASigningKeyBase64),
		logger,
	)
	if err != nil {
		log.Fatalf("failed creating signing key from base64: %s", err)
	}
	auth := auth.NewAuth(
		signingKey,
		config.Config.Auth.ExpireInSecs.Jwt,
		config.Config.Auth.ExpireInSecs.Refresh,
	)

	hasher := hasher.NewBcrypt(bcrypt.DefaultCost)

	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=disable",
		config.Config.Postgres.User,
		config.Config.Postgres.Password,
		config.Config.Postgres.Host,
		config.Config.Postgres.Port,
		config.Config.Postgres.DB,
	)
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("failed opening database connection %s: %s", dsn, err)
	}

	repo := repo.NewRepository(db)
	app := app.NewApplication(repo, auth, hasher)

	router := echo.New()

	server := server.NewServer(
		app,
		router,
		func(ctx echo.Context, s string) (any, error) { return auth.ParseJwtToken(s) },
		nil,
	)

	htServer := httptest.NewServer(server.GetHandler())
	e := httpexpect.Default(t, htServer.URL)

	// run tests

	e.Request(http.MethodGet, "/v1/nocontent").
		Expect().
		Status(http.StatusOK).
		NoContent()

	e.Request(http.MethodGet, "/v1/error").
		Expect().
		Status(http.StatusBadRequest).
		JSON().
		Object().
		Equal(echo.Map{
			"message": http.StatusText(http.StatusBadRequest),
		})

	e.Request(http.MethodGet, "/v1/error-string").
		Expect().
		Status(http.StatusBadRequest).
		JSON().
		Object().
		Equal(echo.Map{
			"message": "error-string",
		})

	e.Request(http.MethodGet, "/v1/error-error").
		Expect().
		Status(http.StatusBadRequest).
		JSON().
		Object().
		Equal(echo.Map{
			// "message": "error-error",
		})

	e.Request(http.MethodGet, "/v1/error-int").
		Expect().
		Status(http.StatusBadRequest).
		JSON().
		Number().
		Equal(22)

	// teardown
}
