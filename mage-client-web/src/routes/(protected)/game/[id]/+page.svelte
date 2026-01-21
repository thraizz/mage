<script lang="ts">
	import { onMount, onDestroy } from 'svelte';
	import { get } from 'svelte/store';
	import { goto } from '$app/navigation';
	import { auth } from '$lib/stores/auth';
	import { websocketStore } from '$lib/stores/websocket';
	import { getSessionIdFromToken } from '$lib/utils/jwt';
	import {
		gameStore,
		players,
		localPlayer,
		opponents,
		hasPriority,
		currentPhase,
		currentStep,
		currentTurn,
		battlefield,
		stack,
		command,
		exile,
		myHand,
		myGraveyard,
		myManaPool,
		pendingPrompt,
		gameOver,
		winner,
		gameError,
		isLoading,
		getCardById
	} from '$lib/stores/game';
	import {
		joinGame,
		fetchGameView,
		passPriority,
		passUntilNextTurn,
		concedeGame,
		sendPlayerBoolean,
		sendPlayerString,
		keepHand,
		mulligan,
		playLand,
		advancePhase,
		activateManaAbility
	} from '$lib/api/game';
	import { CardActionType, type CardView } from '$lib/generated/mage/v1/models';
	import {
		tapUntap,
		untapAll,
		flipCard,
		transformCard,
		moveCard as moveCardToZone,
		modifyCardCounter,
		drawCards,
		shuffleLibrary,
		nextTurn,
		searchLibrary,
		addToStack,
		removeFromStack,
		modifyLife,
		setPlayerCounter
	} from '$lib/api/direct-actions';
	import type { GameCard, GamePhase } from '$lib/types/game';
	import {
		dragDropStore,
		isDragging as isDraggingStore,
		draggedCardId,
		draggedCardName,
		dragPosition,
		isOverValidDropZone,
		currentDropZone,
		getAllValidDropZones,
		type SourceZone,
		type DropZone
	} from '$lib/utils/drag-drop';
	import { moveCard } from '$lib/api/direct-actions';
	import { toast } from '$lib/stores/toast';
	import { getScryfallImageUrl } from '$lib/utils/scryfall';

	// Game components
	import Card from '$lib/components/game/Card.svelte';
	import PlayerHand from '$lib/components/game/PlayerHand.svelte';
	import Graveyard from '$lib/components/game/Graveyard.svelte';
	import ExileZone from '$lib/components/game/ExileZone.svelte';
	import LibraryZone from '$lib/components/game/LibraryZone.svelte';
	import ManaPool from '$lib/components/game/ManaPool.svelte';
	import GameHeader from '$lib/components/game/GameHeader.svelte';
	import MulliganDialog from '$lib/components/game/MulliganDialog.svelte';
	import BattlefieldArea from '$lib/components/game/BattlefieldArea.svelte';

	import PriorityActionBar from '$lib/components/game/PriorityActionBar.svelte';
	import ActionLogOverlay from '$lib/components/game/ActionLogOverlay.svelte';
	import GameChatOverlay from '$lib/components/game/GameChatOverlay.svelte';
	import OpponentPanel from '$lib/components/game/OpponentPanel.svelte';
	import OpponentSection from '$lib/components/game/OpponentSection.svelte';
	import DebugOverlay from '$lib/components/game/DebugOverlay.svelte';
	import ManaPayment from '$lib/components/game/ManaPayment.svelte';
	import XManaSelector from '$lib/components/game/XManaSelector.svelte';
	import LibrarySearch from '$lib/components/game/LibrarySearch.svelte';
	import PlayerInfoRow from '$lib/components/game/PlayerInfoRow.svelte';
	import type {
		GamePlayManaData,
		GamePlayXManaData,
		GameAssignDamageData
	} from '$lib/generated/mage/v1/websocket';
	import Keyboard from '@lucide/svelte/icons/keyboard';
	import KeyboardShortcutsModal from '$lib/components/game/KeyboardShortcutsModal.svelte';

	// Combat components
	import DeclareAttackers from '$lib/components/game/DeclareAttackers.svelte';
	import DeclareBlockers from '$lib/components/game/DeclareBlockers.svelte';
	import AssignDamage from '$lib/components/game/AssignDamage.svelte';

	// Direct player control components
	import DeckContextMenu from '$lib/components/game/DeckContextMenu.svelte'; // Phase 4: Deck context menu (plan lines 854-858)
	import TokenCreator from '$lib/components/game/TokenCreator.svelte';
	import VisualStack from '$lib/components/game/VisualStack.svelte';
	import RollbackConsentDialog from '$lib/components/game/RollbackConsentDialog.svelte';
	import { requestRollback, respondToRollback } from '$lib/api/game';
	import {
		combatStore,
		isInCombat,
		canAttackCardIds,
		declaredAttackerIds,
		canBlockCardIds,
		assignedBlockerIds
	} from '$lib/stores/combat';
	import { visualStackStore, visualStackIsOpen, visualStackCount } from '$lib/stores/visual-stack';
	import {
		parseCombatOptions,
		type DeclaredAttacker,
		type DefenderTarget,
		type DamageAssignmentPrompt,
		type ParsedCombatOptions
	} from '$lib/types/combat';
	import CardContextMenu from '$lib/components/game/CardContextMenu.svelte';

	// Page data from load function
	const { data } = $props();

	// Game ID from load function (more reliable than accessing $page.params directly)
	const gameId = $derived(data.gameId);

	// UI state
	let showActionLog = $state(false);
	let showChat = $state(false);
	let showDebugOverlay = $state(false);
	let showTokenCreator = $state(false);
	let showKeyboardShortcuts = $state(false);
	let showLifeMenu = $state(false);
	let lifeMenuEl: HTMLDivElement | null = $state(null);

	// Phase 4: Deck context menu state (plan lines 860-864)
	let showDeckContextMenu = $state(false);
	let deckContextMenuPosition = $state({ x: 0, y: 0 });

	// Context menu state
	let contextMenuCard = $state<(typeof battlefieldCards)[0] | null>(null);
	let contextMenuPosition = $state({ x: 0, y: 0 });
	let actionLogRef = $state<ActionLogOverlay | undefined>(undefined);

	// Hovered card tracking for keyboard shortcuts
	let hoveredCardId = $state<string | null>(null);
	let gameChatRef = $state<GameChatOverlay | undefined>(undefined);
	let isActionLoading = $state(false);
	let initialized = $state(false);

	// Combat state (from store)
	const inCombat = $derived($isInCombat);
	const canAttackIds = $derived($canAttackCardIds);
	const attackingIds = $derived($declaredAttackerIds);
	const canBlockIds = $derived($canBlockCardIds);
	const blockingIds = $derived($assignedBlockerIds);

	// Damage assignment from special prompt (GAME_ASSIGN_DAMAGE)
	const damageAssignmentPrompt = $derived.by<DamageAssignmentPrompt | null>(() => {
		if (!prompt) return null;
		// Check if this is a damage assignment prompt (sent as 'assignDamage' type)
		if ((prompt as { type: string }).type === 'assignDamage') {
			return prompt.data as DamageAssignmentPrompt;
		}
		return null;
	});

	// Drag-drop state (from store)
	const isDragging = $derived($isDraggingStore);

	const dragCardId = $derived($draggedCardId);
	const dragCardName = $derived($draggedCardName);
	const dragPos = $derived($dragPosition);
	const isOverValidDrop = $derived($isOverValidDropZone);
	const dropZone = $derived($currentDropZone);

	// Drop zone element references
	let battlefieldDropZoneEl: HTMLDivElement | null = $state(null);
	let graveyardDropZoneEl: HTMLElement | null = $state(null);
	let exileDropZoneEl: HTMLElement | null = $state(null);
	let libraryDropZoneEl: HTMLElement | null = $state(null);
	let handDropZoneEl: HTMLElement | null = $state(null);
	let visualStackDropZoneEl: HTMLElement | null = $state(null);
	let dropZoneUnregister: (() => void) | null = null;
	let graveyardDropZoneUnregister: (() => void) | null = null;
	let exileDropZoneUnregister: (() => void) | null = null;
	let handDropZoneUnregister: (() => void) | null = null;
	let visualStackDropZoneUnregister: (() => void) | null = null;

	// Battlefield drag state (plan lines 421-428)
	let battlefieldDragStartPosition = $state<{ x: number; y: number } | null>(null);
	let battlefieldIsDragPending = $state(false);
	let commandDragStartPosition = $state<{ x: number; y: number } | null>(null);
	let commandIsDragPending = $state(false);
	const DRAG_THRESHOLD = 5;

	// Opponent panel states (for collapsing)
	let opponentExpanded = $state<Record<string, boolean>>({});

	// Phase 3: OpponentSection state - plan lines 662-666
	let selectedOpponentId = $state<string | null>(null);
	let showOpponentLifeMenu = $state(false);

	// Auto-pass setting (passes when it's opponent's turn)
	let autoPass = $state(false);

	// Mulligan state
	let mulliganCount = $state(0);
	let isMulliganLoading = $state(false);

	// Guard to prevent double initialization
	let isInitializing = $state(false);

	// Get local player ID from auth
	const localPlayerId = $derived($auth.user?.username || '');

	// Derived state from stores
	const gameState = $derived($gameStore);
	const allPlayers = $derived($players);
	const me = $derived($localPlayer);
	const otherPlayers = $derived($opponents);

	// Ensure opponentExpanded has entries for all opponents (fixes bind:expanded with undefined)
	$effect(() => {
		for (const p of otherPlayers) {
			if (!(p.playerId in opponentExpanded)) {
				opponentExpanded[p.playerId] = true;
			}
		}
	});

	const myCards = $derived($myHand);
	const myGrave = $derived($myGraveyard);
	const myExile = $derived($exile);
	const myMana = $derived($myManaPool);
	const havePriority = $derived($hasPriority);
	const phase = $derived($currentPhase);
	const step = $derived($currentStep);
	const turn = $derived($currentTurn);
	const battlefieldCards = $derived($battlefield);
	// Hovered card (for keyboard shortcuts) - derived after battlefieldCards is defined
	const hoveredCard = $derived(
		hoveredCardId ? battlefieldCards.find((c) => c.id === hoveredCardId) : null
	);
	const stackCards = $derived($stack);
	const commandCards = $derived($command);
	const prompt = $derived($pendingPrompt);

	// Plan document lines 408-419: Derived state for BattlefieldArea component
	const myBattlefield = $derived($battlefield.filter((c) => c.controllerId === localPlayerId));
	const myBattlefieldNonlands = $derived(myBattlefield.filter((c) => !isLandPermanent(c.type)));
	const myBattlefieldLands = $derived(myBattlefield.filter((c) => isLandPermanent(c.type)));
	const myCommandCards = $derived(
		$command.filter((c) => (c.ownerId || c.controllerId) === localPlayerId)
	);
	const isCommanderGame = $derived($command.length > 0);

	// Phase 3: OpponentSection derived state - plan lines 668-696
	const selectedOpponent = $derived.by(() => {
		if (otherPlayers.length === 0) return null;
		if (!selectedOpponentId || !otherPlayers.find((p) => p.playerId === selectedOpponentId)) {
			return otherPlayers[0];
		}
		return otherPlayers.find((p) => p.playerId === selectedOpponentId) || otherPlayers[0];
	});

	const opponentBattlefield = $derived.by(() => {
		return selectedOpponent
			? $battlefield.filter((c) => c.controllerId === selectedOpponent.playerId)
			: [];
	});

	const opponentBattlefieldNonlands = $derived.by(() =>
		opponentBattlefield.filter((c) => !isLandPermanent(c.type))
	);

	const opponentBattlefieldLands = $derived.by(() =>
		opponentBattlefield.filter((c) => isLandPermanent(c.type))
	);

	const opponentCommandCards = $derived.by(() => {
		return selectedOpponent
			? $command.filter((c) => (c.ownerId || c.controllerId) === selectedOpponent.playerId)
			: [];
	});

	// Combat phase detection from prompts
	const combatPromptOptions = $derived.by<ParsedCombatOptions | null>(() => {
		if (!prompt) return null;
		if (prompt.type !== 'choice') return null;
		const data = prompt.data as { choices?: string[] };
		if (!data?.choices) return null;
		const parsed = parseCombatOptions(data.choices);
		if (parsed.type === 'none') return null;
		return parsed;
	});

	const isDeclaringAttackersPhase = $derived(
		step === 'DECLARE_ATTACKERS' && combatPromptOptions?.type === 'attack'
	);

	const isDeclaringBlockersPhase = $derived(
		step === 'DECLARE_BLOCKERS' && combatPromptOptions?.type === 'block'
	);

	// Phase 4: Card+Deck context menu actions (plan lines 876-906)
	// TODO: Fix this
	const cardContextMenuActions = $derived.by(() => {
		return [
			{
				label: 'Draw Cards',
				icon: '🃏',
				onClick: () => handleDrawN(1)
			}
		];
	});
	const deckContextMenuActions = $derived.by(() => {
		type MenuAction = {
			label?: string;
			icon?: string;
			divider?: boolean;
			submenu?: MenuAction[];
			onClick?: () => void;
			disabled?: boolean;
		};

		const actions: MenuAction[] = [
			{
				label: 'Draw Cards',
				icon: '🃏',
				submenu: [
					{ label: 'Draw 1', onClick: () => handleDrawN(1) },
					{ label: 'Draw 2', onClick: () => handleDrawN(2) },
					{ label: 'Draw 3', onClick: () => handleDrawN(3) },
					{ label: 'Draw 5', onClick: () => handleDrawN(5) },
					{ label: 'Draw 7', onClick: () => handleDrawN(7) }
					// DO NOT add Custom... option yet (Phase 5)
				]
			},
			{ divider: true },
			{
				label: 'Search Library',
				icon: '🔍',
				onClick: handleSearchLibrary
			},
			{
				label: 'Shuffle Library',
				icon: '🔀',
				onClick: handleShuffleLibrary
			}
		];

		return actions;
	});

	const isGameOver = $derived($gameOver);
	const gameWinner = $derived($winner);
	const error = $derived($gameError);
	const loading = $derived($isLoading);

	// Visual stack state (client-side only)
	const showVisualStack = $derived($visualStackIsOpen);
	const visualStackItemCount = $derived($visualStackCount);

	// Rollback request state
	const pendingRollbackRequest = $derived(gameState.pendingRollbackRequest);

	// Player name map for display
	const playerNames = $derived(new Map(allPlayers.map((p) => [p.playerId, p.name])));

	// Mulligan phase detection - use server-provided value
	const isMulliganPhase = $derived(
		gameState.gameView?.isMulliganPhase ?? gameState.gameView?.state?.toLowerCase() === 'mulligan'
	);

	// Check if local player has already kept their hand (waiting for other players)
	const hasKeptHand = $derived(me?.keptHand ?? false);

	// Is it the local player's turn?
	const isYourTurn = $derived(gameState.gameView?.activePlayerId === localPlayerId);

	// Does the local player have any available actions? (server-computed)
	const myHasAvailableActions = $derived(me?.hasAvailableActions ?? false);

	// Get active player name - use server-provided value
	const activePlayerName = $derived(gameState.gameView?.activePlayerName || 'Unknown');

	// Game format - use server-provided value
	const gameFormat = $derived(gameState.gameView?.gameFormat || 'Game');

	// Priority player name - who currently has priority
	const priorityPlayerName = $derived(
		gameState.gameView?.priorityPlayerId
			? playerNames.get(gameState.gameView.priorityPlayerId) || gameState.gameView.priorityPlayerId
			: ''
	);

	/**
	 * Convert CardView from proto to GameCard for components
	 */
	function toGameCard(card: CardView): GameCard {
		return {
			id: card.id,
			name: card.name,
			manaCost: card.manaCost,
			cardType: card.type,
			power: card.power,
			toughness: card.toughness,
			imageUrl: '',
			isTapped: card.tapped,
			isSelected: false,
			ownerId: card.ownerId,
			controllerId: card.controllerId
		};
	}

	/**
	 * Convert phase/step string to GamePhase type
	 * Maps server step names to client phase keys for PhaseIndicator
	 *
	 * Server sends steps like: UNTAP, UPKEEP, DRAW, MAIN1, BEGIN_COMBAT, etc.
	 * Client expects: UNTAP, UPKEEP, DRAW, PRECOMBAT_MAIN, COMBAT, etc.
	 */
	function toGamePhase(phaseOrStep: string): GamePhase {
		const mapping: Record<string, GamePhase> = {
			// Direct matches (server step = client key)
			UNTAP: 'UNTAP',
			UPKEEP: 'UPKEEP',
			DRAW: 'DRAW',
			DECLARE_ATTACKERS: 'DECLARE_ATTACKERS',
			DECLARE_BLOCKERS: 'DECLARE_BLOCKERS',
			COMBAT_DAMAGE: 'COMBAT_DAMAGE',
			END: 'END',
			CLEANUP: 'CLEANUP',

			// Server step names that need mapping to client keys
			MAIN1: 'PRECOMBAT_MAIN',
			BEGIN_COMBAT: 'COMBAT',
			END_COMBAT: 'END_OF_COMBAT',
			MAIN2: 'POSTCOMBAT_MAIN',

			// Client-only keys (for backwards compatibility)
			BEGINNING: 'BEGINNING',
			PRECOMBAT_MAIN: 'PRECOMBAT_MAIN',
			COMBAT: 'COMBAT',
			END_OF_COMBAT: 'END_OF_COMBAT',
			POSTCOMBAT_MAIN: 'POSTCOMBAT_MAIN',
			END_OF_TURN: 'END_OF_TURN',

			// Phase names (fallback when step not available)
			ENDING: 'END'
		};
		return mapping[phaseOrStep] || 'PRECOMBAT_MAIN';
	}

	/**
	 * Parse error message to provide user-friendly feedback
	 */
	function parseGameError(err: unknown): { title: string; message: string } {
		const errorMsg = err instanceof Error ? err.message : String(err);
		const lowerMsg = errorMsg.toLowerCase();

		// Game has ended (found in match history)
		if (lowerMsg.includes('game has ended')) {
			return {
				title: 'Game Has Ended',
				message: 'This game has already finished. You can view your match history from the lobby.'
			};
		}

		// Game not found errors
		if (lowerMsg.includes('game not found') || lowerMsg.includes('no game data')) {
			return {
				title: 'Game Not Found',
				message: 'This game does not exist. The link may be invalid or the game was never created.'
			};
		}

		// Player not in game
		if (lowerMsg.includes('not part of this game')) {
			return {
				title: 'Access Denied',
				message:
					'You are not a participant in this game. You may need to join as a spectator instead.'
			};
		}

		// Session/auth errors
		if (lowerMsg.includes('session') || lowerMsg.includes('login') || lowerMsg.includes('auth')) {
			return {
				title: 'Session Expired',
				message: 'Your session has expired. Please log in again to continue.'
			};
		}

		// WebSocket/connection errors
		if (lowerMsg.includes('websocket') || lowerMsg.includes('connection')) {
			return {
				title: 'Connection Failed',
				message:
					'Unable to connect to the game server. Please check your internet connection and try again.'
			};
		}

		// Default error
		return {
			title: 'Error Loading Game',
			message: errorMsg || 'An unexpected error occurred while loading the game.'
		};
	}

	/**
	 * Initialize game connection
	 */
	async function initializeGame() {
		// Guard against double initialization
		if (isInitializing || initialized) {
			console.log('[GamePage] Skipping initialization - already in progress or completed');
			return;
		}
		isInitializing = true;

		// Read player ID directly from auth store to avoid timing issues with $derived
		const playerId = $auth.user?.username || '';

		console.log('[GamePage] initializeGame called', {
			gameId,
			playerId,
			localPlayerId,
			authUser: $auth.user,
			isAuthenticated: $auth.isAuthenticated
		});

		if (!playerId || !gameId) {
			console.error('[GamePage] Missing player ID or game ID', { playerId, gameId });
			isInitializing = false;
			return;
		}

		try {
			console.log('[GamePage] Starting game initialization...', { gameId, playerId });

			// Initialize game store first
			console.log('[GamePage] Initializing game store...');
			gameStore.initGame(gameId, playerId);

			// Try to join the game first (via HTTP) to verify it exists
			// This gives us a cleaner error if the game doesn't exist
			console.log('[GamePage] Joining game...');
			try {
				await joinGame(gameId);
				console.log('[GamePage] Joined game successfully');
			} catch (joinErr) {
				// Re-throw with context about the game not being found
				const errorMsg = joinErr instanceof Error ? joinErr.message : String(joinErr);
				if (errorMsg.toLowerCase().includes('game not found')) {
					throw new Error('Game not found');
				}
				throw joinErr;
			}

			// Now connect WebSocket for real-time updates
			const wsState = $websocketStore;
			if (wsState.state !== 'connected') {
				const token = $auth.token;
				const sessionId = token ? getSessionIdFromToken(token) : null;
				if (sessionId) {
					console.log('[GamePage] Connecting to WebSocket...');
					await websocketStore.connect(sessionId);
					console.log('[GamePage] WebSocket connected');
				} else {
					throw new Error('No session ID available - please log in again');
				}
			}

			console.log('[GamePage] Fetching initial game state...');
			const gameView = await fetchGameView(gameId, playerId);
			console.log('[GamePage] Got game state:', {
				players: gameView.players?.length,
				turn: gameView.turn,
				phase: gameView.phase,
				priorityPlayerId: gameView.priorityPlayerId
			});
			gameStore.setGameView(gameView);

			// Initialize opponent expanded states
			otherPlayers.forEach((p) => {
				opponentExpanded[p.playerId] = true;
			});

			initialized = true;
			console.log('[GamePage] Game initialization complete');
		} catch (err) {
			console.error('[GamePage] Failed to initialize game:', err);
			const parsedError = parseGameError(err);
			gameStore.setError(`${parsedError.title}: ${parsedError.message}`);
		} finally {
			isInitializing = false;
		}
	}

	/**
	 * Handle pass priority
	 */
	async function handlePassPriority() {
		if (!havePriority || isActionLoading || !gameId) return;

		isActionLoading = true;
		try {
			await passPriority(gameId);
			addLogEntry('You passed priority');
		} catch (err) {
			console.error('Failed to pass priority:', err);
		} finally {
			isActionLoading = false;
		}
	}

	/**
	 * Handle pass until end of turn (F6)
	 */
	/**
	 * Handle pass turn - passes through remaining phases until next turn begins
	 * This skips all remaining phases of your turn and lets the next player start fresh
	 */
	async function handlePassUntilEOT() {
		if (!havePriority || isActionLoading || !gameId) return;

		isActionLoading = true;
		try {
			await passUntilNextTurn(gameId);
			addLogEntry('You will pass until the next turn');
		} catch (err) {
			console.error('Failed to pass until next turn:', err);
		} finally {
			isActionLoading = false;
		}
	}

	/**
	 * Handle cast spell / play land - uses server-provided availableActions
	 */
	async function handleCastSpell() {
		console.log('[handleCastSpell] Called', {
			havePriority,
			isActionLoading,
			gameId,
			selectedCardIds: gameState.selectedCardIds,
			myCardsCount: myCards.length
		});

		if (!havePriority) {
			console.log('[handleCastSpell] No priority, returning');
			return;
		}
		if (isActionLoading) {
			console.log('[handleCastSpell] Action loading, returning');
			return;
		}
		if (!gameId) {
			console.log('[handleCastSpell] No gameId, returning');
			return;
		}

		const selectedIds = gameState.selectedCardIds;
		console.log('[handleCastSpell] Selected IDs:', selectedIds);

		if (selectedIds.length === 0) {
			console.log('[handleCastSpell] No cards selected');
			addLogEntry('No card selected');
			return;
		}

		const cardId = selectedIds[0];
		console.log('[handleCastSpell] Looking for card:', cardId);
		console.log(
			'[handleCastSpell] myCards:',
			myCards.map((c) => ({ id: c.id, name: c.name, type: c.type, actions: c.availableActions }))
		);

		const card = myCards.find((c) => c.id === cardId);
		if (!card) {
			console.log('[handleCastSpell] Card not found in hand');
			addLogEntry('Selected card not found in hand');
			return;
		}

		console.log('[handleCastSpell] Found card:', {
			id: card.id,
			name: card.name,
			type: card.type,
			actions: card.availableActions
		});

		// Use server-provided availableActions to determine action type
		const playLandAction = card.availableActions?.find(
			(a) => a.actionType === CardActionType.CARD_ACTION_PLAY_LAND
		);
		const castSpellAction = card.availableActions?.find(
			(a) => a.actionType === CardActionType.CARD_ACTION_CAST_SPELL
		);

		// Check if action is enabled
		if (playLandAction && !playLandAction.isEnabled) {
			console.log('[handleCastSpell] Play land action disabled:', playLandAction.disabledReason);
			addLogEntry(`Cannot play land: ${playLandAction.disabledReason}`);
			return;
		}
		if (castSpellAction && !castSpellAction.isEnabled) {
			console.log('[handleCastSpell] Cast spell action disabled:', castSpellAction.disabledReason);
			addLogEntry(`Cannot cast spell: ${castSpellAction.disabledReason}`);
			return;
		}

		isActionLoading = true;
		try {
			// Use server-provided action type, fall back to type check for backward compatibility
			const isLand =
				playLandAction !== undefined ||
				(!castSpellAction && card.type.toLowerCase().includes('land'));
			console.log(
				'[handleCastSpell] isLand:',
				isLand,
				'playLandAction:',
				!!playLandAction,
				'castSpellAction:',
				!!castSpellAction
			);

			if (isLand) {
				console.log('[handleCastSpell] Calling playLand with gameId:', gameId, 'cardId:', cardId);
				await playLand(gameId, cardId);
				console.log('[handleCastSpell] playLand returned successfully');
				addLogEntry(`Playing land: ${card.name}`);
			} else {
				console.log('[handleCastSpell] Calling sendPlayerString with:', card.name);
				await sendPlayerString(gameId, card.name);
				console.log('[handleCastSpell] sendPlayerString returned successfully');
				addLogEntry(`Casting spell: ${card.name}`);
			}
			gameStore.clearSelection();
		} catch (err) {
			const errorMessage = err instanceof Error ? err.message : 'Unknown error';
			console.error('[handleCastSpell] Failed:', err);
			addLogEntry(`Failed: ${errorMessage}`);
		} finally {
			isActionLoading = false;
		}
	}

	/**
	 * Play a card with optimistic UI update (used by drag-drop)
	 * Shows the card as "being played" immediately while waiting for server confirmation
	 */
	async function playCardOptimistic(cardId: string): Promise<void> {
		console.log('[playCardOptimistic] Playing card:', cardId);

		if (!havePriority || isActionLoading || !gameId) {
			console.log('[playCardOptimistic] Cannot play - no priority or loading');
			return;
		}

		const card = getCardById(cardId);
		if (!card) {
			console.log('[playCardOptimistic] Card not found:', cardId);
			toast.error('Card not found');
			return;
		}

		// Start optimistic update
		gameStore.addPendingCardPlay(cardId, card, 'hand', 'battlefield');
		gameStore.clearSelection();
		isActionLoading = true;

		try {
			// Determine if it's a land or spell
			const isLand = card.type?.toLowerCase().includes('land');

			if (isLand) {
				console.log('[playCardOptimistic] Playing land:', card.name);
				await playLand(gameId, cardId);
				addLogEntry(`Playing land: ${card.name}`);
			} else {
				console.log('[playCardOptimistic] Casting spell:', card.name);
				await sendPlayerString(gameId, card.name);
				addLogEntry(`Casting spell: ${card.name}`);
			}
			// Server will send GAME_UPDATE which clears the pending state
		} catch (err) {
			console.error('[playCardOptimistic] Failed to play card:', err);
			// Rollback on error
			gameStore.rollbackCardPlay(cardId);
			const errorMessage = err instanceof Error ? err.message : 'Unknown error';
			toast.error(`Failed to play ${card.name}: ${errorMessage}`);
		} finally {
			isActionLoading = false;
		}
	}

	/**
	 * Handle card drop on battlefield
	 */
	function handleBattlefieldDrop(cardId: string): void {
		console.log('[handleBattlefieldDrop] Card dropped:', cardId);
		// Get the source zone from drag state
		const dragState = get(dragDropStore);
		const sourceZone = dragState.sourceZone;

		if (sourceZone === 'hand') {
			// Playing from hand - use existing play logic
			playCardOptimistic(cardId);
		} else if (sourceZone && sourceZone !== 'battlefield') {
			// Moving from another zone (graveyard, exile, etc.) - use direct move
			handleZoneDrop(cardId, 'battlefield');
		}
	}

	/**
	 * Handle card drop on graveyard
	 */
	function handleGraveyardDrop(cardId: string): void {
		console.log('[handleGraveyardDrop] Card dropped:', cardId);
		handleZoneDrop(cardId, 'graveyard');
	}

	/**
	 * Handle card drop on exile
	 */
	function handleExileDrop(cardId: string): void {
		console.log('[handleExileDrop] Card dropped:', cardId);
		handleZoneDrop(cardId, 'exile');
	}

	/**
	 * Handle card drop on hand
	 */
	function handleHandDrop(cardId: string): void {
		console.log('[handleHandDrop] Card dropped:', cardId);
		handleZoneDrop(cardId, 'hand');
	}

	/**
	 * Handle card drop on stack - adds to visual stack without moving the card.
	 * The card stays in its current zone but appears on the stack for all players.
	 */
	async function handleVisualStackDrop(cardId: string): Promise<void> {
		console.log('[handleVisualStackDrop] Card dropped:', cardId);
		if (!gameId) return;

		const card = getCardById(cardId);
		try {
			await addToStack(gameId, cardId);
			addLogEntry(`Added ${card?.name || cardId} to stack`);
		} catch (err) {
			console.error('[handleVisualStackDrop] Failed to add to stack:', err);
			const errorMessage = err instanceof Error ? err.message : 'Unknown error';
			toast.error(`Failed to add to stack: ${errorMessage}`);
		}
	}

	/**
	 * Handle stack item resolve/remove - sends to server to remove from stack
	 */
	async function handleStackItemRemove(itemId: string): Promise<void> {
		console.log('[handleStackItemRemove] Removing stack item:', itemId);
		if (!gameId) return;

		try {
			await removeFromStack(gameId, itemId);
			addLogEntry('Resolved stack item');
		} catch (err) {
			console.error('[handleStackItemRemove] Failed to remove from stack:', err);
			const errorMessage = err instanceof Error ? err.message : 'Unknown error';
			toast.error(`Failed to remove from stack: ${errorMessage}`);
		}
	}

	/**
	 * Handle card drop on any zone - uses moveCard API
	 */
	async function handleZoneDrop(cardId: string, targetZone: DropZone): Promise<void> {
		if (!gameId) return;

		console.log(`[handleZoneDrop] Moving card ${cardId} to ${targetZone}`);
		try {
			await moveCard(gameId, cardId, targetZone.toUpperCase());
			addLogEntry(`Moved card to ${targetZone}`);
		} catch (err) {
			console.error('[handleZoneDrop] Failed to move card:', err);
			const errorMessage = err instanceof Error ? err.message : 'Unknown error';
			toast.error(`Failed to move card: ${errorMessage}`);
		}
	}

	/**
	 * Handle mouse down on battlefield card - start drag tracking
	 * Note: No priority check - players can always drag battlefield cards
	 * to manage game state (e.g., moving destroyed creatures to graveyard)
	 */
	function handleBattlefieldCardMouseDown(
		cardId: string,
		cardName: string,
		event: MouseEvent
	): void {
		if (event.button !== 0) return; // Only left click

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
	 * Prevent native drag events on battlefield cards
	 */
	function handleBattlefieldDragStart(event: DragEvent): void {
		event.preventDefault();
	}

	/**
	 * Plan document lines 687-729: Handle command zone card mouse down with drag threshold
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
	 * Check if a battlefield card is being dragged
	 */
	function isBattlefieldCardDragging(cardId: string): boolean {
		return isDragging && dragCardId === cardId;
	}

	/**
	 * Handle advancing to the next phase/step
	 */
	async function handleAdvancePhase() {
		if (!havePriority || isActionLoading || !gameId) return;

		isActionLoading = true;
		try {
			await advancePhase(gameId);
			addLogEntry('Advanced to next phase');
		} catch (err) {
			const errorMessage = err instanceof Error ? err.message : 'Unknown error';
			console.error('Failed to advance phase:', err);
			addLogEntry(`Failed to advance phase: ${errorMessage}`);
		} finally {
			isActionLoading = false;
		}
	}

	/**
	 * Handle untap all permanents (X key shortcut)
	 */
	async function handleUntapAll() {
		if (isActionLoading || !gameId) return;

		isActionLoading = true;
		try {
			await untapAll(gameId);
			addLogEntry('Untapped all permanents');
		} catch (err) {
			const errorMessage = err instanceof Error ? err.message : 'Unknown error';
			console.error('Failed to untap all:', err);
			addLogEntry(`Failed to untap all: ${errorMessage}`);
		} finally {
			isActionLoading = false;
		}
	}

	/**
	 * Handle global keyboard shortcuts
	 * Based on untap.in shortcuts
	 */
	function handleGlobalKeydown(event: KeyboardEvent) {
		// Ignore if typing in an input
		if (event.target instanceof HTMLInputElement || event.target instanceof HTMLTextAreaElement) {
			return;
		}

		const key = event.key.toLowerCase();

		// === GLOBAL HOTKEYS (no card needed) ===
		switch (key) {
			case '?':
				// ? - Keyboard shortcuts
				showKeyboardShortcuts = !showKeyboardShortcuts;
				event.preventDefault();
				return;

			case 'x':
				// X - Untap all permanents
				handleUntapAll();
				event.preventDefault();
				return;

			case 'c':
				// C - Draw a card from deck
				handleDrawCard();
				event.preventDefault();
				return;

			case 'v':
				// V - Shuffle your deck
				handleShuffleLibrary();
				event.preventDefault();
				return;

			case 'm':
				// M - Mulligan hand (only in mulligan phase)
				if (isMulliganPhase && !hasKeptHand) {
					handleMulligan();
					event.preventDefault();
				}
				return;

			case 'b':
				// B - Focus chat input
				showChat = true;
				// Focus will happen when chat opens
				event.preventDefault();
				return;

			case 'n':
				// N - Next phase of turn
				handleAdvancePhase();
				event.preventDefault();
				return;

			case 'e':
				// E - End Turn (pass until next turn)
				handleEndTurn();
				event.preventDefault();
				return;

			case 'w':
				// W - Insert token or card (show token creator)
				showTokenCreator = true;
				event.preventDefault();
				return;

			case 'f':
				// F - Find a card in your main deck (search library)
				handleSearchLibrary();
				event.preventDefault();
				return;
		}

		// === HOVER CARD HOTKEYS (need a hovered card) ===
		if (hoveredCard && gameId) {
			switch (key) {
				case 'j':
					// J - Face Down/Up
					handleFlipCard(hoveredCard.id, hoveredCard.faceDown ?? false);
					event.preventDefault();
					return;

				case 'l':
					// L - Alt/Default Face (Transform)
					handleTransformCard(hoveredCard.id);
					event.preventDefault();
					return;

				case 'd':
					// D - Send card to Graveyard
					handleMoveCard(hoveredCard.id, 'GRAVEYARD');
					event.preventDefault();
					return;

				case 's':
					// S - Send card to Exile
					handleMoveCard(hoveredCard.id, 'EXILE');
					event.preventDefault();
					return;

				case 'r':
					// R - Send card to hand
					handleMoveCard(hoveredCard.id, 'HAND');
					event.preventDefault();
					return;

				case 't':
					// T - Send card to top of Deck
					handleMoveCard(hoveredCard.id, 'LIBRARY');
					event.preventDefault();
					return;

				case '.':
					// . - Send card to bottom of Deck
					handleMoveCard(hoveredCard.id, 'LIBRARY_BOTTOM');
					event.preventDefault();
					return;

				case 'u':
					// U - Add +1/+1 counter
					handleAddCounter(hoveredCard.id);
					event.preventDefault();
					return;
			}
		}
	}

	/**
	 * Handle life change (+ / - buttons)
	 * Phase 3: Updated to accept optional playerId for opponent life changes - plan lines 736-749
	 */
	async function handleLifeChange(delta: number, playerId?: string) {
		const targetPlayerId = playerId || localPlayerId;
		if (!gameId || !targetPlayerId) return;
		try {
			await modifyLife(gameId, targetPlayerId, delta);
		} catch (err) {
			const errorMessage = err instanceof Error ? err.message : 'Unknown error';
			toast.error(`Failed: ${errorMessage}`);
		}
	}

	/**
	 * Handle poison counter change
	 * Phase 3: Updated to accept optional playerId for opponent poison changes - plan lines 736-749
	 */
	async function handlePoisonChange(delta: number, playerId?: string) {
		const targetPlayerId = playerId || localPlayerId;
		if (!gameId || !targetPlayerId) return;
		const player = $players.find((p) => p.playerId === targetPlayerId);
		if (!player) return;
		const currentPoison = player.poison ?? 0;
		const newValue = Math.max(0, currentPoison + delta);
		try {
			await setPlayerCounter(gameId, targetPlayerId, 'poison', newValue);
		} catch (err) {
			const errorMessage = err instanceof Error ? err.message : 'Unknown error';
			toast.error(`Failed: ${errorMessage}`);
		}
	}

	/**
	 * Draw a card from deck (C key)
	 */
	async function handleDrawCard() {
		if (isActionLoading || !gameId || !localPlayerId) return;

		isActionLoading = true;
		try {
			await drawCards(gameId, localPlayerId, 1);
			addLogEntry('Drew a card');
		} catch (err) {
			const errorMessage = err instanceof Error ? err.message : 'Unknown error';
			console.error('Failed to draw card:', err);
			addLogEntry(`Failed to draw: ${errorMessage}`);
		} finally {
			isActionLoading = false;
		}
	}

	/**
	 * Shuffle library (V key)
	 */
	async function handleShuffleLibrary() {
		if (isActionLoading || !gameId) return;

		isActionLoading = true;
		try {
			await shuffleLibrary(gameId, localPlayerId);
			addLogEntry('Shuffled library');
		} catch (err) {
			const errorMessage = err instanceof Error ? err.message : 'Unknown error';
			console.error('Failed to shuffle:', err);
			addLogEntry(`Failed to shuffle: ${errorMessage}`);
		} finally {
			isActionLoading = false;
		}
	}

	/**
	 * End turn (E key)
	 */
	async function handleEndTurn() {
		if (isActionLoading || !gameId) return;

		isActionLoading = true;
		try {
			await nextTurn(gameId);
			addLogEntry('Ended turn');
		} catch (err) {
			const errorMessage = err instanceof Error ? err.message : 'Unknown error';
			console.error('Failed to end turn:', err);
			addLogEntry(`Failed to end turn: ${errorMessage}`);
		} finally {
			isActionLoading = false;
		}
	}

	/**
	 * Search library (F key)
	 */
	async function handleSearchLibrary() {
		if (isActionLoading || !gameId) return;

		isActionLoading = true;
		try {
			await searchLibrary(gameId, 'hand', true);
			addLogEntry('Searching library...');
		} catch (err) {
			const errorMessage = err instanceof Error ? err.message : 'Unknown error';
			console.error('Failed to search library:', err);
			addLogEntry(`Failed to search: ${errorMessage}`);
		} finally {
			isActionLoading = false;
		}
	}

	/**
	 * Phase 4: Handle deck context menu (plan lines 866-874)
	 */
	function handleDeckContextMenu(event: MouseEvent) {
		event.preventDefault();
		deckContextMenuPosition = { x: event.clientX, y: event.clientY };
		showDeckContextMenu = true;
	}

	/**
	 * Phase 4: Draw N cards from library (plan lines 909-915)
	 */
	async function handleDrawN(count: number) {
		if (isActionLoading || !gameId || !localPlayerId) return;

		isActionLoading = true;
		try {
			await drawCards(gameId, localPlayerId, count);
			addLogEntry(`Drew ${count} card${count !== 1 ? 's' : ''}`);
		} catch (err) {
			const errorMessage = err instanceof Error ? err.message : 'Unknown error';
			console.error('Failed to draw cards:', err);
			addLogEntry(`Failed to draw: ${errorMessage}`);
		} finally {
			isActionLoading = false;
		}
	}

	/**
	 * Flip card face up/down (J key)
	 */
	async function handleFlipCard(cardId: string, currentlyFaceDown: boolean) {
		if (isActionLoading || !gameId) return;

		isActionLoading = true;
		try {
			await flipCard(gameId, cardId, !currentlyFaceDown);
			addLogEntry(currentlyFaceDown ? 'Turned card face up' : 'Turned card face down');
		} catch (err) {
			const errorMessage = err instanceof Error ? err.message : 'Unknown error';
			console.error('Failed to flip card:', err);
			addLogEntry(`Failed to flip: ${errorMessage}`);
		} finally {
			isActionLoading = false;
		}
	}

	/**
	 * Transform card (L key)
	 */
	async function handleTransformCard(cardId: string) {
		if (isActionLoading || !gameId) return;

		isActionLoading = true;
		try {
			await transformCard(gameId, cardId);
			addLogEntry('Transformed card');
		} catch (err) {
			const errorMessage = err instanceof Error ? err.message : 'Unknown error';
			console.error('Failed to transform card:', err);
			addLogEntry(`Failed to transform: ${errorMessage}`);
		} finally {
			isActionLoading = false;
		}
	}

	/**
	 * Move card to zone (D/S/R/T/. keys)
	 */
	async function handleMoveCard(cardId: string, zone: string) {
		if (isActionLoading || !gameId) return;

		const zoneNames: Record<string, string> = {
			GRAVEYARD: 'graveyard',
			EXILE: 'exile',
			HAND: 'hand',
			LIBRARY: 'top of library',
			LIBRARY_BOTTOM: 'bottom of library'
		};

		isActionLoading = true;
		try {
			await moveCardToZone(gameId, cardId, zone);
			addLogEntry(`Moved card to ${zoneNames[zone] || zone}`);
		} catch (err) {
			const errorMessage = err instanceof Error ? err.message : 'Unknown error';
			console.error('Failed to move card:', err);
			addLogEntry(`Failed to move: ${errorMessage}`);
		} finally {
			isActionLoading = false;
		}
	}

	/**
	 * Add +1/+1 counter to card (U key)
	 */
	async function handleAddCounter(cardId: string) {
		if (isActionLoading || !gameId) return;

		isActionLoading = true;
		try {
			await modifyCardCounter(gameId, cardId, '+1/+1', 1);
			addLogEntry('Added +1/+1 counter');
		} catch (err) {
			const errorMessage = err instanceof Error ? err.message : 'Unknown error';
			console.error('Failed to add counter:', err);
			addLogEntry(`Failed to add counter: ${errorMessage}`);
		} finally {
			isActionLoading = false;
		}
	}

	/**
	 * Handle concede
	 */
	async function handleConcede() {
		if (!gameId) return;

		const confirmed = window.confirm(
			'Are you sure you want to concede? This action cannot be undone.'
		);
		if (!confirmed) return;

		try {
			await concedeGame(gameId);
			goto('/lobby');
		} catch (err) {
			console.error('Failed to concede:', err);
			alert('Failed to concede game');
		}
	}

	/**
	 * Handle card click (for selection)
	 * Note: PlayerHand already handles selection toggle via gameStore.toggleCardSelection
	 * This handler is for additional logic (logging)
	 */
	function handleCardClick(cardId: string) {
		const card = myCards.find((c) => c.id === cardId);
		if (card) {
			addLogEntry(`Selected: ${card.name}`);
		}
	}

	/**
	 * Handle battlefield card click
	 */
	async function handleBattlefieldCardClick(cardId: string) {
		// Handle declare attackers phase - toggle attacker on click
		if (isDeclaringAttackersPhase && combatPromptOptions) {
			const card = battlefieldCards.find((c) => c.id === cardId);
			if (!card) return;

			// Check if this card can attack
			const canAttack = canAttackIds.has(cardId);
			if (canAttack) {
				// Get defenders for this attacker
				const options = combatPromptOptions;
				const validDefenders = options!.attackOptions
					.filter((opt) => opt.cardId === cardId)
					.map((opt) => opt.defenderId);

				if (validDefenders.length > 0) {
					const isCurrentlyAttacking = attackingIds.has(cardId);
					combatStore.toggleAttacker(cardId, validDefenders[0]);

					if (isCurrentlyAttacking) {
						addLogEntry(`${card.name} will not attack`);
					} else {
						// Get defender name for log
						const defenders = getDefenderTargets();
						const defenderName =
							defenders.find((d) => d.id === validDefenders[0])?.name || 'opponent';
						addLogEntry(`${card.name} attacks ${defenderName}`);
					}
				}
			}
			return;
		}

		// Find the card on the battlefield
		const card = battlefieldCards.find((c) => c.id === cardId);
		if (!card) {
			console.log('[handleBattlefieldCardClick] Card not found:', cardId);
			return;
		}

		// Debug: Log card info for mana ability checking
		console.log('[handleBattlefieldCardClick] Clicked card:', {
			id: card.id,
			name: card.name,
			type: card.type,
			tapped: card.tapped,
			controllerId: card.controllerId,
			availableActions: card.availableActions,
			availableActionsLength: card.availableActions?.length,
			availableActionsJSON: JSON.stringify(card.availableActions),
			havePriority,
			isActionLoading,
			gameId
		});

		// Log each action for debugging
		if (card.availableActions && card.availableActions.length > 0) {
			card.availableActions.forEach((action, i) => {
				console.log(`[handleBattlefieldCardClick] Action ${i}:`, {
					actionType: action.actionType,
					actionTypeValue: typeof action.actionType,
					expectedManaAbilityType: CardActionType.CARD_ACTION_ACTIVATE_MANA_ABILITY,
					isMatch: action.actionType === CardActionType.CARD_ACTION_ACTIVATE_MANA_ABILITY,
					displayText: action.displayText,
					isEnabled: action.isEnabled
				});
			});
		} else {
			console.log('[handleBattlefieldCardClick] No available actions on card');
		}

		// Check for mana ability activation
		// Only for cards we control, and only if we have priority
		if (havePriority && !isActionLoading && gameId) {
			// Note: actionType can be either a number (enum) or string (JSON serialized)
			// so we check for both using string comparison
			const manaAbilityAction = card.availableActions?.find(
				(a) =>
					(a.actionType === CardActionType.CARD_ACTION_ACTIVATE_MANA_ABILITY ||
						String(a.actionType) === 'CARD_ACTION_ACTIVATE_MANA_ABILITY') &&
					a.isEnabled
			);

			console.log('[handleBattlefieldCardClick] Mana ability action found:', manaAbilityAction);

			if (manaAbilityAction) {
				isActionLoading = true;
				try {
					await activateManaAbility(gameId, cardId);
					addLogEntry(`${card.name} - ${manaAbilityAction.displayText}`);
				} catch (err) {
					const errorMessage = err instanceof Error ? err.message : 'Unknown error';
					console.error('Failed to activate mana ability:', err);
					addLogEntry(`Failed to tap for mana: ${errorMessage}`);
				} finally {
					isActionLoading = false;
				}
				return;
			}
		}

		// Default behavior: tap/untap the permanent
		if (gameId) {
			try {
				await tapUntap(gameId, cardId, !card.tapped);
				addLogEntry(`${card.tapped ? 'Untapped' : 'Tapped'} ${card.name}`);
			} catch (err) {
				const errorMessage = err instanceof Error ? err.message : 'Unknown error';
				console.error('Failed to tap/untap:', err);
				addLogEntry(`Failed to tap/untap: ${errorMessage}`);
			}
		}
	}

	/**
	 * Handle right-click context menu on battlefield cards
	 */
	function handleCardContextMenu(event: MouseEvent, card: (typeof battlefieldCards)[0]) {
		event.preventDefault();
		contextMenuCard = card;
		contextMenuPosition = { x: event.clientX, y: event.clientY };
	}

	/**
	 * Close context menu
	 */
	function closeContextMenu() {
		contextMenuCard = null;
	}

	/**
	 * Handle prompt responses
	 */
	async function handlePromptYes() {
		if (prompt?.type === 'ask' && gameId) {
			await sendPlayerBoolean(gameId, true);
			gameStore.clearPrompt();
		}
	}

	async function handlePromptNo() {
		if (prompt?.type === 'ask' && gameId) {
			await sendPlayerBoolean(gameId, false);
			gameStore.clearPrompt();
		}
	}

	async function handleChoiceSelect(choice: string) {
		if (prompt?.type === 'choice' && gameId) {
			await sendPlayerString(gameId, choice);
			gameStore.clearPrompt();
		}
	}

	/**
	 * Handle keeping hand during mulligan
	 */
	async function handleKeepHand() {
		if (!gameId || isMulliganLoading) return;

		isMulliganLoading = true;
		try {
			await keepHand(gameId);
			addLogEntry('You kept your hand');
		} catch (err) {
			console.error('Failed to keep hand:', err);
			gameStore.setError(err instanceof Error ? err.message : 'Failed to keep hand');
		} finally {
			isMulliganLoading = false;
		}
	}

	/**
	 * Handle mulligan
	 */
	async function handleMulligan() {
		if (!gameId || isMulliganLoading) return;

		isMulliganLoading = true;
		try {
			await mulligan(gameId);
			mulliganCount++;
			addLogEntry(`You took mulligan #${mulliganCount}`);
		} catch (err) {
			console.error('Failed to mulligan:', err);
			gameStore.setError(err instanceof Error ? err.message : 'Failed to mulligan');
		} finally {
			isMulliganLoading = false;
		}
	}

	/**
	 * Add message to action log
	 */
	function addLogEntry(message: string) {
		if (actionLogRef) {
			actionLogRef.addAction('system', message, { type: 'system' });
		}
	}

	/**
	 * Handle rollback request from action log
	 */
	async function handleRequestRollback(messageId: number) {
		if (!gameId) return;

		try {
			const result = await requestRollback(gameId, messageId);
			if (result.success) {
				if (result.requiresConsent) {
					addLogEntry('Rollback request sent to opponent(s)');
				} else {
					addLogEntry('Rollback performed');
				}
			} else {
				gameStore.setError(result.error || 'Failed to request rollback');
			}
		} catch (err) {
			console.error('Failed to request rollback:', err);
			gameStore.setError(err instanceof Error ? err.message : 'Failed to request rollback');
		}
	}

	/**
	 * Handle rollback consent response (approve)
	 */
	async function handleApproveRollback() {
		if (!gameId || !pendingRollbackRequest) return;

		try {
			const result = await respondToRollback(gameId, pendingRollbackRequest.requestId, true);
			if (!result.success) {
				gameStore.setError(result.error || 'Failed to approve rollback');
			}
			// Clear the pending request - server will send GAME_ROLLBACK_COMPLETE
			gameStore.rollbackCardPlay();
		} catch (err) {
			console.error('Failed to approve rollback:', err);
			gameStore.setError(err instanceof Error ? err.message : 'Failed to approve rollback');
		}
	}

	/**
	 * Handle rollback consent response (deny)
	 */
	async function handleDenyRollback() {
		if (!gameId || !pendingRollbackRequest) return;

		try {
			const result = await respondToRollback(gameId, pendingRollbackRequest.requestId, false);
			if (!result.success) {
				gameStore.setError(result.error || 'Failed to deny rollback');
			}
			// Clear the pending request
			gameStore.update((s) => ({ ...s, pendingRollbackRequest: null }));
			addLogEntry('Rollback request denied');
		} catch (err) {
			console.error('Failed to deny rollback:', err);
			gameStore.setError(err instanceof Error ? err.message : 'Failed to deny rollback');
		}
	}

	/**
	 * Handle combat phase completion
	 */
	function handleCombatComplete() {
		// Clear any pending prompt after combat action
		gameStore.clearPrompt();
		addLogEntry('Combat action completed');
	}

	/**
	 * Get defender targets for declare attackers phase
	 */
	function getDefenderTargets(): DefenderTarget[] {
		const defenders: DefenderTarget[] = [];

		// Add opponents as defenders
		for (const opponent of otherPlayers) {
			defenders.push({
				id: opponent.playerId,
				name: opponent.name,
				type: 'player',
				life: opponent.life
			});
		}

		// Add planeswalkers controlled by opponents
		for (const card of battlefieldCards) {
			if (
				card.type?.toLowerCase().includes('planeswalker') &&
				card.controllerId !== localPlayerId
			) {
				defenders.push({
					id: card.id,
					name: card.name,
					type: 'planeswalker',
					loyalty: parseInt(card.loyalty || '0', 10)
				});
			}
		}

		return defenders;
	}

	/**
	 * Get attacking creatures info for blockers phase
	 */
	function getAttackingCreatures(): DeclaredAttacker[] {
		const attackers: DeclaredAttacker[] = [];

		// Get combat info from game view if available
		const combatView = gameState.gameView?.combat;
		if (combatView?.groups) {
			for (const group of combatView.groups) {
				for (const attackerId of group.attackers || []) {
					const card = battlefieldCards.find((c) => c.id === attackerId);
					if (card) {
						attackers.push({
							cardId: attackerId,
							cardName: card.name,
							defenderId: group.defendingPlayerId || '',
							defenderName: playerNames.get(group.defendingPlayerId || '') || 'Unknown'
						});
					}
				}
			}
		}

		return attackers;
	}

	/**
	 * Get cards controlled by a specific player on battlefield
	 */
	function getPlayerBattlefieldCards(playerId: string): CardView[] {
		return battlefieldCards.filter((c) => c.controllerId === playerId);
	}

	// Plan document lines 408-419: Helper function for battlefield separation
	function isLandPermanent(cardType?: string | null): boolean {
		return !!cardType && /\bland\b/i.test(cardType);
	}

	/**
	 * Format life total
	 */
	function formatLife(life: number): string {
		return life.toString();
	}

	/**
	 * Get opponent position for layout
	 */
	function getOpponentPosition(index: number, total: number): 'top' | 'left' | 'right' {
		if (total === 1) return 'top';
		if (total === 2) return index === 0 ? 'left' : 'right';
		if (total === 3) {
			if (index === 0) return 'left';
			if (index === 1) return 'top';
			return 'right';
		}
		return 'top';
	}

	// Track if we're in the middle of an auto-pass to prevent double-triggers
	let isAutoPassPending = $state(false);

	/**
	 * Auto-pass effect - passes priority on opponent's turn
	 *
	 * This is a rules-light auto-pass: it only checks whose turn it is,
	 * not what actions are available (which would require rules enforcement).
	 */
	$effect(() => {
		// Skip if not ready for auto-pass
		if (!initialized || isMulliganPhase || isActionLoading || isAutoPassPending || !gameId) {
			return;
		}

		// Must have priority to pass
		if (!havePriority) {
			return;
		}

		// Auto-pass on opponent's turn
		if (autoPass && !isYourTurn) {
			console.log("[AutoPass] Triggering auto-pass: opponent's turn");
			isAutoPassPending = true;

			// Capture gameId in local scope for the async callback
			const currentGameId = gameId;

			// Use setTimeout to avoid synchronous state updates in effect
			setTimeout(async () => {
				try {
					if (currentGameId) {
						await passPriority(currentGameId);
						addLogEntry('Auto-passed');
					}
				} catch (err) {
					console.error('[AutoPass] Failed to auto-pass:', err);
				} finally {
					isAutoPassPending = false;
				}
			}, 50); // Small delay for smoother UX
		}
	});

	// Click outside handler for life menu
	$effect(() => {
		if (!showLifeMenu) return;

		const handleClickOutside = (event: MouseEvent) => {
			if (lifeMenuEl && !lifeMenuEl.contains(event.target as Node)) {
				showLifeMenu = false;
			}
		};

		// Delay to prevent immediate close from the click that opened it
		const timeoutId = setTimeout(() => {
			document.addEventListener('click', handleClickOutside);
		}, 10);

		return () => {
			clearTimeout(timeoutId);
			document.removeEventListener('click', handleClickOutside);
		};
	});

	// Register battlefield drop zone when element is available
	$effect(() => {
		if (battlefieldDropZoneEl && !dropZoneUnregister) {
			console.log('[GamePage] Registering battlefield drop zone');
			dropZoneUnregister = dragDropStore.registerDropZone({
				id: 'battlefield',
				type: 'battlefield',
				element: battlefieldDropZoneEl,
				accepts: (_cardId, sourceZone) => {
					// Accept cards from any zone except battlefield when player has priority
					return sourceZone !== 'battlefield' && havePriority;
				},
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

	// Register graveyard drop zone
	$effect(() => {
		if (graveyardDropZoneEl && !graveyardDropZoneUnregister) {
			console.log('[GamePage] Registering graveyard drop zone', {
				element: graveyardDropZoneEl,
				rect: graveyardDropZoneEl.getBoundingClientRect()
			});
			graveyardDropZoneUnregister = dragDropStore.registerDropZone({
				id: 'graveyard',
				type: 'graveyard',
				element: graveyardDropZoneEl,
				accepts: (_cardId, sourceZone) => {
					// Accept cards from any zone except graveyard when player has priority
					// Use get() to read current store value at call time
					const hasPrio = get(hasPriority);
					console.log('[Graveyard accepts]', {
						sourceZone,
						hasPrio,
						result: sourceZone !== 'graveyard' && hasPrio
					});
					return sourceZone !== 'graveyard' && hasPrio;
				},
				onDrop: handleGraveyardDrop
			});
		}

		return () => {
			if (graveyardDropZoneUnregister) {
				graveyardDropZoneUnregister();
				graveyardDropZoneUnregister = null;
			}
		};
	});

	// Register exile drop zone
	$effect(() => {
		if (exileDropZoneEl && !exileDropZoneUnregister) {
			console.log('[GamePage] Registering exile drop zone', {
				element: exileDropZoneEl,
				rect: exileDropZoneEl.getBoundingClientRect()
			});
			exileDropZoneUnregister = dragDropStore.registerDropZone({
				id: 'exile',
				type: 'exile',
				element: exileDropZoneEl,
				accepts: (_cardId, sourceZone) => {
					// Accept cards from any zone except exile when player has priority
					// Use get() to read current store value at call time
					const hasPrio = get(hasPriority);
					console.log('[Exile accepts]', {
						sourceZone,
						hasPrio,
						result: sourceZone !== 'exile' && hasPrio
					});
					return sourceZone !== 'exile' && hasPrio;
				},
				onDrop: handleExileDrop
			});
		}

		return () => {
			if (exileDropZoneUnregister) {
				exileDropZoneUnregister();
				exileDropZoneUnregister = null;
			}
		};
	});

	// Register hand drop zone
	$effect(() => {
		if (handDropZoneEl && !handDropZoneUnregister) {
			console.log('[GamePage] Registering hand drop zone');
			handDropZoneUnregister = dragDropStore.registerDropZone({
				id: 'hand',
				type: 'hand',
				element: handDropZoneEl,
				accepts: (_cardId, sourceZone) => {
					// Accept cards from any zone except hand when player has priority
					return sourceZone !== 'hand' && havePriority;
				},
				onDrop: handleHandDrop
			});
		}

		return () => {
			if (handDropZoneUnregister) {
				handDropZoneUnregister();
				handDropZoneUnregister = null;
			}
		};
	});

	// Register visual stack drop zone
	$effect(() => {
		if (visualStackDropZoneEl && !visualStackDropZoneUnregister) {
			console.log('[GamePage] Registering visual stack drop zone');
			visualStackDropZoneUnregister = dragDropStore.registerDropZone({
				id: 'visual-stack',
				type: 'stack',
				element: visualStackDropZoneEl,
				accepts: (_cardId, sourceZone) => {
					// Accept cards from any zone except stack when player has priority
					return sourceZone !== 'stack' && havePriority;
				},
				onDrop: handleVisualStackDrop
			});
		}

		return () => {
			if (visualStackDropZoneUnregister) {
				visualStackDropZoneUnregister();
				visualStackDropZoneUnregister = null;
			}
		};
	});

	// Track which server stack item IDs have been synced to visual stack
	let lastSyncedStackIds = new Set<string>();
	let stackUnsubscribe: (() => void) | null = null;

	// Sync server stack cards to visual stack - using store subscription for reliability
	function syncServerStackToVisualStack(currentStack: typeof stackCards) {
		console.log(
			'[GamePage] Server stack sync running, cards:',
			currentStack.length,
			currentStack.map((c) => ({ id: c.id, name: c.name }))
		);
		console.log('[GamePage] Already synced IDs:', [...lastSyncedStackIds]);

		const serverIds = new Set(currentStack.map((c) => c.id));

		// Add new server stack items to visual stack
		for (const card of currentStack) {
			if (!lastSyncedStackIds.has(card.id)) {
				console.log('[GamePage] Syncing server stack item to visual stack:', card.name, card.id);
				visualStackStore.addItem(card.id, card.name, 'stack', {
					imageUrl: card?.imageUrl || getScryfallImageUrl(card.name),
					controllerId: card.controllerId,
					note: 'Spell',
					// Use server's stack item ID as localId for sync
					localId: card.id
				});
				lastSyncedStackIds.add(card.id);

				// Auto-open visual stack when items are added
				visualStackStore.setOpen(true);
			}
		}

		// Remove items from visual stack that are no longer on server stack
		for (const id of lastSyncedStackIds) {
			if (!serverIds.has(id)) {
				console.log('[GamePage] Removing resolved stack item from visual stack:', id);
				visualStackStore.removeItem(id);
				lastSyncedStackIds.delete(id);
			}
		}
	}

	// Initialize on mount
	onMount(() => {
		console.log('[GamePage] onMount called', {
			isAuthenticated: $auth.isAuthenticated,
			user: $auth.user,
			gameId,
			localPlayerId
		});

		if (!$auth.isAuthenticated) {
			console.log('[GamePage] Not authenticated, redirecting to login');
			goto('/login');
			return;
		}

		// Initialize game - gameId should always be available from load function
		initializeGame();

		// Subscribe to stack changes to sync with visual stack
		console.log('[GamePage] Setting up stack subscription');
		stackUnsubscribe = stack.subscribe((currentStack) => {
			console.log('[GamePage] Stack subscription triggered, stack length:', currentStack.length);
			syncServerStackToVisualStack(currentStack);
		});
	});

	// Cleanup on destroy
	onDestroy(() => {
		gameStore.reset();
		if (stackUnsubscribe) {
			stackUnsubscribe();
			stackUnsubscribe = null;
		}
	});
</script>

<svelte:window onkeydown={handleGlobalKeydown} />

<svelte:head>
	<title>{gameId ? `Game ${gameId}` : 'Loading Game'} - MAGE</title>
</svelte:head>

<div class="game-container" class:has-priority={havePriority}>
	{#if loading && !initialized}
		<div class="loading-overlay">
			<div class="spinner"></div>
			<p>Loading game...</p>
		</div>
	{:else if error && !gameState.gameView}
		{@const parsedError = (() => {
			// Check if error already includes title: message format
			if (error.includes(': ')) {
				const [title, ...rest] = error.split(': ');
				return { title, message: rest.join(': ') };
			}
			return parseGameError(new Error(error));
		})()}
		<div class="error-overlay">
			<div class="error-icon">
				{#if parsedError.title.toLowerCase().includes('has ended')}
					🏁
				{:else if parsedError.title.toLowerCase().includes('not found')}
					❓
				{:else if parsedError.title.toLowerCase().includes('access denied')}
					🚫
				{:else if parsedError.title.toLowerCase().includes('session') || parsedError.title
						.toLowerCase()
						.includes('expired')}
					🔐
				{:else if parsedError.title.toLowerCase().includes('connection')}
					📡
				{:else}
					⚠️
				{/if}
			</div>
			<h2 class="error-title">{parsedError.title}</h2>
			<p class="error-message">{parsedError.message}</p>
			<div class="error-actions">
				<button class="btn-primary" onclick={() => goto('/lobby')}>Return to Lobby</button>
				{#if parsedError.title.toLowerCase().includes('connection')}
					<button
						class="btn-secondary"
						onclick={() => {
							isInitializing = false;
							initialized = false;
							initializeGame();
						}}
					>
						Try Again
					</button>
				{/if}
			</div>
		</div>
	{:else if isGameOver}
		<div class="game-over-overlay">
			<div class="game-over-content">
				<h2>Game Over</h2>
				<p class="winner-text">
					{gameWinner ? `Winner: ${playerNames.get(gameWinner) || gameWinner}` : 'Draw'}
				</p>
				<button class="btn-primary" onclick={() => goto('/lobby')}>Return to Lobby</button>
			</div>
		</div>
	{:else if isMulliganPhase}
		<MulliganDialog
			cards={myCards}
			{mulliganCount}
			onKeep={handleKeepHand}
			onMulligan={handleMulligan}
			isLoading={isMulliganLoading}
			{hasKeptHand}
		/>
	{:else}
		<!-- Game Header - Clean UX answering key questions -->
		<GameHeader
			{turn}
			{activePlayerName}
			{priorityPlayerName}
			localPlayerName={localPlayerId}
			hasPriority={havePriority}
			currentPhase={toGamePhase(step || phase)}
			onLogClick={() => (showActionLog = true)}
			onConcedeClick={handleConcede}
		/>

		<!-- Floating action buttons -->
		<div class="floating-actions">
			<button
				class="floating-btn stack-btn"
				class:has-items={visualStackItemCount > 0 || stackCards.length > 0}
				onclick={() => visualStackStore.toggleOpen()}
				title="Stack ({visualStackItemCount} items)"
			>
				📚 {#if visualStackItemCount > 0 || stackCards.length > 0}<span class="badge"
						>{visualStackItemCount || stackCards.length}</span
					>{/if}
			</button>
			<button class="floating-btn" onclick={() => (showChat = true)} title="Game Chat"> 💬 </button>
			<button
				class="floating-btn"
				onclick={() => (showKeyboardShortcuts = true)}
				title="Keyboard shortcuts (?)"
				aria-label="Keyboard shortcuts"
			>
				<Keyboard size={18} aria-hidden="true" />
			</button>
		</div>

		<!-- Mana Payment Modal -->
		{#if prompt && prompt.type === 'mana' && gameId}
			<ManaPayment
				{gameId}
				manaData={prompt.data as GamePlayManaData}
				onComplete={() => gameStore.clearPrompt()}
				onCancel={() => gameStore.clearPrompt()}
			/>
		{/if}

		<!-- X Mana Selector Modal -->
		{#if prompt && prompt.type === 'xmana' && gameId}
			<XManaSelector
				{gameId}
				xManaData={prompt.data as GamePlayXManaData}
				onComplete={() => gameStore.clearPrompt()}
				onCancel={() => gameStore.clearPrompt()}
			/>
		{/if}

		<!-- Library Search Modal -->
		{#if prompt && prompt.type === 'librarySearch' && gameId}
			<LibrarySearch
				{gameId}
				librarySearchData={prompt.data as import('$lib/generated/mage/v1/models').LibrarySearchView}
				onComplete={() => gameStore.clearPrompt()}
				onCancel={() => gameStore.clearPrompt()}
			/>
		{/if}

		<!-- Combat: Declare Attackers -->
		{#if isDeclaringAttackersPhase && combatPromptOptions && gameId}
			<DeclareAttackers
				{gameId}
				options={combatPromptOptions}
				battlefieldCards={getPlayerBattlefieldCards(localPlayerId)}
				defenders={getDefenderTargets()}
				onComplete={handleCombatComplete}
			/>
		{/if}

		<!-- Combat: Declare Blockers -->
		{#if isDeclaringBlockersPhase && combatPromptOptions && gameId}
			<DeclareBlockers
				{gameId}
				options={combatPromptOptions}
				{battlefieldCards}
				attackingCreatures={getAttackingCreatures()}
				onComplete={handleCombatComplete}
			/>
		{/if}

		<!-- Combat: Assign Damage -->
		{#if damageAssignmentPrompt && gameId}
			<AssignDamage {gameId} prompt={damageAssignmentPrompt} onComplete={handleCombatComplete} />
		{/if}

		<!-- Prompt Overlay (non-mana prompts) -->
		{#if prompt && !['mana', 'xmana'].includes(prompt.type)}
			<div class="prompt-overlay">
				<div class="prompt-content">
					<p class="prompt-message">{prompt.message}</p>
					{#if prompt.type === 'ask'}
						<div class="prompt-buttons">
							<button class="btn-yes" onclick={handlePromptYes}>Yes</button>
							<button class="btn-no" onclick={handlePromptNo}>No</button>
						</div>
					{:else if prompt.type === 'choice'}
						{@const choiceData = prompt.data as { choices: string[] }}
						<div class="choice-buttons">
							{#each choiceData.choices as choice}
								<button class="btn-choice" onclick={() => handleChoiceSelect(choice)}>
									{choice}
								</button>
							{/each}
						</div>
					{/if}
				</div>
			</div>
		{/if}

		<!-- Main Game Area with Sidebar -->
		<div class="game-area-wrapper">
			<main class="game-layout" class:four-player={otherPlayers.length >= 3}>
				<!-- Phase 3: OpponentSection integration - plan lines 516-594 -->
				{#if otherPlayers.length === 1}
					<!-- 1v1 layout - single opponent (plan lines 516-538) -->
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
								contextMenuCard = battlefieldCards.find((c) => c.id === cardId) || null;
								contextMenuPosition = { x: 0, y: 0 };
							}}
						/>
					{/if}
				{:else}
					<!-- Multi-opponent layouts (plan lines 542-594) -->
					<!-- Grid layout for large screens -->
					<div class="opponents-grid opponents-grid-large">
						{#each otherPlayers as opponent (opponent.playerId)}
							{@const oppBattlefield = $battlefield.filter(
								(c) => c.controllerId === opponent.playerId
							)}
							{@const oppBattlefieldNonlands = oppBattlefield.filter(
								(c) => !isLandPermanent(c.type)
							)}
							{@const oppBattlefieldLands = oppBattlefield.filter((c) => isLandPermanent(c.type))}
							{@const oppCommandCards = $command.filter(
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
									contextMenuCard = battlefieldCards.find((c) => c.id === cardId) || null;
									contextMenuPosition = { x: 0, y: 0 };
								}}
							/>
						{/each}
					</div>

					<!-- Single opponent with cycling for small screens -->
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
									contextMenuCard = battlefieldCards.find((c) => c.id === cardId) || null;
									contextMenuPosition = { x: 0, y: 0 };
								}}
							/>
						{/if}
					</div>
				{/if}

				<!-- Central Battlefield Area - Plan document lines 322-343 -->
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
						contextMenuCard = battlefieldCards.find((c) => c.id === cardId) || null;
						contextMenuPosition = { x: 0, y: 0 };
					}}
					onCommandCardMouseDown={handleCommandCardMouseDown}
					onCardHover={(cardId) => (hoveredCardId = cardId)}
					battlefieldDropZoneRef={(el) => (battlefieldDropZoneEl = el)}
					commandDropZoneRef={(el) => {
						// Command drop zone not currently used in game page
					}}
				/>

				<!-- Player Info & Zones Row - Using PlayerInfoRow component (see plan lines 35-51) -->
				{#if me}
					<PlayerInfoRow
						player={{
							name: 'You',
							life: me.life,
							poison: me.poison ?? 0,
							libraryCount: me.libraryCount ?? 0
						}}
						graveyard={myGrave}
						exile={myExile}
						mana={myMana}
						{showLifeMenu}
						onLifeChange={handleLifeChange}
						onPoisonChange={handlePoisonChange}
						onToggleLifeMenu={() => (showLifeMenu = !showLifeMenu)}
						onSearchLibrary={handleSearchLibrary}
						onDeckContextMenu={() => {}}
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
						onCardClick={handleCardClick}
						size="normal"
						currentPhase={step || phase}
						canDrag={havePriority}
					/>
				</div>
			</main>

			<!-- Visual Stack Sidebar (inline, not overlay) -->
			<div
				bind:this={visualStackDropZoneEl}
				class="visual-stack-sidebar-container"
				class:drag-active={isDragging}
				class:drag-valid={isDragging && isOverValidDrop && dropZone === 'stack'}
			>
				<VisualStack
					isOpen={showVisualStack}
					onResolve={handleStackItemRemove}
					onRemove={handleStackItemRemove}
				/>
			</div>
		</div>
		<!-- End game-area-wrapper -->

		<!-- Overlay Panels -->
		<ActionLogOverlay
			bind:this={actionLogRef}
			bind:open={showActionLog}
			onRequestRollback={handleRequestRollback}
		/>
		<GameChatOverlay bind:this={gameChatRef} gameId={gameId || ''} bind:open={showChat} />
		<KeyboardShortcutsModal bind:open={showKeyboardShortcuts} mode="game" />

		<!-- Priority Action Bar (Docked at bottom) -->
		<PriorityActionBar
			{gameId}
			hasPriority={havePriority}
			activePlayerId={gameState.gameView?.activePlayerId || ''}
			{localPlayerId}
			{activePlayerName}
			currentPhase={phase}
			canPassPriority={havePriority}
			isLoading={isActionLoading}
			onPassPriority={handlePassPriority}
			onPassUntilEOT={handlePassUntilEOT}
			onCastSpell={handleCastSpell}
			onAdvancePhase={handleAdvancePhase}
			onCreateToken={() => (showTokenCreator = true)}
			bind:autoPass
		/>

		<!-- Floating Debug Button -->
		<button class="debug-fab" onclick={() => (showDebugOverlay = true)} title="Open Debug View">
			🔧
		</button>

		<!-- Debug Overlay Modal -->
		<DebugOverlay
			bind:open={showDebugOverlay}
			gameId={gameId || ''}
			{localPlayerId}
			{gameState}
			{allPlayers}
			{battlefieldCards}
			{stackCards}
			{commandCards}
			{turn}
			{phase}
			{havePriority}
			{isMulliganPhase}
			{gameFormat}
			{isGameOver}
			{gameWinner}
			{activePlayerName}
			{prompt}
			{error}
			onClose={() => (showDebugOverlay = false)}
		/>

		<!-- Card Context Menu -->
		{#if contextMenuCard && gameId}
			<CardContextMenu
				card={contextMenuCard}
				position={contextMenuPosition}
				onClose={closeContextMenu}
				actions={cardContextMenuActions}
			/>
		{/if}

		<!-- Phase 4: Deck Context Menu (plan lines 917-927) -->
		{#if showDeckContextMenu && me}
			<DeckContextMenu
				position={deckContextMenuPosition}
				deckCount={me.libraryCount ?? 0}
				playerName="You"
				actions={deckContextMenuActions}
				onClose={() => (showDeckContextMenu = false)}
			/>
		{/if}

		<!-- Token Creator Dialog -->
		{#if showTokenCreator && gameId}
			<TokenCreator {gameId} onClose={() => (showTokenCreator = false)} />
		{/if}

		<!-- Rollback Consent Dialog -->
		{#if pendingRollbackRequest}
			<RollbackConsentDialog
				request={pendingRollbackRequest}
				onApprove={handleApproveRollback}
				onDeny={handleDenyRollback}
			/>
		{/if}

		<!-- Drag Ghost - Card following the cursor during drag -->
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
	/* Container with priority glow effect */
	.game-container {
		position: fixed;
		inset: 0;
		background: #0a0d12;
		color: white;
		display: flex;
		flex-direction: column;
		overflow: hidden;
		transition: box-shadow 0.5s ease;
	}

	/* Full-screen priority glow */
	.game-container.has-priority {
		box-shadow: inset 0 0 80px rgba(251, 191, 36, 0.08);
	}

	.game-container.has-priority::before {
		content: '';
		position: absolute;
		inset: 0;
		border: 2px solid rgba(251, 191, 36, 0.25);
		pointer-events: none;
		z-index: 1000;
		animation: priority-border-pulse 2s ease-in-out infinite;
	}

	@keyframes priority-border-pulse {
		0%,
		100% {
			border-color: rgba(251, 191, 36, 0.25);
		}
		50% {
			border-color: rgba(251, 191, 36, 0.5);
		}
	}

	/* Loading & Error & Game Over States */
	.loading-overlay,
	.error-overlay,
	.game-over-overlay {
		position: absolute;
		inset: 0;
		display: flex;
		flex-direction: column;
		align-items: center;
		justify-content: center;
		background: #0a0d12;
		z-index: 100;
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
		margin-bottom: 0.5rem;
		opacity: 0.9;
	}

	.error-title {
		font-size: 1.5rem;
		color: #f8fafc;
		margin: 0 0 0.5rem;
		font-weight: 600;
	}

	.error-message {
		color: #94a3b8;
		max-width: 400px;
		text-align: center;
		line-height: 1.6;
		margin: 0;
	}

	.error-actions {
		display: flex;
		gap: 0.75rem;
		margin-top: 1rem;
	}

	.btn-secondary {
		padding: 0.75rem 1.5rem;
		border-radius: 8px;
		font-weight: 500;
		cursor: pointer;
		transition: all 0.2s;
		background: transparent;
		border: 1px solid #374151;
		color: #9ca3af;
	}

	.btn-secondary:hover {
		background: #1f2937;
		border-color: #4b5563;
		color: #f8fafc;
	}

	.game-over-content {
		text-align: center;
		background: #1a1f2e;
		padding: 2rem 3rem;
		border-radius: 12px;
		border: 2px solid #667eea;
	}

	.game-over-content h2 {
		font-size: 2rem;
		margin-bottom: 1rem;
	}

	.winner-text {
		font-size: 1.25rem;
		color: #fbbf24;
		margin-bottom: 1.5rem;
	}

	/* Floating action buttons */
	.floating-actions {
		position: fixed;
		top: 100px;
		right: 16px;
		display: flex;
		flex-direction: column;
		gap: 0.5rem;
		z-index: 50;
	}

	.floating-btn {
		width: 44px;
		height: 44px;
		display: flex;
		align-items: center;
		justify-content: center;
		background: rgba(26, 31, 46, 0.95);
		border: 1px solid #2a3441;
		border-radius: 10px;
		font-size: 1.25rem;
		cursor: pointer;
		transition: all 0.2s;
		color: #fff;
		position: relative;
		background: rgba(0, 0, 0, 0.85);
	}

	.floating-btn:hover {
		background: #2a3441;
		border-color: #374151;
		transform: scale(1.05);
	}

	.floating-btn .badge {
		position: absolute;
		top: -4px;
		right: -4px;
		min-width: 18px;
		height: 18px;
		padding: 0 4px;
		background: #ef4444;
		border-radius: 9px;
		font-size: 0.625rem;
		font-weight: 700;
		display: flex;
		align-items: center;
		justify-content: center;
	}

	.floating-btn.stack-btn {
		background: rgba(102, 126, 234, 0.15);
		border-color: rgba(102, 126, 234, 0.3);
	}

	.floating-btn.stack-btn:hover {
		background: rgba(102, 126, 234, 0.25);
		border-color: rgba(102, 126, 234, 0.5);
	}

	.floating-btn.stack-btn.has-items {
		background: rgba(251, 191, 36, 0.2);
		border-color: rgba(251, 191, 36, 0.4);
	}

	.floating-btn.stack-btn.has-items:hover {
		background: rgba(251, 191, 36, 0.3);
		border-color: rgba(251, 191, 36, 0.6);
	}

	.btn-primary {
		padding: 0.75rem 1.5rem;
		background: #667eea;
		color: white;
		border: none;
		border-radius: 6px;
		font-weight: 600;
		cursor: pointer;
		transition: background 0.2s;
	}

	.btn-primary:hover {
		background: #5568d3;
	}

	/* Prompt Overlay */
	.prompt-overlay {
		position: fixed;
		inset: 0;
		background: rgba(0, 0, 0, 0.85);
		display: flex;
		align-items: center;
		justify-content: center;
		z-index: 50;
	}

	.prompt-content {
		background: #1a1f2e;
		padding: 2rem;
		border-radius: 12px;
		border: 2px solid #667eea;
		max-width: 400px;
		text-align: center;
	}

	.prompt-message {
		font-size: 1.125rem;
		margin-bottom: 1.5rem;
	}

	.prompt-hint {
		color: #94a3b8;
		font-size: 0.875rem;
	}

	.prompt-buttons,
	.choice-buttons {
		display: flex;
		gap: 1rem;
		justify-content: center;
		flex-wrap: wrap;
	}

	.btn-yes {
		background: #22c55e;
		color: white;
		padding: 0.75rem 2rem;
		border: none;
		border-radius: 6px;
		font-weight: 600;
		cursor: pointer;
		transition: background 0.2s;
	}

	.btn-yes:hover {
		background: #16a34a;
	}

	.btn-no {
		background: #ef4444;
		color: white;
		padding: 0.75rem 2rem;
		border: none;
		border-radius: 6px;
		font-weight: 600;
		cursor: pointer;
		transition: background 0.2s;
	}

	.btn-no:hover {
		background: #dc2626;
	}

	.btn-choice {
		background: #374151;
		color: white;
		padding: 0.5rem 1rem;
		border: none;
		border-radius: 6px;
		cursor: pointer;
		transition: background 0.2s;
	}

	.btn-choice:hover {
		background: #4b5563;
	}

	/* Game Area Wrapper - contains game layout + sidebar */
	.game-area-wrapper {
		flex: 1;
		display: flex;
		overflow: hidden;
	}

	/* Main Game Layout - Full Width */
	.game-layout {
		flex: 1;
		display: flex;
		flex-direction: column;
		padding: 0.75rem;
		padding-bottom: 80px; /* Space for action bar */
		gap: 0.75rem;
		overflow: hidden;
		min-width: 0; /* Allow shrinking */
	}

	/* Visual Stack Sidebar Container */
	.visual-stack-sidebar-container {
		height: 100%;
		transition: all 0.2s ease-out;
	}

	.visual-stack-sidebar-container.drag-active {
		background: rgba(102, 126, 234, 0.05);
	}

	.visual-stack-sidebar-container.drag-valid {
		background: rgba(102, 126, 234, 0.15);
	}

	/* Opponents Row */
	.opponents-row {
		display: flex;
		gap: 0.75rem;
		flex-wrap: wrap;
		justify-content: center;
	}

	.opponents-row > :global(*) {
		flex: 1;
		min-width: 250px;
		max-width: 400px;
	}

	/* For 4 players, use different layout */
	.game-layout.four-player .opponents-row {
		justify-content: space-between;
	}

	.game-layout.four-player .opponents-row > :global(*) {
		flex: 1;
		min-width: 200px;
		max-width: 350px;
	}

	/* Battlefield Area */
	.battlefield-area {
		flex: 1;
		display: flex;
		flex-direction: column;
		gap: 0.75rem;
		background: linear-gradient(135deg, #0d1117, #141821);
		border: 1px solid #2a3441;
		border-radius: 12px;
		padding: 1rem;
		overflow: auto;
		min-height: 200px;
		position: relative;
		transition:
			border-color 0.2s,
			box-shadow 0.2s;
	}

	/* Battlefield drag state */
	.battlefield-area.drag-active {
		border-color: rgba(102, 126, 234, 0.5);
		box-shadow: inset 0 0 30px rgba(102, 126, 234, 0.1);
	}

	.battlefield-area.drag-valid {
		border-color: #22c55e;
		box-shadow:
			inset 0 0 40px rgba(34, 197, 94, 0.15),
			0 0 0 2px rgba(34, 197, 94, 0.3);
	}

	/* Drop zone overlay */
	.drop-zone-overlay {
		position: absolute;
		inset: 0;
		display: flex;
		align-items: center;
		justify-content: center;
		background: rgba(102, 126, 234, 0.1);
		border-radius: 12px;
		pointer-events: none;
		z-index: 10;
		transition: background 0.2s;
	}

	.drop-zone-overlay.valid {
		background: rgba(34, 197, 94, 0.15);
	}

	.drop-label {
		font-size: 1.125rem;
		font-weight: 600;
		color: #667eea;
		padding: 0.75rem 1.5rem;
		background: rgba(26, 31, 46, 0.9);
		border-radius: 8px;
		border: 2px dashed rgba(102, 126, 234, 0.5);
	}

	.drop-zone-overlay.valid .drop-label {
		color: #22c55e;
		border-color: rgba(34, 197, 94, 0.5);
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

	.zone-label {
		font-size: 0.6875rem;
		color: #6b7280;
		text-transform: uppercase;
		letter-spacing: 0.5px;
		font-weight: 600;
		margin-bottom: 0.5rem;
	}

	.command-zone {
		padding-bottom: 0.75rem;
		border-bottom: 1px solid #2a3441;
	}

	.command-cards,
	.battlefield-cards {
		display: flex;
		flex-wrap: wrap;
		gap: 0.5rem;
		align-content: flex-start;
	}

	.battlefield-rows {
		display: flex;
		flex-direction: column;
		gap: 0.5rem;
	}

	.battlefield-row--lands {
		margin-top: 0.25rem;
		padding-top: 0.5rem;
		border-top: 1px dashed rgba(148, 163, 184, 0.25);
	}

	.battlefield-card-wrapper {
		user-select: none;
		-webkit-user-select: none;
		transition:
			transform 0.2s ease,
			opacity 0.2s ease;
	}

	.battlefield-card-wrapper.draggable {
		cursor: grab;
	}

	.battlefield-card-wrapper.draggable:active {
		cursor: grabbing;
	}

	.battlefield-card-wrapper.is-dragging {
		opacity: 0.4;
		transform: scale(0.95);
	}

	.battlefield-card-wrapper.is-hovered {
		/* Subtle glow to indicate keyboard shortcuts are active */
		filter: drop-shadow(0 0 4px rgba(100, 200, 255, 0.5));
	}

	.my-battlefield {
		flex: 1;
	}

	/* Player Info Row */
	.player-info-row {
		display: flex;
		justify-content: space-between;
		align-items: center;
		gap: 0.75rem;
		padding: 0.375rem 0.75rem;
		background: rgba(26, 31, 46, 0.8);
		border-radius: 8px;
		border: 1px solid #2a3441;
	}

	/* Compact Player Info Row */
	.player-info-row.compact {
		padding: 0.25rem 0.5rem;
		gap: 0.5rem;
		min-height: 36px;
		background: rgba(26, 31, 46, 0.9);
	}

	.player-identity {
		display: flex;
		align-items: center;
		gap: 0.75rem;
	}

	.player-info-row.compact .player-identity {
		gap: 0.5rem;
	}

	.player-name {
		font-weight: 700;
		font-size: 0.9375rem;
		display: flex;
		align-items: center;
		gap: 0.375rem;
		white-space: nowrap;
	}

	.player-info-row.compact .player-name {
		font-size: 0.8125rem;
	}

	.player-name.has-priority {
		color: #22c55e;
	}

	.priority-dot {
		width: 8px;
		height: 8px;
		background: #22c55e;
		border-radius: 50%;
		animation: pulse 1.5s infinite;
	}

	.player-info-row.compact .priority-dot {
		width: 6px;
		height: 6px;
	}

	@keyframes pulse {
		0%,
		100% {
			opacity: 1;
			transform: scale(1);
		}
		50% {
			opacity: 0.6;
			transform: scale(1.2);
		}
	}

	.player-zones {
		display: flex;
		gap: 0.5rem;
		align-items: center;
	}

	.player-info-row.compact .player-zones {
		gap: 0.375rem;
	}

	/* Player Stats Inline (Life, Poison, Library) */
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

	.stat-btn.life-btn.minus:hover {
		background: rgba(239, 68, 68, 0.3);
		color: #ef4444;
	}

	.stat-btn.life-btn.plus:hover {
		background: rgba(34, 197, 94, 0.3);
		color: #22c55e;
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
		transition: background 0.15s ease;
	}

	.stat-display:hover {
		background: rgba(63, 63, 70, 0.3);
	}

	.stat-display.life {
		color: #f4f4f5;
	}

	.stat-display.poison {
		color: #a855f7;
	}

	.stat-icon {
		font-size: 0.75rem;
	}

	.stat-value {
		font-family: 'JetBrains Mono', monospace;
		min-width: 20px;
		text-align: center;
	}

	/* Quick Menu */
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
		transition: all 0.15s ease;
	}

	.menu-row button:hover {
		background: rgba(201, 162, 39, 0.2);
		border-color: rgba(201, 162, 39, 0.4);
		color: #f4f4f5;
	}

	.menu-row button.search-btn {
		flex: 1;
		background: rgba(139, 92, 246, 0.2);
		border-color: rgba(139, 92, 246, 0.4);
		color: #a78bfa;
	}

	.menu-row button.search-btn:hover {
		background: rgba(139, 92, 246, 0.3);
		border-color: rgba(139, 92, 246, 0.5);
		color: #c4b5fd;
	}

	.menu-value {
		min-width: 24px;
		text-align: center;
		font-weight: 600;
		color: #f4f4f5;
		font-family: 'JetBrains Mono', monospace;
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
		transition: all 0.15s ease;
	}

	.menu-close:hover {
		background: rgba(63, 63, 70, 0.5);
		color: #f4f4f5;
	}

	/* Ensure drop zone wrappers are positioned correctly for hit testing */
	.player-zones > div {
		position: relative;
	}

	/* Graveyard Drop Zone */
	.graveyard-drop-zone {
		transition: all 0.2s ease;
		border-radius: 6px;
		/* Ensure minimum drop target size for drag-drop detection */
		min-width: 70px;
		min-height: 32px;
		/* Ensure element receives pointer events for drop detection */
		position: relative;
	}

	.graveyard-drop-zone.drag-active {
		outline: 2px dashed #6b7280;
		outline-offset: 2px;
	}

	.graveyard-drop-zone.drag-valid {
		outline-color: #22c55e;
		background: rgba(34, 197, 94, 0.1);
	}

	/* Exile Drop Zone */
	.exile-drop-zone {
		transition: all 0.2s ease;
		border-radius: 6px;
		/* Ensure minimum drop target size for drag-drop detection */
		min-width: 70px;
		min-height: 32px;
	}

	.exile-drop-zone.drag-active {
		outline: 2px dashed #6b7280;
		outline-offset: 2px;
	}

	.exile-drop-zone.drag-valid {
		outline-color: #a78bfa;
		background: rgba(167, 139, 250, 0.1);
	}

	/* Hand Area */
	.hand-area {
		flex-shrink: 0;
		transition: all 0.2s ease;
		border-radius: 8px;
	}

	.hand-area.drag-active {
		outline: 2px dashed #6b7280;
		outline-offset: 2px;
	}

	.hand-area.drag-valid {
		outline-color: #22c55e;
		background: rgba(34, 197, 94, 0.1);
	}

	/* Responsive */
	@media (max-width: 900px) {
		.opponents-row > :global(*) {
			min-width: 200px;
		}

		.floating-actions {
			top: 80px;
			right: 8px;
		}

		.floating-btn {
			width: 38px;
			height: 38px;
			font-size: 1.125rem;
		}
	}

	@media (max-width: 600px) {
		.game-layout {
			padding: 0.5rem;
			padding-bottom: 70px;
		}

		.opponents-row > :global(*) {
			min-width: 150px;
			max-width: none;
		}

		.battlefield-area {
			padding: 0.75rem;
		}

		.player-info-row {
			flex-wrap: wrap;
			justify-content: center;
		}
	}

	/* Debug FAB Button */
	.debug-fab {
		position: fixed;
		bottom: 100px;
		right: 20px;
		width: 48px;
		height: 48px;
		border-radius: 50%;
		background: linear-gradient(135deg, #1a1a2e, #16213e);
		border: 2px solid #00ff00;
		color: #00ff00;
		font-size: 1.25rem;
		cursor: pointer;
		z-index: 500;
		box-shadow: 0 4px 20px rgba(0, 255, 0, 0.2);
		transition: all 0.2s;
	}

	.debug-fab:hover {
		transform: scale(1.1);
		box-shadow: 0 6px 30px rgba(0, 255, 0, 0.4);
	}

	/* Drag Ghost - Card following cursor */
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
		box-shadow:
			0 15px 40px rgba(0, 0, 0, 0.6),
			0 0 0 1px rgba(102, 126, 234, 0.3),
			0 0 30px rgba(102, 126, 234, 0.2);
		opacity: 0.95;
		transform: scale(1.1) rotate(-5deg);
		transition:
			border-color 0.15s,
			transform 0.15s,
			box-shadow 0.15s;
	}

	.drag-ghost-card.valid {
		border-color: #22c55e;
		box-shadow:
			0 15px 40px rgba(0, 0, 0, 0.6),
			0 0 30px rgba(34, 197, 94, 0.5),
			0 0 60px rgba(34, 197, 94, 0.2);
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
		overflow: hidden;
		text-overflow: ellipsis;
		display: -webkit-box;
		-webkit-line-clamp: 3;
		-webkit-box-orient: vertical;
		line-height: 1.3;
		padding: 0.25rem;
	}
</style>
