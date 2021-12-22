package main

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"runtime"

	"github.com/prometheus/client_golang/prometheus"
)

type DBStat struct {
	OpenConnections int
}

// primaryDB and secondaryDB represent two example *sql.DB connections we want to instrument.
var primaryDB, secondaryDB interface {
	Stats() *DBStat
}

type prDB struct {
}

func (*prDB) Stats() *DBStat {
	return &DBStat{OpenConnections: 20}
}

type secDB struct {
}

func (*secDB) Stats() struct{ OpenConnections int } {
	return struct{ OpenConnections int }{OpenConnections: 30}
}
func main() {
	prDB := &prDB{}
	if err := prometheus.Register(prometheus.NewGaugeFunc(
		prometheus.GaugeOpts{
			Namespace:   "myapp_mysql",
			Name:        "connections_open",
			Help:        "Number of mysql connections open.",
			ConstLabels: prometheus.Labels{"destination": "primary"},
		},
		func() float64 { return float64(prDB.Stats().OpenConnections) },
	)); err == nil {
		fmt.Println(`GaugeFunc 'connections_open' for primary DB connection registered with labels {destination="primary"}`)
	}
	secDB := &secDB{}
	if err := prometheus.Register(prometheus.NewGaugeFunc(
		prometheus.GaugeOpts{
			Namespace:   "myapp_mysql",
			Name:        "connections_open",
			Help:        "Number of mysql connections open.",
			ConstLabels: prometheus.Labels{"destination": "secondary"},
		},
		func() float64 { return float64(secDB.Stats().OpenConnections) },
	)); err == nil {
		fmt.Println(`GaugeFunc 'connections_open' for secondary DB connection registered with labels {destination="secondary"}`)
	}
	// Note that we can register more than once GaugeFunc with same metric name
	// as long as their const labels are consistent.

	if err := prometheus.Register(prometheus.NewGaugeFunc(
		prometheus.GaugeOpts{
			Subsystem: "myapp_runtime",
			Name:      "goroutines_count",
			Help:      "Number of goroutines that currently exist.",
		},
		func() float64 { return float64(runtime.NumGoroutine()) },
	)); err == nil {
		fmt.Println("GaugeFunc 'goroutines_count' registered.")
	}
	// Note that the count of goroutines is a gauge (and not a counter) as
	// it can go up and down.
	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":2112", nil)
}
