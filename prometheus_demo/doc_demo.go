package main

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

func main() {
	// Imagine you have a worker pool and want to count the tasks completed.
	taskCounter := prometheus.NewCounter(prometheus.CounterOpts{
		Subsystem: "worker_pool",
		Name:      "completed_tasks_total",
		Help:      "Total number of tasks completed.",
	})
	// This will register fine.
	if err := prometheus.Register(taskCounter); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("taskCounter registered.")
	}
	// Don't forget to tell the HTTP server about the Prometheus handler.
	// (In a real program, you still need to start the HTTP server...)
	http.Handle("/metrics", promhttp.Handler())

	// Now you can start workers and give every one of them a pointer to
	// taskCounter and let it increment it whenever it completes a task.
	taskCounter.Inc() // This has to happen somewhere in the worker code.

	// But wait, you want to see how individual workers perform. So you need
	// a vector of counters, with one element for each worker.
	taskCounterVec := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Subsystem: "worker_pool",
			Name:      "completed_tasks_total",
			Help:      "Total number of tasks completed.",
		},
		[]string{"worker_id"},
	)

	// Registering will fail because we already have a metric of that name.
	if err := prometheus.Register(taskCounterVec); err != nil {
		fmt.Println("taskCounterVec not registered:", err)
	} else {
		fmt.Println("taskCounterVec registered.")
	}

	// To fix, first unregister the old taskCounter.
	if prometheus.Unregister(taskCounter) {
		fmt.Println("taskCounter unregistered.")
	}

	// Try registering taskCounterVec again.
	if err := prometheus.Register(taskCounterVec); err != nil {
		fmt.Println("taskCounterVec not registered:", err)
	} else {
		fmt.Println("taskCounterVec registered.")
	}
	// Bummer! Still doesn't work.

	// Prometheus will not allow you to ever export metrics with
	// inconsistent help strings or label names. After unregistering, the
	// unregistered metrics will cease to show up in the /metrics HTTP
	// response, but the registry still remembers that those metrics had
	// been exported before. For this example, we will now choose a
	// different name. (In a real program, you would obviously not export
	// the obsolete metric in the first place.)
	taskCounterVec = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Subsystem: "worker_pool",
			Name:      "completed_tasks_by_id",
			Help:      "Total number of tasks completed counted by worker_id.",
		},
		[]string{"worker_id"},
	)
	if err := prometheus.Register(taskCounterVec); err != nil {
		fmt.Println("taskCounterVec not registered:", err)
	} else {
		fmt.Println("taskCounterVec registered.")
	}
	// Finally it worked!

	// The workers have to tell taskCounterVec their id to increment the
	// right element in the metric vector.
	taskCounterVec.WithLabelValues("31").Inc() // Code from worker 31.

	// Each worker could also keep a reference to their own counter element
	// around. Pick the counter at initialization time of the worker.
	myCounter := taskCounterVec.WithLabelValues("31") // From worker 31 initialization code.
	myCounter.Inc()                                   // Somewhere in the code of that worker.

	myCounter1 := taskCounterVec.WithLabelValues("34") // From worker 34 initialization code.
	myCounter1.Inc()                                   // Somewhere in the code of that worker.

	// Note that something like WithLabelValues("42", "spurious arg") would
	// panic (because you have provided too many label values). If you want
	// to get an error instead, use GetMetricWithLabelValues(...) instead.
	notMyCounter, err := taskCounterVec.GetMetricWithLabelValues("42", "spurious arg")
	if err != nil {
		fmt.Println("Worker initialization failed:", err)
	}
	if notMyCounter == nil {
		fmt.Println("notMyCounter is nil.")
	}

	// A different (and somewhat tricky) approach is to use
	// ConstLabels. ConstLabels are pairs of label names and label values
	// that never change. Each worker creates and registers an own Counter
	// instance where the only difference is in the value of the
	// ConstLabels. Those Counters can all be registered because the
	// different ConstLabel values guarantee that each worker will increment
	// a different Counter metric.
	counterOpts := prometheus.CounterOpts{
		Subsystem:   "worker_pool",
		Name:        "completed_tasks",
		Help:        "Total number of tasks completed.",
		ConstLabels: prometheus.Labels{"worker_id": "42"},
	}
	taskCounterForWorker42 := prometheus.NewCounter(counterOpts)
	if err := prometheus.Register(taskCounterForWorker42); err != nil {
		fmt.Println("taskCounterVForWorker42 not registered:", err)
	} else {
		fmt.Println("taskCounterForWorker42 registered.")
	}
	// Obviously, in real code, taskCounterForWorker42 would be a member
	// variable of a worker struct, and the "42" would be retrieved with a
	// GetId() method or something. The Counter would be created and
	// registered in the initialization code of the worker.

	// For the creation of the next Counter, we can recycle
	// counterOpts. Just change the ConstLabels.
	counterOpts.ConstLabels = prometheus.Labels{"worker_id": "2001"}
	taskCounterForWorker2001 := prometheus.NewCounter(counterOpts)
	if err := prometheus.Register(taskCounterForWorker2001); err != nil {
		fmt.Println("taskCounterVForWorker2001 not registered:", err)
	} else {
		fmt.Println("taskCounterForWorker2001 registered.")
	}

	taskCounterForWorker2001.Inc()
	taskCounterForWorker42.Inc()
	taskCounterForWorker2001.Inc()

	// Yet another approach would be to turn the workers themselves into
	// Collectors and register them. See the Collector example for details.
	http.ListenAndServe(":2112", nil)
}
