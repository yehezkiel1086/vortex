"use client"

import { useState } from "react"
import { useRouter } from "next/navigation"
import { Search, Globe, Hash, Shield, Sparkles, ArrowRight } from "lucide-react"
import { Card, CardHeader, CardTitle, CardDescription, CardContent } from "@/components/ui/card"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"

const SAMPLE_QUERIES = [
  { type: "ip", value: "185.10.20.30", label: "SSH Brute Force Attacker IP" },
  { type: "ip", value: "198.51.100.45", label: "WAF SQLi Exploit Source" },
  { type: "sha256", value: "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855", label: "Trojan Dropper Hash" },
  { type: "domain", value: "evil-payload-drop.org", label: "C2 Malware Domain" },
]

export default function InvestigationIndexPage() {
  const router = useRouter()
  const [query, setQuery] = useState("")

  const handleSearch = (e: React.FormEvent) => {
    e.preventDefault()
    const q = query.trim()
    if (!q) return

    let type = "domain"
    if (/^(?:[0-9]{1,3}\.){3}[0-9]{1,3}$/.test(q)) {
      type = "ip"
    } else if (/^[a-fA-F0-9]{32,64}$/.test(q)) {
      type = "sha256"
    }

    router.push(`/investigation/${type}/${encodeURIComponent(q)}`)
  }

  return (
    <div className="max-w-4xl mx-auto space-y-8 py-6">
      {/* Title */}
      <div className="text-center space-y-2">
        <h1 className="text-2xl font-bold tracking-tight text-foreground flex items-center justify-center gap-2">
          <Search className="size-6 text-primary" />
          Threat Intelligence Investigation Workspace
        </h1>
        <p className="text-xs text-muted-foreground max-w-lg mx-auto">
          Deep-dive into threat indicators, examine multi-factor explainable risk scoring, GeoIP metadata, VirusTotal analysis, and correlated relationships.
        </p>
      </div>

      {/* Main Search Input Card */}
      <Card className="bg-card/80 border-border/80 p-6 shadow-lg">
        <form onSubmit={handleSearch} className="flex gap-2">
          <div className="relative flex-1">
            <Search className="absolute left-3 top-1/2 size-4 -translate-y-1/2 text-muted-foreground" />
            <Input
              type="text"
              placeholder="Enter IP address, domain name, or file hash..."
              value={query}
              onChange={(e) => setQuery(e.target.value)}
              className="pl-9 h-10 text-sm bg-muted/30 border-border/60"
            />
          </div>
          <Button type="submit" size="default" className="h-10 px-5 gap-1.5 font-semibold">
            <span>Investigate</span>
            <ArrowRight className="size-4" />
          </Button>
        </form>

        {/* Quick Sample Queries */}
        <div className="mt-4 pt-4 border-t border-border/40">
          <span className="text-[0.7rem] text-muted-foreground flex items-center gap-1 mb-2 font-medium">
            <Sparkles className="size-3 text-amber-400" /> Quick Investigation Demos
          </span>
          <div className="grid grid-cols-1 sm:grid-cols-2 gap-2">
            {SAMPLE_QUERIES.map((sq) => (
              <button
                key={sq.value}
                onClick={() => router.push(`/investigation/${sq.type}/${encodeURIComponent(sq.value)}`)}
                className="flex items-center justify-between p-2.5 rounded-lg border border-border/40 bg-muted/20 hover:bg-muted/50 hover:border-primary/40 text-left transition-colors text-xs"
              >
                <div className="flex items-center gap-2 truncate">
                  {sq.type === "ip" ? (
                    <Globe className="size-3.5 text-blue-400 shrink-0" />
                  ) : sq.type === "sha256" ? (
                    <Hash className="size-3.5 text-amber-400 shrink-0" />
                  ) : (
                    <Shield className="size-3.5 text-purple-400 shrink-0" />
                  )}
                  <div className="truncate">
                    <div className="font-mono font-bold text-foreground truncate">{sq.value}</div>
                    <div className="text-[0.65rem] text-muted-foreground">{sq.label}</div>
                  </div>
                </div>
                <ArrowRight className="size-3 text-muted-foreground shrink-0 ml-2" />
              </button>
            ))}
          </div>
        </div>
      </Card>
    </div>
  )
}
