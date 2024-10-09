package server

import (
	"context"
	"errors"
	"github.com/chloyka/ginannot"
	"github.com/gin-gonic/gin"
	"github.com/gopybara/utils"
	"github.com/gopybara/utils/server/validator"
	"github.com/gopybara/utils/supervisor"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"net/http"
	"os"
	"time"
)

type HttpServer struct {
	engine   *gin.Engine
	log      *zap.Logger
	handlers []ginannot.Handler
	srv      *http.Server
}

func NewHttpServer(handlers []ginannot.Handler, log *zap.Logger) *HttpServer {
	logsDisabled := os.Getenv("DISABLE_ACCESS_LOGS") == "true"
	r := gin.New()
	r.Use(traceparentMiddleware())
	r.Use(gin.Recovery())
	if !logsDisabled {
		r.Use(logMiddleware(log))
	}
	handlers = append(handlers, ginannot.Handler(&Middlewares{}))
	srv := &http.Server{
		Addr:    ":1489",
		Handler: r,
	}

	return &HttpServer{
		engine:   r,
		log:      log,
		handlers: handlers,
		srv:      srv,
	}
}

func (s *HttpServer) Engine() *gin.Engine {
	return s.engine
}

func RunHTTPServer(
	lifecycle fx.Lifecycle,
	s *HttpServer,
	log *zap.Logger,
	supervisor *supervisor.GoroutineSupervisor,
) {
	lifecycle.Append(validator.ValidatorLifecycleHook())
	lifecycle.Append(
		fx.Hook{
			OnStart: func(context.Context) error {
				annotator := ginannot.New(s.engine)
				annotator.Apply(s.handlers)

				hasHealth := false
				hasReadiness := false
				hasStartup := false

				processed := make(map[string]bool)

				for _, c := range s.engine.Routes() {
					if c.Path == "/health" || c.Path == "/healthz" {
						hasHealth = true
					}

					if c.Path == "/readiness" {
						hasReadiness = true
					}

					if c.Path == "/startup" {
						hasStartup = true
					}

					if processed[c.Path] {
						continue
					}

					processed[c.Path] = true
					s.engine.OPTIONS(c.Path, func(ctx *gin.Context) {
						ctx.Status(200)
					})
				}

				if !hasHealth {
					s.engine.GET("/health")
					s.engine.GET("/healthz")
				}
				if !hasReadiness {
					s.engine.GET("/readiness")
				}

				if !hasStartup {
					s.engine.GET("/startup")
				}
				go func() {
					if err := s.srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
						log.Fatal("failed to listen and serve http server", zap.Error(err))
					}
				}()

				return nil
			},
			OnStop: func(ctx context.Context) error {
				log.Info("gracefully shutting down http server")
				supervisor.WaitForComplete()
				log.Info("all necessary goroutines are stopped, shutting down http server")

				return s.srv.Shutdown(ctx)
			},
		},
	)
}

func logMiddleware(log *zap.Logger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		timeSpent := time.Now().UnixMicro()
		ctx.Next()
		if ctx.Request.URL.Path == "/health" || ctx.Request.URL.Path == "/healthz" || ctx.Request.URL.Path == "/readiness" {
			return
		}
		end := float64(time.Now().UnixMicro()-timeSpent) / 1000
		cloned := ctx.Copy()

		log.Info("incoming request",
			zap.String("method", cloned.Request.Method),
			zap.String("path", cloned.Request.URL.Path),
			zap.Any("query", cloned.Request.URL.Query()),
			zap.Int("responseStatus", cloned.Writer.Status()),
			zap.String("traceparent", cloned.Request.Header.Get("Traceparent")),
			zap.Float64("timeMs", end))
	}
}

func ProvideHttpServer() fx.Option {
	return fx.Provide(
		supervisor.NewGoroutineSupervisor,
		fx.Annotate(
			NewHttpServer,
			fx.ParamTags(`group:"handlers"`),
		),
	)
}

func InvokeHTTPServer() fx.Option {
	return fx.Invoke(RunHTTPServer)
}

type MiddlewaresInterface struct {
	CORSMiddleware ginannot.Middleware `middleware:"name=cors"`
}

type Middlewares struct {
	MiddlewaresInterface
}

func (m *Middlewares) CORSMiddleware(ctx *gin.Context) {
	ctx.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	ctx.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
	ctx.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, ResponseType, accept, origin, Cache-Control, X-Requested-With")
	ctx.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, PATCH, GET, PUT, DELETE")

	ctx.Next()
}

func traceparentMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if ctx.Request.Header.Get("Traceparent") == "" || ctx.Request.Header.Get("Traceparent") == "0" {
			ctx.Request.Header.Set("Traceparent", utils.GenerateTraceparent())
		}

		ctx.Next()
	}
}
