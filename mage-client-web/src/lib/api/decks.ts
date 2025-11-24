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
	const parseCardList = (cards: string[]) => {
		return cards.map((card) => {
			// Parse card names - format is typically just card name
			// or "4x Card Name" or "Card Name (SET)"
			const match = card.match(/^(?:(\d+)x?\s+)?(.+?)(?:\s+\(([A-Z0-9]+)\))?$/i);
			if (match) {
				const quantity = match[1] ? parseInt(match[1]) : 1;
				const cardName = match[2].trim();
				const setCode = match[3];
				return { cardName, quantity, setCode };
			}
			return { cardName: card, quantity: 1 };
		});
	};

	return {
		mainDeck: parseCardList(cardLists.mainDeck),
		sideboard: parseCardList(cardLists.sideboard)
	};
}

/**
 * Fetch user's decks, optionally filtered by format
 */
export async function fetchUserDecks(format?: string): Promise<Deck[]> {
	const client = getMageClient();
	const sessionId = client.getSessionId();

	if (!sessionId) {
		throw new Error('No active session - please login first');
	}

	const request: DeckListRequest = {
		sessionId,
		format: format || ''
	};

	const response = await client.call<DeckListRequest, DeckListResponse>('DeckList', request);

	if (!response.success) {
		throw new Error(response.error || 'Failed to fetch decks');
	}

	return response.decks.map(convertDeckInfoToDeck);
}

/**
 * Get full deck details including card lists
 */
export async function getDeckDetails(deckId: string): Promise<Deck> {
	const client = getMageClient();
	const sessionId = client.getSessionId();

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
	const sessionId = client.getSessionId();

	if (!sessionId) {
		throw new Error('No active session - please login first');
	}

	// Parse the deck list text format
	// Expected format: "4 Lightning Bolt\n20 Mountain\n\nSideboard:\n2 Dragon's Claw"
	const lines = request.deckList.split('\n').map((line) => line.trim());
	const mainDeck: string[] = [];
	const sideboard: string[] = [];
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
			for (let i = 0; i < quantity; i++) {
				if (inSideboard) {
					sideboard.push(cardName);
				} else {
					mainDeck.push(cardName);
				}
			}
		} else if (line) {
			// Single card without quantity prefix
			if (inSideboard) {
				sideboard.push(line);
			} else {
				mainDeck.push(line);
			}
		}
	}

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

	// Return a basic deck object - caller should refetch to get full details
	return {
		id: response.deckId.toString(),
		name: request.name,
		format: request.format,
		cardCount: mainDeck.length + sideboard.length,
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
	const sessionId = client.getSessionId();

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
