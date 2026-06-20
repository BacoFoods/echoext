package echoext

import (
	stdcontext "context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/labstack/gommon/color"
	echoSwagger "github.com/swaggo/echo-swagger"

	validator "github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	emiddleware "github.com/labstack/echo/v4/middleware"
)

// shutdownTimeout bounds how long graceful shutdown waits for in-flight
// requests to drain before forcing termination.
const shutdownTimeout = 10 * time.Second

type EchoMode string

const (
	TestMode     EchoMode = "test"
	StandardMode EchoMode = "standard"
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
	mode    string
}

func New(cl ...ServerConfig) Server {
	c := ServerConfig{}
	if len(cl) > 0 {
		c = cl[0]
	}

	s := echo.New()
	s.HideBanner = true

	s.Validator = &Validator{
		v: validator.New(
			validator.WithRequiredStructEnabled(),
		),
	}

	c.HealthcheckPath = c.escapeHealthcheckSuffix()
	c.PathPrefix = c.escapePrefix()

	s.Use(CustomLogger(c))
	s.Use(CustomRecovery)
	s.Use(CustomCORS(c))

	if !c.MetricsConfig.Disabled {
		s.Use(metricsMiddleware)
	}

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

	if !c.MetricsConfig.Disabled {
		colorer.Printf("[%s] metrics: %s\n", colorer.Green("echoext"), colorer.Blue(fmt.Sprintf("http://%s:%d%s", "0.0.0.0", c.MetricsConfig.escapePort(), c.MetricsConfig.escapePath())))
	}

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

func escapePath(path string) string {
	p := strings.TrimSpace(path)
	// empty case
	if p == "" {
		p = "/"
	} else if p[0] != '/' { // ensure leading /
		p = "/" + p
	}

	// remove trailing /
	if len(p) > 0 && p[len(p)-1] == '/' {
		p = p[:len(p)-1]
	}

	return strings.ToLower(p)
}

func (s extServer) Group(prefix string, mount setupfn, middlewares ...MiddlewareFunc) *Group {
	p := escapePath(prefix)

	s.colorer.Printf("[%s] group prefix: %s\n", s.colorer.Green("echoext"), s.colorer.Blue(s.config.PathPrefix+p))

	g := s.root.NewGroup(p, middlewares...)

	mount(g)

	return g
}

// Start boots the main HTTP server and, unless disabled, the dedicated
// Prometheus metrics server. It blocks until either server fails or an
// interrupt/terminate signal is received, at which point both servers are
// gracefully shut down within shutdownTimeout.
func (s extServer) Start() error {
	// Buffered for both servers so a failing goroutine never blocks on send.
	errCh := make(chan error, 2)

	var metricsSrv *http.Server
	if !s.config.MetricsConfig.Disabled {
		metricsSrv = newMetricsServer(s.config)
		go func() {
			if err := metricsSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				errCh <- err
			}
		}()
	}

	go func() {
		if err := s.Echo.Start(s.config.escapeHost()); err != nil && err != http.ErrServerClosed {
			errCh <- err
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-errCh:
		// One of the servers failed to start; tear down the other.
		s.shutdown(metricsSrv)
		return err
	case <-quit:
		return s.shutdown(metricsSrv)
	}
}

// shutdown gracefully stops the main server and, when present, the metrics
// server, sharing a single timeout-bounded context.
func (s extServer) shutdown(metricsSrv *http.Server) error {
	ctx, cancel := stdcontext.WithTimeout(stdcontext.Background(), shutdownTimeout)
	defer cancel()

	err := s.Echo.Shutdown(ctx)

	if metricsSrv != nil {
		if mErr := metricsSrv.Shutdown(ctx); mErr != nil && err == nil {
			err = mErr
		}
	}

	return err
}

func (s extServer) Engine() *echo.Echo {
	return s.Echo
}
