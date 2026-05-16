-- Add total_price column to reservations
ALTER TABLE reservations ADD COLUMN IF NOT EXISTS total_price DECIMAL(10,2) NOT NULL DEFAULT 0;

-- Create club_schedules table
CREATE TABLE IF NOT EXISTS club_schedules (
    id CHAR(36) PRIMARY KEY,
    day_of_week INT NOT NULL CHECK (day_of_week >= 0 AND day_of_week <= 6) UNIQUE,
    open_time TIME NOT NULL,
    close_time TIME NOT NULL,
    is_open BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create pricing_rules table
CREATE TABLE IF NOT EXISTS pricing_rules (
    id CHAR(36) PRIMARY KEY,
    sport_id CHAR(36) NOT NULL,
    day_of_week INT CHECK (day_of_week IS NULL OR (day_of_week >= 0 AND day_of_week <= 6)),
    start_hour INT NOT NULL CHECK (start_hour >= 0 AND start_hour <= 23),
    end_hour INT NOT NULL CHECK (end_hour >= 0 AND end_hour <= 23),
    price_per_hour DECIMAL(10,2) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (sport_id) REFERENCES sports(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_pricing_rules_sport_id ON pricing_rules(sport_id);
CREATE INDEX IF NOT EXISTS idx_pricing_rules_sport_day ON pricing_rules(sport_id, day_of_week);
