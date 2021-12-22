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
		"A histogram of the HTTP request durations.",
		[]string{"code", "method"},
		prometheus.Labels{"owner": "example"},
	)

	// Create a constant histogram from values we got from a 3rd party telemetry system.
	h := prometheus.MustNewConstHistogram(
		desc,
		4711, 403.34,
		map[float64]uint64{25: 121, 50: 2403, 100: 3221, 200: 4233},
		"200", "get",
	)

	// Just for demonstration, let's check the state of the histogram by
	// (ab)using its Write method (which is usually only used by Prometheus
	// internally).
	metric := &dto.Metric{}
	h.Write(metric)
	fmt.Println(proto.MarshalTextString(metric))

}
