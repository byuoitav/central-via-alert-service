package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"sync"
	//"time"

	//viadriver "github.com/byuoitav/kramer-driver"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/spf13/pflag"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type AlertServer struct {
	Address  string
	Username string
	Password string
	Logger   *zap.SugaredLogger
}

func main() {
	var (
		port     int
		username string
		password string
		logLevel int8
	)

	pflag.IntVarP(&port, "port", "P", 8080, "port to run the server on")
	pflag.StringVarP(&username, "username", "u", "", "username for database")
	pflag.StringVarP(&password, "password", "p", "", "password for database")
	pflag.Int8VarP(&logLevel, "log-level", "L", 0, "Level to log at. Provided by zap logger: https://godoc.org/go.uber.org/zap/zapcore")
	pflag.Parse()

	// Build out the Logger
	config := zap.Config{
		Level:       zap.NewAtomicLevelAt(zapcore.Level(logLevel)),
		Development: false,
		Sampling: &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		},
		Encoding: "json",
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:        "@",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "caller",
			MessageKey:     "msg",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.StringDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		},
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
	}

	plain, err := config.Build()
	if err != nil {
		fmt.Printf("unable to build logger you foolish mortal: %s", err)
		os.Exit(1)
	}

	logger := plain.Sugar()

	addr := fmt.Sprintf(":%d", port)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		fmt.Printf("failed to start server: %s\n", err)
		os.Exit(1)
	}

	logger.Info("Starting Service.....")

	alertService := &sync.Map{}

	handlers := Handlers{
		CreateServer: func(addr string) *AlertServer {
			if vs, ok := alertService.Load(addr); ok {
				return vs.(*AlertServer)
			}

			v := &AlertServer{
				Address:  addr,
				Username: username,
				Password: password,
				Logger:   logger,
			}

			alertService.Store(addr, v)
			return v
		},
	}

	e := echo.New()

	e.Pre(middleware.RemoveTrailingSlash())

	api := e.Group("/api/v1")
	handlers.RegisterRoutes(api)

	log.Printf("Server started on %v", lis.Addr())
	if err := e.Server.Serve(lis); err != nil {
		log.Printf("unable to serve: %s", err)
	}
}
