# OpenCensus Go Prometheus Exporter

[![Build Status](https://travis-ci.org/census-ecosystem/opencensus-go-exporter-prometheus.svg?branch=master)](https://travis-ci.org/census-ecosystem/opencensus-go-exporter-prometheus) [![GoDoc][godoc-image]][godoc-url]

Provides OpenCensus metrics export support for Prometheus.

## Installation

```
$ go get -u contrib.go.opencensus.io/exporter/prometheus
```

## Testing:
- Run the pushgateway:
```bash
docker run -d -p 9091:9091 prom/pushgateway
```

- Run tests using the local environment:
```bash
make test
```

- Run tests using the Docker environment:
```bash
make ci-test
```

[godoc-image]: https://godoc.org/contrib.go.opencensus.io/exporter/prometheus?status.svg
[godoc-url]: https://godoc.org/contrib.go.opencensus.io/exporter/prometheus
