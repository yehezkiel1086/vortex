"use client"

import { useState, useEffect, useCallback } from "react"
import Link from "next/link"
import {
  Radar,
  Search,
  RefreshCw,
  Globe,
  Hash,
  Link as LinkIcon,
  Shield,
  ExternalLink,
  Copy,
  Check,
} from "lucide-react"
import { Card, CardHeader, CardTitle, CardDescription, CardContent } from "@/components/ui/card"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Table, TableHeader, TableBody, TableHead, TableRow, TableCell } from "@/components/ui/table"
import { Badge } from "@/components/ui/badge"
import { RiskMeter } from "@/components/dashboard/risk-meter"
import { api } from "@/lib/api"
import type { Indicator, IndicatorType } from "@/lib/types"

const TYPE_FILTERS: { label: string; value: string; icon?: string }[] = [
  { label: "All Types", value: "all" },
  { label: "IP Addresses", value: "ip" },
  { label: "Domains", value: "domain" },
  { label: "Hashes (SHA256/MD5)", value: "sha256" },
  { label: "URLs", value: "url" },
]

export default function IndicatorsPage() {
  const [indicators, setIndicators] = useState<Indicator[]>([])
  const [selectedType, setSelectedType] = useState<string>("all")
  const [searchQuery, setSearchQuery] = useState("")
  const [isLoading, setIsLoading] = useState(true)
  const [copiedValue, setCopiedValue] = useState<string | null>(null)

  const loadIndicators = useCallback(async () => {
    setIsLoading(true)
    try {
      const typeParam = selectedType === "all" ? undefined : selectedType
      const data = await api.getIndicators(typeParam, 50, 0)
      setIndicators(data || [])
    } catch (err) {
      console.error("Failed to load indicators:", err)
    } finally {
      setIsLoading(false)
    }
  }, [selectedType])

  useEffect(() => {
    loadIndicators()
  }, [loadIndicators])

  const copyToClipboard = (text: string) => {
    navigator.clipboard.writeText(text)
    setCopiedValue(text)
    setTimeout(() => setCopiedValue(null), 2000)
  }

  const filteredIndicators = indicators.filter((ind) => {
    if (!searchQuery) return true
    const q = searchQuery.toLowerCase()
    return ind.value.toLowerCase().includes(q) || ind.type.toLowerCase().includes(q)
  })

  const getTypeIcon = (type: IndicatorType | string) => {
    switch (type) {
      case "ip":
        return <Globe className="size-3.5 text-blue-400" />
      case "domain":
        return <Globe className="size-3.5 text-purple-400" />
      case "sha256":
      case "sha1":
      case "md5":
        return <Hash className="size-3.5 text-amber-400" />
      case "url":
        return <LinkIcon className="size-3.5 text-emerald-400" />
      default:
        return <Shield className="size-3.5 text-slate-400" />
    }
  }

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex flex-col gap-2 sm:flex-row sm:items-center sm:justify-between">
        <div>
          <h1 className="text-xl font-bold tracking-tight text-foreground flex items-center gap-2">
            <Radar className="size-5 text-purple-400" />
            Indicators of Compromise (IOC Explorer)
          </h1>
          <p className="text-xs text-muted-foreground">
            Search, filter, and explore extracted threat indicators and multi-factor risk scores.
          </p>
        </div>

        <Button
          variant="outline"
          size="sm"
          onClick={loadIndicators}
          disabled={isLoading}
          className="text-xs gap-1.5 h-8 border-border/60"
        >
          <RefreshCw className={`size-3 ${isLoading ? "animate-spin" : ""}`} />
          <span>Refresh</span>
        </Button>
      </div>

      {/* Filter and Search Bar */}
      <div className="flex flex-col sm:flex-row items-center justify-between gap-3">
        {/* Search Input */}
        <div className="relative w-full sm:w-80">
          <Search className="absolute left-2.5 top-1/2 size-3.5 -translate-y-1/2 text-muted-foreground" />
          <Input
            type="text"
            placeholder="Search indicator by value..."
            value={searchQuery}
            onChange={(e) => setSearchQuery(e.target.value)}
            className="pl-8 bg-muted/30 text-xs h-8 border-border/60"
          />
        </div>

        {/* Type Filter Buttons */}
        <div className="flex flex-wrap items-center gap-1.5 w-full sm:w-auto">
          {TYPE_FILTERS.map((filter) => (
            <Button
              key={filter.value}
              variant={selectedType === filter.value ? "default" : "outline"}
              size="xs"
              onClick={() => setSelectedType(filter.value)}
              className="text-xs h-7 border-border/60"
            >
              {filter.label}
            </Button>
          ))}
        </div>
      </div>

      {/* Indicators Table Card */}
      <Card className="bg-card/80 border-border/60 overflow-hidden">
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead>Indicator Value</TableHead>
              <TableHead>Type</TableHead>
              <TableHead className="w-48">Risk Score (0–100)</TableHead>
              <TableHead>Confidence</TableHead>
              <TableHead>Status</TableHead>
              <TableHead>First Seen</TableHead>
              <TableHead>Last Seen</TableHead>
              <TableHead className="text-right">Action</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {isLoading ? (
              <TableRow>
                <TableCell colSpan={8} className="h-32 text-center text-muted-foreground">
                  <div className="flex items-center justify-center gap-2 text-xs">
                    <RefreshCw className="size-4 animate-spin text-primary" />
                    <span>Loading indicators from database...</span>
                  </div>
                </TableCell>
              </TableRow>
            ) : filteredIndicators.length === 0 ? (
              <TableRow>
                <TableCell colSpan={8} className="h-32 text-center text-muted-foreground">
                  No indicators found matching the criteria. Ingest security telemetry to extract IOCs.
                </TableCell>
              </TableRow>
            ) : (
              filteredIndicators.map((ind) => (
                <TableRow key={ind.id} className="hover:bg-muted/30">
                  {/* Value */}
                  <TableCell>
                    <div className="flex items-center gap-2">
                      <span className="font-mono text-xs font-semibold text-foreground select-all">
                        {ind.value}
                      </span>
                      <button
                        onClick={() => copyToClipboard(ind.value)}
                        className="text-muted-foreground hover:text-foreground p-0.5 rounded transition-colors"
                        title="Copy to clipboard"
                      >
                        {copiedValue === ind.value ? (
                          <Check className="size-3 text-emerald-400" />
                        ) : (
                          <Copy className="size-3 opacity-60" />
                        )}
                      </button>
                    </div>
                  </TableCell>

                  {/* Type */}
                  <TableCell>
                    <Badge variant="outline" className="text-[0.65rem] font-mono gap-1 uppercase">
                      {getTypeIcon(ind.type)}
                      <span>{ind.type}</span>
                    </Badge>
                  </TableCell>

                  {/* Risk Score */}
                  <TableCell>
                    <RiskMeter score={ind.risk_score} showBar={true} />
                  </TableCell>

                  {/* Confidence */}
                  <TableCell className="font-mono text-xs text-muted-foreground">
                    {Math.round(ind.confidence * 100)}%
                  </TableCell>

                  {/* Status */}
                  <TableCell>
                    <Badge
                      variant={ind.status === "active" ? "success" : "secondary"}
                      className="text-[0.65rem] capitalize font-medium"
                    >
                      {ind.status}
                    </Badge>
                  </TableCell>

                  {/* Timestamps */}
                  <TableCell className="text-[0.7rem] font-mono text-muted-foreground whitespace-nowrap">
                    {new Date(ind.first_seen).toLocaleDateString()}
                  </TableCell>
                  <TableCell className="text-[0.7rem] font-mono text-muted-foreground whitespace-nowrap">
                    {new Date(ind.last_seen).toLocaleDateString()}
                  </TableCell>

                  {/* Action */}
                  <TableCell className="text-right">
                    <Link href={`/investigation/${ind.type}/${encodeURIComponent(ind.value)}`}>
                      <Button
                        variant="ghost"
                        size="xs"
                        className="gap-1 text-[0.65rem] hover:text-primary"
                      >
                        <span>Investigate</span>
                        <ExternalLink className="size-2.5" />
                      </Button>
                    </Link>
                  </TableCell>
                </TableRow>
              ))
            )}
          </TableBody>
        </Table>
      </Card>
    </div>
  )
}
