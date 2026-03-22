package logic

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/zeromicro/go-zero/core/metric"
)

var (
	requestTotal = metric.NewCounterVec(&metric.CounterVecOpts{
		Name: "sentinel_request_total",
		Help: "Total number of requests",
		Labels: []string{"app_id", "resource"},
	})

	rejectTotal = metric.NewCounterVec(&metric.CounterVecOpts{
		Name: "sentinel_reject_total",
		Help: "Total number of rejected requests",
		Labels: []string{"app_id", "resource"},
	})

	l1HitTotal = metric.NewCounterVec(&metric.CounterVecOpts{
		Name: "sentinel_l1_hit_total",
		Help: "Total number of L1 cache hits",
		Labels: []string{"app_id", "resource"},
	})

	checkLatency = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "sentinel_check_latency_seconds",
			Help:    "Latency distribution of quota check",
			Buckets: []float64{.001, .002, .005, .01, .02, .05},
		},
		[]string{"app_id", "resource"},
	)
)

func init() {
	prometheus.MustRegister(checkLatency)
}
