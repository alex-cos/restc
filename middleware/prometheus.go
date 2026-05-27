package middleware

import (
	"strconv"
	"time"

	"github.com/alex-cos/restc"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

const Method = "method"
const StatusCode = "status_code"

type PrometheusMetrics struct {
	RequestsTotal    *prometheus.CounterVec
	RequestDuration  *prometheus.HistogramVec
	RequestsInFlight prometheus.Gauge
}

func NewPrometheusMetrics(prefix string) *PrometheusMetrics {
	p := prefix
	if p == "" {
		p = "restc"
	}
	return &PrometheusMetrics{
		RequestsTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: p + "_requests_total",
				Help: "Total number of HTTP requests made.",
			},
			[]string{Method, StatusCode},
		),
		RequestDuration: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    p + "_request_duration_seconds",
				Help:    "Duration of HTTP requests in seconds.",
				Buckets: prometheus.DefBuckets,
			},
			[]string{Method, StatusCode},
		),
		RequestsInFlight: promauto.NewGauge(
			prometheus.GaugeOpts{
				Name: p + "_requests_in_flight",
				Help: "Number of HTTP requests currently in flight.",
			},
		),
	}
}

func NewPrometheusMetricsWith(prefix string, reg prometheus.Registerer) *PrometheusMetrics {
	factory := promauto.With(reg)
	p := prefix
	if p == "" {
		p = "restc"
	}
	return &PrometheusMetrics{
		RequestsTotal: factory.NewCounterVec(
			prometheus.CounterOpts{
				Name: p + "_requests_total",
				Help: "Total number of HTTP requests made.",
			},
			[]string{Method, StatusCode},
		),
		RequestDuration: factory.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    p + "_request_duration_seconds",
				Help:    "Duration of HTTP requests in seconds.",
				Buckets: prometheus.DefBuckets,
			},
			[]string{Method, StatusCode},
		),
		RequestsInFlight: factory.NewGauge(
			prometheus.GaugeOpts{
				Name: p + "_requests_in_flight",
				Help: "Number of HTTP requests currently in flight.",
			},
		),
	}
}

func (pm *PrometheusMetrics) Handler() restc.Middleware {
	return func(req *restc.Request, next restc.HandlerFunc) (*restc.Response, error) {
		pm.RequestsInFlight.Inc()
		defer pm.RequestsInFlight.Dec()

		start := time.Now()
		resp, err := next(req)

		statusCode := "0"
		if resp != nil {
			statusCode = strconv.Itoa(resp.StatusCode())
		}

		pm.RequestDuration.WithLabelValues(req.Method(), statusCode).Observe(time.Since(start).Seconds())
		pm.RequestsTotal.WithLabelValues(req.Method(), statusCode).Inc()

		return resp, err
	}
}
