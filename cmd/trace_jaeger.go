package cmd

import (
	"fmt"
	"net/http"

	"github.com/spf13/cobra"

	opentracing "github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/uber/jaeger-lib/metrics"

	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	jaegerlog "github.com/uber/jaeger-client-go/log"
)

var traceJaegerCmd = &cobra.Command{
	Use:   "jaeger",
	Short: "Send some traces using Jaeger transport",

	Run: func(cmd *cobra.Command, args []string) {

		serviceName, _ := cmd.Flags().GetString("name")
		udpHost, _ := cmd.Flags().GetString("thrift-udp-host")
		httpURL, _ := cmd.Flags().GetString("thrift-http-url")

		requestURL, _ := cmd.Flags().GetString("http-request-url")

		cfg := jaegercfg.Configuration{
			ServiceName: serviceName,
			Sampler: &jaegercfg.SamplerConfig{
				Type:  jaeger.SamplerTypeConst,
				Param: 1,
			},
			Reporter: &jaegercfg.ReporterConfig{
				LogSpans: true,
			},
		}

		if httpURL != "" {
			cfg.Reporter.CollectorEndpoint = httpURL
		} else {
			cfg.Reporter.LocalAgentHostPort = udpHost
		}

		jLogger := jaegerlog.StdLogger
		jMetricsFactory := metrics.NullFactory

		tracer, closer, _ := cfg.NewTracer(
			jaegercfg.Logger(jLogger),
			jaegercfg.Metrics(jMetricsFactory),
		)

		opentracing.SetGlobalTracer(tracer)
		defer closer.Close()

		//tracer := opentracing.GlobalTracer()

		span := tracer.StartSpan("hello")
		fmt.Println("hello")
		span.Finish()

		parentSpan := tracer.StartSpan("parent")
		defer parentSpan.Finish()

		childSpan := tracer.StartSpan(
			"child",
			opentracing.ChildOf(parentSpan.Context()),
		)
		defer childSpan.Finish()

		if requestURL == "" {
			return
		}

		clientSpan := tracer.StartSpan("http-client")
		defer clientSpan.Finish()

		url := requestURL
		req, _ := http.NewRequest("GET", url, nil)

		ext.SpanKindRPCClient.Set(clientSpan)
		ext.HTTPUrl.Set(clientSpan, url)
		ext.HTTPMethod.Set(clientSpan, "GET")

		tracer.Inject(clientSpan.Context(), opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(req.Header))
		resp, _ := http.DefaultClient.Do(req)
		fmt.Println(resp.Status)
	},
}

func init() {
	traceCmd.AddCommand(traceJaegerCmd)

	traceJaegerCmd.Flags().StringP("name", "n", "tracer-cli", "Service Name")
	traceJaegerCmd.Flags().StringP("thrift-udp-host", "u", "localhost:6831", "The hostname:port for communicating jaeger.thrift via UDP")
	traceJaegerCmd.Flags().StringP("thrift-http-url", "c", "", "The URL for communicating jaeger.thrift via HTTP, e.g. 'http://jaeger-collector:14268/api/traces'")
	traceJaegerCmd.Flags().StringP("http-request-url", "r", "", "Make HTTP GET to this URL")
}
