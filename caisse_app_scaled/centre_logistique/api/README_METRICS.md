# Prometheus Metrics for Centre Logistique API

This document describes the Prometheus metrics implementation for the Centre Logistique API service.

## Overview

The Centre Logistique API now includes comprehensive Prometheus metrics for monitoring and observability. Metrics are exposed at the `/api/v1/metrics` endpoint.

## Available Metrics

### HTTP Request Metrics

- **`logistique_http_requests_total`** (Counter)
  - Labels: `method`, `endpoint`, `status`
  - Description: Total number of HTTP requests

- **`logistique_http_request_duration_seconds`** (Histogram)
  - Labels: `method`, `endpoint`
  - Description: Duration of HTTP requests

### Business Metrics

- **`logistique_products_total`** (Gauge)
  - Description: Total number of products in inventory

- **`logistique_commands_total`** (Counter)
  - Description: Total number of commands created

- **`logistique_commands_accepted_total`** (Counter)
  - Description: Total number of accepted commands

- **`logistique_commands_refused_total`** (Counter)
  - Description: Total number of refused commands

- **`logistique_commands_pending`** (Gauge)
  - Description: Number of pending commands

### Error Metrics

- **`logistique_errors_total`** (Counter)
  - Labels: `type` (invalid_parameter, command_accept_failed, command_refuse_failed)
  - Description: Total number of errors by type

### Database Metrics

- **`logistique_database_operations_total`** (Counter)
  - Labels: `operation`
  - Description: Total number of database operations

- **`logistique_database_operation_duration_seconds`** (Histogram)
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
curl http://localhost:8091/api/v1/metrics
```

### Prometheus Configuration

Add the following to your `prometheus.yml` configuration:

```yaml
scrape_configs:
  - job_name: 'logistique-api'
    static_configs:
      - targets: ['localhost:8091']
    metrics_path: '/api/v1/metrics'
    scrape_interval: 15s
```

### Example Queries

Here are some useful PromQL queries for monitoring:

```promql
# Request rate by endpoint
rate(logistique_http_requests_total[5m])

# Error rate
rate(logistique_errors_total[5m])

# Average response time
histogram_quantile(0.95, rate(logistique_http_request_duration_seconds_bucket[5m]))

# Total commands
logistique_commands_total

# Command acceptance rate
rate(logistique_commands_accepted_total[5m])

# Pending commands
logistique_commands_pending

# Products in inventory
logistique_products_total
```

## Implementation Details

### Metrics Collection

Metrics are automatically collected through:

1. **HTTP Middleware**: All HTTP requests are automatically instrumented with request count and duration metrics
2. **Business Logic Instrumentation**: Key business operations (commands, acceptances, refusals) are instrumented with specific metrics
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
  - name: logistique-api
    rules:
      - alert: HighErrorRate
        expr: rate(logistique_errors_total[5m]) > 0.1
        for: 2m
        labels:
          severity: warning
        annotations:
          summary: "High error rate detected"
          
      - alert: SlowResponseTime
        expr: histogram_quantile(0.95, rate(logistique_http_request_duration_seconds_bucket[5m])) > 1
        for: 2m
        labels:
          severity: warning
        annotations:
          summary: "Slow response time detected"
          
      - alert: HighPendingCommands
        expr: logistique_commands_pending > 10
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: "High number of pending commands"
``` 