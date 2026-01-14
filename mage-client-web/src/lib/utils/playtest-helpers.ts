/**
 * Playtest State Management Helpers
 *
 * Shared utilities for managing playtest game state.
 * Reduces duplication in playtest-game.ts store operations.
 */

import type { CardView } from '$lib/generated/mage/v1/models';
import type { PlaytestPlayer, PlaytestGameState } from '$lib/stores/playtest-game';
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
 * Remove a card from its source zone in the game state
 */
export function removeCardFromZone(
	state: PlaytestGameState,
	cardId: string,
	sourceZone: string
): PlaytestGameState {
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
	card: CardView,
	targetZone: string,
	controllerId: string
): PlaytestGameState {
	const newState = { ...state };
	const upperTargetZone = targetZone.toUpperCase();
	const normalizedZone = normalizeZoneName(upperTargetZone);

	// Prepare card for new zone
	card.faceDown = false;

	switch (normalizedZone) {
		case 'BATTLEFIELD':
			card.zone = ZoneId.BATTLEFIELD;
			card.controllerId = controllerId;
			newState.battlefield = [...newState.battlefield, card];
			break;

		case 'GRAVEYARD':
			card.zone = ZoneId.GRAVEYARD;
			const graveyardOwner = card.ownerId || controllerId;
			newState.players = newState.players.map((p) =>
				p.playerId === graveyardOwner ? { ...p, graveyard: [...p.graveyard, card] } : p
			);
			break;

		case 'EXILE':
			card.zone = ZoneId.EXILE;
			newState.exile = [...newState.exile, card];
			break;

		case 'HAND':
			card.zone = ZoneId.HAND;
			newState.players = newState.players.map((p) =>
				p.playerId === controllerId
					? { ...p, hand: [...p.hand, card], handCount: p.hand.length + 1 }
					: p
			);
			break;

		case 'COMMAND':
			card.zone = ZoneId.COMMAND;
			card.faceDown = false;
			newState.command = [...newState.command, card];
			break;

		case 'LIBRARY':
			card.zone = ZoneId.LIBRARY;
			card.faceDown = true;
			const libraryOwner = card.ownerId || controllerId;
			const addToTop = isLibraryTop(upperTargetZone);
			const addToBottom = isLibraryBottom(upperTargetZone);

			newState.players = newState.players.map((p) => {
				if (p.playerId !== libraryOwner) return p;

				const newLibrary = addToBottom ? [...p.library, card] : [card, ...p.library];

				return {
					...p,
					library: newLibrary,
					libraryCount: p.library.length + 1
				};
			});
			break;

		case 'STACK':
			card.zone = ZoneId.STACK;
			newState.stack = [...newState.stack, card];
			break;
	}

	return newState;
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
