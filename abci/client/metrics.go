package abcicli

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/tendermint/tendermint/libs/metrics"
)

const (
	metricsSubsystem = "client"
)

type Metrics struct {
	LockWaitDuration metrics.HistogramVec
	UnlockedDuration metrics.HistogramVec
	TotalDuration    metrics.HistogramVec
}

func PrometheusMetrics(namespace string, labels ...string) *Metrics {
	return &Metrics{
		LockWaitDuration: prometheus.NewHistogramVec(prometheus.HistogramOpts{}, labels),
		UnlockedDuration: prometheus.NewHistogramVec(prometheus.HistogramOpts{}, labels),
		TotalDuration:    prometheus.NewHistogramVec(prometheus.HistogramOpts{}, labels),
	}
}

func NopMetrics() *Metrics {
	return &Metrics{
		LockWaitDuration: metrics.NoopHistogramVec(),
		UnlockedDuration: metrics.NoopHistogramVec(),
		TotalDuration:    metrics.NoopHistogramVec(),
	}
}
