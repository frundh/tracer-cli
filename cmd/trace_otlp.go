package cmd

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"strings"

	"github.com/spf13/cobra"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	sdkresource "go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
)

var traceOtlpCmd = &cobra.Command{
	Use:   "otlp",
	Short: "Send some traces using OpenTelemetry transport",

	Run: func(cmd *cobra.Command, args []string) {

		serviceName, _ := cmd.Flags().GetString("name")
		grpcUrl, _ := cmd.Flags().GetString("otlp-grpc-url")

		//httpUrl, _ := cmd.Flags().GetString("otlp-http-url")
		requestUrl, _ := cmd.Flags().GetString("http-request-url")

		url, err := url.Parse(grpcUrl)
		if err != nil {
			log.Fatal(err)
		}

		ctx := context.Background()
		opts := []otlptracegrpc.Option{
			otlptracegrpc.WithEndpoint(url.Host),
		}

		if strings.ToLower(url.Scheme) == "http" {
			opts = append(opts, otlptracegrpc.WithInsecure())
		}

		exporter, err := otlptracegrpc.New(ctx, opts...)
		if err != nil {
			log.Fatal(err)
		}

		res, err := sdkresource.New(
			ctx,
			sdkresource.WithHost(),
			sdkresource.WithOS(),
			sdkresource.WithAttributes(semconv.ServiceName(serviceName)),
		)
		if err != nil {
			log.Fatal(err)
		}

		bsp := sdktrace.NewBatchSpanProcessor(exporter)

		tp := sdktrace.NewTracerProvider(
			sdktrace.WithSampler(sdktrace.AlwaysSample()),
			sdktrace.WithSpanProcessor(bsp),
			sdktrace.WithResource(res),
		)

		otel.SetTracerProvider(tp)
		defer func() { _ = tp.Shutdown(ctx) }()

		tr := otel.Tracer("example")
		_, span := tr.Start(ctx, "example-span")
		defer span.End()

		// if requestURL is not empty, make a request to it
		if requestUrl != "" {
			_, _ = otelhttp.Get(ctx, requestUrl)
		}

		fmt.Println("Done!")
	},
}

func init() {
	traceCmd.AddCommand(traceOtlpCmd)

	traceOtlpCmd.Flags().StringP("name", "n", "tracer-cli", "Service Name")
	traceOtlpCmd.Flags().StringP("otlp-grpc-url", "u", "http://localhost:4317", "The URL for communicating otlp via gRPC")
	traceOtlpCmd.Flags().StringP("otlp-http-url", "c", "http://localhost:4318", "The URL for communicating otlp via HTTP")
	traceOtlpCmd.Flags().StringP("http-request-url", "r", "", "Make HTTP GET to this URL")
}
