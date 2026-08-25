# Vortex Dashboard

Web-based analyst interface for the Vortex Threat Intelligence and Security Investigation platform.

Built with **Next.js 16 (App Router)**, **Tailwind CSS v4**, and **shadcn/ui**.

---

## Features

- **SOC Overview (`/`)**: High-level telemetry statistics, real-time event stream with MITRE ATT&CK tagging, and active high-risk alerts.
- **IOC Explorer (`/indicators`)**: Filter and search through extracted indicators (IPs, domains, hashes, URLs) with risk scoring.
- **Deep Investigation (`/investigation/[type]/[value]`)**:
  - Explainable 5-factor risk scoring breakdown (Reputation, Severity, Frequency, Confidence, Correlation).
  - External threat intelligence (GeoIP geolocation, ASN, VirusTotal detection ratios, malware families).
  - Chronological observation timeline mapped to MITRE techniques.
  - Threat correlation graph showing related indicators and payloads.
- **Alert Center (`/alerts`)**: Incident triage dashboard supporting status transitions (`open` → `investigating` → `resolved` → `false_positive`).
- **Telemetry Ingestion Modal**: UI simulator to dispatch sample security telemetry (SSH brute force, SQL injection, malware downloads, port scans) directly to the Go API.

---

## Tech Stack

- **Framework**: Next.js 16 (App Router, Turbopack, React 19)
- **Styling**: Tailwind CSS v4, Lucide Icons
- **Components**: shadcn/ui base components
- **State & Networking**: Native fetch client connecting to the Vortex Go API

---

## Getting Started

### Prerequisites

Ensure the Vortex Go backend is running on `http://localhost:8080` (or configure via environment variable).

### Environment Variables

Create a `.env.local` file (optional, defaults to `http://localhost:8080/api/v1`):

```env
NEXT_PUBLIC_API_URL=http://localhost:8080/api/v1
```

### Installation

```bash
npm install
```

### Development Server

```bash
npm run dev
```

Open [http://localhost:3000](http://localhost:3000) in your browser.

### Production Build

```bash
npm run build
npm run start
```
