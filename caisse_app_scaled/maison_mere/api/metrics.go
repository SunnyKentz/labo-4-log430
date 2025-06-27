package api

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// HTTP request metrics
	httpRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "mere_http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "endpoint", "status"},
	)

	httpRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "mere_http_request_duration_seconds",
			Help:    "Duration of HTTP requests",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "endpoint"},
	)

	// Business metrics
	transactionsTotal = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "mere_transactions_total",
			Help: "Total number of transactions",
		},
	)

	salesTotal = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "mere_sales_total",
			Help: "Total number of sales",
		},
	)

	refundsTotal = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "mere_refunds_total",
			Help: "Total number of refunds",
		},
	)

	magasinsTotal = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "mere_magasins_total",
			Help: "Total number of registered stores",
		},
	)

	notificationsTotal = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "mere_notifications_total",
			Help: "Total number of notifications received",
		},
	)

	alertsTotal = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "mere_alerts_total",
			Help: "Total number of active alerts",
		},
	)

	// Error metrics
	errorsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "mere_errors_total",
			Help: "Total number of errors",
		},
		[]string{"type"},
	)

	// Database metrics
	databaseOperationsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "mere_database_operations_total",
			Help: "Total number of database operations",
		},
		[]string{"operation"},
	)

	databaseOperationDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "mere_database_operation_duration_seconds",
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

// RecordTransaction records a transaction metric
func RecordTransaction() {
	transactionsTotal.Inc()
}

// RecordSale records a sale metric
func RecordSale() {
	salesTotal.Inc()
}

// RecordRefund records a refund metric
func RecordRefund() {
	refundsTotal.Inc()
}

// UpdateMagasinsTotal updates the total magasins metric
func UpdateMagasinsTotal(count int) {
	magasinsTotal.Set(float64(count))
}

// RecordNotification records a notification metric
func RecordNotification() {
	notificationsTotal.Inc()
}

// UpdateAlertsTotal updates the total alerts metric
func UpdateAlertsTotal(count int) {
	alertsTotal.Set(float64(count))
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
