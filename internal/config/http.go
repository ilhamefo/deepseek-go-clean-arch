package config

import (
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	HTTPClientDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "garmin-exporter",
			Subsystem: "http_client",
			Name:      "request_duration_seconds",
			Help:      "Duration of outgoing HTTP requests",
			Buckets:   prometheus.DefBuckets,
		},
		[]string{"method", "code"},
	)

	HTTPClientRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "garmin-exporter",
			Subsystem: "http_client",
			Name:      "requests_total",
			Help:      "Total outgoing HTTP requests",
		},
		[]string{"method", "code"},
	)
)

func NewHTTPClient() *http.Client {
	return &http.Client{
		Timeout: 30 * time.Second,
		Transport: promhttp.InstrumentRoundTripperCounter(
			HTTPClientRequestsTotal,
			promhttp.InstrumentRoundTripperDuration(
				HTTPClientDuration,
				http.DefaultTransport,
			),
		),
	}
}
