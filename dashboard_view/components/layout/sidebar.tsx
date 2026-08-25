"use client"

import Link from "next/link"
import { usePathname } from "next/navigation"
import {
  Shield,
  Activity,
  Radar,
  AlertTriangle,
  Search,
  Server,
  Layers,
} from "lucide-react"
import { cn } from "@/lib/utils"

const navigation = [
  { name: "Overview", href: "/", icon: Activity },
  { name: "IOC Explorer", href: "/indicators", icon: Radar },
  { name: "Alert Center", href: "/alerts", icon: AlertTriangle },
  { name: "Investigation", href: "/investigation", icon: Search },
]

export function Sidebar() {
  const pathname = usePathname()

  return (
    <aside className="fixed inset-y-0 left-0 z-30 flex w-60 flex-col border-r border-border/50 bg-card/80 backdrop-blur-md">
      {/* Brand Header */}
      <div className="flex h-14 items-center gap-2.5 border-b border-border/40 px-5">
        <div className="flex size-8 items-center justify-center rounded-lg bg-primary text-primary-foreground shadow-sm">
          <Shield className="size-4.5 text-primary-foreground" />
        </div>
        <div className="flex flex-col">
          <span className="font-mono text-sm font-bold tracking-tight text-foreground">
            VORTEX
          </span>
          <span className="text-[0.65rem] font-medium text-muted-foreground uppercase tracking-wider">
            Threat Intelligence
          </span>
        </div>
      </div>

      {/* Navigation Links */}
      <div className="flex-1 overflow-y-auto px-3 py-4 space-y-1">
        <div className="px-2 pb-2 text-[0.65rem] font-semibold text-muted-foreground uppercase tracking-wider">
          Platform
        </div>
        {navigation.map((item) => {
          const isActive =
            item.href === "/"
              ? pathname === "/"
              : pathname.startsWith(item.href)

          const Icon = item.icon

          return (
            <Link
              key={item.name}
              href={item.href}
              className={cn(
                "group flex items-center gap-3 rounded-md px-3 py-2 text-xs font-medium transition-colors",
                isActive
                  ? "bg-primary text-primary-foreground font-semibold shadow-xs"
                  : "text-muted-foreground hover:bg-muted hover:text-foreground"
              )}
            >
              <Icon
                className={cn(
                  "size-4 shrink-0 transition-colors",
                  isActive
                    ? "text-primary-foreground"
                    : "text-muted-foreground group-hover:text-foreground"
                )}
              />
              <span>{item.name}</span>
            </Link>
          )
        })}

        <div className="pt-6 px-2 pb-2 text-[0.65rem] font-semibold text-muted-foreground uppercase tracking-wider">
          System Context
        </div>
        <div className="rounded-lg border border-border/40 bg-muted/20 p-3 space-y-2 text-xs">
          <div className="flex items-center justify-between text-[0.7rem]">
            <span className="text-muted-foreground flex items-center gap-1.5">
              <Server className="size-3 text-emerald-400" /> Go Engine
            </span>
            <span className="font-mono text-emerald-400 font-bold">:8080</span>
          </div>
          <div className="flex items-center justify-between text-[0.7rem]">
            <span className="text-muted-foreground flex items-center gap-1.5">
              <Layers className="size-3 text-purple-400" /> Architecture
            </span>
            <span className="font-mono text-purple-400 font-semibold">Hexagonal</span>
          </div>
        </div>
      </div>

      {/* Footer Info */}
      <div className="border-t border-border/40 p-4">
        <div className="flex items-center justify-between text-[0.65rem] text-muted-foreground">
          <span>Vortex v0.1.0-MVP</span>
          <span className="flex items-center gap-1">
            <span className="size-1.5 rounded-full bg-emerald-500 animate-pulse" />
            Live
          </span>
        </div>
      </div>
    </aside>
  )
}
