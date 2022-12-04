package server

import (
	"fmt"
	"io"
	"net/http"

	"go.uber.org/zap"

	"github.com/AlekSi/pointer"
	"github.com/labstack/echo/v4"
	ory "github.com/ory/client-go"
)

func (s *Server) handleConsent(c echo.Context) error {
	challenge := c.QueryParam("consent_challenge")

	if challenge == "" {
		return fmt.Errorf("no consent challenge found")
	}

	consentReq, res, err := s.hydraAdminAPIClient.OAuth2Api.GetOAuth2ConsentRequest(c.Request().Context()).ConsentChallenge(challenge).Execute()
	if err != nil {
		b, _ := io.ReadAll(res.Body)
		s.logger.Debug("failed to get consent request", zap.ByteString("response", b), zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	cookie := c.Request().Header.Get("cookie")
	session, res, err := s.kratosPublicAPIClient.FrontendApi.ToSession(c.Request().Context()).Cookie(cookie).Execute()
	if err != nil {
		b, _ := io.ReadAll(res.Body)
		s.logger.Debug("failed to get session", zap.ByteString("response", b), zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	redirectTo, res, err := s.hydraAdminAPIClient.OAuth2Api.AcceptOAuth2ConsentRequest(c.Request().Context()).ConsentChallenge(challenge).AcceptOAuth2ConsentRequest(
		ory.AcceptOAuth2ConsentRequest{
			GrantScope:  consentReq.RequestedScope,
			Remember:    pointer.ToBool(true),
			RememberFor: pointer.ToInt64(3600),
			Session: &ory.AcceptOAuth2ConsentRequestSession{
				IdToken: session,
			},
		}).Execute()
	if err != nil {
		s.logger.Debug("accept consent request rejected", zap.Any("response", res), zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.Redirect(http.StatusSeeOther, redirectTo.RedirectTo)
}
