package app

import (
	"jwt-enabled-api/internal/auth"
	"jwt-enabled-api/internal/hasher"
	"jwt-enabled-api/internal/repo"
)

type Application struct {
	repo   repo.ServiceTx
	auth   auth.Interface
	hasher hasher.Interface
}

func NewApplication(
	repo repo.ServiceTx,
	auth auth.Interface,
	hasher hasher.Interface,
) *Application {
	return &Application{repo: repo, auth: auth, hasher: hasher}
}
