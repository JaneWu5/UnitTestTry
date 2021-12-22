package main

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

func main() {
	reqCounter := prometheus.NewCounter(prometheus.CounterOpts{
		Name: "myapp_requests_total",
		Help: "The total number of requests served.",
	})
	prometheus.MustRegister(reqCounter)
	reqCounter.Inc()
	if err := prometheus.Register(reqCounter); err != nil {
		if are, ok := err.(prometheus.AlreadyRegisteredError); ok {
			fmt.Println("A counter for that metric has been registered before.")
			// Use the old counter from now on.
			reqCounter = are.ExistingCollector.(prometheus.Counter)
		} else {
			// Something else went wrong!
			panic(err)
		}
	}
	reqCounter.Inc()

	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":2112", nil)
}
