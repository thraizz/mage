-- Drop decks table
DROP TRIGGER IF EXISTS trigger_update_decks_updated_at ON decks;
DROP FUNCTION IF EXISTS update_decks_updated_at();
DROP INDEX IF EXISTS idx_decks_user_name;
DROP INDEX IF EXISTS idx_decks_format;
DROP INDEX IF EXISTS idx_decks_user_id;
DROP TABLE IF EXISTS decks;
