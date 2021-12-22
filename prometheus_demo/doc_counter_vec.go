package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
)

func main() {
	httpReqs := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "myapp_http_requests_total",
			Help: "How many HTTP requests processed, partitioned by status code and HTTP method.",
		},
		[]string{"code", "method"},
	)
	prometheus.MustRegister(httpReqs)

	httpReqs.WithLabelValues("404", "POST").Add(42)

	// If you have to access the same set of labels very frequently, it
	// might be good to retrieve the metric only once and keep a handle to
	// it. But beware of deletion of that metric, see below!
	m := httpReqs.WithLabelValues("200", "GET")
	for i := 0; i < 1000000; i++ {
		m.Inc()
	}
	// Delete a metric from the vector. If you have previously kept a handle
	// to that metric (as above), future updates via that handle will go
	// unseen (even if you re-create a metric with the same label set
	// later).
	httpReqs.DeleteLabelValues("200", "GET")
	// Same thing with the more verbose Labels syntax.
	httpReqs.Delete(prometheus.Labels{"method": "GET", "code": "200"})
	// 下面果真没啥用
	m.Inc()
	// 可是下面的管用了呢
	m = httpReqs.WithLabelValues("200", "GET")
	for i := 0; i < 1000000; i++ {
		m.Inc()
	}
	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(":2112", nil))
}
