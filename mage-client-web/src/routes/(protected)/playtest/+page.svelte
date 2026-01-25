<script lang="ts">
	import GameStateLog from '../../../lib/components/game/GameStateLog.svelte';

	import { goto } from '$app/navigation';
	import { page } from '$app/state';
	import { createGamePageController } from '$lib/controllers/game-page-controller';
	import { initializePlaytest, validateDeckIds } from '$lib/playtest/initializer';
	import { createGameUIState } from '$lib/stores/game-ui-state.svelte';
	import {
		playtestActiveControlSeat,
		playtestBattlefield,
		playtestExile,
		playtestGameStore,
		playtestIsInitialized,
		playtestLocalPlayer,
		playtestOpponents,
		playtestPlayers,
		type PlaytestSessionMeta
	} from '$lib/stores/playtest-game';
	import { toast } from '$lib/stores/toast';
	import { onMount } from 'svelte';
	// Game components
	import BattlefieldArea from '$lib/components/game/BattlefieldArea.svelte';
	import Card from '$lib/components/game/Card.svelte';
	import type { MenuAction } from '$lib/components/game/DeckContextMenu.svelte';
	import DragGhost from '$lib/components/game/DragGhost.svelte';
	import GameDialogs from '$lib/components/game/GameDialogs.svelte';
	import GameMenu from '$lib/components/game/GameMenu.svelte';
	import MulliganDialog from '$lib/components/game/MulliganDialog.svelte';
	import OpponentSection from '$lib/components/game/OpponentSection.svelte';
	import PlayerHand from '$lib/components/game/PlayerHand.svelte';
	import PlayerInfoRow from '$lib/components/game/PlayerInfoRow.svelte';
	import PlaytestHeader from '$lib/components/game/PlaytestHeader.svelte';
	import {
		currentDropZone,
		dragDropStore,
		draggedCardName,
		dragPosition,
		isDragging as isDraggingStore,
		isOverValidDropZone,
		type SourceZone
	} from '$lib/utils/drag-drop';
	import { useDropZones } from '$lib/utils/use-drop-zones.svelte';
	import Copy from '@lucide/svelte/icons/copy';
	import Heart from '@lucide/svelte/icons/heart';

	// State
	let loading = $state(true);
	let error = $state<string | null>(null);

	// Shared UI state
	const uiState = createGameUIState();

	// Mulligan state (playtest-specific)
	let mulliganPlayerIndex = $state<number | null>(null);

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
		if (
			!uiState.selectedOpponentId ||
			!otherPlayers.find((p) => p.playerId === uiState.selectedOpponentId)
		) {
			// Auto-select first opponent
			return otherPlayers[0];
		}
		return otherPlayers.find((p) => p.playerId === uiState.selectedOpponentId) || otherPlayers[0];
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
		const currentId = uiState.selectedCardForCounters?.id;

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
			'[uiState.selectedCardForCountersData] Re-evaluated.',
			`Card: ${card?.name}`,
			`Counters: ${card?.counters}`
		);

		return card;
	});

	// Hovered card
	const hoveredCard = $derived(
		uiState.hoveredCardId ? battlefield.find((c) => c.id === uiState.hoveredCardId) : null
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

	// Create game page controller
	const controller = createGamePageController(
		{
			gameStore: playtestGameStore,
			getState: () => $playtestGameStore,
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
	 * Update the URL with playtestId parameter
	 */
	async function updateUrlWithPlaytestId(playtestId: string): Promise<void> {
		const newSearchParams = new URLSearchParams(page.url.searchParams);
		newSearchParams.set('playtestId', playtestId);
		await goto(`${page.url.pathname}?${newSearchParams.toString()}`, {
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
			const searchParams = new URLSearchParams(page.url.searchParams);
			const deckIds = validateDeckIds(searchParams);

			// Parse mulligan settings from URL
			const mulliganType = (searchParams.get('mulliganType') as 'london') || 'london';
			const freeMulligans = parseInt(searchParams.get('freeMulligans') ?? '0', 10);

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

	// Use controller handlers (delegated to controller)
	const handleLifeChange = controller.handleLifeChange;
	const handlePoisonChange = controller.handlePoisonChange;
	const handleDrawCard = controller.handleDrawCard;
	const handleShuffleLibrary = controller.handleShuffleLibrary;
	const handleUntapAll = controller.handleUntapAll;
	const handleNextTurn = controller.handleNextTurn;

	// Use controller deck handlers
	const handleDeckContextMenu = controller.handleDeckContextMenu;

	// Wrap handleScryComplete to pass currentSession
	function handleScryComplete(
		keepOnTop: import('$lib/generated/mage/v1/models').CardView[],
		putToBottom: import('$lib/generated/mage/v1/models').CardView[]
	): void {
		if (!uiState.currentScrySession) return;
		controller.handleScryComplete(keepOnTop, putToBottom, uiState.currentScrySession);
	}

	// Create deck context menu actions using controller
	const deckContextMenuActions = $derived<MenuAction[]>(
		controller.createDeckContextMenuActions(() => {
			uiState.showDeckSearch = true;
		})
	);

	// Use controller battlefield handlers
	const handleBattlefieldCardClick = controller.handleBattlefieldCardClick;

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
	 * Handle battlefield drop
	 */
	function handleBattlefieldDrop(cardId: string): void {
		const dragState = $dragDropStore;
		controller.handleBattlefieldDrop(cardId, dragState.sourceZone);
	}

	/**
	 * Handle zone drop (graveyard, exile, hand)
	 */
	const handleZoneDrop = controller.handleZoneDrop;

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
					uiState.showMenu = !uiState.showMenu;
					event.preventDefault();
				}
				break;
			case 'escape':
				// Escape - Close menu or modals
				if (uiState.showMenu) {
					uiState.showMenu = false;
					event.preventDefault();
				} else if (uiState.showKeyboardShortcuts) {
					uiState.showKeyboardShortcuts = false;
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
				} else if (uiState.showDebugOverlay) {
					uiState.showDebugOverlay = false;
					event.preventDefault();
				}
				break;
			case '?':
				uiState.showKeyboardShortcuts = !uiState.showKeyboardShortcuts;
				event.preventDefault();
				break;
			case 'f':
				// F - Search your deck
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
		const playtestId = page.url.searchParams.get('playtestId');
		if (playtestId) {
			const success = playtestGameStore.restoreSession(playtestId);
			if (success) {
				loading = false;

				// Restore mulligan phase based on first player who hasn't kept.
				const idx = players.findIndex((p) => !p.keptHand);
				mulliganPlayerIndex = idx === -1 ? null : idx;

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
		const hasUrlDecks = page.url.searchParams.has('d1') || page.url.searchParams.has('d2');
		const hasPlaytestIdInUrl = page.url.searchParams.has('playtestId');

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
			isMultiplayer={false}
			{players}
			{activeControlSeat}
			availableSessions={availableSessions.length}
			{turnNumber}
			{activePlayerName}
			showAllHands={uiState.showAllHands}
			onBack={() => goto('/lobby')}
			onSessionsClick={() => {
				loadAvailableSessions();
				showSessionPicker = true;
			}}
			onSwitchPlayer={switchPlayer}
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

		<!-- Session Picker Overlay for restoring older playtest sessions -->
		{#if showSessionPicker}
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
						<button class="btn-back" onclick={() => (showSessionPicker = false)}> Back</button>
						<button class="btn-primary" onclick={() => goto('/lobby')}> Start New Playtest </button>
					</div>
				</div>
			</div>
		{/if}

		<!-- Menu Overlay -->
		<GameMenu
			isOpen={uiState.showMenu}
			isMultiplayer={false}
			{players}
			{activeControlSeat}
			{turnNumber}
			{activePlayerName}
			availableSessions={availableSessions.length}
			showAllHands={uiState.showAllHands}
			onClose={() => (uiState.showMenu = false)}
			onBackToLobby={() => goto('/lobby')}
			onShowKeyboardShortcuts={() => (uiState.showKeyboardShortcuts = true)}
			onSwitchPlayer={(playerId) => playtestGameStore.switchControlSeat(playerId)}
			onToggleAllHands={() => (uiState.showAllHands = !uiState.showAllHands)}
			onNextTurn={handleNextTurn}
			onShowDebug={() => (uiState.showDebugOverlay = true)}
			onSessionsClick={() => (showSessionPicker = true)}
		/>

		<!-- All Hands Overlay -->
		{#if uiState.showAllHands}
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

		<div class="game-and-log-container">
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
				{:else}
					<!-- Multiplayer (3-4 players): Grid on large screens, cycling on small -->
					<!-- Grid layout (shown on large screens) -->
					<div class="opponents-grid opponents-grid-large">
						{#each otherPlayers as opponent (opponent.playerId)}
							{@const oppBattlefield = battlefield.filter(
								(c) => c.controllerId === opponent.playerId
							)}
							{@const oppBattlefieldNonlands = oppBattlefield.filter(
								(c) => !isLandPermanent(c.type)
							)}
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

				<!-- My Battlefield Area (Editable) -->
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
					<!-- TODO: Verify the props here -->
					<PlayerHand
						cards={me?.hand || []}
						selectedCardIds={[]}
						playingCardIds={[]}
						hasPriority={true}
						size="normal"
						currentPhase="PRECOMBAT_MAIN"
						canDrag={true}
					/>
				</div>
			</main>
			<GameStateLog />
		</div>

		<!-- Game Dialogs Component -->
		<GameDialogs
			{uiState}
			gameId="playtest"
			{me}
			{selectedCardForCountersData}
			{deckContextMenuActions}
			onCreateToken={(name, types, power, toughness, color) => {
				playtestGameStore.createToken(name, types, power, toughness, color);
				uiState.showCreateTokenDialog = false;
			}}
			onAddCounter={(cardId, counterName, amount) =>
				playtestGameStore.addCounter(cardId, counterName, amount)}
			onRemoveCounter={(cardId, counterName, amount) =>
				playtestGameStore.removeCounter(cardId, counterName, amount)}
			onSetCounter={(cardId, counterName, amount) =>
				playtestGameStore.setCounter(cardId, counterName, amount)}
			onScryComplete={handleScryComplete}
			onNumberConfirm={(value) => {
				uiState.numberInputDialogConfig?.onConfirm(value);
				uiState.showNumberInputDialog = false;
			}}
			keyboardShortcutsMode="playtest"
			librarySearchVariant="playtest"
			onLibraryMove={(cardId, zone) => playtestGameStore.moveCardToZone(cardId, zone)}
			onLibraryShuffle={() => me && playtestGameStore.shuffleLibrary(me.playerId)}
			onLibraryClose={() => (uiState.showDeckSearch = false)}
		/>

		<!-- Debug Overlay -->
		{#if uiState.showDebugOverlay}
			<div class="debug-overlay" role="dialog" aria-modal="true">
				<div class="debug-modal">
					<header class="debug-header">
						<div class="debug-header-left">
							<h2>🔧 Playtest Debug View</h2>
							<div class="debug-status connected">● Playtest Mode</div>
						</div>
						<div class="debug-header-right">
							<button class="debug-close" onclick={() => (uiState.showDebugOverlay = false)}
								>✕</button
							>
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
						</section>

						<!-- Game State Log -->
						<GameStateLog />

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
		<DragGhost
			{isDragging}
			cardName={dragCardName}
			position={dragPos}
			{isOverValidDrop}
			imageSize="small"
		/>
	{/if}
</div>
