package api

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// HTTP request metrics
	httpRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "magasin_http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "endpoint", "status"},
	)

	httpRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "magasin_http_request_duration_seconds",
			Help:    "Duration of HTTP requests",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "endpoint"},
	)

	// Business metrics
	productsTotal = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "magasin_products_total",
			Help: "Total number of products in inventory",
		},
	)

	transactionsTotal = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "magasin_transactions_total",
			Help: "Total number of transactions",
		},
	)

	salesTotal = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "magasin_sales_total",
			Help: "Total number of sales",
		},
	)

	refundsTotal = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "magasin_refunds_total",
			Help: "Total number of refunds",
		},
	)

	cartOperationsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "magasin_cart_operations_total",
			Help: "Total number of cart operations",
		},
		[]string{"operation"},
	)

	// Error metrics
	errorsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "magasin_errors_total",
			Help: "Total number of errors",
		},
		[]string{"type"},
	)

	// Database metrics
	databaseOperationsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "magasin_database_operations_total",
			Help: "Total number of database operations",
		},
		[]string{"operation"},
	)

	databaseOperationDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "magasin_database_operation_duration_seconds",
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

// RecordCartOperation records a cart operation metric
func RecordCartOperation(operation string) {
	cartOperationsTotal.WithLabelValues(operation).Inc()
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
