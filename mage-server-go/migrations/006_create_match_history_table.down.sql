-- Drop match_history table
DROP INDEX IF EXISTS idx_match_history_end_time;
DROP INDEX IF EXISTS idx_match_history_tournament_id;
DROP INDEX IF EXISTS idx_match_history_table_id;
DROP INDEX IF EXISTS idx_match_history_game_type;
DROP INDEX IF EXISTS idx_match_history_created_at;
DROP INDEX IF EXISTS idx_match_history_winner_id;
DROP INDEX IF EXISTS idx_match_history_players;
DROP TABLE IF EXISTS match_history;
