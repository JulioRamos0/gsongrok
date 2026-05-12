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

type EngineInfo struct {
	NgrokURL string `json:"ngrok_url"`
	Error    string `json:"error,omitempty"`
	Status   string `json:"status"`
}

var (
	currentInfo = EngineInfo{Status: "offline"}
	infoMu      sync.RWMutex
)

func startTunnel(mux *http.ServeMux) {
	apiKey := os.Getenv("APIKEY")
	host := os.Getenv("HOST")

	if apiKey == "" || apiKey == "your_authtoken_here" {
		log.Printf("No APIKEY provided. Running in local-only mode.")
		return
	}

	updateStatus("starting", "", "")

	go func() {
		ctx := context.Background()

		sess, err := ngrok.Connect(ctx, ngrok.WithAuthtoken(apiKey))
		if err != nil {
			log.Printf("ERROR: Failed to establish ngrok session: %v", err)
			updateStatus("error", err.Error(), "")
			return
		}

		opts := []config.HTTPEndpointOption{}
		if host != "" && host != "your_domain.ngrok-free.app" {
			opts = append(opts, config.WithDomain(host))
		}

		tun, err := sess.Listen(ctx, config.HTTPEndpoint(opts...))
		if err != nil {
			log.Printf("ERROR: Failed to open ngrok tunnel: %v", err)
			updateStatus("error", err.Error(), "")
			return
		}

		updateStatus("online", "", tun.URL())
		log.Printf("SUCCESS: Tunnel is LIVE at %s", tun.URL())

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
