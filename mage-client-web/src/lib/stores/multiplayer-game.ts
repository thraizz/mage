/**
 * Multiplayer Game Store
 *
 * Server-synchronized game state management for multiplayer mode.
 * Based on playtest-game.ts structure but sends operations to server
 * and receives state updates via WebSocket.
 */

import { writable, derived, get } from 'svelte/store';
import type { CardView, ManaPoolView } from '$lib/generated/mage/v1/models';
import { CallbackMethod } from '$lib/generated/mage/v1/websocket';
import type { GameUpdateData, GameInitData } from '$lib/generated/mage/v1/websocket';
import { websocketStore } from './websocket';
import * as directActions from '$lib/api/direct-actions';

/**
 * Player state for multiplayer game
 * From playtest-game.ts lines 25-40
 */
export interface PlaytestPlayer {
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
 * From playtest-game.ts lines 116-124
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
 * From playtest-game.ts lines 129-133
 */
export type ScrySession = {
	sessionId: string;
	playerId: string;
	cards: CardView[];
};

/**
 * Multiplayer game state
 * Based on playtest-game.ts lines 45-59
 * Added server sync fields: isConnected, pendingActions
 */
export interface MultiplayerGameState {
	gameId: string;
	activeControlSeat: string; // Which player perspective you're controlling
	players: PlaytestPlayer[];
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

	// NEW: Server sync fields
	isConnected: boolean;
	pendingActions: string[];
}

const initialState: MultiplayerGameState = {
	gameId: '',
	activeControlSeat: '',
	players: [],
	battlefield: [],
	exile: [],
	stack: [],
	command: [],
	turn: 1,
	activePlayerId: '',
	isInitialized: false,
	log: [],
	mulliganType: 'london',
	freeMulligans: 0,
	isConnected: false,
	pendingActions: []
};

/**
 * Create multiplayer game store
 * Based on playtest-game.ts lines 380-401 (store creation pattern)
 */
function createMultiplayerGameStore() {
	const { subscribe, set, update } = writable<MultiplayerGameState>(initialState);

	// Track unsubscribe functions for WebSocket handlers
	const unsubscribers: Array<() => void> = [];

	/**
	 * Subscribe to WebSocket game state updates
	 * Pattern from game.ts lines 140-200
	 */
	function subscribeToGameEvents() {
		// Unsubscribe from previous subscriptions
		unsubscribeFromEvents();

		// GAME_INIT - Initial game state
		unsubscribers.push(
			websocketStore.on(CallbackMethod.GAME_INIT, (data) => {
				const initData = data as GameInitData;
				console.log('[MultiplayerGame] GAME_INIT received:', initData);

				if (initData.game) {
					// Apply server state to store
					update((state) => ({
						...state,
						// Map server GameView to our state structure
						// TODO: Implement full mapping when server GameView is available
						isConnected: true,
						isInitialized: true,
						pendingActions: []
					}));
				}
			})
		);

		// GAME_UPDATE - State update
		unsubscribers.push(
			websocketStore.on(CallbackMethod.GAME_UPDATE, (data) => {
				const updateData = data as GameUpdateData;
				console.log('[MultiplayerGame] GAME_UPDATE received:', updateData);

				if (updateData.game) {
					// Apply server state to store
					update((state) => ({
						...state,
						// Map server GameView to our state structure
						// TODO: Implement full mapping when server GameView is available
						isConnected: true,
						pendingActions: []
					}));
				}
			})
		);
	}

	/**
	 * Unsubscribe from WebSocket events
	 */
	function unsubscribeFromEvents() {
		unsubscribers.forEach((unsubscribe) => unsubscribe());
		unsubscribers.length = 0;
	}

	/**
	 * Initialize game state
	 * From playtest-game.ts lines 418-467
	 * Modified: Initialization happens server-side, we just subscribe
	 */
	function initialize(gameId: string): void {
		update((state) => ({
			...state,
			gameId,
			isInitialized: false,
			isConnected: false
		}));

		// Subscribe to WebSocket updates for this game
		subscribeToGameEvents();

		console.log('[MultiplayerGame] Initialized for game:', gameId);
	}

	/**
	 * Draw cards for a player
	 * From playtest-game.ts lines 492-527
	 * Modified: Send to server instead of local mutation
	 */
	function drawCards(playerId: string, count: number): void {
		const state = get({ subscribe });
		directActions.drawCards(state.gameId, playerId, count);
		console.log('[MultiplayerGame] drawCards:', { playerId, count });
	}

	/**
	 * Play a card from hand to battlefield
	 * From playtest-game.ts lines 532-569
	 * Modified: Send to server instead of local mutation
	 */
	function playCard(cardId: string, tapped: boolean = false): void {
		const state = get({ subscribe });
		// Use moveCard to move from hand to battlefield
		directActions.moveCard(state.gameId, cardId, 'BATTLEFIELD');
		if (tapped) {
			directActions.tapUntap(state.gameId, cardId, true);
		}
		console.log('[MultiplayerGame] playCard:', { cardId, tapped });
	}

	/**
	 * Move a card to a different zone
	 * From playtest-game.ts lines 574-607
	 * Modified: Send to server instead of local mutation
	 */
	function moveCardToZone(cardId: string, targetZone: string): void {
		const state = get({ subscribe });
		directActions.moveCard(state.gameId, cardId, targetZone);
		console.log('[MultiplayerGame] moveCardToZone:', { cardId, targetZone });
	}

	/**
	 * Tap or untap a card
	 * From playtest-game.ts lines 612-630
	 * Modified: Send to server instead of local mutation
	 */
	function tapCard(cardId: string, tapped: boolean): void {
		const state = get({ subscribe });
		directActions.tapUntap(state.gameId, cardId, tapped);
		console.log('[MultiplayerGame] tapCard:', { cardId, tapped });
	}

	/**
	 * Untap all permanents controlled by a player
	 * From playtest-game.ts lines 634-646
	 * Modified: Send to server instead of local mutation
	 */
	function untapAll(playerId: string): void {
		const state = get({ subscribe });
		directActions.untapAll(state.gameId);
		console.log('[MultiplayerGame] untapAll:', { playerId });
	}

	/**
	 * Flip a card face up/down
	 * From playtest-game.ts lines 650-663
	 * Modified: Send to server instead of local mutation
	 */
	function flipCard(cardId: string, faceDown: boolean): void {
		const state = get({ subscribe });
		directActions.flipCard(state.gameId, cardId, faceDown);
		console.log('[MultiplayerGame] flipCard:', { cardId, faceDown });
	}

	/**
	 * Transform a double-faced card
	 * New operation - not in playtest
	 */
	function transformCard(cardId: string): void {
		const state = get({ subscribe });
		directActions.transformCard(state.gameId, cardId);
		console.log('[MultiplayerGame] transformCard:', { cardId });
	}

	/**
	 * Modify player life
	 * From playtest-game.ts lines 668-679
	 * Modified: Send to server instead of local mutation
	 */
	function modifyLife(playerId: string, delta: number): void {
		const state = get({ subscribe });
		directActions.modifyPlayerLife(state.gameId, playerId, delta);
		console.log('[MultiplayerGame] modifyLife:', { playerId, delta });
	}

	/**
	 * Set player counter (poison, energy, etc.)
	 * From playtest-game.ts lines 684-700
	 * Modified: Send to server instead of local mutation
	 */
	function setPlayerCounter(playerId: string, counterType: string, value: number): void {
		const state = get({ subscribe });
		directActions.setPlayerCounter(state.gameId, playerId, counterType, value);
		console.log('[MultiplayerGame] setPlayerCounter:', { playerId, counterType, value });
	}

	/**
	 * Shuffle a player's library
	 * From playtest-game.ts lines 705-717
	 * Modified: Send to server instead of local mutation
	 */
	function shuffleLibrary(playerId: string): void {
		const state = get({ subscribe });
		directActions.shuffleLibrary(state.gameId, playerId);
		console.log('[MultiplayerGame] shuffleLibrary:', { playerId });
	}

	/**
	 * Add a card to the visual stack
	 * From playtest-game.ts lines 722-740
	 * Modified: Send to server instead of local mutation
	 */
	function addToStack(cardId: string): void {
		const state = get({ subscribe });
		directActions.addToStack(state.gameId, cardId);
		console.log('[MultiplayerGame] addToStack:', { cardId });
	}

	/**
	 * Remove an item from the stack
	 * From playtest-game.ts lines 745-754
	 * Modified: Send to server instead of local mutation
	 */
	function removeFromStack(itemId: string): void {
		const state = get({ subscribe });
		directActions.removeFromStack(state.gameId, itemId);
		console.log('[MultiplayerGame] removeFromStack:', { itemId });
	}

	/**
	 * Create a token on the battlefield
	 * From playtest-game.ts lines 759-805
	 * Modified: Send to server instead of local mutation
	 */
	function createToken(
		name: string,
		types: string,
		power: string,
		toughness: string,
		color: string,
		abilities: string[] = []
	): void {
		const state = get({ subscribe });
		directActions.createToken(state.gameId, name, types, power, toughness, color, abilities);
		console.log('[MultiplayerGame] createToken:', { name, types, power, toughness, color });
	}

	/**
	 * Destroy a token (remove it from the game)
	 * New operation - maps to directActions
	 */
	function destroyToken(cardId: string): void {
		const state = get({ subscribe });
		directActions.destroyToken(state.gameId, cardId);
		console.log('[MultiplayerGame] destroyToken:', { cardId });
	}

	/**
	 * Add counters to a card
	 * From playtest-game.ts lines 810-840
	 * Modified: Send to server instead of local mutation
	 */
	function addCounter(cardId: string, counterName: string, amount: number = 1): void {
		const state = get({ subscribe });
		directActions.modifyCardCounter(state.gameId, cardId, counterName, amount);
		console.log('[MultiplayerGame] addCounter:', { cardId, counterName, amount });
	}

	/**
	 * Remove counters from a card
	 * From playtest-game.ts lines 845-882
	 * Modified: Send to server instead of local mutation
	 */
	function removeCounter(cardId: string, counterName: string, amount: number = 1): void {
		const state = get({ subscribe });
		directActions.modifyCardCounter(state.gameId, cardId, counterName, -amount);
		console.log('[MultiplayerGame] removeCounter:', { cardId, counterName, amount });
	}

	/**
	 * Set a counter to a specific value
	 * From playtest-game.ts lines 887-918
	 * Modified: Send to server instead of local mutation
	 */
	function setCounter(cardId: string, counterName: string, amount: number): void {
		const state = get({ subscribe });
		directActions.setCardCounter(state.gameId, cardId, counterName, amount);
		console.log('[MultiplayerGame] setCounter:', { cardId, counterName, amount });
	}

	/**
	 * Mill cards (move top N cards from library to graveyard)
	 * From playtest-game.ts lines 923-957
	 * Modified: Send to server via moveCard (not directly implemented in direct-actions)
	 * Note: This may need server implementation for atomic mill operation
	 */
	function millCards(playerId: string, count: number): void {
		console.warn('[MultiplayerGame] millCards not yet implemented server-side:', {
			playerId,
			count
		});
		// TODO: Implement MILL command in direct-actions API
	}

	/**
	 * Reveal top N cards (for temporary view)
	 * From playtest-game.ts lines 962-978
	 * Modified: This is a read operation, not implemented yet
	 */
	function revealTopCards(playerId: string, count: number): CardView[] {
		console.warn('[MultiplayerGame] revealTopCards not yet implemented server-side:', {
			playerId,
			count
		});
		// TODO: Implement REVEAL_TOP command in direct-actions API
		return [];
	}

	/**
	 * Start a scry session
	 * From playtest-game.ts lines 983-1011
	 * Modified: This requires server support for scry state
	 */
	function scryCards(playerId: string, count: number): ScrySession | null {
		console.warn('[MultiplayerGame] scryCards not yet implemented server-side:', {
			playerId,
			count
		});
		// TODO: Implement SCRY command in direct-actions API
		return null;
	}

	/**
	 * Apply scry decision
	 * From playtest-game.ts lines 1016-1053
	 * Modified: This requires server support for scry state
	 */
	function applyScryDecision(
		playerId: string,
		scryCount: number,
		keepOnTop: CardView[],
		putToBottom: CardView[]
	): void {
		console.warn('[MultiplayerGame] applyScryDecision not yet implemented server-side:', {
			playerId,
			scryCount,
			keepOnTop: keepOnTop.length,
			putToBottom: putToBottom.length
		});
		// TODO: Implement SCRY_APPLY command in direct-actions API
	}

	/**
	 * Set revealed top card state
	 * From playtest-game.ts lines 1058-1071
	 * Modified: This requires server state tracking
	 */
	function setRevealedTop(playerId: string, revealed: boolean): void {
		console.warn('[MultiplayerGame] setRevealedTop not yet implemented server-side:', {
			playerId,
			revealed
		});
		// TODO: Implement REVEAL_TOP_PERMANENT command in direct-actions API
	}

	/**
	 * Next turn
	 * From playtest-game.ts lines 1076-1090
	 * Modified: Send to server instead of local mutation
	 */
	function nextTurn(): void {
		const state = get({ subscribe });
		directActions.nextTurn(state.gameId);
		console.log('[MultiplayerGame] nextTurn');
	}

	/**
	 * Clear combat state
	 * New operation - maps to directActions
	 */
	function clearCombat(): void {
		const state = get({ subscribe });
		directActions.clearCombat(state.gameId);
		console.log('[MultiplayerGame] clearCombat');
	}

	/**
	 * Search library
	 * New operation - maps to directActions
	 */
	function searchLibrary(
		destination: 'hand' | 'battlefield' | 'top' | 'graveyard' = 'hand',
		shuffle: boolean = true,
		message?: string
	): void {
		const state = get({ subscribe });
		directActions.searchLibrary(state.gameId, destination, shuffle, message);
		console.log('[MultiplayerGame] searchLibrary:', { destination, shuffle, message });
	}

	/**
	 * Select a card from library search
	 * New operation - maps to directActions
	 */
	function selectLibraryCard(cardId: string): void {
		const state = get({ subscribe });
		directActions.selectLibraryCard(state.gameId, cardId);
		console.log('[MultiplayerGame] selectLibraryCard:', { cardId });
	}

	/**
	 * Mulligan for a player
	 * From playtest-game.ts lines 1095-1146
	 * Modified: This requires server mulligan implementation
	 */
	function mulligan(playerId: string): void {
		console.warn('[MultiplayerGame] mulligan not yet implemented server-side:', { playerId });
		// TODO: Implement MULLIGAN command in direct-actions API
	}

	/**
	 * Keep hand (no mulligan)
	 * From playtest-game.ts lines 1151-1161
	 * Modified: This requires server mulligan state tracking
	 */
	function keepHand(playerId: string): void {
		console.warn('[MultiplayerGame] keepHand not yet implemented server-side:', { playerId });
		// TODO: Implement KEEP_HAND command in direct-actions API
	}

	/**
	 * Switch active control seat (which player you're controlling)
	 * From playtest-game.ts lines 482-487
	 * Modified: This is client-side only for view perspective
	 */
	function switchControlSeat(playerId: string): void {
		update((state) => ({
			...state,
			activeControlSeat: playerId
		}));
		console.log('[MultiplayerGame] switchControlSeat:', { playerId });
	}

	/**
	 * Set command zone cards
	 * From playtest-game.ts lines 472-477
	 * Modified: This requires server support
	 */
	function setCommand(cards: CardView[]): void {
		console.warn('[MultiplayerGame] setCommand not yet implemented server-side:', {
			cardCount: cards.length
		});
		// TODO: Implement SET_COMMAND command in direct-actions API
	}

	/**
	 * Reset to initial state
	 * From playtest-game.ts lines 1166-1170
	 */
	function reset(): void {
		unsubscribeFromEvents();
		set(initialState);
		console.log('[MultiplayerGame] reset');
	}

	/**
	 * Cleanup on destroy
	 */
	function cleanup(): void {
		unsubscribeFromEvents();
	}

	return {
		subscribe,
		initialize,
		drawCards,
		playCard,
		moveCardToZone,
		tapCard,
		untapAll,
		flipCard,
		transformCard,
		modifyLife,
		setPlayerCounter,
		shuffleLibrary,
		addToStack,
		removeFromStack,
		createToken,
		destroyToken,
		addCounter,
		removeCounter,
		setCounter,
		millCards,
		revealTopCards,
		scryCards,
		applyScryDecision,
		setRevealedTop,
		nextTurn,
		clearCombat,
		searchLibrary,
		selectLibraryCard,
		mulligan,
		keepHand,
		switchControlSeat,
		setCommand,
		reset,
		cleanup
	};
}

/**
 * Global multiplayer game store instance
 */
export const multiplayerGameStore = createMultiplayerGameStore();

// Derived stores for convenient access
// Based on playtest-game.ts lines 1233-1256

export const multiplayerPlayers = derived(multiplayerGameStore, ($game) => $game.players);

export const multiplayerLocalPlayer = derived(multiplayerGameStore, ($game) => {
	return $game.players.find((p) => p.playerId === $game.activeControlSeat) || null;
});

export const multiplayerOpponents = derived(multiplayerGameStore, ($game) => {
	return $game.players.filter((p) => p.playerId !== $game.activeControlSeat);
});

export const multiplayerBattlefield = derived(multiplayerGameStore, ($game) => $game.battlefield);

export const multiplayerExile = derived(multiplayerGameStore, ($game) => $game.exile);

export const multiplayerStack = derived(multiplayerGameStore, ($game) => $game.stack);

export const multiplayerActiveControlSeat = derived(
	multiplayerGameStore,
	($game) => $game.activeControlSeat
);

export const multiplayerIsInitialized = derived(
	multiplayerGameStore,
	($game) => $game.isInitialized
);

export const multiplayerIsConnected = derived(multiplayerGameStore, ($game) => $game.isConnected);
