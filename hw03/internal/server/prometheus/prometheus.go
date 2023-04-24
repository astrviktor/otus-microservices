package prometheus

import (
	"github.com/prometheus/client_golang/prometheus"
)

type MetricsStruct struct {
	ResponseTime *prometheus.HistogramVec
}

func NewPrometheus() {
	responseTime := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "response_duration_ms",
			Help: "response time, ms",
		},
		[]string{"method", "path", "code"},
	)

	prometheus.MustRegister(responseTime)

	//responseTime.Set(0)

	Metrics.ResponseTime = responseTime
}

var Metrics MetricsStruct
