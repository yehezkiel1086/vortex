"use client"

import { useState } from "react"
import { Send, Terminal, Zap, CheckCircle2, AlertCircle, X } from "lucide-react"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { api } from "@/lib/api"
import type { Event, Severity } from "@/lib/types"

interface IngestEventModalProps {
  isOpen: boolean
  onClose: () => void
  onSuccess?: () => void
}

const PRESET_SCENARIOS = [
  {
    name: "SSH Brute Force",
    icon: "🔐",
    event: {
      source: "honeypot-ssh-01",
      source_ip: "185.10.20.30",
      destination_port: 22,
      protocol: "tcp",
      attack_type: "ssh_bruteforce",
      severity: "high" as Severity,
      confidence: 0.92,
      username: "root",
      raw_payload: { attempts: 32, method: "password", last_user: "root" },
    },
  },
  {
    name: "SQL Injection",
    icon: "💉",
    event: {
      source: "waf-ingress",
      source_ip: "198.51.100.45",
      destination_port: 443,
      protocol: "https",
      domain: "portal.example-corp.com",
      url: "https://portal.example-corp.com/api/v1/users?id=1' UNION SELECT username, password_hash FROM admin_users--",
      attack_type: "sqli",
      severity: "high" as Severity,
      confidence: 0.95,
    },
  },
  {
    name: "Malware Download (C2)",
    icon: "🦠",
    event: {
      source: "ids-sensor-02",
      source_ip: "203.0.113.88",
      domain: "evil-payload-drop.org",
      url: "http://evil-payload-drop.org/bin/trojan_dropper.exe",
      file_hash: "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
      attack_type: "malware_download",
      severity: "critical" as Severity,
      confidence: 0.98,
    },
  },
  {
    name: "Port Scanning",
    icon: "📡",
    event: {
      source: "firewall-edge",
      source_ip: "185.10.20.30", // Co-occurring with SSH Brute Force
      protocol: "tcp",
      attack_type: "port_scan",
      severity: "medium" as Severity,
      confidence: 0.85,
    },
  },
]

export function IngestEventModal({ isOpen, onClose, onSuccess }: IngestEventModalProps) {
  const [source, setSource] = useState("custom-sensor")
  const [sourceIP, setSourceIP] = useState("185.10.20.30")
  const [attackType, setAttackType] = useState("ssh_bruteforce")
  const [severity, setSeverity] = useState<Severity>("high")
  const [url, setURL] = useState("")
  const [fileHash, setFileHash] = useState("")
  const [isLoading, setIsLoading] = useState(false)
  const [statusMessage, setStatusMessage] = useState<{ text: string; success: boolean } | null>(null)

  if (!isOpen) return null

  const handleSend = async (eventData: Partial<Event>) => {
    setIsLoading(true)
    setStatusMessage(null)
    try {
      await api.ingestEvent({
        ...eventData,
        timestamp: new Date().toISOString(),
      })
      setStatusMessage({ text: "Security event dispatched to Vortex Go pipeline!", success: true })
      onSuccess?.()
      setTimeout(() => {
        setStatusMessage(null)
        onClose()
      }, 1500)
    } catch (err: unknown) {
      const msg = err instanceof Error ? err.message : "Failed to ingest event"
      setStatusMessage({ text: msg, success: false })
    } finally {
      setIsLoading(false)
    }
  }

  const handleCustomSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    handleSend({
      source,
      source_ip: sourceIP,
      attack_type: attackType,
      severity,
      url: url || undefined,
      file_hash: fileHash || undefined,
      confidence: 0.9,
    })
  }

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/60 backdrop-blur-xs p-4">
      <div className="relative w-full max-w-lg rounded-xl border border-border/80 bg-card p-6 shadow-2xl space-y-5 animate-in fade-in zoom-in-95 duration-150">
        {/* Modal Header */}
        <div className="flex items-center justify-between border-b border-border/50 pb-3">
          <div className="flex items-center gap-2">
            <div className="flex size-7 items-center justify-center rounded-md bg-primary/20 text-primary">
              <Zap className="size-4" />
            </div>
            <div>
              <h3 className="text-sm font-semibold text-foreground">Security Telemetry Ingest</h3>
              <p className="text-[0.7rem] text-muted-foreground">Simulate incoming security sensor detections</p>
            </div>
          </div>
          <button
            onClick={onClose}
            className="rounded-md p-1 text-muted-foreground hover:bg-muted hover:text-foreground"
          >
            <X className="size-4" />
          </button>
        </div>

        {/* Status Alert */}
        {statusMessage && (
          <div
            className={`flex items-center gap-2 rounded-md p-2.5 text-xs ${
              statusMessage.success
                ? "bg-emerald-500/10 text-emerald-400 border border-emerald-500/20"
                : "bg-red-500/10 text-red-400 border border-red-500/20"
            }`}
          >
            {statusMessage.success ? (
              <CheckCircle2 className="size-4 shrink-0" />
            ) : (
              <AlertCircle className="size-4 shrink-0" />
            )}
            <span>{statusMessage.text}</span>
          </div>
        )}

        {/* Quick Presets */}
        <div className="space-y-2">
          <label className="text-[0.7rem] font-semibold uppercase tracking-wider text-muted-foreground">
            ⚡ Quick Simulation Presets
          </label>
          <div className="grid grid-cols-2 gap-2">
            {PRESET_SCENARIOS.map((p) => (
              <Button
                key={p.name}
                type="button"
                variant="outline"
                size="sm"
                disabled={isLoading}
                onClick={() => handleSend(p.event)}
                className="justify-start text-left h-auto py-2 px-3 border-border/60 hover:bg-muted/60"
              >
                <span className="mr-2 text-sm">{p.icon}</span>
                <div className="flex flex-col">
                  <span className="text-xs font-semibold">{p.name}</span>
                  <span className="text-[0.65rem] text-muted-foreground">{p.event.source}</span>
                </div>
              </Button>
            ))}
          </div>
        </div>

        {/* Custom Form */}
        <form onSubmit={handleCustomSubmit} className="space-y-3 pt-2 border-t border-border/40">
          <label className="text-[0.7rem] font-semibold uppercase tracking-wider text-muted-foreground flex items-center gap-1">
            <Terminal className="size-3" /> Custom Telemetry Payload
          </label>

          <div className="grid grid-cols-2 gap-3">
            <div>
              <label className="text-[0.65rem] text-muted-foreground">Source</label>
              <Input
                value={source}
                onChange={(e) => setSource(e.target.value)}
                placeholder="honeypot / waf / ids"
                required
              />
            </div>
            <div>
              <label className="text-[0.65rem] text-muted-foreground">Source IP</label>
              <Input
                value={sourceIP}
                onChange={(e) => setSourceIP(e.target.value)}
                placeholder="185.10.20.30"
              />
            </div>
          </div>

          <div className="grid grid-cols-2 gap-3">
            <div>
              <label className="text-[0.65rem] text-muted-foreground">Attack Type</label>
              <Input
                value={attackType}
                onChange={(e) => setAttackType(e.target.value)}
                placeholder="ssh_bruteforce / sqli / port_scan"
              />
            </div>
            <div>
              <label className="text-[0.65rem] text-muted-foreground">Severity</label>
              <select
                value={severity}
                onChange={(e) => setSeverity(e.target.value as Severity)}
                className="flex h-8 w-full rounded-md border border-input bg-transparent px-3 py-1 text-xs shadow-sm focus-visible:outline-none"
              >
                <option value="critical" className="bg-card text-foreground">Critical</option>
                <option value="high" className="bg-card text-foreground">High</option>
                <option value="medium" className="bg-card text-foreground">Medium</option>
                <option value="low" className="bg-card text-foreground">Low</option>
                <option value="informational" className="bg-card text-foreground">Informational</option>
              </select>
            </div>
          </div>

          <div className="flex justify-end gap-2 pt-2">
            <Button type="button" variant="ghost" size="sm" onClick={onClose}>
              Cancel
            </Button>
            <Button type="submit" size="sm" disabled={isLoading} className="gap-1.5">
              <Send className="size-3" />
              <span>{isLoading ? "Ingesting..." : "Send Event"}</span>
            </Button>
          </div>
        </form>
      </div>
    </div>
  )
}
