package main

import (
	"fmt"

	"github.com/golang/protobuf/proto"
	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
)

func main() {
	desc := prometheus.NewDesc(
		"http_request_duration_seconds",
		"A summary of the HTTP request durations.",
		[]string{"code", "method"},
		prometheus.Labels{"owner": "example"},
	)

	// Create a constant summary from values we got from a 3rd party telemetry system.
	s := prometheus.MustNewConstSummary(
		desc,
		4711, 403.34,
		map[float64]float64{0.5: 42.3, 0.9: 323.3},
		"200", "get",
	)

	// Just for demonstration, let's check the state of the summary by
	// (ab)using its Write method (which is usually only used by Prometheus
	// internally).
	metric := &dto.Metric{}
	s.Write(metric)
	fmt.Println(proto.MarshalTextString(metric))

}
