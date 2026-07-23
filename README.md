# Vortex

> Open-source threat intelligence and security investigation platform built with Go.

Vortex is a threat intelligence platform designed to collect security observations, extract Indicators of Compromise (IOCs), enrich them with external intelligence, correlate related activity, and provide analysts with actionable threat context.

The project explores the intersection of **Go backend engineering, distributed systems, DevOps, and cybersecurity**.

---

## ⚠️ Project Status

**Vortex is currently under active development.**

The initial MVP is being developed as a focused 2-week engineering project.

Vortex is **not intended to replace enterprise platforms** such as Splunk, Microsoft Sentinel, CrowdStrike, or commercial Threat Intelligence Platforms.

Instead, the MVP focuses on demonstrating the core mechanics behind a modern threat intelligence workflow:

```text
Security Telemetry
        │
        ▼
     Ingestion
        │
        ▼
     Detection
        │
        ▼
   IOC Extraction
        │
        ▼
     Enrichment
        │
        ▼
     Correlation
        │
        ▼
   Risk Assessment
        │
        ▼
      Alert
        │
        ▼
 Analyst Investigation
```

---

# Features

## 🔍 Security Event Ingestion

Vortex accepts normalized security events from different security sources.

Planned sources include:

- Honeypots
- IDS
- WAF
- Web application logs
- Authentication logs
- Network security tools
- Custom security sensors

All incoming events are normalized into a common internal event format.

Example:

```json
{
  "timestamp": "2026-08-25T10:00:00Z",
  "source": "honeypot",
  "source_ip": "185.10.20.30",
  "destination_port": 22,
  "protocol": "tcp",
  "attack_type": "ssh_bruteforce",
  "severity": "high",
  "confidence": 0.90
}
```

---

# 🛡️ Attack Detection

Vortex includes lightweight rule-based detection for common attack patterns.

Initial MVP detections include:

- SSH brute force
- Port scanning
- Credential stuffing
- SQL injection
- Cross-site scripting (XSS)
- Path traversal
- Malware download observations
- Suspicious command execution

Vortex intentionally does not attempt to become a complete IDS.

Instead, detection produces **security observations** that can subsequently be enriched and correlated.

---

# 🧩 IOC Extraction

Vortex automatically extracts Indicators of Compromise from security events.

Supported IOC types include:

| Indicator | Example |
|---|---|
| IPv4 | `185.10.20.30` |
| IPv6 | `2001:db8::1` |
| Domain | `evil-example.com` |
| URL | `https://evil-example.com/payload` |
| SHA256 | `abc123...` |
| SHA1 | `abc123...` |
| MD5 | `abc123...` |
| Email | `attacker@example.com` |
| ASN | `AS12345` |

Each indicator receives its own historical record.

---

# 🌐 Threat Intelligence Enrichment

Vortex enriches discovered indicators with external intelligence.

For example:

```text
185.10.20.30
       │
       ├── Reputation
       ├── GeoIP
       ├── ASN
       ├── WHOIS
       ├── Threat Feeds
       └── Historical Observations
```

Hash enrichment can provide:

```text
SHA256
  │
  └── External Malware Intelligence
        │
        ├── Detection ratio
        ├── File type
        ├── Reputation
        └── Malware family
```

External intelligence is cached to reduce unnecessary API requests.

---

# 🔗 Threat Correlation

Vortex correlates indicators, events, and observations to identify relationships between seemingly separate security events.

Example:

```text
                185.10.20.30
                 /    |    \
                /     |     \
               ▼      ▼      ▼
          Port Scan SSH BF   SQLi
                       │
                       ▼
                  malware.exe
                       │
                       ▼
                 SHA256 abc123
                       │
                       ▼
              evil-example.com
```

This allows analysts to move from an individual IOC toward a broader understanding of the activity surrounding it.

---

# 📊 Risk Scoring

Vortex calculates an explainable risk score from multiple signals.

Example:

```text
Reputation       0–30
Attack Severity  0–25
Frequency        0–20
Confidence       0–15
Correlation      0–10
----------------------
Maximum          100
```

Risk levels:

| Score | Level |
|---:|---|
| 0–29 | Informational |
| 30–49 | Low |
| 50–69 | Medium |
| 70–84 | High |
| 85–100 | Critical |

Risk and confidence are intentionally treated as separate concepts.

For example:

```text
Risk:       91/100
Confidence: 65%
```

means the potential threat is significant, but the available evidence is not yet highly conclusive.

---

# 🎯 MITRE ATT&CK Mapping

Observed behaviors can be mapped to MITRE ATT&CK techniques.

Examples:

```text
SSH Brute Force
      ↓
T1110 — Brute Force

Port Scanning
      ↓
T1046 — Network Service Scanning

Exploit Public-Facing Application
      ↓
T1190
```

This provides additional context for analysts investigating an indicator or attack.

---

# 🚨 Alerting

Vortex generates alerts when indicators or observations meet configured risk thresholds.

Example:

```text
┌─────────────────────────────────────┐
│ 🚨 HIGH RISK THREAT                 │
├─────────────────────────────────────┤
│ IP:         185.10.20.30            │
│ Risk:       91 / 100                │
│ Confidence: 94%                     │
│                                     │
│ Observed:                           │
│ • Port scanning                     │
│ • SSH brute force                   │
│ • Malware download                  │
│                                     │
│ ATT&CK:                             │
│ • T1046                             │
│ • T1110                             │
└─────────────────────────────────────┘
```

Alert states:

- `open`
- `investigating`
- `resolved`
- `false_positive`

---

# 🕸️ Relationship Investigation

Vortex provides an investigation view for exploring relationships between indicators.

Example:

```text
                    ┌──────────────┐
                    │ 185.10.20.30 │
                    └──────┬───────┘
                 ┌─────────┼─────────┐
                 │         │         │
                 ▼         ▼         ▼
             Port Scan   SSH BF     SQLi
                           │
                           ▼
                      malware.exe
                           │
                           ▼
                     SHA256 abc123
                           │
                           ▼
                  evil-example.com
```

The goal is to allow an analyst to answer:

> "What else is associated with this indicator?"

---

# 🏗️ Architecture

High-level MVP architecture:

```text
                         ┌──────────────────┐
                         │     Next.js      │
                         │    Dashboard     │
                         └────────┬─────────┘
                                  │
                                  │ REST API
                                  ▼
                         ┌──────────────────┐
                         │    Vortex API    │
                         │       Go         │
                         └────────┬─────────┘
                                  │
                 ┌────────────────┼────────────────┐
                 │                │                │
                 ▼                ▼                ▼
          ┌────────────┐   ┌───────────┐   ┌────────────┐
          │ PostgreSQL │   │   Redis   │   │ RabbitMQ   │
          └────────────┘   └───────────┘   └─────┬──────┘
                                                  │
                                                  ▼
                                         ┌────────────────┐
                                         │ Vortex Worker  │
                                         ├────────────────┤
                                         │ Detection      │
                                         │ Enrichment     │
                                         │ Correlation    │
                                         │ Risk Scoring   │
                                         └───────┬────────┘
                                                 │
                                  ┌──────────────┼──────────────┐
                                  ▼              ▼              ▼
                              GeoIP        External TI    Threat Feeds
```

Security telemetry enters through:

```text
Honeypot ──────┐
IDS ───────────┤
WAF ───────────┤
Application ───┤
Auth Logs ─────┤
Network Logs ──┤
               ▼
        Vortex Ingestion
```

---

# 🔄 Core Data Flow

A typical Vortex investigation looks like this:

```text
1. Attacker generates suspicious activity
                    │
                    ▼
2. Security sensor produces an event
                    │
                    ▼
3. Vortex ingests and normalizes event
                    │
                    ▼
4. Detection engine classifies behavior
                    │
                    ▼
5. IOC extraction identifies IP/domain/hash
                    │
                    ▼
6. IOC is stored
                    │
                    ▼
7. Enrichment worker queries external intelligence
                    │
                    ▼
8. Intelligence is cached and stored
                    │
                    ▼
9. Correlation engine searches historical activity
                    │
                    ▼
10. Risk + confidence are calculated
                    │
                    ▼
11. MITRE ATT&CK techniques are assigned
                    │
                    ▼
12. Alert is generated
                    │
                    ▼
13. Analyst investigates IOC
```

---

# 🧱 Technology Stack

## Backend

- Go
- Gin
- PostgreSQL
- Redis
- RabbitMQ

## Frontend

- Next.js
- TypeScript

## Infrastructure

- Docker
- Docker Compose
- Linux
- Nginx / reverse proxy

## Observability

- Prometheus
- Grafana

## Security Intelligence

- VirusTotal
- GeoIP
- ASN / WHOIS
- Public threat intelligence feeds
- MITRE ATT&CK

---

# 🗄️ Data Model

Core entities:

```text
Event
 │
 ├── Observation
 │
 └── Indicator
       │
       ├── Enrichment
       │
       └── Relationship
              │
              ├── Domain
              ├── IP
              ├── Hash
              ├── Malware
              └── Attack
```

Primary database tables:

```text
events
indicators
observations
enrichments
relationships
alerts
```

---

# 🚀 Running Locally

## Requirements

- Go
- Docker
- Docker Compose
- Node.js
- npm / pnpm

## Clone

```bash
git clone https://github.com/<your-username>/vortex-ti.git
cd vortex-ti
```

## Configure Environment

```bash
cp .env.example .env
```

Configure the required API keys and database credentials.

Example:

```env
DATABASE_URL=postgres://vortex:vortex@postgres:5432/vortex
REDIS_URL=redis://redis:6379
RABBITMQ_URL=amqp://vortex:vortex@rabbitmq:5672/

VIRUSTOTAL_API_KEY=
```

## Start Infrastructure

```bash
docker compose up -d
```

## Run Migrations

```bash
make migrate
```

## Start API

```bash
go run ./cmd/api
```

## Start Worker

```bash
go run ./cmd/worker
```

## Start Frontend

```bash
cd web
npm install
npm run dev
```

---

# 🧪 Example Event

Send a security event:

```bash
curl -X POST http://localhost:8080/api/v1/events \
  -H "Content-Type: application/json" \
  -d '{
    "timestamp": "2026-08-25T10:00:00Z",
    "source": "honeypot",
    "source_ip": "185.10.20.30",
    "destination_port": 22,
    "protocol": "tcp",
    "attack_type": "ssh_bruteforce",
    "severity": "high",
    "confidence": 0.9
  }'
```

Vortex processes the event:

```text
Event
  ↓
SSH brute-force detection
  ↓
IOC extraction
  ↓
185.10.20.30
  ↓
Enrichment
  ↓
Correlation
  ↓
Risk scoring
  ↓
Alert
```

---

# 📡 API

Base URL:

```text
/api/v1
```

## Events

```http
POST /events
GET  /events
GET  /events/:id
```

## Indicators

```http
GET /indicators
GET /indicators/:type/:value
GET /indicators/:id/timeline
GET /indicators/:id/relationships
```

## Alerts

```http
GET   /alerts
GET   /alerts/:id
PATCH /alerts/:id
```

## Statistics

```http
GET /stats
```

API documentation will be provided through OpenAPI/Swagger.

---

# 🔐 Security Considerations

Vortex is intended to process security telemetry and potentially untrusted data.

Security considerations include:

- Input validation
- Parameterized SQL queries
- API authentication
- Rate limiting
- Secret management
- Container hardening
- Dependency scanning
- Structured logging
- Avoiding unnecessary storage of sensitive information
- External API failure handling

Vortex should not be deployed directly to the public internet without appropriate authentication, TLS, network controls, and operational hardening.

---

# 📈 Observability

Vortex exposes application metrics for Prometheus.

Example metrics:

```text
vortex_events_ingested_total
vortex_events_processed_total
vortex_detection_total
vortex_enrichment_total
vortex_enrichment_errors_total
vortex_alerts_total
vortex_processing_latency
vortex_queue_depth
```

Grafana is used to visualize system health and processing performance.

---

# 🧪 Testing

Vortex uses multiple layers of testing.

## Unit Tests

- Event normalization
- IOC extraction
- Detection rules
- Risk scoring
- Confidence scoring
- Correlation

## Integration Tests

- API → PostgreSQL
- API → RabbitMQ
- Worker → Redis
- Worker → External Intelligence Providers

## End-to-End Tests

```text
Security Event
      ↓
Detection
      ↓
IOC
      ↓
Enrichment
      ↓
Correlation
      ↓
Risk
      ↓
Alert
      ↓
Investigation
```

---

# 🗺️ Roadmap

## MVP

- [x] Project architecture
- [ ] Event ingestion
- [ ] Security event normalization
- [ ] SSH brute-force detection
- [ ] Port scan detection
- [ ] Credential stuffing detection
- [ ] SQLi detection
- [ ] XSS detection
- [ ] Path traversal detection
- [ ] IOC extraction
- [ ] IP enrichment
- [ ] Domain enrichment
- [ ] Hash enrichment
- [ ] Redis caching
- [ ] RabbitMQ workers
- [ ] IOC correlation
- [ ] Risk scoring
- [ ] MITRE ATT&CK mapping
- [ ] Alerting
- [ ] Analyst dashboard
- [ ] IOC investigation
- [ ] Relationship visualization
- [ ] Prometheus metrics
- [ ] Grafana dashboard
- [ ] Docker deployment
- [ ] CI/CD

---

# 🔮 Future Roadmap

## V2

Potential V2 features include:

- STIX 2.1 support
- TAXII integration
- Additional threat intelligence feeds
- Zeek integration
- Suricata integration
- Automated IOC expiration
- Advanced correlation
- Threat campaign tracking
- Improved detection engine
- YARA integration
- Sigma rule support

## V3

Longer-term possibilities include:

- Distributed collectors
- Kubernetes deployment
- Multi-tenancy
- RBAC
- Advanced threat graphs
- ML-assisted anomaly detection
- Malware analysis pipeline
- Automated response / SOAR
- Threat actor and campaign modeling

---

# 🧠 Design Philosophy

Vortex is intentionally built around a separation of responsibilities.

```text
IDS
 ↓
"What happened?"

SIEM
 ↓
"What events are happening?"

Threat Intelligence
 ↓
"What do we know about these indicators?"

Vortex
 ↓
"How can we connect those observations and intelligence
into useful context for an analyst?"
```

Vortex therefore does not attempt to replace an IDS or SIEM.

Instead, it can consume observations from those systems and transform raw security events into contextualized threat intelligence.

---

# 🎯 MVP Demo

The intended demonstration is:

```text
                    ATTACKER
                       │
                       │
              ┌────────▼────────┐
              │    Honeypot     │
              └────────┬────────┘
                       │
                  SSH brute force
                       │
                       ▼
              ┌─────────────────┐
              │ Vortex Ingest   │
              └────────┬────────┘
                       │
                       ▼
              ┌─────────────────┐
              │    Detection    │
              └────────┬────────┘
                       │
                       ▼
                   IOC: IP
                       │
                       ▼
              ┌─────────────────┐
              │   Enrichment    │
              └────────┬────────┘
                       │
             ┌─────────┼─────────┐
             ▼         ▼         ▼
          GeoIP      ASN     Reputation
                       │
                       ▼
              ┌─────────────────┐
              │   Correlation   │
              └────────┬────────┘
                       │
                       ▼
                Risk: 91/100
                       │
                       ▼
              ┌─────────────────┐
              │      Alert      │
              └────────┬────────┘
                       │
                       ▼
                Analyst View
                       │
              ┌────────┼────────┐
              ▼        ▼        ▼
           Timeline  ATT&CK   Graph
```

---

# 🤝 Contributing

Contributions, issues, and security-related discussions are welcome.

Before contributing:

1. Fork the repository.
2. Create a feature branch.
3. Make your changes.
4. Add tests where appropriate.
5. Run the test suite.
6. Open a pull request.

---

# ⚠️ Disclaimer

Vortex is an educational and research-oriented security project.

Only use Vortex and its associated security testing components against systems and infrastructure that you own or have explicit authorization to test.

The authors are not responsible for misuse of the software.

---

# 📜 License

Vortex is intended to be released under the **Apache License 2.0**.

See [`LICENSE`](LICENSE) for details.

---

# 👤 Author

Built by **Yehezkiel Wiradhika**.

Focused on:

```text
Go Backend Engineering
        +
DevOps / Cloud
        +
Cybersecurity
        +
Distributed Systems
```

Vortex is part of an ongoing effort to explore production-oriented security engineering and backend infrastructure.
