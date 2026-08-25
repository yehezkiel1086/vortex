import type {
  Alert,
  DashboardStats,
  Event,
  Indicator,
  InvestigationContext,
} from "./types"

const API_BASE = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080/api/v1"

async function fetcher<T>(endpoint: string, options?: RequestInit): Promise<T> {
  const url = `${API_BASE}${endpoint}`
  try {
    const res = await fetch(url, {
      ...options,
      headers: {
        "Content-Type": "application/json",
        ...options?.headers,
      },
      cache: "no-store",
    })

    if (!res.ok) {
      const errBody = await res.json().catch(() => ({}))
      throw new Error(errBody.error || `HTTP error ${res.status}: ${res.statusText}`)
    }

    const data = await res.json()
    return data.data !== undefined ? data.data : data
  } catch (err: unknown) {
    const message = err instanceof Error ? err.message : "Network request failed"
    console.error(`[API Error] ${endpoint}:`, message)
    throw new Error(message)
  }
}

export const api = {
  // Statistics
  getStats: async (): Promise<DashboardStats> => {
    return fetcher<DashboardStats>("/stats")
  },

  // Events
  getEvents: async (limit = 20, offset = 0): Promise<Event[]> => {
    return fetcher<Event[]>(`/events?limit=${limit}&offset=${offset}`)
  },

  getEventByID: async (id: string): Promise<Event> => {
    return fetcher<Event>(`/events/${id}`)
  },

  ingestEvent: async (event: Partial<Event>): Promise<Event> => {
    return fetcher<Event>("/events", {
      method: "POST",
      body: JSON.stringify(event),
    })
  },

  // Indicators
  getIndicators: async (type?: string, limit = 50, offset = 0): Promise<Indicator[]> => {
    const query = new URLSearchParams()
    if (type && type !== "all") query.set("type", type)
    query.set("limit", limit.toString())
    query.set("offset", offset.toString())
    return fetcher<Indicator[]>(`/indicators?${query.toString()}`)
  },

  getIndicatorDetails: async (type: string, value: string): Promise<InvestigationContext> => {
    return fetcher<InvestigationContext>(`/indicators/${type}/${encodeURIComponent(value)}`)
  },

  // Alerts
  getAlerts: async (status?: string, limit = 50, offset = 0): Promise<Alert[]> => {
    const query = new URLSearchParams()
    if (status && status !== "all") query.set("status", status)
    query.set("limit", limit.toString())
    query.set("offset", offset.toString())
    return fetcher<Alert[]>(`/alerts?${query.toString()}`)
  },

  getAlertByID: async (id: string): Promise<Alert> => {
    return fetcher<Alert>(`/alerts/${id}`)
  },

  updateAlertStatus: async (id: string, status: string): Promise<Alert> => {
    return fetcher<Alert>(`/alerts/${id}`, {
      method: "PATCH",
      body: JSON.stringify({ status }),
    })
  },

  // Health check
  checkHealth: async (): Promise<boolean> => {
    try {
      const res = await fetch("http://localhost:8080/health", { cache: "no-store" })
      return res.ok
    } catch {
      return false
    }
  },
}
