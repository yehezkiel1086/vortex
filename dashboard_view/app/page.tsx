"use client"

import { useState, useEffect, useCallback } from "react"
import Link from "next/link"
import {
  Activity,
  Radar,
  AlertTriangle,
  Server,
  RefreshCw,
  ArrowUpRight,
  ShieldAlert,
  Clock,
  Radio,
  ExternalLink,
} from "lucide-react"
import { Card, CardHeader, CardTitle, CardDescription, CardContent } from "@/components/ui/card"
import { Button } from "@/components/ui/button"
import { Table, TableHeader, TableBody, TableHead, TableRow, TableCell } from "@/components/ui/table"
import { Badge } from "@/components/ui/badge"
import { SeverityBadge } from "@/components/dashboard/severity-badge"
import { RiskMeter } from "@/components/dashboard/risk-meter"
import { MitreTag } from "@/components/dashboard/mitre-tag"
import { api } from "@/lib/api"
import type { DashboardStats, Event, Alert } from "@/lib/types"

export default function SOCDashboardPage() {
  const [stats, setStats] = useState<DashboardStats | null>(null)
  const [events, setEvents] = useState<Event[]>([])
  const [alerts, setAlerts] = useState<Alert[]>([])
  const [isLoading, setIsLoading] = useState(true)
  const [lastRefreshed, setLastRefreshed] = useState<Date>(new Date())

  const loadDashboardData = useCallback(async () => {
    setIsLoading(true)
    try {
      const [statsData, eventsData, alertsData] = await Promise.allSettled([
        api.getStats(),
        api.getEvents(10, 0),
        api.getAlerts("open", 5, 0),
      ])

      if (statsData.status === "fulfilled") setStats(statsData.value)
      if (eventsData.status === "fulfilled") setEvents(eventsData.value || [])
      if (alertsData.status === "fulfilled") setAlerts(alertsData.value || [])

      setLastRefreshed(new Date())
    } catch (err) {
      console.error("Failed to load dashboard data:", err)
    } finally {
      setIsLoading(false)
    }
  }, [])

  useEffect(() => {
    loadDashboardData()
    const timer = setInterval(loadDashboardData, 15000)
    return () => clearInterval(timer)
  }, [loadDashboardData])

  const openAlertsCount =
    stats?.alerts_by_status?.open ?? alerts.filter((a) => a.status === "open").length

  return (
    <div className="space-y-6">
      {/* Page Header */}
      <div className="flex flex-col gap-2 sm:flex-row sm:items-center sm:justify-between">
        <div>
          <h1 className="text-xl font-bold tracking-tight text-foreground flex items-center gap-2">
            <Activity className="size-5 text-emerald-400" />
            Security Operations Overview
          </h1>
          <p className="text-xs text-muted-foreground">
            Real-time threat telemetry, indicator extraction, and autonomous risk scoring.
          </p>
        </div>

        <div className="flex items-center gap-3">
          <div className="text-[0.7rem] text-muted-foreground flex items-center gap-1">
            <Clock className="size-3" />
            <span>Updated {lastRefreshed.toLocaleTimeString()}</span>
          </div>
          <Button
            variant="outline"
            size="sm"
            onClick={loadDashboardData}
            disabled={isLoading}
            className="text-xs gap-1.5 h-8 border-border/60"
          >
            <RefreshCw className={`size-3 ${isLoading ? "animate-spin" : ""}`} />
            <span>Refresh</span>
          </Button>
        </div>
      </div>

      {/* Metric Cards Grid */}
      <div className="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-4">
        {/* Total Events */}
        <Card className="bg-card/70 border-border/60">
          <CardHeader className="flex flex-row items-center justify-between pb-2 p-4">
            <CardTitle className="text-xs font-medium text-muted-foreground">
              Total Ingested Events
            </CardTitle>
            <Server className="size-4 text-blue-400" />
          </CardHeader>
          <CardContent className="p-4 pt-0">
            <div className="text-2xl font-bold font-mono text-foreground">
              {stats?.total_events ?? events.length}
            </div>
            <p className="text-[0.65rem] text-muted-foreground mt-1 flex items-center gap-1">
              <span className="text-emerald-400 font-semibold">Normalized</span> across all sensors
            </p>
          </CardContent>
        </Card>

        {/* Active Indicators */}
        <Card className="bg-card/70 border-border/60">
          <CardHeader className="flex flex-row items-center justify-between pb-2 p-4">
            <CardTitle className="text-xs font-medium text-muted-foreground">
              Active Indicators (IOCs)
            </CardTitle>
            <Radar className="size-4 text-purple-400" />
          </CardHeader>
          <CardContent className="p-4 pt-0">
            <div className="text-2xl font-bold font-mono text-foreground">
              {stats?.total_indicators ?? "—"}
            </div>
            <p className="text-[0.65rem] text-muted-foreground mt-1">
              Extracted IPs, Domains & Hashes
            </p>
          </CardContent>
        </Card>

        {/* Open Critical/High Alerts */}
        <Card className="bg-card/70 border-border/60">
          <CardHeader className="flex flex-row items-center justify-between pb-2 p-4">
            <CardTitle className="text-xs font-medium text-muted-foreground">
              Open Security Alerts
            </CardTitle>
            <AlertTriangle className="size-4 text-red-400" />
          </CardHeader>
          <CardContent className="p-4 pt-0">
            <div className="text-2xl font-bold font-mono text-red-400">
              {openAlertsCount}
            </div>
            <p className="text-[0.65rem] text-muted-foreground mt-1">
              Requiring analyst triage
            </p>
          </CardContent>
        </Card>

        {/* Processing Pipeline Engine */}
        <Card className="bg-card/70 border-border/60">
          <CardHeader className="flex flex-row items-center justify-between pb-2 p-4">
            <CardTitle className="text-xs font-medium text-muted-foreground">
              Autonomous Pipeline
            </CardTitle>
            <Radio className="size-4 text-emerald-400 animate-pulse" />
          </CardHeader>
          <CardContent className="p-4 pt-0">
            <div className="text-sm font-bold font-mono text-emerald-400">
              Active &amp; Ingesting
            </div>
            <p className="text-[0.65rem] text-muted-foreground mt-1">
              Redis 24h Cache + RabbitMQ
            </p>
          </CardContent>
        </Card>
      </div>

      {/* Main Grid: Active Alerts & Live Telemetry Feed */}
      <div className="grid grid-cols-1 gap-6 lg:grid-cols-3">
        {/* Left Column: High-Risk Alerts */}
        <div className="lg:col-span-1 space-y-4">
          <div className="flex items-center justify-between">
            <h2 className="text-sm font-semibold text-foreground flex items-center gap-1.5">
              <ShieldAlert className="size-4 text-red-400" />
              High Risk Alerts
            </h2>
            <Link
              href="/alerts"
              className="text-[0.7rem] text-muted-foreground hover:text-foreground flex items-center gap-0.5"
            >
              View all <ArrowUpRight className="size-3" />
            </Link>
          </div>

          <div className="space-y-3">
            {alerts.length === 0 ? (
              <Card className="p-6 text-center text-xs text-muted-foreground bg-card/40 border-dashed">
                No open high-risk alerts. Use the telemetry feeder to simulate incoming attacks.
              </Card>
            ) : (
              alerts.map((alert) => (
                <Card
                  key={alert.id}
                  className="bg-card/80 border-border/60 hover:border-red-500/40 transition-colors"
                >
                  <CardHeader className="p-3 pb-2">
                    <div className="flex items-start justify-between gap-2">
                      <SeverityBadge severity={alert.severity} />
                      <span className="text-[0.65rem] text-muted-foreground font-mono">
                        {new Date(alert.created_at).toLocaleTimeString()}
                      </span>
                    </div>
                    <CardTitle className="text-xs font-semibold mt-1.5 leading-snug line-clamp-1">
                      {alert.title}
                    </CardTitle>
                  </CardHeader>
                  <CardContent className="p-3 pt-0 space-y-2">
                    <p className="text-[0.7rem] text-muted-foreground line-clamp-2">
                      {alert.description}
                    </p>
                    <RiskMeter score={alert.risk_score} />
                    <div className="pt-2 flex justify-end">
                      <Link href="/alerts">
                        <Button variant="outline" size="xs" className="gap-1 text-[0.65rem]">
                          <span>Triage Alert</span>
                          <ArrowUpRight className="size-2.5" />
                        </Button>
                      </Link>
                    </div>
                  </CardContent>
                </Card>
              ))
            )}
          </div>
        </div>

        {/* Right Column: Live Telemetry Stream */}
        <div className="lg:col-span-2 space-y-4">
          <div className="flex items-center justify-between">
            <h2 className="text-sm font-semibold text-foreground flex items-center gap-1.5">
              <Server className="size-4 text-blue-400" />
              Live Security Telemetry Stream
            </h2>
            <Badge variant="outline" className="text-[0.65rem] font-mono">
              Latest {events.length} events
            </Badge>
          </div>

          <Card className="bg-card/80 border-border/60 overflow-hidden">
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>Time</TableHead>
                  <TableHead>Source</TableHead>
                  <TableHead>Attacker (Source IP)</TableHead>
                  <TableHead>Attack Type</TableHead>
                  <TableHead>Severity</TableHead>
                  <TableHead className="text-right">Action</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {events.length === 0 ? (
                  <TableRow>
                    <TableCell colSpan={6} className="h-28 text-center text-muted-foreground">
                      No security telemetry received yet.
                    </TableCell>
                  </TableRow>
                ) : (
                  events.map((event) => (
                    <TableRow key={event.id} className="hover:bg-muted/30">
                      <TableCell className="font-mono text-[0.7rem] text-muted-foreground whitespace-nowrap">
                        {new Date(event.timestamp).toLocaleTimeString()}
                      </TableCell>
                      <TableCell className="font-medium text-xs whitespace-nowrap">
                        <Badge variant="outline" className="text-[0.65rem]">
                          {event.source}
                        </Badge>
                      </TableCell>
                      <TableCell className="font-mono text-xs text-foreground font-semibold">
                        {event.source_ip ? (
                          <Link
                            href={`/investigation/ip/${event.source_ip}`}
                            className="hover:underline hover:text-primary flex items-center gap-1"
                          >
                            <span>{event.source_ip}</span>
                            <ExternalLink className="size-2.5 opacity-60" />
                          </Link>
                        ) : (
                          <span className="text-muted-foreground">—</span>
                        )}
                      </TableCell>
                      <TableCell>
                        <div className="flex flex-col gap-0.5">
                          <span className="text-xs font-semibold capitalize">
                            {event.attack_type?.replace(/_/g, " ") || "Observation"}
                          </span>
                          {event.attack_type === "ssh_bruteforce" && <MitreTag techniqueId="T1110" />}
                          {event.attack_type === "port_scan" && <MitreTag techniqueId="T1046" />}
                          {event.attack_type === "sqli" && <MitreTag techniqueId="T1190" />}
                          {event.attack_type === "malware_download" && <MitreTag techniqueId="T1105" />}
                        </div>
                      </TableCell>
                      <TableCell>
                        <SeverityBadge severity={event.severity} />
                      </TableCell>
                      <TableCell className="text-right">
                        {event.source_ip && (
                          <Link href={`/investigation/ip/${event.source_ip}`}>
                            <Button variant="ghost" size="xs" className="h-6 px-2 text-[0.65rem]">
                              Investigate
                            </Button>
                          </Link>
                        )}
                      </TableCell>
                    </TableRow>
                  ))
                )}
              </TableBody>
            </Table>
          </Card>
        </div>
      </div>
    </div>
  )
}
