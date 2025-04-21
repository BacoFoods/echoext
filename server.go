package echoext

import (
	"net/http"
	"os"
	"strings"

	"github.com/labstack/gommon/color"
	echoSwagger "github.com/swaggo/echo-swagger"

	"github.com/labstack/echo/v4"
	emiddleware "github.com/labstack/echo/v4/middleware"
)

type Server interface {
	Start() error
	Group(string, setupfn, ...MiddlewareFunc) *Group
	Engine() *echo.Echo
}

type Mountable interface {
	Mount(*echo.Group)
}

type extServer struct {
	*echo.Echo
	config  ServerConfig
	colorer *color.Color
	appEnv  string
	root    *Group
}

func New(cl ...ServerConfig) Server {
	c := ServerConfig{}
	if len(cl) > 0 {
		c = cl[0]
	}

	s := echo.New()
	s.HideBanner = true

	c.HealthcheckPath = c.escapeHealthcheckSuffix()
	c.PathPrefix = c.escapePrefix()

	s.Use(CustomLogger(c))
	s.Use(CustomRecovery)
	s.Use(CustomCORS(c))

	root := s.Group(c.PathPrefix)
	root.GET(c.escapeHealthcheckSuffix(), func(ctx echo.Context) error {
		return ctx.JSON(http.StatusOK, echo.Map{"status": "ok"})
	})

	colorer := color.New()
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "local"
	}

	colorer.Printf("[%s] app enviroment: %s\n", colorer.Green("echoext"), colorer.Blue(env))

	colorer.Printf("[%s] server prefix: %s\n", colorer.Green("echoext"), colorer.Blue(c.PathPrefix))
	colorer.Printf("[%s] healthcheck path: %s\n", colorer.Green("echoext"), colorer.Blue(c.healthcheckFullPath()))

	if env != "production" {
		sp := c.swaggerPath()
		colorer.Printf("[%s] swagger docs: %s\n", colorer.Green("echoext"), colorer.Blue("http://"+c.escapeHost()+sp+"/index.html"))
		swaggerCreds := os.Getenv("SWAGGER_CREDENTIALS")

		swaggerUser := ""
		swaggerPass := ""

		parts := strings.Split(swaggerCreds, ":")
		if len(parts) == 2 {
			swaggerUser = parts[0]
			swaggerPass = parts[1]
		}

		swaggerAuth := emiddleware.BasicAuthWithConfig(emiddleware.BasicAuthConfig{
			Skipper: nil,
			Validator: func(u string, p string, ctx echo.Context) (bool, error) {
				return u == swaggerUser && p == swaggerPass, nil
			},
			Realm: "",
		})

		s.GET(sp+"/*", echoSwagger.EchoWrapHandler(), swaggerAuth)
	}

	colorer.Println()

	return extServer{
		Echo:    s,
		config:  c,
		colorer: colorer,
		appEnv:  env,
		root:    &Group{root},
	}
}

type setupfn func(*Group)

func (s extServer) Group(prefix string, mount setupfn, middlewares ...MiddlewareFunc) *Group {
	p := prefix
	// ensure leading /
	if p[0] != '/' {
		p = "/" + p
	}

	// remove trailing /
	if p[len(p)-1] == '/' {
		p = p[:len(p)-1]
	}

	s.colorer.Printf("[%s] group prefix: %s\n", s.colorer.Green("echoext"), s.colorer.Blue(s.config.PathPrefix+p))

	g := s.root.NewGroup(p, middlewares...)

	mount(g)

	return g
}

func (s extServer) Start() error {
	return s.Echo.Start(s.config.escapeHost())
}

func (s extServer) Engine() *echo.Echo {
	return s.Echo
}
