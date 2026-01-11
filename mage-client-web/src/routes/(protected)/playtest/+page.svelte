<script lang="ts">
	import { onMount } from 'svelte';
	import { page } from '$app/stores';
	import { goto, beforeNavigate } from '$app/navigation';
	import { gameStore } from '$lib/stores/game';
	import {
		playtestGameStore,
		playtestPlayers,
		playtestLocalPlayer,
		playtestOpponents,
		playtestBattlefield,
		playtestExile,
		playtestActiveControlSeat,
		playtestIsInitialized
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
	import MulliganDialog from '$lib/components/game/MulliganDialog.svelte';
	import Keyboard from '@lucide/svelte/icons/keyboard';
	import Clock from '@lucide/svelte/icons/clock';
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
	let showKeyboardShortcuts = $state(false);
	let showAllHands = $state(false);
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

	// Drop zone elements
	let battlefieldDropZoneEl: HTMLDivElement | null = $state(null);
	let graveyardDropZoneEl: HTMLElement | null = $state(null);
	let exileDropZoneEl: HTMLElement | null = $state(null);
	let handDropZoneEl: HTMLElement | null = $state(null);
	let libraryDropZoneEl: HTMLElement | null = $state(null);
	let dropZoneUnregister: (() => void) | null = null;
	let graveyardDropZoneUnregister: (() => void) | null = null;
	let exileDropZoneUnregister: (() => void) | null = null;
	let handDropZoneUnregister: (() => void) | null = null;
	let libraryDropZoneUnregister: (() => void) | null = null;

	// Battlefield drag state
	let battlefieldDragStartPosition = $state<{ x: number; y: number } | null>(null);
	let battlefieldIsDragPending = $state(false);
	const DRAG_THRESHOLD = 5;
	
	// Command zone drag state
	let commandDragStartPosition = $state<{ x: number; y: number } | null>(null);
	let commandIsDragPending = $state(false);

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
		if (!selectedOpponentId || !otherPlayers.find(p => p.playerId === selectedOpponentId)) {
			// Auto-select first opponent
			return otherPlayers[0];
		}
		return otherPlayers.find(p => p.playerId === selectedOpponentId) || otherPlayers[0];
	});

	// Split battlefield by controller
	const myBattlefield = $derived(battlefield.filter(c => c.controllerId === activeControlSeat));
	const opponentBattlefield = $derived(() => {
		const opponent = selectedOpponent();
		return opponent ? battlefield.filter(c => c.controllerId === opponent.playerId) : [];
	});

	function isLandPermanent(cardType?: string | null): boolean {
		// Mage type strings are typically like: "Land", "Legendary Land", "Artifact Land", etc.
		return !!cardType && /\bland\b/i.test(cardType);
	}

	// Split battlefield rows: nonlands (top) + lands (bottom)
	const myBattlefieldNonlands = $derived(myBattlefield.filter(c => !isLandPermanent(c.type)));
	const myBattlefieldLands = $derived(myBattlefield.filter(c => isLandPermanent(c.type)));
	const opponentBattlefieldNonlands = $derived(() => opponentBattlefield().filter(c => !isLandPermanent(c.type)));
	const opponentBattlefieldLands = $derived(() => opponentBattlefield().filter(c => isLandPermanent(c.type)));

	// My cards (from controlling player perspective)
	const myGrave = $derived(me?.graveyard || []);
	const myMana = $derived(me?.manaPool || { white: 0, blue: 0, black: 0, red: 0, green: 0, colorless: 0 });

	// Hovered card
	const hoveredCard = $derived(hoveredCardId ? battlefield.find(c => c.id === hoveredCardId) : null);

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
	const allPlayersKept = $derived(players.every(p => p.keptHand));

	// Command zone (Commander): show for currently controlled player
	const commandCards = $derived($playtestGameStore.command || []);
	const myCommandCards = $derived(commandCards.filter(c => (c.ownerId || c.controllerId) === activeControlSeat));

	/**
	 * Initialize playtest from URL params
	 */
	async function initializeFromUrl(): Promise<void> {
		loading = true;
		error = null;

		try {
			const searchParams = $page.url.searchParams;
			const deckIds = validateDeckIds(searchParams);

			console.log('[Playtest] Initializing with deck IDs:', deckIds);

			const init = await initializePlaytest(deckIds);
			const initializedPlayers = init.players;
			const gameId = `playtest-${Date.now()}`;
			
			playtestGameStore.initialize(gameId, initializedPlayers);
			playtestGameStore.setCommand(init.command);

			// Initialize the normal game store with playtest data so PlayerHand works
			gameStore.initGame(gameId, initializedPlayers[0].playerId);
			syncPlaytestToGameStore();

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
		const state = $playtestGameStore;
		const controllingPlayer = players.find(p => p.playerId === activeControlSeat);
		
		if (!controllingPlayer) return;

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
			players: players.map(p => ({
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
		mulliganCount++;
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
		mulliganCount = 0;
		
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
		const player = players.find(p => p.playerId === playerId);
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
		const player = players.find(p => p.playerId === targetPlayerId);
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
		const newActivePlayer = players.find(p => p.playerId === $playtestGameStore.activePlayerId);
		if (newActivePlayer) {
			toast.info(`${newActivePlayer.name}'s turn`);
		}
	}

	/**
	 * Handle battlefield card click
	 */
	function handleBattlefieldCardClick(cardId: string): void {
		const card = battlefield.find(c => c.id === cardId);
		if (!card) return;

		// Toggle tap/untap
		playtestGameStore.tapCard(cardId, !card.tapped);
	}

	/**
	 * Handle battlefield card mouse down (for drag)
	 */
	function handleBattlefieldCardMouseDown(cardId: string, cardName: string, event: MouseEvent): void {
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
	 * Handle keyboard shortcuts
	 */
	function handleGlobalKeydown(event: KeyboardEvent): void {
		if (event.target instanceof HTMLInputElement || event.target instanceof HTMLTextAreaElement) {
			return;
		}

		const key = event.key.toLowerCase();

		switch (key) {
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
				showTokenCreator = true;
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

	// Initialize on mount
	onMount(() => {
		// Prefer restoring from persisted playtest state (refresh-safe).
		// If no persisted state exists, fall back to URL-based initialization.
		if ($playtestGameStore.isInitialized) {
			loading = false;

			// Restore mulligan phase based on first player who hasn't kept.
			const idx = players.findIndex((p) => !p.keptHand);
			mulliganPlayerIndex = idx === -1 ? null : idx;
			mulliganCount = 0;

			// Ensure the normal game store is initialized for shared components.
			gameStore.initGame($playtestGameStore.gameId, $playtestGameStore.activeControlSeat);
			syncPlaytestToGameStore();
			return;
		}

		initializeFromUrl();
	});

	// Cleanup ONLY when navigating away (client-side). This preserves state across refresh.
	beforeNavigate(({ from, to }) => {
		if (!from) return;
		if (from.url.pathname !== '/playtest') return;
		if (!to) return;
		if (to.url.pathname === '/playtest') return;

		playtestGameStore.reset();
		gameStore.reset();
	});
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
			<button class="btn-primary" onclick={() => goto('/lobby')}>
				Return to Lobby
			</button>
		</div>
	{:else if !isInitialized}
		<div class="loading-overlay">
			<p>Initializing game state...</p>
		</div>
	{:else if mulliganPlayerIndex !== null && !allPlayersKept}
		<MulliganDialog
			cards={players[mulliganPlayerIndex]?.hand || []}
			mulliganCount={mulliganCount}
			playerName={players[mulliganPlayerIndex]?.name}
			onKeep={handleKeepHand}
			onMulligan={handleMulligan}
			isLoading={false}
			hasKeptHand={false}
		/>
	{:else}
		<!-- Playtest Header -->
		<div class="playtest-header">
			<div class="header-left">
				<button class="btn-back" onclick={() => goto('/lobby')}>
					← Back to Lobby
				</button>
				<span class="mode-badge">Playtest Mode</span>
			</div>
			
			<div class="playtest-controls">
				<label for="playtest-controlling-select">Controlling:</label>
				<select id="playtest-controlling-select" class="player-select" value={activeControlSeat} onchange={(e) => switchPlayer(e.currentTarget.value)}>
					{#each players as player}
						<option value={player.playerId}>{player.name}</option>
					{/each}
				</select>
				<button class="btn-toggle" onclick={() => showAllHands = !showAllHands}>
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
					onclick={() => showKeyboardShortcuts = true}
					title="Keyboard shortcuts (?)"
					aria-label="Keyboard shortcuts"
				>
					<Keyboard size={20} aria-hidden="true" />
				</button>
				<button class="btn-debug" onclick={() => showDebugOverlay = true} title="Debug View">
					🔧
				</button>
			</div>
		</div>

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
					<!-- Opponent Header -->
					<div class="opponent-header-bar">
						<div class="opponent-identity">
							{#if otherPlayers.length > 1}
								<select 
									class="opponent-select" 
									value={opponent.playerId} 
									onchange={(e) => selectedOpponentId = e.currentTarget.value}
								>
									{#each otherPlayers as opp}
										<option value={opp.playerId}>{opp.name}</option>
									{/each}
								</select>
							{:else}
								<span class="opponent-name-label">{opponent.name}</span>
							{/if}
						</div>
						<div class="opponent-controls">
							<div class="life-group">
								<button class="stat-btn minus" onclick={() => handleLifeChange(-1, opponent.playerId)}>−</button>
								<button class="stat-display life" onclick={() => showOpponentLifeMenu = !showOpponentLifeMenu}>
									<span class="stat-icon">❤️</span>
									<span class="stat-value">{opponent.life}</span>
								</button>
								<button class="stat-btn plus" onclick={() => handleLifeChange(1, opponent.playerId)}>+</button>
							</div>

							{#if opponent.poison > 0}
								<div class="stat-display poison">
									<span class="stat-icon">☠️</span>
									<span class="stat-value">{opponent.poison}</span>
								</div>
							{/if}

							{#if showOpponentLifeMenu}
								<div bind:this={opponentLifeMenuEl} class="quick-menu opponent-menu">
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
											<button onclick={() => handlePoisonChange(-1, opponent.playerId)}>−1</button>
											<span class="menu-value">{opponent.poison}</span>
											<button onclick={() => handlePoisonChange(1, opponent.playerId)}>+1</button>
										</div>
									</div>
									<button class="menu-close" onclick={() => showOpponentLifeMenu = false}>✕</button>
								</div>
							{/if}

							<div class="opponent-stats">
								<span class="opponent-stat">Hand: {opponent.handCount}</span>
								<span class="opponent-stat">Library: {opponent.libraryCount}</span>
								<span class="opponent-stat">Graveyard: {opponent.graveyard.length}</span>
							</div>
						</div>
					</div>

					<!-- Opponent Battlefield (Non-editable) -->
					<div class="battlefield-area opponent-battlefield">
						<span class="zone-label">{opponent.name}'s Battlefield</span>
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
												imageUrl=""
												isTapped={card.tapped}
												isSelected={false}
												size="normal"
												onclick={() => {}}
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
												imageUrl=""
												isTapped={card.tapped}
												isSelected={false}
												size="normal"
												onclick={() => {}}
											/>
										</div>
									{/each}
								</div>
							{/if}

							{#if opponentBattlefieldNonlands().length === 0 && opponentBattlefieldLands().length === 0}
								<div class="empty-battlefield">
									No permanents
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
				<!-- Command Zone (Commander) -->
				{#if myCommandCards.length > 0}
					<div class="command-zone">
						<span class="zone-label">Command Zone</span>
						<div class="command-cards">
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
									onmouseenter={() => hoveredCardId = card.id}
									onmouseleave={() => { if (hoveredCardId === card.id) hoveredCardId = null; }}
								>
									<Card
										cardId={card.id}
										cardName={card.name}
										manaCost={card.manaCost}
										cardType={card.type}
										power={card.power}
										toughness={card.toughness}
										imageUrl=""
										isTapped={card.tapped}
										isSelected={false}
										size="normal"
										onclick={() => handleBattlefieldCardClick(card.id)}
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
									onmouseenter={() => hoveredCardId = card.id}
									onmouseleave={() => { if (hoveredCardId === card.id) hoveredCardId = null; }}
								>
									<Card
										cardId={card.id}
										cardName={card.name}
										manaCost={card.manaCost}
										cardType={card.type}
										power={card.power}
										toughness={card.toughness}
										imageUrl=""
										isTapped={card.tapped}
										isSelected={false}
										size="normal"
										onclick={() => handleBattlefieldCardClick(card.id)}
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

			<!-- Player Info Row -->
			{#if me}
				<div class="player-info-row">
					<div class="player-identity">
						<span class="player-name">{me.name}</span>
					</div>
					
					<div class="player-stats-inline">
						<div class="life-group">
							<button class="stat-btn minus" onclick={() => handleLifeChange(-1)}>−</button>
							<button class="stat-display life" onclick={() => showLifeMenu = !showLifeMenu}>
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
								onSearch={() => { showDeckSearch = true; }}
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
								<button class="menu-close" onclick={() => showLifeMenu = false}>✕</button>
							</div>
						{/if}
					</div>

					<div class="player-zones">
						<div bind:this={graveyardDropZoneEl} class="graveyard-drop-zone">
							<Graveyard
								cards={myGrave.map(c => ({
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
								cards={exile.map(c => ({
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
			<TokenCreator
				gameId="playtest"
				onClose={() => showTokenCreator = false}
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
				onClose={() => showDeckSearch = false}
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
							<button class="debug-close" onclick={() => showDebugOverlay = false}>✕</button>
						</div>
					</header>

					<main class="debug-content">
						<!-- Game State Overview -->
						<section class="debug-section">
							<div class="debug-section-header">
								<span>Game State Overview</span>
							</div>
							<div class="debug-code">
								<pre><code>{@html `<span class="dk">activeControlSeat:</span> <span class="ds">"${activeControlSeat}"</span>
<span class="dk">turn:</span> <span class="dn">${$playtestGameStore.turn}</span>
<span class="dk">activePlayerId:</span> <span class="ds">"${$playtestGameStore.activePlayerId}"</span>
<span class="dk">isInitialized:</span> <span class="db">${isInitialized}</span>`}</code></pre>
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
										<pre><code>{@html `<span class="dk">playerId:</span> <span class="ds">"${player.playerId}"</span>
<span class="dk">life:</span> <span class="dn">${player.life}</span>
<span class="dk">poison:</span> <span class="dn">${player.poison}</span>
<span class="dk">libraryCount:</span> <span class="dn">${player.libraryCount}</span>
<span class="dk">handCount:</span> <span class="dn">${player.handCount}</span>
<span class="dk">hand:</span> [${player.hand.map(c => `\n  <span class="ds">"${c.name}"</span> <span class="dc">// ${c.id}</span>`).join(',') || ''}
]
<span class="dk">graveyard:</span> [${player.graveyard.map(c => `\n  <span class="ds">"${c.name}"</span>`).join(',') || ''}
]
<span class="dk">library (first 5):</span> [${player.library.slice(0, 5).map(c => `\n  <span class="ds">"${c.name}"</span>`).join(',') || ''}
]
<span class="dk">keptHand:</span> <span class="db">${player.keptHand}</span>`}</code></pre>
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
										<pre><code>{battlefield.length > 0 
											? JSON.stringify(battlefield.map(c => ({
												id: c.id, name: c.name, controller: c.controllerId, tapped: c.tapped
											})), null, 2)
											: '[]'}</code></pre>
									</div>
								</div>
								<div class="debug-zone">
									<h4>🚫 Exile ({exile.length})</h4>
									<div class="debug-code small">
										<pre><code>{exile.length > 0 
											? JSON.stringify(exile.map(c => ({ id: c.id, name: c.name })), null, 2)
											: '[]'}</code></pre>
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
		to { transform: rotate(360deg); }
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
		margin-top: 80px;
		padding: 1rem;
		background: rgba(26, 31, 46, 0.95);
		border-bottom: 1px solid #2a3441;
		gap: 1rem;
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
		gap: 0.5rem;
		flex-shrink: 0;
	}

	.opponent-header-bar {
		display: flex;
		justify-content: space-between;
		align-items: center;
		padding: 0.5rem 0.75rem;
		background: rgba(26, 31, 46, 0.8);
		border: 1px solid #2a3441;
		border-radius: 8px;
		gap: 1rem;
	}

	.opponent-identity {
		display: flex;
		align-items: center;
		gap: 0.5rem;
		flex-shrink: 0;
	}

	.opponent-controls {
		display: flex;
		align-items: center;
		gap: 1rem;
		flex: 1;
		position: relative;
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

	.opponent-stats {
		display: flex;
		gap: 0.75rem;
		font-size: 0.875rem;
		align-items: center;
		margin-left: auto;
	}

	.opponent-stat {
		color: #94a3b8;
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
		flex: 1;
		min-height: 200px;
	}

	.opponent-battlefield {
		min-height: 150px;
		max-height: 250px;
		background: linear-gradient(135deg, #1a1217, #1c1428);
		border-color: rgba(200, 100, 100, 0.3);
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

	/* Command zone (Commander) */
	.command-zone {
		margin-bottom: 0.75rem;
	}

	.command-cards {
		display: flex;
		flex-wrap: wrap;
		gap: 0.5rem;
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
		box-shadow: 0 15px 40px rgba(0, 0, 0, 0.6), 0 0 30px rgba(34, 197, 94, 0.5);
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
		0%, 100% { opacity: 0.7; }
		50% { opacity: 1; }
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
	:global(.debug-code .dk) { color: #9cdcfe; }
	:global(.debug-code .ds) { color: #ce9178; }
	:global(.debug-code .dn) { color: #b5cea8; }
	:global(.debug-code .db) { color: #569cd6; }
	:global(.debug-code .dc) { color: #6a9955; font-style: italic; }
</style>
