"use client"

import { useState } from "react"
import { Sidebar } from "./sidebar"
import { Header } from "./header"
import { IngestEventModal } from "@/components/modals/ingest-event-modal"

export function AppShell({ children }: { children: React.ReactNode }) {
  const [isIngestModalOpen, setIsIngestModalOpen] = useState(false)

  return (
    <div className="flex min-h-screen bg-background text-foreground">
      {/* Fixed Sidebar */}
      <Sidebar />

      {/* Main Content Area */}
      <div className="flex flex-1 flex-col pl-60 min-w-0">
        <Header onIngestClick={() => setIsIngestModalOpen(true)} />
        <main className="flex-1 p-6 max-w-7xl w-full mx-auto">{children}</main>
      </div>

      {/* Telemetry Ingest Modal */}
      <IngestEventModal
        isOpen={isIngestModalOpen}
        onClose={() => setIsIngestModalOpen(false)}
        onSuccess={() => {
          // Trigger any global refresh if needed
        }}
      />
    </div>
  )
}
