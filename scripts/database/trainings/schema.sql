-- Enable the uuid-ossp extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- This is a temporary table for local development
-- DO NOT USE THIS IN PRODUCTION
DROP TABLE IF EXISTS trainings;

CREATE TABLE IF NOT EXISTS trainings (
    uuid uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id uuid NOT NULL,
    training_datetime TIMESTAMP NOT NULL,
    notes TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);