package main

import (
	"github.com/prometheus/client_golang/prometheus/collectors"
	"log"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// ClusterManager is an example for a system that might have been built without
// Prometheus in mind. It models a central manager of jobs running in a
// cluster. Thus, we implement a custom Collector called
// ClusterManagerCollector, which collects information from a ClusterManager
// using its provided methods and turns them into Prometheus Metrics for
// collection.
//
// An additional challenge is that multiple instances of the ClusterManager are
// run within the same binary, each in charge of a different zone. We need to
// make use of wrapping Registerers to be able to register each
// ClusterManagerCollector instance with Prometheus.
type ClusterManager struct {
	Zone string
	// Contains many more fields not listed in this example.
}

// ReallyExpensiveAssessmentOfTheSystemState is a mock for the data gathering a
// real cluster manager would have to do. Since it may actually be really
// expensive, it must only be called once per collection. This implementation,
// obviously, only returns some made-up data.
func (c *ClusterManager) ReallyExpensiveAssessmentOfTheSystemState() (
	oomCountByHost map[string]int64, ramUsageByHost map[string]float64,
) {
	// Just example fake data.
	oomCountByHost = map[string]int64{
		"foo.example.org": time.Now().Unix(),
		"bar.example.org": time.Now().Unix(),
	}
	ramUsageByHost = map[string]float64{
		"foo.example.org": 6.023e23,
		"bar.example.org": 3.14,
	}
	return
}

// ClusterManagerCollector implements the Collector interface.
type ClusterManagerCollector struct {
	ClusterManager *ClusterManager
}

// Descriptors used by the ClusterManagerCollector below.
var (
	oomCountDesc = prometheus.NewDesc(
		"clustermanager_oom_crashes_total",
		"Number of OOM crashes.",
		[]string{"host"}, nil,
	)
	ramUsageDesc = prometheus.NewDesc(
		"clustermanager_ram_usage_bytes",
		"RAM usage as reported to the cluster manager.",
		[]string{"host"}, nil,
	)
)

// Describe is implemented with DescribeByCollect. That's possible because the
// Collect method will always return the same two metrics with the same two
// descriptors.
func (cc ClusterManagerCollector) Describe(ch chan<- *prometheus.Desc) {
	prometheus.DescribeByCollect(cc, ch)
}

// Collect first triggers the ReallyExpensiveAssessmentOfTheSystemState. Then it
// creates constant metrics for each host on the fly based on the returned data.
//
// Note that Collect could be called concurrently, so we depend on
// ReallyExpensiveAssessmentOfTheSystemState to be concurrency-safe.
func (cc ClusterManagerCollector) Collect(ch chan<- prometheus.Metric) {
	oomCountByHost, ramUsageByHost := cc.ClusterManager.ReallyExpensiveAssessmentOfTheSystemState()
	for host, oomCount := range oomCountByHost {
		ch <- prometheus.MustNewConstMetric(
			oomCountDesc,
			prometheus.CounterValue,
			float64(oomCount),
			host,
		)
	}
	for host, ramUsage := range ramUsageByHost {
		ch <- prometheus.MustNewConstMetric(
			ramUsageDesc,
			prometheus.GaugeValue,
			ramUsage,
			host,
		)
	}
}

// NewClusterManager first creates a Prometheus-ignorant ClusterManager
// instance. Then, it creates a ClusterManagerCollector for the just created
// ClusterManager. Finally, it registers the ClusterManagerCollector with a
// wrapping Registerer that adds the zone as a label. In this way, the metrics
// collected by different ClusterManagerCollectors do not collide.
func NewClusterManager(zone string, reg prometheus.Registerer) *ClusterManager {
	c := &ClusterManager{
		Zone: zone,
	}
	cc := ClusterManagerCollector{ClusterManager: c}
	prometheus.WrapRegistererWith(prometheus.Labels{"zone": zone}, reg).MustRegister(cc)
	return c
}

func main() {
	// Since we are dealing with custom Collector implementations, it might
	// be a good idea to try it out with a pedantic registry.
	reg := prometheus.NewPedanticRegistry()

	// Construct cluster managers. In real code, we would assign them to
	// variables to then do something with them.
	NewClusterManager("db", reg)
	NewClusterManager("ca", reg)

	// Add the standard process and Go metrics to the custom registry.
	reg.MustRegister(
		collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
		collectors.NewGoCollector(),
	)

	http.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{}))
	log.Fatal(http.ListenAndServe(":2112", nil))
}
