package main

import (
	"admin-dashboard/cmd/app/handlers"
	"admin-dashboard/cmd/app/middlewares"
	"admin-dashboard/cmd/app/repositories"
	webRoute "admin-dashboard/cmd/app/routes/webRoutes"
	"admin-dashboard/cmd/app/services"
	"admin-dashboard/cmd/common/mailService"
	"admin-dashboard/cmd/common/renderer"
	"admin-dashboard/cmd/internal/applicationData"
	"admin-dashboard/cmd/internal/cacheSystem"
	"admin-dashboard/cmd/internal/config"
	"admin-dashboard/cmd/internal/databases"
	"admin-dashboard/cmd/internal/mailer"
	appSession "admin-dashboard/cmd/internal/session"
	"admin-dashboard/cmd/internal/storage"
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {

	e := echo.New()

	cfg, err := config.Init()

	if err != nil {
		e.Logger.Fatal(err)
	}

	db, err := databases.ConnectDatabase(cfg)

	if err != nil {
		e.Logger.Fatal(err)
	}

	cache, err := cacheSystem.NewManager(cfg)

	if err != nil {
		e.Logger.Fatal(err)
	}

	mail, err := mailer.ConnectMailer(cfg, e.Logger, db)

	if err != nil {
		e.Logger.Fatal(err)
	}

	repository := repositories.Repository{Db: db, Logger: e.Logger, CacheSystem: cache}

	service := services.Service{Db: db, Logger: e.Logger, CacheSystem: cache, Storage: storage.New(cfg), Config: *cfg}

	workerCtx, workerCancel := context.WithCancel(context.Background())
	defer workerCancel()
	go service.StartWorker(workerCtx, time.Duration(cfg.QueueWorkerInterval)*time.Second)

	sendMail := mailService.SendMailService{Db: db, Logger: e.Logger, Mailer: mail}

	middlewareData := middlewares.Middleware{Db: db, Logger: e.Logger, Config: *cfg, Repository: repository, Service: service, CacheSystem: cache}

	handler := handlers.Handler{Logger: e.Logger, Db: db, Repository: repository, Service: service, SendMail: sendMail, Config: *cfg, Storage: storage.New(cfg), CacheSystem: cache}

	data := &applicationData.Application{Logger: e.Logger, Server: e, Handler: handler, Middleware: middlewareData, Repository: repository, Service: service, Mailer: sendMail, Storage: storage.New(cfg), CacheSystem: cache}

	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI:     true,
		LogStatus:  true,
		LogMethod:  true,
		LogLatency: true,
		LogError:   true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			fmt.Printf(
				"METHOD=%s URI=%s STATUS=%d LATENCY=%s ERROR=%v\n",
				v.Method,
				v.URI,
				v.Status,
				v.Latency,
				v.Error,
			)
			return nil
		},
	}))

	e.Use(middleware.Recover())

	sessionStore, err := appSession.NewSessionStore(cfg, db.GetDb())
	if err != nil {
		e.Logger.Fatal(err)
	}

	e.Use(session.Middleware(sessionStore))

	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("config", cfg)
			return next(c)
		}
	})

	web := e.Group("/")
	webRoute.MainRoute(web, data)

	templateRenderer, err := renderer.NewTemplate()

	if err != nil {
		e.Logger.Fatal(err)
	}
	e.Renderer = templateRenderer
	e.Static("/assets", "assets")

	e.HTTPErrorHandler = func(err error, c echo.Context) {

		var he *echo.HTTPError

		if errors.As(err, &he) {

			if he.Code == http.StatusNotFound {

				_ = handler.NotFoundError(c)
			}
		}

		_ = handler.ServerError(c)
	}

	port, err := strconv.Atoi(cfg.AppPort)

	if err != nil {

		e.Logger.Fatal(err.Error())
	}

	appAddress := fmt.Sprintf("%s:%d", cfg.AppDomain, port)

	e.Logger.Fatal(e.Start(appAddress))
}
