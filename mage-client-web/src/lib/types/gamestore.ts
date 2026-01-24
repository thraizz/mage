/**
 * Shared Game Store Types
 *
 * Common types and interfaces used across playtest and multiplayer game stores.
 */

import type { CardView, ManaPoolView } from '$lib/generated/mage/v1/models';

/**
 * Player state for game stores
 */
export interface Player {
	playerId: string;
	name: string;
	life: number;
	poison: number;
	energy: number;
	libraryCount: number;
	handCount: number;
	hand: CardView[];
	library: CardView[];
	graveyard: CardView[];
	manaPool: ManaPoolView;
	keptHand: boolean;
	mulliganCount: number;
	revealedTopCard: boolean; // When true, top card of library is permanently visible
}

/**
 * Game log entry
 */
export type PlaytestLogEntry = {
	id: string;
	at: number; // unix ms
	turn: number;
	activePlayerId: string;
	controlSeat: string; // activeControlSeat at time of event
	kind: string; // "draw" | "move" | "life" | ...
	message: string;
};

/**
 * Scry session for tracking ongoing scry operations
 */
export type ScrySession = {
	sessionId: string;
	playerId: string;
	cards: CardView[];
};

/**
 * Base game state shared by both playtest and multiplayer stores
 */
export interface BaseGameState {
	gameId: string;
	activeControlSeat: string; // Which player perspective you're controlling
	players: Player[];
	battlefield: CardView[];
	exile: CardView[];
	stack: CardView[];
	command: CardView[];
	turn: number;
	activePlayerId: string;
	isInitialized: boolean;
	log: PlaytestLogEntry[];
	mulliganType: 'london';
	freeMulligans: number;
}

/**
 * Common interface for game store methods
 * Implemented by both playtestGameStore and multiplayerGameStore
 */
export interface GameStore {
	// Svelte store contract
	subscribe: (
		run: (value: BaseGameState) => void,
		invalidate?: (value?: BaseGameState) => void
	) => () => void;

	// Initialization
	initialize: (...args: any[]) => void | Promise<void>;
	reset: () => void;

	// Player actions
	drawCards: (playerId: string, count: number) => void;
	playCard: (cardId: string, tapped?: boolean) => void;
	moveCardToZone: (cardId: string, targetZone: string) => void;
	tapCard: (cardId: string, tapped: boolean) => void;
	untapAll: (playerId: string) => void;
	flipCard: (cardId: string, faceDown: boolean) => void;
	modifyLife: (playerId: string, delta: number) => void;
	setPlayerCounter: (playerId: string, counterType: string, value: number) => void;
	shuffleLibrary: (playerId: string) => void;

	// Stack operations
	addToStack: (cardId: string) => void;
	removeFromStack: (itemId: string) => void;

	// Token operations
	createToken: (
		name: string,
		types: string,
		power: string,
		toughness: string,
		color: string,
		abilities?: string[]
	) => void;

	// Counter operations
	addCounter: (cardId: string, counterName: string, amount?: number) => void;
	removeCounter: (cardId: string, counterName: string, amount?: number) => void;
	setCounter: (cardId: string, counterName: string, amount: number) => void;

	// Library operations
	millCards: (playerId: string, count: number) => void;
	revealTopCards: (playerId: string, count: number) => CardView[];
	scryCards: (playerId: string, count: number) => ScrySession | null;
	applyScryDecision: (
		playerId: string,
		scryCount: number,
		keepOnTop: CardView[],
		putToBottom: CardView[]
	) => void;
	setRevealedTop: (playerId: string, revealed: boolean) => void;

	// Turn management
	nextTurn: () => void;

	// Mulligan
	mulligan: (playerId: string) => void | Promise<void>;
	keepHand: (playerId: string) => void | Promise<void>;

	// View control
	switchControlSeat: (playerId: string) => void;

	// Command zone
	setCommand: (cards: CardView[]) => void;
}
