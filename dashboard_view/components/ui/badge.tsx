import * as React from "react"
import { cva, type VariantProps } from "class-variance-authority"
import { cn } from "@/lib/utils"

const badgeVariants = cva(
  "inline-flex items-center justify-center rounded-md border px-2 py-0.5 text-xs font-medium w-fit whitespace-nowrap shrink-0 transition-colors focus:outline-none focus:ring-2 focus:ring-ring focus:ring-offset-2",
  {
    variants: {
      variant: {
        default:
          "border-transparent bg-primary text-primary-foreground shadow hover:bg-primary/80",
        secondary:
          "border-transparent bg-secondary text-secondary-foreground hover:bg-secondary/80",
        destructive:
          "border-transparent bg-destructive text-destructive-foreground shadow hover:bg-destructive/80",
        outline: "text-foreground",
        critical:
          "border-red-500/30 bg-red-500/10 text-red-400 dark:bg-red-950/40 dark:text-red-400 font-semibold",
        high: "border-orange-500/30 bg-orange-500/10 text-orange-400 dark:bg-orange-950/40 dark:text-orange-400 font-semibold",
        medium:
          "border-amber-500/30 bg-amber-500/10 text-amber-400 dark:bg-amber-950/40 dark:text-amber-400",
        low: "border-blue-500/30 bg-blue-500/10 text-blue-400 dark:bg-blue-950/40 dark:text-blue-400",
        info: "border-slate-500/30 bg-slate-500/10 text-slate-400 dark:bg-slate-950/40 dark:text-slate-400",
        success:
          "border-emerald-500/30 bg-emerald-500/10 text-emerald-400 dark:bg-emerald-950/40 dark:text-emerald-400",
      },
    },
    defaultVariants: {
      variant: "default",
    },
  }
)

export interface BadgeProps
  extends React.HTMLAttributes<HTMLDivElement>,
    VariantProps<typeof badgeVariants> {}

function Badge({ className, variant, ...props }: BadgeProps) {
  return (
    <div className={cn(badgeVariants({ variant }), className)} {...props} />
  )
}

export { Badge, badgeVariants }
