package server

import (
	"errors"
	"fmt"
	"io/fs"
	"net/http"

	"jwt-enabled-api/internal/app"
	"jwt-enabled-api/internal/auth"

	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

type Server struct {
	app    *app.Application
	router *echo.Echo
	// logger *zap.Logger
	parseTokenFunc func(echo.Context, string) (any, error)
	openapiSubFS   fs.FS
}

func NewServer(
	app *app.Application,
	router *echo.Echo,
	parseTokenFunc func(echo.Context, string) (any, error),
	openapiSubFS fs.FS,
) *Server {
	server := &Server{
		app:            app,
		router:         router,
		parseTokenFunc: parseTokenFunc,
		openapiSubFS:   openapiSubFS,
	}
	server.setHandlers()
	return server
}

func (s *Server) GetHandler() http.Handler {
	return s.router.Server.Handler
}

func (s *Server) setHandlers() {
	v1 := s.router.Group("/v1")

	////// test routes /////////////////////////////////////////////////////////
	v1.GET("/nocontent", func(c echo.Context) error {
		return c.NoContent(http.StatusOK)
	})
	v1.GET("/error", func(c echo.Context) error {
		httpError := echo.NewHTTPError(http.StatusBadRequest)
		fmt.Println(httpError.Error())
		return httpError
	})
	v1.GET("/error-string", func(c echo.Context) error {
		httpError := echo.NewHTTPError(http.StatusBadRequest, "error-string")
		fmt.Println(httpError.Error())
		return httpError
	})
	v1.GET("/error-empty-string", func(c echo.Context) error {
		httpError := echo.NewHTTPError(http.StatusBadRequest, "")
		fmt.Println(httpError.Error())
		return httpError
	})
	v1.GET("/error-error", func(c echo.Context) error {
		httpError := echo.NewHTTPError(
			http.StatusBadRequest,
			errors.New("error-error"),
		)
		fmt.Println(httpError.Error())
		return httpError
	})
	v1.GET("/error-error-internal", func(c echo.Context) error {
		err := errors.New("error-error")
		httpError := echo.NewHTTPError(
			http.StatusBadRequest,
			err,
		).SetInternal(err)
		fmt.Println(httpError.Error())
		return httpError
	})
	v1.GET("/error-int", func(c echo.Context) error {
		httpError := echo.NewHTTPError(
			http.StatusBadRequest,
			22,
		)
		fmt.Println(httpError.Error())
		return httpError
	})
	////////////////////////////////////////////////////////////////////////////

	// serve openapi spec documentation
	v1.StaticFS("/openapi", s.openapiSubFS)

	public := v1.Group("/public")
	publicUser := public.Group("/user")
	publicUser.POST("/", s.HandleUserCreate)
	publicUser.POST("/login", s.HandleUserLogin)
	publicUserID := publicUser.Group("/:id")
	publicUserID.POST("/logout", s.HandleUserLogout)
	publicUserID.POST("/refresh", s.HandleUserRefreshToken)

	authorized := v1.Group("/authorized", echojwt.WithConfig(echojwt.Config{
		ContextKey:     contextKey,
		ParseTokenFunc: s.parseTokenFunc,
	}))

	authorized.GET("/check", func(c echo.Context) error {
		payload, httpError := s.getUserPayload(c)
		if httpError != nil {
			return httpError
		}
		return c.JSON(http.StatusOK, echo.Map{
			"id": payload.UserID,
		})
	})

	authorizedUser := authorized.Group("/user")
	authorizedUser.GET("/:id", s.HandleUserGet)
}

func (s *Server) Run(addr string) error {
	return s.router.Start(addr)
}

const contextKey string = "token_payload"

func (s *Server) getUserPayload(
	c echo.Context,
) (*auth.Payload, *echo.HTTPError) {
	payload, ok := c.Get(contextKey).(*auth.Payload)
	if !ok {
		return nil, echo.NewHTTPError(
			http.StatusInternalServerError,
			"context key '"+contextKey+"' not set",
		)
	}
	return payload, nil
}
