package main

import (
	"fmt"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
)

func main() {
	desc := prometheus.NewDesc(
		"temperature_kelvin",
		"Current temperature in Kelvin.",
		nil, nil,
	)

	// Create a constant gauge from values we got from an external
	// temperature reporting system. Those values are reported with a slight
	// delay, so we want to add the timestamp of the actual measurement.
	temperatureReportedByExternalSystem := 298.15
	timeReportedByExternalSystem := time.Date(2009, time.November, 10, 23, 0, 0, 12345678, time.UTC)
	s := prometheus.NewMetricWithTimestamp(
		timeReportedByExternalSystem,
		prometheus.MustNewConstMetric(
			desc, prometheus.GaugeValue, temperatureReportedByExternalSystem,
		),
	)

	// Just for demonstration, let's check the state of the gauge by
	// (ab)using its Write method (which is usually only used by Prometheus
	// internally).
	metric := &dto.Metric{}
	s.Write(metric)
	fmt.Println(proto.MarshalTextString(metric))
}
