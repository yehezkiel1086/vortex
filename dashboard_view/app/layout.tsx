import type { Metadata } from "next"
import { Inter, JetBrains_Mono } from "next/font/google"
import "./globals.css"
import { AppShell } from "@/components/layout/app-shell"

const inter = Inter({
  subsets: ["latin"],
  variable: "--font-sans",
})

const mono = JetBrains_Mono({
  subsets: ["latin"],
  variable: "--font-geist-mono",
})

export const metadata: Metadata = {
  title: "Vortex — Threat Intelligence & Security Investigation",
  description:
    "Open-source threat intelligence and security investigation platform built with Go and Next.js.",
}

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode
}>) {
  return (
    <html lang="en" className="dark">
      <body className={`${inter.variable} ${mono.variable} font-sans antialiased bg-background text-foreground`}>
        <AppShell>{children}</AppShell>
      </body>
    </html>
  )
}
