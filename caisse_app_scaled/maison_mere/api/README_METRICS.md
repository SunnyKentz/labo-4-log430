# Prometheus Metrics for Maison Mere API

This document describes the Prometheus metrics implementation for the Maison Mere API service.

## Overview

The Maison Mere API now includes comprehensive Prometheus metrics for monitoring and observability. Metrics are exposed at the `/api/v1/metrics` endpoint.

## Available Metrics

### HTTP Request Metrics

- **`mere_http_requests_total`** (Counter)
  - Labels: `method`, `endpoint`, `status`
  - Description: Total number of HTTP requests

- **`mere_http_request_duration_seconds`** (Histogram)
  - Labels: `method`, `endpoint`
  - Description: Duration of HTTP requests

### Business Metrics

- **`mere_transactions_total`** (Counter)
  - Description: Total number of transactions

- **`mere_sales_total`** (Counter)
  - Description: Total number of sales

- **`mere_refunds_total`** (Counter)
  - Description: Total number of refunds

- **`mere_magasins_total`** (Gauge)
  - Description: Total number of registered stores

- **`mere_notifications_total`** (Counter)
  - Description: Total number of notifications received

- **`mere_alerts_total`** (Gauge)
  - Description: Total number of active alerts

### Error Metrics

- **`mere_errors_total`** (Counter)
  - Labels: `type` (invalid_parameter, database_error, transaction_creation_failed, refund_failed)
  - Description: Total number of errors by type

### Database Metrics

- **`mere_database_operations_total`** (Counter)
  - Labels: `operation`
  - Description: Total number of database operations

- **`mere_database_operation_duration_seconds`** (Histogram)
  - Labels: `operation`
  - Description: Duration of database operations

## Usage

### Accessing Metrics

To view the metrics, make a GET request to:
```
GET /api/v1/metrics
```

Example:
```bash
curl http://localhost:8090/api/v1/metrics
```

### Prometheus Configuration

Add the following to your `prometheus.yml` configuration:

```yaml
scrape_configs:
  - job_name: 'mere-api'
    static_configs:
      - targets: ['localhost:8090']
    metrics_path: '/api/v1/metrics'
    scrape_interval: 15s
```

### Example Queries

Here are some useful PromQL queries for monitoring:

```promql
# Request rate by endpoint
rate(mere_http_requests_total[5m])

# Error rate
rate(mere_errors_total[5m])

# Average response time
histogram_quantile(0.95, rate(mere_http_request_duration_seconds_bucket[5m]))

# Total transactions
mere_transactions_total

# Sales rate
rate(mere_sales_total[5m])

# Refunds rate
rate(mere_refunds_total[5m])

# Number of registered stores
mere_magasins_total

# Notifications received
rate(mere_notifications_total[5m])

# Active alerts
mere_alerts_total
```

## Implementation Details

### Metrics Collection

Metrics are automatically collected through:

1. **HTTP Middleware**: All HTTP requests are automatically instrumented with request count and duration metrics
2. **Business Logic Instrumentation**: Key business operations (transactions, sales, refunds, notifications) are instrumented with specific metrics
3. **Error Tracking**: Errors are categorized and tracked by type

### Files

- `metrics.go`: Defines all Prometheus metrics and helper functions
- `api_data.go`: Contains the metrics endpoint and middleware implementation

### Adding New Metrics

To add new metrics:

1. Define the metric in `metrics.go`
2. Add recording calls in the appropriate handlers
3. Update this documentation

## Monitoring Best Practices

1. **Set up alerts** for high error rates or slow response times
2. **Monitor business metrics** to track application health
3. **Use dashboards** to visualize key metrics
4. **Set retention policies** appropriate for your metrics volume

## Example Alerts

```yaml
groups:
  - name: mere-api
    rules:
      - alert: HighErrorRate
        expr: rate(mere_errors_total[5m]) > 0.1
        for: 2m
        labels:
          severity: warning
        annotations:
          summary: "High error rate detected"
          
      - alert: SlowResponseTime
        expr: histogram_quantile(0.95, rate(mere_http_request_duration_seconds_bucket[5m])) > 1
        for: 2m
        labels:
          severity: warning
        annotations:
          summary: "Slow response time detected"
          
      - alert: HighRefundRate
        expr: rate(mere_refunds_total[5m]) > rate(mere_sales_total[5m]) * 0.1
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: "High refund rate detected"
          
      - alert: NoStoresRegistered
        expr: mere_magasins_total == 0
        for: 1m
        labels:
          severity: critical
        annotations:
          summary: "No stores registered"
``` 