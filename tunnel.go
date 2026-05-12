package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"sync"

	"golang.ngrok.com/ngrok"
	"golang.ngrok.com/ngrok/config"
)

// EngineInfo holds the current state of the gsongrok tunnel
type EngineInfo struct {
	NgrokURL string `json:"ngrok_url"`
	Error    string `json:"error,omitempty"`
	Status   string `json:"status"` // "offline", "online", "error", "starting"
}

var (
	currentInfo = EngineInfo{Status: "offline"}
	infoMu      sync.RWMutex
)

// startTunnel establishes the ngrok connection using an explicit session
func startTunnel(mux *http.ServeMux) {
	apiKey := os.Getenv("APIKEY")
	host := os.Getenv("HOST")

	if apiKey == "" {
		log.Printf("No APIKEY provided. Running in local-only mode.")
		return
	}

	updateStatus("starting", "", "")

	go func() {
		ctx := context.Background()
		
		// 1. Establish the session
		sess, err := ngrok.Connect(ctx, ngrok.WithAuthtoken(apiKey))
		if err != nil {
			log.Printf("ERROR: Failed to establish ngrok session: %v", err)
			updateStatus("error", err.Error(), "")
			return
		}

		// 2. Prepare endpoint options
		opts := []config.HTTPEndpointOption{}
		if host != "" {
			opts = append(opts, config.WithDomain(host))
		}

		// 3. Start listening on the tunnel
		tun, err := sess.Listen(ctx, config.HTTPEndpoint(opts...))
		if err != nil {
			log.Printf("ERROR: Failed to open ngrok tunnel: %v", err)
			updateStatus("error", err.Error(), "")
			return
		}

		updateStatus("online", "", tun.URL())
		log.Printf("SUCCESS: Tunnel is LIVE at %s", tun.URL())

		// 4. Serve requests
		if err := http.Serve(tun, mux); err != nil {
			log.Printf("Tunnel server stopped: %v", err)
			updateStatus("error", err.Error(), "")
		}
	}()
}

func updateStatus(status, errStr, url string) {
	infoMu.Lock()
	defer infoMu.Unlock()
	currentInfo.Status = status
	currentInfo.Error = errStr
	currentInfo.NgrokURL = url
}
