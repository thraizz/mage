/**
 * Deck type definitions for deck management
 */

export interface Deck {
	id: string;
	name: string;
	format: string;
	cardCount: number;
	createdAt: number;
	updatedAt: number;
	lastModified?: string; // Human-readable format
	isValid: boolean;
	mainDeck: DeckCard[];
	mainDeckCount?: number; // Server-provided count for summary views
	sideboard: DeckCard[];
	commanders: DeckCard[]; // For Commander format
}

export interface DeckCard {
	cardName: string;
	quantity: number;
	setCode?: string;
	manaCost?: string;
	cardType?: string;
	types?: string[];
	colors?: string[];
	power?: string;
	toughness?: string;
}

export interface DeckListState {
	decks: Deck[];
	isLoading: boolean;
	error: string | null;
}

export interface DeckUploadRequest {
	name: string;
	format: string;
	deckList: string; // Text format deck list
}
