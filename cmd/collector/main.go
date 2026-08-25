package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/yehezkiel1086/vortex/internal/core/domain"
)

type scenario struct {
	Name  string
	Event domain.Event
}

func getScenarios() []scenario {
	now := time.Now().UTC()

	return []scenario{
		{
			Name: "SSH Brute Force Attack",
			Event: domain.Event{
				Timestamp:       now.Add(-2 * time.Minute),
				Source:          "honeypot-ssh-01",
				SourceIP:        "185.10.20.30",
				DestinationIP:   "192.168.1.50",
				SourcePort:      49152,
				DestinationPort: 22,
				Protocol:        "tcp",
				Username:        "root",
				AttackType:      "ssh_bruteforce",
				Severity:        domain.SeverityHigh,
				Confidence:      0.92,
				RawPayload:      json.RawMessage(`{"failed_attempts": 28, "auth_method": "password", "last_user": "root"}`),
			},
		},
		{
			Name: "Port Scan Activity",
			Event: domain.Event{
				Timestamp:       now.Add(-1 * time.Minute),
				Source:          "firewall-edge",
				SourceIP:        "185.10.20.30", // Same attacker IP -> Triggers Correlation
				DestinationIP:   "192.168.1.50",
				Protocol:        "tcp",
				AttackType:      "port_scan",
				Severity:        domain.SeverityMedium,
				Confidence:      0.85,
				RawPayload:      json.RawMessage(`{"ports_scanned": [21, 22, 23, 80, 443, 3306, 8080], "flags": "SYN"}`),
			},
		},
		{
			Name: "SQL Injection Web Exploit",
			Event: domain.Event{
				Timestamp:       now.Add(-30 * time.Second),
				Source:          "waf-ingress",
				SourceIP:        "198.51.100.45",
				DestinationIP:   "10.0.0.10",
				SourcePort:      54321,
				DestinationPort: 443,
				Protocol:        "https",
				Domain:          "portal.example-corp.com",
				URL:             "https://portal.example-corp.com/api/v1/users?id=1' UNION SELECT username, password_hash FROM admin_users--",
				AttackType:      "sqli",
				Severity:        domain.SeverityHigh,
				Confidence:      0.95,
				RawPayload:      json.RawMessage(`{"rule_id": "CRS-942-100", "payload_matched": "UNION SELECT"}`),
			},
		},
		{
			Name: "Malware Payload Download (C2)",
			Event: domain.Event{
				Timestamp:       now,
				Source:          "ids-sensor-02",
				SourceIP:        "203.0.113.88",
				DestinationIP:   "192.168.1.105",
				SourcePort:      44321,
				DestinationPort: 80,
				Protocol:        "http",
				Domain:          "evil-payload-drop.org",
				URL:             "http://evil-payload-drop.org/bin/trojan_dropper.exe",
				FileHash:        "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
				AttackType:      "malware_download",
				Severity:        domain.SeverityCritical,
				Confidence:      0.98,
				RawPayload:      json.RawMessage(`{"file_name": "trojan_dropper.exe", "mime_type": "application/x-dosexec"}`),
			},
		},
	}
}

func main() {
	apiURL := flag.String("url", "http://localhost:8080/api/v1/events", "Vortex API events ingestion endpoint")
	scenarioFlag := flag.String("scenario", "all", "Scenario to simulate: 'all', 'ssh', 'sqli', 'malware', 'scan'")
	delayMs := flag.Int("delay", 1000, "Delay in milliseconds between events when running 'all'")
	flag.Parse()

	log.Printf("[Collector] Vortex Security Telemetry Feeder starting...")
	log.Printf("[Collector] Target API: %s\n", *apiURL)

	scenarios := getScenarios()
	client := &http.Client{Timeout: 5 * time.Second}

	for _, s := range scenarios {
		if *scenarioFlag != "all" {
			match := false
			switch *scenarioFlag {
			case "ssh":
				match = s.Event.AttackType == "ssh_bruteforce"
			case "scan":
				match = s.Event.AttackType == "port_scan"
			case "sqli":
				match = s.Event.AttackType == "sqli"
			case "malware":
				match = s.Event.AttackType == "malware_download"
			}
			if !match {
				continue
			}
		}

		log.Printf("[Collector] 📡 Dispatching Scenario: '%s'...", s.Name)
		sendEvent(client, *apiURL, s.Event)

		if *scenarioFlag == "all" && *delayMs > 0 {
			time.Sleep(time.Duration(*delayMs) * time.Millisecond)
		}
	}

	log.Println("[Collector] ✅ Telemetry ingestion completed.")
}

func sendEvent(client *http.Client, url string, event domain.Event) {
	data, err := json.Marshal(event)
	if err != nil {
		log.Printf("[Collector] Error serializing event: %v", err)
		return
	}

	resp, err := client.Post(url, "application/json", bytes.NewBuffer(data))
	if err != nil {
		log.Printf("[Collector] ❌ Request failed: %v", err)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode == http.StatusCreated {
		log.Printf("[Collector] ✔️ [%d Created] Event ingested successfully: %s\n", resp.StatusCode, string(body))
	} else {
		log.Printf("[Collector] ⚠️ [%d] Unexpected response: %s\n", resp.StatusCode, string(body))
	}
}
