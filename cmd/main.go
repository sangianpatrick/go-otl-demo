package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"path"
	"runtime"
	"strings"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/sangianpatrick/go-otl-demo/config"
	"github.com/sangianpatrick/go-otl-demo/domain/album"
	"github.com/sirupsen/logrus"
	"github.com/uptrace/opentelemetry-go-extra/otellogrus"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"

	// "go.opentelemetry.io/otel/exporters/stdout/stdouttrace"

	_ "github.com/joho/godotenv/autoload"
	"github.com/sangianpatrick/go-otl-demo/monitoring"
)

var cfg *config.Config

func init() {
	cfg = config.Get()
}

func main() {
	// exporter, _ := stdouttrace.New(stdouttrace.WithPrettyPrint())
	motel := monitoring.NewOpenTelemetry(cfg.Application.Name, cfg.Opentelemetry.CollectorHost, cfg.Application.Enviroment)
	motel.Start()

	logger := logrus.New()
	logger.AddHook(otellogrus.NewHook(otellogrus.WithLevels(
		logrus.PanicLevel,
		logrus.FatalLevel,
		logrus.ErrorLevel,
		logrus.WarnLevel,
	)))
	logger.SetReportCaller(true)
	logger.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: time.RFC3339,
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			s := strings.Split(f.Function, ".")
			funcname := s[len(s)-1]
			_, filename := path.Split(f.File)
			filename = fmt.Sprintf("%s:%d", filename, f.Line)
			return funcname, filename
		},
	})

	httpClient := &http.Client{Transport: otelhttp.NewTransport(http.DefaultTransport)}

	router := mux.NewRouter()
	router.Use(otelmux.Middleware(cfg.Application.Name))

	// init domain
	albumRepository := album.NewAlbumRepository(logger, httpClient, cfg.JSONPlaceHolderAPI.Host)
	albumService := album.NewAlbumService(albumRepository)
	album.NewAlbumHTTPHandler(router, albumService)

	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Application.Port),
		Handler: router,
	}

	go func() {
		logger.Infof("application is running on port :%d", cfg.Application.Port)
		httpServer.ListenAndServe()
	}()

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGTERM, syscall.SIGINT)
	<-sigterm

	shutdownCtx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	httpServer.Shutdown(shutdownCtx)
	motel.Stop()
}
