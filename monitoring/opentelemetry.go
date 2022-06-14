package monitoring

import (
	"context"
	"fmt"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/metric/global"
	"go.opentelemetry.io/otel/propagation"
	controller "go.opentelemetry.io/otel/sdk/metric/controller/basic"
	processor "go.opentelemetry.io/otel/sdk/metric/processor/basic"
	"go.opentelemetry.io/otel/sdk/metric/selector/simple"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.10.0"
)

type OpenTelemetry interface {
	Start()
	Stop()
}

type openTelemetryAdapter struct {
	ctx            context.Context
	service        string
	environment    string
	endpoint       string
	resource       *resource.Resource
	controller     *controller.Controller
	tracerExporter *otlptrace.Exporter
}

func NewOpenTelemetry(service string, endpoint string, environment string) OpenTelemetry {
	ctx := context.Background()
	res, _ := resource.New(ctx,
		resource.WithFromEnv(),
		resource.WithProcess(),
		resource.WithTelemetrySDK(),
		resource.WithHost(),
		resource.WithSchemaURL(semconv.SchemaURL),
		resource.WithAttributes(
			// the service name used to display traces in backends,
			semconv.ServiceNameKey.String(service),
			attribute.String("environment", environment),
		),
	)

	return &openTelemetryAdapter{
		ctx:         context.Background(),
		resource:    res,
		service:     service,
		environment: environment,
		endpoint:    endpoint,
	}
}

func (om *openTelemetryAdapter) onMetric() (err error) {
	client := otlpmetricgrpc.NewClient(
		otlpmetricgrpc.WithInsecure(),
		otlpmetricgrpc.WithEndpoint(om.endpoint),
	)
	exporter, err := otlpmetric.New(om.ctx, client)
	if err != nil {
		return err
	}

	controller := controller.New(
		processor.NewFactory(
			simple.NewWithHistogramDistribution(),
			exporter,
		),
		controller.WithExporter(exporter),
		controller.WithCollectPeriod(time.Second*5),
	)

	global.SetMeterProvider(controller)

	om.controller = controller

	return controller.Start(om.ctx)
}

func (om *openTelemetryAdapter) onTrace() (err error) {
	fmt.Println(om.endpoint)
	client := otlptracegrpc.NewClient(
		otlptracegrpc.WithInsecure(),
		otlptracegrpc.WithEndpoint(om.endpoint),
		otlptracegrpc.WithDialOption())

	exporter, err := otlptrace.New(om.ctx, client)
	if err != nil {
		return err
	}

	bsp := sdktrace.NewBatchSpanProcessor(exporter)
	provider := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithResource(om.resource),
		sdktrace.WithSpanProcessor(bsp),
	)

	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	otel.SetTracerProvider(provider)

	om.tracerExporter = exporter

	return
}

func (om *openTelemetryAdapter) Start() {
	if err := om.onMetric(); err != nil {
		fmt.Printf("monitoring metric: %s", err.Error())
	}

	if err := om.onTrace(); err != nil {
		fmt.Printf("monitoring tracer: %s", err.Error())
	}
}

func (om *openTelemetryAdapter) Stop() {
	cxt, cancel := context.WithTimeout(om.ctx, time.Second*2)
	defer cancel()
	if err := om.tracerExporter.Shutdown(cxt); err != nil {
		otel.Handle(err)
	}
	// pushes any last exports to the receiver
	if err := om.controller.Stop(cxt); err != nil {
		otel.Handle(err)
	}
}
