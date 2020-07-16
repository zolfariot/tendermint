package p2p

import (
	"github.com/go-kit/kit/metrics"
	"github.com/go-kit/kit/metrics/discard"
	"github.com/go-kit/kit/metrics/prometheus"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
)

const (
	// MetricsSubsystem is a subsystem shared by all metrics exposed by this
	// package.
	MetricsSubsystem = "p2p"
)

// Metrics contains metrics exposed by this package.
type Metrics struct {
	// Number of peers.
	Peers metrics.Gauge
	// Number of bytes received from a given peer.
	PeerReceiveBytesTotal metrics.Counter
	// Number of bytes sent to a given peer.
	PeerSendBytesTotal metrics.Counter
	// Pending bytes to be sent to a given peer.
	PeerPendingSendBytes metrics.Gauge
	// Number of transactions submitted by each peer.
	NumTxs metrics.Gauge

	ReactorReceiveDuration           metrics.Histogram
	ReactorConsensusReceiveDuration  metrics.Histogram
	ReactorBlockchainReceiveDuration metrics.Histogram
	ReactorMempoolReceiveDuration    metrics.Histogram
}

// PrometheusMetrics returns Metrics build using Prometheus client library.
// Optionally, labels can be provided along with their values ("foo",
// "fooValue").
func PrometheusMetrics(namespace string, labelsAndValues ...string) *Metrics {
	labels := []string{}
	for i := 0; i < len(labelsAndValues); i += 2 {
		labels = append(labels, labelsAndValues[i])
	}
	return &Metrics{
		Peers: prometheus.NewGaugeFrom(stdprometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: MetricsSubsystem,
			Name:      "peers",
			Help:      "Number of peers.",
		}, labels).With(labelsAndValues...),
		PeerReceiveBytesTotal: prometheus.NewCounterFrom(stdprometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: MetricsSubsystem,
			Name:      "peer_receive_bytes_total",
			Help:      "Number of bytes received from a given peer.",
		}, append(labels, "peer_id", "chID")).With(labelsAndValues...),
		PeerSendBytesTotal: prometheus.NewCounterFrom(stdprometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: MetricsSubsystem,
			Name:      "peer_send_bytes_total",
			Help:      "Number of bytes sent to a given peer.",
		}, append(labels, "peer_id", "chID")).With(labelsAndValues...),
		PeerPendingSendBytes: prometheus.NewGaugeFrom(stdprometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: MetricsSubsystem,
			Name:      "peer_pending_send_bytes",
			Help:      "Number of pending bytes to be sent to a given peer.",
		}, append(labels, "peer_id")).With(labelsAndValues...),
		NumTxs: prometheus.NewGaugeFrom(stdprometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: MetricsSubsystem,
			Name:      "num_txs",
			Help:      "Number of transactions submitted by each peer.",
		}, append(labels, "peer_id")).With(labelsAndValues...),
		ReactorReceiveDuration: prometheus.NewHistogramFrom(stdprometheus.HistogramOpts{
			Namespace: namespace,
			Subsystem: MetricsSubsystem,
			Name:      "reactor_receive_duration",
			Help:      "duration of execution of reactor.accept",
			Buckets:   stdprometheus.ExponentialBuckets(1, 5, 7),
		}, labels).With(labelsAndValues...),
		ReactorConsensusReceiveDuration: prometheus.NewHistogramFrom(stdprometheus.HistogramOpts{
			Namespace: namespace,
			Subsystem: MetricsSubsystem,
			Name:      "reactor_consensus_receive_duration",
			Help:      "duration of execution of reactor.accept",
			Buckets:   stdprometheus.ExponentialBuckets(1, 5, 7),
		}, labels).With(labelsAndValues...),
		ReactorBlockchainReceiveDuration: prometheus.NewHistogramFrom(stdprometheus.HistogramOpts{
			Namespace: namespace,
			Subsystem: MetricsSubsystem,
			Name:      "reactor_blockchain_receive_duration",
			Help:      "duration of execution of reactor.accept",
			Buckets:   stdprometheus.ExponentialBuckets(1, 5, 7),
		}, labels).With(labelsAndValues...),
		ReactorMempoolReceiveDuration: prometheus.NewHistogramFrom(stdprometheus.HistogramOpts{
			Namespace: namespace,
			Subsystem: MetricsSubsystem,
			Name:      "reactor_mempool_receive_duration",
			Help:      "duration of execution of reactor.accept",
			Buckets:   stdprometheus.ExponentialBuckets(1, 5, 7),
		}, labels).With(labelsAndValues...),
	}
}

// NopMetrics returns no-op Metrics.
func NopMetrics() *Metrics {
	return &Metrics{
		Peers:                            discard.NewGauge(),
		PeerReceiveBytesTotal:            discard.NewCounter(),
		PeerSendBytesTotal:               discard.NewCounter(),
		PeerPendingSendBytes:             discard.NewGauge(),
		NumTxs:                           discard.NewGauge(),
		ReactorReceiveDuration:           discard.NewHistogram(),
		ReactorConsensusReceiveDuration:  discard.NewHistogram(),
		ReactorBlockchainReceiveDuration: discard.NewHistogram(),
		ReactorMempoolReceiveDuration:    discard.NewHistogram(),
	}
}
