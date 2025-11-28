-- Drop active_games table and indexes
DROP INDEX IF EXISTS idx_active_games_updated_at;
DROP INDEX IF EXISTS idx_active_games_state;
DROP INDEX IF EXISTS idx_active_games_table_id;
DROP INDEX IF EXISTS idx_active_games_game_id;
DROP INDEX IF EXISTS idx_active_games_players;
DROP TABLE IF EXISTS active_games;


