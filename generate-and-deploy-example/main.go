package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Prometheus metrics
var (
	requestCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total HTTP Requests",
		},
		[]string{"method", "endpoint", "http_status"},
	)

	requestLatency = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "HTTP Request Latency",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "endpoint"},
	)
)

var logger *log.Logger
var webLogger *ServiceLogger
var dashboardLogger *ServiceLogger

func init() {
	// Register Prometheus metrics
	prometheus.MustRegister(requestCount)
	prometheus.MustRegister(requestLatency)

	configureLogger()
	webLogger = NewServiceLogger(logger, "web-service")
	dashboardLogger = NewServiceLogger(logger, "dashboard-service")
}

func configureLogger() *os.File {
	// Use the mounted volume path
	logDir := "/tmp/logs"

	// Create logs directory with permissive permissions
	err := os.MkdirAll(logDir, 0777)
	if err != nil {
		log.Fatal("Failed to create log directory:", err)
	}

	// Open log file with permissive permissions
	logFile, err := os.OpenFile(filepath.Join(logDir, "app.log"), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("Failed to open log file:", err)
	}

	// Configure logger to write to both file and console
	logger = log.New(io.MultiWriter(os.Stdout, logFile), "", log.LstdFlags)

	return logFile
}

// Middleware for logging and metrics
func metricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Wrap the ResponseWriter to capture the status code
		rr := &responseRecorder{ResponseWriter: w, statusCode: http.StatusOK}
		next.ServeHTTP(rr, r)

		// Calculate latency
		latency := time.Since(start).Seconds()

		// Log request details with status and latency
		logMessage := fmt.Sprintf("%s %s %d %s [%.0fms]", r.Method, r.URL.Path, rr.statusCode, http.StatusText(rr.statusCode), latency*1000)
		if rr.statusCode >= 500 {
			webLogger.Error(logMessage)
		} else {
			webLogger.Info(logMessage)
		}

		// Update Prometheus metrics
		requestCount.WithLabelValues(r.Method, r.URL.Path, strconv.Itoa(rr.statusCode)).Inc()
		requestLatency.WithLabelValues(r.Method, r.URL.Path).Observe(latency)
	})
}

// Healthcheck handler
func healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	response := map[string]string{"status": "healthy"}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Dummy data handler
func dataHandler(w http.ResponseWriter, r *http.Request) {
	dummyData := map[string]interface{}{
		"message": "Hello, world!",
		"value":   42,
	}

	errorRate := 0.1 // 10% error rate
	if rand.Float64() < errorRate {
		// Simulate an error response
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Internal Server Error"})
		return
	}

	// Introduce an artificial delay between 0.25 - 2 seconds
	delay := time.Duration(250+rand.Intn(1750)) * time.Millisecond
	time.Sleep(delay)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(dummyData)
}

// Main function
func main() {
	// Define a flag to control whether to deploy the dashboard
	deployDashboard := flag.Bool("deploy-dashboard", false, "Deploy the dashboard to Grafana")
	flag.Parse()

	// Deploy the dashboard if the flag is set
	if *deployDashboard {
		GenerateAndPublishDashboard()
		return
	}

	mux := http.NewServeMux()

	// Define routes
	mux.HandleFunc("/health", healthcheckHandler)
	mux.HandleFunc("/data", dataHandler)
	mux.Handle("/metrics", promhttp.Handler())

	// Wrap the mux with the metrics middleware
	webLogger.Println("Starting server on :5001")
	if err := http.ListenAndServe(":5001", metricsMiddleware(mux)); err != nil {
		webLogger.Error(fmt.Sprintf("Error starting server: %v", err))
		os.Exit(1)
	}
}

// responseRecorder is a wrapper around http.ResponseWriter to capture the status code
type responseRecorder struct {
	http.ResponseWriter
	statusCode int
}

func (rr *responseRecorder) WriteHeader(code int) {
	rr.statusCode = code
	rr.ResponseWriter.WriteHeader(code)
}
