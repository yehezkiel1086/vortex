import { Clock, ShieldAlert, Radio } from "lucide-react"
import { Card, CardHeader, CardTitle, CardContent } from "@/components/ui/card"
import { SeverityBadge } from "@/components/dashboard/severity-badge"
import { MitreTag } from "@/components/dashboard/mitre-tag"
import { Badge } from "@/components/ui/badge"
import type { Observation } from "@/lib/types"

interface TimelineProps {
  observations: Observation[]
}

export function Timeline({ observations }: TimelineProps) {
  if (observations.length === 0) {
    return (
      <Card className="bg-card/80 border-border/60">
        <CardHeader className="p-4 pb-2">
          <CardTitle className="text-xs font-semibold text-foreground flex items-center gap-1.5">
            <Clock className="size-3.5 text-purple-400" />
            Attack Observations Timeline
          </CardTitle>
        </CardHeader>
        <CardContent className="p-4 pt-0 text-xs text-muted-foreground">
          No observations recorded for this indicator yet.
        </CardContent>
      </Card>
    )
  }

  return (
    <Card className="bg-card/80 border-border/60">
      <CardHeader className="p-4 pb-3 flex flex-row items-center justify-between">
        <CardTitle className="text-xs font-semibold text-foreground flex items-center gap-1.5">
          <Clock className="size-3.5 text-purple-400" />
          Attack Observations Timeline ({observations.length})
        </CardTitle>
        <Badge variant="outline" className="text-[0.65rem] font-mono">
          Chronological
        </Badge>
      </CardHeader>

      <CardContent className="p-4 pt-0 space-y-3">
        <div className="relative border-l border-border/60 ml-2.5 space-y-4 py-1">
          {observations.map((obs) => (
            <div key={obs.id} className="relative pl-5 group">
              {/* Bullet node */}
              <div className="absolute -left-1.5 top-1.5 size-3 rounded-full border-2 border-background bg-purple-500 group-hover:scale-125 transition-transform" />

              <div className="rounded-lg border border-border/40 bg-muted/20 p-3 space-y-2 hover:border-purple-500/30 transition-colors">
                <div className="flex flex-wrap items-center justify-between gap-2">
                  <div className="flex items-center gap-2">
                    <span className="font-semibold text-xs text-foreground capitalize">
                      {obs.attack_type.replace(/_/g, " ")}
                    </span>
                    <SeverityBadge severity={obs.severity} />
                  </div>
                  <span className="font-mono text-[0.65rem] text-muted-foreground">
                    {new Date(obs.timestamp).toLocaleString()}
                  </span>
                </div>

                <div className="flex flex-wrap items-center gap-2 text-xs">
                  {obs.technique_id && <MitreTag techniqueId={obs.technique_id} />}
                  {obs.source && (
                    <Badge variant="outline" className="text-[0.65rem] gap-1">
                      <Radio className="size-2.5" />
                      <span>{obs.source}</span>
                    </Badge>
                  )}
                  <span className="text-[0.65rem] text-muted-foreground font-mono">
                    Confidence: {Math.round(obs.confidence * 100)}%
                  </span>
                </div>
              </div>
            </div>
          ))}
        </div>
      </CardContent>
    </Card>
  )
}
