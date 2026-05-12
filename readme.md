# gsongrok v1.0.0

**gsongrok** is a powerful, ultra-lightweight Go engine designed to serve dynamic JSON mocks and static assets with zero friction. It features an embedded **ngrok tunnel** and a premium **Admin Dashboard** for real-time API management and traffic inspection.

## ✨ Features
- **Standalone Engine**: A single Go binary that manages its own ngrok tunnel.
- **Security First**: Separate local and public multiplexers. The public tunnel **only** exposes your mocks; the dashboard and logs are restricted to `localhost`.
- **Zero-Config CLI**: Automatically detects, creates, and loads `.env` files. Simply run the binary and it's ready.
- **Dynamic JSON Mocking**: Map any path to a JSON response in real-time via the dashboard.
- **Premium Dashboard**: Manage your endpoints and inspect traffic from a sleek, dark-mode SPA.
- **Live Traffic Inspector**: Circular buffer (last 50 requests) with raw dumps, headers, and interactive JSON tree viewer.
- **Noise Filtering**: Automatically ignores noisy system/browser requests (favicon, .well-known, etc.).
- **Asset Embedding**: All frontend assets are embedded in the single binary (Go 1.16+ embed).
- **Docker Ready**: Optimized multi-stage build (~15MB image).

## 🚀 Installation

### Using Go
If you have Go installed (1.26+), you can install the binary globally:
```bash
go install github.com/JulioRamos0/gsongrok@latest
```

### From Source
```bash
git clone https://github.com/JulioRamos0/gsongrok.git
cd gsongrok
go build -o gsongrok .
./gsongrok
```

## 🛠️ Usage

### Quick Start
Just run `./gsongrok`. On the first run, it will create a `.env` file for you. Fill in your `APIKEY` and restart to enable the tunnel.

### Private Dashboard
For security, the dashboard and management APIs are **only accessible locally**:
- **Dashboard**: [http://localhost:8080/](http://localhost:8080/)
- **Public API**: `https://your-domain.ngrok-free.app/api/path`

### Environment Variables
| Variable | Description | Default |
| :--- | :--- | :--- |
| `APIKEY` | Ngrok Authtoken | (Required for tunnel) |
| `HOST` | Custom Ngrok Domain | (Optional) |
| `PORT` | Local Server Port | `8080` |
| `PATHS_JSON_PATH` | Path to mocks file | `data/paths.json` |

## 🐳 Docker Deployment

### Docker Compose
```yaml
services:
  gsongrok:
    image: ramosisw/gsongrok
    container_name: gsongrok
    volumes:
      - ./data:/data
    environment:
      - APIKEY=${APIKEY}
      - HOST=${HOST}
    ports:
      - "8080:8080"
    restart: always
```

---
Built with ❤️ using **Go**, **Preact**, and **Ngrok-Go SDK**.