import { Card, CardHeader, CardTitle, CardContent } from "@/components/ui/card"
import { Progress } from "@/components/ui/progress"
import { getRiskLevel } from "@/components/dashboard/risk-meter"
import type { Indicator } from "@/lib/types"

interface ScoreBreakdownProps {
  indicator: Indicator
  observationCount: number
  relationshipCount: number
  hasMaliciousEnrichment: boolean
}

export function ScoreBreakdown({
  indicator,
  observationCount,
  relationshipCount,
  hasMaliciousEnrichment,
}: ScoreBreakdownProps) {
  const score = indicator.risk_score
  const level = getRiskLevel(score)

  // Calculate estimated factors based on the backend formula
  const reputationScore = hasMaliciousEnrichment ? Math.min(30, 25) : 0
  const severityScore = observationCount > 0 ? 20 : 0
  const freqScore = Math.min(20, observationCount * 5)
  const confScore = Math.round((indicator.confidence || 0.8) * 15 * 10) / 10
  const corrScore = Math.min(10, relationshipCount * 3.5)

  const factors = [
    { name: "External Reputation", score: reputationScore, max: 30, desc: "VirusTotal / Threat Feeds" },
    { name: "Attack Severity", score: severityScore, max: 25, desc: "Highest observed threat severity" },
    { name: "Attack Frequency", score: freqScore, max: 20, desc: `${observationCount} telemetry observations` },
    { name: "Confidence Weight", score: confScore, max: 15, desc: `${Math.round(indicator.confidence * 100)}% detection confidence` },
    { name: "Graph Correlation", score: corrScore, max: 10, desc: `${relationshipCount} linked indicators` },
  ]

  return (
    <Card className="bg-card/80 border-border/60">
      <CardHeader className="p-4 pb-3">
        <div className="flex items-center justify-between">
          <CardTitle className="text-xs font-semibold text-foreground uppercase tracking-wider">
            Explainable Risk Scoring Model
          </CardTitle>
          <div className="flex items-center gap-2">
            <span className={`px-2 py-0.5 rounded border text-[0.7rem] font-bold uppercase ${level.color}`}>
              {level.label}
            </span>
            <span className="font-mono text-sm font-bold text-foreground">
              {Math.round(score * 10) / 10}
              <span className="text-[0.65rem] text-muted-foreground font-normal">/100</span>
            </span>
          </div>
        </div>
      </CardHeader>
      <CardContent className="p-4 pt-0 space-y-3">
        <p className="text-[0.7rem] text-muted-foreground">
          Vortex calculates risk transparently from multi-source observations, telemetry frequency, and threat intelligence.
        </p>

        <div className="space-y-2.5 pt-1">
          {factors.map((f) => (
            <div key={f.name} className="space-y-1">
              <div className="flex items-center justify-between text-xs">
                <div>
                  <span className="font-medium text-foreground">{f.name}</span>
                  <span className="ml-2 text-[0.65rem] text-muted-foreground">({f.desc})</span>
                </div>
                <span className="font-mono text-[0.7rem] text-muted-foreground">
                  <span className="font-semibold text-foreground">{f.score}</span>/{f.max}
                </span>
              </div>
              <Progress
                value={(f.score / f.max) * 100}
                className="h-1.5 bg-muted/60"
                indicatorClassName={level.barColor}
              />
            </div>
          ))}
        </div>
      </CardContent>
    </Card>
  )
}
