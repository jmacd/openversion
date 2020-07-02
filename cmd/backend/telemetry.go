package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"go.opentelemetry.io/otel/api/global"
	"go.opentelemetry.io/otel/exporters/otlp"
	"go.opentelemetry.io/otel/sdk/metric/controller/push"
	"go.opentelemetry.io/otel/sdk/metric/selector/simple"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

func initProviders() *push.Controller {
	collectorAddr, ok := os.LookupEnv("OTEL_RECEIVER_ENDPOINT")
	if !ok {
		collectorAddr = fmt.Sprint(otlp.DefaultCollectorHost, ":", otlp.DefaultCollectorPort)
	}
	exporter, err := otlp.NewExporter(otlp.WithAddress(collectorAddr), otlp.WithInsecure())

	if err != nil {
		log.Fatal(err)
	}

	tp, err := sdktrace.NewProvider(sdktrace.WithConfig(sdktrace.Config{DefaultSampler: sdktrace.AlwaysSample()}),
		sdktrace.WithSyncer(exporter))
	if err != nil {
		log.Fatal(err)
	}

	global.SetTraceProvider(tp)

	pusher := push.New(
		simple.NewWithInexpensiveDistribution(),
		exporter,
		push.WithPeriod(2*time.Second),
	)
	global.SetMeterProvider(pusher.Provider())
	pusher.Start()

	return pusher
}
