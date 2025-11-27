/**
 * Game Store
 * Manages game state and WebSocket event subscriptions for real-time game updates
 */

import { writable, derived, get } from 'svelte/store';
import { websocketStore } from './websocket';
import { toast } from './toast';
import { CallbackMethod } from '$lib/generated/mage/v1/websocket';
import type {
	GameInitData,
	GameUpdateData,
	GameUpdateAndInformData,
	GameInformPersonalData,
	GameErrorData,
	GameTargetData,
	GameChooseAbilityData,
	GameChoosePileData,
	GameChoiceData,
	GameAskData,
	GameSelectData,
	GamePlayManaData,
	GamePlayXManaData,
	GameGetAmountData,
	GameGetMultiAmountData,
	GameOverData,
	StartGameData
} from '$lib/generated/mage/v1/websocket';
import type {
	GameView,
	PlayerView,
	CardView,
	ManaPoolView
} from '$lib/generated/mage/v1/models';

/**
 * Pending prompt types for user interaction
 */
export type PromptType =
	| 'target'
	| 'ability'
	| 'pile'
	| 'choice'
	| 'ask'
	| 'select'
	| 'mana'
	| 'xmana'
	| 'amount'
	| 'multiAmount';

export interface GamePrompt {
	type: PromptType;
	message: string;
	data: unknown;
}

/**
 * Game store state
 */
export interface GameStoreState {
	// Connection state
	gameId: string | null;
	localPlayerId: string | null;
	isConnected: boolean;
	isLoading: boolean;
	error: string | null;

	// Game state from server
	gameView: GameView | null;

	// Pending prompts requiring user action
	pendingPrompt: GamePrompt | null;

	// Game over state
	gameOver: boolean;
	winner: string | null;
	results: Array<{ playerName: string; wins: number; losses: number; quit: boolean }>;

	// UI state
	selectedCardIds: string[];
	showStack: boolean;
}

const initialState: GameStoreState = {
	gameId: null,
	localPlayerId: null,
	isConnected: false,
	isLoading: false,
	error: null,
	gameView: null,
	pendingPrompt: null,
	gameOver: false,
	winner: null,
	results: [],
	selectedCardIds: [],
	showStack: false
};

/**
 * Create game store
 */
function createGameStore() {
	const { subscribe, set, update } = writable<GameStoreState>(initialState);

	// Track unsubscribe functions for cleanup
	let unsubscribers: Array<() => void> = [];

	/**
	 * Subscribe to all game-related WebSocket events
	 */
	function subscribeToGameEvents() {
		// Unsubscribe from previous subscriptions
		unsubscribeFromEvents();
		
		console.log('[GameStore] Subscribing to game WebSocket events...');

		// START_GAME - Game is starting
		unsubscribers.push(
			websocketStore.on(CallbackMethod.START_GAME, (data) => {
				const startData = data as StartGameData;
				console.log('[GameStore] START_GAME:', startData);
				update((s) => ({
					...s,
					isLoading: true,
					error: null
				}));
			})
		);

		// GAME_INIT - Initial game state
		unsubscribers.push(
			websocketStore.on(CallbackMethod.GAME_INIT, (data) => {
				const initData = data as GameInitData;
				console.log('[GameStore] GAME_INIT:', initData);
				if (initData.game) {
					const normalized = normalizeGameView(initData.game);
					update((s) => ({
						...s,
						gameView: normalized,
						isLoading: false,
						isConnected: true,
						error: null
					}));
				}
			})
		);

		// GAME_UPDATE - State update
		unsubscribers.push(
			websocketStore.on(CallbackMethod.GAME_UPDATE, (data) => {
				const updateData = data as GameUpdateData;
				console.log('[GameStore] GAME_UPDATE received:', {
					hasGame: !!updateData.game,
					state: updateData.game?.state,
					turn: updateData.game?.turn,
					phase: updateData.game?.phase,
					playerCount: updateData.game?.players?.length
				});
				if (updateData.game) {
					// Log each player's hand
					updateData.game.players?.forEach((player) => {
						console.log('[GameStore] Player hand data:', {
							playerId: player.playerId,
							playerName: player.name,
							handCount: player.handCount,
							handCards: player.hand?.map((c) => c.name) || [],
							handCardIds: player.hand?.map((c) => c.id) || []
						});
					});

					const normalized = normalizeGameView(updateData.game);
					update((s) => ({
						...s,
						gameView: normalized,
						error: null
					}));
				}
			})
		);

		// GAME_UPDATE_AND_INFORM - State update with message
		unsubscribers.push(
			websocketStore.on(CallbackMethod.GAME_UPDATE_AND_INFORM, (data) => {
				const updateData = data as GameUpdateAndInformData;
				console.log('[GameStore] GAME_UPDATE_AND_INFORM:', updateData);
				if (updateData.game) {
					const normalized = normalizeGameView(updateData.game);
					update((s) => ({
						...s,
						gameView: normalized,
						error: null
					}));
				}
			})
		);

		// GAME_INFORM_PERSONAL - Personal message
		unsubscribers.push(
			websocketStore.on(CallbackMethod.GAME_INFORM_PERSONAL, (data) => {
				const informData = data as GameInformPersonalData;
				console.log('[GameStore] GAME_INFORM_PERSONAL:', informData);
			})
		);

		// GAME_ERROR - Error message
		unsubscribers.push(
			websocketStore.on(CallbackMethod.GAME_ERROR, (data) => {
				const errorData = data as GameErrorData;
				console.error('[GameStore] GAME_ERROR:', errorData);
				
				// Show toast notification for game errors
				if (errorData.error) {
					// Clean up the error message for display
					let errorMessage = errorData.error;
					// Remove the "action failed and state restored: " prefix if present
					if (errorMessage.includes('action failed and state restored: ')) {
						errorMessage = errorMessage.replace('action failed and state restored: ', '');
					}
					toast.error(errorMessage);
				}
				
				update((s) => ({
					...s,
					error: errorData.error
				}));
			})
		);

		// GAME_TARGET - Target selection required
		unsubscribers.push(
			websocketStore.on(CallbackMethod.GAME_TARGET, (data) => {
				const targetData = data as GameTargetData;
				console.log('[GameStore] GAME_TARGET:', targetData);
				update((s) => ({
					...s,
					pendingPrompt: {
						type: 'target',
						message: targetData.message,
						data: targetData
					}
				}));
			})
		);

		// GAME_CHOOSE_ABILITY - Ability selection required
		unsubscribers.push(
			websocketStore.on(CallbackMethod.GAME_CHOOSE_ABILITY, (data) => {
				const abilityData = data as GameChooseAbilityData;
				console.log('[GameStore] GAME_CHOOSE_ABILITY:', abilityData);
				update((s) => ({
					...s,
					pendingPrompt: {
						type: 'ability',
						message: abilityData.message,
						data: abilityData
					}
				}));
			})
		);

		// GAME_CHOOSE_PILE - Pile selection required
		unsubscribers.push(
			websocketStore.on(CallbackMethod.GAME_CHOOSE_PILE, (data) => {
				const pileData = data as GameChoosePileData;
				console.log('[GameStore] GAME_CHOOSE_PILE:', pileData);
				update((s) => ({
					...s,
					pendingPrompt: {
						type: 'pile',
						message: pileData.message,
						data: pileData
					}
				}));
			})
		);

		// GAME_CHOOSE_CHOICE - Choice selection required
		unsubscribers.push(
			websocketStore.on(CallbackMethod.GAME_CHOOSE_CHOICE, (data) => {
				const choiceData = data as GameChoiceData;
				console.log('[GameStore] GAME_CHOOSE_CHOICE:', choiceData);
				update((s) => ({
					...s,
					pendingPrompt: {
						type: 'choice',
						message: choiceData.message,
						data: choiceData
					}
				}));
			})
		);

		// GAME_ASK - Yes/No question
		unsubscribers.push(
			websocketStore.on(CallbackMethod.GAME_ASK, (data) => {
				const askData = data as GameAskData;
				console.log('[GameStore] GAME_ASK:', askData);
				update((s) => ({
					...s,
					pendingPrompt: {
						type: 'ask',
						message: askData.message,
						data: askData
					}
				}));
			})
		);

		// GAME_SELECT - Selection required
		unsubscribers.push(
			websocketStore.on(CallbackMethod.GAME_SELECT, (data) => {
				const selectData = data as GameSelectData;
				console.log('[GameStore] GAME_SELECT:', selectData);
				update((s) => ({
					...s,
					pendingPrompt: {
						type: 'select',
						message: selectData.message,
						data: selectData
					}
				}));
			})
		);

		// GAME_PLAY_MANA - Mana payment required
		unsubscribers.push(
			websocketStore.on(CallbackMethod.GAME_PLAY_MANA, (data) => {
				const manaData = data as GamePlayManaData;
				console.log('[GameStore] GAME_PLAY_MANA:', manaData);
				update((s) => ({
					...s,
					pendingPrompt: {
						type: 'mana',
						message: manaData.message,
						data: manaData
					}
				}));
			})
		);

		// GAME_PLAY_XMANA - X mana selection
		unsubscribers.push(
			websocketStore.on(CallbackMethod.GAME_PLAY_XMANA, (data) => {
				const xmanaData = data as GamePlayXManaData;
				console.log('[GameStore] GAME_PLAY_XMANA:', xmanaData);
				update((s) => ({
					...s,
					pendingPrompt: {
						type: 'xmana',
						message: xmanaData.message,
						data: xmanaData
					}
				}));
			})
		);

		// GAME_GET_AMOUNT - Amount selection
		unsubscribers.push(
			websocketStore.on(CallbackMethod.GAME_GET_AMOUNT, (data) => {
				const amountData = data as GameGetAmountData;
				console.log('[GameStore] GAME_GET_AMOUNT:', amountData);
				update((s) => ({
					...s,
					pendingPrompt: {
						type: 'amount',
						message: amountData.message,
						data: amountData
					}
				}));
			})
		);

		// GAME_GET_MULTI_AMOUNT - Multiple amount selection
		unsubscribers.push(
			websocketStore.on(CallbackMethod.GAME_GET_MULTI_AMOUNT, (data) => {
				const multiAmountData = data as GameGetMultiAmountData;
				console.log('[GameStore] GAME_GET_MULTI_AMOUNT:', multiAmountData);
				update((s) => ({
					...s,
					pendingPrompt: {
						type: 'multiAmount',
						message: multiAmountData.message,
						data: multiAmountData
					}
				}));
			})
		);

		// GAME_OVER - Game ended
		unsubscribers.push(
			websocketStore.on(CallbackMethod.GAME_OVER, (data) => {
				const overData = data as GameOverData;
				console.log('[GameStore] GAME_OVER:', overData);
				update((s) => ({
					...s,
					gameOver: true,
					winner: overData.winner,
					results: overData.results || [],
					pendingPrompt: null
				}));
			})
		);
		
		console.log(`[GameStore] Subscribed to ${unsubscribers.length} game event types`);
	}

	/**
	 * Unsubscribe from all events
	 */
	function unsubscribeFromEvents() {
		unsubscribers.forEach((unsub) => unsub());
		unsubscribers = [];
	}

	/**
	 * Initialize the store for a game
	 */
	function initGame(gameId: string, localPlayerId: string) {
		update((s) => ({
			...s,
			gameId,
			localPlayerId,
			isLoading: true,
			error: null,
			gameOver: false,
			winner: null,
			results: [],
			pendingPrompt: null,
			selectedCardIds: []
		}));

		subscribeToGameEvents();
	}

	/**
	 * Normalize game view to ensure all arrays are initialized
	 * This is needed because protojson omits empty arrays
	 */
	function normalizeGameView(gameView: GameView): GameView {
		return {
			...gameView,
			players: (gameView.players || []).map((player) => ({
				...player,
				hand: player.hand || [],
				graveyard: player.graveyard || [],
				manaPool: player.manaPool || { white: 0, blue: 0, black: 0, red: 0, green: 0, colorless: 0 }
			})),
			battlefield: gameView.battlefield || [],
			stack: gameView.stack || [],
			exile: gameView.exile || [],
			command: gameView.command || [],
			messages: gameView.messages || []
		};
	}

	/**
	 * Update game view directly (for initial load via RPC)
	 */
	function setGameView(gameView: GameView) {
		console.log('[GameStore] setGameView called:', {
			state: gameView.state,
			turn: gameView.turn,
			phase: gameView.phase,
			playerCount: gameView.players?.length
		});

		// Log each player's hand
		gameView.players?.forEach((player) => {
			console.log('[GameStore] setGameView - Player hand:', {
				playerId: player.playerId,
				playerName: player.name,
				handCount: player.handCount,
				handCards: player.hand?.map((c) => c.name) || [],
				handCardIds: player.hand?.map((c) => c.id) || []
			});
		});

		const normalized = normalizeGameView(gameView);
		update((s) => ({
			...s,
			gameView: normalized,
			isLoading: false,
			isConnected: true,
			error: null
		}));
	}

	/**
	 * Clear pending prompt
	 */
	function clearPrompt() {
		update((s) => ({
			...s,
			pendingPrompt: null
		}));
	}

	/**
	 * Set error
	 */
	function setError(error: string) {
		update((s) => ({
			...s,
			error,
			isLoading: false
		}));
	}

	/**
	 * Toggle card selection
	 */
	function toggleCardSelection(cardId: string) {
		update((s) => {
			const idx = s.selectedCardIds.indexOf(cardId);
			if (idx >= 0) {
				return {
					...s,
					selectedCardIds: s.selectedCardIds.filter((id) => id !== cardId)
				};
			} else {
				return {
					...s,
					selectedCardIds: [...s.selectedCardIds, cardId]
				};
			}
		});
	}

	/**
	 * Clear card selection
	 */
	function clearSelection() {
		update((s) => ({
			...s,
			selectedCardIds: []
		}));
	}

	/**
	 * Toggle stack visibility
	 */
	function toggleStack() {
		update((s) => ({
			...s,
			showStack: !s.showStack
		}));
	}

	/**
	 * Reset store to initial state
	 */
	function reset() {
		unsubscribeFromEvents();
		set(initialState);
	}

	return {
		subscribe,
		initGame,
		setGameView,
		clearPrompt,
		setError,
		toggleCardSelection,
		clearSelection,
		toggleStack,
		reset
	};
}

/**
 * Global game store instance
 */
export const gameStore = createGameStore();

// Derived stores for convenient access to specific parts of game state

/**
 * Current game view
 */
export const gameView = derived(gameStore, ($game) => $game.gameView);

/**
 * All players in the game
 */
export const players = derived(gameStore, ($game) => $game.gameView?.players || []);

/**
 * Local player (the user playing)
 */
export const localPlayer = derived(gameStore, ($game) => {
	if (!$game.gameView || !$game.localPlayerId) return null;
	return $game.gameView.players.find((p) => p.playerId === $game.localPlayerId) || null;
});

/**
 * Opponents (all other players)
 */
export const opponents = derived(gameStore, ($game) => {
	if (!$game.gameView || !$game.localPlayerId) return [];
	return $game.gameView.players.filter((p) => p.playerId !== $game.localPlayerId);
});

/**
 * Active player (whose turn it is)
 */
export const activePlayer = derived(gameStore, ($game) => {
	if (!$game.gameView) return null;
	return (
		$game.gameView.players.find((p) => p.playerId === $game.gameView!.activePlayerId) || null
	);
});

/**
 * Priority player (who has priority)
 */
export const priorityPlayer = derived(gameStore, ($game) => {
	if (!$game.gameView) return null;
	return (
		$game.gameView.players.find((p) => p.playerId === $game.gameView!.priorityPlayerId) || null
	);
});

/**
 * Whether the local player has priority
 */
export const hasPriority = derived(gameStore, ($game) => {
	if (!$game.gameView || !$game.localPlayerId) return false;
	return $game.gameView.priorityPlayerId === $game.localPlayerId;
});

/**
 * Whether it's the local player's turn
 */
export const isMyTurn = derived(gameStore, ($game) => {
	if (!$game.gameView || !$game.localPlayerId) return false;
	return $game.gameView.activePlayerId === $game.localPlayerId;
});

/**
 * Current phase
 */
export const currentPhase = derived(gameStore, ($game) => $game.gameView?.phase || '');

/**
 * Current step
 */
export const currentStep = derived(gameStore, ($game) => $game.gameView?.step || '');

/**
 * Current turn number
 */
export const currentTurn = derived(gameStore, ($game) => $game.gameView?.turn || 0);

/**
 * Battlefield cards
 */
export const battlefield = derived(gameStore, ($game) => $game.gameView?.battlefield || []);

/**
 * Stack
 */
export const stack = derived(gameStore, ($game) => $game.gameView?.stack || []);

/**
 * Exile zone
 */
export const exile = derived(gameStore, ($game) => $game.gameView?.exile || []);

/**
 * Command zone
 */
export const command = derived(gameStore, ($game) => $game.gameView?.command || []);

/**
 * Combat state
 */
export const combat = derived(gameStore, ($game) => $game.gameView?.combat || null);

/**
 * Game messages/log
 */
export const gameMessages = derived(gameStore, ($game) => $game.gameView?.messages || []);

/**
 * Local player's hand
 */
export const myHand = derived(gameStore, ($game) => {
	if (!$game.gameView || !$game.localPlayerId) return [];
	const player = $game.gameView.players.find((p) => p.playerId === $game.localPlayerId);
	return player?.hand || [];
});

/**
 * Local player's graveyard
 */
export const myGraveyard = derived(gameStore, ($game) => {
	if (!$game.gameView || !$game.localPlayerId) return [];
	const player = $game.gameView.players.find((p) => p.playerId === $game.localPlayerId);
	return player?.graveyard || [];
});

/**
 * Local player's mana pool
 */
export const myManaPool = derived(gameStore, ($game): ManaPoolView => {
	if (!$game.gameView || !$game.localPlayerId) {
		return { white: 0, blue: 0, black: 0, red: 0, green: 0, colorless: 0 };
	}
	const player = $game.gameView.players.find((p) => p.playerId === $game.localPlayerId);
	return player?.manaPool || { white: 0, blue: 0, black: 0, red: 0, green: 0, colorless: 0 };
});

/**
 * Local player's life total
 */
export const myLife = derived(gameStore, ($game) => {
	if (!$game.gameView || !$game.localPlayerId) return 0;
	const player = $game.gameView.players.find((p) => p.playerId === $game.localPlayerId);
	return player?.life || 0;
});

/**
 * Pending prompt requiring user action
 */
export const pendingPrompt = derived(gameStore, ($game) => $game.pendingPrompt);

/**
 * Game over state
 */
export const gameOver = derived(gameStore, ($game) => $game.gameOver);

/**
 * Game winner
 */
export const winner = derived(gameStore, ($game) => $game.winner);

/**
 * Error state
 */
export const gameError = derived(gameStore, ($game) => $game.error);

/**
 * Loading state
 */
export const isLoading = derived(gameStore, ($game) => $game.isLoading);

/**
 * Selected card IDs
 */
export const selectedCards = derived(gameStore, ($game) => $game.selectedCardIds);

/**
 * Get player by ID
 */
export function getPlayerById(playerId: string): PlayerView | null {
	const state = get(gameStore);
	if (!state.gameView) return null;
	return state.gameView.players.find((p) => p.playerId === playerId) || null;
}

/**
 * Get card by ID from any zone
 */
export function getCardById(cardId: string): CardView | null {
	const state = get(gameStore);
	if (!state.gameView) return null;

	// Search battlefield
	const bfCard = state.gameView.battlefield.find((c) => c.id === cardId);
	if (bfCard) return bfCard;

	// Search stack
	const stackCard = state.gameView.stack.find((c) => c.id === cardId);
	if (stackCard) return stackCard;

	// Search all players' hands and graveyards
	for (const player of state.gameView.players) {
		const handCard = player.hand.find((c) => c.id === cardId);
		if (handCard) return handCard;

		const gyCard = player.graveyard.find((c) => c.id === cardId);
		if (gyCard) return gyCard;
	}

	// Search exile
	const exileCard = state.gameView.exile.find((c) => c.id === cardId);
	if (exileCard) return exileCard;

	// Search command zone
	const cmdCard = state.gameView.command.find((c) => c.id === cardId);
	if (cmdCard) return cmdCard;

	return null;
}
