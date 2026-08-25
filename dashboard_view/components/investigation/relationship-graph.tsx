import Link from "next/link"
import { Share2, ArrowRight, Globe, Hash, Shield } from "lucide-react"
import { Card, CardHeader, CardTitle, CardContent } from "@/components/ui/card"
import { Badge } from "@/components/ui/badge"
import type { Relationship } from "@/lib/types"

interface RelationshipGraphProps {
  currentIndicatorId: string
  currentValue: string
  relationships: Relationship[]
}

export function RelationshipGraph({
  currentIndicatorId,
  currentValue,
  relationships,
}: RelationshipGraphProps) {
  if (relationships.length === 0) {
    return (
      <Card className="bg-card/80 border-border/60">
        <CardHeader className="p-4 pb-2">
          <CardTitle className="text-xs font-semibold text-foreground flex items-center gap-1.5">
            <Share2 className="size-3.5 text-blue-400" />
            Threat Correlation &amp; Relationship Graph
          </CardTitle>
        </CardHeader>
        <CardContent className="p-4 pt-0 text-xs text-muted-foreground">
          No correlated relationships discovered for this indicator yet.
        </CardContent>
      </Card>
    )
  }

  return (
    <Card className="bg-card/80 border-border/60">
      <CardHeader className="p-4 pb-3 flex flex-row items-center justify-between">
        <CardTitle className="text-xs font-semibold text-foreground flex items-center gap-1.5">
          <Share2 className="size-3.5 text-blue-400" />
          Threat Correlation Graph ({relationships.length} Links)
        </CardTitle>
        <Badge variant="outline" className="text-[0.65rem] font-mono">
          Correlated
        </Badge>
      </CardHeader>

      <CardContent className="p-4 pt-0 space-y-3">
        <p className="text-[0.7rem] text-muted-foreground">
          Entities correlated with <span className="font-mono text-foreground font-semibold">{currentValue}</span> across multi-sensor observations:
        </p>

        <div className="space-y-2">
          {relationships.map((rel) => {
            const isSource = rel.source_indicator_id === currentIndicatorId
            const relType = rel.relationship_type

            return (
              <div
                key={rel.id}
                className="flex items-center justify-between gap-3 rounded-lg border border-border/40 bg-muted/20 p-3 hover:border-blue-500/30 transition-colors text-xs"
              >
                {/* Source Node */}
                <div className="flex items-center gap-2">
                  <span className="font-mono font-bold text-foreground">
                    {isSource ? currentValue : "Related Entity"}
                  </span>
                </div>

                {/* Relationship Direction Badge */}
                <div className="flex items-center gap-1 px-2 py-0.5 rounded bg-blue-500/10 border border-blue-500/20 text-blue-300 font-mono text-[0.65rem] font-semibold">
                  <span>{relType}</span>
                  <ArrowRight className="size-3" />
                </div>

                {/* Target Node */}
                <div className="flex items-center gap-2">
                  <span className="font-mono font-bold text-purple-300">
                    {isSource ? "Target Entity" : currentValue}
                  </span>
                </div>

                {/* Confidence */}
                <div className="text-[0.65rem] font-mono text-muted-foreground">
                  Conf: {Math.round(rel.confidence * 100)}%
                </div>
              </div>
            )
          })}
        </div>
      </CardContent>
    </Card>
  )
}
