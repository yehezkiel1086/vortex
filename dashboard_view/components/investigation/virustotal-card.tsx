import { ShieldAlert, Bug, Tag, ExternalLink } from "lucide-react"
import { Card, CardHeader, CardTitle, CardContent } from "@/components/ui/card"
import { Badge } from "@/components/ui/badge"
import type { ThreatIntelData } from "@/lib/types"

interface VirusTotalCardProps {
  data?: ThreatIntelData
  fetchedAt?: string
  indicatorValue: string
  indicatorType: string
}

export function VirusTotalCard({
  data,
  fetchedAt,
  indicatorValue,
  indicatorType,
}: VirusTotalCardProps) {
  const isMalicious = (data?.malicious_votes ?? 0) > 0 || (data?.reputation ?? 0) < 0

  const vtURL =
    indicatorType === "ip"
      ? `https://www.virustotal.com/gui/ip-address/${indicatorValue}`
      : `https://www.virustotal.com/gui/file/${indicatorValue}`

  return (
    <Card className="bg-card/80 border-border/60">
      <CardHeader className="p-4 pb-3 flex flex-row items-center justify-between">
        <CardTitle className="text-xs font-semibold text-foreground flex items-center gap-1.5">
          <ShieldAlert className="size-3.5 text-red-400" />
          VirusTotal Threat Intelligence
        </CardTitle>
        <div className="flex items-center gap-2">
          <Badge variant="outline" className="text-[0.65rem] font-mono text-emerald-400 border-emerald-500/30">
            Cached 24h
          </Badge>
          <a
            href={vtURL}
            target="_blank"
            rel="noopener noreferrer"
            className="text-muted-foreground hover:text-foreground"
            title="View on VirusTotal"
          >
            <ExternalLink className="size-3" />
          </a>
        </div>
      </CardHeader>

      <CardContent className="p-4 pt-0 space-y-3">
        {/* Detection Ratios */}
        <div className="grid grid-cols-3 gap-2">
          <div className="rounded-lg border border-border/40 bg-muted/20 p-2.5 text-center">
            <span className="text-[0.65rem] text-muted-foreground block">Malicious Detections</span>
            <span
              className={`text-base font-bold font-mono ${
                (data?.malicious_votes ?? 0) > 0 ? "text-red-400" : "text-muted-foreground"
              }`}
            >
              {data?.malicious_votes ?? 0}
            </span>
          </div>

          <div className="rounded-lg border border-border/40 bg-muted/20 p-2.5 text-center">
            <span className="text-[0.65rem] text-muted-foreground block">Harmless Votes</span>
            <span className="text-base font-bold font-mono text-emerald-400">
              {data?.harmless_votes ?? 0}
            </span>
          </div>

          <div className="rounded-lg border border-border/40 bg-muted/20 p-2.5 text-center">
            <span className="text-[0.65rem] text-muted-foreground block">Reputation Score</span>
            <span
              className={`text-base font-bold font-mono ${
                (data?.reputation ?? 0) < 0 ? "text-red-400" : "text-foreground"
              }`}
            >
              {data?.reputation ?? 0}
            </span>
          </div>
        </div>

        {/* Malware Family */}
        {data?.malware_family && (
          <div className="flex items-center gap-2 rounded-md border border-red-500/30 bg-red-500/10 p-2 text-xs text-red-300">
            <Bug className="size-3.5 shrink-0 text-red-400" />
            <span className="font-semibold">Malware Family:</span>
            <span className="font-mono font-bold text-red-200">{data.malware_family}</span>
          </div>
        )}

        {/* Tags */}
        {data?.tags && data.tags.length > 0 && (
          <div className="space-y-1">
            <span className="text-[0.65rem] text-muted-foreground flex items-center gap-1">
              <Tag className="size-2.5" /> Threat Tags
            </span>
            <div className="flex flex-wrap gap-1">
              {data.tags.map((tag) => (
                <Badge key={tag} variant="secondary" className="text-[0.65rem] font-mono">
                  {tag}
                </Badge>
              ))}
            </div>
          </div>
        )}

        {fetchedAt && (
          <div className="border-t border-border/40 pt-2 text-[0.65rem] text-muted-foreground">
            Analysis Timestamp: {new Date(fetchedAt).toLocaleString()}
          </div>
        )}
      </CardContent>
    </Card>
  )
}
