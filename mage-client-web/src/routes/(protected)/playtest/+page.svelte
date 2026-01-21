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
	import PlaytestLibrarySearch from '$lib/components/game/PlaytestLibrarySearch.svelte';
	import TokenCreator from '$lib/components/game/TokenCreator.svelte';
	import CreateTokenDialog from '$lib/components/game/CreateTokenDialog.svelte';
	import CounterDialog from '$lib/components/game/CounterDialog.svelte';
	import MulliganDialog from '$lib/components/game/MulliganDialog.svelte';
	import DeckContextMenu from '$lib/components/game/DeckContextMenu.svelte';
	import NumberInputDialog from '$lib/components/game/NumberInputDialog.svelte';
	import ScryDialog from '$lib/components/game/ScryDialog.svelte';
	import RevealTopDialog from '$lib/components/game/RevealTopDialog.svelte';
	import PlaytestHeader from '$lib/components/game/PlaytestHeader.svelte';
	import PlayerInfoRow from '$lib/components/game/PlayerInfoRow.svelte';
	import OpponentSection from '$lib/components/game/OpponentSection.svelte';
	import BattlefieldArea from '$lib/components/game/BattlefieldArea.svelte';
	import type { MenuAction } from '$lib/components/game/DeckContextMenu.svelte';
	import type { ScrySession } from '$lib/stores/playtest-game';
	import Keyboard from '@lucide/svelte/icons/keyboard';
	import Clock from '@lucide/svelte/icons/clock';
	import Copy from '@lucide/svelte/icons/copy';
	import Menu from '@lucide/svelte/icons/menu';
	import X from '@lucide/svelte/icons/x';
	import ArrowDown from '@lucide/svelte/icons/arrow-down';
	import RotateCcw from '@lucide/svelte/icons/rotate-ccw';
	import Shuffle from '@lucide/svelte/icons/shuffle';
	import Search from '@lucide/svelte/icons/search';
	import Plus from '@lucide/svelte/icons/plus';
	import Hand from '@lucide/svelte/icons/hand';
	import BookOpen from '@lucide/svelte/icons/book-open';
	import Eye from '@lucide/svelte/icons/eye';
	import EyeOff from '@lucide/svelte/icons/eye-off';
	import Heart from '@lucide/svelte/icons/heart';
	import Skull from '@lucide/svelte/icons/skull';
	import GalleryVertical from '@lucide/svelte/icons/gallery-vertical';
	import FastForward from '@lucide/svelte/icons/fast-forward';
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

	// Deck context menu and dialog state
	let showDeckContextMenu = $state(false);
	let deckContextMenuPosition = $state<{ x: number; y: number }>({ x: 0, y: 0 });
	let showNumberInputDialog = $state(false);
	let numberInputDialogConfig = $state<{
		title: string;
		defaultValue: number;
		min: number;
		max: number;
		onConfirm: (value: number) => void;
	} | null>(null);
	let showScryDialog = $state(false);
	let currentScrySession = $state<ScrySession | null>(null);
	let showRevealTopDialog = $state(false);
	let revealedCards = $state<import('$lib/generated/mage/v1/models').CardView[]>([]);

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
	const selectedOpponent = $derived.by(() => {
		if (otherPlayers.length === 0) return null;
		if (!selectedOpponentId || !otherPlayers.find((p) => p.playerId === selectedOpponentId)) {
			// Auto-select first opponent
			return otherPlayers[0];
		}
		return otherPlayers.find((p) => p.playerId === selectedOpponentId) || otherPlayers[0];
	});

	// Split battlefield by controller
	const myBattlefield = $derived(battlefield.filter((c) => c.controllerId === activeControlSeat));
	const opponentBattlefield = $derived.by(() => {
		const opponent = selectedOpponent;
		return opponent ? battlefield.filter((c) => c.controllerId === opponent.playerId) : [];
	});

	function isLandPermanent(cardType?: string | null): boolean {
		// Mage type strings are typically like: "Land", "Legendary Land", "Artifact Land", etc.
		return !!cardType && /\bland\b/i.test(cardType);
	}

	// Split battlefield rows: nonlands (top) + lands (bottom)
	const myBattlefieldNonlands = $derived(myBattlefield.filter((c) => !isLandPermanent(c.type)));
	const myBattlefieldLands = $derived(myBattlefield.filter((c) => isLandPermanent(c.type)));
	const opponentBattlefieldNonlands = $derived.by(() =>
		opponentBattlefield.filter((c) => !isLandPermanent(c.type))
	);
	const opponentBattlefieldLands = $derived.by(() =>
		opponentBattlefield.filter((c) => isLandPermanent(c.type))
	);

	// My cards (from controlling player perspective)
	const myGrave = $derived(me?.graveyard || []);
	const myMana = $derived(
		me?.manaPool || { white: 0, blue: 0, black: 0, red: 0, green: 0, colorless: 0 }
	);

	// Reactive card lookup for counter dialog
	const selectedCardForCountersData = $derived.by(() => {
		// 1. Capture the value in a local variable for "Type Narrowing"
		const currentId = selectedCardForCounters?.id;

		// 2. If no ID exists, exit early
		if (!currentId) return null;

		// 3. Search through collections
		// Using the local currentId ensures TS knows it's a string/number, not null
		const card =
			$playtestBattlefield.find((c) => c.id === currentId) ||
			me?.hand.find((c) => c.id === currentId) ||
			me?.graveyard.find((c) => c.id === currentId) ||
			null;

		console.log(
			'[selectedCardForCountersData] Re-evaluated.',
			`Card: ${card?.name}`,
			`Counters: ${card?.counters}`
		);

		return card;
	});

	// Hovered card
	const hoveredCard = $derived(
		hoveredCardId ? battlefield.find((c) => c.id === hoveredCardId) : null
	);

	const activePlayerName = $derived.by(() => {
		return players.find((p) => p.playerId === $playtestGameStore.activePlayerId)?.name ?? '';
	});

	// Playtest store "turn" currently increments each time we advance priority seat.
	// For display, we want "Turn 2" only when player 1 starts their second turn (i.e. rounds).
	const turnNumber = $derived.by(() => {
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
	const opponentCommandCards = $derived.by(() => {
		const opponent = selectedOpponent;
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
	 * Deck context menu handlers
	 */
	function handleDeckContextMenu(event: MouseEvent): void {
		if (!me) return;
		deckContextMenuPosition = { x: event.clientX, y: event.clientY };
		showDeckContextMenu = true;
	}

	function handleDrawN(count: number): void {
		if (!me) return;
		playtestGameStore.drawCards(me.playerId, count);
		toast.success(`Drew ${count} card(s)`);
	}

	function handleMill(count: number): void {
		if (!me) return;
		playtestGameStore.millCards(me.playerId, count);
		toast.success(`Milled ${count} card(s)`);
	}

	function handleScry(count: number): void {
		if (!me) return;
		const session = playtestGameStore.scryCards(me.playerId, count);
		if (session) {
			currentScrySession = session;
			showScryDialog = true;
		} else {
			toast.error('No cards to scry');
		}
	}

	function handleScryComplete(
		keepOnTop: import('$lib/generated/mage/v1/models').CardView[],
		putToBottom: import('$lib/generated/mage/v1/models').CardView[]
	): void {
		if (!me || !currentScrySession) return;
		const scryCount = currentScrySession.cards.length;
		playtestGameStore.applyScryDecision(me.playerId, scryCount, keepOnTop, putToBottom);
		showScryDialog = false;
		currentScrySession = null;
		toast.success(`Scry ${scryCount} complete`);
	}

	function handleRevealTop(count: number): void {
		if (!me) return;
		const cards = playtestGameStore.revealTopCards(me.playerId, count);
		revealedCards = cards;
		showRevealTopDialog = true;
	}

	function handleToggleRevealedTop(): void {
		if (!me) return;
		const willReveal = !me.revealedTopCard;
		playtestGameStore.setRevealedTop(me.playerId, willReveal);
		toast.info(willReveal ? 'Top card revealed permanently' : 'Top card hidden');
	}

	function showNumberInput(
		title: string,
		defaultValue: number,
		onConfirm: (value: number) => void
	): void {
		numberInputDialogConfig = {
			title,
			defaultValue,
			min: 1,
			max: 99,
			onConfirm: (value) => {
				onConfirm(value);
				showNumberInputDialog = false;
				numberInputDialogConfig = null;
			}
		};
		showNumberInputDialog = true;
	}

	const deckContextMenuActions = $derived<MenuAction[]>(
		!me
			? []
			: [
					{
						label: 'Draw Cards',
						submenu: [
							{ label: '1 Card', onClick: () => handleDrawN(1) },
							{ label: '2 Cards', onClick: () => handleDrawN(2) },
							{ label: '3 Cards', onClick: () => handleDrawN(3) },
							{ label: '7 Cards', onClick: () => handleDrawN(7) },
							{ label: 'Custom...', onClick: () => showNumberInput('Draw N Cards', 1, handleDrawN) }
						]
					},
					{
						label: 'Scry',
						submenu: [
							{ label: '1 Card', onClick: () => handleScry(1) },
							{ label: '2 Cards', onClick: () => handleScry(2) },
							{ label: '3 Cards', onClick: () => handleScry(3) },
							{ label: 'Custom...', onClick: () => showNumberInput('Scry N Cards', 1, handleScry) }
						]
					},
					{
						label: 'Mill Cards',
						submenu: [
							{ label: '1 Card', onClick: () => handleMill(1) },
							{ label: '2 Cards', onClick: () => handleMill(2) },
							{ label: '3 Cards', onClick: () => handleMill(3) },
							{ label: '5 Cards', onClick: () => handleMill(5) },
							{ label: 'Custom...', onClick: () => showNumberInput('Mill N Cards', 1, handleMill) }
						]
					},
					{ divider: true },
					{
						label: 'Reveal Top Card',
						onClick: () => handleRevealTop(1)
					},
					{
						label: me.revealedTopCard ? 'Hide Revealed Top' : 'Reveal Top Permanently',
						onClick: handleToggleRevealedTop
					},
					{ divider: true },
					{
						label: 'Search Library',
						onClick: () => {
							showDeckSearch = true;
						}
					},
					{
						label: 'Shuffle Library',
						onClick: handleShuffleLibrary
					}
				]
	);

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
	 * Build visible gamestate JSON for export
	 */
	function buildVisibleGamestate() {
		const state = $playtestGameStore;

		// Helper to find card name by ID
		const findCardName = (cardId: string): string => {
			const card =
				battlefield.find((c) => c.id === cardId) ||
				exile.find((c) => c.id === cardId) ||
				players
					.flatMap((p) => [...p.hand, ...p.graveyard, ...p.library])
					.find((c) => c.id === cardId);
			return card?.name || 'Unknown';
		};

		// Helper to extract card details
		const extractCardDetails = (card: import('$lib/generated/mage/v1/models').CardView) => {
			const details: any = {
				name: card.name,
				type: card.type || '',
				manaCost: card.manaCost || ''
			};

			if (card.power !== undefined && card.power !== null) details.power = card.power;
			if (card.toughness !== undefined && card.toughness !== null)
				details.toughness = card.toughness;
			if (card.abilities && card.abilities.length > 0) {
				details.abilities = card.abilities.map((a) => a.text || '').filter((t) => t);
			}
			if (card.counters && card.counters.length > 0) {
				details.counters = card.counters;
			}
			if (card.tapped) details.tapped = true;
			if (card.faceDown) details.faceDown = true;
			if (card.summoningSickness) details.summoningSickness = true;
			if (card.attachedTo && card.attachedTo.length > 0) {
				details.attachedTo = card.attachedTo.map((id) => findCardName(id));
			}

			return details;
		};

		// Get my player (active control seat)
		const myPlayer = players.find((p) => p.playerId === activeControlSeat);
		if (!myPlayer) return null;

		// Build my state
		const myState = {
			hand: myPlayer.hand.map(extractCardDetails),
			battlefield: battlefield
				.filter((c) => c.controllerId === activeControlSeat)
				.map(extractCardDetails),
			graveyard: myPlayer.graveyard.map(extractCardDetails),
			exile: exile
				.filter((c) => (c.controllerId || c.ownerId) === activeControlSeat)
				.map(extractCardDetails),
			command: commandCards
				.filter((c) => (c.controllerId || c.ownerId) === activeControlSeat)
				.map(extractCardDetails),
			manaPool: myPlayer.manaPool,
			life: myPlayer.life,
			poison: myPlayer.poison || 0,
			energy: myPlayer.energy || 0,
			libraryCount: myPlayer.libraryCount,
			libraryContents: myPlayer.library.map(extractCardDetails)
		};

		// Build opponents state
		const opponents = players
			.filter((p) => p.playerId !== activeControlSeat)
			.map((opponent) => ({
				name: opponent.name,
				battlefield: battlefield
					.filter((c) => c.controllerId === opponent.playerId)
					.map(extractCardDetails),
				graveyard: opponent.graveyard.map(extractCardDetails),
				exile: exile
					.filter((c) => (c.controllerId || c.ownerId) === opponent.playerId)
					.map(extractCardDetails),
				command: commandCards
					.filter((c) => (c.controllerId || c.ownerId) === opponent.playerId)
					.map(extractCardDetails),
				life: opponent.life,
				poison: opponent.poison || 0,
				energy: opponent.energy || 0,
				handCount: opponent.handCount,
				libraryCount: opponent.libraryCount
			}));

		// Build game state
		const gameState = {
			turnNumber: turnNumber,
			activePlayer: activePlayerName || '',
			currentPhase: 'N/A', // Playtest doesn't track phases
			stack: (state.stack || []).map(extractCardDetails)
		};

		// Get game log text
		const gameLogText = playtestGameStore.buildLogText(state);

		return {
			myState,
			opponents,
			gameState,
			gameLog: gameLogText
		};
	}

	/**
	 * Copy gamestate to clipboard
	 */
	async function copyGamestateToClipboard(): Promise<void> {
		const gamestate = buildVisibleGamestate();
		if (!gamestate) {
			toast.error('No gamestate available');
			return;
		}

		const json = JSON.stringify(gamestate, null, 2);

		try {
			if (typeof navigator !== 'undefined' && navigator.clipboard?.writeText) {
				await navigator.clipboard.writeText(json);
				toast.success('Gamestate copied to clipboard!');
				return;
			}
			// Fallback for older browsers
			const textarea = document.createElement('textarea');
			textarea.value = json;
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
				toast.success('Gamestate copied to clipboard!');
			} else {
				toast.error('Failed to copy gamestate');
			}
		} catch (err) {
			console.error('Failed to copy gamestate to clipboard:', err);
			toast.error('Failed to copy gamestate');
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
			case 'j':
				// Shift+J - Copy gamestate to clipboard
				if (event.shiftKey) {
					copyGamestateToClipboard();
					event.preventDefault();
				}
				break;
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
		<PlaytestHeader
			{players}
			{activeControlSeat}
			availableSessions={availableSessions.length}
			{turnNumber}
			{activePlayerName}
			{showAllHands}
			onBack={() => goto('/lobby')}
			onSessionsClick={() => {
				loadAvailableSessions();
				showSessionPicker = true;
			}}
			onSwitchPlayer={switchPlayer}
			onToggleAllHands={() => (showAllHands = !showAllHands)}
			onDrawCard={handleDrawCard}
			onUntapAll={handleUntapAll}
			onShuffleLibrary={handleShuffleLibrary}
			onSearchLibrary={() => (showDeckSearch = true)}
			onCreateToken={() => (showCreateTokenDialog = true)}
			onNextTurn={handleNextTurn}
			onShowKeyboardShortcuts={() => (showKeyboardShortcuts = true)}
			onShowDebug={() => (showDebugOverlay = true)}
			onToggleMenu={() => (showMenu = !showMenu)}
		/>

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
								{#if showAllHands}
									<EyeOff size={16} />
									Hide
								{:else}
									<Eye size={16} />
									Show
								{/if}
								All Hands
							</button>
						</div>
					</div>

					<!-- Turn Info Section -->
					<div class="menu-section">
						<h3 class="menu-section-title">Turn Info</h3>
						<div class="menu-section-content">
							<div class="turn-info">
								<Clock size={18} />
								<span>Turn {turnNumber}</span>
								{#if activePlayerName}
									<span class="active-player">· {activePlayerName}</span>
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
							<span class="life-compact"><Heart size={14} /> {player.life}</span>
						</div>
						<div class="cards-compact">
							{#each player.hand as card}
								<Card cardId={card.id} cardName={card.name} size="large" manaCost={card.manaCost} />
							{/each}
						</div>
					</div>
				{/each}
			</div>
		{/if}

		<!-- Main Game Area -->
		<main class="game-layout">
			<!-- Opponent Section(s) -->
			{#if otherPlayers.length === 1}
				<!-- 1v1 Layout: Single opponent -->
				{#if selectedOpponent}
					{@const opponent = selectedOpponent}
					<OpponentSection
						{opponent}
						{otherPlayers}
						battlefieldNonlands={opponentBattlefieldNonlands}
						battlefieldLands={opponentBattlefieldLands}
						commandCards={opponentCommandCards}
						{isCommanderGame}
						showLifeMenu={showOpponentLifeMenu}
						onSelectOpponent={(playerId) => (selectedOpponentId = playerId)}
						onLifeChange={handleLifeChange}
						onPoisonChange={handlePoisonChange}
						onToggleLifeMenu={() => (showOpponentLifeMenu = !showOpponentLifeMenu)}
						onCardContextMenu={(cardId, cardName) => {
							selectedCardForCounters = { id: cardId, name: cardName };
							showCounterDialog = true;
						}}
					/>
				{/if}
			{:else}
				<!-- Multiplayer (3-4 players): Grid on large screens, cycling on small -->
				<!-- Grid layout (shown on large screens) -->
				<div class="opponents-grid opponents-grid-large">
					{#each otherPlayers as opponent (opponent.playerId)}
						{@const oppBattlefield = battlefield.filter(
							(c) => c.controllerId === opponent.playerId
						)}
						{@const oppBattlefieldNonlands = oppBattlefield.filter((c) => !isLandPermanent(c.type))}
						{@const oppBattlefieldLands = oppBattlefield.filter((c) => isLandPermanent(c.type))}
						{@const oppCommandCards = commandCards.filter(
							(c) => (c.ownerId || c.controllerId) === opponent.playerId
						)}
						<OpponentSection
							{opponent}
							otherPlayers={[]}
							battlefieldNonlands={oppBattlefieldNonlands}
							battlefieldLands={oppBattlefieldLands}
							commandCards={oppCommandCards}
							{isCommanderGame}
							showLifeMenu={false}
							onSelectOpponent={undefined}
							onLifeChange={handleLifeChange}
							onPoisonChange={handlePoisonChange}
							onToggleLifeMenu={() => {}}
							onCardContextMenu={(cardId, cardName) => {
								selectedCardForCounters = { id: cardId, name: cardName };
								showCounterDialog = true;
							}}
						/>
					{/each}
				</div>
				<!-- Single opponent with cycling (shown on small screens) -->
				<div class="opponents-grid-small">
					{#if selectedOpponent}
						{@const opponent = selectedOpponent}
						<OpponentSection
							{opponent}
							{otherPlayers}
							battlefieldNonlands={opponentBattlefieldNonlands}
							battlefieldLands={opponentBattlefieldLands}
							commandCards={opponentCommandCards}
							{isCommanderGame}
							showLifeMenu={showOpponentLifeMenu}
							onSelectOpponent={(playerId) => (selectedOpponentId = playerId)}
							onLifeChange={handleLifeChange}
							onPoisonChange={handlePoisonChange}
							onToggleLifeMenu={() => (showOpponentLifeMenu = !showOpponentLifeMenu)}
							onCardContextMenu={(cardId, cardName) => {
								selectedCardForCounters = { id: cardId, name: cardName };
								showCounterDialog = true;
							}}
						/>
					{/if}
				</div>
			{/if}

			<!-- My Battlefield Area (Editable) -->
			<BattlefieldArea
				battlefieldNonlands={myBattlefieldNonlands}
				battlefieldLands={myBattlefieldLands}
				commandCards={myCommandCards}
				{isCommanderGame}
				{isDragging}
				{isOverValidDrop}
				{dropZone}
				{hoveredCardId}
				onCardClick={handleBattlefieldCardClick}
				onCardMouseDown={handleBattlefieldCardMouseDown}
				onCardContextMenu={(cardId, cardName) => {
					selectedCardForCounters = { id: cardId, name: cardName };
					showCounterDialog = true;
				}}
				onCommandCardMouseDown={handleCommandCardMouseDown}
				onCardHover={(cardId) => (hoveredCardId = cardId)}
				battlefieldDropZoneRef={(el) => (battlefieldDropZoneEl = el)}
				commandDropZoneRef={(el) => (commandDropZoneEl = el)}
			/>

			<!-- Player Info Row -->
			{#if me}
				<PlayerInfoRow
					player={{
						name: me.name,
						life: me.life,
						poison: me.poison,
						libraryCount: me.libraryCount
					}}
					graveyard={myGrave}
					{exile}
					mana={myMana}
					{showLifeMenu}
					onLifeChange={handleLifeChange}
					onPoisonChange={handlePoisonChange}
					onToggleLifeMenu={() => (showLifeMenu = !showLifeMenu)}
					onSearchLibrary={() => (showDeckSearch = true)}
					onDeckContextMenu={handleDeckContextMenu}
					libraryDropZoneRef={(el) => (libraryDropZoneEl = el)}
					graveyardDropZoneRef={(el) => (graveyardDropZoneEl = el)}
					exileDropZoneRef={(el) => (exileDropZoneEl = el)}
				/>
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
		{#if showCounterDialog && selectedCardForCounters && selectedCardForCountersData}
			<CounterDialog
				cardName={selectedCardForCountersData.name}
				cardId={selectedCardForCountersData.id}
				currentCounters={selectedCardForCountersData.counters}
				onAddCounter={(counterName, amount) => {
					const card = selectedCardForCountersData;
					playtestGameStore.addCounter(card.id, counterName, amount);
					syncPlaytestToGameStore();
				}}
				onRemoveCounter={(counterName, amount) => {
					const card = selectedCardForCountersData;
					playtestGameStore.removeCounter(card.id, counterName, amount);
					syncPlaytestToGameStore();
				}}
				onSetCounter={(counterName, amount) => {
					const card = selectedCardForCountersData;
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

		<!-- Deck Context Menu -->
		{#if showDeckContextMenu}
			<DeckContextMenu
				position={deckContextMenuPosition}
				deckCount={me?.libraryCount || 0}
				playerName={me?.name || 'You'}
				onClose={() => (showDeckContextMenu = false)}
				actions={deckContextMenuActions}
			/>
		{/if}

		<!-- Number Input Dialog -->
		{#if showNumberInputDialog && numberInputDialogConfig}
			<NumberInputDialog
				title={numberInputDialogConfig.title}
				defaultValue={numberInputDialogConfig.defaultValue}
				min={numberInputDialogConfig.min}
				max={numberInputDialogConfig.max}
				onConfirm={numberInputDialogConfig.onConfirm}
				onCancel={() => {
					showNumberInputDialog = false;
					numberInputDialogConfig = null;
				}}
			/>
		{/if}

		<!-- Scry Dialog -->
		{#if showScryDialog && currentScrySession}
			<ScryDialog
				cards={currentScrySession.cards}
				onComplete={handleScryComplete}
				onCancel={() => {
					showScryDialog = false;
					currentScrySession = null;
				}}
			/>
		{/if}

		<!-- Reveal Top Dialog -->
		{#if showRevealTopDialog}
			<RevealTopDialog
				cards={revealedCards}
				onClose={() => {
					showRevealTopDialog = false;
					revealedCards = [];
				}}
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

						<!-- Visible Gamestate (JSON Export) -->
						<section class="debug-section">
							<div class="debug-section-header">
								<span>Visible Gamestate</span>
								<button
									class="debug-copy-btn"
									onclick={copyGamestateToClipboard}
									title="Copy complete gamestate to clipboard"
									aria-label="Copy gamestate to clipboard"
								>
									<Copy size={16} aria-hidden="true" />
									<span>Copy JSON</span>
								</button>
							</div>
							<div class="debug-code">
								<pre><code
										>{@html `<span class="dc">// Complete game state in JSON format</span>
<span class="dc">// Includes all players, zones, cards, and game info</span>
<span class="dc">// Click "Copy JSON" to copy to clipboard</span>

<span class="dk">myState:</span> { hand, battlefield, graveyard, exile, command, manaPool, life, ... }
<span class="dk">opponents:</span> [{ name, battlefield, graveyard, life, ... }]
<span class="dk">gameState:</span> { turnNumber, activePlayer, currentPhase, stack }
<span class="dk">gameLog:</span> "T1: Player drew 7 cards\\nT1: Player kept hand\\n..."`}</code
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
										{#each gameLog.slice().reverse() as entry (entry.id)}
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
