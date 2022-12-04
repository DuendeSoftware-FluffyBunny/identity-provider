package server

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	ory "github.com/ory/client-go"
	"github.com/voice0726/identity-provider/src/config"
	"go.uber.org/zap"
	"golang.org/x/oauth2"
)

type Server struct {
	kratosPublicAPIClient *ory.APIClient
	kratosAdminAPIClient  *ory.APIClient
	kratosPublicEndpoint  string
	hydraPublicAPIClient  *ory.APIClient
	hydraAdminAPIClient   *ory.APIClient
	oAuth2Config          *oauth2.Config
	hydraPublicEndpoint   string
	port                  string
	echo                  *echo.Echo
	logger                *zap.Logger
}

func NewServer(config *config.Config, logger *zap.Logger) *Server {
	kratosPublicAPIConf := ory.Configuration{Servers: ory.ServerConfigurations{{URL: config.KratosPublicEndpoint}}}
	kratosAdminAPIConf := ory.Configuration{Servers: ory.ServerConfigurations{{URL: config.KratosAdminEndpoint}}}
	hydraPublicAPIConf := ory.Configuration{Servers: ory.ServerConfigurations{{URL: config.HydraPublicEndpoint}}}
	hydraAdminAPIConf := ory.Configuration{Servers: ory.ServerConfigurations{{URL: config.HydraAdminEndpoint}}}

	oauth2Conf := oauth2.Config{
		ClientID:     config.OAuth2ClientID,
		ClientSecret: config.OAuth2ClientSecret,
		Endpoint: oauth2.Endpoint{
			AuthURL:  fmt.Sprintf("%s/oauth2/auth", config.HydraPublicEndpoint),
			TokenURL: fmt.Sprintf("%s/oauth2/token", config.HydraPublicEndpoint),
		},
		RedirectURL: "http://localhost:4455/callback",
		Scopes:      []string{"openid", "offline_access"},
	}

	return &Server{
		kratosPublicAPIClient: ory.NewAPIClient(&kratosPublicAPIConf),
		kratosAdminAPIClient:  ory.NewAPIClient(&kratosAdminAPIConf),
		kratosPublicEndpoint:  config.KratosPublicEndpoint,
		hydraPublicAPIClient:  ory.NewAPIClient(&hydraPublicAPIConf),
		hydraAdminAPIClient:   ory.NewAPIClient(&hydraAdminAPIConf),
		hydraPublicEndpoint:   config.HydraPublicEndpoint,
		oAuth2Config:          &oauth2Conf,
		echo:                  echo.New(),
		port:                  config.AppPort,
		logger:                logger,
	}
}

func (s *Server) Register() {
	s.echo.Renderer = NewTemplate()
	s.echo.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI:       true,
		LogStatus:    true,
		LogMethod:    true,
		LogRemoteIP:  true,
		LogUserAgent: true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			s.logger.Info(
				"request",
				zap.String("ip", v.RemoteIP),
				zap.String("method", v.Method),
				zap.String("URI", v.URI),
				zap.Int("status", v.Status),
				zap.String("user_agent", v.UserAgent),
			)
			return nil
		},
	}))
	s.echo.GET("/login", s.handleLogin)
	s.echo.GET("/callback", s.handleCallback)
	s.echo.GET("/auth/consent", s.handleConsent)
	s.echo.GET("/dashboard", s.handleDashboard)
	s.echo.GET("/logout", s.handleLogout)
	s.echo.GET("/error", s.handleError)
}

func (s *Server) Start() error {
	return s.echo.Start(":" + s.port)
}

func convertCookies(cookies []*http.Cookie) string {
	var cookieStrArr []string
	for _, cookie := range cookies {
		cookieStrArr = append(cookieStrArr, cookie.Raw)
	}
	return strings.Join(cookieStrArr, " ")
}

func getUiInputNodeByAttributeName(name string, nodes []ory.UiNode) *ory.UiNode {
	for _, node := range nodes {
		if node.Type != "input" {
			continue
		}
		attr := node.Attributes.GetActualInstance().(*ory.UiNodeInputAttributes)
		if attr.Name == name {
			return &node
		}
	}
	return nil
}

func parseSetCookies(cookie string) []*http.Cookie {
	h := http.Header{}
	h.Add("Set-Cookie", cookie)
	res := http.Response{Header: h}
	return res.Cookies()
}
