-- Create active_games table for persisting ongoing game state
-- Enables game recovery after server crashes and player reconnection
CREATE TABLE IF NOT EXISTS active_games (
    id SERIAL PRIMARY KEY,
    game_id VARCHAR(255) UNIQUE NOT NULL,
    table_id VARCHAR(255) NOT NULL,
    game_type VARCHAR(100) NOT NULL,
    players JSONB NOT NULL,           -- ["player1", "player2"]
    game_state BYTEA NOT NULL,        -- Serialized gameStateSnapshot (gob encoded)
    turn_number INT NOT NULL DEFAULT 0,
    state VARCHAR(50) NOT NULL,       -- STARTING, MULLIGAN, IN_PROGRESS, PAUSED, FINISHED
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Index for looking up games by player (searches in JSONB array)
CREATE INDEX idx_active_games_players ON active_games USING GIN (players);

-- Index for looking up by game_id (already unique, but explicit index for clarity)
CREATE INDEX idx_active_games_game_id ON active_games(game_id);

-- Index for looking up by table_id
CREATE INDEX idx_active_games_table_id ON active_games(table_id);

-- Index for querying by state
CREATE INDEX idx_active_games_state ON active_games(state);

-- Index for cleanup queries (finding old/stale games)
CREATE INDEX idx_active_games_updated_at ON active_games(updated_at);


