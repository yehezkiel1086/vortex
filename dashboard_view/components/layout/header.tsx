"use client"

import { useState, useEffect } from "react"
import { useRouter } from "next/navigation"
import { Search, Radio, RefreshCw, Send } from "lucide-react"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { api } from "@/lib/api"

interface HeaderProps {
  onIngestClick?: () => void
}

export function Header({ onIngestClick }: HeaderProps) {
  const router = useRouter()
  const [searchQuery, setSearchQuery] = useState("")
  const [isBackendHealthy, setIsBackendHealthy] = useState<boolean | null>(null)
  const [isChecking, setIsChecking] = useState(false)

  const checkStatus = async () => {
    setIsChecking(true)
    const healthy = await api.checkHealth()
    setIsBackendHealthy(healthy)
    setIsChecking(false)
  }

  useEffect(() => {
    checkStatus()
    const interval = setInterval(checkStatus, 30000)
    return () => clearInterval(interval)
  }, [])

  const handleSearch = (e: React.FormEvent) => {
    e.preventDefault()
    const q = searchQuery.trim()
    if (!q) return

    // Auto-detect type
    let type = "domain"
    if (/^(?:[0-9]{1,3}\.){3}[0-9]{1,3}$/.test(q)) {
      type = "ip"
    } else if (/^[a-fA-F0-9]{32,64}$/.test(q)) {
      type = "sha256"
    }

    router.push(`/investigation/${type}/${encodeURIComponent(q)}`)
  }

  return (
    <header className="sticky top-0 z-20 flex h-14 w-full items-center justify-between border-b border-border/40 bg-background/80 px-6 backdrop-blur-md">
      {/* Global IOC Search Bar */}
      <form onSubmit={handleSearch} className="relative w-72 sm:w-96">
        <Search className="absolute left-2.5 top-1/2 size-3.5 -translate-y-1/2 text-muted-foreground pointer-events-none" />
        <Input
          type="text"
          placeholder="Investigate IP, Domain, Hash (e.g. 185.10.20.30)..."
          value={searchQuery}
          onChange={(e) => setSearchQuery(e.target.value)}
          className="pl-8 bg-muted/30 text-xs h-8 focus-visible:bg-background border-border/60"
        />
      </form>

      {/* Right Controls */}
      <div className="flex items-center gap-3">
        {/* Backend Status Indicator */}
        <div
          onClick={checkStatus}
          className="hidden sm:flex items-center gap-1.5 rounded-full border border-border/60 bg-muted/40 px-2.5 py-1 text-[0.65rem] text-muted-foreground cursor-pointer hover:bg-muted transition-colors"
          title="Click to re-check Go backend connection"
        >
          <Radio
            className={`size-3 ${
              isBackendHealthy === true
                ? "text-emerald-400"
                : isBackendHealthy === false
                ? "text-red-400"
                : "text-amber-400"
            }`}
          />
          <span className="font-mono">
            {isBackendHealthy === true
              ? "API Connected"
              : isBackendHealthy === false
              ? "API Offline"
              : "Checking..."}
          </span>
          <RefreshCw className={`size-2.5 ml-0.5 ${isChecking ? "animate-spin" : ""}`} />
        </div>

        {/* Action Button: Ingest Telemetry */}
        <Button
          size="sm"
          onClick={onIngestClick}
          className="bg-primary text-primary-foreground hover:bg-primary/90 text-xs gap-1.5 shadow-xs"
        >
          <Send className="size-3" />
          <span>Ingest Telemetry</span>
        </Button>
      </div>
    </header>
  )
}
