/**
 * Multiplayer Game Store
 *
 * Server-synchronized game state management for multiplayer mode.
 * Based on playtest-game.ts structure but sends operations to server
 * and receives state updates via WebSocket.
 */

import { writable, derived, get } from 'svelte/store';
import type { CardView, GameView, PlayerView } from '$lib/generated/mage/v1/models';
import { CallbackMethod } from '$lib/generated/mage/v1/websocket';
import type { GameUpdateData, GameInitData, StartGameData } from '$lib/generated/mage/v1/websocket';
import { websocketStore } from './websocket';
import * as directActions from '$lib/api/direct-actions';
import { joinGame, fetchGameView } from '$lib/api/game';
import { auth } from './auth';
import { goto } from '$app/navigation';
import type { PlaytestPlayer, PlaytestLogEntry, ScrySession } from '$lib/types/gamestore';

/**
 * Multiplayer game state
 * Based on playtest-game.ts lines 45-59
 * Added server sync fields: isConnected, pendingActions
 */
export interface MultiplayerGameState {
	gameId: string;
	activeControlSeat: string; // Which player perspective you're controlling (only the logged-in player in multiplayer mode, but other perspectives in playtest mode)
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

/**
 * Converts a proto PlayerView to a PlaytestPlayer object.
 * Handles differences between proto structure and client expectations.
 * From Task 1.3 - Proto GameView to State Mapping
 */
function convertPlayerViewToPlaytestPlayer(pv: PlayerView): PlaytestPlayer {
	return {
		playerId: pv.playerId || '',
		name: pv.name || '',
		life: Number(pv.life) || 0,
		poison: Number(pv.poison) || 0,
		energy: Number(pv.energy) || 0,
		libraryCount: Number(pv.libraryCount) || 0,
		handCount: Number(pv.handCount) || 0,
		hand: pv.hand || [],
		library: pv.library || [], // Will be populated for viewing player, empty for opponents
		graveyard: pv.graveyard || [],
		manaPool: pv.manaPool || {
			white: 0,
			blue: 0,
			black: 0,
			red: 0,
			green: 0,
			colorless: 0
		},
		keptHand: pv.keptHand || false,
		mulliganCount: Number(pv.mulliganCount) || 0,
		revealedTopCard: false // Not in proto yet, default to false
	};
}

/**
 * Normalizes a proto GameView by ensuring all array fields are initialized.
 * Proto messages omit empty arrays, so we need to provide defaults.
 * From Task 1.3 - Proto GameView to State Mapping
 */
function normalizeProtoGameView(game: GameView): GameView {
	return {
		...game,
		players: game.players || [],
		battlefield: game.battlefield || [],
		stack: game.stack || [],
		exile: game.exile || [],
		command: game.command || []
	};
}

/**
 * Maps a normalized proto GameView to MultiplayerGameState structure.
 * Converts proto PlayerView[] to PlaytestPlayer[] and maps all fields.
 * From Task 1.3 - Proto GameView to State Mapping
 */
function mapProtoGameViewToState(game: GameView): Partial<MultiplayerGameState> {
	return {
		gameId: game.gameId || '',
		turn: Number(game.turn) || 0,
		activePlayerId: game.activePlayerId || '',
		activeControlSeat: game.activeControlSeat || '', // From Phase 1 Task 1.1
		battlefield: game.battlefield || [],
		exile: game.exile || [],
		stack: game.stack || [],
		command: game.command || [],
		players: (game.players || []).map(convertPlayerViewToPlaytestPlayer),
		mulliganType: 'london', // Default, could be in proto later
		freeMulligans: 0, // Default, could be in proto later
		log: [] // Not in current proto, default to empty
	};
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
	 * Pattern from game.legacy.ts lines 140-200
	 */
	function subscribeToGameEvents() {
		// Unsubscribe from previous subscriptions
		unsubscribeFromEvents();

		// START_GAME - Game is starting (from game.legacy.ts lines 147-156)
		unsubscribers.push(
			websocketStore.on(CallbackMethod.START_GAME, (data) => {
				const startData = data as StartGameData;
				console.log('[MultiplayerGame] START_GAME received:', startData);

				// Prepare for initialization - game is starting
				update((state) => ({
					...state,
					isInitialized: false,
					isConnected: false
				}));
			})
		);

		// GAME_INIT - Initial game state (from game.legacy.ts lines 158-174)
		// Note: Server may not send this, but we handle it for compatibility
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

		// GAME_UPDATE - State update (from game.legacy.ts lines 176-226)
		unsubscribers.push(
			websocketStore.on(CallbackMethod.GAME_UPDATE, (data) => {
				const updateData = data as GameUpdateData;
				console.log('[MultiplayerGame] GAME_UPDATE received:', updateData);

				if (updateData.game) {
					// Normalize proto GameView (handle empty arrays)
					const normalized = normalizeProtoGameView(updateData.game);

					// Map proto structure to our state structure
					const mappedState = mapProtoGameViewToState(normalized);

					// Apply server state to store
					update((state) => ({
						...state,
						...mappedState,
						isConnected: true,
						isInitialized: true,
						// Clear pending actions (server is source of truth)
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
	 * Modified: Initialization happens server-side, we subscribe and join the game
	 * Pattern from debug/+page.svelte lines 87-127
	 */
	async function initialize(gameId: string): Promise<void> {
		update((state) => ({
			...state,
			gameId,
			isInitialized: false,
			isConnected: false
		}));

		// Subscribe to WebSocket updates for this game
		subscribeToGameEvents();

		try {
			console.log('[MultiplayerGame] Joining game:', gameId);
			// Join the game to register with the server
			try {
				await joinGame(gameId);
			} catch (err) {
				console.error('[MultiplayerGame] Failed to join game:', err);
				goto('/lobby');
				return;
			}
			console.log('[MultiplayerGame] Joined game successfully');

			// Fetch initial game view to get current state
			// This ensures we have state even if we missed the START_GAME event
			const currentAuth = get(auth);
			const playerId = currentAuth.user?.username || currentAuth.user?.id || '';
			if (playerId) {
				console.log('[MultiplayerGame] Fetching initial game view...');

				// Poll for game view with exponential backoff
				// This handles cases where game engine is still initializing
				let retries = 0;
				const maxRetries = 5;
				const baseDelay = 500; // ms
				let gameView: GameView | null = null;

				while (retries < maxRetries) {
					try {
						gameView = await fetchGameView(gameId, playerId);
						console.log(`[MultiplayerGame] Fetch attempt ${retries + 1}:`, gameView);

						// Check if we got valid game state
						if (gameView && gameView.players && gameView.players.length > 0) {
							console.log('[MultiplayerGame] Got valid game view with players');
							break;
						}

						// If gameView is null or has no players, retry after delay
						retries++;
						if (retries < maxRetries) {
							const delay = baseDelay * Math.pow(2, retries - 1);
							console.log(`[MultiplayerGame] Game state not ready, retrying in ${delay}ms...`);
							await new Promise((resolve) => setTimeout(resolve, delay));
						}
					} catch (err) {
						console.error(`[MultiplayerGame] Fetch attempt ${retries + 1} failed:`, err);
						retries++;
						if (retries < maxRetries) {
							const delay = baseDelay * Math.pow(2, retries - 1);
							await new Promise((resolve) => setTimeout(resolve, delay));
						}
					}
				}

				if (gameView && gameView.players && gameView.players.length > 0) {
					// Normalize and map the fetched GameView
					const normalized = normalizeProtoGameView(gameView);
					const mappedState = mapProtoGameViewToState(normalized);

					// Apply initial state
					update((state) => ({
						...state,
						...mappedState,
						isConnected: true,
						isInitialized: true,
						pendingActions: []
					}));
					console.log('[MultiplayerGame] Initialized with polled game state');
				} else {
					// After all retries, still no valid game state
					console.warn(
						'[MultiplayerGame] No valid game state after polling, waiting for GAME_UPDATE'
					);
					update((state) => ({
						...state,
						isConnected: true,
						isInitialized: false // Will be initialized when GAME_UPDATE arrives
					}));
				}
			} else {
				console.warn('[MultiplayerGame] No player ID available, waiting for GAME_UPDATE');
				// Will be initialized when first GAME_UPDATE arrives
			}
		} catch (err) {
			console.error('[MultiplayerGame] Failed to join game:', err);
			// Don't throw - let WebSocket events handle initialization
			// If join fails, we'll wait for GAME_UPDATE events
		}

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
	 * Modified: Send to server via direct-actions API
	 */
	function millCards(playerId: string, count: number): void {
		const state = get({ subscribe });
		directActions.millCards(state.gameId, playerId, count);
		console.log('[MultiplayerGame] millCards:', { playerId, count });
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
	 * Modified: Send to server via direct-actions API
	 * Note: This is a simplified scry that just initiates the action.
	 * Full scry UI with card selection would need additional implementation.
	 */
	function scryCards(playerId: string, count: number): ScrySession | null {
		const state = get({ subscribe });
		directActions.scryCards(state.gameId, playerId, count);
		console.log('[MultiplayerGame] scryCards:', { playerId, count });
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
	 * Modified: Send to server via direct-actions API
	 */
	function setRevealedTop(playerId: string, revealed: boolean): void {
		const state = get({ subscribe });
		directActions.setRevealedTop(state.gameId, playerId, revealed);
		console.log('[MultiplayerGame] setRevealedTop:', { playerId, revealed });
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
	 * Modified: Send to server via direct-actions API
	 */
	function mulligan(playerId: string): Promise<void> {
		const state = get({ subscribe });
		console.log('[MultiplayerGame] mulligan:', { playerId });
		return directActions.mulligan(state.gameId, playerId);
	}

	/**
	 * Keep hand (no mulligan)
	 * From playtest-game.ts lines 1151-1161
	 * Modified: Send to server via direct-actions API
	 */
	function keepHand(playerId: string): Promise<void> {
		const state = get({ subscribe });
		console.log('[MultiplayerGame] keepHand:', { playerId });
		return directActions.keepHand(state.gameId, playerId);
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

// Re-export types for backward compatibility
export type { PlaytestPlayer, PlaytestLogEntry, ScrySession } from '$lib/types/gamestore';
