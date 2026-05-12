## Gsongrok Engine Architecture v1.0.0

### Folder Structure
```
gsongrok/
├── data/             # Persistent data (mocks)
├── public/           # Frontend assets (Dashboard)
├── main.go           # Entry point & Dual-Mux setup (Security Layer)
├── engine.go         # JSON logic, Traffic Inspector, & Management Handlers
├── tunnel.go         # Ngrok SDK integration & Session management
├── embedded.go       # Asset embedding logic (Go 1.16+ embed)
├── .env              # Local credentials (auto-generated)
├── Dockerfile        # Multi-stage optimized build (~15MB)
└── docker-compose.yml # Orchestration with env injection
```

### v1.0.0 Key Features
- **Security Isolation**: Uses separate `localMux` (Full access) and `publicMux` (Mocks only). The ngrok tunnel never exposes the dashboard or traffic logs.
- **Zero-Config CLI**: Automatically detects, creates, and loads `.env` files. Compatible with `go install`.
- **Hybrid FileSystem**: Serves dashboard files from the local `public/` directory if present, otherwise falls back to the embedded version inside the binary.
- **Embedded Tunneling**: Native Ngrok SDK integration (no external ngrok binary needed).
- **Traffic Inspector**: Circular buffer (50 events) with raw dump, header inspection, and interactive JSON tree viewer.
- **Dynamic Reloading**: Mocks are updated in real-time via the dashboard and reloaded on each request.
- **Noise Filtering**: Automatically ignores system/browser requests (favicon, .well-known, etc.) in the inspector.

### Environment Variables
- `APIKEY`: Required for ngrok tunneling.
- `HOST`: Optional custom ngrok domain.
- `PORT`: Local server port (defaults to `8080`).
- `PATHS_JSON_PATH`: Custom path for the mocks file.