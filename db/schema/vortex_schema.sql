CREATE TABLE IF NOT EXISTS "events" (
  "id" uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  "timestamp" timestamptz NOT NULL,
  "source" varchar(100) NOT NULL,
  "source_ip" inet,
  "destination_ip" inet,
  "source_port" int,
  "destination_port" int,
  "protocol" varchar(20),
  "domain" varchar(255),
  "url" text,
  "file_hash" varchar(128),
  "username" varchar(100),
  "attack_type" varchar(100),
  "severity" varchar(20) NOT NULL DEFAULT 'informational',
  "confidence" float NOT NULL DEFAULT 0.0,
  "raw_payload" jsonb,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE IF NOT EXISTS "indicators" (
  "id" uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  "type" varchar(50) NOT NULL,
  "value" text NOT NULL,
  "first_seen" timestamptz NOT NULL DEFAULT (now()),
  "last_seen" timestamptz NOT NULL DEFAULT (now()),
  "risk_score" float NOT NULL DEFAULT 0.0,
  "confidence" float NOT NULL DEFAULT 0.0,
  "status" varchar(50) NOT NULL DEFAULT 'active',
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now()),
  CONSTRAINT indicators_type_value_unique UNIQUE ("type", "value")
);

CREATE TABLE IF NOT EXISTS "observations" (
  "id" uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  "indicator_id" uuid NOT NULL REFERENCES "indicators" ("id") ON DELETE CASCADE,
  "event_id" uuid NOT NULL REFERENCES "events" ("id") ON DELETE CASCADE,
  "attack_type" varchar(100) NOT NULL,
  "technique_id" varchar(50),
  "timestamp" timestamptz NOT NULL,
  "severity" varchar(20) NOT NULL DEFAULT 'medium',
  "confidence" float NOT NULL DEFAULT 0.0,
  "source" varchar(100),
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE IF NOT EXISTS "enrichments" (
  "id" uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  "indicator_id" uuid NOT NULL REFERENCES "indicators" ("id") ON DELETE CASCADE,
  "provider" varchar(100) NOT NULL,
  "data" jsonb NOT NULL,
  "fetched_at" timestamptz NOT NULL DEFAULT (now()),
  "expires_at" timestamptz,
  CONSTRAINT enrichments_indicator_provider_unique UNIQUE ("indicator_id", "provider")
);

CREATE TABLE IF NOT EXISTS "relationships" (
  "id" uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  "source_indicator_id" uuid NOT NULL REFERENCES "indicators" ("id") ON DELETE CASCADE,
  "target_indicator_id" uuid NOT NULL REFERENCES "indicators" ("id") ON DELETE CASCADE,
  "relationship_type" varchar(100) NOT NULL,
  "confidence" float NOT NULL DEFAULT 0.0,
  "first_seen" timestamptz NOT NULL DEFAULT (now()),
  "last_seen" timestamptz NOT NULL DEFAULT (now()),
  CONSTRAINT relationships_unique UNIQUE ("source_indicator_id", "target_indicator_id", "relationship_type")
);

CREATE TABLE IF NOT EXISTS "alerts" (
  "id" uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  "indicator_id" uuid NOT NULL REFERENCES "indicators" ("id") ON DELETE CASCADE,
  "event_id" uuid REFERENCES "events" ("id") ON DELETE SET NULL,
  "severity" varchar(20) NOT NULL DEFAULT 'high',
  "risk_score" float NOT NULL DEFAULT 0.0,
  "confidence" float NOT NULL DEFAULT 0.0,
  "title" varchar(255) NOT NULL,
  "description" text,
  "status" varchar(50) NOT NULL DEFAULT 'open',
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "resolved_at" timestamptz
);

-- Indexes for performance
CREATE INDEX IF NOT EXISTS idx_events_timestamp ON "events" ("timestamp" DESC);
CREATE INDEX IF NOT EXISTS idx_events_source_ip ON "events" ("source_ip");
CREATE INDEX IF NOT EXISTS idx_events_attack_type ON "events" ("attack_type");

CREATE INDEX IF NOT EXISTS idx_indicators_type_value ON "indicators" ("type", "value");
CREATE INDEX IF NOT EXISTS idx_indicators_risk_score ON "indicators" ("risk_score" DESC);
CREATE INDEX IF NOT EXISTS idx_indicators_status ON "indicators" ("status");
CREATE INDEX IF NOT EXISTS idx_indicators_last_seen ON "indicators" ("last_seen" DESC);

CREATE INDEX IF NOT EXISTS idx_observations_indicator ON "observations" ("indicator_id");
CREATE INDEX IF NOT EXISTS idx_observations_event ON "observations" ("event_id");
CREATE INDEX IF NOT EXISTS idx_observations_timestamp ON "observations" ("timestamp" DESC);

CREATE INDEX IF NOT EXISTS idx_enrichments_indicator ON "enrichments" ("indicator_id");
CREATE INDEX IF NOT EXISTS idx_enrichments_expires_at ON "enrichments" ("expires_at");

CREATE INDEX IF NOT EXISTS idx_relationships_source ON "relationships" ("source_indicator_id");
CREATE INDEX IF NOT EXISTS idx_relationships_target ON "relationships" ("target_indicator_id");

CREATE INDEX IF NOT EXISTS idx_alerts_status ON "alerts" ("status");
CREATE INDEX IF NOT EXISTS idx_alerts_indicator ON "alerts" ("indicator_id");
CREATE INDEX IF NOT EXISTS idx_alerts_created_at ON "alerts" ("created_at" DESC);
