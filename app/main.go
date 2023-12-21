package main

import (
	"context"
	"net"
	"os"

	"github.com/kokoichi206-sandbox/url-shortener/config"
	"github.com/kokoichi206-sandbox/url-shortener/handler"
	"github.com/kokoichi206-sandbox/url-shortener/repository/database"
	"github.com/kokoichi206-sandbox/url-shortener/usecase"
	"github.com/kokoichi206-sandbox/url-shortener/util"
	"github.com/kokoichi206-sandbox/url-shortener/util/logger"
	"github.com/opentracing/opentracing-go"
)

const (
	service = "server-template"
)

func main() {
	// config
	cfg := config.New()

	// logger
	logger := logger.NewBasicLogger(os.Stdout, "ubuntu", service)

	// tracer
	tracer, traceCloser, err := util.NewJaegerTracer(cfg.AgentHost, cfg.AgentPort, service)
	if err != nil {
		logger.Errorf(context.Background(), "cannot initialize jaeger tracer: ", err)
	} else {
		defer traceCloser.Close()
		opentracing.SetGlobalTracer(tracer)
	}

	// database
	database, err := database.New(
		cfg.DBDriver, cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword,
		cfg.DBName, cfg.DBSSLMode, logger,
	)
	if err != nil {
		logger.Errorf(context.Background(), "failed to db.New: ", err)
	}

	// usecase
	usecase := usecase.New(database, logger)

	// handler
	h := handler.New(logger, usecase)
	addr := net.JoinHostPort(cfg.ServerHost, cfg.ServerPort)

	// run
	if err := h.Engine.Run(addr); err != nil {
		logger.Critical(context.Background(), "failed to serve http")
	}
}
