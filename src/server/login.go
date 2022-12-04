package server

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"go.uber.org/zap"
)

func (s *Server) handleLogin(c echo.Context) error {
	loginChallenge := c.QueryParam("login_challenge")
	flowID := c.QueryParam("flow")
	cookie := c.Request().Header.Get("cookie")

	if loginChallenge == "" && flowID == "" {
		// create oauth2 state and store in session
		b := make([]byte, 32)
		_, err := rand.Read(b)
		if err != nil {
			log.Error("generate state failed: %v", err)
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}
		state := base64.StdEncoding.EncodeToString(b)
		sc := http.Cookie{
			Name:     "oauth2_state",
			Value:    state,
			Path:     "/",
			Secure:   true,
			HttpOnly: true,
		}
		c.SetCookie(&sc)
		url := s.oAuth2Config.AuthCodeURL(state)

		return c.Redirect(http.StatusSeeOther, url)
	}

	if flowID == "" {
		param := url.Values{
			"login_challenge": []string{loginChallenge},
		}
		return c.Redirect(http.StatusSeeOther, fmt.Sprintf("%s/self-service/login/browser?%s", s.kratosPublicEndpoint, param.Encode()))
	}

	loginFlow, res, err := s.kratosPublicAPIClient.FrontendApi.GetLoginFlow(c.Request().Context()).Id(flowID).Cookie(cookie).Execute()
	if err != nil {
		b, _ := io.ReadAll(res.Body)
		s.logger.Debug("failed to create browser login flow", zap.ByteString("response", b), zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	ui := loginFlow.Ui
	return c.Render(http.StatusOK, "index.html", map[string]interface{}{
		"Title": "Login",
		"UI":    ui,
	})
}
