import type { Deck, DeckUploadRequest } from '$lib/types/deck';

/**
 * Mock deck data for development
 */
const MOCK_DECKS: Deck[] = [
	{
		id: 'd1',
		name: 'Mono-Red Aggro',
		format: 'Standard',
		cardCount: 60,
		createdAt: Date.now() - 86400000 * 7,
		updatedAt: Date.now() - 86400000 * 2,
		isValid: true,
		mainDeck: [
			{ cardName: 'Mountain', quantity: 20 },
			{ cardName: 'Lightning Bolt', quantity: 4 },
			{ cardName: 'Goblin Guide', quantity: 4 }
		],
		sideboard: []
	},
	{
		id: 'd2',
		name: 'Blue Control',
		format: 'Standard',
		cardCount: 60,
		createdAt: Date.now() - 86400000 * 5,
		updatedAt: Date.now() - 86400000,
		isValid: true,
		mainDeck: [
			{ cardName: 'Island', quantity: 24 },
			{ cardName: 'Counterspell', quantity: 4 }
		],
		sideboard: []
	},
	{
		id: 'd3',
		name: 'Atraxa Commander',
		format: 'Commander',
		cardCount: 100,
		createdAt: Date.now() - 86400000 * 14,
		updatedAt: Date.now() - 86400000 * 3,
		isValid: true,
		mainDeck: [
			{ cardName: 'Atraxa, Praetors\' Voice', quantity: 1 },
			{ cardName: 'Sol Ring', quantity: 1 },
			{ cardName: 'Command Tower', quantity: 1 }
		],
		sideboard: []
	},
	{
		id: 'd4',
		name: 'Izzet Spells',
		format: 'Modern',
		cardCount: 60,
		createdAt: Date.now() - 86400000 * 10,
		updatedAt: Date.now() - 86400000 * 5,
		isValid: true,
		mainDeck: [
			{ cardName: 'Steam Vents', quantity: 4 },
			{ cardName: 'Opt', quantity: 4 }
		],
		sideboard: []
	}
];

/**
 * Fetch user's decks, optionally filtered by format
 */
export async function fetchUserDecks(format?: string): Promise<Deck[]> {
	// Simulate network delay
	await new Promise((resolve) => setTimeout(resolve, 400));

	// In production, this would be:
	// const response = await grpcCall(deckService.listUserDecks, { format }, 'DeckService.listUserDecks');
	// return response.decks;

	if (format) {
		return MOCK_DECKS.filter((deck) => deck.format === format);
	}

	return MOCK_DECKS;
}

/**
 * Upload a new deck
 */
export async function uploadDeck(request: DeckUploadRequest): Promise<Deck> {
	// Simulate network delay
	await new Promise((resolve) => setTimeout(resolve, 800));

	// In production, this would be:
	// const response = await grpcCall(deckService.uploadDeck, request, 'DeckService.uploadDeck');
	// return response.deck;

	const newDeck: Deck = {
		id: `deck-${Date.now()}`,
		name: request.name,
		format: request.format,
		cardCount: 60, // Would be parsed from deckList
		createdAt: Date.now(),
		updatedAt: Date.now(),
		isValid: true,
		mainDeck: [],
		sideboard: []
	};

	return newDeck;
}

/**
 * Delete a deck
 */
export async function deleteDeck(deckId: string): Promise<void> {
	// Simulate network delay
	await new Promise((resolve) => setTimeout(resolve, 300));

	// In production, this would be:
	// await grpcCall(deckService.deleteDeck, { deckId }, 'DeckService.deleteDeck');

	console.log(`Deleting deck ${deckId}`);
}
