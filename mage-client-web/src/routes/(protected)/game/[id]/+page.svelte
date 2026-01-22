<script lang="ts">
	// From playtest/+page.svelte lines 1-110: Imports and setup
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { auth } from '$lib/stores/auth';

	// CHANGE: Import multiplayerGameStore instead of playtestGameStore
	import {
		multiplayerGameStore,
		multiplayerPlayers,
		multiplayerLocalPlayer,
		multiplayerOpponents,
		multiplayerBattlefield,
		multiplayerExile,
		multiplayerActiveControlSeat,
		multiplayerIsInitialized,
		type PlaytestPlayer,
		type ScrySession
	} from '$lib/stores/multiplayer-game';

	import { toast } from '$lib/stores/toast';

	// Game components (from playtest/+page.svelte lines 20-46)
	import Card from '$lib/components/game/Card.svelte';
	import PlayerHand from '$lib/components/game/PlayerHand.svelte';
	import TokenCreator from '$lib/components/game/TokenCreator.svelte';
	import CreateTokenDialog from '$lib/components/game/CreateTokenDialog.svelte';
	import CounterDialog from '$lib/components/game/CounterDialog.svelte';
	import DeckContextMenu from '$lib/components/game/DeckContextMenu.svelte';
	import NumberInputDialog from '$lib/components/game/NumberInputDialog.svelte';
	import ScryDialog from '$lib/components/game/ScryDialog.svelte';
	import RevealTopDialog from '$lib/components/game/RevealTopDialog.svelte';
	import PlaytestHeader from '$lib/components/game/PlaytestHeader.svelte';
	import PlayerInfoRow from '$lib/components/game/PlayerInfoRow.svelte';
	import OpponentSection from '$lib/components/game/OpponentSection.svelte';
	import BattlefieldArea from '$lib/components/game/BattlefieldArea.svelte';
	import KeyboardShortcutsModal from '$lib/components/game/KeyboardShortcutsModal.svelte';

	// ADD: Multiplayer components
	import GameChatOverlay from '$lib/components/game/GameChatOverlay.svelte';

	import type { MenuAction } from '$lib/components/game/DeckContextMenu.svelte';
	import Keyboard from '@lucide/svelte/icons/keyboard';
	import X from '@lucide/svelte/icons/x';
	import Eye from '@lucide/svelte/icons/eye';
	import EyeOff from '@lucide/svelte/icons/eye-off';
	import Heart from '@lucide/svelte/icons/heart';

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

	// Page data from load function (CHANGE: Initialize from URL params)
	const { data } = $props<{ data: { gameId: string } }>();

	// State (from playtest/+page.svelte lines 58-92)
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

	// Deck context menu and dialog state (from playtest/+page.svelte lines 77-92)
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

	// Drag-drop state (from playtest/+page.svelte lines 96-102)
	const isDragging = $derived($isDraggingStore);
	const dragCardName = $derived($draggedCardName);
	const dragPos = $derived($dragPosition);
	const isOverValidDrop = $derived($isOverValidDropZone);
	const dropZone = $derived($currentDropZone);

	// Game log (from playtest/+page.svelte line 104)
	const gameLog = $derived($multiplayerGameStore.log || []);

	// Drop zone elements (from playtest/+page.svelte lines 106-118)
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
		if (!selectedOpponentId || !otherPlayers.find((p) => p.playerId === selectedOpponentId)) {
			return otherPlayers[0];
		}
		return otherPlayers.find((p) => p.playerId === selectedOpponentId) || otherPlayers[0];
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

	// Reactive card lookup for counter dialog (from playtest/+page.svelte lines 180-203)
	const selectedCardForCountersData = $derived.by(() => {
		const currentId = selectedCardForCounters?.id;
		if (!currentId) return null;

		const card =
			$multiplayerBattlefield.find((c) => c.id === currentId) ||
			me?.hand.find((c) => c.id === currentId) ||
			me?.graveyard.find((c) => c.id === currentId) ||
			null;

		return card;
	});

	// Hovered card (from playtest/+page.svelte lines 205-208)
	const hoveredCard = $derived(
		hoveredCardId ? battlefield.find((c) => c.id === hoveredCardId) : null
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

	/**
	 * Handle life change (from playtest/+page.svelte lines 439-443)
	 * CHANGE: Use multiplayerGameStore
	 */
	function handleLifeChange(delta: number, playerId?: string): void {
		const targetPlayerId = playerId || me?.playerId;
		if (!targetPlayerId) return;
		multiplayerGameStore.modifyLife(targetPlayerId, delta);
	}

	/**
	 * Handle poison counter change (from playtest/+page.svelte lines 446-456)
	 * CHANGE: Use multiplayerGameStore
	 */
	function handlePoisonChange(delta: number, playerId?: string): void {
		const targetPlayerId = playerId || me?.playerId;
		if (!targetPlayerId) return;
		const player = players.find((p) => p.playerId === targetPlayerId);
		if (!player) return;
		const newValue = Math.max(0, (player.poison || 0) + delta);
		multiplayerGameStore.setPlayerCounter(targetPlayerId, 'poison', newValue);
	}

	/**
	 * Draw a card (from playtest/+page.svelte lines 459-465)
	 * CHANGE: Use multiplayerGameStore
	 */
	function handleDrawCard(): void {
		if (!me) return;
		multiplayerGameStore.drawCards(me.playerId, 1);
		toast.success('Drew a card');
	}

	/**
	 * Shuffle library (from playtest/+page.svelte lines 467-474)
	 * CHANGE: Use multiplayerGameStore
	 */
	function handleShuffleLibrary(): void {
		if (!me) return;
		multiplayerGameStore.shuffleLibrary(me.playerId);
		toast.success('Shuffled library');
	}

	/**
	 * Untap all permanents (from playtest/+page.svelte lines 476-483)
	 * CHANGE: Use multiplayerGameStore
	 */
	function handleUntapAll(): void {
		if (!me) return;
		multiplayerGameStore.untapAll(me.playerId);
		toast.success('Untapped all');
	}

	/**
	 * Next turn (from playtest/+page.svelte lines 485-494)
	 * CHANGE: Use multiplayerGameStore
	 */
	function handleNextTurn(): void {
		multiplayerGameStore.nextTurn();
		const newActivePlayer = players.find(
			(p) => p.playerId === $multiplayerGameStore.activePlayerId
		);
		if (newActivePlayer) {
			toast.info(`${newActivePlayer.name}'s turn`);
		}
	}

	/**
	 * Deck context menu handlers (from playtest/+page.svelte lines 497-626)
	 * CHANGE: Use multiplayerGameStore
	 */
	function handleDeckContextMenu(event: MouseEvent): void {
		if (!me) return;
		deckContextMenuPosition = { x: event.clientX, y: event.clientY };
		showDeckContextMenu = true;
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
		multiplayerGameStore.applyScryDecision(me.playerId, scryCount, keepOnTop, putToBottom);
		showScryDialog = false;
		currentScrySession = null;
		toast.success(`Scry ${scryCount} complete`);
	}

	function handleRevealTop(count: number): void {
		if (!me) return;
		const cards = multiplayerGameStore.revealTopCards(me.playerId, count);
		revealedCards = cards;
		showRevealTopDialog = true;
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
				showMenu = !showMenu;
				event.preventDefault();
				break;
			case 'escape':
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
				}
				break;
			case '?':
				showKeyboardShortcuts = !showKeyboardShortcuts;
				event.preventDefault();
				break;
			case 'f':
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

		// Hover card shortcuts (from playtest/+page.svelte lines 1032-1056)
		if (hoveredCard) {
			switch (key) {
				case 'd':
					multiplayerGameStore.moveCardToZone(hoveredCard.id, 'GRAVEYARD');
					event.preventDefault();
					break;
				case 's':
					multiplayerGameStore.moveCardToZone(hoveredCard.id, 'EXILE');
					event.preventDefault();
					break;
				case 'r':
					multiplayerGameStore.moveCardToZone(hoveredCard.id, 'HAND');
					event.preventDefault();
					break;
				case 't':
					multiplayerGameStore.moveCardToZone(hoveredCard.id, 'LIBRARY');
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
	 * Register drop zones (from playtest/+page.svelte lines 1059-1168)
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
			{players}
			{activeControlSeat}
			availableSessions={0}
			{turnNumber}
			{activePlayerName}
			{showAllHands}
			onBack={() => goto('/lobby')}
			onSessionsClick={() => {}}
			onSwitchPlayer={(playerId) => multiplayerGameStore.switchControlSeat(playerId)}
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

		<!-- Menu Overlay (from playtest/+page.svelte lines 1403-1514) -->
		{#if showMenu}
			<div
				class="menu-backdrop"
				role="button"
				tabindex="0"
				onclick={() => (showMenu = false)}
				onkeydown={(e) => e.key === 'Escape' && (showMenu = false)}
			></div>

			<div class="menu-overlay open">
				<div class="menu-header">
					<h2>Menu</h2>
					<button class="menu-close-btn" onclick={() => (showMenu = false)} aria-label="Close menu">
						<X size={24} />
					</button>
				</div>

				<div class="menu-content">
					<div class="menu-section">
						<h3 class="menu-section-title">Controls</h3>
						<div class="menu-section-content">
							<label>
								<span class="menu-label">Controlling:</span>
								<select
									class="control-select"
									value={activeControlSeat}
									onchange={(e) => multiplayerGameStore.switchControlSeat(e.currentTarget.value)}
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
						</div>
					</div>

					<div class="menu-section">
						<h3 class="menu-section-title">Navigation</h3>
						<div class="menu-section-content">
							<button class="menu-btn" onclick={() => goto('/lobby')}> ← Back to Lobby </button>
						</div>
					</div>
				</div>
			</div>
		{/if}

		<!-- All Hands Overlay (from playtest/+page.svelte lines 1516-1533) -->
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
								selectedCardForCounters = { id: cardId, name: cardName };
								showCounterDialog = true;
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

			<!-- My Battlefield Area -->
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
			<TokenCreator gameId={data.gameId} onClose={() => (showTokenCreator = false)} />
		{/if}

		<!-- Create Token Dialog -->
		{#if showCreateTokenDialog}
			<CreateTokenDialog
				onCreateToken={(name, types, power, toughness, color) => {
					multiplayerGameStore.createToken(name, types, power, toughness, color);
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
					multiplayerGameStore.addCounter(card.id, counterName, amount);
				}}
				onRemoveCounter={(counterName, amount) => {
					const card = selectedCardForCountersData;
					multiplayerGameStore.removeCounter(card.id, counterName, amount);
				}}
				onSetCounter={(counterName, amount) => {
					const card = selectedCardForCountersData;
					multiplayerGameStore.setCounter(card.id, counterName, amount);
				}}
				onClose={() => {
					showCounterDialog = false;
					selectedCardForCounters = null;
				}}
			/>
		{/if}

		<!-- Deck Search -->
		<!-- TODO: Re-enable when LibrarySearch component supports local data or when server-side library search is implemented
		{#if showDeckSearch && me}
			<LibrarySearch
				gameId={data.gameId}
				librarySearchData={{
					playerId: me.playerId,
					message: 'Search your library',
					destination: 'hand',
					cards: me.library,
					canCancel: true
				}}
				onComplete={() => (showDeckSearch = false)}
				onCancel={() => (showDeckSearch = false)}
			/>
		{/if}
		-->

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

		<KeyboardShortcutsModal bind:open={showKeyboardShortcuts} mode="game" />

		<!-- ADD: GameChatOverlay for multiplayer -->
		<GameChatOverlay gameId={data.gameId} />

		<!-- Drag Ghost (from playtest/+page.svelte lines 1978-1990) -->
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

	.menu-backdrop {
		position: fixed;
		top: 0;
		left: 0;
		right: 0;
		bottom: 0;
		background: rgba(0, 0, 0, 0.5);
		z-index: 999;
	}

	.menu-overlay {
		position: fixed;
		top: 0;
		right: -400px;
		width: 400px;
		height: 100vh;
		background: #1a1a2e;
		box-shadow: -2px 0 10px rgba(0, 0, 0, 0.3);
		z-index: 1000;
		transition: right 0.3s ease;
		overflow-y: auto;
	}

	.menu-overlay.open {
		right: 0;
	}

	.menu-header {
		display: flex;
		justify-content: space-between;
		align-items: center;
		padding: 20px;
		border-bottom: 1px solid rgba(255, 255, 255, 0.1);
	}

	.menu-header h2 {
		margin: 0;
		color: white;
		font-size: 24px;
	}

	.menu-close-btn {
		background: none;
		border: none;
		color: white;
		cursor: pointer;
		padding: 5px;
	}

	.menu-content {
		padding: 20px;
	}

	.menu-section {
		margin-bottom: 30px;
	}

	.menu-section-title {
		color: #4a90e2;
		font-size: 14px;
		text-transform: uppercase;
		margin-bottom: 10px;
		letter-spacing: 1px;
	}

	.menu-section-content {
		display: flex;
		flex-direction: column;
		gap: 10px;
	}

	.menu-btn {
		padding: 10px 15px;
		background: rgba(255, 255, 255, 0.1);
		border: 1px solid rgba(255, 255, 255, 0.2);
		border-radius: 4px;
		color: white;
		cursor: pointer;
		display: flex;
		align-items: center;
		gap: 8px;
		transition: background 0.2s;
	}

	.menu-btn:hover {
		background: rgba(255, 255, 255, 0.2);
	}

	.menu-label {
		color: #aaa;
		font-size: 14px;
		margin-bottom: 5px;
		display: block;
	}

	.control-select {
		width: 100%;
		padding: 8px;
		background: rgba(0, 0, 0, 0.3);
		border: 1px solid rgba(255, 255, 255, 0.2);
		border-radius: 4px;
		color: white;
		font-size: 14px;
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

	.drag-ghost {
		position: fixed;
		pointer-events: none;
		z-index: 10000;
		transform: translate(-50%, -50%);
	}

	.drag-ghost-card {
		opacity: 0.8;
		transition: opacity 0.2s;
	}

	.drag-ghost-card.valid {
		opacity: 1;
	}

	.drag-ghost-image {
		width: 100px;
		height: 140px;
		object-fit: cover;
		border-radius: 8px;
		box-shadow: 0 4px 12px rgba(0, 0, 0, 0.5);
	}

	.drag-ghost-name {
		display: block;
		padding: 10px;
		background: rgba(0, 0, 0, 0.9);
		color: white;
		border-radius: 4px;
		font-size: 14px;
	}
</style>
