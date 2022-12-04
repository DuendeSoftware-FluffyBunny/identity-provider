package server

import (
	"encoding/json"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (s *Server) handleDashboard(c echo.Context) error {
	// get cookie from headers
	cookie := c.Request().Header.Get("cookie")
	// get session details
	session, _, err := s.kratosPublicAPIClient.FrontendApi.ToSession(c.Request().Context()).Cookie(cookie).Execute()
	if err != nil {
		return c.Redirect(http.StatusSeeOther, "/login")
	}

	// marshal session to json
	sessionJSON, err := json.MarshalIndent(session, "", "  ")
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.Render(http.StatusOK, "index.html", map[string]interface{}{
		"Title":   "Session details",
		"Details": string(sessionJSON),
	})
}
