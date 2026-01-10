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
	StartGameData,
	GameAssignDamageData,
	GameRollbackRequestData,
	GameRollbackCompleteData
} from '$lib/generated/mage/v1/websocket';
import type {
	GameView,
	PlayerView,
	CardView,
	ManaPoolView
} from '$lib/generated/mage/v1/models';

/**
 * Pending card play action for optimistic updates
 */
export interface PendingCardPlay {
	cardId: string;
	card: CardView;
	timestamp: number;
	sourceZone: 'hand' | 'battlefield';
	targetZone: 'battlefield' | 'stack';
}

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
	| 'multiAmount'
	| 'assignDamage'
	| 'librarySearch';

export interface GamePrompt {
	type: PromptType;
	message: string;
	data: unknown;
}

/**
 * Pending rollback request awaiting consent
 */
export interface PendingRollbackRequest {
	requestId: string;
	requestingPlayerId: string;
	requestingPlayerName: string;
	targetMessageId: number;
	targetMessageText: string;
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

	// Optimistic UI state
	pendingCardPlays: Map<string, PendingCardPlay>;
	cardsBeingPlayed: string[]; // Card IDs currently animating out

	// Rollback state
	pendingRollbackRequest: PendingRollbackRequest | null;
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
	showStack: false,
	pendingCardPlays: new Map(),
	cardsBeingPlayed: [],
	pendingRollbackRequest: null
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
					playerCount: updateData.game?.players?.length,
					stackCount: updateData.game?.stack?.length || 0,
					stackCards: updateData.game?.stack?.map(c => ({ id: c.id, name: c.name })) || []
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
					update((s) => {
						// Clear pending plays for cards that are now on battlefield or stack
						// (server state is source of truth)
						const newPending = new Map(s.pendingCardPlays);
						const newPlaying = [...s.cardsBeingPlayed];
						
						for (const [cardId] of s.pendingCardPlays) {
							// Check if card is now in battlefield or stack
							const onBattlefield = normalized.battlefield.some((c) => c.id === cardId);
							const onStack = normalized.stack.some((c) => c.id === cardId);
							if (onBattlefield || onStack) {
								newPending.delete(cardId);
								const playingIdx = newPlaying.indexOf(cardId);
								if (playingIdx >= 0) {
									newPlaying.splice(playingIdx, 1);
								}
							}
						}

						// Check for pending library search from server
						let newPrompt = s.pendingPrompt;
						if (normalized.pendingLibrarySearch) {
							console.log('[GameStore] Library search detected:', normalized.pendingLibrarySearch);
							newPrompt = {
								type: 'librarySearch',
								message: normalized.pendingLibrarySearch.message || 'Search your library',
								data: normalized.pendingLibrarySearch
							};
						} else if (s.pendingPrompt?.type === 'librarySearch') {
							// Clear library search prompt if it's no longer pending
							newPrompt = null;
						}

						return {
							...s,
							gameView: normalized,
							error: null,
							pendingCardPlays: newPending,
							cardsBeingPlayed: newPlaying,
							pendingPrompt: newPrompt
						};
					});
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
				
				update((s) => {
					// Rollback all pending card plays on error
					// The server has rejected the action(s)
					if (s.pendingCardPlays.size > 0) {
						console.log('[GameStore] Rolling back pending plays due to error:', 
							Array.from(s.pendingCardPlays.keys())
						);
					}
					
					return {
						...s,
						error: errorData.error,
						// Clear all pending plays - server state is truth
						pendingCardPlays: new Map(),
						cardsBeingPlayed: []
					};
				});
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

		// GAME_ASSIGN_DAMAGE - Combat damage assignment
		unsubscribers.push(
			websocketStore.on(CallbackMethod.GAME_ASSIGN_DAMAGE, (data) => {
				const assignDamageData = data as GameAssignDamageData;
				console.log('[GameStore] GAME_ASSIGN_DAMAGE:', assignDamageData);
				update((s) => ({
					...s,
					pendingPrompt: {
						type: 'assignDamage',
						message: assignDamageData.message,
						data: {
							attackerId: assignDamageData.attackerId,
							attackerName: assignDamageData.attackerName,
							attackerPower: assignDamageData.attackerPower,
							blockers: assignDamageData.blockers.map((b, i) => ({
								id: b.id,
								name: b.name,
								toughness: b.toughness,
								damage: b.markedDamage,
								order: b.order ?? i
							})),
							hasTrample: assignDamageData.hasTrample,
							defendingPlayerId: assignDamageData.defendingPlayerId,
							defendingPlayerName: assignDamageData.defendingPlayerName
						}
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

		// GAME_ROLLBACK_REQUEST - Opponent requesting rollback consent
		unsubscribers.push(
			websocketStore.on(CallbackMethod.GAME_ROLLBACK_REQUEST, (data) => {
				const rollbackData = data as GameRollbackRequestData;
				console.log('[GameStore] GAME_ROLLBACK_REQUEST:', rollbackData);
				update((s) => ({
					...s,
					pendingRollbackRequest: {
						requestId: rollbackData.requestId,
						requestingPlayerId: rollbackData.requestingPlayerId,
						requestingPlayerName: rollbackData.requestingPlayerName,
						targetMessageId: rollbackData.targetMessageId,
						targetMessageText: rollbackData.targetMessageText
					}
				}));
				toast.info(`${rollbackData.requestingPlayerName} wants to rollback to: "${rollbackData.targetMessageText}"`);
			})
		);

		// GAME_ROLLBACK_COMPLETE - Rollback was performed
		unsubscribers.push(
			websocketStore.on(CallbackMethod.GAME_ROLLBACK_COMPLETE, (data) => {
				const completeData = data as GameRollbackCompleteData;
				console.log('[GameStore] GAME_ROLLBACK_COMPLETE:', completeData);
				update((s) => ({
					...s,
					pendingRollbackRequest: null
				}));
				toast.success(`Game rolled back by ${completeData.initiatedByName}`);
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
			selectedCardIds: [],
			pendingCardPlays: new Map(),
			cardsBeingPlayed: []
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

	/**
	 * Add a pending card play for optimistic UI updates
	 * This immediately shows the card as "being played" before server confirmation
	 */
	function addPendingCardPlay(
		cardId: string,
		card: CardView,
		sourceZone: 'hand' | 'battlefield' = 'hand',
		targetZone: 'battlefield' | 'stack' = 'battlefield'
	): void {
		update((s) => {
			const newPending = new Map(s.pendingCardPlays);
			newPending.set(cardId, {
				cardId,
				card,
				timestamp: Date.now(),
				sourceZone,
				targetZone
			});
			return {
				...s,
				pendingCardPlays: newPending,
				cardsBeingPlayed: [...s.cardsBeingPlayed, cardId]
			};
		});
	}

	/**
	 * Remove a pending card play (on success or failure)
	 */
	function removePendingCardPlay(cardId: string): void {
		update((s) => {
			const newPending = new Map(s.pendingCardPlays);
			newPending.delete(cardId);
			return {
				...s,
				pendingCardPlays: newPending,
				cardsBeingPlayed: s.cardsBeingPlayed.filter((id) => id !== cardId)
			};
		});
	}

	/**
	 * Mark a card as animating (being played)
	 */
	function setCardPlaying(cardId: string): void {
		update((s) => ({
			...s,
			cardsBeingPlayed: s.cardsBeingPlayed.includes(cardId)
				? s.cardsBeingPlayed
				: [...s.cardsBeingPlayed, cardId]
		}));
	}

	/**
	 * Clear the playing state for a card (animation complete)
	 */
	function clearCardPlaying(cardId: string): void {
		update((s) => ({
			...s,
			cardsBeingPlayed: s.cardsBeingPlayed.filter((id) => id !== cardId)
		}));
	}

	/**
	 * Rollback a pending card play on server rejection
	 * Returns the card to its original zone
	 */
	function rollbackCardPlay(cardId: string): PendingCardPlay | null {
		const state = get({ subscribe });
		const pending = state.pendingCardPlays.get(cardId);

		if (!pending) {
			console.warn('[GameStore] Attempted to rollback non-pending card:', cardId);
			return null;
		}

		console.log('[GameStore] Rolling back card play:', {
			cardId,
			cardName: pending.card.name,
			sourceZone: pending.sourceZone
		});

		// Remove from pending and playing state
		removePendingCardPlay(cardId);

		return pending;
	}

	/**
	 * Get a pending card play by ID
	 */
	function getPendingCardPlay(cardId: string): PendingCardPlay | null {
		const state = get({ subscribe });
		return state.pendingCardPlays.get(cardId) || null;
	}

	/**
	 * Check if a card is currently being played (pending or animating)
	 */
	function isCardBeingPlayed(cardId: string): boolean {
		const state = get({ subscribe });
		return state.cardsBeingPlayed.includes(cardId) || state.pendingCardPlays.has(cardId);
	}

	/**
	 * Clear all pending card plays (e.g., on game reset)
	 */
	function clearAllPendingPlays(): void {
		update((s) => ({
			...s,
			pendingCardPlays: new Map(),
			cardsBeingPlayed: []
		}));
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
		reset,
		// Optimistic UI methods
		addPendingCardPlay,
		removePendingCardPlay,
		setCardPlaying,
		clearCardPlaying,
		rollbackCardPlay,
		getPendingCardPlay,
		isCardBeingPlayed,
		clearAllPendingPlays
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
 * Battlefield cards (sorted by UUID ascending for consistent ordering)
 */
export const battlefield = derived(gameStore, ($game) =>
	[...($game.gameView?.battlefield || [])].sort((a, b) => a.id.localeCompare(b.id))
);

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
 * Local player's hand (sorted by UUID ascending for consistent ordering)
 */
export const myHand = derived(gameStore, ($game) => {
	if (!$game.gameView || !$game.localPlayerId) return [];
	const player = $game.gameView.players.find((p) => p.playerId === $game.localPlayerId);
	return [...(player?.hand || [])].sort((a, b) => a.id.localeCompare(b.id));
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
 * Cards currently being played (animating)
 */
export const cardsBeingPlayed = derived(gameStore, ($game) => $game.cardsBeingPlayed);

/**
 * Pending card plays (for optimistic UI)
 */
export const pendingCardPlays = derived(gameStore, ($game) => $game.pendingCardPlays);

/**
 * Check if any cards are pending
 */
export const hasPendingPlays = derived(
	gameStore,
	($game) => $game.pendingCardPlays.size > 0 || $game.cardsBeingPlayed.length > 0
);

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
