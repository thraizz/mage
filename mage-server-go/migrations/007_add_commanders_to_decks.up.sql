-- Add commanders column to decks table for Commander format support
ALTER TABLE decks ADD COLUMN IF NOT EXISTS commanders JSONB DEFAULT '[]'::jsonb;

-- Add comment explaining the column
COMMENT ON COLUMN decks.commanders IS 'Array of commander card names for Commander format decks';
