-- Create decks table for storing user deck collections
CREATE TABLE IF NOT EXISTS decks (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES authorized_users(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    format VARCHAR(50) NOT NULL, -- e.g., 'Standard', 'Modern', 'Commander', etc.
    description TEXT,
    main_deck TEXT NOT NULL, -- JSON array of card names
    sideboard TEXT, -- JSON array of card names
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),

    CONSTRAINT deck_name_length CHECK (LENGTH(name) >= 1 AND LENGTH(name) <= 255),
    CONSTRAINT deck_format_not_empty CHECK (LENGTH(format) > 0)
);

-- Index for querying user's decks
CREATE INDEX idx_decks_user_id ON decks(user_id);

-- Index for searching by format
CREATE INDEX idx_decks_format ON decks(format);

-- Index for searching by name (for user's decks)
CREATE INDEX idx_decks_user_name ON decks(user_id, name);

-- Updated_at trigger
CREATE OR REPLACE FUNCTION update_decks_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_update_decks_updated_at
    BEFORE UPDATE ON decks
    FOR EACH ROW
    EXECUTE FUNCTION update_decks_updated_at();
