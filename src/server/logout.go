package server

import (
	"fmt"
	"io"
	"net/http"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

func (s *Server) handleLogout(c echo.Context) error {
	cookie := c.Request().Header.Get("cookie")
	s.logger.Debug("cookie", zap.String("cookie", cookie))
	logoutFlow, res, err := s.kratosPublicAPIClient.FrontendApi.CreateBrowserLogoutFlow(c.Request().Context()).Cookie(cookie).Execute()
	if err != nil {
		b, _ := io.ReadAll(res.Body)
		s.logger.Debug("failed to create browser logout flow", zap.ByteString("response", b), zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	s.logger.Debug(fmt.Sprintf("logout url: %s", logoutFlow.LogoutUrl))

	return c.Redirect(http.StatusSeeOther, logoutFlow.LogoutUrl)
}
