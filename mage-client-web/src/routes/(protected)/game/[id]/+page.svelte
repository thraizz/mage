<script lang="ts">
	// From playtest/+page.svelte lines 1-110: Imports and setup
	import { goto } from '$app/navigation';
	import { createGamePageController } from '$lib/controllers/game-page-controller';
	import { auth } from '$lib/stores/auth';
	import { onMount } from 'svelte';
	// CHANGE: Import multiplayerGameStore instead of playtestGameStore
	import {
		multiplayerActiveControlSeat,
		multiplayerBattlefield,
		multiplayerExile,
		multiplayerGameStore,
		multiplayerIsInitialized,
		multiplayerLocalPlayer,
		multiplayerOpponents,
		multiplayerPlayers
	} from '$lib/stores/multiplayer-game';

	import { createGameUIState } from '$lib/stores/game-ui-state.svelte';
	import { toast } from '$lib/stores/toast';
	// Game components (from playtest/+page.svelte lines 20-46)
	import BattlefieldArea from '$lib/components/game/BattlefieldArea.svelte';
	import DragGhost from '$lib/components/game/DragGhost.svelte';
	import GameDialogs from '$lib/components/game/GameDialogs.svelte';
	import GameMenu from '$lib/components/game/GameMenu.svelte';
	import OpponentSection from '$lib/components/game/OpponentSection.svelte';
	import PlayerHand from '$lib/components/game/PlayerHand.svelte';
	import PlayerInfoRow from '$lib/components/game/PlayerInfoRow.svelte';
	import PlaytestHeader from '$lib/components/game/PlaytestHeader.svelte';

	import type { MenuAction } from '$lib/components/game/DeckContextMenu.svelte';

	import MulliganDialog from '$lib/components/game/MulliganDialog.svelte';
	import {
		currentDropZone,
		dragDropStore,
		draggedCardName,
		dragPosition,
		isDragging as isDraggingStore,
		isOverValidDropZone,
		type SourceZone
	} from '$lib/utils/drag-drop';
	import { getScryfallImageUrl } from '$lib/utils/scryfall';
	import { useDropZones } from '$lib/utils/use-drop-zones.svelte';

	// Page data from load function (CHANGE: Initialize from URL params)
	const { data } = $props<{ data: { gameId: string } }>();

	// State
	let loading = $state(true);
	let error = $state<string | null>(null);

	// Shared UI state
	const uiState = createGameUIState();

	// Drag-drop state (from playtest/+page.svelte lines 96-102)
	const isDragging = $derived($isDraggingStore);
	const dragCardName = $derived($draggedCardName);
	const dragPos = $derived($dragPosition);
	const isOverValidDrop = $derived($isOverValidDropZone);
	const dropZone = $derived($currentDropZone);

	// Game log (from playtest/+page.svelte line 104)
	const gameLog = $derived($multiplayerGameStore.log || []);

	// Drop zone elements
	let battlefieldDropZoneEl: HTMLDivElement | null = $state(null);
	let graveyardDropZoneEl: HTMLElement | null = $state(null);
	let exileDropZoneEl: HTMLElement | null = $state(null);
	let handDropZoneEl: HTMLElement | null = $state(null);
	let libraryDropZoneEl: HTMLElement | null = $state(null);
	let commandDropZoneEl: HTMLElement | null = $state(null);

	// Battlefield drag state (from playtest/+page.svelte lines 120-123)
	let battlefieldDragStartPosition = $state<{ x: number; y: number } | null>(null);
	let battlefieldIsDragPending = $state(false);
	const DRAG_THRESHOLD = 5;

	// Command zone drag state (from playtest/+page.svelte lines 125-127)
	let commandDragStartPosition = $state<{ x: number; y: number } | null>(null);
	let commandIsDragPending = $state(false);

	// Derived state from stores (from playtest/+page.svelte lines 133-140)
	// CHANGE: Use multiplayer stores instead of playtest stores
	const players = $derived($multiplayerPlayers);
	const me = $derived($multiplayerLocalPlayer);
	const otherPlayers = $derived($multiplayerOpponents);
	const battlefield = $derived($multiplayerBattlefield);
	const exile = $derived($multiplayerExile);
	const activeControlSeat = $derived($multiplayerActiveControlSeat);
	const isInitialized = $derived($multiplayerIsInitialized);

	// Selected opponent (from playtest/+page.svelte lines 142-150)
	const selectedOpponent = $derived.by(() => {
		if (otherPlayers.length === 0) return null;
		if (
			!uiState.selectedOpponentId ||
			!otherPlayers.find((p) => p.playerId === uiState.selectedOpponentId)
		) {
			return otherPlayers[0];
		}
		return otherPlayers.find((p) => p.playerId === uiState.selectedOpponentId) || otherPlayers[0];
	});

	// Split battlefield by controller (from playtest/+page.svelte lines 152-172)
	const myBattlefield = $derived(battlefield.filter((c) => c.controllerId === activeControlSeat));
	const opponentBattlefield = $derived.by(() => {
		const opponent = selectedOpponent;
		return opponent ? battlefield.filter((c) => c.controllerId === opponent.playerId) : [];
	});

	function isLandPermanent(cardType?: string | null): boolean {
		return !!cardType && /\bland\b/i.test(cardType);
	}

	const myBattlefieldNonlands = $derived(myBattlefield.filter((c) => !isLandPermanent(c.type)));
	const myBattlefieldLands = $derived(myBattlefield.filter((c) => isLandPermanent(c.type)));
	const opponentBattlefieldNonlands = $derived.by(() =>
		opponentBattlefield.filter((c) => !isLandPermanent(c.type))
	);
	const opponentBattlefieldLands = $derived.by(() =>
		opponentBattlefield.filter((c) => isLandPermanent(c.type))
	);

	// My cards (from playtest/+page.svelte lines 174-178)
	const myGrave = $derived(me?.graveyard || []);
	const myMana = $derived(
		me?.manaPool || { white: 0, blue: 0, black: 0, red: 0, green: 0, colorless: 0 }
	);

	// Reactive card lookup for counter dialog
	const selectedCardForCountersData = $derived.by(() => {
		const currentId = uiState.selectedCardForCounters?.id;
		if (!currentId) return null;

		const card =
			$multiplayerBattlefield.find((c) => c.id === currentId) ||
			me?.hand.find((c) => c.id === currentId) ||
			me?.graveyard.find((c) => c.id === currentId) ||
			null;

		return card;
	});

	// Hovered card
	const hoveredCard = $derived(
		uiState.hoveredCardId ? battlefield.find((c) => c.id === uiState.hoveredCardId) : null
	);

	const activePlayerName = $derived.by(() => {
		return players.find((p) => p.playerId === $multiplayerGameStore.activePlayerId)?.name ?? '';
	});

	const turnNumber = $derived.by(() => {
		const step = Math.max(1, $multiplayerGameStore.turn);
		const n = players.length;
		if (n <= 0) return step;
		return Math.floor((step - 1) / n) + 1;
	});

	// Command zone (from playtest/+page.svelte lines 229-244)
	const commandCards = $derived($multiplayerGameStore.command || []);
	const myCommandCards = $derived(
		commandCards.filter((c) => (c.ownerId || c.controllerId) === activeControlSeat)
	);

	const opponentCommandCards = $derived.by(() => {
		const opponent = selectedOpponent;
		return opponent
			? commandCards.filter((c) => (c.ownerId || c.controllerId) === opponent.playerId)
			: [];
	});

	const isCommanderGame = $derived(commandCards.length > 0);

	let selectedCardIds = $state<string[]>([]);
	let playingCardIds = $state<string[]>([]);

	// NOTE: We are not consindering priority anymore in the rules-light engine
	// const hasPriority = $derived($multiplayerGameStore.activePlayerId === me?.playerId);

	// Create game page controller
	const controller = createGamePageController(
		{
			gameStore: multiplayerGameStore,
			getState: () => $multiplayerGameStore,
			getLocalPlayer: () => me,
			getPlayers: () => players,
			getBattlefield: () => battlefield
		},
		{
			setScryDialog: (show, session) => {
				uiState.showScryDialog = show;
				uiState.currentScrySession = session;
			},
			setRevealTopDialog: (show, cards) => {
				uiState.showRevealTopDialog = show;
				uiState.revealedCards = cards;
			},
			setNumberInputDialog: (config) => {
				uiState.showNumberInputDialog = config.show;
				if (config.show && config.title && config.onConfirm) {
					uiState.numberInputDialogConfig = {
						title: config.title,
						defaultValue: config.defaultValue || 1,
						min: 1,
						max: 99,
						onConfirm: config.onConfirm
					};
				} else {
					uiState.numberInputDialogConfig = null;
				}
			},
			setDeckContextMenu: (show, position) => {
				uiState.showDeckContextMenu = show;
				if (position) {
					uiState.deckContextMenuPosition = position;
				}
			}
		}
	);

	/**
	 * Initialize from game ID (CHANGE: Server-based initialization)
	 */
	async function initializeFromGameId(): Promise<void> {
		loading = true;
		error = null;

		try {
			// Check authentication
			if (!$auth.user?.id) {
				error = 'Not authenticated';
				return;
			}

			console.log('[Multiplayer] Initializing with game ID:', data.gameId);

			// CHANGE: Initialize multiplayer store with game ID (only takes gameId)
			await multiplayerGameStore.initialize(data.gameId);

			loading = false;
		} catch (err) {
			console.error('[Multiplayer] Initialization failed:', err);
			error = err instanceof Error ? err.message : 'Failed to initialize game';
			loading = false;

			setTimeout(() => {
				goto('/lobby');
			}, 3000);
		}
	}

	// Use controller handlers (delegated to controller)
	const handleLifeChange = controller.handleLifeChange;
	const handlePoisonChange = controller.handlePoisonChange;
	const handleDrawCard = controller.handleDrawCard;
	const handleShuffleLibrary = controller.handleShuffleLibrary;
	const handleUntapAll = controller.handleUntapAll;
	const handleNextTurn = controller.handleNextTurn;

	/**
	 * Mulligan handlers
	 */
	const showMulliganDialog = $derived.by(() => {
		// Show if game is initialized, player exists, and hasn't kept hand yet
		return isInitialized && me && !me.keptHand;
	});

	async function handleKeepHand(): Promise<void> {
		if (!me) return;
		try {
			await multiplayerGameStore.keepHand(me.playerId);
			// Server will broadcast GAME_UPDATE with keptHand=true, which will hide the dialog
		} catch (err) {
			console.error('Failed to keep hand:', err);
			toast.error('Failed to keep hand');
		}
	}

	async function handleMulligan(): Promise<void> {
		if (!me) return;
		try {
			await multiplayerGameStore.mulligan(me.playerId);
			// Server will broadcast GAME_UPDATE with new hand, keptHand will remain false
		} catch (err) {
			console.error('Failed to mulligan:', err);
			toast.error('Failed to mulligan');
		}
	}

	/**
	 * Deck context menu handlers (from playtest/+page.svelte lines 497-626)
	 * CHANGE: Use multiplayerGameStore
	 */
	function handleDeckContextMenu(event: MouseEvent): void {
		if (!me) return;
		uiState.deckContextMenuPosition = { x: event.clientX, y: event.clientY };
		uiState.showDeckContextMenu = true;
	}

	function handleDrawN(count: number): void {
		if (!me) return;
		multiplayerGameStore.drawCards(me.playerId, count);
		toast.success(`Drew ${count} card(s)`);
	}

	function handleMill(count: number): void {
		if (!me) return;
		multiplayerGameStore.millCards(me.playerId, count);
		toast.success(`Milled ${count} card(s)`);
	}

	function handleScry(count: number): void {
		if (!me) return;
		const session = multiplayerGameStore.scryCards(me.playerId, count);
		if (session) {
			uiState.currentScrySession = session;
			uiState.showScryDialog = true;
		} else {
			toast.error('No cards to scry');
		}
	}

	function handleScryComplete(
		keepOnTop: import('$lib/generated/mage/v1/models').CardView[],
		putToBottom: import('$lib/generated/mage/v1/models').CardView[]
	): void {
		if (!me || !uiState.currentScrySession) return;
		const scryCount = uiState.currentScrySession.cards.length;
		multiplayerGameStore.applyScryDecision(me.playerId, scryCount, keepOnTop, putToBottom);
		uiState.showScryDialog = false;
		uiState.currentScrySession = null;
		toast.success(`Scry ${scryCount} complete`);
	}

	function handleRevealTop(count: number): void {
		if (!me) return;
		const cards = multiplayerGameStore.revealTopCards(me.playerId, count);
		uiState.revealedCards = cards;
		uiState.showRevealTopDialog = true;
	}

	function handleToggleRevealedTop(): void {
		if (!me) return;
		const willReveal = !me.revealedTopCard;
		multiplayerGameStore.setRevealedTop(me.playerId, willReveal);
		toast.info(willReveal ? 'Top card revealed permanently' : 'Top card hidden');
	}

	function showNumberInput(
		title: string,
		defaultValue: number,
		onConfirm: (value: number) => void
	): void {
		uiState.numberInputDialogConfig = {
			title,
			defaultValue,
			min: 1,
			max: 99,
			onConfirm: (value) => {
				onConfirm(value);
				uiState.showNumberInputDialog = false;
				uiState.numberInputDialogConfig = null;
			}
		};
		uiState.showNumberInputDialog = true;
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
							uiState.showDeckSearch = true;
						}
					},
					{
						label: 'Shuffle Library',
						onClick: handleShuffleLibrary
					}
				]
	);

	/**
	 * Handle battlefield card click (from playtest/+page.svelte lines 629-638)
	 * CHANGE: Use multiplayerGameStore
	 */
	function handleBattlefieldCardClick(cardId: string): void {
		const card = battlefield.find((c) => c.id === cardId);
		if (!card) return;
		multiplayerGameStore.tapCard(cardId, !card.tapped);
	}

	/**
	 * Handle battlefield card mouse down (from playtest/+page.svelte lines 640-689)
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
				dragDropStore.startDrag(
					cardId,
					cardName,
					'battlefield' as SourceZone,
					moveEvent.clientX,
					moveEvent.clientY
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
	 * Handle command zone card mouse down (from playtest/+page.svelte lines 691-736)
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
				dragDropStore.startDrag(
					cardId,
					cardName,
					'command' as SourceZone,
					moveEvent.clientX,
					moveEvent.clientY
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
	 * Handle battlefield drop (from playtest/+page.svelte lines 738-756)
	 * CHANGE: Use multiplayerGameStore
	 */
	function handleBattlefieldDrop(cardId: string): void {
		const dragState = $dragDropStore;
		const sourceZone = dragState.sourceZone;

		if (sourceZone === 'hand') {
			multiplayerGameStore.moveCardToZone(cardId, 'BATTLEFIELD');
		} else if (sourceZone && sourceZone !== 'battlefield') {
			multiplayerGameStore.moveCardToZone(cardId, 'BATTLEFIELD');
		}
	}

	/**
	 * Handle zone drop (from playtest/+page.svelte lines 758-765)
	 * CHANGE: Use multiplayerGameStore
	 */
	function handleZoneDrop(cardId: string, zone: string): void {
		multiplayerGameStore.moveCardToZone(cardId, zone);
	}

	/**
	 * Handle keyboard shortcuts (from playtest/+page.svelte lines 954-1057)
	 */
	function handleGlobalKeydown(event: KeyboardEvent): void {
		if (event.target instanceof HTMLInputElement || event.target instanceof HTMLTextAreaElement) {
			return;
		}

		const key = event.key.toLowerCase();

		switch (key) {
			case 'm':
				uiState.showMenu = !uiState.showMenu;
				event.preventDefault();
				break;
			case 'escape':
				if (uiState.showMenu) {
					uiState.showMenu = false;
					event.preventDefault();
				} else if (uiState.showKeyboardShortcuts) {
					uiState.showKeyboardShortcuts = false;
					event.preventDefault();
				} else if (uiState.showScryDialog) {
					uiState.showScryDialog = false;
					uiState.currentScrySession = null;
					event.preventDefault();
				} else if (uiState.showRevealTopDialog) {
					uiState.showRevealTopDialog = false;
					uiState.revealedCards = [];
					event.preventDefault();
				} else if (uiState.showNumberInputDialog) {
					uiState.showNumberInputDialog = false;
					uiState.numberInputDialogConfig = null;
					event.preventDefault();
				} else if (uiState.showDeckContextMenu) {
					uiState.showDeckContextMenu = false;
					event.preventDefault();
				} else if (uiState.showTokenCreator) {
					uiState.showTokenCreator = false;
					event.preventDefault();
				} else if (uiState.showCreateTokenDialog) {
					uiState.showCreateTokenDialog = false;
					event.preventDefault();
				} else if (uiState.showCounterDialog) {
					uiState.showCounterDialog = false;
					uiState.selectedCardForCounters = null;
					event.preventDefault();
				}
				break;
			case '?':
				uiState.showKeyboardShortcuts = !uiState.showKeyboardShortcuts;
				event.preventDefault();
				break;
			case 'f':
				uiState.showDeckSearch = true;
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
				uiState.showCreateTokenDialog = true;
				event.preventDefault();
				break;
		}

		// Hover card shortcuts - use controller
		if (hoveredCard) {
			const handled = controller.handleHoveredCardShortcut(key, hoveredCard.id);
			if (handled) {
				event.preventDefault();
			} else if (key === 'k') {
				// Counter dialog shortcut (not in controller)
				uiState.selectedCardForCounters = { id: hoveredCard.id, name: hoveredCard.name };
				uiState.showCounterDialog = true;
				event.preventDefault();
			}
		}
	}

	function handleSelectCard(cardId: string, isMultiSelect: boolean) {
		if (isMultiSelect) {
			// Toggle in selection
			if (selectedCardIds.includes(cardId)) {
				selectedCardIds = selectedCardIds.filter((id) => id !== cardId);
			} else {
				selectedCardIds = [...selectedCardIds, cardId];
			}
		} else {
			// Single select - toggle or replace
			selectedCardIds = selectedCardIds.includes(cardId) ? [] : [cardId];
		}
	}

	function handleClearSelection() {
		selectedCardIds = [];
	}

	/**
	 * Register all drop zones using shared helper
	 */
	useDropZones(
		() => ({
			battlefield: battlefieldDropZoneEl,
			graveyard: graveyardDropZoneEl,
			exile: exileDropZoneEl,
			hand: handDropZoneEl,
			library: libraryDropZoneEl,
			command: commandDropZoneEl
		}),
		{
			onBattlefieldDrop: handleBattlefieldDrop,
			onGraveyardDrop: (cardId) => handleZoneDrop(cardId, 'GRAVEYARD'),
			onExileDrop: (cardId) => handleZoneDrop(cardId, 'EXILE'),
			onHandDrop: (cardId) => handleZoneDrop(cardId, 'HAND'),
			onLibraryDrop: (cardId) => handleZoneDrop(cardId, 'LIBRARY'),
			onCommandDrop: (cardId) => handleZoneDrop(cardId, 'COMMAND')
		}
	);

	// Initialize on mount (CHANGE: Server-based initialization)
	onMount(() => {
		initializeFromGameId();
	});
</script>

<svelte:window onkeydown={handleGlobalKeydown} />

<svelte:head>
	<title>Multiplayer Game - MAGE</title>
</svelte:head>

<div class="game-container">
	{#if loading}
		<div class="loading-overlay">
			<div class="spinner"></div>
			<p>Loading game...</p>
		</div>
	{:else if error}
		<div class="error-overlay">
			<div class="error-icon">⚠️</div>
			<h2>Error</h2>
			<p>{error}</p>
			<button class="btn-primary" onclick={() => goto('/lobby')}> Return to Lobby </button>
		</div>
	{:else if !isInitialized}
		<div class="loading-overlay">
			<p>Initializing game state...</p>
		</div>
	{:else}
		<!-- Template from playtest/+page.svelte lines 1378+ adapted for multiplayer -->
		<PlaytestHeader
			isMultiplayer={true}
			{players}
			{activeControlSeat}
			availableSessions={0}
			{turnNumber}
			{activePlayerName}
			showAllHands={uiState.showAllHands}
			onBack={() => goto('/lobby')}
			onSessionsClick={() => {}}
			onSwitchPlayer={(playerId) => multiplayerGameStore.switchControlSeat(playerId)}
			onToggleAllHands={() => (uiState.showAllHands = !uiState.showAllHands)}
			onDrawCard={handleDrawCard}
			onUntapAll={handleUntapAll}
			onShuffleLibrary={handleShuffleLibrary}
			onSearchLibrary={() => (uiState.showDeckSearch = true)}
			onCreateToken={() => (uiState.showCreateTokenDialog = true)}
			onNextTurn={handleNextTurn}
			onShowKeyboardShortcuts={() => (uiState.showKeyboardShortcuts = true)}
			onShowDebug={() => (uiState.showDebugOverlay = true)}
			onToggleMenu={() => (uiState.showMenu = !uiState.showMenu)}
		/>

		<!-- Menu Overlay -->
		<GameMenu
			isOpen={uiState.showMenu}
			isMultiplayer={true}
			onClose={() => (uiState.showMenu = false)}
			onBackToLobby={() => goto('/lobby')}
			onShowKeyboardShortcuts={() => (uiState.showKeyboardShortcuts = true)}
		/>

		<!-- Main Game Area (from playtest/+page.svelte lines 1535-1677) -->
		<main class="game-layout">
			<!-- Opponent Section(s) -->
			{#if otherPlayers.length === 1}
				{#if selectedOpponent}
					{@const opponent = selectedOpponent}
					<OpponentSection
						{opponent}
						{otherPlayers}
						battlefieldNonlands={opponentBattlefieldNonlands}
						battlefieldLands={opponentBattlefieldLands}
						commandCards={opponentCommandCards}
						{isCommanderGame}
						showLifeMenu={uiState.showOpponentLifeMenu}
						onSelectOpponent={(playerId) => (uiState.selectedOpponentId = playerId)}
						onLifeChange={handleLifeChange}
						onPoisonChange={handlePoisonChange}
						onToggleLifeMenu={() => (uiState.showOpponentLifeMenu = !uiState.showOpponentLifeMenu)}
						onCardContextMenu={(cardId, cardName) => {
							uiState.selectedCardForCounters = { id: cardId, name: cardName };
							uiState.showCounterDialog = true;
						}}
					/>
				{/if}
			{:else}
				<!-- Multiplayer grid -->
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
								uiState.selectedCardForCounters = { id: cardId, name: cardName };
								uiState.showCounterDialog = true;
							}}
						/>
					{/each}
				</div>
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
							showLifeMenu={uiState.showOpponentLifeMenu}
							onSelectOpponent={(playerId) => (uiState.selectedOpponentId = playerId)}
							onLifeChange={handleLifeChange}
							onPoisonChange={handlePoisonChange}
							onToggleLifeMenu={() =>
								(uiState.showOpponentLifeMenu = !uiState.showOpponentLifeMenu)}
							onCardContextMenu={(cardId, cardName) => {
								uiState.selectedCardForCounters = { id: cardId, name: cardName };
								uiState.showCounterDialog = true;
							}}
						/>
					{/if}
				</div>
			{/if}

			<!-- My Battlefield Area -->
			<BattlefieldArea
				battlefieldNonlands={myBattlefieldNonlands}
				battlefieldLands={myBattlefieldLands}
				commandCards={myCommandCards}
				{isCommanderGame}
				{isDragging}
				{isOverValidDrop}
				{dropZone}
				hoveredCardId={uiState.hoveredCardId}
				onCardClick={handleBattlefieldCardClick}
				onCardMouseDown={handleBattlefieldCardMouseDown}
				onCardContextMenu={(cardId, cardName) => {
					uiState.selectedCardForCounters = { id: cardId, name: cardName };
					uiState.showCounterDialog = true;
				}}
				onCommandCardMouseDown={handleCommandCardMouseDown}
				onCardHover={(cardId) => (uiState.hoveredCardId = cardId)}
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
					showLifeMenu={uiState.showLifeMenu}
					onLifeChange={handleLifeChange}
					onPoisonChange={handlePoisonChange}
					onToggleLifeMenu={() => (uiState.showLifeMenu = !uiState.showLifeMenu)}
					onSearchLibrary={() => (uiState.showDeckSearch = true)}
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
					cards={me?.hand || []}
					{selectedCardIds}
					{playingCardIds}
					hasPriority={true}
					onSelectCard={handleSelectCard}
					onClearSelection={handleClearSelection}
					onCardClick={() => {}}
					size="normal"
					currentPhase="PRECOMBAT_MAIN"
					canDrag={true}
				/>
			</div>
		</main>

		<!-- Mulligan Dialog (not part of GameDialogs - shown before game starts) -->
		{#if showMulliganDialog && me}
			<MulliganDialog
				cards={me.hand}
				mulliganCount={me.mulliganCount}
				freeMulligans={$multiplayerGameStore.freeMulligans}
				onKeep={handleKeepHand}
				onMulligan={handleMulligan}
				hasKeptHand={me.keptHand}
				playerName={me.name}
			/>
		{/if}

		<!-- Game Dialogs Component -->
		<GameDialogs
			{uiState}
			gameId={data.gameId}
			{me}
			{selectedCardForCountersData}
			{deckContextMenuActions}
			onCreateToken={(name, types, power, toughness, color) => {
				multiplayerGameStore.createToken(name, types, power, toughness, color);
				uiState.showCreateTokenDialog = false;
			}}
			onAddCounter={(cardId, counterName, amount) =>
				multiplayerGameStore.addCounter(cardId, counterName, amount)}
			onRemoveCounter={(cardId, counterName, amount) =>
				multiplayerGameStore.removeCounter(cardId, counterName, amount)}
			onSetCounter={(cardId, counterName, amount) =>
				multiplayerGameStore.setCounter(cardId, counterName, amount)}
			onScryComplete={handleScryComplete}
			onNumberConfirm={(value) => {
				uiState.numberInputDialogConfig?.onConfirm(value);
				uiState.showNumberInputDialog = false;
			}}
			keyboardShortcutsMode="game"
			librarySearchVariant="multiplayer"
			showGameChatOverlay={true}
			onLibraryComplete={() => (uiState.showDeckSearch = false)}
			onLibraryCancel={() => (uiState.showDeckSearch = false)}
		/>

		<!-- Drag Ghost -->
		<DragGhost
			{isDragging}
			cardName={dragCardName}
			position={dragPos}
			{isOverValidDrop}
			imageSize="normal"
		/>
	{/if}
</div>

<style>
	/* Copy styles from playtest page */
	.game-container {
		width: 100vw;
		height: 100vh;
		overflow: hidden;
		background: linear-gradient(135deg, #1a1a2e 0%, #16213e 100%);
		position: relative;
	}

	.loading-overlay,
	.error-overlay {
		position: fixed;
		top: 0;
		left: 0;
		right: 0;
		bottom: 0;
		display: flex;
		flex-direction: column;
		align-items: center;
		justify-content: center;
		background: rgba(0, 0, 0, 0.9);
		z-index: 10000;
		color: white;
	}

	.spinner {
		width: 50px;
		height: 50px;
		border: 4px solid rgba(255, 255, 255, 0.1);
		border-top-color: #4a90e2;
		border-radius: 50%;
		animation: spin 1s linear infinite;
		margin-bottom: 20px;
	}

	@keyframes spin {
		to {
			transform: rotate(360deg);
		}
	}

	.error-icon {
		font-size: 4rem;
		margin-bottom: 20px;
	}

	.btn-primary {
		margin-top: 20px;
		padding: 10px 20px;
		background: #4a90e2;
		color: white;
		border: none;
		border-radius: 4px;
		cursor: pointer;
		font-size: 16px;
	}

	.btn-primary:hover {
		background: #357abd;
	}

	.all-hands-overlay {
		position: fixed;
		top: 80px;
		left: 0;
		right: 0;
		bottom: 0;
		background: rgba(0, 0, 0, 0.95);
		z-index: 900;
		padding: 20px;
		overflow-y: auto;
		display: flex;
		flex-direction: column;
		gap: 20px;
	}

	.player-hand-compact {
		background: rgba(255, 255, 255, 0.05);
		border: 2px solid rgba(255, 255, 255, 0.1);
		border-radius: 8px;
		padding: 15px;
	}

	.player-hand-compact.active {
		border-color: #4a90e2;
	}

	.compact-header {
		display: flex;
		justify-content: space-between;
		align-items: center;
		margin-bottom: 10px;
	}

	.player-name-compact {
		color: white;
		font-weight: bold;
		font-size: 18px;
	}

	.life-compact {
		color: #ff6b6b;
		display: flex;
		align-items: center;
		gap: 5px;
		font-size: 16px;
	}

	.cards-compact {
		display: flex;
		gap: 10px;
		flex-wrap: wrap;
	}

	.game-layout {
		height: 100vh;
		display: flex;
		flex-direction: column;
		padding-top: 60px;
	}

	.opponents-grid {
		padding: 10px;
	}

	.opponents-grid-large {
		display: grid;
		grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
		gap: 10px;
	}

	.opponents-grid-small {
		display: none;
	}

	@media (max-width: 1024px) {
		.opponents-grid-large {
			display: none;
		}

		.opponents-grid-small {
			display: block;
		}
	}

	.hand-area {
		position: relative;
		min-height: 150px;
		transition: background 0.2s;
	}

	.hand-area.drag-active {
		background: rgba(255, 255, 255, 0.05);
	}

	.hand-area.drag-valid {
		background: rgba(74, 144, 226, 0.2);
		border: 2px dashed #4a90e2;
	}
</style>
