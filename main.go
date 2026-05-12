package main

import (
	"encoding/json"
	"io/fs"
	"log"
	"net/http"
	"os"
)

func main() {
	// Initial load of paths
	if err := loadPaths(); err != nil {
		log.Printf("Warning: could not load initial paths: %v", err)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Mux setup
	mux := http.NewServeMux()

	// Hybrid FileServer: Check local "public" first, then fallback to embedded
	var fileServer http.Handler
	if _, err := os.Stat("public"); err == nil {
		fileServer = http.FileServer(http.Dir("public"))
	} else {
		subFS, _ := fs.Sub(publicFS, "public")
		fileServer = http.FileServer(http.FS(subFS))
	}

	// Dynamic handler (JSON mocks + Static fallback)
	mux.HandleFunc("/", dynamicHandler(fileServer))

	// Management API
	mux.HandleFunc("/gsongrok.json", managementHandler)

	// Traffic Inspector API
	mux.HandleFunc("/gsongrok/traffic", trafficHandler)

	// Health check endpoint
	mux.HandleFunc("/gsongrok/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Engine Info endpoint
	mux.HandleFunc("/gsongrok/info", func(w http.ResponseWriter, r *http.Request) {
		infoMu.RLock()
		defer infoMu.RUnlock()
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(currentInfo)
	})

	// Background Tunnel startup
	startTunnel(mux)

	// Local HTTP server
	log.Printf("gsongrok engine starting on port %s...", port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
