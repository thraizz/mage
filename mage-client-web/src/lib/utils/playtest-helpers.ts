/**
 * Playtest State Management Helpers
 *
 * Shared utilities for managing playtest game state.
 * Reduces duplication in playtest-game.ts store operations.
 */

import type { AbilityView, CardAction, CardView, CounterView } from '$lib/generated/mage/v1/models';
import type { PlaytestPlayer, PlaytestGameState, PlaytestLogEntry } from '$lib/stores/playtest-game';
import { ZoneId, isLibraryTop, isLibraryBottom, normalizeZoneName } from './zones';

/**
 * Update a specific player in the players array
 */
export function updatePlayer(
	players: PlaytestPlayer[],
	playerId: string,
	updater: (player: PlaytestPlayer) => Partial<PlaytestPlayer>
): PlaytestPlayer[] {
	return players.map((p) => (p.playerId === playerId ? { ...p, ...updater(p) } : p));
}

/**
 * Find a card in the game state and return it with its source zone
 */
export function findCardInState(
	state: PlaytestGameState,
	cardId: string
): { card: CardView; sourceZone: string } | null {
	// Check battlefield
	const bfCard = state.battlefield.find((c) => c.id === cardId);
	if (bfCard) return { card: bfCard, sourceZone: 'battlefield' };

	// Check exile
	const exileCard = state.exile.find((c) => c.id === cardId);
	if (exileCard) return { card: exileCard, sourceZone: 'exile' };

	// Check command
	const cmdCard = state.command.find((c) => c.id === cardId);
	if (cmdCard) return { card: cmdCard, sourceZone: 'command' };

	// Check stack
	const stackCard = state.stack.find((c) => c.id === cardId);
	if (stackCard) return { card: stackCard, sourceZone: 'stack' };

	// Check player zones
	for (const player of state.players) {
		const handCard = player.hand.find((c) => c.id === cardId);
		if (handCard) return { card: handCard, sourceZone: `hand:${player.playerId}` };

		const libraryCard = player.library.find((c) => c.id === cardId);
		if (libraryCard) return { card: libraryCard, sourceZone: `library:${player.playerId}` };

		const graveyardCard = player.graveyard.find((c) => c.id === cardId);
		if (graveyardCard) return { card: graveyardCard, sourceZone: `graveyard:${player.playerId}` };
	}

	return null;
}

/**
 * Check if a card is a token (tokens have IDs starting with "token-")
 */
export function isToken(cardId: string): boolean {
	return cardId.startsWith('token-');
}

/**
 * Update a card in its zone by creating a new card object
 */
export function updateCardInZone(
	state: PlaytestGameState,
	cardId: string,
	sourceZone: string,
	updater: (card: CardView) => CardView
): PlaytestGameState {
	const newState = { ...state };

	if (sourceZone === 'battlefield') {
		newState.battlefield = state.battlefield.map((c) => (c.id === cardId ? updater(c) : c));
	} else if (sourceZone === 'exile') {
		newState.exile = state.exile.map((c) => (c.id === cardId ? updater(c) : c));
	} else if (sourceZone === 'command') {
		newState.command = state.command.map((c) => (c.id === cardId ? updater(c) : c));
	} else if (sourceZone === 'stack') {
		newState.stack = state.stack.map((c) => (c.id === cardId ? updater(c) : c));
	} else if (sourceZone.startsWith('hand:')) {
		const playerId = sourceZone.split(':')[1];
		newState.players = state.players.map((p) =>
			p.playerId === playerId ? { ...p, hand: p.hand.map((c) => (c.id === cardId ? updater(c) : c)) } : p
		);
	} else if (sourceZone.startsWith('library:')) {
		const playerId = sourceZone.split(':')[1];
		newState.players = state.players.map((p) =>
			p.playerId === playerId
				? { ...p, library: p.library.map((c) => (c.id === cardId ? updater(c) : c)) }
				: p
		);
	} else if (sourceZone.startsWith('graveyard:')) {
		const playerId = sourceZone.split(':')[1];
		newState.players = state.players.map((p) =>
			p.playerId === playerId
				? { ...p, graveyard: p.graveyard.map((c) => (c.id === cardId ? updater(c) : c)) }
				: p
		);
	}

	return newState;
}

/**
 * Remove a card from its source zone in the game state
 */
export function removeCardFromZone(
	state: PlaytestGameState,
	cardId: string,
	sourceZone: string | undefined
): PlaytestGameState {
	// Handle undefined sourceZone gracefully
	if (!sourceZone) {
		console.warn('[playtest-helpers] removeCardFromZone called with undefined sourceZone for card:', cardId);
		return state;
	}

	const newState = { ...state };

	if (sourceZone === 'battlefield') {
		newState.battlefield = state.battlefield.filter((c) => c.id !== cardId);
	} else if (sourceZone === 'exile') {
		newState.exile = state.exile.filter((c) => c.id !== cardId);
	} else if (sourceZone === 'command') {
		newState.command = state.command.filter((c) => c.id !== cardId);
	} else if (sourceZone === 'stack') {
		newState.stack = state.stack.filter((c) => c.id !== cardId);
	} else if (sourceZone.startsWith('hand:')) {
		const playerId = sourceZone.split(':')[1];
		newState.players = state.players.map((p) =>
			p.playerId === playerId
				? { ...p, hand: p.hand.filter((c) => c.id !== cardId), handCount: p.hand.length - 1 }
				: p
		);
	} else if (sourceZone.startsWith('library:')) {
		const playerId = sourceZone.split(':')[1];
		newState.players = state.players.map((p) =>
			p.playerId === playerId
				? {
						...p,
						library: p.library.filter((c) => c.id !== cardId),
						libraryCount: p.library.length - 1
					}
				: p
		);
	} else if (sourceZone.startsWith('graveyard:')) {
		const playerId = sourceZone.split(':')[1];
		newState.players = state.players.map((p) =>
			p.playerId === playerId ? { ...p, graveyard: p.graveyard.filter((c) => c.id !== cardId) } : p
		);
	}

	return newState;
}

/**
 * Add a card to a target zone in the game state
 */
export function addCardToZone(
	state: PlaytestGameState,
	card: CardView | undefined,
	targetZone: string,
	controllerId: string
): PlaytestGameState {
	// Validate card exists
	if (!card) {
		console.warn('[playtest-helpers] addCardToZone called with undefined card');
		return state;
	}

	const newState = { ...state };
	const upperTargetZone = targetZone.toUpperCase();
	const normalizedZone = normalizeZoneName(upperTargetZone);

	// Create a copy of the card to avoid mutating the original
	const cardCopy = { ...card };

	// Prepare card for new zone
	cardCopy.faceDown = false;

	switch (normalizedZone) {
		case 'BATTLEFIELD':
			cardCopy.zone = ZoneId.BATTLEFIELD;
			cardCopy.controllerId = controllerId;
			newState.battlefield = [...newState.battlefield, cardCopy];
			break;

		case 'GRAVEYARD':
			cardCopy.zone = ZoneId.GRAVEYARD;
			const graveyardOwner = cardCopy.ownerId || controllerId;
			newState.players = newState.players.map((p) =>
				p.playerId === graveyardOwner ? { ...p, graveyard: [...p.graveyard, cardCopy] } : p
			);
			break;

		case 'EXILE':
			cardCopy.zone = ZoneId.EXILE;
			newState.exile = [...newState.exile, cardCopy];
			break;

		case 'HAND':
			cardCopy.zone = ZoneId.HAND;
			newState.players = newState.players.map((p) =>
				p.playerId === controllerId
					? { ...p, hand: [...p.hand, cardCopy], handCount: p.hand.length + 1 }
					: p
			);
			break;

		case 'COMMAND':
			cardCopy.zone = ZoneId.COMMAND;
			cardCopy.faceDown = false;
			newState.command = [...newState.command, cardCopy];
			break;

		case 'LIBRARY':
			newState.players = addCardToLibrary(cardCopy, controllerId, upperTargetZone, newState);
			break;

		case 'STACK':
			cardCopy.zone = ZoneId.STACK;
			newState.stack = [...newState.stack, cardCopy];
			break;
	}

	return newState;
}

function addCardToLibrary(cardCopy: CardView, controllerId: string, upperTargetZone: string, newState: PlaytestGameState) {
	cardCopy.zone = ZoneId.LIBRARY;
	cardCopy.faceDown = true;
	const libraryOwner = cardCopy.ownerId || controllerId;
	const addToTop = isLibraryTop(upperTargetZone);
	const addToBottom = isLibraryBottom(upperTargetZone);

	return newState.players.map((p) => {
		if (p.playerId !== libraryOwner) return p;

		const newLibrary = addToBottom ? [...p.library, cardCopy] : [cardCopy, ...p.library];

		return {
			...p,
			library: newLibrary,
			libraryCount: p.library.length + 1
		};
	});
}

/**
 * Shuffle an array using Fisher-Yates algorithm
 */
export function shuffleArray<T>(array: T[]): T[] {
	const shuffled = [...array];
	for (let i = shuffled.length - 1; i > 0; i--) {
		const j = Math.floor(Math.random() * (i + 1));
		[shuffled[i], shuffled[j]] = [shuffled[j], shuffled[i]];
	}
	return shuffled;
}

/**
 * Update card zone and visibility properties
 */
export function updateCardZone(card: CardView, zoneId: ZoneId, faceDown: boolean = false): void {
	card.zone = zoneId;
	card.faceDown = faceDown;
}

/**
 * Find the next player in turn order
 */
export function getNextPlayer(
	players: PlaytestPlayer[],
	currentPlayerId: string
): PlaytestPlayer | null {
	if (players.length === 0) return null;

	const currentIndex = players.findIndex((p) => p.playerId === currentPlayerId);
	const nextIndex = (currentIndex + 1) % players.length;
	return players[nextIndex];
}
