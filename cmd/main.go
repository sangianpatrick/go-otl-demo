package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/sangianpatrick/go-otl-demo/config"
	"github.com/sangianpatrick/go-otl-demo/domain/album"
	"github.com/sirupsen/logrus"
	"github.com/uptrace/opentelemetry-go-extra/otellogrus"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"

	"go.opentelemetry.io/otel/attribute"
	// "go.opentelemetry.io/otel/exporters/stdout/stdouttrace"

	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.10.0"

	_ "github.com/joho/godotenv/autoload"
)

var cfg *config.Config

func init() {
	cfg = config.Get()
}

func main() {
	// exporter, _ := stdouttrace.New(stdouttrace.WithPrettyPrint())
	exporter, _ := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(cfg.Jaeger.Host)))
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(cfg.Application.Name),
			attribute.String("environment", cfg.Application.Enviroment),
		)),
	)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(
		propagation.NewCompositeTextMapPropagator(
			propagation.TraceContext{},
			propagation.Baggage{},
		),
	)

	logger := logrus.New()
	logger.AddHook(otellogrus.NewHook(otellogrus.WithLevels(
		logrus.PanicLevel,
		logrus.FatalLevel,
		logrus.ErrorLevel,
		logrus.WarnLevel,
	)))

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
	tp.Shutdown(shutdownCtx)
}
