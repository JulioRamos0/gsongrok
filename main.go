package main

import (
	"encoding/json"
	"io/fs"
	"log"
	"net/http"
	"os"
	"strings"
)

func init() {
	loadEnv()
}

func loadEnv() {
	if os.Getenv("APIKEY") != "" {
		return
	}

	envFile := ".env"
	if _, err := os.Stat(envFile); os.IsNotExist(err) {
		log.Println(".env not found, initializing...")
		if _, err := os.Stat(".env.example"); err == nil {
			data, _ := os.ReadFile(".env.example")
			os.WriteFile(envFile, data, 0644)
			log.Println("Successfully created .env from .env.example")
		} else {
			defaultEnv := []byte("APIKEY=\nHOST=\nPORT=8080\n")
			os.WriteFile(envFile, defaultEnv, 0644)
			log.Println("Successfully created a fresh .env file")
		}
	}

	content, err := os.ReadFile(envFile)
	if err != nil {
		return
	}

	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			val := strings.TrimSpace(parts[1])
			val = strings.Trim(val, `"'`)
			if os.Getenv(key) == "" {
				os.Setenv(key, val)
			}
		}
	}
}

func main() {
	if err := loadPaths(); err != nil {
		log.Printf("Warning: could not load initial paths: %v", err)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	localMux := http.NewServeMux()
	publicMux := http.NewServeMux()

	var fileServer http.Handler
	if _, err := os.Stat("public"); err == nil {
		fileServer = http.FileServer(http.Dir("public"))
	} else {
		subFS, _ := fs.Sub(publicFS, "public")
		fileServer = http.FileServer(http.FS(subFS))
	}

	localMux.HandleFunc("/", dynamicHandler(fileServer))
	publicMux.HandleFunc("/", dynamicHandler(http.HandlerFunc(http.NotFound)))

	localMux.HandleFunc("/gsongrok.json", managementHandler)
	localMux.HandleFunc("/gsongrok/traffic", trafficHandler)
	localMux.HandleFunc("/gsongrok/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})
	localMux.HandleFunc("/gsongrok/info", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("ngrok-skip-browser-warning", "true")
		infoMu.RLock()
		defer infoMu.RUnlock()
		json.NewEncoder(w).Encode(currentInfo)
	})

	startTunnel(publicMux)

	log.Printf("gsongrok engine starting on port %s...", port)
	log.Printf("Dashboard available at: http://localhost:%s", port)
	if err := http.ListenAndServe(":"+port, localMux); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
