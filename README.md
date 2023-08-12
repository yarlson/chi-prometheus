# Prometheus Middleware for Chi Router

## Introduction

The `chiprom` library provides an elegant integration of Prometheus metrics into the chi router, offering detailed insights into HTTP request handling within a chi-based service. The metrics that can be exposed include request count, latency, and response size, categorized by status code, HTTP method, and path. This integration allows for effortless observability of web services, ensuring that operators can measure, analyze, and troubleshoot their applications' performance with ease.

## Usage

### Simple Usage

The most straightforward usage is to instrument your HTTP handlers with basic request metrics.

#### Example

```go
package main
import (
    "github.com/go-chi/chi/v5"
    "github.com/yarlson/chiprom"
)

func main() {
    r := chi.NewRouter()
    r.Use(chiprom.NewMiddleware("myservice"))

    // ... your handlers ...

    http.ListenAndServe(":8080", r)
}
```

This example demonstrates how you can start instrumenting your chi-based service with default Prometheus metrics. The middleware will track the request count, latency, and response size, categorized by the status code, HTTP method, and path.

### Using Custom Latency Buckets

You can customize the latency histograms with your desired buckets.

#### Example

```go
buckets := []float64{100, 500, 2000}
r.Use(chiprom.NewMiddleware("myservice", buckets...))
```

### Monitoring With Routing Patterns

To group requests by chi routing patterns, such as monitoring paths like `/users/{firstName}` instead of individual instances like `/users/bob`, you can utilize the `NewPatternMiddleware`.

#### Example

```go
r.Use(chiprom.NewPatternMiddleware("myservice"))
```

### Advanced Usage

You may want to combine both pattern monitoring and custom buckets. Here is a more complex example:

```go
buckets := []float64{100, 500, 2000}
r.Use(chiprom.NewPatternMiddleware("myservice", buckets...))

// ... your handlers ...

http.ListenAndServe(":8080", r)
```

## Metrics

The library exposes the following metrics:

- `chi_requests_total`: How many HTTP requests processed, partitioned by status code, method, and HTTP path.
- `chi_request_duration_milliseconds`: How long it took to process the request, partitioned by status code, method, and HTTP path.
- `chi_pattern_requests_total`: Similar to `chi_requests_total`, but with patterns.
- `chi_pattern_request_duration_milliseconds`: Similar to `chi_request_duration_milliseconds`, but with patterns.

## Installation

To add `chiprom` to your Go project, you can use `go get`:

```shell
go get -u github.com/yarlson/chiprom
```

## Contribute

Contributions are welcome. Feel free to open a pull request or file an issue on the [GitHub repository](https://github.com/yarlson/chiprom).

## License

This library is distributed under the Apache 2.0 License. See the [LICENSE](LICENSE) file for details.
