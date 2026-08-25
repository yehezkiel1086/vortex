import { Badge } from "@/components/ui/badge"

const MITRE_NAMES: Record<string, string> = {
  T1110: "Brute Force",
  T1046: "Network Service Scanning",
  T1190: "Exploit Public-Facing App",
  "T1059.007": "JavaScript Execution",
  T1083: "File & Directory Discovery",
  T1105: "Ingress Tool Transfer",
  T1000: "Security Observation",
}

interface MitreTagProps {
  techniqueId?: string
  className?: string
}

export function MitreTag({ techniqueId, className }: MitreTagProps) {
  if (!techniqueId) return null

  const name = MITRE_NAMES[techniqueId] || "ATT&CK Technique"

  return (
    <a
      href={`https://attack.mitre.org/techniques/${techniqueId.replace(".", "/")}/`}
      target="_blank"
      rel="noopener noreferrer"
      className="inline-block hover:opacity-80 transition-opacity"
      title={`MITRE ATT&CK: ${name}`}
    >
      <Badge
        variant="outline"
        className={`font-mono text-[0.65rem] border-purple-500/30 bg-purple-500/10 text-purple-300 ${className}`}
      >
        <span className="font-semibold">{techniqueId}</span>
        <span className="ml-1 text-muted-foreground hidden sm:inline">— {name}</span>
      </Badge>
    </a>
  )
}
