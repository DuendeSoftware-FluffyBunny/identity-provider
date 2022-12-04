package server

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

func (s *Server) handleError(c echo.Context) error {
	// get url query parameters
	errorID := c.QueryParam("id")
	// get error details
	errorDetails, res, err := s.kratosPublicAPIClient.FrontendApi.GetFlowError(c.Request().Context()).Id(errorID).Execute()
	if err != nil {
		b, _ := io.ReadAll(res.Body)
		s.logger.Debug("failed to get flow error", zap.ByteString("response", b), zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	// marshal errorDetails to json
	errorDetailsJSON, err := json.MarshalIndent(errorDetails, "", "  ")
	if err != nil {
		s.logger.Debug("failed to marshal flow error", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.Render(http.StatusOK, "index.html", map[string]interface{}{
		"Title":   "Error",
		"Details": string(errorDetailsJSON),
	})
}
