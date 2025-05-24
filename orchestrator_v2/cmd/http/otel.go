package main

import (
	"context"
	"errors"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/log/global"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/log"

	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/trace"
)

var (
	ENDPOINT = "observability" // "localhost"
)

func SetupOTelSDK(ctx context.Context) (shutdown func(context.Context) error, err error) {
	var shutdownFuncs []func(context.Context) error

	shutdown = func(ctx context.Context) error {
		var err error

		for _, fn := range shutdownFuncs {
			err = errors.Join(err, fn(ctx))
		}

		shutdownFuncs = nil
		return err
	}

	handleErr := func(inError error) {
		err = errors.Join(err, inError)
	}

	prop := newPropagator()
	otel.SetTextMapPropagator(prop)

	tracerProvider, err := newTracerProvider()
	if err != nil {
		handleErr(err)
		return
	}

	shutdownFuncs = append(shutdownFuncs, tracerProvider.Shutdown)
	otel.SetTracerProvider(tracerProvider)

	meterProvider, err := newMeterProvider()
	if err != nil {
		handleErr(err)
		return
	}
	shutdownFuncs = append(shutdownFuncs, meterProvider.Shutdown)
	otel.SetMeterProvider(meterProvider)

	// Set up logger provider.
	loggerProvider, err := newLoggerProvider()
	if err != nil {
		handleErr(err)
		return
	}
	shutdownFuncs = append(shutdownFuncs, loggerProvider.Shutdown)
	global.SetLoggerProvider(loggerProvider)

	return
}

func newPropagator() propagation.TextMapPropagator {
	return propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	)
}

func newTracerProvider() (*trace.TracerProvider, error) {

	exporter, err := otlptracegrpc.New(context.Background(),
		otlptracegrpc.WithEndpoint(ENDPOINT+":4317"), // adjust this to your setup
		// otlptracegrpc.WithDialOption(grpc.WithBlock()),
		otlptracegrpc.WithInsecure(), // if you're not using TLS
	)

	if err != nil {
		return nil, err
	}

	tracerProvider := trace.NewTracerProvider(
		trace.WithBatcher(
			exporter,
			trace.WithBatchTimeout(time.Second*5),
		),
	)

	return tracerProvider, nil
}

func newMeterProvider() (*metric.MeterProvider, error) {
	metricExporter, err := otlpmetricgrpc.New(context.Background(),
		otlpmetricgrpc.WithDialOption(),
		otlpmetricgrpc.WithEndpoint(ENDPOINT+":4317"),
		otlpmetricgrpc.WithInsecure(),
	)

	if err != nil {
		return nil, err
	}

	meterProvider := metric.NewMeterProvider(
		metric.WithReader(
			metric.NewPeriodicReader(
				metricExporter,
				metric.WithInterval(time.Minute*1),
			),
		),
	)
	return meterProvider, nil
}

func newLoggerProvider() (*log.LoggerProvider, error) {

	logExporter, err := otlploggrpc.New(context.Background(),
		otlploggrpc.WithEndpoint(ENDPOINT+":4317"), // Adjust this
		otlploggrpc.WithInsecure(),                 // Remove if using TLS
		otlploggrpc.WithDialOption(),
	)

	if err != nil {
		return nil, err
	}

	loggerProvider := log.NewLoggerProvider(
		log.WithProcessor(
			log.NewBatchProcessor(
				logExporter,
				log.WithExportInterval(time.Second*5),
			),
		),
	)
	return loggerProvider, nil
}
