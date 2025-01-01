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
	shopCurrentSpeed := prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "shop_current_speed",
			Help: "The current speed value",
		},
	)
	shopQueueLength := prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "shop_queue_length",
			Help: "The current length of the queue",
		},
	)
	shopRequestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "shop_request_duration_seconds",
			Help:    "Histogram of request durations in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "endpoint"},
	)

	// Register metrics
	prometheus.MustRegister(httpRequestsTotal)
	prometheus.MustRegister(shopCurrentSpeed)
	prometheus.MustRegister(shopQueueLength)
	prometheus.MustRegister(shopRequestDuration)

	// Update queue length periodically
	go func() {
		for {
			updateQueue(shopQueueLength)
			time.Sleep(1 * time.Second)
		}
	}()

	// Wrap the default mux with logging middleware
	http.Handle("/", logAndMeasureRequest(shopRequestDuration, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Increment the counter
		httpRequestsTotal.WithLabelValues(r.Method, r.URL.Path).Inc()
		w.Write([]byte("Welcome to the Speed Controller! Use /faster or /slower to adjust speed."))
	})))

	http.Handle("/faster", logAndMeasureRequest(shopRequestDuration, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		httpRequestsTotal.WithLabelValues(r.Method, r.URL.Path).Inc()
		speedMutex.Lock()
		speed += 1.0
		shopCurrentSpeed.Set(speed)
		speedMutex.Unlock()
		w.Write([]byte(fmt.Sprintf("Increased speed to %.2f", speed)))
	})))

	http.Handle("/slower", logAndMeasureRequest(shopRequestDuration, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		httpRequestsTotal.WithLabelValues(r.Method, r.URL.Path).Inc()
		speedMutex.Lock()
		if speed > 0 {
			speed *= 0.7
		}
		shopCurrentSpeed.Set(speed)
		speedMutex.Unlock()
		w.Write([]byte(fmt.Sprintf("Decreased speed to %.2f", speed)))
	})))

	http.Handle("/metrics", logAndMeasureRequest(shopRequestDuration, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		randomDelay := time.Duration(rand.Intn(1000)) * time.Millisecond // Random delay between 0 and 300ms
		time.Sleep(randomDelay)
		promhttp.Handler().ServeHTTP(w, r)
	})))

	// Start HTTP server
	log.Println("Starting server on :8080")
	http.ListenAndServe(":8080", nil)
}

// Middleware to log and measure requests
func logAndMeasureRequest(histogram *prometheus.HistogramVec, handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		log.Printf("Request: %s %s", r.Method, r.URL.Path)
		handler.ServeHTTP(w, r)
		duration := time.Since(start).Seconds()
		histogram.WithLabelValues(r.Method, r.URL.Path).Observe(duration)
	})
}

func updateQueue(queueLength prometheus.Gauge) {
	speedMutex.Lock()
	currentSpeed := speed
	speedMutex.Unlock()

	mean := 5.0
	stdDev := 2.0
	randomComponent := rand.NormFloat64()*stdDev + mean // Random value with normal distribution
	queueValue := 100/(currentSpeed+0.1) + randomComponent
	queueLength.Set(queueValue)
}
