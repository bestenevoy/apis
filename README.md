# WRZ API

A minimal Go API service with modular structure. It provides:
- Current time info
- Page info parsing (title, icon, description)
- OpenAPI spec and Swagger UI

## Run

```powershell
go run ./cmd/server
```

Server defaults to port 8080. Set `PORT` to override.

## Endpoints

- `GET /healthz`
- `GET /api/time`
- `GET /api/page-info?url=https://example.com`
- `GET /openapi.yaml`
- `GET /docs`

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
- UI: `http://localhost:8080/docs`

## CI/CD (GitHub Actions + systemd)

### Trigger

Push a tag like `v1.0.0` to GitHub. The workflow builds the Linux binary and deploys it to your server, then restarts the systemd service.

### GitHub Secrets

Add these secrets in your GitHub repo:
- `DEPLOY_HOST` (e.g. `1.2.3.4`)
- `DEPLOY_PORT` (e.g. `22`)
- `DEPLOY_USER` (ssh user, must be able to `sudo systemctl restart wrzapi`)
- `DEPLOY_SSH_KEY` (private key content)

### Server setup (one-time)

Run this script from this repo on the server:

```bash
bash scripts/server_setup.sh
```

If your deploy user is not root, make sure it can run `sudo systemctl restart wrzapi` without a password (sudoers).

### Deploy flow

1) Tag and push:

```bash
git tag v1.0.0
git push origin v1.0.0
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
