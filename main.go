package main

import (
	"database/sql"
	"embed"
	"fmt"
	"log"

	"jwt-enabled-api/internal/app"
	"jwt-enabled-api/internal/auth"
	"jwt-enabled-api/internal/config"
	"jwt-enabled-api/internal/hasher"
	"jwt-enabled-api/internal/repo"
	"jwt-enabled-api/internal/server"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"

	_ "github.com/lib/pq"
)

//go:embed openapi
var openapiFS embed.FS

func main() {
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

	server := server.NewServer(
		app,
		echo.New(),
		func(ctx echo.Context, s string) (any, error) { return auth.ParseJwtToken(s) },
		echo.MustSubFS(openapiFS, "openapi"),
	)
	log.Fatal(server.Run(":8080"))
}
