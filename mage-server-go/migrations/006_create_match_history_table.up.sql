-- Create match_history table for tracking completed games
CREATE TABLE IF NOT EXISTS match_history (
    id BIGSERIAL PRIMARY KEY,
    game_id VARCHAR(255) NOT NULL,
    table_id VARCHAR(255),
    tournament_id VARCHAR(255),

    -- Players involved (JSON array of player info)
    players JSONB NOT NULL,

    -- Game information
    game_type VARCHAR(100) NOT NULL,
    start_time TIMESTAMP NOT NULL,
    end_time TIMESTAMP NOT NULL,
    duration_seconds INTEGER NOT NULL,

    -- Winner information
    winner_id BIGINT REFERENCES authorized_users(id) ON DELETE SET NULL,
    winner_name VARCHAR(255),

    -- Game settings
    match_options JSONB, -- MatchOptions from proto

    -- Replay data (optional, could be large)
    replay_data TEXT, -- Compressed game log for replay

    created_at TIMESTAMP NOT NULL DEFAULT NOW(),

    CONSTRAINT match_duration_positive CHECK (duration_seconds >= 0)
);

-- Index for querying by player (searches in JSONB array)
CREATE INDEX idx_match_history_players ON match_history USING GIN (players);

-- Index for querying by winner
CREATE INDEX idx_match_history_winner_id ON match_history(winner_id);

-- Index for querying recent matches
CREATE INDEX idx_match_history_created_at ON match_history(created_at DESC);

-- Index for querying by game type
CREATE INDEX idx_match_history_game_type ON match_history(game_type);

-- Index for querying by table/tournament
CREATE INDEX idx_match_history_table_id ON match_history(table_id);
CREATE INDEX idx_match_history_tournament_id ON match_history(tournament_id);

-- Index for time-based queries
CREATE INDEX idx_match_history_end_time ON match_history(end_time DESC);
