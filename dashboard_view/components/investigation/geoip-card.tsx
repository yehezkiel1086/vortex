import { Globe, MapPin, Server, Radio } from "lucide-react"
import { Card, CardHeader, CardTitle, CardContent } from "@/components/ui/card"
import { Badge } from "@/components/ui/badge"
import type { GeoIPData } from "@/lib/types"

interface GeoIPCardProps {
  data?: GeoIPData
  fetchedAt?: string
}

export function GeoIPCard({ data, fetchedAt }: GeoIPCardProps) {
  if (!data || (!data.country && !data.asn)) {
    return (
      <Card className="bg-card/80 border-border/60">
        <CardHeader className="p-4 pb-2">
          <CardTitle className="text-xs font-semibold text-foreground flex items-center gap-1.5">
            <Globe className="size-3.5 text-blue-400" />
            GeoIP &amp; Network Intelligence
          </CardTitle>
        </CardHeader>
        <CardContent className="p-4 pt-0 text-xs text-muted-foreground">
          No GeoIP records available for this indicator.
        </CardContent>
      </Card>
    )
  }

  return (
    <Card className="bg-card/80 border-border/60">
      <CardHeader className="p-4 pb-3 flex flex-row items-center justify-between">
        <CardTitle className="text-xs font-semibold text-foreground flex items-center gap-1.5">
          <Globe className="size-3.5 text-blue-400" />
          GeoIP &amp; Network Intelligence
        </CardTitle>
        <Badge variant="outline" className="text-[0.65rem] font-mono text-emerald-400 border-emerald-500/30">
          Cached 24h
        </Badge>
      </CardHeader>

      <CardContent className="p-4 pt-0 space-y-3">
        <div className="grid grid-cols-2 gap-3 text-xs">
          {/* Location */}
          <div className="space-y-1">
            <span className="text-[0.65rem] text-muted-foreground flex items-center gap-1">
              <MapPin className="size-2.5 text-muted-foreground" /> Location
            </span>
            <div className="font-medium text-foreground">
              {data.city ? `${data.city}, ` : ""}
              {data.country || "Unknown Country"}
              {data.country_code && (
                <span className="ml-1 text-[0.65rem] font-mono text-muted-foreground">
                  ({data.country_code})
                </span>
              )}
            </div>
          </div>

          {/* ASN */}
          <div className="space-y-1">
            <span className="text-[0.65rem] text-muted-foreground flex items-center gap-1">
              <Server className="size-2.5 text-muted-foreground" /> Autonomous System (ASN)
            </span>
            <div className="font-mono text-xs font-semibold text-purple-300 truncate">
              {data.asn || "N/A"}
            </div>
          </div>

          {/* ISP */}
          <div className="space-y-1">
            <span className="text-[0.65rem] text-muted-foreground flex items-center gap-1">
              <Radio className="size-2.5 text-muted-foreground" /> ISP / Operator
            </span>
            <div className="text-xs text-foreground truncate">
              {data.isp || data.org || "N/A"}
            </div>
          </div>

          {/* Coordinates */}
          <div className="space-y-1">
            <span className="text-[0.65rem] text-muted-foreground">Coordinates</span>
            <div className="font-mono text-[0.7rem] text-muted-foreground">
              {data.latitude && data.longitude
                ? `${data.latitude.toFixed(4)}, ${data.longitude.toFixed(4)}`
                : "N/A"}
            </div>
          </div>
        </div>

        {fetchedAt && (
          <div className="border-t border-border/40 pt-2 text-[0.65rem] text-muted-foreground">
            Fetched: {new Date(fetchedAt).toLocaleString()}
          </div>
        )}
      </CardContent>
    </Card>
  )
}
