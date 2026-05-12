# gsongrok

**gsongrok** is a powerful, ultra-lightweight Go engine designed to serve dynamic JSON mocks and static assets with zero friction. It features an embedded **ngrok tunnel** and a premium **Admin Dashboard** for real-time API management and traffic inspection.

## ✨ Features
- **Standalone Engine**: A single Go binary that manages its own ngrok tunnel.
- **Dynamic JSON Mocking**: Map any path to a JSON response in real-time without restarts.
- **Premium Dashboard**: Manage your endpoints and inspect traffic from a sleek, dark-mode SPA.
- **Live Traffic Inspector**: Built-in replacement for ngrok's 4040 UI, including:
  - **Full Request Dumps**: View headers, raw request bodies, and "plain" request strings.
  - **JSON Tree Explorer**: Interactive exploration of JSON payloads (expand/collapse).
  - **Noise Filtering**: Automatically hides internal management and system requests (favicon, .well-known, etc.).
- **Static Asset Fallback**: Automatically serves files from the `public/` folder if no JSON mock matches the path.
- **Docker Ready**: Optimized multi-stage build (~15MB image).

## Deployment

### Using Docker (Standalone)
Run it directly with your ngrok authtoken:
```bash
docker run -d \
  -e APIKEY="your_ngrok_authtoken" \
  -e HOST="your_custom_domain.ngrok-free.app" \
  -v ${PWD}/data:/data \
  -p 8080:8080 \
  --name gsongrok \
  ramosisw/gsongrok
```

### Using Docker Compose (Recommended)
1. Create a `.env` file from the example:
   ```bash
   cp .env.example .env
   ```
2. Fill in your `NGROK_AUTHTOKEN` and `NGROK_DOMAIN` in the `.env` file.
3. Use the following `docker-compose.yml`:
```yaml
version: '3.8'
services:
  gsongrok:
    image: ramosisw/gsongrok
    container_name: gsongrok
    volumes:
      - ./data:/data
    environment:
      - APIKEY=${NGROK_AUTHTOKEN}
      - HOST=${NGROK_DOMAIN}
    ports:
      - "8080:8080"
    restart: always
```
4. Start the engine:
   ```bash
   docker-compose up -d
   ```

## 🛠️ Configuration
All routing is managed via `data/paths.json` or the Admin Dashboard at `http://localhost:8080/`.

### Example `paths.json`
```json
{
  "/api/health": { "status": "ok", "engine": "gsongrok" },
  "/api/user/1": { "id": 1, "name": "John Doe", "active": true }
}
```

## 📋 Management API
- `GET /gsongrok.json`: Retrieve the current JSON configuration.
- `POST /gsongrok.json`: Update the entire configuration (requires valid JSON).
- `GET /gsongrok/traffic`: Access the last 50 captured requests.
- `GET /gsongrok/info`: Check tunnel status and public URL.

## 💡 Pro Tips
> [!TIP]
> **Ngrok Warning Page**: On the free tier, ngrok shows a security warning the first time you access the tunnel in a browser. Just click **"Visit Site"** to reach the dashboard.
> 
> **Inspector**: Go to the "Inspector" tab in the dashboard to see live traffic details, including raw headers and parsed JSON bodies.

---
Built with ❤️ using **Go**, **Preact**, and **Ngrok-Go SDK**.