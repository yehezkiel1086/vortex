# Vortex

> Open-source threat intelligence and security investigation platform built with Go and Next.js.

Vortex is a lightweight, distributed threat intelligence platform designed to ingest raw security telemetry, extract Indicators of Compromise (IOCs), enrich them with external threat intelligence, correlate related activity, and calculate explainable risk scores for security analysts.

---

## 🏗️ Architecture

Vortex is engineered using **Hexagonal Architecture (Ports and Adapters)**, ensuring core domain logic remains decoupled from databases, caches, message brokers, and transport layers.

```text
               INBOUND (Driving)                      CORE                      OUTBOUND (Driven)
         ┌───────────────────────────┐      ┌──────────────────────┐      ┌───────────────────────────┐
         │                           │      │                      │      │                           │
         │  • Gin REST API Handlers  │ ---> │  [Primary Ports]     │      │  • PostgreSQL Adapter     │
HTTP ──> │    (internal/adapter/     │      │   - IngestionService │      │    (sqlc + pgx/v5)        │
         │     handler/http)         │      │   - QueryService     │      │                           │
         │                           │      │                      │      │  • Redis Adapter          │
         ├───────────────────────────┤      │  [Domain & Services] │ ---> │    (24h Caching)          │
         │                           │      │   - Detection Rules  │      │                           │
RabbitMQ │  • Queue Consumer Handler │ ---> │   - Risk Scoring     │      │  • RabbitMQ Adapter       │
   ───>  │    (cmd/worker)           │      │   - IOC Extractor    │      │    (Persistent Pub/Sub)   │
         │                           │      │   - Correlation      │      │                           │
         │                           │      │                      │      │  • External TI Adapters   │
         └───────────────────────────┘      │  [Secondary Ports]   │      │    (GeoIP, VirusTotal)    │
                                            │   - Repositories     │      │                           │
                                            │   - Cache / Broker   │      └───────────────────────────┘
                                            └──────────────────────┘
```

---

## 🔄 Core Pipeline

```text
Security Telemetry (Honeypot, WAF, IDS, Firewall)
       │
       ▼
1. Ingestion & Normalization ──► PostgreSQL + RabbitMQ
       │
       ▼
2. Worker Pipeline
   ├── IOC Extraction (IP, Domain, URL, SHA256, MD5)
   ├── Detection Engine (SSH Brute Force, Port Scan, SQLi, XSS, Malware)
   │     └─► MITRE ATT&CK Mapping (T1110, T1046, T1190, T1059, T1083, T1105)
   ├── External Enrichment (ip-api.com GeoIP + VirusTotal v3)
   ├── Redis 24h Intelligence Caching
   ├── Graph Correlation (IP ↔ Domain ↔ File Hash)
   └── 5-Factor Risk & Confidence Scoring
       │
       ▼
3. Autonomous Alerting (Threshold ≥ 70)
       │
       ▼
4. Analyst Web Dashboard (Next.js 16 + shadcn/ui)
```

---

## 📊 Explainable Risk Scoring

Vortex evaluates threat risk using a transparent 5-factor calculation:

$$\text{Risk Score} = \text{Reputation}(0\text{--}30) + \text{Severity}(0\text{--}25) + \text{Frequency}(0\text{--}20) + \text{Confidence}(0\text{--}15) + \text{Correlation}(0\text{--}10)$$

| Score Range | Risk Level | Action |
|:---:|:---:|:---|
| **85 – 100** | 🔴 Critical | Immediate triage & alert dispatch |
| **70 – 84** | 🟠 High | High-priority security alert |
| **50 – 69** | 🟡 Medium | Logged & monitored |
| **30 – 49** | 🔵 Low | Informational record |
| **0 – 29** | ⚪ Info | Baseline telemetry |

---

## 🧱 Tech Stack

- **Backend**: Go (1.25+), Gin, `pgx/v5`, `sqlc`, `amqp091-go`, `go-redis/v9`
- **Frontend**: Next.js 16 (App Router), React 19, TypeScript, Tailwind CSS v4, shadcn/ui, Lucide Icons
- **Database & Storage**: PostgreSQL 17, Redis 8
- **Message Broker**: RabbitMQ 3.13
- **Intelligence Providers**: Free GeoIP (`ip-api.com`), VirusTotal v3 API, MITRE ATT&CK Framework

---

## 📂 Project Structure

```text
vortex/
├── cmd/
│   ├── api/main.go              # REST API server entrypoint (:8080)
│   ├── worker/main.go           # Pipeline background consumer daemon
│   └── collector/main.go        # Telemetry simulator CLI
├── db/
│   ├── schema/vortex_schema.sql # PostgreSQL DDL schema & indexes
│   └── queries/                 # SQL queries for sqlc code generation
├── internal/
│   ├── core/                    # Pure Domain & Hexagonal Ports
│   │   ├── domain/              # Entities (Event, Indicator, Alert, RiskScore)
│   │   ├── port/                # Driving & Driven interfaces
│   │   ├── service/             # Extractor, Detection, Enrichment, Correlation, Risk, Alert
│   │   └── util/                # Network validators & IP/hash regex
│   └── adapter/                 # Infrastructure Adapters
│       ├── config/              # Environment config container
│       ├── handler/http/        # Gin REST API controllers & router
│       ├── storage/postgres/    # sqlc generated code & repository adapters
│       ├── storage/redis/       # Redis cache repository adapter
│       ├── storage/rabbitmq/    # RabbitMQ publisher & consumer
│       └── enrichment/          # GeoIP and VirusTotal API clients
└── dashboard_view/              # Next.js 16 SOC web dashboard
    ├── app/                     # App router pages (Overview, IOCs, Investigation, Alerts)
    ├── components/              # shadcn/ui & cybersecurity widgets
    └── lib/                     # Typed API client & interfaces
```

---

## 🚀 Quick Start

### 1. Prerequisites

- Go (1.23+)
- Node.js (18+) & npm
- Docker & Docker Compose

### 2. Configure Environment

```bash
cp .env.example .env
```

### 3. Start Infrastructure

```bash
docker compose up -d
```

### 4. Initialize Database Schema

```bash
docker exec -i postgres psql -U vortex -d vortex < db/schema/vortex_schema.sql
```

### 5. Run the Application

In separate terminal windows:

```bash
# Terminal 1: Background Worker Daemon
go run ./cmd/worker

# Terminal 2: REST API Server
go run ./cmd/api

# Terminal 3: Next.js Frontend Dashboard
cd dashboard_view
npm install
npm run dev
```

Open **`http://localhost:3000`** to access the dashboard.

### 6. Send Test Telemetry

Dispatch realistic attack scenarios (SSH brute force, SQL injection, malware downloads, port scans):

```bash
go run ./cmd/collector -scenario=all
```

Or trigger events directly using the **"Ingest Telemetry"** button in the dashboard UI.

---

## 📡 REST API Endpoints

Base URL: `/api/v1`

| Method | Endpoint | Description |
|:---|:---|:---|
| `POST` | `/events` | Ingest normalized security telemetry |
| `GET` | `/events` | List ingested events (paginated) |
| `GET` | `/events/:id` | Get event by UUID |
| `GET` | `/indicators` | List extracted IOCs (filterable by `type`) |
| `GET` | `/indicators/:type/:value` | Get deep investigation context (observations, enrichments, relationships) |
| `GET` | `/alerts` | List security alerts (filterable by `status`) |
| `GET` | `/alerts/:id` | Get alert by UUID |
| `PATCH`| `/alerts/:id` | Update alert triage status (`investigating`, `resolved`, `false_positive`) |
| `GET` | `/stats` | Aggregated dashboard statistics |
| `GET` | `/health` | System health check probe |

---

## 🗺️ Roadmap

### MVP (Completed)
- [x] Hexagonal Architecture (Ports and Adapters)
- [x] Normalized security event ingestion
- [x] Rule-based threat detection (SSH brute force, port scan, SQLi, XSS, malware downloads)
- [x] Automated IOC extraction (IP, Domain, URL, SHA256, MD5)
- [x] Multi-source threat intelligence enrichment (GeoIP + VirusTotal v3)
- [x] Redis 24h caching layer
- [x] RabbitMQ asynchronous worker pipeline
- [x] Multi-factor explainable risk scoring (0–100)
- [x] Threat relationship graph correlation
- [x] MITRE ATT&CK technique mapping
- [x] Autonomous threshold-based alerting ($\text{Risk} \ge 70$)
- [x] Next.js SOC analyst dashboard (Overview, IOC Explorer, Deep Investigation, Alert Center)
- [x] Telemetry simulator CLI (`cmd/collector`)
- [x] Docker deployment setup

### Next Iterations
- [ ] Prometheus metrics & Grafana dashboard integration
- [ ] STIX 2.1 / TAXII 2.1 threat feed consumer
- [ ] Suricata & Zeek log parsers
- [ ] Role-Based Access Control (RBAC) & authentication

---

## 📜 License

Licensed under the [Apache 2.0 License](LICENSE).
