package main

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Global variable to represent speed
var (
	speed      float64 = 1.0
	speedMutex         = &sync.Mutex{} // Mutex to ensure thread-safe access to the speed variable
)

func main() {
	// Define custom metrics
	httpRequestsTotal := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "endpoint"},
	)
	currentSpeed := prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "current_speed",
			Help: "The current speed value",
		},
	)

	// Register metrics
	prometheus.MustRegister(httpRequestsTotal)
	prometheus.MustRegister(currentSpeed)

	// Sample handler
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Increment the counter
		httpRequestsTotal.WithLabelValues(r.Method, r.URL.Path).Inc()
		w.Write([]byte("Welcome to the Speed Controller! Use /faster or /slower to adjust speed."))
	})

	// Faster handler
	http.HandleFunc("/faster", func(w http.ResponseWriter, r *http.Request) {
		httpRequestsTotal.WithLabelValues(r.Method, r.URL.Path).Inc()
		speedMutex.Lock()
		speed += 1.0
		currentSpeed.Set(speed)
		speedMutex.Unlock()
		w.Write([]byte(fmt.Sprintf("Increased speed to %.2f", speed)))
	})

	// Slower handler
	http.HandleFunc("/slower", func(w http.ResponseWriter, r *http.Request) {
		httpRequestsTotal.WithLabelValues(r.Method, r.URL.Path).Inc()
		speedMutex.Lock()
		if speed > 0 {
			speed -= 1.0
		}
		currentSpeed.Set(speed)
		speedMutex.Unlock()
		w.Write([]byte(fmt.Sprintf("Decreased speed to %.2f", speed)))
	})

	// Expose metrics endpoint
	http.Handle("/metrics", promhttp.Handler())

	// Start HTTP server
	http.ListenAndServe(":8080", nil)
}
