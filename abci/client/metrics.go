package abcicli

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/tendermint/tendermint/libs/metrics"
)

const (
	metricsSubsystem = "client"

	metricsMethodKey = "method"
)

type Metrics struct {
	LockWaitDuration metrics.HistogramVec
	UnlockedDuration metrics.HistogramVec
	TotalDuration    metrics.HistogramVec
}

func PrometheusMetrics(namespace string, labels ...string) *Metrics {
	lnames := make([]string, 0, len(labels)/2+1)
	for i := 0; i+1 < len(labels); i += 2 {
		lnames = append(lnames, labels[i])
	}
	lnames = append(lnames, metricsMethodKey)

	return &Metrics{
		LockWaitDuration: metrics.NewHistogramVecFrom(prometheus.HistogramOpts{
			Namespace: namespace,
			Subsystem: metricsSubsystem,
			Name:      "lock_wait_duration",
			Help:      "time spent waiting for lock",
			Buckets:   prometheus.ExponentialBuckets(1, 5, 10),
		}, lnames).WithLabelValues(labels...),
		UnlockedDuration: metrics.NewHistogramVecFrom(prometheus.HistogramOpts{
			Namespace: namespace,
			Subsystem: metricsSubsystem,
			Name:      "unlocked_duration",
			Help:      "execution time sans lock wait",
			Buckets:   prometheus.ExponentialBuckets(1, 5, 10),
		}, lnames).WithLabelValues(labels...),
		TotalDuration: metrics.NewHistogramVecFrom(prometheus.HistogramOpts{
			Namespace: namespace,
			Subsystem: metricsSubsystem,
			Name:      "total_duration",
			Help:      "total execution time",
			Buckets:   prometheus.ExponentialBuckets(1, 5, 10),
		}, lnames).WithLabelValues(labels...),
	}
}

func NopMetrics() *Metrics {
	return &Metrics{
		LockWaitDuration: metrics.NoopHistogramVec(),
		UnlockedDuration: metrics.NoopHistogramVec(),
		TotalDuration:    metrics.NoopHistogramVec(),
	}
}
