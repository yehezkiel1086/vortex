"use client"

import * as React from "react"
import { cn } from "@/lib/utils"

interface TabsContextValue {
  value: string
  onValueChange: (value: string) => void
}

const TabsContext = React.createContext<TabsContextValue | undefined>(undefined)

function Tabs({
  defaultValue,
  value,
  onValueChange,
  className,
  children,
  ...props
}: {
  defaultValue?: string
  value?: string
  onValueChange?: (val: string) => void
} & React.ComponentProps<"div">) {
  const [tabValue, setTabValue] = React.useState(defaultValue || "")
  const current = value !== undefined ? value : tabValue

  const handleValueChange = React.useCallback(
    (val: string) => {
      if (value === undefined) setTabValue(val)
      onValueChange?.(val)
    },
    [value, onValueChange]
  )

  return (
    <TabsContext.Provider value={{ value: current, onValueChange: handleValueChange }}>
      <div data-slot="tabs" className={cn("w-full space-y-4", className)} {...props}>
        {children}
      </div>
    </TabsContext.Provider>
  )
}

function TabsList({ className, ...props }: React.ComponentProps<"div">) {
  return (
    <div
      data-slot="tabs-list"
      className={cn(
        "inline-flex h-8 items-center justify-center rounded-lg bg-muted p-1 text-muted-foreground border border-border/50",
        className
      )}
      {...props}
    />
  )
}

function TabsTrigger({
  value,
  className,
  children,
  ...props
}: { value: string } & React.ComponentProps<"button">) {
  const context = React.useContext(TabsContext)
  const isSelected = context?.value === value

  return (
    <button
      type="button"
      data-slot="tabs-trigger"
      data-state={isSelected ? "active" : "inactive"}
      className={cn(
        "inline-flex items-center justify-center whitespace-nowrap rounded-md px-3 py-1 text-xs font-medium ring-offset-background transition-all outline-none disabled:pointer-events-none disabled:opacity-50",
        isSelected
          ? "bg-background text-foreground shadow-sm font-semibold"
          : "hover:bg-background/50 hover:text-foreground",
        className
      )}
      onClick={() => context?.onValueChange(value)}
      {...props}
    >
      {children}
    </button>
  )
}

function TabsContent({
  value,
  className,
  children,
  ...props
}: { value: string } & React.ComponentProps<"div">) {
  const context = React.useContext(TabsContext)
  if (context?.value !== value) return null

  return (
    <div
      data-slot="tabs-content"
      className={cn("mt-2 ring-offset-background focus-visible:outline-none", className)}
      {...props}
    >
      {children}
    </div>
  )
}

export { Tabs, TabsList, TabsTrigger, TabsContent }
