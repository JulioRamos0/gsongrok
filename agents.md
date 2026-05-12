## Gsongrok Engine Architecture

### Folder Structure
```
gsongrok/
├── data/             # Persistent data
│   └── paths.json    # JSON route mappings
├── main.go           # Entry point & Mux setup
├── engine.go         # JSON logic, Traffic Inspector, & Handlers
├── tunnel.go         # Ngrok SDK integration & Session management
├── public/          # Frontend & Assets
│   └── index.html   # Premium Dashboard (Preact SPA)
├── Dockerfile       # Multi-stage optimized build
└── docker-compose.yml # Orchestration
```

### Backend Implementation
- **Modular Design**: Separated concerns into `engine` (logic) and `tunnel` (connectivity).
- **Traffic Inspector**: Captures `TrafficEvent` objects in a circular buffer (last 50 requests).
- **Resilient Tunneling**: Uses an explicit `ngrok.Connect` session. If the tunnel fails (e.g. domain taken), the server continues to run locally.
- **Dynamic Reloading**: Reads `paths.json` from disk on every request to ensure zero-restart updates.
- **Static Fallback**: If a path isn't in JSON, it checks the `public/` directory.

### Frontend Dashboard
- **Technology**: Preact + HTM (no build step, pure ESM).
- **Inspector Tab**: Advanced view with Raw Request Dumps, Header lists, and an interactive **JsonView** component for body exploration.
- **Config Tab**: Raw JSON editor with real-time validation and Toast feedback (replacing alerts).
- **Responsive Status**: Monitors engine health and ngrok URL status via polling.

### Environment Variables
- `APIKEY`: (Required for tunnel) Ngrok Authtoken.
- `HOST`: (Optional) Custom ngrok domain.
- `PATHS_JSON_PATH`: (Optional) Defaults to `/data/paths.json`.
- `PORT`: (Optional) Internal port, defaults to `8080`.