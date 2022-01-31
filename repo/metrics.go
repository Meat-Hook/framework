package repo

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/Meat-Hook/framework/reflectx"
)

// MetricCollector is a helper for easy collecting metrics for every handler.
type MetricCollector interface {
	// Collecting collects Metrics information for handlers.
	Collecting(method string, f func() error) func() error
}

const (
	labelFunc = "func" // Value: caller's func/method name.
)

var _ MetricCollector = Metrics{}

// Metrics contains general metrics for DAL methods.
type Metrics struct {
	callErrTotal *prometheus.CounterVec
	callDuration *prometheus.HistogramVec
}

// NewMetrics registers and returns common DAL metrics used by all
// services (namespace).
func NewMetrics(reg *prometheus.Registry, namespace, subsystem string, methodsFrom interface{}) (metric Metrics) {
	metric.callErrTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      "errors_total",
			Help:      "Amount of DAL errors.",
		},
		[]string{labelFunc},
	)
	reg.MustRegister(metric.callErrTotal)
	metric.callDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      "call_duration_seconds",
			Help:      "DAL call latency.",
		},
		[]string{labelFunc},
	)
	reg.MustRegister(metric.callDuration)

	for _, methodName := range reflectx.MethodsOf(methodsFrom) {
		l := prometheus.Labels{
			labelFunc: methodName,
		}
		metric.callErrTotal.With(l)
		metric.callDuration.With(l)
	}

	return metric
}

// Collecting implements MetricCollector.
func (m Metrics) Collecting(method string, f func() error) func() error {
	return func() (err error) {
		start := time.Now()
		l := prometheus.Labels{labelFunc: method}
		defer func() {
			m.callDuration.With(l).Observe(time.Since(start).Seconds())
			if err != nil {
				m.callErrTotal.With(l).Inc()
			} else if err := recover(); err != nil {
				m.callErrTotal.With(l).Inc()
				panic(err)
			}
		}()
		return f()
	}
}

var _ MetricCollector = NoMetric{}

// NoMetric if you want to turn off metrics.
type NoMetric struct{}

// Collecting implements MetricCollector.
func (n NoMetric) Collecting(_ string, _ func() error) func() error {
	return func() error { return nil }
}
