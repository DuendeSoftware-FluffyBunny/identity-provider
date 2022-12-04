package server

import (
	"fmt"
	"io"
	"net/http"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

func (s *Server) handleCallback(c echo.Context) error {
	cookie := c.Request().Header.Get("cookie")

	session, res, err := s.kratosPublicAPIClient.FrontendApi.ToSession(c.Request().Context()).Cookie(cookie).Execute()
	if err != nil {
		b, _ := io.ReadAll(res.Body)
		s.logger.Debug("failed to create browser login flow", zap.ByteString("response", b), zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	s.logger.Debug("session", zap.Any("session", session))

	state, err := c.Cookie("oauth2_state")
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("invalid oauth2 state"))
	}

	if c.QueryParam("state") != state.Value {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("invalid oauth2 state"))
	}

	code := c.QueryParam("code")
	token, err := s.oAuth2Config.Exchange(c.Request().Context(), code)
	if err != nil {
		s.logger.Debug("failed to exchange token", zap.Error(err))
		return echo.NewHTTPError(http.StatusUnauthorized, fmt.Errorf("invalid oauth2 state"))
	}

	s.logger.Info("callback completed", zap.String("id_token", token.AccessToken))

	return c.Redirect(http.StatusSeeOther, "http://localhost:4455/dashboard")
}
