<script lang="ts">
	import { onMount } from 'svelte';
	import { page } from '$app/stores';
	import { goto } from '$app/navigation';
	import { gameStore } from '$lib/stores/game';
	import {
		playtestGameStore,
		playtestPlayers,
		playtestLocalPlayer,
		playtestOpponents,
		playtestBattlefield,
		playtestExile,
		playtestActiveControlSeat,
		playtestIsInitialized,
		type PlaytestSessionMeta
	} from '$lib/stores/playtest-game';
	import { initializePlaytest, validateDeckIds } from '$lib/playtest/initializer';
	import { toast } from '$lib/stores/toast';

	// Game components
	import Card from '$lib/components/game/Card.svelte';
	import PlayerHand from '$lib/components/game/PlayerHand.svelte';
	import Graveyard from '$lib/components/game/Graveyard.svelte';
	import ExileZone from '$lib/components/game/ExileZone.svelte';
	import LibraryZone from '$lib/components/game/LibraryZone.svelte';
	import PlaytestLibrarySearch from '$lib/components/game/PlaytestLibrarySearch.svelte';
	import ManaPool from '$lib/components/game/ManaPool.svelte';
	import TokenCreator from '$lib/components/game/TokenCreator.svelte';
	import CreateTokenDialog from '$lib/components/game/CreateTokenDialog.svelte';
	import CounterDialog from '$lib/components/game/CounterDialog.svelte';
	import MulliganDialog from '$lib/components/game/MulliganDialog.svelte';
	import Keyboard from '@lucide/svelte/icons/keyboard';
	import Clock from '@lucide/svelte/icons/clock';
	import Copy from '@lucide/svelte/icons/copy';
	import Menu from '@lucide/svelte/icons/menu';
	import X from '@lucide/svelte/icons/x';
	import KeyboardShortcutsModal from '$lib/components/game/KeyboardShortcutsModal.svelte';
	import {
		dragDropStore,
		isDragging as isDraggingStore,
		draggedCardName,
		dragPosition,
		isOverValidDropZone,
		currentDropZone,
		getAllValidDropZones,
		type SourceZone
	} from '$lib/utils/drag-drop';
	import { getScryfallImageUrl } from '$lib/utils/scryfall';

	// State
	let loading = $state(true);
	let error = $state<string | null>(null);
	let showTokenCreator = $state(false);
	let showCreateTokenDialog = $state(false);
	let showCounterDialog = $state(false);
	let selectedCardForCounters = $state<{ id: string; name: string } | null>(null);
	let showKeyboardShortcuts = $state(false);
	let showAllHands = $state(false);
	let showMenu = $state(false);
	let hoveredCardId = $state<string | null>(null);
	let showLifeMenu = $state(false);
	let lifeMenuEl: HTMLDivElement | null = $state(null);
	let showDebugOverlay = $state(false);
	let selectedOpponentId = $state<string | null>(null);
	let showOpponentLifeMenu = $state(false);
	let opponentLifeMenuEl: HTMLDivElement | null = $state(null);
	let showDeckSearch = $state(false);

	// Mulligan state
	let mulliganPlayerIndex = $state<number | null>(null);
	let mulliganCount = $state(0);

	// Drag-drop state
	const isDragging = $derived($isDraggingStore);
	const dragCardName = $derived($draggedCardName);
	const dragPos = $derived($dragPosition);
	const isOverValidDrop = $derived($isOverValidDropZone);
	const dropZone = $derived($currentDropZone);

	// Game log
	const gameLog = $derived($playtestGameStore.log || []);

	// Drop zone elements
	let battlefieldDropZoneEl: HTMLDivElement | null = $state(null);
	let graveyardDropZoneEl: HTMLElement | null = $state(null);
	let exileDropZoneEl: HTMLElement | null = $state(null);
	let handDropZoneEl: HTMLElement | null = $state(null);
	let libraryDropZoneEl: HTMLElement | null = $state(null);
	let commandDropZoneEl: HTMLElement | null = $state(null);
	let dropZoneUnregister: (() => void) | null = null;
	let graveyardDropZoneUnregister: (() => void) | null = null;
	let exileDropZoneUnregister: (() => void) | null = null;
	let handDropZoneUnregister: (() => void) | null = null;
	let libraryDropZoneUnregister: (() => void) | null = null;
	let commandDropZoneUnregister: (() => void) | null = null;

	// Battlefield drag state
	let battlefieldDragStartPosition = $state<{ x: number; y: number } | null>(null);
	let battlefieldIsDragPending = $state(false);
	const DRAG_THRESHOLD = 5;

	// Command zone drag state
	let commandDragStartPosition = $state<{ x: number; y: number } | null>(null);
	let commandIsDragPending = $state(false);

	// State for session restore UI
	let showSessionPicker = $state(false);
	let availableSessions = $state<PlaytestSessionMeta[]>([]);

	// Derived state from stores
	const players = $derived($playtestPlayers);
	const me = $derived($playtestLocalPlayer);
	const otherPlayers = $derived($playtestOpponents);
	const battlefield = $derived($playtestBattlefield);
	const exile = $derived($playtestExile);
	const activeControlSeat = $derived($playtestActiveControlSeat);
	const isInitialized = $derived($playtestIsInitialized);

	// Selected opponent (auto-select first opponent if not set)
	const selectedOpponent = $derived(() => {
		if (otherPlayers.length === 0) return null;
		if (!selectedOpponentId || !otherPlayers.find((p) => p.playerId === selectedOpponentId)) {
			// Auto-select first opponent
			return otherPlayers[0];
		}
		return otherPlayers.find((p) => p.playerId === selectedOpponentId) || otherPlayers[0];
	});

	// Split battlefield by controller
	const myBattlefield = $derived(battlefield.filter((c) => c.controllerId === activeControlSeat));
	const opponentBattlefield = $derived(() => {
		const opponent = selectedOpponent();
		return opponent ? battlefield.filter((c) => c.controllerId === opponent.playerId) : [];
	});

	function isLandPermanent(cardType?: string | null): boolean {
		// Mage type strings are typically like: "Land", "Legendary Land", "Artifact Land", etc.
		return !!cardType && /\bland\b/i.test(cardType);
	}

	// Split battlefield rows: nonlands (top) + lands (bottom)
	const myBattlefieldNonlands = $derived(myBattlefield.filter((c) => !isLandPermanent(c.type)));
	const myBattlefieldLands = $derived(myBattlefield.filter((c) => isLandPermanent(c.type)));
	const opponentBattlefieldNonlands = $derived(() =>
		opponentBattlefield().filter((c) => !isLandPermanent(c.type))
	);
	const opponentBattlefieldLands = $derived(() =>
		opponentBattlefield().filter((c) => isLandPermanent(c.type))
	);

	// My cards (from controlling player perspective)
	const myGrave = $derived(me?.graveyard || []);
	const myMana = $derived(
		me?.manaPool || { white: 0, blue: 0, black: 0, red: 0, green: 0, colorless: 0 }
	);

	// Reactive card lookup for counter dialog
	const selectedCardForCountersData = $derived(() => {
		if (!selectedCardForCounters) return null;
		const card =
			$playtestBattlefield.find((c) => c.id === selectedCardForCounters.id) ||
			me?.hand.find((c) => c.id === selectedCardForCounters.id) ||
			me?.graveyard.find((c) => c.id === selectedCardForCounters.id) ||
			me?.exile.find((c) => c.id === selectedCardForCounters.id) ||
			me?.commandZone.find((c) => c.id === selectedCardForCounters.id) ||
			null;
		console.log(
			'[selectedCardForCountersData] Re-evaluated. Card:',
			card?.name,
			'Counters:',
			card?.counters
		);
		return card;
	});

	// Hovered card
	const hoveredCard = $derived(
		hoveredCardId ? battlefield.find((c) => c.id === hoveredCardId) : null
	);

	const activePlayerName = $derived(() => {
		return players.find((p) => p.playerId === $playtestGameStore.activePlayerId)?.name ?? '';
	});

	// Playtest store "turn" currently increments each time we advance priority seat.
	// For display, we want "Turn 2" only when player 1 starts their second turn (i.e. rounds).
	const turnNumber = $derived(() => {
		const step = Math.max(1, $playtestGameStore.turn);
		const n = players.length;
		if (n <= 0) return step;
		return Math.floor((step - 1) / n) + 1;
	});

	// Check if all players have kept hands
	const allPlayersKept = $derived(players.every((p) => p.keptHand));

	// Game started state (for header visibility)
	const isGameStarted = $derived(allPlayersKept && mulliganPlayerIndex === null);

	// Command zone (Commander): show for currently controlled player
	const commandCards = $derived($playtestGameStore.command || []);
	const myCommandCards = $derived(
		commandCards.filter((c) => (c.ownerId || c.controllerId) === activeControlSeat)
	);

	// Opponent command cards
	const opponentCommandCards = $derived(() => {
		const opponent = selectedOpponent();
		return opponent
			? commandCards.filter((c) => (c.ownerId || c.controllerId) === opponent.playerId)
			: [];
	});

	// Check if this is a Commander game (has or had command zone cards)
	const isCommanderGame = $derived(commandCards.length > 0);

	/**
	 * Update the URL with playtestId parameter
	 */
	async function updateUrlWithPlaytestId(playtestId: string): Promise<void> {
		const newSearchParams = new URLSearchParams($page.url.searchParams);
		newSearchParams.set('playtestId', playtestId);
		await goto(`${$page.url.pathname}?${newSearchParams.toString()}`, {
			replaceState: true,
			noScroll: true,
			keepFocus: true
		});
	}

	/**
	 * Initialize playtest from URL params
	 */
	async function initializeFromUrl(): Promise<void> {
		loading = true;
		error = null;

		try {
			const searchParams = $page.url.searchParams;
			const deckIds = validateDeckIds(searchParams);

			// Parse mulligan settings from URL
			const mulliganType = (searchParams.get('mulliganType') as 'london') || 'london';
			const freeMulligans = parseInt(searchParams.get('freeMulligans') || '0', 10);

			console.log('[Playtest] Initializing with deck IDs:', deckIds);
			console.log('[Playtest] Mulligan settings:', { mulliganType, freeMulligans });

			const init = await initializePlaytest(deckIds);
			const initializedPlayers = init.players;
			const gameId = `playtest-${Date.now()}`;

			playtestGameStore.initialize(gameId, initializedPlayers, {
				mulliganType,
				freeMulligans
			});
			playtestGameStore.setCommand(init.command);

			// Initialize the normal game store with playtest data so PlayerHand works
			gameStore.initGame(gameId, initializedPlayers[0].playerId);
			syncPlaytestToGameStore();

			// Set playtestId in URL
			await updateUrlWithPlaytestId(gameId);

			// Start mulligan phase for first player
			mulliganPlayerIndex = 0;

			loading = false;
		} catch (err) {
			console.error('[Playtest] Initialization failed:', err);
			error = err instanceof Error ? err.message : 'Failed to initialize playtest';
			loading = false;

			// Redirect back to lobby after showing error
			setTimeout(() => {
				goto('/lobby');
			}, 3000);
		}
	}

	/**
	 * Sync playtest store to game store for component compatibility
	 */
	function syncPlaytestToGameStore(): void {
		console.log('[syncPlaytestToGameStore] Called');
		const state = $playtestGameStore;
		const controllingPlayer = players.find((p) => p.playerId === activeControlSeat);

		if (!controllingPlayer) {
			console.log('[syncPlaytestToGameStore] No controlling player found');
			return;
		}

		// Convert playtest state to GameView format
		const gameView = {
			gameId: state.gameId,
			state: 'IN_PROGRESS',
			turn: state.turn,
			phase: '',
			step: '',
			activePlayerId: state.activePlayerId,
			activePlayerName: controllingPlayer.name,
			priorityPlayerId: activeControlSeat,
			priorityPlayerName: controllingPlayer.name,
			players: players.map((p) => ({
				playerId: p.playerId,
				name: p.name,
				life: p.life,
				poison: p.poison,
				energy: p.energy,
				libraryCount: p.libraryCount,
				handCount: p.handCount,
				hand: p.hand,
				graveyard: p.graveyard,
				manaPool: p.manaPool,
				keptHand: p.keptHand,
				hasPriority: p.playerId === activeControlSeat,
				hasAvailableActions: false,
				passed: false,
				stateOrdinal: 0,
				lost: false,
				left: false,
				wins: 0,
				losses: 0
			})),
			battlefield: state.battlefield,
			stack: state.stack,
			exile: state.exile,
			command: state.command,
			messages: [],
			revealed: [],
			lookedAt: [],
			special: false,
			isMulliganPhase: mulliganPlayerIndex !== null,
			gameFormat: 'Playtest',
			landsPlayedThisTurn: 0,
			landsAllowedThisTurn: 1
		};

		// Reinitialize the game with the current controlling player
		gameStore.initGame(state.gameId, activeControlSeat);
		gameStore.setGameView(gameView);
	}

	// Sync playtest state changes to game store
	$effect(() => {
		if (isInitialized) {
			// React to any state changes by accessing the store
			void $playtestGameStore;
			syncPlaytestToGameStore();
		}
	});

	/**
	 * Handle mulligan decision
	 */
	function handleMulligan(): void {
		if (mulliganPlayerIndex === null) return;

		const player = players[mulliganPlayerIndex];
		playtestGameStore.mulligan(player.playerId);
		toast.info(`${player.name} mulliganed`);
	}

	/**
	 * Handle keep hand decision
	 */
	function handleKeepHand(): void {
		if (mulliganPlayerIndex === null) return;

		const player = players[mulliganPlayerIndex];
		playtestGameStore.keepHand(player.playerId);
		toast.success(`${player.name} kept their hand`);

		// Move to next player
		mulliganPlayerIndex++;

		// If all players have decided, end mulligan phase
		if (mulliganPlayerIndex >= players.length) {
			mulliganPlayerIndex = null;
		}
	}

	/**
	 * Switch active control seat
	 */
	function switchPlayer(playerId: string): void {
		playtestGameStore.switchControlSeat(playerId);
		const player = players.find((p) => p.playerId === playerId);
		if (player) {
			toast.info(`Now controlling ${player.name}`);
		}
	}

	/**
	 * Handle life change
	 */
	function handleLifeChange(delta: number, playerId?: string): void {
		const targetPlayerId = playerId || me?.playerId;
		if (!targetPlayerId) return;
		playtestGameStore.modifyLife(targetPlayerId, delta);
	}

	/**
	 * Handle poison counter change
	 */
	function handlePoisonChange(delta: number, playerId?: string): void {
		const targetPlayerId = playerId || me?.playerId;
		if (!targetPlayerId) return;
		const player = players.find((p) => p.playerId === targetPlayerId);
		if (!player) return;
		const newValue = Math.max(0, (player.poison || 0) + delta);
		playtestGameStore.setPlayerCounter(targetPlayerId, 'poison', newValue);
	}

	/**
	 * Draw a card
	 */
	function handleDrawCard(): void {
		if (!me) return;
		playtestGameStore.drawCards(me.playerId, 1);
		toast.success('Drew a card');
	}

	/**
	 * Shuffle library
	 */
	function handleShuffleLibrary(): void {
		if (!me) return;
		playtestGameStore.shuffleLibrary(me.playerId);
		toast.success('Shuffled library');
	}

	/**
	 * Untap all permanents
	 */
	function handleUntapAll(): void {
		if (!me) return;
		playtestGameStore.untapAll(me.playerId);
		toast.success('Untapped all');
	}

	/**
	 * Next turn
	 */
	function handleNextTurn(): void {
		playtestGameStore.nextTurn();
		const newActivePlayer = players.find((p) => p.playerId === $playtestGameStore.activePlayerId);
		if (newActivePlayer) {
			toast.info(`${newActivePlayer.name}'s turn`);
		}
	}

	/**
	 * Handle battlefield card click
	 */
	function handleBattlefieldCardClick(cardId: string): void {
		const card = battlefield.find((c) => c.id === cardId);
		if (!card) return;

		// Toggle tap/untap
		playtestGameStore.tapCard(cardId, !card.tapped);
	}

	/**
	 * Handle battlefield card mouse down (for drag)
	 */
	function handleBattlefieldCardMouseDown(
		cardId: string,
		cardName: string,
		event: MouseEvent
	): void {
		if (event.button !== 0) return;

		event.preventDefault();
		event.stopPropagation();

		battlefieldDragStartPosition = { x: event.clientX, y: event.clientY };
		battlefieldIsDragPending = true;

		const handleMouseMove = (moveEvent: MouseEvent) => {
			if (!battlefieldDragStartPosition || !battlefieldIsDragPending) return;

			const dx = moveEvent.clientX - battlefieldDragStartPosition.x;
			const dy = moveEvent.clientY - battlefieldDragStartPosition.y;
			const distance = Math.sqrt(dx * dx + dy * dy);

			if (distance >= DRAG_THRESHOLD) {
				battlefieldIsDragPending = false;
				const validZones = getAllValidDropZones('battlefield' as SourceZone);
				dragDropStore.startDrag(
					cardId,
					cardName,
					'battlefield' as SourceZone,
					moveEvent.clientX,
					moveEvent.clientY,
					validZones
				);

				document.removeEventListener('mousemove', handleMouseMove);
				document.removeEventListener('mouseup', handleMouseUp);
			}
		};

		const handleMouseUp = () => {
			battlefieldIsDragPending = false;
			battlefieldDragStartPosition = null;
			document.removeEventListener('mousemove', handleMouseMove);
			document.removeEventListener('mouseup', handleMouseUp);
		};

		document.addEventListener('mousemove', handleMouseMove);
		document.addEventListener('mouseup', handleMouseUp);
	}

	/**
	 * Handle command zone card mouse down (for drag)
	 */
	function handleCommandCardMouseDown(cardId: string, cardName: string, event: MouseEvent): void {
		if (event.button !== 0) return;

		event.preventDefault();
		event.stopPropagation();

		commandDragStartPosition = { x: event.clientX, y: event.clientY };
		commandIsDragPending = true;

		const handleMouseMove = (moveEvent: MouseEvent) => {
			if (!commandDragStartPosition || !commandIsDragPending) return;

			const dx = moveEvent.clientX - commandDragStartPosition.x;
			const dy = moveEvent.clientY - commandDragStartPosition.y;
			const distance = Math.sqrt(dx * dx + dy * dy);

			if (distance >= DRAG_THRESHOLD) {
				commandIsDragPending = false;
				const validZones = getAllValidDropZones('command' as SourceZone);
				dragDropStore.startDrag(
					cardId,
					cardName,
					'command' as SourceZone,
					moveEvent.clientX,
					moveEvent.clientY,
					validZones
				);

				document.removeEventListener('mousemove', handleMouseMove);
				document.removeEventListener('mouseup', handleMouseUp);
			}
		};

		const handleMouseUp = () => {
			commandIsDragPending = false;
			commandDragStartPosition = null;
			document.removeEventListener('mousemove', handleMouseMove);
			document.removeEventListener('mouseup', handleMouseUp);
		};

		document.addEventListener('mousemove', handleMouseMove);
		document.addEventListener('mouseup', handleMouseUp);
	}

	/**
	 * Handle battlefield drop
	 */
	function handleBattlefieldDrop(cardId: string): void {
		const dragState = $dragDropStore;
		const sourceZone = dragState.sourceZone;

		if (sourceZone === 'hand') {
			// Move from hand to battlefield
			playtestGameStore.moveCardToZone(cardId, 'BATTLEFIELD');
			// Sync to game store
			syncPlaytestToGameStore();
		} else if (sourceZone && sourceZone !== 'battlefield') {
			playtestGameStore.moveCardToZone(cardId, 'BATTLEFIELD');
			syncPlaytestToGameStore();
		}
	}

	/**
	 * Handle zone drop (graveyard, exile, hand)
	 */
	function handleZoneDrop(cardId: string, zone: string): void {
		playtestGameStore.moveCardToZone(cardId, zone);
		syncPlaytestToGameStore();
	}

	/**
	 * Copy game log to clipboard
	 */
	async function handleCopyLog(): Promise<void> {
		const logText = playtestGameStore.buildLogText($playtestGameStore);
		try {
			if (typeof navigator !== 'undefined' && navigator.clipboard?.writeText) {
				await navigator.clipboard.writeText(logText);
				toast.success('Game log copied to clipboard!');
				return;
			}
			// Fallback for older browsers
			const textarea = document.createElement('textarea');
			textarea.value = logText;
			textarea.style.position = 'fixed';
			textarea.style.top = '0';
			textarea.style.left = '0';
			textarea.style.opacity = '0';
			document.body.appendChild(textarea);
			textarea.focus();
			textarea.select();
			const ok = document.execCommand('copy');
			document.body.removeChild(textarea);
			if (ok) {
				toast.success('Game log copied to clipboard!');
			} else {
				toast.error('Failed to copy log');
			}
		} catch (err) {
			console.error('Failed to copy log to clipboard:', err);
			toast.error('Failed to copy log');
		}
	}

	/**
	 * Handle keyboard shortcuts
	 */
	function handleGlobalKeydown(event: KeyboardEvent): void {
		if (event.target instanceof HTMLInputElement || event.target instanceof HTMLTextAreaElement) {
			return;
		}

		const key = event.key.toLowerCase();

		switch (key) {
			case 'm':
				// M - Toggle menu (only when game started)
				if (isGameStarted) {
					showMenu = !showMenu;
					event.preventDefault();
				}
				break;
			case 'escape':
				// Escape - Close menu or modals
				if (showMenu) {
					showMenu = false;
					event.preventDefault();
				} else if (showKeyboardShortcuts) {
					showKeyboardShortcuts = false;
					event.preventDefault();
				} else if (showTokenCreator) {
					showTokenCreator = false;
					event.preventDefault();
				} else if (showCreateTokenDialog) {
					showCreateTokenDialog = false;
					event.preventDefault();
				} else if (showCounterDialog) {
					showCounterDialog = false;
					selectedCardForCounters = null;
					event.preventDefault();
				} else if (showDebugOverlay) {
					showDebugOverlay = false;
					event.preventDefault();
				}
				break;
			case '?':
				showKeyboardShortcuts = !showKeyboardShortcuts;
				event.preventDefault();
				break;
			case 'f':
				// F - Search your deck
				showDeckSearch = true;
				event.preventDefault();
				break;
			case 'x':
				handleUntapAll();
				event.preventDefault();
				break;
			case 'c':
				handleDrawCard();
				event.preventDefault();
				break;
			case 'v':
				handleShuffleLibrary();
				event.preventDefault();
				break;
			case 'e':
				handleNextTurn();
				event.preventDefault();
				break;
			case 'w':
				showCreateTokenDialog = true;
				event.preventDefault();
				break;
		}

		// Hover card shortcuts
		if (hoveredCard) {
			switch (key) {
				case 'd':
					playtestGameStore.moveCardToZone(hoveredCard.id, 'GRAVEYARD');
					event.preventDefault();
					break;
				case 's':
					playtestGameStore.moveCardToZone(hoveredCard.id, 'EXILE');
					event.preventDefault();
					break;
				case 'r':
					playtestGameStore.moveCardToZone(hoveredCard.id, 'HAND');
					event.preventDefault();
					break;
				case 't':
					playtestGameStore.moveCardToZone(hoveredCard.id, 'LIBRARY');
					event.preventDefault();
					break;
				case 'k':
					selectedCardForCounters = { id: hoveredCard.id, name: hoveredCard.name };
					showCounterDialog = true;
					event.preventDefault();
					break;
			}
		}
	}

	/**
	 * Register drop zones
	 */
	$effect(() => {
		if (battlefieldDropZoneEl && !dropZoneUnregister) {
			dropZoneUnregister = dragDropStore.registerDropZone({
				id: 'battlefield',
				type: 'battlefield',
				element: battlefieldDropZoneEl,
				accepts: (_cardId, sourceZone) => sourceZone !== 'battlefield',
				onDrop: handleBattlefieldDrop
			});
		}
		return () => {
			if (dropZoneUnregister) {
				dropZoneUnregister();
				dropZoneUnregister = null;
			}
		};
	});

	$effect(() => {
		if (graveyardDropZoneEl && !graveyardDropZoneUnregister) {
			graveyardDropZoneUnregister = dragDropStore.registerDropZone({
				id: 'graveyard',
				type: 'graveyard',
				element: graveyardDropZoneEl,
				accepts: (_cardId, sourceZone) => sourceZone !== 'graveyard',
				onDrop: (cardId) => handleZoneDrop(cardId, 'GRAVEYARD')
			});
		}
		return () => {
			if (graveyardDropZoneUnregister) {
				graveyardDropZoneUnregister();
				graveyardDropZoneUnregister = null;
			}
		};
	});

	$effect(() => {
		if (exileDropZoneEl && !exileDropZoneUnregister) {
			exileDropZoneUnregister = dragDropStore.registerDropZone({
				id: 'exile',
				type: 'exile',
				element: exileDropZoneEl,
				accepts: (_cardId, sourceZone) => sourceZone !== 'exile',
				onDrop: (cardId) => handleZoneDrop(cardId, 'EXILE')
			});
		}
		return () => {
			if (exileDropZoneUnregister) {
				exileDropZoneUnregister();
				exileDropZoneUnregister = null;
			}
		};
	});

	$effect(() => {
		if (handDropZoneEl && !handDropZoneUnregister) {
			handDropZoneUnregister = dragDropStore.registerDropZone({
				id: 'hand',
				type: 'hand',
				element: handDropZoneEl,
				accepts: (_cardId, sourceZone) => sourceZone !== 'hand',
				onDrop: (cardId) => handleZoneDrop(cardId, 'HAND')
			});
		}
		return () => {
			if (handDropZoneUnregister) {
				handDropZoneUnregister();
				handDropZoneUnregister = null;
			}
		};
	});

	$effect(() => {
		if (libraryDropZoneEl && !libraryDropZoneUnregister) {
			libraryDropZoneUnregister = dragDropStore.registerDropZone({
				id: 'library',
				type: 'library',
				element: libraryDropZoneEl,
				accepts: (_cardId, sourceZone) => sourceZone !== 'library',
				onDrop: (cardId) => handleZoneDrop(cardId, 'LIBRARY')
			});
		}
		return () => {
			if (libraryDropZoneUnregister) {
				libraryDropZoneUnregister();
				libraryDropZoneUnregister = null;
			}
		};
	});

	$effect(() => {
		if (commandDropZoneEl && !commandDropZoneUnregister) {
			commandDropZoneUnregister = dragDropStore.registerDropZone({
				id: 'command',
				type: 'command',
				element: commandDropZoneEl,
				accepts: (_cardId, sourceZone) => sourceZone !== 'command',
				onDrop: (cardId) => handleZoneDrop(cardId, 'COMMAND')
			});
		}
		return () => {
			if (commandDropZoneUnregister) {
				commandDropZoneUnregister();
				commandDropZoneUnregister = null;
			}
		};
	});

	/**
	 * Load available sessions for restoration
	 */
	function loadAvailableSessions(): void {
		availableSessions = playtestGameStore.listSessions();
	}

	/**
	 * Restore a specific session
	 */
	async function restoreSession(sessionId: string): Promise<void> {
		const success = playtestGameStore.restoreSession(sessionId);
		if (success) {
			loading = false;
			showSessionPicker = false;

			// Restore mulligan phase based on first player who hasn't kept.
			const idx = players.findIndex((p) => !p.keptHand);
			mulliganPlayerIndex = idx === -1 ? null : idx;
			mulliganCount = 0;

			// Ensure the normal game store is initialized for shared components.
			gameStore.initGame($playtestGameStore.gameId, $playtestGameStore.activeControlSeat);
			syncPlaytestToGameStore();

			// Set playtestId in URL
			await updateUrlWithPlaytestId($playtestGameStore.gameId);

			toast.success('Session restored');
		} else {
			toast.error('Failed to restore session');
		}
	}

	/**
	 * Delete a session
	 */
	function deleteSession(sessionId: string): void {
		playtestGameStore.deleteSession(sessionId);
		loadAvailableSessions();
		toast.info('Session deleted');
	}

	async function initializeFromPlaytestId(): Promise<void> {
		const playtestId = $page.url.searchParams.get('playtestId');
		if (playtestId) {
			const success = playtestGameStore.restoreSession(playtestId);
			if (success) {
				loading = false;

				// Restore mulligan phase based on first player who hasn't kept.
				const idx = players.findIndex((p) => !p.keptHand);
				mulliganPlayerIndex = idx === -1 ? null : idx;
				mulliganCount = 0;

				// Ensure the normal game store is initialized for shared components.
				gameStore.initGame($playtestGameStore.gameId, $playtestGameStore.activeControlSeat);
				syncPlaytestToGameStore();

				// URL already has playtestId, but ensure it's set correctly
				await updateUrlWithPlaytestId(playtestId);
			} else {
				error = 'Failed to restore playtest session';
				loading = false;
			}
		}
	}

	// Initialize on mount
	onMount(() => {
		// Check if URL has deck params (user wants to start a new playtest)
		const hasUrlDecks = $page.url.searchParams.has('d1') || $page.url.searchParams.has('d2');
		const hasPlaytestIdInUrl = $page.url.searchParams.has('playtestId');

		// Otherwise, check for existing sessions
		loadAvailableSessions();

		if (hasPlaytestIdInUrl) {
			initializeFromPlaytestId();
			return;
		}

		// If URL has decks, start new playtest
		if (hasUrlDecks) {
			initializeFromUrl();
			return;
		}

		if (availableSessions.length === 0) {
			// No sessions and no URL decks - redirect to lobby
			error = 'No playtest sessions found. Please configure a new playtest from the lobby.';
			setTimeout(() => {
				goto('/lobby');
			}, 2000);
			return;
		}

		// If there's an active session, restore it
		if ($playtestGameStore.isInitialized) {
			loading = false;

			// Restore mulligan phase based on first player who hasn't kept.
			const idx = players.findIndex((p) => !p.keptHand);
			mulliganPlayerIndex = idx === -1 ? null : idx;
			mulliganCount = 0;

			// Ensure the normal game store is initialized for shared components.
			gameStore.initGame($playtestGameStore.gameId, $playtestGameStore.activeControlSeat);
			syncPlaytestToGameStore();

			// Set playtestId in URL
			updateUrlWithPlaytestId($playtestGameStore.gameId);
			return;
		}

		// Show session picker if we have sessions but none active
		if (availableSessions.length > 0) {
			loading = false;
			showSessionPicker = true;
		}
	});

	// Remove the beforeNavigate cleanup - we want to persist sessions across navigation
</script>

<svelte:window onkeydown={handleGlobalKeydown} />

<svelte:head>
	<title>Playtest Mode - MAGE</title>
</svelte:head>

<div class="playtest-container">
	{#if loading}
		<div class="loading-overlay">
			<div class="spinner"></div>
			<p>Setting up playtest...</p>
		</div>
	{:else if error}
		<div class="error-overlay">
			<div class="error-icon">⚠️</div>
			<h2>Error</h2>
			<p>{error}</p>
			<button class="btn-primary" onclick={() => goto('/lobby')}> Return to Lobby </button>
		</div>
	{:else if showSessionPicker}
		<div class="session-picker-overlay">
			<div class="session-picker-modal">
				<h2>Restore Playtest Session</h2>
				<p class="session-picker-hint">
					Select a recent playtest session to continue, or start a new one.
				</p>

				{#if availableSessions.length > 0}
					<div class="sessions-list">
						{#each availableSessions as session (session.id)}
							<div class="session-card">
								<div class="session-info">
									<div class="session-label">{session.label}</div>
									<div class="session-meta">
										{session.playerCount} players · Turn {session.turn} ·
										{new Date(session.savedAt).toLocaleDateString()}
										{new Date(session.savedAt).toLocaleTimeString([], {
											hour: '2-digit',
											minute: '2-digit'
										})}
									</div>
								</div>
								<div class="session-actions">
									<button class="btn-restore" onclick={() => restoreSession(session.id)}>
										Restore
									</button>
									<button
										class="btn-delete"
										onclick={() => deleteSession(session.id)}
										title="Delete session"
									>
										✕
									</button>
								</div>
							</div>
						{/each}
					</div>
				{:else}
					<p class="no-sessions">No saved sessions found.</p>
				{/if}

				<div class="session-picker-actions">
					<button class="btn-secondary" onclick={() => goto('/lobby')}> Back to Lobby </button>
					<button class="btn-primary" onclick={() => goto('/lobby')}> Start New Playtest </button>
				</div>
			</div>
		</div>
	{:else if !isInitialized}
		<div class="loading-overlay">
			<p>Initializing game state...</p>
		</div>
	{:else if mulliganPlayerIndex !== null && !allPlayersKept}
		<MulliganDialog
			cards={players[mulliganPlayerIndex]?.hand || []}
			mulliganCount={players[mulliganPlayerIndex]?.mulliganCount || 0}
			freeMulligans={$playtestGameStore.freeMulligans}
			playerName={players[mulliganPlayerIndex]?.name}
			onKeep={handleKeepHand}
			onMulligan={handleMulligan}
			isLoading={false}
			hasKeptHand={false}
		/>
	{:else}
		<!-- Hamburger Menu Button (shown when game started) -->
		{#if isGameStarted}
			<button class="menu-toggle-btn" onclick={() => (showMenu = !showMenu)} aria-label="Open menu">
				<Menu size={24} />
			</button>
		{/if}

		<!-- Playtest Header -->
		<div class="playtest-header" class:hidden={isGameStarted}>
			<div class="header-left">
				<button class="btn-back" onclick={() => goto('/lobby')}> ← Back to Lobby </button>
				<span class="mode-badge">Playtest Mode</span>
				<button
					class="btn-sessions"
					onclick={() => {
						loadAvailableSessions();
						showSessionPicker = true;
					}}
					title="Manage sessions"
				>
					Sessions ({availableSessions.length})
				</button>
			</div>

			<div class="playtest-controls">
				<label for="playtest-controlling-select">Controlling:</label>
				<select
					id="playtest-controlling-select"
					class="player-select"
					value={activeControlSeat}
					onchange={(e) => switchPlayer(e.currentTarget.value)}
				>
					{#each players as player}
						<option value={player.playerId}>{player.name}</option>
					{/each}
				</select>
				<button class="btn-toggle" onclick={() => (showAllHands = !showAllHands)}>
					{showAllHands ? '🙈 Hide' : '👁️ Show'} All Hands
				</button>
			</div>

			<div class="header-right">
				<div class="turn-indicator" title="Current turn">
					<Clock size={16} aria-hidden="true" />
					<span class="turn-text">
						Turn {turnNumber()}{activePlayerName() ? ` · ${activePlayerName()}` : ''}
					</span>
				</div>
				<button class="btn-action" onclick={handleNextTurn}>Next Turn</button>
				<button
					class="btn-debug"
					onclick={() => (showKeyboardShortcuts = true)}
					title="Keyboard shortcuts (?)"
					aria-label="Keyboard shortcuts"
				>
					<Keyboard size={20} aria-hidden="true" />
				</button>
				<button class="btn-debug" onclick={() => (showDebugOverlay = true)} title="Debug View">
					🔧
				</button>
			</div>
		</div>

		<!-- Menu Overlay (slide-in from right) -->
		{#if showMenu}
			<!-- Backdrop -->
			<div
				class="menu-backdrop"
				role="button"
				tabindex="0"
				onclick={() => (showMenu = false)}
				onkeydown={(e) => e.key === 'Escape' && (showMenu = false)}
			></div>

			<!-- Menu Panel -->
			<div class="menu-overlay open">
				<div class="menu-header">
					<h2>Menu</h2>
					<button class="menu-close-btn" onclick={() => (showMenu = false)} aria-label="Close menu">
						<X size={24} />
					</button>
				</div>

				<div class="menu-content">
					<!-- Controls Section -->
					<div class="menu-section">
						<h3 class="menu-section-title">Controls</h3>
						<div class="menu-section-content">
							<label>
								<span class="menu-label">Controlling:</span>
								<select
									class="control-select"
									value={activeControlSeat}
									onchange={(e) => playtestGameStore.switchControlSeat(e.currentTarget.value)}
								>
									{#each players as player}
										<option value={player.playerId}>{player.name}</option>
									{/each}
								</select>
							</label>

							<button class="menu-btn" onclick={() => (showAllHands = !showAllHands)}>
								{showAllHands ? '🙈 Hide' : '👁️ Show'} All Hands
							</button>
						</div>
					</div>

					<!-- Turn Info Section -->
					<div class="menu-section">
						<h3 class="menu-section-title">Turn Info</h3>
						<div class="menu-section-content">
							<div class="turn-info">
								<Clock size={18} />
								<span>Turn {turnNumber()}</span>
								{#if activePlayerName()}
									<span class="active-player">· {activePlayerName()}</span>
								{/if}
							</div>
							<button class="menu-btn primary" onclick={handleNextTurn}>Next Turn</button>
						</div>
					</div>

					<!-- Utility Section -->
					<div class="menu-section">
						<h3 class="menu-section-title">Utilities</h3>
						<div class="menu-section-content">
							<button
								class="menu-btn"
								onclick={() => {
									showKeyboardShortcuts = true;
									showMenu = false;
								}}
							>
								<Keyboard size={18} />
								Keyboard Shortcuts
							</button>
							<button
								class="menu-btn"
								onclick={() => {
									showDebugOverlay = true;
									showMenu = false;
								}}
							>
								🔧 Debug View
							</button>
						</div>
					</div>

					<!-- Navigation Section -->
					<div class="menu-section">
						<h3 class="menu-section-title">Navigation</h3>
						<div class="menu-section-content">
							<button class="menu-btn" onclick={() => goto('/lobby')}> ← Back to Lobby </button>
							<button
								class="menu-btn"
								onclick={() => {
									showSessionPicker = true;
									showMenu = false;
								}}
							>
								<Clock size={18} />
								Sessions
							</button>
						</div>
					</div>
				</div>
			</div>
		{/if}

		<!-- All Hands Overlay -->
		{#if showAllHands}
			<div class="all-hands-overlay">
				{#each players as player}
					<div class="player-hand-compact" class:active={player.playerId === activeControlSeat}>
						<div class="compact-header">
							<span class="player-name-compact">{player.name}</span>
							<span class="life-compact">❤️ {player.life}</span>
						</div>
						<div class="cards-compact">
							{#each player.hand as card}
								<div class="card-mini">
									<img
										src={getScryfallImageUrl(card.name, 'small')}
										alt={card.name}
										title={card.name}
									/>
								</div>
							{/each}
						</div>
					</div>
				{/each}
			</div>
		{/if}

		<!-- Main Game Area -->
		<main class="game-layout">
			<!-- Opponent Section -->
			{#if selectedOpponent()}
				{@const opponent = selectedOpponent()!}
				<div class="opponent-section">
					<!-- Opponent Battlefield (Non-editable) -->
					<div class="battlefield-area opponent-battlefield">
						<!-- Opponent Info Overlay -->
						<div class="opponent-info-overlay">
							<div class="opponent-identity">
								{#if otherPlayers.length > 1}
									<select
										class="opponent-select"
										value={opponent.playerId}
										onchange={(e) => (selectedOpponentId = e.currentTarget.value)}
									>
										{#each otherPlayers as opp}
											<option value={opp.playerId}>{opp.name}</option>
										{/each}
									</select>
								{:else}
									<span class="opponent-name-label">{opponent.name}</span>
								{/if}
							</div>

							<div class="opponent-stats-compact">
								<div class="life-group">
									<button
										class="stat-btn minus"
										onclick={() => handleLifeChange(-1, opponent.playerId)}
										title="Decrease life"
									>
										−
									</button>
									<button
										class="stat-display life"
										onclick={() => (showOpponentLifeMenu = !showOpponentLifeMenu)}
										title="Life total"
									>
										<span class="stat-icon">❤️</span>
										<span class="stat-value">{opponent.life}</span>
									</button>
									<button
										class="stat-btn plus"
										onclick={() => handleLifeChange(1, opponent.playerId)}
										title="Increase life"
									>
										+
									</button>
								</div>

								{#if opponent.poison > 0}
									<div class="stat-display poison" title="Poison counters">
										<span class="stat-icon">☠️</span>
										<span class="stat-value">{opponent.poison}</span>
									</div>
								{/if}

								<div class="opponent-counts">
									<span class="opponent-count" title="Hand cards">🃏 {opponent.handCount}</span>
									<span class="opponent-count" title="Library cards"
										>📚 {opponent.libraryCount}</span
									>
									<span class="opponent-count" title="Graveyard cards"
										>🪦 {opponent.graveyard.length}</span
									>
								</div>

								{#if showOpponentLifeMenu}
									<div class="quick-menu opponent-menu">
										<div class="menu-section">
											<span class="menu-label">Life</span>
											<div class="menu-row">
												<button onclick={() => handleLifeChange(-5, opponent.playerId)}>−5</button>
												<button onclick={() => handleLifeChange(-1, opponent.playerId)}>−1</button>
												<button onclick={() => handleLifeChange(1, opponent.playerId)}>+1</button>
												<button onclick={() => handleLifeChange(5, opponent.playerId)}>+5</button>
											</div>
										</div>
										<div class="menu-section">
											<span class="menu-label">Poison</span>
											<div class="menu-row">
												<button onclick={() => handlePoisonChange(-1, opponent.playerId)}>−1</button
												>
												<span class="menu-value">{opponent.poison}</span>
												<button onclick={() => handlePoisonChange(1, opponent.playerId)}>+1</button>
											</div>
										</div>
										<button class="menu-close" onclick={() => (showOpponentLifeMenu = false)}>
											✕
										</button>
									</div>
								{/if}
							</div>
						</div>

						<!-- Opponent Battlefield Content Wrapper -->
						<div class="battlefield-content-wrapper">
							<!-- Opponent Battlefield Main -->
							<div class="battlefield-main">
								<div class="battlefield-rows">
									{#if opponentBattlefieldNonlands().length > 0}
										<div class="battlefield-row battlefield-row--nonlands">
											{#each opponentBattlefieldNonlands() as card (card.id)}
												<div
													class="battlefield-card-wrapper readonly"
													title="{card.name} (controlled by {opponent.name})"
												>
													<Card
														cardId={card.id}
														cardName={card.name}
														manaCost={card.manaCost}
														cardType={card.type}
														power={card.power}
														toughness={card.toughness}
														color={card.color}
														imageUrl=""
														isTapped={card.tapped}
														isSelected={false}
														counters={card.counters}
														size="normal"
														onclick={() => {}}
														oncontextmenu={(e) => {
															e.preventDefault();
															selectedCardForCounters = { id: card.id, name: card.name };
															showCounterDialog = true;
														}}
													/>
												</div>
											{/each}
										</div>
									{/if}

									{#if opponentBattlefieldLands().length > 0}
										<div class="battlefield-row battlefield-row--lands">
											{#each opponentBattlefieldLands() as card (card.id)}
												<div
													class="battlefield-card-wrapper readonly"
													title="{card.name} (controlled by {opponent.name})"
												>
													<Card
														cardId={card.id}
														cardName={card.name}
														manaCost={card.manaCost}
														cardType={card.type}
														power={card.power}
														toughness={card.toughness}
														color={card.color}
														imageUrl=""
														isTapped={card.tapped}
														isSelected={false}
														counters={card.counters}
														size="normal"
														onclick={() => {}}
														oncontextmenu={(e) => {
															e.preventDefault();
															selectedCardForCounters = { id: card.id, name: card.name };
															showCounterDialog = true;
														}}
													/>
												</div>
											{/each}
										</div>
									{/if}

									{#if opponentBattlefieldNonlands().length === 0 && opponentBattlefieldLands().length === 0}
										<div class="empty-battlefield">No permanents</div>
									{/if}
								</div>
							</div>

							<!-- Opponent Command Zone (right side) -->
							{#if isCommanderGame}
								<div class="command-zone opponent-command-zone">
									<span class="zone-label">Command Zone</span>
									<div class="command-cards">
										{#if opponentCommandCards().length === 0}
											<div class="command-zone-empty">
												<span class="zone-empty-text">Empty</span>
											</div>
										{/if}
										{#each opponentCommandCards() as card (card.id)}
											<div class="command-card-wrapper readonly" title={card.name}>
												<Card
													cardId={card.id}
													cardName={card.name}
													manaCost={card.manaCost}
													cardType={card.type}
													power={card.power}
													toughness={card.toughness}
													color={card.color}
													imageUrl=""
													isTapped={card.tapped}
													isSelected={false}
													size="small"
													onclick={() => {}}
												/>
											</div>
										{/each}
									</div>
								</div>
							{/if}
						</div>
					</div>
				</div>
			{/if}

			<!-- My Battlefield Area (Editable) -->
			<div
				bind:this={battlefieldDropZoneEl}
				class="battlefield-area my-battlefield"
				class:drag-active={isDragging}
				class:drag-valid={isDragging && isOverValidDrop && dropZone === 'battlefield'}
			>
				<div class="battlefield-content-wrapper">
					<!-- Main Battlefield (left side) -->
					<div class="battlefield-main">
						<span class="zone-label">Your Battlefield</span>
						<div class="battlefield-rows">
							{#if myBattlefieldNonlands.length > 0}
								<div class="battlefield-row battlefield-row--nonlands">
									{#each myBattlefieldNonlands as card (card.id)}
										<div
											class="battlefield-card-wrapper"
											class:is-hovered={hoveredCardId === card.id}
											role="button"
											tabindex="0"
											aria-label={card.name}
											onmousedown={(e) => handleBattlefieldCardMouseDown(card.id, card.name, e)}
											onkeydown={(e) => {
												if (e.key === 'Enter' || e.key === ' ') {
													e.preventDefault();
													handleBattlefieldCardClick(card.id);
												}
											}}
											onmouseenter={() => (hoveredCardId = card.id)}
											onmouseleave={() => {
												if (hoveredCardId === card.id) hoveredCardId = null;
											}}
										>
											<Card
												cardId={card.id}
												cardName={card.name}
												manaCost={card.manaCost}
												cardType={card.type}
												power={card.power}
												toughness={card.toughness}
												color={card.color}
												imageUrl=""
												isTapped={card.tapped}
												isSelected={false}
												counters={card.counters}
												size="normal"
												onclick={() => handleBattlefieldCardClick(card.id)}
												oncontextmenu={(e) => {
													e.preventDefault();
													selectedCardForCounters = { id: card.id, name: card.name };
													showCounterDialog = true;
												}}
											/>
										</div>
									{/each}
								</div>
							{/if}

							{#if myBattlefieldLands.length > 0}
								<div class="battlefield-row battlefield-row--lands">
									{#each myBattlefieldLands as card (card.id)}
										<div
											class="battlefield-card-wrapper"
											class:is-hovered={hoveredCardId === card.id}
											role="button"
											tabindex="0"
											aria-label={card.name}
											onmousedown={(e) => handleBattlefieldCardMouseDown(card.id, card.name, e)}
											onkeydown={(e) => {
												if (e.key === 'Enter' || e.key === ' ') {
													e.preventDefault();
													handleBattlefieldCardClick(card.id);
												}
											}}
											onmouseenter={() => (hoveredCardId = card.id)}
											onmouseleave={() => {
												if (hoveredCardId === card.id) hoveredCardId = null;
											}}
										>
											<Card
												cardId={card.id}
												cardName={card.name}
												manaCost={card.manaCost}
												cardType={card.type}
												power={card.power}
												toughness={card.toughness}
												color={card.color}
												imageUrl=""
												isTapped={card.tapped}
												isSelected={false}
												counters={card.counters}
												size="normal"
												onclick={() => handleBattlefieldCardClick(card.id)}
												oncontextmenu={(e) => {
													e.preventDefault();
													selectedCardForCounters = { id: card.id, name: card.name };
													showCounterDialog = true;
												}}
											/>
										</div>
									{/each}
								</div>
							{/if}

							{#if myBattlefield.length === 0}
								<div class="empty-battlefield">
									{#if isDragging}
										<span class="drop-hint">Drop card here to play</span>
									{:else}
										No permanents
									{/if}
								</div>
							{/if}
						</div>
					</div>

					<!-- Command Zone (right side) -->
					{#if isCommanderGame}
						<div
							bind:this={commandDropZoneEl}
							class="command-zone"
							class:drag-valid={isDragging && isOverValidDrop && dropZone === 'command'}
						>
							<span class="zone-label">Command Zone</span>
							<div class="command-cards">
								{#if myCommandCards.length === 0}
									<div class="command-zone-empty">
										<span class="zone-empty-text">Empty</span>
										<span class="zone-empty-hint">Drag commander here</span>
									</div>
								{/if}
								{#each myCommandCards as card (card.id)}
									<div
										class="command-card-wrapper"
										title={card.name}
										role="button"
										tabindex="0"
										aria-label={card.name}
										onmousedown={(e) => handleCommandCardMouseDown(card.id, card.name, e)}
									>
										<Card
											cardId={card.id}
											cardName={card.name}
											manaCost={card.manaCost}
											cardType={card.type}
											power={card.power}
											toughness={card.toughness}
											color={card.color}
											imageUrl=""
											isTapped={card.tapped}
											isSelected={false}
											size="small"
											onclick={() => {}}
										/>
									</div>
								{/each}
							</div>
						</div>
					{/if}
				</div>
			</div>

			<!-- Player Info Row -->
			{#if me}
				<div class="player-info-row">
					<div class="player-identity">
						<span class="player-name">{me.name}</span>
					</div>

					<div class="player-stats-inline">
						<div class="life-group">
							<button class="stat-btn minus" onclick={() => handleLifeChange(-1)}>−</button>
							<button class="stat-display life" onclick={() => (showLifeMenu = !showLifeMenu)}>
								<span class="stat-icon">❤️</span>
								<span class="stat-value">{me.life}</span>
							</button>
							<button class="stat-btn plus" onclick={() => handleLifeChange(1)}>+</button>
						</div>

						{#if me.poison > 0}
							<div class="stat-display poison">
								<span class="stat-icon">☠️</span>
								<span class="stat-value">{me.poison}</span>
							</div>
						{/if}

						<div bind:this={libraryDropZoneEl} class="library-drop-zone">
							<LibraryZone
								libraryCount={me.libraryCount}
								playerName="You"
								onSearch={() => {
									showDeckSearch = true;
								}}
							/>
						</div>

						{#if showLifeMenu}
							<div bind:this={lifeMenuEl} class="quick-menu">
								<div class="menu-section">
									<span class="menu-label">Life</span>
									<div class="menu-row">
										<button onclick={() => handleLifeChange(-5)}>−5</button>
										<button onclick={() => handleLifeChange(-1)}>−1</button>
										<button onclick={() => handleLifeChange(1)}>+1</button>
										<button onclick={() => handleLifeChange(5)}>+5</button>
									</div>
								</div>
								<div class="menu-section">
									<span class="menu-label">Poison</span>
									<div class="menu-row">
										<button onclick={() => handlePoisonChange(-1)}>−1</button>
										<span class="menu-value">{me.poison}</span>
										<button onclick={() => handlePoisonChange(1)}>+1</button>
									</div>
								</div>
								<button class="menu-close" onclick={() => (showLifeMenu = false)}>✕</button>
							</div>
						{/if}
					</div>

					<div class="player-zones">
						<div bind:this={graveyardDropZoneEl} class="graveyard-drop-zone">
							<Graveyard
								cards={myGrave.map((c) => ({
									id: c.id,
									name: c.name,
									manaCost: c.manaCost,
									cardType: c.type,
									power: c.power,
									toughness: c.toughness,
									imageUrl: '',
									isTapped: false,
									isSelected: false
								}))}
								playerName="You"
								isOpponent={false}
								canDrag={true}
								onCardClick={() => {}}
							/>
						</div>
						<div bind:this={exileDropZoneEl} class="exile-drop-zone">
							<ExileZone
								cards={exile.map((c) => ({
									id: c.id,
									name: c.name,
									manaCost: c.manaCost,
									cardType: c.type,
									power: c.power,
									toughness: c.toughness,
									imageUrl: '',
									isTapped: false,
									isSelected: false
								}))}
								playerName="You"
								isOpponent={false}
								canDrag={true}
								onCardClick={() => {}}
								compact={true}
							/>
						</div>
						<ManaPool mana={myMana} showEmpty={false} size="small" onManaClick={() => {}} />
					</div>
				</div>
			{/if}

			<!-- Player Hand -->
			<div
				bind:this={handDropZoneEl}
				class="hand-area"
				class:drag-active={isDragging}
				class:drag-valid={isDragging && isOverValidDrop && dropZone === 'hand'}
			>
				<PlayerHand
					onCardClick={() => {}}
					size="normal"
					currentPhase="PRECOMBAT_MAIN"
					canDrag={true}
				/>
			</div>
		</main>

		<!-- Token Creator -->
		{#if showTokenCreator}
			<TokenCreator gameId="playtest" onClose={() => (showTokenCreator = false)} />
		{/if}

		<!-- Create Token Dialog (New) -->
		{#if showCreateTokenDialog}
			<CreateTokenDialog
				onCreateToken={(name, types, power, toughness, color) => {
					playtestGameStore.createToken(name, types, power, toughness, color);
					syncPlaytestToGameStore();
					showCreateTokenDialog = false;
				}}
				onClose={() => (showCreateTokenDialog = false)}
			/>
		{/if}

		<!-- Counter Dialog -->
		{#if showCounterDialog && selectedCardForCounters && selectedCardForCountersData()}
			<CounterDialog
				cardName={selectedCardForCountersData().name}
				cardId={selectedCardForCountersData().id}
				currentCounters={selectedCardForCountersData().counters}
				onAddCounter={(counterName, amount) => {
					const card = selectedCardForCountersData();
					playtestGameStore.addCounter(card.id, counterName, amount);
					syncPlaytestToGameStore();
				}}
				onRemoveCounter={(counterName, amount) => {
					const card = selectedCardForCountersData();
					playtestGameStore.removeCounter(card.id, counterName, amount);
					syncPlaytestToGameStore();
				}}
				onSetCounter={(counterName, amount) => {
					const card = selectedCardForCountersData();
					playtestGameStore.setCounter(card.id, counterName, amount);
					syncPlaytestToGameStore();
				}}
				onClose={() => {
					showCounterDialog = false;
					selectedCardForCounters = null;
				}}
			/>
		{/if}

		<!-- Deck Search -->
		{#if showDeckSearch && me}
			<PlaytestLibrarySearch
				cards={me.library}
				playerName="You"
				onMove={(cardId, zone) => {
					playtestGameStore.moveCardToZone(cardId, zone);
					syncPlaytestToGameStore();
				}}
				onShuffle={() => {
					playtestGameStore.shuffleLibrary(me.playerId);
					syncPlaytestToGameStore();
				}}
				onClose={() => (showDeckSearch = false)}
			/>
		{/if}

		<KeyboardShortcutsModal bind:open={showKeyboardShortcuts} mode="playtest" />

		<!-- Debug Overlay -->
		{#if showDebugOverlay}
			<div class="debug-overlay" role="dialog" aria-modal="true">
				<div class="debug-modal">
					<header class="debug-header">
						<div class="debug-header-left">
							<h2>🔧 Playtest Debug View</h2>
							<div class="debug-status connected">● Playtest Mode</div>
						</div>
						<div class="debug-header-right">
							<button class="debug-close" onclick={() => (showDebugOverlay = false)}>✕</button>
						</div>
					</header>

					<main class="debug-content">
						<!-- Game State Overview -->
						<section class="debug-section">
							<div class="debug-section-header">
								<span>Game State Overview</span>
							</div>
							<div class="debug-code">
								<pre><code
										>{@html `<span class="dk">activeControlSeat:</span> <span class="ds">"${activeControlSeat}"</span>
<span class="dk">turn:</span> <span class="dn">${$playtestGameStore.turn}</span>
<span class="dk">activePlayerId:</span> <span class="ds">"${$playtestGameStore.activePlayerId}"</span>
<span class="dk">isInitialized:</span> <span class="db">${isInitialized}</span>`}</code
									></pre>
							</div>
						</section>

						<!-- Game State Log -->
						<section class="debug-section">
							<div class="debug-section-header">
								<span>Game State Log ({gameLog.length} events)</span>
								<button
									class="debug-copy-btn"
									onclick={handleCopyLog}
									title="Copy log to clipboard"
									aria-label="Copy log to clipboard"
								>
									<Copy size={16} aria-hidden="true" />
									<span>Copy</span>
								</button>
							</div>
							<div class="debug-log-container">
								{#if gameLog.length === 0}
									<div class="debug-log-empty">No events logged yet</div>
								{:else}
									<div class="debug-log-entries">
										{#each gameLog as entry (entry.id)}
											<div class="debug-log-entry">
												<span class="debug-log-time">
													{new Date(entry.at).toLocaleTimeString([], {
														hour: '2-digit',
														minute: '2-digit',
														second: '2-digit'
													})}
												</span>
												<span class="debug-log-turn">T{entry.turn}</span>
												<span class="debug-log-kind">{entry.kind}</span>
												<span class="debug-log-message">{entry.message}</span>
											</div>
										{/each}
									</div>
								{/if}
							</div>
						</section>

						<!-- Players -->
						<section class="debug-section">
							<div class="debug-section-header">
								<span>Players ({players.length})</span>
							</div>
							{#each players as player (player.playerId)}
								<div class="debug-player">
									<div class="debug-player-header">
										<span class="debug-badge" class:local={player.playerId === activeControlSeat}>
											{player.playerId === activeControlSeat ? '👤 Controlling' : '🎮 Other'}
										</span>
										<span>{player.name}</span>
									</div>
									<div class="debug-code">
										<pre><code
												>{@html `<span class="dk">playerId:</span> <span class="ds">"${player.playerId}"</span>
<span class="dk">life:</span> <span class="dn">${player.life}</span>
<span class="dk">poison:</span> <span class="dn">${player.poison}</span>
<span class="dk">libraryCount:</span> <span class="dn">${player.libraryCount}</span>
<span class="dk">handCount:</span> <span class="dn">${player.handCount}</span>
<span class="dk">hand:</span> [${player.hand.map((c) => `\n  <span class="ds">"${c.name}"</span> <span class="dc">// ${c.id}</span>`).join(',') || ''}
]
<span class="dk">graveyard:</span> [${player.graveyard.map((c) => `\n  <span class="ds">"${c.name}"</span>`).join(',') || ''}
]
<span class="dk">library (first 5):</span> [${
													player.library
														.slice(0, 5)
														.map((c) => `\n  <span class="ds">"${c.name}"</span>`)
														.join(',') || ''
												}
]
<span class="dk">keptHand:</span> <span class="db">${player.keptHand}</span>`}</code
											></pre>
									</div>
								</div>
							{/each}
						</section>

						<!-- Zones -->
						<section class="debug-section">
							<div class="debug-section-header">
								<span>Zones</span>
							</div>
							<div class="debug-zones-grid">
								<div class="debug-zone">
									<h4>🏟️ Battlefield ({battlefield.length})</h4>
									<div class="debug-code small">
										<pre><code
												>{battlefield.length > 0
													? JSON.stringify(
															battlefield.map((c) => ({
																id: c.id,
																name: c.name,
																controller: c.controllerId,
																tapped: c.tapped
															})),
															null,
															2
														)
													: '[]'}</code
											></pre>
									</div>
								</div>
								<div class="debug-zone">
									<h4>🚫 Exile ({exile.length})</h4>
									<div class="debug-code small">
										<pre><code
												>{exile.length > 0
													? JSON.stringify(
															exile.map((c) => ({ id: c.id, name: c.name })),
															null,
															2
														)
													: '[]'}</code
											></pre>
									</div>
								</div>
							</div>
						</section>
					</main>
				</div>
			</div>
		{/if}

		<!-- Drag Ghost -->
		{#if isDragging && dragCardName}
			{@const dragImageUrl = getScryfallImageUrl(dragCardName, 'small')}
			<div class="drag-ghost" style="left: {dragPos.x}px; top: {dragPos.y}px;">
				<div class="drag-ghost-card" class:valid={isOverValidDrop}>
					{#if dragImageUrl}
						<img src={dragImageUrl} alt={dragCardName} class="drag-ghost-image" draggable="false" />
					{:else}
						<span class="drag-ghost-name">{dragCardName}</span>
					{/if}
				</div>
			</div>
		{/if}
	{/if}
</div>

<style>
	.playtest-container {
		position: fixed;
		inset: 0;
		background: #0a0d12;
		color: white;
		display: flex;
		flex-direction: column;
		overflow: hidden;
	}

	.loading-overlay,
	.error-overlay {
		position: absolute;
		inset: 0;
		display: flex;
		flex-direction: column;
		align-items: center;
		justify-content: center;
		background: #0a0d12;
		gap: 1rem;
	}

	.spinner {
		width: 48px;
		height: 48px;
		border: 4px solid #2a2a2a;
		border-top-color: #667eea;
		border-radius: 50%;
		animation: spin 0.8s linear infinite;
	}

	@keyframes spin {
		to {
			transform: rotate(360deg);
		}
	}

	.error-icon {
		font-size: 4rem;
	}

	.btn-primary {
		padding: 0.75rem 1.5rem;
		background: #667eea;
		color: white;
		border: none;
		border-radius: 6px;
		font-weight: 600;
		cursor: pointer;
	}

	.playtest-header {
		display: flex;
		align-items: center;
		justify-content: space-between;
		margin-top: 0;
		padding: 1rem;
		background: rgba(26, 31, 46, 0.95);
		border-bottom: 1px solid #2a3441;
		gap: 1rem;
		transition: all 0.3s ease;
	}

	.playtest-header.hidden {
		transform: translateY(-100%);
		opacity: 0;
		pointer-events: none;
	}

	/* Hamburger Menu Button */
	.menu-toggle-btn {
		position: fixed;
		top: 1rem;
		right: 1rem;
		z-index: 200;
		width: 48px;
		height: 48px;
		background: rgba(26, 31, 46, 0.95);
		border: 1px solid #2a3441;
		border-radius: 8px;
		color: #f4f4f5;
		cursor: pointer;
		display: flex;
		align-items: center;
		justify-content: center;
		transition: all 0.2s ease;
		backdrop-filter: blur(8px);
	}

	.menu-toggle-btn:hover {
		background: rgba(102, 126, 234, 0.2);
		border-color: #667eea;
		transform: scale(1.05);
	}

	/* Menu Backdrop */
	.menu-backdrop {
		position: fixed;
		inset: 0;
		background: rgba(0, 0, 0, 0.5);
		z-index: 149;
		animation: fadeIn 0.3s ease;
		cursor: pointer;
	}

	@keyframes fadeIn {
		from {
			opacity: 0;
		}
		to {
			opacity: 1;
		}
	}

	/* Menu Overlay */
	.menu-overlay {
		position: fixed;
		top: 0;
		right: 0;
		width: 400px;
		max-width: 90vw;
		height: 100vh;
		background: rgba(26, 31, 46, 0.98);
		border-left: 1px solid #2a3441;
		z-index: 150;
		padding: 0;
		overflow-y: auto;
		transform: translateX(100%);
		transition: transform 0.3s ease;
		backdrop-filter: blur(16px);
	}

	.menu-overlay.open {
		transform: translateX(0);
	}

	.menu-header {
		display: flex;
		align-items: center;
		justify-content: space-between;
		padding: 1.5rem;
		border-bottom: 1px solid #2a3441;
		position: sticky;
		top: 0;
		background: rgba(26, 31, 46, 0.98);
		backdrop-filter: blur(16px);
		z-index: 10;
	}

	.menu-header h2 {
		font-size: 1.25rem;
		font-weight: 700;
		margin: 0;
		color: #f4f4f5;
	}

	.menu-close-btn {
		width: 36px;
		height: 36px;
		border-radius: 6px;
		background: rgba(63, 63, 70, 0.4);
		border: 1px solid rgba(63, 63, 70, 0.6);
		color: #f4f4f5;
		cursor: pointer;
		display: flex;
		align-items: center;
		justify-content: center;
		transition: all 0.2s ease;
	}

	.menu-close-btn:hover {
		background: rgba(239, 68, 68, 0.2);
		border-color: #ef4444;
	}

	.menu-content {
		padding: 1.5rem;
		display: flex;
		flex-direction: column;
		gap: 1.5rem;
	}

	.menu-section {
		display: flex;
		flex-direction: column;
		gap: 0.75rem;
	}

	.menu-section-title {
		font-size: 0.75rem;
		font-weight: 700;
		text-transform: uppercase;
		letter-spacing: 0.05em;
		color: #71717a;
		margin: 0;
	}

	.menu-section-content {
		display: flex;
		flex-direction: column;
		gap: 0.5rem;
	}

	.menu-label {
		font-size: 0.875rem;
		font-weight: 500;
		color: #a1a1aa;
		display: block;
		margin-bottom: 0.25rem;
	}

	.control-select {
		width: 100%;
		padding: 0.5rem 0.75rem;
		background: rgba(36, 40, 51, 0.8);
		border: 1px solid rgba(63, 63, 70, 0.5);
		border-radius: 6px;
		color: #f4f4f5;
		font-size: 0.875rem;
		cursor: pointer;
	}

	.menu-btn {
		width: 100%;
		padding: 0.625rem 0.75rem;
		background: rgba(36, 40, 51, 0.8);
		border: 1px solid rgba(63, 63, 70, 0.5);
		border-radius: 6px;
		color: #f4f4f5;
		font-size: 0.875rem;
		font-weight: 500;
		cursor: pointer;
		display: flex;
		align-items: center;
		gap: 0.5rem;
		justify-content: center;
		transition: all 0.2s ease;
	}

	.menu-btn:hover {
		background: rgba(63, 63, 70, 0.6);
		border-color: #667eea;
	}

	.menu-btn.primary {
		background: rgba(102, 126, 234, 0.2);
		border-color: #667eea;
		color: #667eea;
		font-weight: 600;
	}

	.menu-btn.primary:hover {
		background: rgba(102, 126, 234, 0.3);
	}

	.turn-info {
		display: flex;
		align-items: center;
		gap: 0.5rem;
		padding: 0.5rem 0.75rem;
		background: rgba(36, 40, 51, 0.6);
		border-radius: 6px;
		font-size: 0.875rem;
		color: #f4f4f5;
	}

	.turn-info .active-player {
		color: #a1a1aa;
	}

	.header-left,
	.header-right {
		display: flex;
		align-items: center;
		gap: 0.75rem;
	}

	.turn-indicator {
		display: inline-flex;
		align-items: center;
		gap: 0.375rem;
		padding: 0.375rem 0.625rem;
		border-radius: 999px;
		border: 1px solid rgba(148, 163, 184, 0.25);
		background: rgba(17, 24, 39, 0.5);
		color: #cbd5e1;
		font-size: 0.875rem;
		font-weight: 600;
		white-space: nowrap;
	}

	.turn-text {
		font-family: 'JetBrains Mono', monospace;
		font-size: 0.8125rem;
	}

	.btn-back {
		padding: 0.5rem 1rem;
		background: transparent;
		border: 1px solid #374151;
		color: #9ca3af;
		border-radius: 6px;
		cursor: pointer;
		transition: all 0.2s;
	}

	.btn-back:hover {
		background: #1f2937;
		color: #f8fafc;
	}

	.mode-badge {
		padding: 0.375rem 0.75rem;
		background: rgba(102, 126, 234, 0.15);
		border: 1px solid rgba(102, 126, 234, 0.3);
		border-radius: 6px;
		color: #667eea;
		font-size: 0.875rem;
		font-weight: 600;
	}

	.btn-sessions {
		padding: 0.5rem 1rem;
		background: transparent;
		border: 1px solid rgba(102, 126, 234, 0.3);
		color: #667eea;
		border-radius: 6px;
		font-size: 0.875rem;
		font-weight: 600;
		cursor: pointer;
		transition: all 0.2s;
	}

	.btn-sessions:hover {
		background: rgba(102, 126, 234, 0.1);
		border-color: rgba(102, 126, 234, 0.5);
	}

	.playtest-controls {
		display: flex;
		align-items: center;
		gap: 0.75rem;
	}

	.playtest-controls label {
		font-size: 0.875rem;
		color: #94a3b8;
	}

	.player-select {
		padding: 0.5rem 1rem;
		background: #1a1f2e;
		border: 1px solid #2a3441;
		border-radius: 6px;
		color: #f8fafc;
		font-weight: 600;
		cursor: pointer;
	}

	.btn-toggle,
	.btn-action {
		padding: 0.5rem 1rem;
		background: #667eea;
		color: white;
		border: none;
		border-radius: 6px;
		font-weight: 600;
		cursor: pointer;
		transition: background 0.2s;
	}

	.btn-toggle:hover,
	.btn-action:hover {
		background: #5568d3;
	}

	.all-hands-overlay {
		position: fixed;
		top: 5rem;
		right: 1rem;
		width: 300px;
		max-height: calc(100vh - 6rem);
		overflow-y: auto;
		background: rgba(26, 31, 46, 0.98);
		border: 1px solid #2a3441;
		border-radius: 8px;
		padding: 0.75rem;
		z-index: 100;
		display: flex;
		flex-direction: column;
		gap: 0.75rem;
	}

	.player-hand-compact {
		padding: 0.75rem;
		background: rgba(17, 24, 39, 0.8);
		border: 1px solid #2a3441;
		border-radius: 6px;
	}

	.player-hand-compact.active {
		border-color: #667eea;
		background: rgba(102, 126, 234, 0.1);
	}

	.compact-header {
		display: flex;
		justify-content: space-between;
		align-items: center;
		margin-bottom: 0.5rem;
	}

	.player-name-compact {
		font-weight: 600;
		font-size: 0.875rem;
	}

	.life-compact {
		font-size: 0.875rem;
		color: #94a3b8;
	}

	.cards-compact {
		display: flex;
		flex-wrap: wrap;
		gap: 0.25rem;
	}

	.card-mini {
		width: 40px;
		height: 56px;
		border-radius: 3px;
		overflow: hidden;
	}

	.card-mini img {
		width: 100%;
		height: 100%;
		object-fit: cover;
	}

	.game-layout {
		flex: 1;
		display: flex;
		flex-direction: column;
		padding: 0.75rem;
		gap: 0.75rem;
		overflow: hidden;
	}

	.opponent-section {
		display: flex;
		flex-direction: column;
		gap: 0;
		flex: 1;
	}

	.opponent-info-overlay {
		position: absolute;
		top: 0.5rem;
		left: 0.5rem;
		right: 0.5rem;
		display: flex;
		justify-content: space-between;
		align-items: center;
		background: rgba(26, 31, 46, 0.9);
		backdrop-filter: blur(12px);
		border: 1px solid rgba(42, 52, 65, 0.6);
		border-radius: 8px;
		padding: 0.5rem 0.75rem;
		z-index: 10;
		gap: 0.75rem;
	}

	.opponent-identity {
		display: flex;
		align-items: center;
		gap: 0.5rem;
		flex-shrink: 0;
	}

	.opponent-select {
		padding: 0.375rem 0.75rem;
		background: #1a1f2e;
		border: 1px solid #2a3441;
		border-radius: 6px;
		color: #f8fafc;
		font-weight: 600;
		font-size: 0.9375rem;
		cursor: pointer;
		transition: all 0.2s;
	}

	.opponent-select:hover {
		border-color: #667eea;
		background: rgba(102, 126, 234, 0.1);
	}

	.opponent-name-label {
		font-weight: 700;
		font-size: 0.9375rem;
		color: #f8fafc;
	}

	.opponent-stats-compact {
		display: flex;
		align-items: center;
		gap: 0.75rem;
		position: relative;
	}

	.opponent-counts {
		display: flex;
		gap: 0.5rem;
		align-items: center;
	}

	.opponent-count {
		font-size: 0.75rem;
		color: #a1a1aa;
		white-space: nowrap;
	}

	.opponent-menu {
		top: calc(100% + 8px);
		bottom: auto;
	}

	.battlefield-area {
		background: linear-gradient(135deg, #0d1117, #141821);
		border: 1px solid #2a3441;
		border-radius: 12px;
		padding: 1rem;
		overflow: auto;
		position: relative;
	}

	.my-battlefield {
		min-height: 473px;
	}

	.opponent-battlefield {
		position: relative;
		min-height: 150px;
		max-height: none;
		background: linear-gradient(135deg, #1a1217, #1c1428);
		border-color: rgba(200, 100, 100, 0.3);
		padding-top: 4rem;
		flex: 1;
	}

	.battlefield-area.drag-active {
		border-color: rgba(102, 126, 234, 0.5);
	}

	.battlefield-area.drag-valid {
		border-color: #22c55e;
		box-shadow: inset 0 0 40px rgba(34, 197, 94, 0.15);
	}

	.zone-label {
		font-size: 0.6875rem;
		color: #6b7280;
		text-transform: uppercase;
		letter-spacing: 0.5px;
		font-weight: 600;
		margin-bottom: 0.5rem;
		display: block;
	}

	.battlefield-rows {
		display: flex;
		flex-direction: column;
		gap: 0.5rem;
	}

	.battlefield-row {
		display: flex;
		flex-wrap: wrap;
		gap: 0.5rem;
	}

	.battlefield-row--lands {
		margin-top: 0.25rem;
		padding-top: 0.5rem;
		border-top: 1px dashed rgba(148, 163, 184, 0.25);
	}

	/* Battlefield content wrapper - flexbox for main battlefield + command zone */
	.battlefield-content-wrapper {
		display: flex;
		gap: 1rem;
		height: 100%;
		align-items: flex-start;
	}

	.battlefield-main {
		flex: 1;
		min-width: 0;
		display: flex;
		flex-direction: column;
		gap: 0.5rem;
	}

	/* Command zone (Commander) - right side */
	.command-zone {
		width: fit-content;
		flex-shrink: 0;
		padding: 0.75rem;
		background: rgba(26, 31, 46, 0.4);
		border: 1px solid rgba(102, 126, 234, 0.3);
		border-radius: 8px;
		display: flex;
		align-items: center;
		flex-direction: column;
		gap: 0.5rem;
		transition: all 0.2s ease;
	}

	.command-zone.drag-valid {
		border-color: #22c55e;
		background: rgba(34, 197, 94, 0.1);
		box-shadow: 0 0 0 2px rgba(34, 197, 94, 0.2);
	}

	/* Opponent command zone - slightly different styling */
	.opponent-command-zone {
		background: rgba(26, 18, 23, 0.4);
		border-color: rgba(200, 100, 100, 0.3);
	}

	.command-cards {
		display: flex;
		flex-direction: column;
		gap: 0.5rem;
	}

	.command-zone-empty {
		display: flex;
		flex-direction: column;
		align-items: center;
		justify-content: center;
		padding: 1.5rem 0.5rem;
		border: 1px dashed rgba(102, 126, 234, 0.3);
		border-radius: 6px;
		background: rgba(26, 31, 46, 0.2);
		min-height: 120px;
	}

	.zone-empty-text {
		font-size: 0.75rem;
		color: #71717a;
		text-transform: uppercase;
		letter-spacing: 0.05em;
		font-weight: 600;
	}

	.zone-empty-hint {
		font-size: 0.6875rem;
		color: #52525b;
		margin-top: 0.25rem;
		font-style: italic;
	}

	.command-card-wrapper {
		display: flex;
		cursor: grab;
	}

	.battlefield-card-wrapper {
		cursor: grab;
		user-select: none;
		transition: filter 0.2s;
	}

	.battlefield-card-wrapper.readonly {
		cursor: default;
		opacity: 0.9;
	}

	.battlefield-card-wrapper.is-hovered {
		filter: drop-shadow(0 0 4px rgba(100, 200, 255, 0.5));
	}

	.empty-battlefield {
		color: #4b5563;
		font-style: italic;
		font-size: 0.875rem;
		padding: 2rem;
		text-align: center;
	}

	.player-info-row {
		display: flex;
		justify-content: space-between;
		align-items: center;
		gap: 0.75rem;
		padding: 0.5rem 0.75rem;
		background: rgba(26, 31, 46, 0.9);
		border-radius: 8px;
		border: 1px solid #2a3441;
	}

	.player-identity {
		display: flex;
		align-items: center;
		gap: 0.5rem;
	}

	.player-name {
		font-weight: 700;
		font-size: 0.9375rem;
	}

	.player-stats-inline {
		display: flex;
		align-items: center;
		gap: 0.5rem;
		position: relative;
		margin-left: auto;
		margin-right: 0.5rem;
	}

	.life-group {
		display: flex;
		align-items: center;
		gap: 0.125rem;
	}

	.stat-btn {
		width: 24px;
		height: 24px;
		border: 1px solid rgba(63, 63, 70, 0.4);
		border-radius: 4px;
		background: rgba(36, 40, 51, 0.6);
		color: #a1a1aa;
		font-size: 14px;
		font-weight: 600;
		cursor: pointer;
		transition: all 0.15s ease;
		display: flex;
		align-items: center;
		justify-content: center;
	}

	.stat-btn:hover {
		background: rgba(63, 63, 70, 0.6);
		color: #f4f4f5;
	}

	.stat-display {
		display: flex;
		align-items: center;
		gap: 0.25rem;
		padding: 0.25rem 0.5rem;
		border-radius: 4px;
		font-size: 0.8125rem;
		font-weight: 600;
		background: transparent;
		border: none;
		cursor: pointer;
	}

	.stat-display:hover {
		background: rgba(63, 63, 70, 0.3);
	}

	.stat-icon {
		font-size: 0.75rem;
	}

	.stat-value {
		font-family: 'JetBrains Mono', monospace;
		min-width: 20px;
		text-align: center;
		color: #f4f4f5;
	}

	.quick-menu {
		position: absolute;
		bottom: calc(100% + 8px);
		left: 50%;
		transform: translateX(-50%);
		min-width: 200px;
		background: rgba(18, 20, 26, 0.98);
		border: 1px solid rgba(63, 63, 70, 0.6);
		border-radius: 8px;
		padding: 0.75rem;
		box-shadow: 0 -8px 24px rgba(0, 0, 0, 0.5);
		display: flex;
		flex-direction: column;
		gap: 0.625rem;
		z-index: 100;
	}

	.menu-section {
		display: flex;
		flex-direction: column;
		gap: 0.25rem;
	}

	.menu-label {
		font-size: 0.625rem;
		font-weight: 700;
		text-transform: uppercase;
		letter-spacing: 0.05em;
		color: #71717a;
	}

	.menu-row {
		display: flex;
		gap: 0.25rem;
		align-items: center;
	}

	.menu-row button {
		padding: 0.375rem 0.5rem;
		border: 1px solid rgba(63, 63, 70, 0.5);
		border-radius: 4px;
		background: rgba(36, 40, 51, 0.8);
		color: #a1a1aa;
		font-size: 0.75rem;
		font-weight: 500;
		cursor: pointer;
	}

	.menu-value {
		min-width: 24px;
		text-align: center;
		font-weight: 600;
		color: #f4f4f5;
	}

	.menu-close {
		position: absolute;
		top: 0.375rem;
		right: 0.375rem;
		width: 20px;
		height: 20px;
		display: flex;
		align-items: center;
		justify-content: center;
		background: transparent;
		border: none;
		color: #71717a;
		font-size: 0.75rem;
		cursor: pointer;
		border-radius: 4px;
	}

	.player-zones {
		display: flex;
		gap: 0.5rem;
		align-items: center;
	}

	.graveyard-drop-zone,
	.exile-drop-zone {
		transition: all 0.2s ease;
		border-radius: 6px;
		min-width: 70px;
		min-height: 32px;
	}

	.library-drop-zone {
		transition: all 0.2s ease;
		border-radius: 6px;
		min-width: 120px;
		min-height: 32px;
	}

	.hand-area {
		flex-shrink: 0;
		transition: all 0.2s ease;
		border-radius: 8px;
	}

	.hand-area.drag-valid {
		outline: 2px solid #22c55e;
		background: rgba(34, 197, 94, 0.1);
	}

	.drag-ghost {
		position: fixed;
		pointer-events: none;
		z-index: 10000;
		transform: translate(-50%, -60%);
	}

	.drag-ghost-card {
		width: 80px;
		height: 112px;
		background: linear-gradient(135deg, #1a1f2e, #0d1117);
		border: 2px solid #667eea;
		border-radius: 6px;
		display: flex;
		align-items: center;
		justify-content: center;
		overflow: hidden;
		box-shadow: 0 15px 40px rgba(0, 0, 0, 0.6);
		opacity: 0.95;
		transform: scale(1.1) rotate(-5deg);
		transition: all 0.15s;
	}

	.drag-ghost-card.valid {
		border-color: #22c55e;
		box-shadow:
			0 15px 40px rgba(0, 0, 0, 0.6),
			0 0 30px rgba(34, 197, 94, 0.5);
		transform: scale(1.15) rotate(0deg);
	}

	.drag-ghost-image {
		width: 100%;
		height: 100%;
		object-fit: cover;
		border-radius: 4px;
	}

	.drag-ghost-name {
		font-size: 0.625rem;
		font-weight: 600;
		color: white;
		text-align: center;
		padding: 0.25rem;
	}

	.drop-hint {
		color: #667eea;
		animation: drop-hint-pulse 1.5s ease-in-out infinite;
	}

	@keyframes drop-hint-pulse {
		0%,
		100% {
			opacity: 0.7;
		}
		50% {
			opacity: 1;
		}
	}

	/* Debug Overlay Styles */
	.debug-overlay {
		position: fixed;
		inset: 0;
		top: 9rem;
		background: #2d2d2d;
		z-index: 2000;
		overflow: hidden;
		display: flex;
		flex-direction: column;
		font-family: 'Fira Code', 'JetBrains Mono', 'Consolas', monospace;
	}

	.debug-modal {
		flex: 1;
		display: flex;
		flex-direction: column;
		overflow: hidden;
	}

	.debug-header {
		display: flex;
		justify-content: space-between;
		align-items: center;
		padding: 0.75rem 1rem;
		background: #1a1a1a;
		border-bottom: 1px solid #444;
		flex-shrink: 0;
	}

	.debug-header-left {
		display: flex;
		align-items: center;
		gap: 1rem;
	}

	.debug-header h2 {
		font-size: 1rem;
		font-weight: 600;
		margin: 0;
		color: #00ff00;
	}

	.debug-status {
		font-size: 0.75rem;
		padding: 0.25rem 0.5rem;
		border-radius: 999px;
		background: #003a00;
		color: #00ff00;
		border: 1px solid #00ff00;
	}

	.debug-header-right {
		display: flex;
		align-items: center;
		gap: 1rem;
	}

	.debug-close {
		width: 32px;
		height: 32px;
		border-radius: 6px;
		background: #333;
		border: 1px solid #555;
		color: #fff;
		font-size: 1rem;
		cursor: pointer;
		transition: all 0.2s;
	}

	.debug-close:hover {
		background: #ef4444;
		border-color: #ef4444;
	}

	.debug-content {
		flex: 1;
		overflow-y: auto;
		padding: 1rem;
	}

	.debug-section {
		margin-bottom: 1rem;
		background: #1e1e1e;
		border: 1px solid #444;
		border-radius: 8px;
		overflow: hidden;
	}

	.debug-section-header {
		display: flex;
		align-items: center;
		justify-content: space-between;
		gap: 0.5rem;
		width: 100%;
		padding: 0.75rem 1rem;
		background: #252525;
		color: #00ff00;
		font-family: inherit;
		font-size: 0.875rem;
		text-align: left;
	}

	.debug-code {
		background: #0d0d0d;
		padding: 0.75rem 1rem;
		overflow-x: auto;
		max-height: 300px;
		overflow-y: auto;
	}

	.debug-code.small {
		max-height: 150px;
	}

	.debug-code pre {
		margin: 0;
		font-size: 0.75rem;
		line-height: 1.4;
	}

	.debug-code code {
		color: #00ff00;
	}

	.debug-player {
		border-top: 1px solid #333;
	}

	.debug-player-header {
		display: flex;
		align-items: center;
		gap: 0.5rem;
		padding: 0.5rem 1rem;
		background: #1a1a1a;
	}

	.debug-badge {
		font-size: 0.625rem;
		padding: 0.125rem 0.375rem;
		border-radius: 4px;
		background: #333;
		color: #888;
	}

	.debug-badge.local {
		background: #003a00;
		color: #00ff00;
	}

	.debug-zones-grid {
		display: grid;
		grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
		gap: 1px;
		background: #333;
	}

	.debug-zone {
		background: #1e1e1e;
	}

	.debug-zone h4 {
		font-size: 0.75rem;
		font-weight: 600;
		margin: 0;
		padding: 0.5rem 0.75rem;
		background: #252525;
		color: #ffff00;
		border-bottom: 1px solid #333;
	}

	.btn-debug {
		width: 40px;
		height: 40px;
		display: flex;
		align-items: center;
		justify-content: center;
		background: rgba(0, 255, 0, 0.1);
		border: 1px solid rgba(0, 255, 0, 0.3);
		border-radius: 6px;
		color: #00ff00;
		cursor: pointer;
		transition: all 0.2s;
		font-size: 1.25rem;
	}

	.btn-debug:hover {
		background: rgba(0, 255, 0, 0.2);
		border-color: rgba(0, 255, 0, 0.5);
	}

	/* Syntax highlighting */
	:global(.debug-code .dk) {
		color: #9cdcfe;
	}
	:global(.debug-code .ds) {
		color: #ce9178;
	}
	:global(.debug-code .dn) {
		color: #b5cea8;
	}
	:global(.debug-code .db) {
		color: #569cd6;
	}
	:global(.debug-code .dc) {
		color: #6a9955;
		font-style: italic;
	}

	/* Game State Log Styles */
	.debug-log-container {
		max-height: 400px;
		overflow-y: auto;
		background: #0d0d0d;
		padding: 0.75rem 1rem;
	}

	.debug-log-empty {
		color: #6b7280;
		font-style: italic;
		padding: 1rem;
		text-align: center;
	}

	.debug-log-entries {
		display: flex;
		flex-direction: column;
		gap: 0.25rem;
	}

	.debug-log-entry {
		display: flex;
		align-items: center;
		gap: 0.75rem;
		padding: 0.5rem;
		font-size: 0.75rem;
		line-height: 1.4;
		border-bottom: 1px solid rgba(255, 255, 255, 0.05);
		font-family: 'JetBrains Mono', monospace;
	}

	.debug-log-entry:last-child {
		border-bottom: none;
	}

	.debug-log-time {
		color: #6b7280;
		min-width: 70px;
		font-size: 0.6875rem;
	}

	.debug-log-turn {
		color: #9cdcfe;
		min-width: 30px;
		font-weight: 600;
	}

	.debug-log-kind {
		color: #ce9178;
		min-width: 80px;
		text-transform: uppercase;
		font-size: 0.6875rem;
		font-weight: 600;
	}

	.debug-log-message {
		color: #d4d4d4;
		flex: 1;
	}

	.debug-copy-btn {
		display: flex;
		align-items: center;
		gap: 0.375rem;
		padding: 0.375rem 0.75rem;
		background: rgba(102, 126, 234, 0.2);
		border: 1px solid rgba(102, 126, 234, 0.4);
		border-radius: 4px;
		color: #667eea;
		font-size: 0.75rem;
		font-weight: 600;
		cursor: pointer;
		transition: all 0.2s;
	}

	.debug-copy-btn:hover {
		background: rgba(102, 126, 234, 0.3);
		border-color: rgba(102, 126, 234, 0.6);
	}

	.debug-copy-btn span {
		font-family: inherit;
	}

	/* Session Picker Styles */
	.session-picker-overlay {
		position: absolute;
		inset: 0;
		display: flex;
		align-items: center;
		justify-content: center;
		background: #0a0d12;
		padding: 2rem;
	}

	.session-picker-modal {
		background: rgba(26, 31, 46, 0.98);
		border: 1px solid #2a3441;
		border-radius: 12px;
		padding: 2rem;
		max-width: 700px;
		width: 100%;
		max-height: 80vh;
		display: flex;
		flex-direction: column;
		gap: 1.5rem;
	}

	.session-picker-modal h2 {
		margin: 0;
		color: #f8fafc;
		font-size: 1.5rem;
	}

	.session-picker-hint {
		margin: 0;
		color: #94a3b8;
		font-size: 0.875rem;
	}

	.sessions-list {
		display: flex;
		flex-direction: column;
		gap: 0.75rem;
		max-height: 400px;
		overflow-y: auto;
		padding: 0.5rem;
		margin: -0.5rem;
	}

	.session-card {
		display: flex;
		justify-content: space-between;
		align-items: center;
		padding: 1rem;
		background: rgba(17, 24, 39, 0.8);
		border: 1px solid #2a3441;
		border-radius: 8px;
		transition: all 0.2s;
		gap: 1rem;
	}

	.session-card:hover {
		border-color: #667eea;
		background: rgba(102, 126, 234, 0.1);
	}

	.session-info {
		flex: 1;
		display: flex;
		flex-direction: column;
		gap: 0.25rem;
	}

	.session-label {
		font-weight: 600;
		color: #f8fafc;
		font-size: 0.9375rem;
	}

	.session-meta {
		font-size: 0.8125rem;
		color: #94a3b8;
	}

	.session-actions {
		display: flex;
		gap: 0.5rem;
		align-items: center;
	}

	.btn-restore {
		padding: 0.5rem 1rem;
		background: #667eea;
		color: white;
		border: none;
		border-radius: 6px;
		font-weight: 600;
		cursor: pointer;
		transition: background 0.2s;
	}

	.btn-restore:hover {
		background: #5568d3;
	}

	.btn-delete {
		width: 32px;
		height: 32px;
		display: flex;
		align-items: center;
		justify-content: center;
		background: rgba(239, 68, 68, 0.1);
		border: 1px solid rgba(239, 68, 68, 0.3);
		border-radius: 6px;
		color: #ef4444;
		cursor: pointer;
		transition: all 0.2s;
		font-size: 1rem;
	}

	.btn-delete:hover {
		background: rgba(239, 68, 68, 0.2);
		border-color: rgba(239, 68, 68, 0.5);
	}

	.no-sessions {
		color: #94a3b8;
		text-align: center;
		padding: 2rem;
		font-style: italic;
	}

	.session-picker-actions {
		display: flex;
		gap: 0.75rem;
		justify-content: flex-end;
		padding-top: 1rem;
		border-top: 1px solid #2a3441;
	}

	/* Responsive Design */
	@media (max-width: 1200px) {
		.command-zone,
		.opponent-command-zone {
			min-width: 150px;
			max-width: 150px;
		}

		.menu-overlay {
			width: 350px;
		}
	}

	@media (max-width: 768px) {
		.battlefield-content-wrapper {
			flex-direction: column;
		}

		.command-zone,
		.opponent-command-zone {
			width: 100%;
			max-width: 100%;
		}

		.command-cards {
			flex-direction: row;
			flex-wrap: wrap;
		}

		.menu-overlay {
			width: 100%;
		}

		.opponent-info-overlay {
			flex-direction: column;
			align-items: flex-start;
			gap: 0.5rem;
		}

		.opponent-stats-compact {
			width: 100%;
			justify-content: space-between;
		}
	}
</style>
