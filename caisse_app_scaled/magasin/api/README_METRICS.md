# Prometheus Metrics for Magasin API

This document describes the Prometheus metrics implementation for the Magasin API service.

## Overview

The Magasin API now includes comprehensive Prometheus metrics for monitoring and observability. Metrics are exposed at the `/api/v1/metrics` endpoint.

## Available Metrics

### HTTP Request Metrics

- **`magasin_http_requests_total`** (Counter)
  - Labels: `method`, `endpoint`, `status`
  - Description: Total number of HTTP requests

- **`magasin_http_request_duration_seconds`** (Histogram)
  - Labels: `method`, `endpoint`
  - Description: Duration of HTTP requests

### Business Metrics

- **`magasin_products_total`** (Gauge)
  - Description: Total number of products in inventory

- **`magasin_transactions_total`** (Counter)
  - Description: Total number of transactions

- **`magasin_sales_total`** (Counter)
  - Description: Total number of sales

- **`magasin_refunds_total`** (Counter)
  - Description: Total number of refunds

- **`magasin_cart_operations_total`** (Counter)
  - Labels: `operation` (add, remove)
  - Description: Total number of cart operations

### Error Metrics

- **`magasin_errors_total`** (Counter)
  - Labels: `type` (database_error, invalid_parameter, product_not_found, sale_failed, refund_failed)
  - Description: Total number of errors by type

### Database Metrics

- **`magasin_database_operations_total`** (Counter)
  - Labels: `operation`
  - Description: Total number of database operations

- **`magasin_database_operation_duration_seconds`** (Histogram)
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
curl http://localhost:8080/api/v1/metrics
```

### Prometheus Configuration

Add the following to your `prometheus.yml` configuration:

```yaml
scrape_configs:
  - job_name: 'magasin-api'
    static_configs:
      - targets: ['localhost:8080']
    metrics_path: '/api/v1/metrics'
    scrape_interval: 15s
```

### Example Queries

Here are some useful PromQL queries for monitoring:

```promql
# Request rate by endpoint
rate(magasin_http_requests_total[5m])

# Error rate
rate(magasin_errors_total[5m])

# Average response time
histogram_quantile(0.95, rate(magasin_http_request_duration_seconds_bucket[5m]))

# Total sales
magasin_sales_total

# Products in inventory
magasin_products_total

# Cart operations
rate(magasin_cart_operations_total[5m])
```

## Implementation Details

### Metrics Collection

Metrics are automatically collected through:

1. **HTTP Middleware**: All HTTP requests are automatically instrumented with request count and duration metrics
2. **Business Logic Instrumentation**: Key business operations (sales, refunds, cart operations) are instrumented with specific metrics
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
  - name: magasin-api
    rules:
      - alert: HighErrorRate
        expr: rate(magasin_errors_total[5m]) > 0.1
        for: 2m
        labels:
          severity: warning
        annotations:
          summary: "High error rate detected"
          
      - alert: SlowResponseTime
        expr: histogram_quantile(0.95, rate(magasin_http_request_duration_seconds_bucket[5m])) > 1
        for: 2m
        labels:
          severity: warning
        annotations:
          summary: "Slow response time detected"
``` 