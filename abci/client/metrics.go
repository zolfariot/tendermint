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
	buckets := prometheus.ExponentialBuckets(1, 5, 10)
	return &Metrics{
		LockWaitDuration: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Namespace: namespace,
			Subsystem: metricsSubsystem,
			Name:      "lock_wait_duration",
			Help:      "time spent waiting for lock",
			Buckets:   buckets,
		}, labels),
		UnlockedDuration: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Namespace: namespace,
			Subsystem: metricsSubsystem,
			Name:      "unlocked_duration",
			Help:      "execution time sans lock wait",
			Buckets:   buckets,
		}, labels),
		TotalDuration: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Namespace: namespace,
			Subsystem: metricsSubsystem,
			Name:      "total_duration",
			Help:      "total execution time",
			Buckets:   buckets,
		}, labels),
	}
}

func NopMetrics() *Metrics {
	return &Metrics{
		LockWaitDuration: metrics.NoopHistogramVec(),
		UnlockedDuration: metrics.NoopHistogramVec(),
		TotalDuration:    metrics.NoopHistogramVec(),
	}
}
