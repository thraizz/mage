/**
 * Playtest Game Store
 *
 * Client-side only game state management for playtest mode.
 * No server communication - all state is local.
 */

import { writable, derived } from 'svelte/store';
import { browser } from '$app/environment';
import type { CardView, ManaPoolView } from '$lib/generated/mage/v1/models';
import { ZoneId } from '$lib/utils/zones';
import {
	updatePlayer,
	findCardInState,
	removeCardFromZone,
	addCardToZone,
	shuffleArray,
	getNextPlayer,
	updateCardInZone
} from '$lib/utils/playtest-helpers';

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
	mulliganCount: number;
	revealedTopCard: boolean; // When true, top card of library is permanently visible
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
	log: PlaytestLogEntry[];
	mulliganType: 'london';
	freeMulligans: number;
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
	isInitialized: false,
	log: [],
	mulliganType: 'london',
	freeMulligans: 0
};

// Store up to the last 10 playtest sessions
const PLAYTEST_SESSIONS_STORAGE_KEY = 'mage.playtest.sessions.v1';
const PLAYTEST_ACTIVE_SESSION_ID_KEY = 'mage.playtest.activeSessionId.v1';
const PLAYTEST_STORAGE_VERSION = 1;

// Legacy (single-session) key migration:
const PLAYTEST_LEGACY_STORAGE_KEY = 'mage.playtest.state.v1';

export type PlaytestSessionMeta = {
	id: string;
	createdAt: number;
	savedAt: number;
	label: string;
	playerCount: number;
	turn: number;
};

type PersistedPlaytestStateV1 = {
	version: number;
	savedAt: number;
	state: PlaytestGameState;
};

type PersistedPlaytestSession = {
	id: string;
	createdAt: number;
	savedAt: number;
	label: string;
	state: PlaytestGameState;
};

type PersistedPlaytestSessionsPayload = {
	version: number;
	sessions: PersistedPlaytestSession[];
};

/**
 * Game log (for playtest analysis)
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

const PLAYTEST_LOG_MAX_ENTRIES = 1000;

function makeLogId(): string {
	return `log-${Date.now()}-${Math.random().toString(36).slice(2, 10)}`;
}

function formatLogEntry(e: PlaytestLogEntry, state: PlaytestGameState): string {
	const t = new Date(e.at);
	const hh = t.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit', second: '2-digit' });

	const activeName =
		state.players.find((p) => p.playerId === e.activePlayerId)?.name ?? e.activePlayerId;
	const controlName =
		state.players.find((p) => p.playerId === e.controlSeat)?.name ?? e.controlSeat;

	return `[${hh}] Turn ${e.turn} · Active: ${activeName} · Controlling: ${controlName} · ${e.message}`;
}

function appendLogToState(
	state: PlaytestGameState,
	entry: Omit<PlaytestLogEntry, 'id' | 'at' | 'turn' | 'activePlayerId' | 'controlSeat'>
): PlaytestGameState {
	// Avoid logging if the session isn't initialized yet.
	if (!state.isInitialized) return state;

	const full: PlaytestLogEntry = {
		id: makeLogId(),
		at: Date.now(),
		turn: state.turn,
		activePlayerId: state.activePlayerId,
		controlSeat: state.activeControlSeat,
		...entry
	};

	const next = [...(state.log ?? []), full];
	const trimmed =
		next.length > PLAYTEST_LOG_MAX_ENTRIES
			? next.slice(next.length - PLAYTEST_LOG_MAX_ENTRIES)
			: next;

	return { ...state, log: trimmed };
}

function playerName(state: PlaytestGameState, playerId: string): string {
	return state.players.find((p) => p.playerId === playerId)?.name ?? playerId;
}

function isValidPlaytestState(state: PlaytestGameState): boolean {
	if (!state || typeof state !== 'object') return false;
	if (state.isInitialized !== true) return false;
	if (!Array.isArray(state.players) || state.players.length === 0) return false;
	if (typeof state.gameId !== 'string' || state.gameId.length === 0) return false;
	return true;
}

function buildSessionLabel(state: PlaytestGameState): string {
	const names = (state.players || []).map((p) => p.name).filter(Boolean);
	if (names.length === 0) return 'Playtest session';
	return names.join(' vs ');
}

function readSessionsPayload(): PersistedPlaytestSessionsPayload | null {
	if (!browser) return null;
	try {
		const raw = localStorage.getItem(PLAYTEST_SESSIONS_STORAGE_KEY);
		if (!raw) return null;
		const parsed = JSON.parse(raw) as PersistedPlaytestSessionsPayload;
		if (!parsed || typeof parsed !== 'object') return null;
		if (parsed.version !== PLAYTEST_STORAGE_VERSION) return null;
		if (!Array.isArray(parsed.sessions)) return null;
		return parsed;
	} catch (err) {
		console.warn('[PlaytestGame] Failed to load persisted sessions:', err);
		return null;
	}
}

function writeSessionsPayload(payload: PersistedPlaytestSessionsPayload): void {
	if (!browser) return;
	try {
		localStorage.setItem(PLAYTEST_SESSIONS_STORAGE_KEY, JSON.stringify(payload));
	} catch (err) {
		console.warn('[PlaytestGame] Failed to persist sessions:', err);
	}
}

function readActiveSessionId(): string | null {
	if (!browser) return null;
	try {
		return localStorage.getItem(PLAYTEST_ACTIVE_SESSION_ID_KEY);
	} catch (err) {
		console.warn('[PlaytestGame] Failed to read active session id:', err);
		return null;
	}
}

function writeActiveSessionId(sessionId: string | null): void {
	if (!browser) return;
	try {
		if (!sessionId) {
			localStorage.removeItem(PLAYTEST_ACTIVE_SESSION_ID_KEY);
			return;
		}
		localStorage.setItem(PLAYTEST_ACTIVE_SESSION_ID_KEY, sessionId);
	} catch (err) {
		console.warn('[PlaytestGame] Failed to write active session id:', err);
	}
}

function loadPersistedSessions(): PersistedPlaytestSession[] {
	const payload = readSessionsPayload();
	const sessions = payload?.sessions ?? [];
	return sessions
		.filter(
			(s) =>
				s &&
				typeof s === 'object' &&
				typeof s.id === 'string' &&
				s.id.length > 0 &&
				isValidPlaytestState(s.state)
		)
		.sort((a, b) => (b.savedAt ?? 0) - (a.savedAt ?? 0))
		.slice(0, 10);
}

function savePersistedSessions(sessions: PersistedPlaytestSession[]): void {
	writeSessionsPayload({
		version: PLAYTEST_STORAGE_VERSION,
		sessions: sessions.slice(0, 10)
	});
}

function migrateLegacySingleSessionIfNeeded(): void {
	if (!browser) return;
	try {
		const legacyRaw = localStorage.getItem(PLAYTEST_LEGACY_STORAGE_KEY);
		if (!legacyRaw) return;

		const legacyParsed = JSON.parse(legacyRaw) as PersistedPlaytestStateV1;
		if (!legacyParsed || legacyParsed.version !== PLAYTEST_STORAGE_VERSION) {
			localStorage.removeItem(PLAYTEST_LEGACY_STORAGE_KEY);
			return;
		}
		if (!legacyParsed.state || !isValidPlaytestState(legacyParsed.state)) {
			localStorage.removeItem(PLAYTEST_LEGACY_STORAGE_KEY);
			return;
		}

		const state = legacyParsed.state;
		const id = state.gameId || `playtest-${legacyParsed.savedAt || Date.now()}`;
		const now = Date.now();
		const createdAt = legacyParsed.savedAt || now;
		const savedAt = legacyParsed.savedAt || now;
		const label = buildSessionLabel(state);

		const existing = loadPersistedSessions();
		const without = existing.filter((s) => s.id !== id);
		const migrated: PersistedPlaytestSession = { id, createdAt, savedAt, label, state };
		savePersistedSessions([migrated, ...without]);
		writeActiveSessionId(id);

		localStorage.removeItem(PLAYTEST_LEGACY_STORAGE_KEY);
	} catch (err) {
		console.warn('[PlaytestGame] Failed to migrate legacy playtest state:', err);
	}
}

function upsertSessionFromState(state: PlaytestGameState, opts?: { bumpSavedAt?: boolean }): void {
	if (!browser) return;
	if (!isValidPlaytestState(state)) return;

	const now = Date.now();
	const id = state.gameId;

	const sessions = loadPersistedSessions();
	const existing = sessions.find((s) => s.id === id);

	const createdAt = existing?.createdAt ?? now;
	const savedAt = opts?.bumpSavedAt === false ? (existing?.savedAt ?? now) : now;
	const label = buildSessionLabel(state);

	const updated: PersistedPlaytestSession = {
		id,
		createdAt,
		savedAt,
		label,
		state
	};

	savePersistedSessions([updated, ...sessions.filter((s) => s.id !== id)]);
	writeActiveSessionId(id);
}

export function getPlaytestSessionsMeta(): PlaytestSessionMeta[] {
	if (!browser) return [];
	migrateLegacySingleSessionIfNeeded();
	const sessions = loadPersistedSessions();
	return sessions.map((s) => ({
		id: s.id,
		createdAt: s.createdAt,
		savedAt: s.savedAt,
		label: s.label,
		playerCount: s.state.players.length,
		turn: s.state.turn
	}));
}

export function deletePlaytestSession(sessionId: string): void {
	if (!browser) return;
	migrateLegacySingleSessionIfNeeded();
	const sessions = loadPersistedSessions().filter((s) => s.id !== sessionId);
	savePersistedSessions(sessions);

	const active = readActiveSessionId();
	if (active === sessionId) {
		writeActiveSessionId(sessions[0]?.id ?? null);
	}
}

export function clearPlaytestSessions(): void {
	if (!browser) return;
	try {
		localStorage.removeItem(PLAYTEST_SESSIONS_STORAGE_KEY);
		localStorage.removeItem(PLAYTEST_ACTIVE_SESSION_ID_KEY);
	} catch (err) {
		console.warn('[PlaytestGame] Failed to clear playtest sessions:', err);
	}
}

function loadPersistedPlaytestState(): PlaytestGameState | null {
	if (!browser) return null;

	migrateLegacySingleSessionIfNeeded();

	const activeId = readActiveSessionId();
	const sessions = loadPersistedSessions();
	const active = (activeId ? sessions.find((s) => s.id === activeId) : null) ?? sessions[0];

	if (!active || !isValidPlaytestState(active.state)) return null;
	return active.state;
}

/**
 * Create playtest game store
 */
function createPlaytestGameStore() {
	function clearLog(): void {
		update((state) => ({ ...state, log: [] }));
	}

	function addLog(message: string, kind: string = 'note'): void {
		update((state) => appendLogToState(state, { kind, message }));
	}

	function buildLogText(state: PlaytestGameState): string {
		const entries = Array.isArray(state.log) ? state.log : [];
		return entries.map((e) => formatLogEntry(e, state)).join('\n');
	}

	const hydrated = loadPersistedPlaytestState();
	const normalizedHydrated = hydrated
		? { ...hydrated, log: Array.isArray(hydrated.log) ? hydrated.log : [] }
		: null;

	const { subscribe, set, update } = writable<PlaytestGameState>(
		normalizedHydrated ?? initialState
	);

	// Persist any meaningful state changes (client-only). Debounced to avoid excessive writes.
	let persistTimer: ReturnType<typeof setTimeout> | null = null;
	subscribe((state) => {
		if (!browser) return;
		if (!isValidPlaytestState(state)) return;
		if (persistTimer) clearTimeout(persistTimer);
		persistTimer = setTimeout(() => {
			upsertSessionFromState(state, { bumpSavedAt: true });
			persistTimer = null;
		}, 200);
	});

	/**
	 * Initialize game state with players and their decks
	 */
	function initialize(
		gameId: string,
		players: PlaytestPlayer[],
		options?: {
			mulliganType?: 'london';
			freeMulligans?: number;
		}
	): void {
		if (players.length === 0) {
			console.error('[PlaytestGame] Cannot initialize with no players');
			return;
		}

		const nextState: PlaytestGameState = {
			gameId,
			activeControlSeat: players[0].playerId,
			players,
			battlefield: [],
			exile: [],
			stack: [],
			command: [],
			turn: 1,
			activePlayerId: players[0].playerId,
			isInitialized: true,
			log: [
				{
					id: makeLogId(),
					at: Date.now(),
					turn: 1,
					activePlayerId: players[0].playerId,
					controlSeat: players[0].playerId,
					kind: 'init',
					message: `Game initialized (${
						players
							.map((p) => p.name)
							.filter(Boolean)
							.join(' vs ') || 'Playtest'
					})`
				}
			],
			mulliganType: options?.mulliganType ?? 'london',
			freeMulligans: options?.freeMulligans ?? 0
		};

		set(nextState);
		// Persist immediately as a new/active session.
		upsertSessionFromState(nextState, { bumpSavedAt: true });

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
		update((state) => ({
			...state,
			activeControlSeat: playerId
		}));
	}

	/**
	 * Draw cards for a player
	 */
	function drawCards(playerId: string, count: number): void {
		update((state) => {
			const playerIndex = state.players.findIndex((p) => p.playerId === playerId);
			if (playerIndex === -1) {
				console.error('[PlaytestGame] Player not found:', playerId);
				return state;
			}

			const player = state.players[playerIndex];
			const before = player.library.length;
			const drawn = player.library.splice(0, Math.min(count, player.library.length));

			// Update zone and make cards visible
			drawn.forEach((card) => {
				card.zone = ZoneId.HAND;
				card.faceDown = false;
			});

			const newPlayers = [...state.players];
			newPlayers[playerIndex] = {
				...player,
				hand: [...player.hand, ...drawn],
				handCount: player.hand.length + drawn.length,
				libraryCount: player.library.length
			};

			const msg = `${playerName(state, playerId)} draws ${drawn.length}${count !== drawn.length ? ` (requested ${count})` : ''}. Library: ${before} → ${player.library.length}.`;

			const next = {
				...state,
				players: newPlayers
			};

			return appendLogToState(next, { kind: 'draw', message: msg });
		});
	}

	/**
	 * Play a card from hand to battlefield
	 */
	function playCard(cardId: string, tapped: boolean = false): void {
		update((state) => {
			const controllingPlayer = state.players.find((p) => p.playerId === state.activeControlSeat);
			if (!controllingPlayer) return state;

			const cardIndex = controllingPlayer.hand.findIndex((c) => c.id === cardId);
			if (cardIndex === -1) {
				console.error('[PlaytestGame] Card not found in hand:', cardId);
				return state;
			}

			const card = controllingPlayer.hand[cardIndex];
			const newHand = [...controllingPlayer.hand];
			newHand.splice(cardIndex, 1);

			// Update card properties
			card.zone = ZoneId.BATTLEFIELD;
			card.controllerId = state.activeControlSeat;
			card.tapped = tapped;
			card.faceDown = false;

			// Update player
			const newPlayers = state.players.map((p) =>
				p.playerId === state.activeControlSeat
					? { ...p, hand: newHand, handCount: newHand.length }
					: p
			);

			const msg = `${playerName(state, state.activeControlSeat)} plays ${card.name}${tapped ? ' (tapped)' : ''}. Hand: ${controllingPlayer.hand.length} → ${newHand.length}.`;

			const next = {
				...state,
				players: newPlayers,
				battlefield: [...state.battlefield, card]
			};
			return appendLogToState(next, { kind: 'play', message: msg });
		});
	}

	/**
	 * Move a card to a different zone
	 */
	function moveCardToZone(cardId: string, targetZone: string): void {
		update((state) => {
			// Find the card in any zone
			const found = findCardInState(state, cardId);
			if (!found) {
				console.error('[PlaytestGame] Card not found:', cardId);
				return state;
			}

			console.log('[PlaytestGame] Found card:', found);
			const { card, sourceZone } = found;
			if (!card) {
				console.error('[PlaytestGame] Card not found:', cardId);
				return state;
			}

			// MTG Rule: Tokens cease to exist when they leave the battlefield
			if (cardId.startsWith('token-') && sourceZone === 'battlefield') {
				// Token is leaving the battlefield - it ceases to exist
				const next = removeCardFromZone(state, cardId, sourceZone);
				const msg = `${playerName(state, state.activeControlSeat)} moves ${card.name} to ${targetZone}. Token ceases to exist.`;
				return appendLogToState(next, { kind: 'move', message: msg });
			}

			// Remove from source zone
			let next = removeCardFromZone(state, cardId, sourceZone);

			// Add to target zone
			next = addCardToZone(next, card, targetZone, state.activeControlSeat);

			const msg = `${playerName(state, state.activeControlSeat)} moves ${card.name} from ${sourceZone} to ${targetZone}.`;
			return appendLogToState(next, { kind: 'move', message: msg });
		});
	}

	/**
	 * Tap or untap a card
	 */
	function tapCard(cardId: string, tapped: boolean): void {
		// update((state) => ({
		// 	...state,
		// 	battlefield: state.battlefield.map((card) =>
		// 		card.id === cardId ? { ...card, tapped } : card
		// 	)
		// }));
		update((state) => {
			const card = state.battlefield.find((c) => c.id === cardId);
			if (!card) return state;

			const msg = `${playerName(state, state.activeControlSeat)} ${tapped ? 'taps' : 'untaps'} ${card.name}.`;
			const next = {
				...state,
				battlefield: state.battlefield.map((c) => (c.id === cardId ? { ...c, tapped } : c))
			};
			return appendLogToState(next, { kind: 'tap', message: msg });
		});
	}

	/**
	 * Untap all permanents controlled by a player
	 */
	function untapAll(playerId: string): void {
		update((state) => {
			const msg = `${playerName(state, state.activeControlSeat)} untaps all permanents.`;
			const next = {
				...state,
				battlefield: state.battlefield.map((c) =>
					c.controllerId === playerId ? { ...c, tapped: false } : c
				)
			};
			return appendLogToState(next, { kind: 'untap', message: msg });
		});
	}

	/**
	 * Flip a card face up/down
	 */
	function flipCard(cardId: string, faceDown: boolean): void {
		update((state) => {
			const card = state.battlefield.find((c) => c.id === cardId);
			if (!card) return state;

			const msg = `${playerName(state, state.activeControlSeat)} flips ${card.name} face ${faceDown ? 'down' : 'up'}.`;
			const next = {
				...state,
				battlefield: state.battlefield.map((c) => (c.id === cardId ? { ...c, faceDown } : c))
			};
			return appendLogToState(next, { kind: 'flip', message: msg });
		});
	}

	/**
	 * Modify player life
	 */
	function modifyLife(playerId: string, delta: number): void {
		update((state) => {
			const msg = `${playerName(state, state.activeControlSeat)} modifies life by ${delta}.`;
			const next = {
				...state,
				players: updatePlayer(state.players, playerId, (p) => ({
					life: Math.max(0, p.life + delta)
				}))
			};
			return appendLogToState(next, { kind: 'life', message: msg });
		});
	}

	/**
	 * Set player counter (poison, energy, etc.)
	 */
	function setPlayerCounter(playerId: string, counterType: string, value: number): void {
		update((state) => {
			const updates: Partial<PlaytestPlayer> = {};
			if (counterType === 'poison') {
				updates.poison = Math.max(0, value);
			} else if (counterType === 'energy') {
				updates.energy = Math.max(0, value);
			}

			const msg = `${playerName(state, state.activeControlSeat)} sets ${counterType} to ${value}.`;
			const next = {
				...state,
				players: updatePlayer(state.players, playerId, () => updates)
			};
			return appendLogToState(next, { kind: 'counter', message: msg });
		});
	}

	/**
	 * Shuffle a player's library (Fisher-Yates algorithm)
	 */
	function shuffleLibrary(playerId: string): void {
		update((state) => {
			const msg = `${playerName(state, state.activeControlSeat)} shuffles their library.`;
			const next = {
				...state,
				players: updatePlayer(state.players, playerId, (p) => ({
					library: shuffleArray(p.library),
					revealedTopCard: false // Clear revealed top when shuffling
				}))
			};
			return appendLogToState(next, { kind: 'shuffle', message: msg });
		});
	}

	/**
	 * Add a card to the visual stack
	 */
	function addToStack(cardId: string): void {
		update((state) => {
			const card =
				state.battlefield.find((c) => c.id === cardId) ||
				state.players.flatMap((p) => p.hand).find((c) => c.id === cardId);

			if (!card) {
				console.error('[PlaytestGame] Card not found for stack:', cardId);
				return state;
			}

			const msg = `${playerName(state, state.activeControlSeat)} adds ${card.name} to the stack.`;
			const next = {
				...state,
				stack: [...state.stack, { ...card }]
			};
			return appendLogToState(next, { kind: 'add', message: msg });
		});
	}

	/**
	 * Remove an item from the stack
	 */
	function removeFromStack(itemId: string): void {
		update((state) => {
			const msg = `${playerName(state, state.activeControlSeat)} removes ${itemId} from the stack.`;
			const next = {
				...state,
				stack: state.stack.filter((item) => item.id !== itemId)
			};
			return appendLogToState(next, { kind: 'remove', message: msg });
		});
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
		update((state) => {
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
				zone: ZoneId.BATTLEFIELD,
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

			const msg = `${playerName(state, state.activeControlSeat)} creates ${token.name} token.`;
			const next = {
				...state,
				battlefield: [...state.battlefield, token]
			};
			return appendLogToState(next, { kind: 'create', message: msg });
		});
	}

	/**
	 * Add counters to a card
	 */
	function addCounter(cardId: string, counterName: string, amount: number = 1): void {
		console.log('[addCounter] Called with:', { cardId, counterName, amount });
		update((state) => {
			// Find the card in any zone
			const found = findCardInState(state, cardId);
			if (!found) {
				console.log('[addCounter] Card not found:', cardId);
				return state;
			}

			const { card, sourceZone } = found;
			console.log('[addCounter] Found card:', card.name, 'Current counters:', card.counters);

			// Create new card with updated counters
			const existingCounter = card.counters.find((c) => c.name === counterName);
			const newCounters = existingCounter
				? card.counters.map((c) => (c.name === counterName ? { ...c, count: c.count + amount } : c))
				: [...card.counters, { name: counterName, count: amount }];

			console.log('[addCounter] New counters:', newCounters);

			// Update card in zone with new card object
			const updatedState = updateCardInZone(state, cardId, sourceZone, (c) => ({
				...c,
				counters: newCounters
			}));

			const msg = `Added ${amount} ${counterName} counter(s) to ${card.name}.`;
			return appendLogToState(updatedState, { kind: 'counter', message: msg });
		});
	}

	/**
	 * Remove counters from a card
	 */
	function removeCounter(cardId: string, counterName: string, amount: number = 1): void {
		console.log('[removeCounter] Called with:', { cardId, counterName, amount });
		update((state) => {
			const found = findCardInState(state, cardId);
			if (!found) {
				console.log('[removeCounter] Card not found:', cardId);
				return state;
			}

			const { card, sourceZone } = found;
			console.log('[removeCounter] Found card:', card.name, 'Current counters:', card.counters);

			const counter = card.counters.find((c) => c.name === counterName);
			if (!counter) {
				console.log('[removeCounter] Counter not found:', counterName);
				return state;
			}

			const newCount = Math.max(0, counter.count - amount);

			// Create new counters array - remove if 0, otherwise update count
			const newCounters =
				newCount === 0
					? card.counters.filter((c) => c.name !== counterName)
					: card.counters.map((c) => (c.name === counterName ? { ...c, count: newCount } : c));

			console.log('[removeCounter] New counters:', newCounters);

			// Update card in zone with new card object
			const updatedState = updateCardInZone(state, cardId, sourceZone, (c) => ({
				...c,
				counters: newCounters
			}));

			const msg = `Removed ${amount} ${counterName} counter(s) from ${card.name}.`;
			return appendLogToState(updatedState, { kind: 'counter', message: msg });
		});
	}

	/**
	 * Set a counter to a specific value
	 */
	function setCounter(cardId: string, counterName: string, amount: number): void {
		console.log('[setCounter] Called with:', { cardId, counterName, amount });
		update((state) => {
			const found = findCardInState(state, cardId);
			if (!found) {
				console.log('[setCounter] Card not found:', cardId);
				return state;
			}

			const { card, sourceZone } = found;
			console.log('[setCounter] Found card:', card.name, 'Current counters:', card.counters);

			// Create new counters array
			const newCounters =
				amount <= 0
					? card.counters.filter((c) => c.name !== counterName)
					: card.counters.find((c) => c.name === counterName)
						? card.counters.map((c) => (c.name === counterName ? { ...c, count: amount } : c))
						: [...card.counters, { name: counterName, count: amount }];

			console.log('[setCounter] New counters:', newCounters);

			// Update card in zone with new card object
			const updatedState = updateCardInZone(state, cardId, sourceZone, (c) => ({
				...c,
				counters: newCounters
			}));

			const msg = `Set ${counterName} counters on ${card.name} to ${amount}.`;
			return appendLogToState(updatedState, { kind: 'counter', message: msg });
		});
	}


	/**
	 * Mill cards (move top N cards from library to graveyard)
	 */
	function millCards(playerId: string, count: number): void {
		update((state) => {
			const playerIndex = state.players.findIndex((p) => p.playerId === playerId);
			if (playerIndex === -1) {
				console.error('[PlaytestGame] Player not found:', playerId);
				return state;
			}

			const player = state.players[playerIndex];
			const actualCount = Math.min(count, player.library.length);
			const milled = player.library.splice(0, actualCount);

			// Update zone for milled cards
			milled.forEach((card) => {
				card.zone = ZoneId.GRAVEYARD;
				card.faceDown = false;
			});

			const newPlayers = [...state.players];
			newPlayers[playerIndex] = {
				...player,
				graveyard: [...player.graveyard, ...milled],
				libraryCount: player.library.length
			};

			const msg = `${playerName(state, playerId)} mills ${actualCount} card(s). Library: ${player.library.length + actualCount} → ${player.library.length}.`;

			const next = {
				...state,
				players: newPlayers
			};

			return appendLogToState(next, { kind: 'mill', message: msg });
		});
	}

	/**
	 * Reveal top N cards (for temporary view, like Reveal the First Card)
	 */
	function revealTopCards(playerId: string, count: number): CardView[] {
		let revealed: CardView[] = [];
		update((state) => {
			const player = state.players.find((p) => p.playerId === playerId);
			if (!player) {
				console.error('[PlaytestGame] Player not found:', playerId);
				return state;
			}

			const actualCount = Math.min(count, player.library.length);
			revealed = player.library.slice(0, actualCount);

			const msg = `${playerName(state, playerId)} reveals top ${actualCount} card(s) of their library.`;
			return appendLogToState(state, { kind: 'reveal', message: msg });
		});
		return revealed;
	}

	/**
	 * Start a scry session (extract top N cards for scry decision)
	 */
	function scryCards(playerId: string, count: number): ScrySession | null {
		let session: ScrySession | null = null;
		update((state) => {
			const player = state.players.find((p) => p.playerId === playerId);
			if (!player) {
				console.error('[PlaytestGame] Player not found:', playerId);
				return state;
			}

			const actualCount = Math.min(count, player.library.length);
			if (actualCount === 0) {
				console.warn('[PlaytestGame] No cards to scry');
				return state;
			}

			const cards = player.library.slice(0, actualCount);
			const sessionId = `scry-${Date.now()}-${Math.random().toString(36).substr(2, 9)}`;

			session = {
				sessionId,
				playerId,
				cards
			};

			const msg = `${playerName(state, playerId)} scries ${actualCount}.`;
			return appendLogToState(state, { kind: 'scry', message: msg });
		});
		return session;
	}

	/**
	 * Apply scry decision (reorder library based on user choices)
	 */
	function applyScryDecision(
		playerId: string,
		scryCount: number,
		keepOnTop: CardView[],
		putToBottom: CardView[]
	): void {
		update((state) => {
			const playerIndex = state.players.findIndex((p) => p.playerId === playerId);
			if (playerIndex === -1) {
				console.error('[PlaytestGame] Player not found:', playerId);
				return state;
			}

			const player = state.players[playerIndex];

			// Remove the scried cards from the library
			const scryCardIds = new Set([...keepOnTop, ...putToBottom].map((c) => c.id));
			const remaining = player.library.filter((c) => !scryCardIds.has(c.id));

			// Rebuild library: keep on top, remaining cards, put to bottom
			const newLibrary = [...keepOnTop, ...remaining, ...putToBottom];

			const newPlayers = [...state.players];
			newPlayers[playerIndex] = {
				...player,
				library: newLibrary
			};

			const msg = `${playerName(state, playerId)} completes scry ${scryCount} (${keepOnTop.length} kept on top, ${putToBottom.length} to bottom).`;

			const next = {
				...state,
				players: newPlayers
			};

			return appendLogToState(next, { kind: 'scry', message: msg });
		});
	}

	/**
	 * Set revealed top card state (like Courser of Kruphix)
	 */
	function setRevealedTop(playerId: string, revealed: boolean): void {
		update((state) => {
			const msg = revealed
				? `${playerName(state, playerId)} reveals the top card of their library permanently.`
				: `${playerName(state, playerId)} hides the revealed top card.`;

			const next = {
				...state,
				players: updatePlayer(state.players, playerId, () => ({ revealedTopCard: revealed }))
			};

			return appendLogToState(next, { kind: 'reveal', message: msg });
		});
	}

	/**
	 * Next turn
	 */
	function nextTurn(): void {
		update((state) => {
			const nextPlayer = getNextPlayer(state.players, state.activePlayerId);
			if (!nextPlayer) return state;

			const msg = `${playerName(state, state.activeControlSeat)} ends their turn.`;
			const next = {
				...state,
				turn: state.turn + 1,
				activePlayerId: nextPlayer.playerId,
				activeControlSeat: nextPlayer.playerId
			};
			return appendLogToState(next, { kind: 'endTurn', message: msg });
		});
	}

	/**
	 * Mulligan for a player
	 */
	function mulligan(playerId: string): void {
		update((state) => {
			const playerIndex = state.players.findIndex((p) => p.playerId === playerId);
			if (playerIndex === -1) return state;

			const player = state.players[playerIndex];
			const newMulliganCount = player.mulliganCount + 1;

			// Return hand to library
			const returnedCards = player.hand.map((card) => ({
				...card,
				zone: ZoneId.LIBRARY,
				faceDown: true
			}));

			// Shuffle library with returned cards
			const newLibrary = shuffleArray([...returnedCards, ...player.library]);

			// Calculate new hand size based on free mulligans
			// If mulliganCount < freeMulligans, draw 7 cards (no penalty)
			// Otherwise, draw 7 - (mulliganCount - freeMulligans + 1) cards
			let newHandSize: number;
			if (newMulliganCount <= state.freeMulligans) {
				newHandSize = 7; // Free mulligan - draw full 7
			} else {
				const penaltyMulligans = newMulliganCount - state.freeMulligans;
				newHandSize = Math.max(0, 7 - penaltyMulligans);
			}

			const newHand = newLibrary.splice(0, newHandSize).map((card) => ({
				...card,
				zone: ZoneId.HAND,
				faceDown: false
			}));

			const newPlayers = [...state.players];
			newPlayers[playerIndex] = {
				...player,
				hand: newHand,
				handCount: newHand.length,
				library: newLibrary,
				libraryCount: newLibrary.length,
				keptHand: false,
				mulliganCount: newMulliganCount
			};

			const msg = `${playerName(state, state.activeControlSeat)} mulls their hand.`;
			const next = { ...state, players: newPlayers };

			return appendLogToState(next, { kind: 'mulligan', message: msg });
		});
	}

	/**
	 * Keep hand (no mulligan)
	 */
	function keepHand(playerId: string): void {
		update((state) => {
			const msg = `${playerName(state, state.activeControlSeat)} keeps their hand.`;
			const next = {
				...state,
				players: updatePlayer(state.players, playerId, () => ({ keptHand: true }))
			};

			return appendLogToState(next, { kind: 'keepHand', message: msg });
		});
	}

	/**
	 * Reset to initial state
	 */
	function reset(): void {
		set(initialState);
		// Do not delete session history; just drop active in-memory session.
		writeActiveSessionId(null);
	}

	function restoreSession(sessionId: string): boolean {
		if (!browser) return false;
		migrateLegacySingleSessionIfNeeded();
		const sessions = loadPersistedSessions();
		const found = sessions.find((s) => s.id === sessionId);
		if (!found) return false;

		// Bump recency and mark active.
		const now = Date.now();
		const updated: PersistedPlaytestSession = { ...found, savedAt: now };
		savePersistedSessions([updated, ...sessions.filter((s) => s.id !== sessionId)]);
		writeActiveSessionId(sessionId);
		set(updated.state);
		return true;
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
		addCounter,
		removeCounter,
		setCounter,
		millCards,
		revealTopCards,
		scryCards,
		applyScryDecision,
		setRevealedTop,
		nextTurn,
		mulligan,
		keepHand,
		reset,
		restoreSession,
		listSessions: (): PlaytestSessionMeta[] => getPlaytestSessionsMeta(),
		deleteSession: (sessionId: string): void => deletePlaytestSession(sessionId),
		clearSessions: (): void => clearPlaytestSessions(),
		// Log API methods
		clearLog,
		addLog,
		buildLogText
	};
}

/**
 * Global playtest game store instance
 */
export const playtestGameStore = createPlaytestGameStore();

// Derived stores for convenient access

export const playtestPlayers = derived(playtestGameStore, ($game) => $game.players);

export const playtestLocalPlayer = derived(playtestGameStore, ($game) => {
	return $game.players.find((p) => p.playerId === $game.activeControlSeat) || null;
});

export const playtestOpponents = derived(playtestGameStore, ($game) => {
	return $game.players.filter((p) => p.playerId !== $game.activeControlSeat);
});

export const playtestBattlefield = derived(playtestGameStore, ($game) => $game.battlefield);

export const playtestExile = derived(playtestGameStore, ($game) => $game.exile);

export const playtestStack = derived(playtestGameStore, ($game) => $game.stack);

export const playtestActiveControlSeat = derived(
	playtestGameStore,
	($game) => $game.activeControlSeat
);

export const playtestIsInitialized = derived(playtestGameStore, ($game) => $game.isInitialized);
