/**
 * Playtest Game Store
 * 
 * Client-side only game state management for playtest mode.
 * No server communication - all state is local.
 */

import { writable, derived } from 'svelte/store';
import { browser } from '$app/environment';
import type { CardView, ManaPoolView } from '$lib/generated/mage/v1/models';

/**
 * Local player state for playtest
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
}

/**
 * Playtest game state
 */
export interface PlaytestGameState {
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
}

const initialState: PlaytestGameState = {
	gameId: '',
	activeControlSeat: '',
	players: [],
	battlefield: [],
	exile: [],
	stack: [],
	command: [],
	turn: 1,
	activePlayerId: '',
	isInitialized: false
};

const PLAYTEST_STORAGE_KEY = 'mage.playtest.state.v1';
const PLAYTEST_STORAGE_VERSION = 1;

type PersistedPlaytestState = {
	version: number;
	savedAt: number;
	state: PlaytestGameState;
};

function loadPersistedPlaytestState(): PlaytestGameState | null {
	if (!browser) return null;

	try {
		const raw = localStorage.getItem(PLAYTEST_STORAGE_KEY);
		if (!raw) return null;

		const parsed = JSON.parse(raw) as PersistedPlaytestState;
		if (!parsed || typeof parsed !== 'object') return null;
		if (parsed.version !== PLAYTEST_STORAGE_VERSION) return null;
		if (!parsed.state || typeof parsed.state !== 'object') return null;

		// Basic validation to avoid crashing on bad/stale data
		if (parsed.state.isInitialized !== true) return null;
		if (!Array.isArray(parsed.state.players) || parsed.state.players.length === 0) return null;
		if (typeof parsed.state.gameId !== 'string' || parsed.state.gameId.length === 0) return null;

		return parsed.state;
	} catch (err) {
		console.warn('[PlaytestGame] Failed to load persisted state:', err);
		return null;
	}
}

function persistPlaytestState(state: PlaytestGameState): void {
	if (!browser) return;
	if (!state.isInitialized) return;

	try {
		const payload: PersistedPlaytestState = {
			version: PLAYTEST_STORAGE_VERSION,
			savedAt: Date.now(),
			state
		};
		localStorage.setItem(PLAYTEST_STORAGE_KEY, JSON.stringify(payload));
	} catch (err) {
		console.warn('[PlaytestGame] Failed to persist state:', err);
	}
}

function clearPersistedPlaytestState(): void {
	if (!browser) return;
	try {
		localStorage.removeItem(PLAYTEST_STORAGE_KEY);
	} catch (err) {
		console.warn('[PlaytestGame] Failed to clear persisted state:', err);
	}
}

/**
 * Create playtest game store
 */
function createPlaytestGameStore() {
	const hydrated = loadPersistedPlaytestState();
	const { subscribe, set, update } = writable<PlaytestGameState>(hydrated ?? initialState);

	// Persist any meaningful state changes (client-only)
	subscribe((state) => {
		persistPlaytestState(state);
	});

	/**
	 * Initialize game state with players and their decks
	 */
	function initialize(gameId: string, players: PlaytestPlayer[]): void {
		if (players.length === 0) {
			console.error('[PlaytestGame] Cannot initialize with no players');
			return;
		}

		set({
			gameId,
			activeControlSeat: players[0].playerId,
			players,
			battlefield: [],
			exile: [],
			stack: [],
			command: [],
			turn: 1,
			activePlayerId: players[0].playerId,
			isInitialized: true
		});

		console.log('[PlaytestGame] Initialized with', players.length, 'players');
	}

	/**
	 * Set command zone cards (e.g. commanders) for the current playtest.
	 */
	function setCommand(cards: CardView[]): void {
		update((state) => ({
			...state,
			command: cards || []
		}));
	}

	/**
	 * Switch active control seat (which player you're controlling)
	 */
	function switchControlSeat(playerId: string): void {
		update(state => ({
			...state,
			activeControlSeat: playerId
		}));
	}

	/**
	 * Draw cards for a player
	 */
	function drawCards(playerId: string, count: number): void {
		update(state => {
			const playerIndex = state.players.findIndex(p => p.playerId === playerId);
			if (playerIndex === -1) {
				console.error('[PlaytestGame] Player not found:', playerId);
				return state;
			}

			const player = state.players[playerIndex];
			const drawn = player.library.splice(0, Math.min(count, player.library.length));
			
			// Update zone and make cards visible
			drawn.forEach(card => {
				card.zone = 1; // HAND
				card.faceDown = false;
			});

			const newPlayers = [...state.players];
			newPlayers[playerIndex] = {
				...player,
				hand: [...player.hand, ...drawn],
				handCount: player.hand.length + drawn.length,
				libraryCount: player.library.length
			};

			return {
				...state,
				players: newPlayers
			};
		});
	}

	/**
	 * Play a card from hand to battlefield
	 */
	function playCard(cardId: string, tapped: boolean = false): void {
		update(state => {
			const controllingPlayer = state.players.find(p => p.playerId === state.activeControlSeat);
			if (!controllingPlayer) return state;

			const cardIndex = controllingPlayer.hand.findIndex(c => c.id === cardId);
			if (cardIndex === -1) {
				console.error('[PlaytestGame] Card not found in hand:', cardId);
				return state;
			}

			const card = controllingPlayer.hand[cardIndex];
			const newHand = [...controllingPlayer.hand];
			newHand.splice(cardIndex, 1);

			// Update card properties
			card.zone = 2; // BATTLEFIELD
			card.controllerId = state.activeControlSeat;
			card.tapped = tapped;
			card.faceDown = false;

			// Update player
			const newPlayers = state.players.map(p => 
				p.playerId === state.activeControlSeat
					? { ...p, hand: newHand, handCount: newHand.length }
					: p
			);

			return {
				...state,
				players: newPlayers,
				battlefield: [...state.battlefield, card]
			};
		});
	}

	/**
	 * Move a card to a different zone
	 */
	function moveCardToZone(cardId: string, targetZone: string): void {
		update(state => {
			// Find the card in any zone
			let card: CardView | null = null;
			let sourceZone: string | null = null;

			// Check battlefield
			const bfIndex = state.battlefield.findIndex(c => c.id === cardId);
			if (bfIndex !== -1) {
				card = state.battlefield[bfIndex];
				sourceZone = 'battlefield';
			}

			// Check player zones
			if (!card) {
				for (const player of state.players) {
					const handIndex = player.hand.findIndex(c => c.id === cardId);
					if (handIndex !== -1) {
						card = player.hand[handIndex];
						sourceZone = `hand:${player.playerId}`;
						break;
					}

					const libraryIndex = player.library.findIndex(c => c.id === cardId);
					if (libraryIndex !== -1) {
						card = player.library[libraryIndex];
						sourceZone = `library:${player.playerId}`;
						break;
					}
					
					const graveyardIndex = player.graveyard.findIndex(c => c.id === cardId);
					if (graveyardIndex !== -1) {
						card = player.graveyard[graveyardIndex];
						sourceZone = `graveyard:${player.playerId}`;
						break;
					}
				}
			}

			// Check exile
			if (!card) {
				const exileIndex = state.exile.findIndex(c => c.id === cardId);
				if (exileIndex !== -1) {
					card = state.exile[exileIndex];
					sourceZone = 'exile';
				}
			}

			// Check command zone
			if (!card) {
				const cmdIndex = state.command.findIndex(c => c.id === cardId);
				if (cmdIndex !== -1) {
					card = state.command[cmdIndex];
					sourceZone = 'command';
				}
			}

			if (!card || !sourceZone) {
				console.error('[PlaytestGame] Card not found:', cardId);
				return state;
			}

			// Remove from source zone
			let newState = { ...state };
			if (sourceZone === 'battlefield') {
				newState.battlefield = state.battlefield.filter(c => c.id !== cardId);
			} else if (sourceZone === 'exile') {
				newState.exile = state.exile.filter(c => c.id !== cardId);
			} else if (sourceZone === 'command') {
				newState.command = state.command.filter(c => c.id !== cardId);
			} else if (sourceZone.startsWith('hand:')) {
				const playerId = sourceZone.split(':')[1];
				newState.players = state.players.map(p =>
					p.playerId === playerId
						? { ...p, hand: p.hand.filter(c => c.id !== cardId), handCount: p.hand.length - 1 }
						: p
				);
			} else if (sourceZone.startsWith('library:')) {
				const playerId = sourceZone.split(':')[1];
				newState.players = state.players.map(p =>
					p.playerId === playerId
						? { ...p, library: p.library.filter(c => c.id !== cardId), libraryCount: p.library.length - 1 }
						: p
				);
			} else if (sourceZone.startsWith('graveyard:')) {
				const playerId = sourceZone.split(':')[1];
				newState.players = state.players.map(p =>
					p.playerId === playerId
						? { ...p, graveyard: p.graveyard.filter(c => c.id !== cardId) }
						: p
				);
			}

			// Add to target zone
			card.faceDown = false;
			const upperTargetZone = targetZone.toUpperCase();
			
			if (upperTargetZone === 'BATTLEFIELD') {
				card.zone = 2;
				card.controllerId = state.activeControlSeat;
				newState.battlefield = [...newState.battlefield, card];
			} else if (upperTargetZone === 'GRAVEYARD') {
				card.zone = 3;
				const owner = card.ownerId || state.activeControlSeat;
				newState.players = newState.players.map(p =>
					p.playerId === owner
						? { ...p, graveyard: [...p.graveyard, card] }
						: p
				);
			} else if (upperTargetZone === 'EXILE') {
				card.zone = 4;
				newState.exile = [...newState.exile, card];
			} else if (upperTargetZone === 'HAND') {
				card.zone = 1;
				newState.players = newState.players.map(p =>
					p.playerId === state.activeControlSeat
						? { ...p, hand: [...p.hand, card], handCount: p.hand.length + 1 }
						: p
				);
			} else if (upperTargetZone === 'COMMAND') {
				card.zone = 5;
				card.faceDown = false;
				newState.command = [...newState.command, card];
			} else if (upperTargetZone === 'LIBRARY' || upperTargetZone === 'LIBRARY_TOP') {
				card.zone = 0;
				card.faceDown = true;
				const owner = card.ownerId || state.activeControlSeat;
				newState.players = newState.players.map(p =>
					p.playerId === owner
						? { ...p, library: [card, ...p.library], libraryCount: p.library.length + 1 }
						: p
				);
			} else if (upperTargetZone === 'LIBRARY_BOTTOM') {
				card.zone = 0;
				card.faceDown = true;
				const owner = card.ownerId || state.activeControlSeat;
				newState.players = newState.players.map(p =>
					p.playerId === owner
						? { ...p, library: [...p.library, card], libraryCount: p.library.length + 1 }
						: p
				);
			}

			return newState;
		});
	}

	/**
	 * Tap or untap a card
	 */
	function tapCard(cardId: string, tapped: boolean): void {
		update(state => ({
			...state,
			battlefield: state.battlefield.map(card =>
				card.id === cardId ? { ...card, tapped } : card
			)
		}));
	}

	/**
	 * Untap all permanents controlled by a player
	 */
	function untapAll(playerId: string): void {
		update(state => ({
			...state,
			battlefield: state.battlefield.map(card =>
				card.controllerId === playerId ? { ...card, tapped: false } : card
			)
		}));
	}

	/**
	 * Flip a card face up/down
	 */
	function flipCard(cardId: string, faceDown: boolean): void {
		update(state => ({
			...state,
			battlefield: state.battlefield.map(card =>
				card.id === cardId ? { ...card, faceDown } : card
			)
		}));
	}

	/**
	 * Modify player life
	 */
	function modifyLife(playerId: string, delta: number): void {
		update(state => ({
			...state,
			players: state.players.map(p =>
				p.playerId === playerId
					? { ...p, life: Math.max(0, p.life + delta) }
					: p
			)
		}));
	}

	/**
	 * Set player counter (poison, energy, etc.)
	 */
	function setPlayerCounter(playerId: string, counterType: string, value: number): void {
		update(state => ({
			...state,
			players: state.players.map(p => {
				if (p.playerId !== playerId) return p;
				
				if (counterType === 'poison') {
					return { ...p, poison: Math.max(0, value) };
				} else if (counterType === 'energy') {
					return { ...p, energy: Math.max(0, value) };
				}
				return p;
			})
		}));
	}

	/**
	 * Shuffle a player's library (Fisher-Yates algorithm)
	 */
	function shuffleLibrary(playerId: string): void {
		update(state => {
			const newPlayers = state.players.map(p => {
				if (p.playerId !== playerId) return p;

				const shuffled = [...p.library];
				for (let i = shuffled.length - 1; i > 0; i--) {
					const j = Math.floor(Math.random() * (i + 1));
					[shuffled[i], shuffled[j]] = [shuffled[j], shuffled[i]];
				}

				return { ...p, library: shuffled };
			});

			return { ...state, players: newPlayers };
		});
	}

	/**
	 * Add a card to the visual stack
	 */
	function addToStack(cardId: string): void {
		update(state => {
			const card = state.battlefield.find(c => c.id === cardId) ||
			              state.players.flatMap(p => p.hand).find(c => c.id === cardId);
			
			if (!card) {
				console.error('[PlaytestGame] Card not found for stack:', cardId);
				return state;
			}

			return {
				...state,
				stack: [...state.stack, { ...card }]
			};
		});
	}

	/**
	 * Remove an item from the stack
	 */
	function removeFromStack(itemId: string): void {
		update(state => ({
			...state,
			stack: state.stack.filter(item => item.id !== itemId)
		}));
	}

	/**
	 * Create a token on the battlefield
	 */
	function createToken(
		name: string,
		types: string,
		power: string,
		toughness: string,
		color: string
	): void {
		update(state => {
			const tokenId = `token-${Date.now()}-${Math.random().toString(36).substr(2, 9)}`;
			const token: CardView = {
				id: tokenId,
				name,
				displayName: name,
				subTypes: '',
				superTypes: '',
				color,
				type: types,
				power,
				toughness,
				loyalty: '',
				manaCost: '',
				cardNumber: 0,
				expansionSetCode: '',
				rarity: '',
				rulesText: '',
				abilities: [],
				zone: 2, // BATTLEFIELD
				ownerId: state.activeControlSeat,
				controllerId: state.activeControlSeat,
				tapped: false,
				flipped: false,
				transformed: false,
				faceDown: false,
				counters: [],
				attachedTo: [],
				summoningSickness: true,
				availableActions: []
			};

			return {
				...state,
				battlefield: [...state.battlefield, token]
			};
		});
	}

	/**
	 * Next turn
	 */
	function nextTurn(): void {
		update(state => {
			// Find next player in turn order
			const currentIndex = state.players.findIndex(p => p.playerId === state.activePlayerId);
			const nextIndex = (currentIndex + 1) % state.players.length;
			const nextPlayer = state.players[nextIndex];

			return {
				...state,
				turn: state.turn + 1,
				activePlayerId: nextPlayer.playerId,
				activeControlSeat: nextPlayer.playerId
			};
		});
	}

	/**
	 * Mulligan for a player
	 */
	function mulligan(playerId: string): void {
		update(state => {
			const playerIndex = state.players.findIndex(p => p.playerId === playerId);
			if (playerIndex === -1) return state;

			const player = state.players[playerIndex];
			
			// Return hand to library
			const returnedCards = player.hand.map(card => ({
				...card,
				zone: 0,
				faceDown: true
			}));
			
			const newLibrary = [...returnedCards, ...player.library];
			
			// Shuffle
			for (let i = newLibrary.length - 1; i > 0; i--) {
				const j = Math.floor(Math.random() * (i + 1));
				[newLibrary[i], newLibrary[j]] = [newLibrary[j], newLibrary[i]];
			}
			
			// Draw one less card
			const newHandSize = Math.max(0, player.handCount - 1);
			const newHand = newLibrary.splice(0, newHandSize).map(card => ({
				...card,
				zone: 1,
				faceDown: false
			}));

			const newPlayers = [...state.players];
			newPlayers[playerIndex] = {
				...player,
				hand: newHand,
				handCount: newHand.length,
				library: newLibrary,
				libraryCount: newLibrary.length,
				keptHand: false
			};

			return { ...state, players: newPlayers };
		});
	}

	/**
	 * Keep hand (no mulligan)
	 */
	function keepHand(playerId: string): void {
		update(state => ({
			...state,
			players: state.players.map(p =>
				p.playerId === playerId ? { ...p, keptHand: true } : p
			)
		}));
	}

	/**
	 * Reset to initial state
	 */
	function reset(): void {
		set(initialState);
		clearPersistedPlaytestState();
	}

	return {
		subscribe,
		initialize,
		setCommand,
		switchControlSeat,
		drawCards,
		playCard,
		moveCardToZone,
		tapCard,
		untapAll,
		flipCard,
		modifyLife,
		setPlayerCounter,
		shuffleLibrary,
		addToStack,
		removeFromStack,
		createToken,
		nextTurn,
		mulligan,
		keepHand,
		reset,
		clearPersisted: clearPersistedPlaytestState
	};
}

/**
 * Global playtest game store instance
 */
export const playtestGameStore = createPlaytestGameStore();

// Derived stores for convenient access

export const playtestPlayers = derived(playtestGameStore, $game => $game.players);

export const playtestLocalPlayer = derived(playtestGameStore, $game => {
	return $game.players.find(p => p.playerId === $game.activeControlSeat) || null;
});

export const playtestOpponents = derived(playtestGameStore, $game => {
	return $game.players.filter(p => p.playerId !== $game.activeControlSeat);
});

export const playtestBattlefield = derived(playtestGameStore, $game => $game.battlefield);

export const playtestExile = derived(playtestGameStore, $game => $game.exile);

export const playtestStack = derived(playtestGameStore, $game => $game.stack);

export const playtestActiveControlSeat = derived(playtestGameStore, $game => $game.activeControlSeat);

export const playtestIsInitialized = derived(playtestGameStore, $game => $game.isInitialized);
