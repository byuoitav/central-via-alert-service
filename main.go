package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"sync"
	"time"

	viadriver "github.com/byuoitav/kramer-driver"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/spf13/pflag"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	var (
		port     int
		username string
		password string
	)

	pflag.IntVarP(&port, "port", "P", 8080, "port to run the server on")
	pflag.StringVarP(&username, "username", "u", "", "username for device")
	pflag.StringVarP(&password, "password", "p", "", "password for device")
	pflag.Parse()

	addr := fmt.Sprintf(":%d", port)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		fmt.Printf("failed to start server: %s\n", err)
		os.Exit(1)
	}

	vias := &sync.Map{}

	cfg := zap.NewProductionConfig()
	cfg.Level.SetLevel(zapcore.DebugLevel)
	zapLog, _ := cfg.Build()

	handlers := Handlers{
		CreateServer: func(addr string) *viadriver.Via {
			if vs, ok := switchers.Load(addr); ok {
				return vs.(*viadriver.Via)
			}

			v := &viadriver.Via{
				Address:  addr,
				Username: username,
				Password: password,
				Logger:   zapLog,
			}

			vias.Store(addr, v)
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
