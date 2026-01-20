package config

import (
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	httpClientDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Subsystem: "http_client",
			Name:      "request_duration_seconds",
			Help:      "Duration of outgoing HTTP requests",
			Buckets:   prometheus.DefBuckets,
		},
		[]string{"method", "code"},
	)
)

func init() {
	prometheus.MustRegister(httpClientDuration)
}

func NewHTTPClient() *http.Client {
	return &http.Client{
		Transport: promhttp.InstrumentRoundTripperDuration(
			httpClientDuration,
			http.DefaultTransport,
		),
		Timeout: 30 * time.Second,
	}
}
