package metrics

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
)

var (
	// RequestsPerSecond is a counter incremented per request. Use rate() in Prometheus to compute RPS.
	RequestsPerSecond = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "requests_per_second",
		Help: "Counter incremented once per request (use rate() for RPS)",
	})

	// RequestDuration measures total HTTP request duration per route
	RequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Total HTTP request duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"route"},
	)

	// DeliveryDuration measures time spent in the delivery layer (handlers)
	DeliveryDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "delivery_layer_duration_seconds",
			Help:    "Time spent in the delivery layer (handlers)",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"route"},
	)

	// DomainDuration measures time spent in domain logic within handlers (excluding repo where possible)
	DomainDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "domain_layer_duration_seconds",
			Help:    "Time spent in domain logic within handlers",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"route"},
	)

	// RepositoryDuration measures time spent in repository operations
	RepositoryDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "repository_layer_duration_seconds",
			Help:    "Time spent in repository operations",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"operation"},
	)
)

func Register() {
	// default Go and process collectors
	prometheus.MustRegister(collectors.NewGoCollector())
	prometheus.MustRegister(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}))

	// custom metrics
	prometheus.MustRegister(RequestsPerSecond)
	prometheus.MustRegister(RequestDuration)
	prometheus.MustRegister(DeliveryDuration)
	prometheus.MustRegister(DomainDuration)
	prometheus.MustRegister(RepositoryDuration)
}

func IncRequestsPerSecond() {
	RequestsPerSecond.Inc()
}

func ObserveRequestDuration(route string, d time.Duration) {
	RequestDuration.WithLabelValues(route).Observe(d.Seconds())
}

func ObserveDeliveryDuration(route string, d time.Duration) {
	DeliveryDuration.WithLabelValues(route).Observe(d.Seconds())
}

func ObserveDomainDuration(route string, d time.Duration) {
	DomainDuration.WithLabelValues(route).Observe(d.Seconds())
}

func ObserveRepositoryDuration(operation string, d time.Duration) {
	RepositoryDuration.WithLabelValues(operation).Observe(d.Seconds())
}
