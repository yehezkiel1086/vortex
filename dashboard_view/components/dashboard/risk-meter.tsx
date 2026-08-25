import { cn } from "@/lib/utils"
import { Progress } from "@/components/ui/progress"

interface RiskMeterProps {
  score: number
  showBar?: boolean
  className?: string
}

export function getRiskLevel(score: number): {
  label: string
  color: string
  barColor: string
  textColor: string
} {
  if (score >= 85) {
    return {
      label: "Critical",
      color: "bg-red-500/20 text-red-400 border-red-500/30",
      barColor: "bg-red-500",
      textColor: "text-red-400",
    }
  }
  if (score >= 70) {
    return {
      label: "High",
      color: "bg-orange-500/20 text-orange-400 border-orange-500/30",
      barColor: "bg-orange-500",
      textColor: "text-orange-400",
    }
  }
  if (score >= 50) {
    return {
      label: "Medium",
      color: "bg-amber-500/20 text-amber-400 border-amber-500/30",
      barColor: "bg-amber-500",
      textColor: "text-amber-400",
    }
  }
  if (score >= 30) {
    return {
      label: "Low",
      color: "bg-blue-500/20 text-blue-400 border-blue-500/30",
      barColor: "bg-blue-500",
      textColor: "text-blue-400",
    }
  }
  return {
    label: "Info",
    color: "bg-slate-500/20 text-slate-400 border-slate-500/30",
    barColor: "bg-slate-500",
    textColor: "text-slate-400",
  }
}

export function RiskMeter({ score, showBar = true, className }: RiskMeterProps) {
  const level = getRiskLevel(score)
  const displayScore = Math.round(score * 10) / 10

  return (
    <div className={cn("flex flex-col gap-1.5", className)}>
      <div className="flex items-center justify-between gap-2">
        <span
          className={cn(
            "inline-flex items-center gap-1 rounded border px-1.5 py-0.5 text-[0.65rem] font-bold uppercase tracking-wider",
            level.color
          )}
        >
          {level.label}
        </span>
        <span className={cn("text-xs font-mono font-bold", level.textColor)}>
          {displayScore}
          <span className="text-[0.65rem] text-muted-foreground font-normal">/100</span>
        </span>
      </div>

      {showBar && (
        <Progress
          value={score}
          className="h-1.5 w-full bg-muted/60"
          indicatorClassName={level.barColor}
        />
      )}
    </div>
  )
}
