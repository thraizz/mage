import type { Deck, DeckUploadRequest } from '$lib/types/deck';
import { getMageClient } from '$lib/grpc/client';
import type {
	DeckListRequest,
	DeckListResponse,
	DeckGetRequest,
	DeckGetResponse,
	DeckDeleteRequest,
	DeckDeleteResponse,
	DeckSaveRequest,
	DeckSaveResponse,
	DeckInfo,
	DeckCardLists
} from '$lib/generated/mage/v1/table';

/**
 * Convert DeckInfo from proto to our Deck type
 */
function convertDeckInfoToDeck(deckInfo: DeckInfo): Deck {
	return {
		id: deckInfo.id.toString(),
		name: deckInfo.name,
		format: deckInfo.format,
		cardCount: deckInfo.mainDeckCount + deckInfo.sideboardCount,
		createdAt: deckInfo.createdAt * 1000, // Convert seconds to milliseconds
		updatedAt: deckInfo.updatedAt * 1000, // Convert seconds to milliseconds
		isValid: true, // Assume valid if returned from server
		mainDeck: [], // Summary view doesn't include card details
		sideboard: []
	};
}

/**
 * Convert DeckCardLists from proto to our DeckCard array format
 */
function convertCardListsToDeckCards(cardLists: DeckCardLists) {
	const convertCardList = (cards: any[]) => {
		return cards.map((card) => ({
			cardName: card.name,
			quantity: card.quantity || 1,
			setCode: card.setCode,
			manaCost: card.manaCost,
			cardType: card.cardType,
			types: card.types || [],
			colors: card.colors || [],
			power: card.power,
			toughness: card.toughness
		}));
	};

	return {
		mainDeck: convertCardList(cardLists.mainDeck),
		sideboard: convertCardList(cardLists.sideboard)
	};
}

/**
 * Fetch user's decks, optionally filtered by format
 */
export async function fetchUserDecks(format?: string): Promise<Deck[]> {
	const client = getMageClient();
	const sessionId = await client.ensureSessionId();

	if (!sessionId) {
		throw new Error('No active session - please login first');
	}

	const request: DeckListRequest = {
		sessionId,
		format: format || ''
	};

	try {
		const response = await client.call<DeckListRequest, DeckListResponse>('DeckList', request);

		if (!response.success) {
			// Check if error indicates session expired
			const errorMsg = response.error || 'Failed to fetch decks';
			if (errorMsg.toLowerCase().includes('session') || errorMsg.toLowerCase().includes('expired')) {
				throw new Error('Session expired - please login again');
			}
			throw new Error(errorMsg);
		}

		return response.decks.map(convertDeckInfoToDeck);
	} catch (error) {
		// Handle network or other errors
		if (error instanceof Error) {
			throw error;
		}
		throw new Error('Failed to fetch decks - network error');
	}
}

/**
 * Get full deck details including card lists
 */
export async function getDeckDetails(deckId: string): Promise<Deck> {
	const client = getMageClient();
	const sessionId = await client.ensureSessionId();

	if (!sessionId) {
		throw new Error('No active session - please login first');
	}

	const request: DeckGetRequest = {
		sessionId,
		deckId: parseInt(deckId)
	};

	const response = await client.call<DeckGetRequest, DeckGetResponse>('DeckGet', request);

	if (!response.success) {
		throw new Error(response.error || 'Failed to fetch deck details');
	}

	if (!response.info || !response.deck) {
		throw new Error('Deck not found');
	}

	const baseDeck = convertDeckInfoToDeck(response.info);
	const cardData = convertCardListsToDeckCards(response.deck);

	return {
		...baseDeck,
		mainDeck: cardData.mainDeck,
		sideboard: cardData.sideboard
	};
}

/**
 * Upload a new deck
 */
export async function uploadDeck(request: DeckUploadRequest): Promise<Deck> {
	const client = getMageClient();
	const sessionId = await client.ensureSessionId();

	if (!sessionId) {
		throw new Error('No active session - please login first');
	}

	// Parse the deck list text format
	// Expected format: "4 Lightning Bolt\n20 Mountain\n\nSideboard:\n2 Dragon's Claw"
	const lines = request.deckList.split('\n').map((line) => line.trim());
	const mainDeckCards = new Map<string, number>();
	const sideboardCards = new Map<string, number>();
	let inSideboard = false;

	for (const line of lines) {
		if (!line || line.startsWith('#') || line.startsWith('//')) {
			continue; // Skip empty lines and comments
		}

		if (line.toLowerCase().includes('sideboard')) {
			inSideboard = true;
			continue;
		}

		// Parse line: "4 Lightning Bolt" or "Lightning Bolt" or "4x Lightning Bolt"
		const match = line.match(/^(\d+)x?\s+(.+)$/i);
		if (match) {
			const quantity = parseInt(match[1]);
			const cardName = match[2].trim();
			if (inSideboard) {
				sideboardCards.set(cardName, (sideboardCards.get(cardName) || 0) + quantity);
			} else {
				mainDeckCards.set(cardName, (mainDeckCards.get(cardName) || 0) + quantity);
			}
		} else if (line) {
			// Single card without quantity prefix
			if (inSideboard) {
				sideboardCards.set(line, (sideboardCards.get(line) || 0) + 1);
			} else {
				mainDeckCards.set(line, (mainDeckCards.get(line) || 0) + 1);
			}
		}
	}

	// Convert to DeckCard format (server will populate metadata)
	const mainDeck = Array.from(mainDeckCards.entries()).map(([name, quantity]) => ({
		name,
		quantity,
		manaCost: '',
		cardType: '',
		types: [],
		colors: [],
		power: '',
		toughness: ''
	}));

	const sideboard = Array.from(sideboardCards.entries()).map(([name, quantity]) => ({
		name,
		quantity,
		manaCost: '',
		cardType: '',
		types: [],
		colors: [],
		power: '',
		toughness: ''
	}));

	const deckCardLists: DeckCardLists = {
		mainDeck,
		sideboard
	};

	const saveRequest: DeckSaveRequest = {
		sessionId,
		deckName: request.name,
		deck: deckCardLists,
		format: request.format,
		description: ''
	};

	const response = await client.call<DeckSaveRequest, DeckSaveResponse>('DeckSave', saveRequest);

	if (!response.success) {
		throw new Error(response.error || 'Failed to save deck');
	}

	if (!response.deckId) {
		throw new Error('Deck saved but no deck ID returned from server');
	}

	// Calculate total card count from quantities
	const mainDeckCount = mainDeck.reduce((sum, card) => sum + card.quantity, 0);
	const sideboardCount = sideboard.reduce((sum, card) => sum + card.quantity, 0);

	// Return a basic deck object - caller should refetch to get full details
	return {
		id: response.deckId.toString(),
		name: request.name,
		format: request.format,
		cardCount: mainDeckCount + sideboardCount,
		createdAt: Date.now(),
		updatedAt: Date.now(),
		isValid: true,
		mainDeck: [],
		sideboard: []
	};
}

/**
 * Delete a deck
 */
export async function deleteDeck(deckId: string): Promise<void> {
	const client = getMageClient();
	const sessionId = await client.ensureSessionId();

	if (!sessionId) {
		throw new Error('No active session - please login first');
	}

	const request: DeckDeleteRequest = {
		sessionId,
		deckId: parseInt(deckId)
	};

	const response = await client.call<DeckDeleteRequest, DeckDeleteResponse>(
		'DeckDelete',
		request
	);

	if (!response.success) {
		throw new Error(response.error || 'Failed to delete deck');
	}
}
