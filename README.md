# WRZ API

A minimal Go API service with modular structure. It provides:
- Current time info
- Page info parsing (title, icon, description)
- OpenAPI spec and Swagger UI
- Navigation page + admin backend (embedded Vue frontend)

## Run

```powershell
go run ./cmd/server
```

Server defaults to port 8080. Set `PORT` to override.
Nav data defaults to `data.json`. Set `NAV_DATA` or `--nav-data` to override.

## Endpoints

- `GET /healthz`
- `GET /api/time`
- `GET /api/page-info?url=https://example.com`
- `GET /openapi.yaml`
- `GET /openapi.json`
- `GET /docs`
- `GET /` (导航页面)
- `GET /admin` (后台管理)
- `GET /login` (登录页)
- `GET /api/data`
- `POST /api/login`
- `POST /api/logout`
- `PUT /api/password`
- `POST /api/category`
- `PUT /api/category/{id}`
- `DELETE /api/category/{id}`
- `POST /api/item`
- `PUT /api/item/{id}`
- `DELETE /api/item/{id}`

## Responses

`GET /api/time`
```json
{
  "utc": "2026-01-26T15:40:12Z",
  "unix": 1769442012,
  "local": "2026-01-26T23:40:12+08:00"
}
```

`GET /api/page-info?url=https://example.com`
```json
{
  "url": "https://example.com",
  "title": "Example Domain",
  "description": "Example Domain",
  "icon": "https://example.com/favicon.ico"
}
```

## OpenAPI

- Spec: `http://localhost:8080/openapi.yaml`
- JSON: `http://localhost:8080/openapi.json`
- UI: `http://localhost:8080/docs`

## Navigation App (Nav)

默认账号密码：
```
admin / admin
```

前端构建（用于嵌入 Go 二进制）：
```bash
cd frontend
npm install
npm run build
```

开发模式（直接从磁盘读取 `frontend/dist`）：
```powershell
go run ./cmd/server --nav-dev
```

## Local build (one-click)

Windows:
```powershell
powershell -ExecutionPolicy Bypass -File scripts\\build_local.ps1
```

macOS/Linux:
```bash
bash scripts/build_local.sh
```

## CI/CD (GitHub Actions + systemd)

### Trigger

Push a tag like `v1.0.0` to GitHub. The workflow builds the Linux binary and deploys it to your server, then restarts the systemd service.

### GitHub Secrets

Add these secrets in your GitHub repo:
- `DEPLOY_HOST` (e.g. `1.2.3.4`)
- `DEPLOY_PORT` (e.g. `22`)
- `DEPLOY_USER` (ssh user, must be able to `sudo systemctl restart wrzapi`)
- `DEPLOY_SSH_KEY` (private key content)

### OpenAPI server URL

The Swagger UI "Servers" dropdown uses the current request host by default. You can override it by setting `SERVER_URL` (e.g. `https://api.example.com`).

Two options:
1) Start command override: `./wrzapi --server-url https://api.example.com`
2) systemd/env: set `SERVER_URL` in `wrzapi.service` (or environment)

### Nav data file

Two options:
1) Start command override: `./wrzapi --nav-data /path/to/data.json`
2) systemd/env: set `NAV_DATA` in `wrzapi.service` (or environment)

### Server setup (one-time)

Run this script from this repo on the server:

```bash
bash scripts/server_setup.sh
```

If your deploy user is not root, make sure it can run `sudo systemctl restart wrzapi` without a password (sudoers).

### Deploy flow

1) Tag and push:

```bash
git tag v1.1.0
git push origin v1.1.0
```

2) GitHub Actions builds and uploads to `/opt/wrzapi/releases`

3) Server switches `current` symlink and restarts systemd

### Rollback

On the server:

```bash
# rollback to previous release (based on timestamp order)
bash scripts/rollback.sh

# or rollback to a specific tag
bash scripts/rollback.sh v1.0.0
```
