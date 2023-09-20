package cmd

import (
	"fmt"
	"net/http"

	"github.com/spf13/cobra"

	"github.com/openzipkin/zipkin-go"
	zipkinhttp "github.com/openzipkin/zipkin-go/middleware/http"
	"github.com/openzipkin/zipkin-go/model"
	reporterhttp "github.com/openzipkin/zipkin-go/reporter/http"
)

var traceZipkinCmd = &cobra.Command{
	Use:   "zipkin",
	Short: "Send some traces using Zipkin transport",

	Run: func(cmd *cobra.Command, args []string) {
		serviceName, _ := cmd.Flags().GetString("name")
		httpURL, _ := cmd.Flags().GetString("http-collector-url")

		requestURL, _ := cmd.Flags().GetString("http-request-url")

		reporter := reporterhttp.NewReporter(httpURL)
		defer reporter.Close()

		// Local endpoint represent the local service information
		localEndpoint := &model.Endpoint{ServiceName: serviceName}

		// Sampler tells you which traces are going to be sampled or not. In this case we will record 100% (1.00) of traces.
		sampler, err := zipkin.NewCountingSampler(1)
		if err != nil {
			fmt.Println(err)
		}

		tracer, err := zipkin.NewTracer(
			reporter,
			zipkin.WithSampler(sampler),
			zipkin.WithLocalEndpoint(localEndpoint),
		)
		if err != nil {
			fmt.Println(err)
		}

		span := tracer.StartSpan("hello")
		span.Finish()

		if requestURL == "" {
			return
		}

		client, _ := zipkinhttp.NewClient(tracer, zipkinhttp.ClientTrace(true))
		req, _ := http.NewRequest("GET", requestURL, nil)
		res, _ := client.DoWithAppSpan(req, "http-client")
		res.Body.Close()

		fmt.Println("Done!")
	},
}

func init() {
	traceCmd.AddCommand(traceZipkinCmd)

	traceZipkinCmd.Flags().StringP("name", "n", "tracer-cli", "Service Name")
	traceZipkinCmd.Flags().StringP("http-collector-url", "c", "http://localhost:9411/api/v2/spans", "The URL for communicating thrift or JSON via HTTP")
	traceZipkinCmd.Flags().StringP("http-request-url", "r", "", "Make HTTP GET to this URL")
}
