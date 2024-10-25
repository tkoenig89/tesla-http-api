# tesla-http-api

REST API based on [tesla-http-proxy](https://github.com/teslamotors/vehicle-command/tree/main/cmd/tesla-http-proxy) and Tesla's [proxy](https://pkg.go.dev/github.com/teslamotors/vehicle-command/pkg/proxy) Go package. [tesla-http-api](https://github.com/L480/tesla-http-api/pkgs/container/tesla-http-api) sends end-to-end authenticated commands to your vehicle and takes care about the authorization flow and token refresh.

Good for controlling vehicles through incoming webhooks from tools like Grafana.

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
curl -H "Authorization: Bearer $API_TOKEN" \
     -H "Content-Type: application/json" \
     --data "{}" \
     -X POST \
     -i https://tesla-http-api.example.com/api/1/vehicles/{VIN}/command/flash_lights
```

You can find all API endpoints in [Tesla's Fleet API documentation](https://developer.tesla.com/docs/fleet-api/endpoints/vehicle-commands).

### Health check

```bash
curl -X GET \
     -i https://tesla-http-api.example.com/health
```