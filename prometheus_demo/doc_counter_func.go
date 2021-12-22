package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
)

var accessTimes float64

func main() {
	httpReqs := prometheus.NewCounterFunc(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "How many HTTP requests processed, partitioned by status code and HTTP method.",
		},
		func() float64 {
			accessTimes++
			log.Printf("accessTimes: %f", accessTimes)
			return accessTimes
		},
	)
	prometheus.MustRegister(httpReqs)

	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(":2112", nil))
}
