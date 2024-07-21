package prometheus

import (
	"awesomeProject/web"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"math/rand"
	"net/http"
	"testing"
	"time"
)

func Counter() {
	counter := prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: "my_namespace",
		Subsystem: "my_counter",
		Name:      "test_counter",
	})
	prometheus.MustRegister(counter)
	// +1
	counter.Inc()
	// 必须是整数
	counter.Add(10.2)
}

func Gauge() {
	gauge := prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: "my_namespace",
		Subsystem: "my_gauge",
		Name:      "test_gauge",
	})
	prometheus.MustRegister(gauge)
	gauge.Set(12)
	gauge.Add(10.2)
	gauge.Add(-3)
	gauge.Sub(3)
}

func Histogram() {
	histogram := prometheus.NewHistogram(prometheus.HistogramOpts{
		Namespace: "my_namespace",
		Subsystem: "my_histogram",
		Name:      "test_histogram",
		// 按照这个来分桶
		Buckets: []float64{10, 50, 100, 200, 500, 1000, 10000},
	})
	prometheus.MustRegister(histogram)
	histogram.Observe(12.4)
}

func Summary() {
	summary := prometheus.NewSummary(prometheus.SummaryOpts{
		Namespace: "my_namespace",
		Subsystem: "my_summary",
		Name:      "test_summary",
		// key 是百分比，value 是误差，比如说 0.5 - 0.01，
		// 那么它实际上计算的可能是 0.49 - 0.51 之间的
		Objectives: map[float64]float64{
			0.5:   0.01,
			0.75:  0.01,
			0.9:   0.01,
			0.99:  0.001,
			0.999: 0.0001,
		},
	})
	prometheus.MustRegister(summary)
	summary.Observe(12.3)
}

func Vector() {
	summary := prometheus.NewSummaryVec(prometheus.SummaryOpts{
		Namespace: "my_namespace",
		Subsystem: "my_summary",
		Name:      "test_summary",
		ConstLabels: map[string]string{
			"server":  "localhost",
			"env":     "test",
			"appname": "test_app",
		},
		Help: "The statics info for http request",
	}, []string{"pattern", "method", "status"})
	prometheus.MustRegister(summary)
	summary.WithLabelValues("/user/:id", "GET", "200").Observe(128)
}

func TestMiddleware_RecordsMetricsSuccessfully(t *testing.T) {
	builder := NewMiddlewareBuilder("test", "test", map[string]string{"appname": "test_app"}, "The statics info for http request")

	middleware := builder.Build()
	server := web.NewHTTPServer(
		web.ServerWithMiddleware(middleware),
	)

	server.Get("/test", func(ctx *web.Context) {
		// 睡眠随机的时间
		time.Sleep(time.Duration(100+rand.Intn(1000)) * time.Millisecond)
		ctx.RespJson(http.StatusOK, map[string]string{
			"message": "success",
		})
	})

	go func() {
		http.Handle("/metrics", promhttp.Handler())
		http.ListenAndServe(":8082", nil)
	}()

	server.Start(":8081")
}
