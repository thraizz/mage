-- Remove commanders column from decks table
ALTER TABLE decks DROP COLUMN IF EXISTS commanders;
