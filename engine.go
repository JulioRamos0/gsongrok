package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"sync"
	"time"
)

// Paths holds the mapping between URL paths and JSON objects
type Paths map[string]interface{}

// TrafficEvent represents a single HTTP request captured by the inspector
type TrafficEvent struct {
	ID         int64               `json:"id"`
	Timestamp  string              `json:"timestamp"`
	Method     string              `json:"method"`
	Path       string              `json:"path"`
	Host       string              `json:"host"`
	Status     int                 `json:"status"`
	Headers    map[string][]string `json:"headers"`
	Body       string              `json:"body"`
	RawRequest string              `json:"raw_request"`
}

var (
	paths      Paths
	pathsMu    sync.RWMutex
	traffic    []TrafficEvent
	trafficMu  sync.RWMutex
	maxTraffic = 50
)

// statusRecorder is a wrapper to capture the HTTP status code
type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (r *statusRecorder) WriteHeader(status int) {
	r.status = status
	r.ResponseWriter.WriteHeader(status)
}

// recordTraffic adds a new event with full details
func recordTraffic(r *http.Request, status int, body string) {
	// Ignore noisy paths
	path := r.URL.Path
	if path == "/" || 
	   path == "/favicon.ico" || 
	   (len(path) >= 10 && path[:10] == "/gsongrok/") || 
	   path == "/gsongrok.json" ||
	   (len(path) >= 13 && path[:13] == "/.well-known/") {
		return
	}

	trafficMu.Lock()
	defer trafficMu.Unlock()

	// Dump the request for the "Plain Request" view
	dump, _ := httputil.DumpRequest(r, true)

	event := TrafficEvent{
		ID:         time.Now().UnixNano(),
		Timestamp:  time.Now().Format("15:04:05"),
		Method:     r.Method,
		Path:       r.URL.Path,
		Host:       r.Host,
		Status:     status,
		Headers:    r.Header,
		Body:       body,
		RawRequest: string(dump),
	}

	traffic = append([]TrafficEvent{event}, traffic...)
	if len(traffic) > maxTraffic {
		traffic = traffic[:maxTraffic]
	}
}

// loadPaths reads the JSON configuration from data/paths.json
func loadPaths() error {
	filePath := os.Getenv("PATHS_JSON_PATH")
	if filePath == "" {
		filePath = "data/paths.json"
	}

	// Ensure directory exists
	dir := "data"
	if filePath != "data/paths.json" {
		// If custom path, we don't necessarily want to create the dir, 
		// but let's be safe if it's the default.
	} else {
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			os.MkdirAll(dir, 0755)
		}
	}

	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		log.Printf("Creating default paths file at %s", filePath)
		if err := os.WriteFile(filePath, defaultPaths, 0644); err != nil {
			return err
		}
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	var newPaths Paths
	if err := json.Unmarshal(data, &newPaths); err != nil {
		return err
	}

	pathsMu.Lock()
	paths = newPaths
	pathsMu.Unlock()
	return nil
}

// dynamicHandler resolves request paths using paths.json or falls back to public folder
func dynamicHandler(fileServer http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Read body and put it back for the handler
		var bodyString string
		if r.Body != nil {
			bodyBytes, _ := io.ReadAll(r.Body)
			bodyString = string(bodyBytes)
			r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		}

		recorder := &statusRecorder{ResponseWriter: w, status: 200}
		
		recorder.Header().Set("Access-Control-Allow-Origin", "*")
		recorder.Header().Set("ngrok-skip-browser-warning", "true")

		defer func() {
			recordTraffic(r, recorder.status, bodyString)
		}()

		if err := loadPaths(); err != nil {
			log.Printf("Error reloading paths: %v", err)
		}

		pathsMu.RLock()
		defer pathsMu.RUnlock()

		if obj, ok := paths[r.URL.Path]; ok {
			recorder.Header().Set("Content-Type", "application/json")
			json.NewEncoder(recorder).Encode(obj)
			return
		}

		fileServer.ServeHTTP(recorder, r)
	}
}

// managementHandler handles GET and POST to /gsongrok.json
func managementHandler(w http.ResponseWriter, r *http.Request) {
	filePath := os.Getenv("PATHS_JSON_PATH")
	if filePath == "" {
		filePath = "data/paths.json"
	}

	switch r.Method {
	case http.MethodGet:
		data, err := os.ReadFile(filePath)
		if err != nil {
			http.Error(w, "Failed to read config", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(data)

	case http.MethodPost:
		var newPaths Paths
		if err := json.NewDecoder(r.Body).Decode(&newPaths); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		data, err := json.MarshalIndent(newPaths, "", "  ")
		if err != nil {
			http.Error(w, "Failed to encode JSON", http.StatusInternalServerError)
			return
		}

		if err := os.WriteFile(filePath, data, 0644); err != nil {
			http.Error(w, "Failed to write config", http.StatusInternalServerError)
			return
		}

		loadPaths()
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message": "Config updated successfully"}`))

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// trafficHandler returns the captured traffic events
func trafficHandler(w http.ResponseWriter, r *http.Request) {
	trafficMu.RLock()
	defer trafficMu.RUnlock()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(traffic)
}
