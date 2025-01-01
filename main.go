package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"sync"
	"time"

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
	queueLength := prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "queue_length",
			Help: "The current length of the queue",
		},
	)

	// Register metrics
	prometheus.MustRegister(httpRequestsTotal)
	prometheus.MustRegister(currentSpeed)
	prometheus.MustRegister(queueLength)

	// Update queue length periodically
	go func() {
		for {
			updateQueue(queueLength)
			time.Sleep(1 * time.Second)
		}
	}()

	// Wrap the default mux with logging middleware
	http.Handle("/", logRequest(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Increment the counter
		httpRequestsTotal.WithLabelValues(r.Method, r.URL.Path).Inc()
		w.Write([]byte("Welcome to the Speed Controller! Use /faster or /slower to adjust speed."))
	})))

	http.Handle("/faster", logRequest(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		httpRequestsTotal.WithLabelValues(r.Method, r.URL.Path).Inc()
		speedMutex.Lock()
		speed += 1.0
		currentSpeed.Set(speed)
		speedMutex.Unlock()
		w.Write([]byte(fmt.Sprintf("Increased speed to %.2f", speed)))
	})))

	http.Handle("/slower", logRequest(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		httpRequestsTotal.WithLabelValues(r.Method, r.URL.Path).Inc()
		speedMutex.Lock()
		if speed > 0 {
			speed -= 1.0
		}
		currentSpeed.Set(speed)
		speedMutex.Unlock()
		w.Write([]byte(fmt.Sprintf("Decreased speed to %.2f", speed)))
	})))

	http.Handle("/metrics", logRequest(promhttp.Handler()))

	// Start HTTP server
	log.Println("Starting server on :8080")
	http.ListenAndServe(":8080", nil)
}

// Middleware to log requests
func logRequest(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Request: %s %s", r.Method, r.URL.Path)
		handler.ServeHTTP(w, r)
	})
}

func updateQueue(queueLength prometheus.Gauge) {
	speedMutex.Lock()
	currentSpeed := speed
	speedMutex.Unlock()

	randomComponent := rand.Float64() * 10 // Random value between 0 and 10
	queueValue := currentSpeed*2 + randomComponent
	queueLength.Set(queueValue)
}
