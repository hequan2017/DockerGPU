[简体中文](README.md) | [English](README.en.md)

# DockerGPU

> A Docker-based GPU computing-power rental system

![Go](https://img.shields.io/badge/Go-1.23-00ADD8?logo=go&logoColor=white)
![Vue](https://img.shields.io/badge/Vue-3.x-4FC08D?logo=vuedotjs&logoColor=white)
![Docker](https://img.shields.io/badge/Docker%20API-26.x-2496ED?logo=docker&logoColor=white)
![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)
![gin-vue-admin](https://img.shields.io/badge/based%20on-gin--vue--admin%202.8.6-green)

## Introduction

DockerGPU is a front-end/back-end separated platform for renting and managing GPU computing power: administrators register GPU servers as **compute nodes** (connected remotely through the Docker API, with optional mutual-TLS certificates), then define shippable **product specs** by GPU model, CPU, memory and disk. When a user places an order, they pick an image and a spec; the system automatically matches a node with enough remaining capacity and creates a resource-limited Docker container on it, which becomes a directly operable "GPU instance".

The platform ships with full instance lifecycle management: container start, restart, stop, log viewing, command execution and a **Web interactive terminal** (xterm.js + WebSocket). Container status is synced automatically after every operation and list query. Regular users can only see and operate the instances they created, while administrators (authority 888) manage everyone's instances — a natural fit for internal GPU sharing, small-team GPU rental, and classroom lab environments.

The project is built on top of the [gin-vue-admin](https://github.com/flipped-aurora/gin-vue-admin) 2.8.6 scaffold and inherits its full admin toolbox: users/roles/API permissions (Casbin), menu management, operation records, code generator, form designer, plugin system, multi-cloud object storage, etc. The business modules (image registry, compute nodes, product specs, instance management) are custom-developed.

## ✨ Features

- **Image registry**: manage rentable container images (name, address, description, type) with an on/off-shelf switch
- **Compute node management**: register GPU servers (region, CPU, memory, system/data disk capacity, public/private IP, SSH info, GPU name and count), configure the Docker endpoint, and optional TLS (CA certificate, client certificate, private key) for encrypted connections
- **Product spec management**: define sellable specs by GPU model/count, CPU cores, memory and disk size, with hourly pricing and on/off-shelf control
- **One-click instance creation**: creates and starts a container on the selected node via the Docker API, injecting resource limits from the spec
  - CPU `NanoCPUs`, memory `Memory`, system-disk `StorageOpt` limits
  - NVIDIA GPU passthrough (`DeviceRequests` with the requested GPU count), plus `GPU_MODEL` / `GPU_COUNT` environment variables and spec labels
  - System/data disks mounted as named volumes at `/system` and `/data` inside the container
  - Container ID is written back automatically on success (status: running); failures are recorded as error
- **Smart node matching**: when creating an instance, nodes that are on-shelf are filtered against the spec requirements (GPU model/count, CPU, memory, disk) **after deducting capacity already consumed by existing instances**
- **Full container lifecycle operations**: restart, stop, view logs (`tail`-truncated, auto-refreshed every 2 seconds on the frontend), and execute commands inside the container
- **Web interactive terminal**: xterm.js connects over WebSocket to the backend, which pipes a binary stream through the Docker SDK's `ExecCreate + ExecAttach` (TTY) — operate the container as if it were SSH
- **Automatic status sync**: real container state is `inspect`-ed and written back to the database on list queries and after every operation (creating / running / exited / error)
- **Multi-tenant permission isolation**: regular users can only query and operate their own instances; administrators (authorityId 888) manage all; endpoints are guarded by Casbin and write operations are logged
- **Deletion protection**: an instance can only be deleted after its container (together with the mounted data volumes) has been removed; otherwise deletion is rejected
- **Full gin-vue-admin capabilities inherited**: RBAC roles, menu/API management, code generator, form designer, scheduled tasks, plugin system (email/announcement), multi-cloud uploads (Alibaba OSS, Tencent COS, AWS S3, Huawei OBS, MinIO, Cloudflare R2, local), AI-assisted MCP development tools, and more

## 🛠 Tech Stack

**Backend** (`server/`, Go 1.23)

- Web framework: [Gin](https://github.com/gin-gonic/gin) v1.10.0
- ORM: [GORM](https://github.com/go-gorm/gorm) v1.25.12 (MySQL / PostgreSQL / SQLite / Oracle / SQL Server)
- Auth: Casbin v2.103.0 + JWT (golang-jwt/v5)
- Container ops: [Docker SDK](https://github.com/moby/moby) (docker/docker v26.1.4) with TLS support for remote nodes
- Terminal transport: gorilla/websocket
- Others: Viper (config), Zap (logging), Swagger (swaggo docs), go-redis, mark3labs/mcp-go (MCP tools)

**Frontend** (`web/`)

- [Vue 3](https://github.com/vuejs/core) 3.5 + [Element Plus](https://github.com/element-plus/element-plus) 2.10
- Build: Vite 6 (dev proxy enables `ws: true` for WebSocket)
- State: Pinia 2.2; routing: Vue Router 4
- Terminal: xterm 5.3 (+ xterm-addon-fit)
- Others: ECharts 5.5, UnoCSS, wangEditor 5, Axios 1.8.2

**Deployment**

- Dockerfiles: `server/Dockerfile` (multi-stage golang:alpine), `web/Dockerfile` (node:20 + nginx)
- Orchestration: `deploy/docker-compose/docker-compose.yaml` (web + server + MySQL 8 + Redis), `deploy/kubernetes/` (K8s manifests)
- Build: root `Makefile` (in-container packaging, image builds, Swagger docs, plugin packaging)

## 🚀 Quick Start

### Requirements

- Go 1.23+, Node.js 20+, Docker
- MySQL 8 (other GORM-supported databases also work)
- Compute nodes need Docker installed with the Docker API exposed (TLS recommended)

### 1. Initialize the database (MySQL example)

```sql
-- Create database and user (example)
CREATE DATABASE IF NOT EXISTS gva DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;
CREATE USER IF NOT EXISTS 'gva'@'%' IDENTIFIED BY 'gva123456';
GRANT ALL PRIVILEGES ON gva.* TO 'gva'@'%';
FLUSH PRIVILEGES;
```

### 2. Start the backend (port 8888 by default)

```bash
cd server

# Edit config.yaml and set the database connection created above
# (mysql section: path/port/db-name/username/password)
# Business tables (registry, compute nodes, product specs, instances, etc.)
# are migrated automatically on first startup

go mod tidy
go run main.go
```

### 3. Start the frontend (port 8080 by default)

```bash
cd web
npm install --registry=https://registry.npmjs.org
npm run dev
```

Open `http://localhost:8080/`, finish the database initialization page, then log in with the default admin account `admin / 123456`.

### 4. Offer compute power

1. On the "Compute Nodes" page, register your GPU server: fill in the Docker endpoint (e.g. `tcp://1.2.3.4:2376`) and TLS certificates if enabled, and turn on the "on-shelf" switch
2. Publish images in the "Image Registry" and publish specs with pricing in the "Product Specs" page
3. When creating an instance, pick an image and a spec — the system matches a node with remaining capacity and creates the container automatically

### Frontend environment variables (`web/.env.development`)

| Variable | Description | Default |
| --- | --- | --- |
| `VITE_CLI_PORT` | Frontend dev port | `8080` |
| `VITE_SERVER_PORT` | Backend port | `8888` |
| `VITE_BASE_API` | API proxy prefix | `/api` |
| `VITE_BASE_PATH` | Proxy target | `http://127.0.0.1` |

### Docker deployment

```bash
# Build backend/frontend images (override the tag with the TAGS_OPT variable)
make build-image-server
make build-image-web

# Package artifacts locally and build the all-in-one image
make image

# One-command docker-compose deployment (web:8080 / server:8888 / mysql:13306 / redis:16379)
cd deploy/docker-compose
docker-compose up -d
```

Kubernetes manifests live in `deploy/kubernetes/` (ConfigMap, Deployment, Service, Ingress).

### Instance operation endpoints

| Method | Path | Description |
| --- | --- | --- |
| POST | `/inst/createInstance` | Create an instance (user ID taken from the login session) |
| DELETE | `/inst/deleteInstance` | Delete an instance (container must be removed first) |
| POST | `/inst/restartContainer` | Restart the container |
| POST | `/inst/stopContainer` | Stop the container |
| POST | `/inst/execContainerCmd` | Execute a command inside the container |
| GET | `/inst/getContainerLogs?ID=&tail=` | View container logs |
| GET | `/inst/getMatchedComputeNodes` | Match available nodes by spec |
| GET | `/inst/terminal?ID=&token=` | WebSocket interactive terminal (token passed as a query parameter) |

Permission records for these endpoints are created with the seed data; assign them per role in "Role Management → API Permissions". Auditing and whitelisting of terminal commands is recommended.

## 📁 Directory Structure

```
DockerGPU
├── server/                 # Backend (Go / Gin / GORM)
│   ├── api/v1/             # API layer (system/example + registry/compute/product/instance)
│   ├── model/              # Models (registry, compute, product, instance)
│   ├── service/            # Services (instance service with Docker API orchestration)
│   ├── router/             # Routers (business routes registered in initialize/router_biz.go)
│   ├── initialize/         # Init (router_biz.go route registration, gorm_biz.go auto-migration)
│   ├── mcp/                # AI-assisted MCP development tools
│   ├── config.yaml         # Backend configuration
│   └── Dockerfile
├── web/                    # Frontend (Vue3 / Element Plus / Vite)
│   ├── src/view/           # Pages (registry/compute/product/instance are business pages)
│   ├── vite.config.js      # Dev proxy (ws: true enabled)
│   └── Dockerfile
├── deploy/                 # docker / docker-compose / kubernetes deployment files
├── docs/                   # Screenshots
└── Makefile                # Build / images / docs / plugin packaging
```

## 📸 Screenshots

![System](./docs/gpu4.png)
![System](./docs/gpu3.png)
![System](./docs/gpu2.png)
![System](./docs/gpu.png)

## 🔗 Related Projects

- [gin-vue-admin](https://github.com/flipped-aurora/gin-vue-admin) — the admin scaffold this project is built on (Go + Vue3)

## 📄 License

This project is licensed under the [Apache License 2.0](./LICENSE).
