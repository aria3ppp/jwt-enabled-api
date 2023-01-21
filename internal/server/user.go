package server

import (
	"net/http"

	"jwt-enabled-api/internal/app"
	"jwt-enabled-api/internal/dto"

	"github.com/labstack/echo/v4"
)

func (s *Server) HandleUserGet(c echo.Context) error {
	// bind and validate param
	var param dto.IDParam
	if httpError := (&echo.DefaultBinder{}).BindPathParams(c, &param); httpError != nil {
		return httpError
	}
	if err := param.Validate(); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// get
	user, err := s.app.UserGet(c.Request().Context(), param.ID)
	if err != nil {
		if err == app.ErrNotFound {
			return echo.NewHTTPError(http.StatusNotFound)
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, user)
}

func (s *Server) HandleUserCreate(c echo.Context) error {
	// bind and validate request
	var req dto.UserCreateRequest
	if httpError := (&echo.DefaultBinder{}).BindBody(c, &req); httpError != nil {
		return httpError
	}
	if err := req.Validate(); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// create
	userID, err := s.app.UserCreate(c.Request().Context(), &req)
	if err != nil {
		if err == app.ErrUsedEmail {
			return echo.NewHTTPError(http.StatusBadRequest, "email used")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, echo.Map{
		"id": userID,
	})
}

func (s *Server) HandleUserLogin(c echo.Context) error {
	// bind and validate request
	var req dto.UserLoginRequest
	if httpError := (&echo.DefaultBinder{}).BindBody(c, &req); httpError != nil {
		return httpError
	}
	if err := req.Validate(); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// login
	resp, err := s.app.UserLogin(c.Request().Context(), &req)
	if err != nil {
		if err == app.ErrNotFound {
			return echo.NewHTTPError(http.StatusNotFound, "email not found")
		}
		if err == app.ErrIncorrectPassword {
			return echo.NewHTTPError(
				http.StatusUnauthorized,
				"incorrect password",
			)
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, resp)
}

func (s *Server) HandleUserLogout(c echo.Context) error {
	// bind and validate param
	var param dto.IDParam
	if httpError := (&echo.DefaultBinder{}).BindPathParams(c, &param); httpError != nil {
		return httpError
	}
	if err := param.Validate(); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// bind and validate request
	var req dto.UserLogoutRequest
	if httpError := (&echo.DefaultBinder{}).BindBody(c, &req); httpError != nil {
		return httpError
	}
	if err := req.Validate(); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// logout
	err := s.app.UserLogout(c.Request().Context(), param.ID, req.Token)
	if err != nil {
		if err == app.ErrNotFound {
			return echo.NewHTTPError(http.StatusNotFound, "token not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.NoContent(http.StatusNoContent)
}

func (s *Server) HandleUserRefreshToken(c echo.Context) error {
	// bind and validate param
	var param dto.IDParam
	if httpError := (&echo.DefaultBinder{}).BindPathParams(c, &param); httpError != nil {
		return httpError
	}
	if err := param.Validate(); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// bind and validate request
	var req dto.UserRefreshTokenRequest
	if httpError := (&echo.DefaultBinder{}).BindBody(c, &req); httpError != nil {
		return httpError
	}
	if err := req.Validate(); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// refresh token
	resp, err := s.app.UserRefreshToken(
		c.Request().Context(),
		param.ID,
		req.Token,
	)
	if err != nil {
		if err == app.ErrNotFound {
			return echo.NewHTTPError(http.StatusNotFound, "token not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, resp)
}
