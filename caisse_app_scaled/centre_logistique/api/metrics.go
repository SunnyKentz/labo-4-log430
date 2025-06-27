package api

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// HTTP request metrics
	httpRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "logistique_http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "endpoint", "status"},
	)

	httpRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "logistique_http_request_duration_seconds",
			Help:    "Duration of HTTP requests",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "endpoint"},
	)

	// Business metrics
	productsTotal = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "logistique_products_total",
			Help: "Total number of products in inventory",
		},
	)

	commandsTotal = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "logistique_commands_total",
			Help: "Total number of commands",
		},
	)

	commandsAccepted = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "logistique_commands_accepted_total",
			Help: "Total number of accepted commands",
		},
	)

	commandsRefused = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "logistique_commands_refused_total",
			Help: "Total number of refused commands",
		},
	)

	commandsPending = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "logistique_commands_pending",
			Help: "Number of pending commands",
		},
	)

	// Error metrics
	errorsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "logistique_errors_total",
			Help: "Total number of errors",
		},
		[]string{"type"},
	)

	// Database metrics
	databaseOperationsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "logistique_database_operations_total",
			Help: "Total number of database operations",
		},
		[]string{"operation"},
	)

	databaseOperationDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "logistique_database_operation_duration_seconds",
			Help:    "Duration of database operations",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"operation"},
	)
)

// RecordHTTPRequest records an HTTP request metric
func RecordHTTPRequest(method, endpoint, status string) {
	httpRequestsTotal.WithLabelValues(method, endpoint, status).Inc()
}

// RecordHTTPRequestDuration records the duration of an HTTP request
func RecordHTTPRequestDuration(method, endpoint string, duration float64) {
	httpRequestDuration.WithLabelValues(method, endpoint).Observe(duration)
}

// UpdateProductsTotal updates the total products metric
func UpdateProductsTotal(count int) {
	productsTotal.Set(float64(count))
}

// RecordCommand records a command metric
func RecordCommand() {
	commandsTotal.Inc()
}

// RecordCommandAccepted records an accepted command metric
func RecordCommandAccepted() {
	commandsAccepted.Inc()
}

// RecordCommandRefused records a refused command metric
func RecordCommandRefused() {
	commandsRefused.Inc()
}

// UpdateCommandsPending updates the pending commands metric
func UpdateCommandsPending(count int) {
	commandsPending.Set(float64(count))
}

// RecordError records an error metric
func RecordError(errorType string) {
	errorsTotal.WithLabelValues(errorType).Inc()
}

// RecordDatabaseOperation records a database operation metric
func RecordDatabaseOperation(operation string) {
	databaseOperationsTotal.WithLabelValues(operation).Inc()
}

// RecordDatabaseOperationDuration records the duration of a database operation
func RecordDatabaseOperationDuration(operation string, duration float64) {
	databaseOperationDuration.WithLabelValues(operation).Observe(duration)
}
