package metrics

import (
	"strconv"
	"time"

	"runtime"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Metrics struct {
	RequestsTotal    *prometheus.CounterVec
	RequestDuration  *prometheus.HistogramVec
	GoroutinesGauge  prometheus.Gauge
	MemoryUsageGauge prometheus.Gauge
	ErrorCounter     *prometheus.CounterVec
}

var m *Metrics

func Init() {
	m = &Metrics{
		RequestsTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "application_http_requests_total",
				Help: "Total number of HTTP requests",
			},
			[]string{"method", "path", "status"},
		),

		RequestDuration: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "application_http_request_duration_seconds",
				Help:    "Duration of HTTP requests",
				Buckets: []float64{.005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5},
			},
			[]string{"method", "path"},
		),

		GoroutinesGauge: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Name: "application_goroutines_current",
				Help: "Current number of goroutines",
			},
		),

		MemoryUsageGauge: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Name: "application_memory_usage_bytes",
				Help: "Current memory usage",
			},
		),

		ErrorCounter: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "application_errors_total",
				Help: "Total number of application errors",
			},
			[]string{"type", "code"},
		),
	}

	// Register metrics
	prometheus.MustRegister(m.RequestsTotal)
	prometheus.MustRegister(m.RequestDuration)
	prometheus.MustRegister(m.GoroutinesGauge)
	prometheus.MustRegister(m.MemoryUsageGauge)
	prometheus.MustRegister(m.ErrorCounter)
}

// Handler returns the Prometheus metrics handler
func Handler() gin.HandlerFunc {
	h := promhttp.Handler()
	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

// Middleware returns a middleware that collects metrics
func Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if m == nil {
			// Initialize metrics if not already initialized
			Init()
		}

		start := time.Now()
		path := c.FullPath()
		if path == "" {
			path = "undefined"
		}

		// Process request
		c.Next()

		// Record metrics after request is processed
		duration := time.Since(start).Seconds()
		status := strconv.Itoa(c.Writer.Status())

		m.RequestsTotal.WithLabelValues(
			c.Request.Method,
			path,
			status,
		).Inc()

		m.RequestDuration.WithLabelValues(
			c.Request.Method,
			path,
		).Observe(duration)

		// Update system metrics
		m.GoroutinesGauge.Set(float64(runtime.NumGoroutine()))

		var mem runtime.MemStats
		runtime.ReadMemStats(&mem)
		m.MemoryUsageGauge.Set(float64(mem.Alloc))

		// Record errors if any
		if len(c.Errors) > 0 {
			m.ErrorCounter.WithLabelValues(
				"http",
				status,
			).Inc()
		}
	}
}

// RecordError records a custom error metric
func RecordError(errorType, code string) {
	if m == nil {
		Init()
	}
	m.ErrorCounter.WithLabelValues(errorType, code).Inc()
}
