package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

type HistogramVec interface {
	prometheus.Observer
	WithLabelValues(...string) HistogramVec
}

func NoopHistogramVec() HistogramVec {
	return noopHistogramVec{}
}

type noopHistogramVec struct{}

func (hv noopHistogramVec) WithLabelValues(_ ...string) HistogramVec {
	return hv
}

func (noopHistogramVec) Observe(_ float64) {
}

type histogramVec struct {
	parent *prometheus.HistogramVec
	lvs    []string
}

func NewHistogramVecFrom(opts prometheus.HistogramOpts, lnames []string) HistogramVec {
	parent := prometheus.NewHistogramVec(opts, lnames)
	prometheus.MustRegister(parent)
	return &histogramVec{
		parent: parent,
	}
}

func (h *histogramVec) WithLabelValues(lvs ...string) HistogramVec {
	return &histogramVec{
		parent: h.parent,
		lvs:    append(h.lvs, lvs...),
	}
}

func (h *histogramVec) Observe(val float64) {
	labels := make(map[string]string, len(h.lvs)/2)
	for i := 0; i+1 < len(h.lvs); i += 2 {
		labels[h.lvs[i]] = h.lvs[i+1]
	}
	h.parent.With(labels).Observe(val)
}
