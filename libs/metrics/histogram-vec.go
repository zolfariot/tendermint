package metrics

import "github.com/prometheus/client_golang/prometheus"

type HistogramVec interface {
	WithLabelValues(...string) prometheus.Observer
}

func NoopHistogramVec() HistogramVec {
	return noopHistogramVec{}
}

type noopHistogramVec struct{}

func (hv noopHistogramVec) WithLabelValues(_ ...string) prometheus.Observer {
	return hv
}

func (noopHistogramVec) Observe(_ float64) {
}
