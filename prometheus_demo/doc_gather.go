package main

import (
	"bytes"
	"fmt"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
	"github.com/prometheus/common/expfmt"
)

func main() {
	reg := prometheus.NewRegistry()
	temp := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "temperature_kelvin",
			Help: "Temperature in Kelvin.",
		},
		[]string{"location"},
	)
	reg.MustRegister(temp)
	temp.WithLabelValues("outside").Set(273.14)
	temp.WithLabelValues("inside").Set(298.44)

	var parser expfmt.TextParser

	text := `
# TYPE humidity_percent gauge
# HELP humidity_percent Humidity in %.
humidity_percent{location="outside"} 45.4
humidity_percent{location="inside"} 33.2
# TYPE temperature_kelvin gauge
# HELP temperature_kelvin Temperature in Kelvin.
temperature_kelvin{location="somewhere else"} 4.5
`

	parseText := func() ([]*dto.MetricFamily, error) {
		parsed, err := parser.TextToMetricFamilies(strings.NewReader(text))
		if err != nil {
			return nil, err
		}
		var result []*dto.MetricFamily
		for _, mf := range parsed {
			result = append(result, mf)
		}
		return result, nil
	}

	gatherers := prometheus.Gatherers{
		reg,
		prometheus.GathererFunc(parseText),
	}

	gathering, err := gatherers.Gather()
	if err != nil {
		fmt.Println(err)
	}

	out := &bytes.Buffer{}
	for _, mf := range gathering {
		if _, err := expfmt.MetricFamilyToText(out, mf); err != nil {
			panic(err)
		}
	}
	fmt.Print(out.String())
	fmt.Println("----------")

	// Note how the temperature_kelvin metric family has been merged from
	// different sources. Now try
	text = `
# TYPE humidity_percent gauge
# HELP humidity_percent Humidity in %.
humidity_percent{location="outside"} 45.4
humidity_percent{location="inside"} 33.2
# TYPE temperature_kelvin gauge
# HELP temperature_kelvin Temperature in Kelvin.
# Duplicate metric:
temperature_kelvin{location="outside"} 265.3
# Missing location label (note that this is undesirable but valid):
temperature_kelvin 4.5
`

	gathering, err = gatherers.Gather()
	if err != nil {
		fmt.Println(err)
	}
	// Note that still as many metrics as possible are returned:
	out.Reset()
	for _, mf := range gathering {
		if _, err := expfmt.MetricFamilyToText(out, mf); err != nil {
			panic(err)
		}
	}
	fmt.Print(out.String())

	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":2112", nil)
}
