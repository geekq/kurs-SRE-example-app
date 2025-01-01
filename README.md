# Shop Metrics

`shop-metrics` is a simple Go application to demonstrate how to expose Prometheus metrics. It allows you to adjust a global `speed` variable via `/faster` and `/slower` endpoints.

## Endpoints
- `/` - Welcome message
- `/faster` - Increases the `speed`
- `/slower` - Decreases the `speed`
- `/metrics` - Exposes Prometheus metrics

## Running the Application
```bash
go run main.go

