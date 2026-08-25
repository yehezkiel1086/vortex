"use client"

import { useState, useEffect, use } from "react"
import Link from "next/link"
import {
  ArrowLeft,
  Radar,
  RefreshCw,
  Globe,
  Hash,
  Link as LinkIcon,
  Shield,
  Clock,
  AlertTriangle,
  ExternalLink,
} from "lucide-react"
import { Button } from "@/components/ui/button"
import { Badge } from "@/components/ui/badge"
import { Card, CardHeader, CardTitle, CardContent } from "@/components/ui/card"
import { ScoreBreakdown } from "@/components/investigation/score-breakdown"
import { GeoIPCard } from "@/components/investigation/geoip-card"
import { VirusTotalCard } from "@/components/investigation/virustotal-card"
import { Timeline } from "@/components/investigation/timeline"
import { RelationshipGraph } from "@/components/investigation/relationship-graph"
import { api } from "@/lib/api"
import type { InvestigationContext, GeoIPData, ThreatIntelData } from "@/lib/types"

export default function InvestigationDetailPage({
  params,
}: {
  params: Promise<{ type: string; value: string }>
}) {
  const resolvedParams = use(params)
  const indType = decodeURIComponent(resolvedParams.type)
  const indValue = decodeURIComponent(resolvedParams.value)

  const [context, setContext] = useState<InvestigationContext | null>(null)
  const [isLoading, setIsLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  const loadData = async () => {
    setIsLoading(true)
    setError(null)
    try {
      const res = await api.getIndicatorDetails(indType, indValue)
      setContext(res)
    } catch (err: unknown) {
      const msg = err instanceof Error ? err.message : "Failed to load indicator details"
      setError(msg)
    } finally {
      setIsLoading(false)
    }
  }

  useEffect(() => {
    loadData()
  }, [indType, indValue])

  const geoIPEnrichment = context?.enrichments?.find((e) => e.provider === "geoip")
  const vtEnrichment = context?.enrichments?.find((e) => e.provider === "virustotal")

  return (
    <div className="space-y-6">
      {/* Top Breadcrumb & Action Header */}
      <div className="flex flex-col gap-2 sm:flex-row sm:items-center sm:justify-between">
        <div className="flex items-center gap-3">
          <Link href="/indicators">
            <Button variant="outline" size="xs" className="gap-1 text-xs border-border/60">
              <ArrowLeft className="size-3" />
              <span>Back to IOCs</span>
            </Button>
          </Link>

          <div>
            <div className="flex items-center gap-2">
              <Badge variant="outline" className="text-[0.65rem] uppercase font-mono">
                {indType}
              </Badge>
              <h1 className="text-lg font-bold font-mono text-foreground select-all">
                {indValue}
              </h1>
            </div>
            <p className="text-[0.7rem] text-muted-foreground">
              Deep Threat Intelligence Context &amp; Correlation Graph
            </p>
          </div>
        </div>

        <Button
          variant="outline"
          size="sm"
          onClick={loadData}
          disabled={isLoading}
          className="text-xs gap-1.5 h-8 border-border/60"
        >
          <RefreshCw className={`size-3 ${isLoading ? "animate-spin" : ""}`} />
          <span>Refresh Intelligence</span>
        </Button>
      </div>

      {/* Loading State */}
      {isLoading && (
        <Card className="p-12 text-center text-xs text-muted-foreground bg-card/50">
          <div className="flex flex-col items-center justify-center gap-2">
            <RefreshCw className="size-6 animate-spin text-primary" />
            <span>Querying Vortex threat graph and intelligence cache...</span>
          </div>
        </Card>
      )}

      {/* Error / Not Found State */}
      {!isLoading && error && (
        <Card className="p-8 text-center bg-red-500/5 border-red-500/20 space-y-3">
          <AlertTriangle className="size-8 text-red-400 mx-auto" />
          <div className="space-y-1">
            <h3 className="text-sm font-semibold text-red-400">Indicator Not Found in Database</h3>
            <p className="text-xs text-muted-foreground">
              {indValue} has not been observed yet. Ingest telemetry containing this indicator to generate analysis.
            </p>
          </div>
        </Card>
      )}

      {/* Investigation Details Workspace */}
      {!isLoading && context?.indicator && (
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
          {/* Left Column: Risk Score Breakdown & External Intelligence */}
          <div className="space-y-6">
            {/* 5-Factor Score Model */}
            <ScoreBreakdown
              indicator={context.indicator}
              observationCount={context.observations?.length || 0}
              relationshipCount={context.relationships?.length || 0}
              hasMaliciousEnrichment={
                (vtEnrichment?.data as ThreatIntelData)?.malicious_votes ? true : false
              }
            />

            {/* GeoIP Intelligence */}
            {indType === "ip" && (
              <GeoIPCard
                data={geoIPEnrichment?.data as GeoIPData}
                fetchedAt={geoIPEnrichment?.fetched_at}
              />
            )}

            {/* VirusTotal Intelligence */}
            <VirusTotalCard
              data={vtEnrichment?.data as ThreatIntelData}
              fetchedAt={vtEnrichment?.fetched_at}
              indicatorValue={indValue}
              indicatorType={indType}
            />
          </div>

          {/* Right Column: Observation Timeline & Relationship Graph */}
          <div className="space-y-6">
            {/* Observation Timeline */}
            <Timeline observations={context.observations || []} />

            {/* Relationship Graph */}
            <RelationshipGraph
              currentIndicatorId={context.indicator.id}
              currentValue={indValue}
              relationships={context.relationships || []}
            />
          </div>
        </div>
      )}
    </div>
  )
}
