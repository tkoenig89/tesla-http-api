# tesla-http-api

REST API based on [tesla-http-proxy](https://github.com/teslamotors/vehicle-command/tree/main/cmd/tesla-http-proxy) and Tesla's [proxy](https://pkg.go.dev/github.com/teslamotors/vehicle-command/pkg/proxy) Go package. [tesla-http-api](https://github.com/L480/tesla-http-api/pkgs/container/tesla-http-api) sends end-to-end authenticated commands to your vehicle and takes care about the authorization flow and token refresh.

![tesla-http-api request flow](/images/request-flow.png "tesla-http-api request flow")

## Usage

### Run tesla-http-api

```bash
# Step 1: Clone repo
git clone https://github.com/L480/tesla-http-api.git
# Step 2: Modify environment variables in docker-compose.yml
# Check out https://developer.tesla.com/docs/fleet-api/authentication/overview
nano docker-compose.yml
# Step 3: Create working directory
mkdir -p /opt/tesla-http-api
# Step 4: Copy your private key to working directory
cp my-private-key.pem /opt/tesla-http-api/private-key.pem
# Step 5: Change directory ownership as container image does not run as root
chown -R 93761:93761 /opt/tesla-http-api
# Step 6: Pull and start container
docker compose up -d
```

> [!IMPORTANT]  
> [tesla-http-api](https://github.com/L480/tesla-http-api/pkgs/container/tesla-http-api) is designed to run behind a proxy (like nginx, Cloudflare etc.).

### Use tesla-http-api

```bash
curl -H "X-Authorization: Bearer $API_TOKEN" \
     -H "Content-Type: application/json" \
     --data "{}" \
     -X POST \
     -i https://tesla-http-api.example.com/api/1/vehicles/{VIN}/command/flash_lights
```

You can find all API endpoints in [Tesla's Fleet API documentation](https://developer.tesla.com/docs/fleet-api/endpoints/vehicle-commands).

> Authentication header compatibility: Preferred header is `X-Authorization: Bearer <token>`. If only `Authorization: Bearer <token>` is sent, it is still accepted. If both headers are present they must match or the request is rejected (HTTP 400). Missing headers result in HTTP 401 and invalid tokens in HTTP 403.

### Health check

```bash
curl -X GET \
     -i https://tesla-http-api.example.com/health
```

### Development Container (VS Code / Dev Containers)

This repository includes a `.devcontainer` configuration so you can develop without installing Go locally.

Steps:
1. Install the "Dev Containers" extension in VS Code.
2. Open the folder (`tesla-http-api`).
3. When prompted, "Reopen in Container" (or press F1 and choose: Dev Containers: Reopen in Container).
4. The container will build using Go 1.23 and automatically run `go mod download`.
5. Run a quick build:
     ```bash
     go build ./cmd/tesla-http-api
     ```
6. Start the API (ensure env vars are set or use a dev `.env`):
     ```bash
     ENABLE_API_TOKEN=true API_TOKEN=devtoken TESLA_REFRESH_TOKEN=... TESLA_CLIENT_ID=... go run ./cmd/tesla-http-api
     ```

Port 8080 is auto-forwarded. Adjust `devcontainer.json` for additional tooling as needed.
