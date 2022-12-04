package server

import (
	"crypto/rand"
	"encoding/base64"
	"io"
	"net/http"
	"strings"

	"go.uber.org/zap"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
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
		loginFlow, res, err := s.kratosPublicAPIClient.FrontendApi.CreateBrowserLoginFlow(c.Request().Context()).LoginChallenge(loginChallenge).Execute()
		if err != nil {
			b, _ := io.ReadAll(res.Body)
			s.logger.Debug("failed to create browser login flow", zap.ByteString("response", b), zap.Error(err))
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}
		setCookie := res.Header.Get("set-cookie")
		s.logger.Debug("cookie string", zap.String("raw_cookie", setCookie))
		cc := parseSetCookies(setCookie)

		ui := loginFlow.Ui

		for _, cookie := range cc {
			s.logger.Debug("cookie", zap.String(cookie.Name, cookie.Value))
			if strings.HasPrefix(cookie.Name, "csrf_token") {
				c.SetCookie(cookie)
			}
		}

		return c.Render(http.StatusOK, "index.html", map[string]interface{}{
			"Title": "Login",
			"UI":    ui,
		})
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
