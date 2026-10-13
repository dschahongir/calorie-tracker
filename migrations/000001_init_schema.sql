-- migrations/000001_init_schema.sql

CREATE TABLE IF NOT EXISTS food_entries (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    calories DOUBLE PRECISION NOT NULL DEFAULT 0,
    protein DOUBLE PRECISION NOT NULL DEFAULT 0,
    carbs DOUBLE PRECISION NOT NULL DEFAULT 0,
    fat DOUBLE PRECISION NOT NULL DEFAULT 0,
    weight_grams DOUBLE PRECISION NOT NULL DEFAULT 0,
    is_from_image BOOLEAN NOT NULL DEFAULT FALSE,
    image_url TEXT,
    consumed_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_food_entries_consumed_at ON food_entries(consumed_at);

CREATE TABLE IF NOT EXISTS water_entries (
    id BIGSERIAL PRIMARY KEY,
    amount_ml INT NOT NULL CHECK (amount_ml > 0),
    consumed_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_water_entries_consumed_at ON water_entries(consumed_at);

CREATE TABLE IF NOT EXISTS body_metrics (
    id BIGSERIAL PRIMARY KEY,
    weight_kg DOUBLE PRECISION NOT NULL CHECK (weight_kg > 0),
    recorded_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_body_metrics_recorded_at ON body_metrics(recorded_at);
