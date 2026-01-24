/**
 * Playtest Initializer
 *
 * Handles initialization of playtest games from deck lists
 */

import type { Deck } from '$lib/types/deck';
import type { CardView } from '$lib/generated/mage/v1/models';
import type { PlaytestPlayer } from '$lib/types/gamestore';
import { getDeckDetails } from '$lib/api/decks';
import { ZoneId } from '$lib/utils/zones';
import { shuffleArray } from '$lib/utils/playtest-helpers';

export type PlaytestInitResult = {
	players: PlaytestPlayer[];
	command: CardView[];
};

/**
 * Initialize a playtest session from deck IDs
 */
export async function initializePlaytest(deckIds: string[]): Promise<PlaytestInitResult> {
	if (deckIds.length < 2 || deckIds.length > 4) {
		throw new Error('Playtest requires 2-4 decks');
	}

	console.log('[PlaytestInit] Fetching deck details for:', deckIds);

	// Fetch all deck details in parallel
	const deckPromises = deckIds.map((deckId) => getDeckDetails(deckId));
	const decks = await Promise.all(deckPromises);

	console.log(
		'[PlaytestInit] Decks loaded:',
		decks.map((d) => d.name)
	);
	console.log(
		'[PlaytestInit] Deck meta sample:',
		decks.map((d) => ({
			id: d.id,
			name: d.name,
			format: d.format,
			mainDeckCount: d.mainDeck.length,
			firstCard: d.mainDeck[0]
				? {
						cardName: d.mainDeck[0].cardName,
						cardType: d.mainDeck[0].cardType,
						manaCost: d.mainDeck[0].manaCost
					}
				: null
		}))
	);

	// Create players from decks
	const perDeck = decks.map((deck, index) => createPlayerFromDeck(deck, index + 1));
	const players: PlaytestPlayer[] = perDeck.map((x) => x.player);
	const command: CardView[] = perDeck.flatMap((x) => x.commanders);

	// Shuffle each player's library
	players.forEach((player) => {
		player.library = shuffleArray(player.library);
	});

	// Draw opening hands (7 cards each)
	players.forEach((player) => {
		const hand = player.library.splice(0, 7);
		hand.forEach((card) => {
			card.zone = ZoneId.HAND;
			card.faceDown = false;
		});
		player.hand = hand;
		player.handCount = hand.length;
		player.libraryCount = player.library.length;
	});

	console.log('[PlaytestInit] Players initialized with opening hands');

	return { players, command };
}

/**
 * Create a player from a deck
 */
function createPlayerFromDeck(
	deck: Deck,
	playerNumber: number
): { player: PlaytestPlayer; commanders: CardView[] } {
	const playerId = `player${playerNumber}`;

	// Create card objects from deck list
	const library: CardView[] = [];
	let cardIndex = 0;

	// Add main deck cards
	for (const deckCard of deck.mainDeck) {
		for (let i = 0; i < deckCard.quantity; i++) {
			library.push(createCardView(playerId, deckCard.cardName, deckCard, cardIndex++));
		}
	}

	// Commanders go to command zone (handled separately in game)
	const commanders: CardView[] = [];
	if (deck.commanders && deck.commanders.length > 0) {
		for (const deckCard of deck.commanders) {
			for (let i = 0; i < deckCard.quantity; i++) {
				const commander = createCardView(playerId, deckCard.cardName, deckCard, cardIndex++);
				commander.zone = ZoneId.COMMAND;
				commander.faceDown = false;
				commanders.push(commander);
			}
		}
	}

	// Determine starting life based on format
	const startingLife = getStartingLife(deck.format);

	const player: PlaytestPlayer = {
		playerId,
		name: `${deck.name} (P${playerNumber})`,
		life: startingLife,
		poison: 0,
		energy: 0,
		libraryCount: library.length,
		handCount: 0,
		hand: [],
		library,
		graveyard: [],
		manaPool: {
			white: 0,
			blue: 0,
			black: 0,
			red: 0,
			green: 0,
			colorless: 0
		},
		keptHand: false,
		mulliganCount: 0,
		revealedTopCard: false
	};

	return { player, commanders };
}

/**
 * Create a CardView from deck card data
 */
function normalizeCardName(rawName: string): string {
	if (!rawName) return '';
	// Some upstream sources encode extra metadata after "@@@"
	// Example: "Swamp@@@<meta>" → "Swamp"
	return rawName.split('@@@')[0].trim();
}

function createCardView(
	ownerId: string,
	cardName: string,
	deckCard: {
		manaCost?: string;
		cardType?: string;
		power?: string;
		toughness?: string;
		colors?: string[];
	},
	index: number
): CardView {
	const cleanedName = normalizeCardName(cardName);

	return {
		id: `${ownerId}-card-${index}`,
		name: cleanedName,
		displayName: cleanedName,
		manaCost: deckCard.manaCost || '',
		type: deckCard.cardType || '',
		subTypes: '',
		superTypes: '',
		color: (deckCard.colors || []).join(''),
		power: deckCard.power || '',
		toughness: deckCard.toughness || '',
		loyalty: '',
		cardNumber: 0,
		expansionSetCode: '',
		rarity: '',
		rulesText: '',
		abilities: [],
		zone: ZoneId.LIBRARY,
		ownerId,
		controllerId: ownerId,
		tapped: false,
		flipped: false,
		transformed: false,
		faceDown: true,
		counters: [],
		attachedTo: [],
		summoningSickness: false,
		availableActions: []
	};
}

/**
 * Get starting life total based on format
 */
function getStartingLife(format: string): number {
	const normalizedFormat = format.toLowerCase();

	if (normalizedFormat.includes('commander') || normalizedFormat.includes('edh')) {
		return 40;
	}

	// Standard, Modern, Legacy, Vintage, etc.
	return 20;
}

/**
 * Validate deck IDs from URL params
 */
export function validateDeckIds(searchParams: URLSearchParams): string[] {
	const deckIds: string[] = [];

	for (let i = 1; i <= 4; i++) {
		const deckId = searchParams.get(`d${i}`);
		if (deckId) {
			deckIds.push(deckId);
		}
	}

	if (deckIds.length < 2) {
		throw new Error('At least 2 decks are required for playtest mode');
	}

	return deckIds;
}
