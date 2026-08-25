import { Badge } from "@/components/ui/badge"
import type { Severity } from "@/lib/types"

interface SeverityBadgeProps {
  severity: Severity | string
  className?: string
}

export function SeverityBadge({ severity, className }: SeverityBadgeProps) {
  const norm = (severity || "informational").toLowerCase()

  switch (norm) {
    case "critical":
      return (
        <Badge variant="critical" className={className}>
          <span className="mr-1 inline-block size-1.5 rounded-full bg-red-500 animate-pulse" />
          Critical
        </Badge>
      )
    case "high":
      return (
        <Badge variant="high" className={className}>
          <span className="mr-1 inline-block size-1.5 rounded-full bg-orange-500" />
          High
        </Badge>
      )
    case "medium":
      return (
        <Badge variant="medium" className={className}>
          <span className="mr-1 inline-block size-1.5 rounded-full bg-amber-500" />
          Medium
        </Badge>
      )
    case "low":
      return (
        <Badge variant="low" className={className}>
          <span className="mr-1 inline-block size-1.5 rounded-full bg-blue-500" />
          Low
        </Badge>
      )
    default:
      return (
        <Badge variant="info" className={className}>
          <span className="mr-1 inline-block size-1.5 rounded-full bg-slate-400" />
          Info
        </Badge>
      )
  }
}
