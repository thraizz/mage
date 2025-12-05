<script lang="ts">
	import { onMount, onDestroy } from 'svelte';
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
		sendPlayerUUID,
		sendPlayerBoolean,
		sendPlayerString,
		keepHand,
		mulligan,
		playLand,
		advancePhase,
		activateManaAbility,
		activateAbility
	} from '$lib/api/game';
	import { CardActionType, type CardView } from '$lib/generated/mage/v1/models';
	import type { GameCard, GamePhase } from '$lib/types/game';
	import {
		dragDropStore,
		isDragging as isDraggingStore,
		draggedCardId,
		draggedCardName,
		dragPosition,
		isOverValidDropZone,
		currentDropZone
	} from '$lib/utils/drag-drop';
	import { toast } from '$lib/stores/toast';
	import { getScryfallImageUrl } from '$lib/utils/scryfall';

	// Game components
	import Card from '$lib/components/game/Card.svelte';
	import PlayerHand from '$lib/components/game/PlayerHand.svelte';
	import Graveyard from '$lib/components/game/Graveyard.svelte';
	import ManaPool from '$lib/components/game/ManaPool.svelte';
	import GameHeader from '$lib/components/game/GameHeader.svelte';
	import Stack from '$lib/components/game/Stack.svelte';
	import MulliganDialog from '$lib/components/game/MulliganDialog.svelte';
	import TargetingMode from '$lib/components/game/TargetingMode.svelte';
	
	import PriorityActionBar from '$lib/components/game/PriorityActionBar.svelte';
	import ActionLogOverlay from '$lib/components/game/ActionLogOverlay.svelte';
	import GameChatOverlay from '$lib/components/game/GameChatOverlay.svelte';
	import OpponentPanel from '$lib/components/game/OpponentPanel.svelte';
	import DebugOverlay from '$lib/components/game/DebugOverlay.svelte';
	import ManaPayment from '$lib/components/game/ManaPayment.svelte';
	import XManaSelector from '$lib/components/game/XManaSelector.svelte';
	import AbilitiesPanel from '$lib/components/game/AbilitiesPanel.svelte';
	import type { GamePlayManaData, GamePlayXManaData, GameAssignDamageData } from '$lib/generated/mage/v1/websocket';
	
	// Combat components
	import DeclareAttackers from '$lib/components/game/DeclareAttackers.svelte';
	import DeclareBlockers from '$lib/components/game/DeclareBlockers.svelte';
	import AssignDamage from '$lib/components/game/AssignDamage.svelte';
	
	// Direct player control components
	import CardContextMenu from '$lib/components/game/CardContextMenu.svelte';
	import TokenCreator from '$lib/components/game/TokenCreator.svelte';
	import { combatStore, isInCombat, canAttackCardIds, declaredAttackerIds, canBlockCardIds, assignedBlockerIds } from '$lib/stores/combat';
	import { parseCombatOptions, type DeclaredAttacker, type DefenderTarget, type DamageAssignmentPrompt, type ParsedCombatOptions } from '$lib/types/combat';

	// Targeting store
	import {
		targetingStore,
		isTargetingActive,
		validTargetIds,
		selectedTargetIds,
		syncWithGamePrompt
	} from '$lib/stores/game-targeting';

	// Page data from load function
	const { data } = $props();
	
	// Game ID from load function (more reliable than accessing $page.params directly)
	const gameId = $derived(data.gameId);

	// UI state
	let showActionLog = $state(false);
	let showChat = $state(false);
	let showDebugOverlay = $state(false);
	let showTokenCreator = $state(false);
	
	// Context menu state
	let contextMenuCard = $state<typeof battlefieldCards[0] | null>(null);
	let contextMenuPosition = $state({ x: 0, y: 0 });
	let actionLogRef = $state<ActionLogOverlay | undefined>(undefined);
	let gameChatRef = $state<GameChatOverlay | undefined>(undefined);
	let isActionLoading = $state(false);
	let showStackOverlay = $state(false);
	let showAbilitiesPanel = $state(false);
	let abilitiesPanelCardId = $state<string | null>(null);
	let initialized = $state(false);

	// Targeting state (from store)
	const isTargeting = $derived($isTargetingActive);
	const validTargets = $derived($validTargetIds);
	const selectedTargets = $derived($selectedTargetIds);

	// Combat state (from store)
	const inCombat = $derived($isInCombat);
	const canAttackIds = $derived($canAttackCardIds);
	const attackingIds = $derived($declaredAttackerIds);
	const canBlockIds = $derived($canBlockCardIds);
	const blockingIds = $derived($assignedBlockerIds);

	// Combat phase detection from prompts
	const combatPromptOptions = $derived<ParsedCombatOptions | null>(() => {
		if (!prompt) return null;
		if (prompt.type !== 'choice') return null;
		const data = prompt.data as { choices?: string[] };
		if (!data?.choices) return null;
		const parsed = parseCombatOptions(data.choices);
		if (parsed.type === 'none') return null;
		return parsed;
	});

	const isDeclaringAttackersPhase = $derived(
		step === 'DECLARE_ATTACKERS' && combatPromptOptions()?.type === 'attack'
	);

	const isDeclaringBlockersPhase = $derived(
		step === 'DECLARE_BLOCKERS' && combatPromptOptions()?.type === 'block'
	);

	// Damage assignment from special prompt (GAME_ASSIGN_DAMAGE)
	const damageAssignmentPrompt = $derived<DamageAssignmentPrompt | null>(() => {
		if (!prompt) return null;
		// Check if this is a damage assignment prompt (sent as 'assignDamage' type)
		if ((prompt as {type: string}).type === 'assignDamage') {
			return prompt.data as DamageAssignmentPrompt;
		}
		return null;
	});

	// Drag-drop state (from store)
	const isDragging = $derived($isDraggingStore);
	// eslint-disable-next-line no-unused-vars
	const dragCardId = $derived($draggedCardId);
	const dragCardName = $derived($draggedCardName);
	const dragPos = $derived($dragPosition);
	const isOverValidDrop = $derived($isOverValidDropZone);
	const dropZone = $derived($currentDropZone);

	// Track targeting sync unsubscriber
	let targetingSyncUnsub: (() => void) | null = null;

	// Battlefield drop zone element reference
	let battlefieldDropZoneEl: HTMLDivElement | null = $state(null);
	let dropZoneUnregister: (() => void) | null = null;

	// Opponent panel states (for collapsing)
	let opponentExpanded = $state<Record<string, boolean>>({});

	// Auto-pass settings
	let autoPassSettings = $state({
		noActions: false,
		opponentTurn: false
	});

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
	const myCards = $derived($myHand);
	const myGrave = $derived($myGraveyard);
	const myMana = $derived($myManaPool);
	const havePriority = $derived($hasPriority);
	const phase = $derived($currentPhase);
	const step = $derived($currentStep);
	const turn = $derived($currentTurn);
	const battlefieldCards = $derived($battlefield);
	const stackCards = $derived($stack);
	const commandCards = $derived($command);
	const prompt = $derived($pendingPrompt);
	const isGameOver = $derived($gameOver);
	const gameWinner = $derived($winner);
	const error = $derived($gameError);
	const loading = $derived($isLoading);

	// Player name map for display
	const playerNames = $derived(
		new Map(allPlayers.map((p) => [p.playerId, p.name]))
	);

	// Mulligan phase detection - use server-provided value
	const isMulliganPhase = $derived(
		gameState.gameView?.isMulliganPhase ?? gameState.gameView?.state?.toLowerCase() === 'mulligan'
	);

	// Check if local player has already kept their hand (waiting for other players)
	const hasKeptHand = $derived(me?.keptHand ?? false);

	// Is it the local player's turn?
	const isYourTurn = $derived(
		gameState.gameView?.activePlayerId === localPlayerId
	);

	// Does the local player have any available actions? (server-computed)
	const myHasAvailableActions = $derived(me?.hasAvailableActions ?? false);

	// Get active player name - use server-provided value
	const activePlayerName = $derived(
		gameState.gameView?.activePlayerName || 'Unknown'
	);

	// Game format - use server-provided value
	const gameFormat = $derived(
		gameState.gameView?.gameFormat || 'Game'
	);

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
			'UNTAP': 'UNTAP',
			'UPKEEP': 'UPKEEP',
			'DRAW': 'DRAW',
			'DECLARE_ATTACKERS': 'DECLARE_ATTACKERS',
			'DECLARE_BLOCKERS': 'DECLARE_BLOCKERS',
			'COMBAT_DAMAGE': 'COMBAT_DAMAGE',
			'END': 'END',
			'CLEANUP': 'CLEANUP',
			
			// Server step names that need mapping to client keys
			'MAIN1': 'PRECOMBAT_MAIN',
			'BEGIN_COMBAT': 'COMBAT',
			'END_COMBAT': 'END_OF_COMBAT',
			'MAIN2': 'POSTCOMBAT_MAIN',
			
			// Client-only keys (for backwards compatibility)
			'BEGINNING': 'BEGINNING',
			'PRECOMBAT_MAIN': 'PRECOMBAT_MAIN',
			'COMBAT': 'COMBAT',
			'END_OF_COMBAT': 'END_OF_COMBAT',
			'POSTCOMBAT_MAIN': 'POSTCOMBAT_MAIN',
			'END_OF_TURN': 'END_OF_TURN',
			
			// Phase names (fallback when step not available)
			'ENDING': 'END'
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
				message: 'You are not a participant in this game. You may need to join as a spectator instead.'
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
				message: 'Unable to connect to the game server. Please check your internet connection and try again.'
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
			otherPlayers.forEach(p => {
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
		console.log('[handleCastSpell] myCards:', myCards.map(c => ({ id: c.id, name: c.name, type: c.type, actions: c.availableActions })));

		const card = myCards.find((c) => c.id === cardId);
		if (!card) {
			console.log('[handleCastSpell] Card not found in hand');
			addLogEntry('Selected card not found in hand');
			return;
		}

		console.log('[handleCastSpell] Found card:', { id: card.id, name: card.name, type: card.type, actions: card.availableActions });

		// Use server-provided availableActions to determine action type
		const playLandAction = card.availableActions?.find(a => a.actionType === CardActionType.CARD_ACTION_PLAY_LAND);
		const castSpellAction = card.availableActions?.find(a => a.actionType === CardActionType.CARD_ACTION_CAST_SPELL);

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
			const isLand = playLandAction !== undefined || (!castSpellAction && card.type.toLowerCase().includes('land'));
			console.log('[handleCastSpell] isLand:', isLand, 'playLandAction:', !!playLandAction, 'castSpellAction:', !!castSpellAction);

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
	 * Handle activate ability button
	 * Shows the abilities panel for the selected permanent
	 * For mana abilities, click directly on the land to tap it
	 */
	function handleActivateAbility() {
		console.log('[handleActivateAbility] Called', { havePriority, isActionLoading });
		if (!havePriority || isActionLoading) {
			console.log('[handleActivateAbility] Blocked - no priority or loading');
			return;
		}

		const gameState = $gameStore;
		const selectedIds = gameState.selectedCardIds;
		console.log('[handleActivateAbility] Selected IDs:', selectedIds);

		// Check if a permanent is selected
		if (selectedIds.length === 0) {
			addLogEntry('Select a permanent first, then press A to activate abilities');
			toast.info('Select a permanent first');
			return;
		}

		if (selectedIds.length > 1) {
			addLogEntry('Select only one permanent to activate abilities');
			toast.info('Select only one permanent');
			return;
		}

		const cardId = selectedIds[0];
		console.log('[handleActivateAbility] Looking for card:', cardId);
		console.log('[handleActivateAbility] battlefieldCards:', battlefieldCards.map(c => ({ id: c.id, name: c.name })));
		
		const card = battlefieldCards.find((c) => c.id === cardId);

		if (!card) {
			console.log('[handleActivateAbility] Card not found on battlefield');
			addLogEntry('Selected card not found on battlefield');
			toast.error('Card not found');
			return;
		}

		console.log('[handleActivateAbility] Found card:', { 
			id: card.id, 
			name: card.name, 
			availableActions: card.availableActions,
			availableActionsCount: card.availableActions?.length 
		});

		// Filter for non-mana activated abilities
		// Note: actionType can be either a number (enum) or string (JSON serialized)
		const activatedAbilities = card.availableActions?.filter(
			(a) => {
				const isMatch = a.actionType === CardActionType.CARD_ACTION_ACTIVATE_ABILITY || 
				                String(a.actionType) === 'CARD_ACTION_ACTIVATE_ABILITY' ||
				                a.actionType === 3; // Enum value for CARD_ACTION_ACTIVATE_ABILITY
				console.log('[handleActivateAbility] Checking action:', { 
					actionType: a.actionType, 
					actionTypeString: String(a.actionType),
					expected: CardActionType.CARD_ACTION_ACTIVATE_ABILITY,
					isMatch 
				});
				return isMatch;
			}
		) || [];

		console.log('[handleActivateAbility] activatedAbilities:', activatedAbilities);

		if (activatedAbilities.length === 0) {
			// Check if it's a land with mana ability - give hint
			const hasManaAbility = card.availableActions?.some(
				(a) => a.actionType === CardActionType.CARD_ACTION_ACTIVATE_MANA_ABILITY ||
				       String(a.actionType) === 'CARD_ACTION_ACTIVATE_MANA_ABILITY' ||
				       a.actionType === 4
			);
			if (hasManaAbility) {
				addLogEntry('Click on the land to tap it for mana');
				toast.info('Click on the land to tap it for mana');
			} else {
				addLogEntry(`${card.name} has no activated abilities`);
				toast.info('No activated abilities');
			}
			return;
		}

		// Show the abilities panel
		console.log('[handleActivateAbility] Showing abilities panel for card:', cardId);
		abilitiesPanelCardId = cardId;
		showAbilitiesPanel = true;
		console.log('[handleActivateAbility] Panel state set:', { showAbilitiesPanel, abilitiesPanelCardId });
	}

	/**
	 * Handle ability activation from the panel
	 */
	async function handleAbilityActivate(abilityId: string) {
		console.log('[handleAbilityActivate] Called with abilityId:', abilityId);
		if (!gameId || !abilitiesPanelCardId) {
			console.log('[handleAbilityActivate] Missing gameId or cardId:', { gameId, abilitiesPanelCardId });
			return;
		}

		const card = battlefieldCards.find((c) => c.id === abilitiesPanelCardId);
		const cardName = card?.name || 'permanent';

		console.log('[handleAbilityActivate] Activating ability:', {
			gameId,
			cardId: abilitiesPanelCardId,
			abilityId,
			cardName
		});

		showAbilitiesPanel = false;
		isActionLoading = true;

		try {
			console.log('[handleAbilityActivate] Calling activateAbility API...');
			await activateAbility(gameId, abilitiesPanelCardId, abilityId);
			console.log('[handleAbilityActivate] API call successful');
			addLogEntry(`Activated ability of ${cardName}`);
			toast.success(`Activated ability of ${cardName}`);
		} catch (err) {
			const message = err instanceof Error ? err.message : 'Failed to activate ability';
			console.error('[handleAbilityActivate] API call failed:', err);
			toast.error(message);
			addLogEntry(`Failed to activate ability: ${message}`);
		} finally {
			isActionLoading = false;
			abilitiesPanelCardId = null;
		}
	}

	/**
	 * Close abilities panel
	 */
	function closeAbilitiesPanel() {
		showAbilitiesPanel = false;
		abilitiesPanelCardId = null;
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
		playCardOptimistic(cardId);
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
	 * Handle card click (for selection/targeting)
	 * Note: PlayerHand already handles selection toggle via gameStore.toggleCardSelection
	 * This handler is for additional logic (target prompts, logging)
	 */
	function handleCardClick(cardId: string) {
		// Handle targeting mode - toggle target selection
		if (isTargeting) {
			const toggled = targetingStore.toggleTarget(cardId);
			if (toggled) {
				const card = myCards.find((c) => c.id === cardId);
				if (card) {
					addLogEntry(`Target: ${card.name}`);
				}
			}
			return;
		}

		const card = myCards.find((c) => c.id === cardId);
		if (card) {
			addLogEntry(`Selected: ${card.name}`);
		}
	}

	/**
	 * Handle battlefield card click
	 */
	async function handleBattlefieldCardClick(cardId: string) {
		// Handle targeting mode - toggle target selection
		if (isTargeting) {
			const toggled = targetingStore.toggleTarget(cardId);
			if (toggled) {
				const card = battlefieldCards.find((c) => c.id === cardId);
				if (card) {
					addLogEntry(`Target: ${card.name}`);
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
				(a) => (a.actionType === CardActionType.CARD_ACTION_ACTIVATE_MANA_ABILITY || 
				        String(a.actionType) === 'CARD_ACTION_ACTIVATE_MANA_ABILITY') && a.isEnabled
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

		// Default behavior: toggle selection
		gameStore.toggleCardSelection(cardId);
	}

	/**
	 * Handle right-click context menu on battlefield cards
	 */
	function handleCardContextMenu(event: MouseEvent, card: typeof battlefieldCards[0]) {
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
	 * Handle target selection confirmation
	 * Sends selected target UUID(s) to the server
	 * 
	 * Note: The server typically handles multi-target by:
	 * 1. Sending separate GAME_TARGET events for each target slot
	 * 2. Or accepting targets sequentially via SendPlayerUUID
	 * 
	 * Current implementation sends targets one at a time.
	 * If server expects all targets at once, use SendPlayerUUIDs API instead.
	 */
	async function handleTargetConfirm(targetIds: string[]) {
		if (!gameId || targetIds.length === 0) return;

		isActionLoading = true;
		try {
			// Send each target to the server
			// For most MTG implementations, the server expects one target per call
			// and will send another GAME_TARGET if more targets are needed
			for (const targetId of targetIds) {
				await sendPlayerUUID(gameId, targetId);
			}
			gameStore.clearPrompt();
			addLogEntry(`Confirmed ${targetIds.length} target(s)`);
		} catch (err) {
			console.error('Failed to confirm target:', err);
		} finally {
			isActionLoading = false;
		}
	}

	/**
	 * Handle target selection cancellation
	 */
	async function handleTargetCancel() {
		if (!gameId) return;

		isActionLoading = true;
		try {
			// Send empty UUID to cancel targeting
			await sendPlayerUUID(gameId, '');
			gameStore.clearPrompt();
			addLogEntry('Cancelled target selection');
		} catch (err) {
			console.error('Failed to cancel targeting:', err);
		} finally {
			isActionLoading = false;
		}
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
			if (card.type?.toLowerCase().includes('planeswalker') && 
			    card.controllerId !== localPlayerId) {
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
					const card = battlefieldCards.find(c => c.id === attackerId);
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
	 * Toggle stack overlay
	 */
	function toggleStack() {
		showStackOverlay = !showStackOverlay;
	}

	/**
	 * Get cards controlled by a specific player on battlefield
	 */
	function getPlayerBattlefieldCards(playerId: string): CardView[] {
		return battlefieldCards.filter((c) => c.controllerId === playerId);
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
	 * Auto-pass effect - triggers when auto-pass conditions are met
	 * Conditions checked:
	 * 1. Player has priority
	 * 2. Not in mulligan phase
	 * 3. Not already loading an action
	 * 4. Game is initialized
	 * 5. Either:
	 *    - "Auto-pass on opponent's turn" is enabled AND it's not your turn
	 *    - "Auto-pass when no actions" is enabled AND server says no available actions
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

		// Check auto-pass conditions
		let shouldAutoPass = false;
		let reason = '';

		// Condition 1: Auto-pass on opponent's turn
		if (autoPassSettings.opponentTurn && !isYourTurn) {
			shouldAutoPass = true;
			reason = "opponent's turn";
		}

		// Condition 2: Auto-pass when no available actions (server-computed)
		if (autoPassSettings.noActions && !myHasAvailableActions) {
			shouldAutoPass = true;
			reason = 'no available actions';
		}

		if (shouldAutoPass) {
			console.log(`[AutoPass] Triggering auto-pass: ${reason}`);
			isAutoPassPending = true;

			// Capture gameId in local scope for the async callback
			const currentGameId = gameId;

			// Use setTimeout to avoid synchronous state updates in effect
			setTimeout(async () => {
				try {
					if (currentGameId) {
						await passPriority(currentGameId);
						addLogEntry(`Auto-passed (${reason})`);
					}
				} catch (err) {
					console.error('[AutoPass] Failed to auto-pass:', err);
				} finally {
					isAutoPassPending = false;
				}
			}, 50); // Small delay for smoother UX
		}
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
					// Accept cards from hand when player has priority
					return sourceZone === 'hand' && havePriority;
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

		// Sync targeting store with game prompts
		targetingSyncUnsub = syncWithGamePrompt();

		// Initialize game - gameId should always be available from load function
		initializeGame();
	});

	// Cleanup on destroy
	onDestroy(() => {
		// Cleanup targeting sync
		if (targetingSyncUnsub) {
			targetingSyncUnsub();
			targetingSyncUnsub = null;
		}
		// Reset targeting state
		targetingStore.exitTargetingMode();
		gameStore.reset();
	});
</script>

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
				{:else if parsedError.title.toLowerCase().includes('session') || parsedError.title.toLowerCase().includes('expired')}
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
					<button class="btn-secondary" onclick={() => { isInitializing = false; initialized = false; initializeGame(); }}>
						Try Again
					</button>
				{/if}
			</div>
		</div>
	{:else if isGameOver}
		<div class="game-over-overlay">
			<div class="game-over-content">
				<h2>Game Over</h2>
				<p class="winner-text">{gameWinner ? `Winner: ${playerNames.get(gameWinner) || gameWinner}` : 'Draw'}</p>
				<button class="btn-primary" onclick={() => goto('/lobby')}>Return to Lobby</button>
			</div>
		</div>
	{:else if isMulliganPhase}
		<MulliganDialog
			cards={myCards}
			mulliganCount={mulliganCount}
			onKeep={handleKeepHand}
			onMulligan={handleMulligan}
			isLoading={isMulliganLoading}
			hasKeptHand={hasKeptHand}
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
			onLogClick={() => showActionLog = true}
			onConcedeClick={handleConcede}
		/>

		<!-- Floating action buttons -->
		<div class="floating-actions">
			{#if stackCards.length > 0}
				<button class="floating-btn stack-btn" onclick={toggleStack} title="View Stack">
					📚 <span class="badge">{stackCards.length}</span>
				</button>
			{/if}
			<button class="floating-btn" onclick={() => showChat = true} title="Game Chat">
				💬
			</button>
		</div>

		<!-- Target Selection Mode Overlay -->
		<TargetingMode
			onConfirm={handleTargetConfirm}
			onCancel={handleTargetCancel}
		/>

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

		<!-- Combat: Declare Attackers -->
		{#if isDeclaringAttackersPhase && combatPromptOptions() && gameId}
			<DeclareAttackers
				{gameId}
				options={combatPromptOptions()!}
				battlefieldCards={getPlayerBattlefieldCards(localPlayerId)}
				defenders={getDefenderTargets()}
				onComplete={handleCombatComplete}
			/>
		{/if}

		<!-- Combat: Declare Blockers -->
		{#if isDeclaringBlockersPhase && combatPromptOptions() && gameId}
			<DeclareBlockers
				{gameId}
				options={combatPromptOptions()!}
				{battlefieldCards}
				attackingCreatures={getAttackingCreatures()}
				onComplete={handleCombatComplete}
			/>
		{/if}

		<!-- Combat: Assign Damage -->
		{#if damageAssignmentPrompt() && gameId}
			<AssignDamage
				{gameId}
				prompt={damageAssignmentPrompt()!}
				onComplete={handleCombatComplete}
			/>
		{/if}

		<!-- Abilities Panel -->
		{@const _debugPanelRender = (() => { 
			if (showAbilitiesPanel || abilitiesPanelCardId) {
				console.log('[RENDER] AbilitiesPanel condition check:', { 
					showAbilitiesPanel, 
					abilitiesPanelCardId,
					battlefieldCardsCount: battlefieldCards.length,
					cardFound: battlefieldCards.find(c => c.id === abilitiesPanelCardId)?.name
				});
			}
			return null;
		})()}
		{#if showAbilitiesPanel && abilitiesPanelCardId}
			{@const selectedCard = battlefieldCards.find(c => c.id === abilitiesPanelCardId)}
			{@const _debugCardFound = console.log('[RENDER] Inside panel block, selectedCard:', selectedCard?.name)}
			{#if selectedCard}
				<AbilitiesPanel
					cardId={abilitiesPanelCardId}
					cardName={selectedCard.name}
					abilities={selectedCard.availableActions?.filter(
						a => a.actionType === CardActionType.CARD_ACTION_ACTIVATE_ABILITY ||
						     String(a.actionType) === 'CARD_ACTION_ACTIVATE_ABILITY' ||
						     a.actionType === 3
					) || []}
					onActivate={handleAbilityActivate}
					onClose={closeAbilitiesPanel}
				/>
			{/if}
		{/if}

		<!-- Prompt Overlay (non-target, non-mana prompts) -->
		{#if prompt && !['target', 'mana', 'xmana'].includes(prompt.type)}
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

		<!-- Main Game Layout - Full Width -->
		<main class="game-layout" class:four-player={otherPlayers.length >= 3}>
			<!-- Opponents Row -->
			<div class="opponents-row">
				{#each otherPlayers as opponent, idx (opponent.playerId)}
					{@const position = getOpponentPosition(idx, otherPlayers.length)}
					<OpponentPanel
						{opponent}
						battlefieldCards={getPlayerBattlefieldCards(opponent.playerId)}
						selectedCardIds={gameState.selectedCardIds}
						bind:expanded={opponentExpanded[opponent.playerId]}
						{position}
						onCardClick={handleBattlefieldCardClick}
					/>
				{/each}
			</div>

			<!-- Central Battlefield Area (Drop Zone) -->
			<div
				bind:this={battlefieldDropZoneEl}
				class="battlefield-area"
				class:drag-active={isDragging}
				class:drag-valid={isDragging && isOverValidDrop && dropZone === 'battlefield'}
			>
				<!-- Command Zone -->
				{#if commandCards.length > 0}
					<div class="command-zone">
						<span class="zone-label">Command Zone</span>
						<div class="command-cards">
							{#each commandCards as card (card.id)}
								<Card
									cardId={card.id}
									cardName={card.name}
									manaCost={card.manaCost}
									cardType={card.type}
									power={card.power}
									toughness={card.toughness}
									imageUrl=""
									isTapped={card.tapped}
									isSelected={gameState.selectedCardIds.includes(card.id)}
									size="small"
									onclick={() => handleBattlefieldCardClick(card.id)}
									isTargetingActive={isTargeting}
									isValidTarget={validTargets.has(card.id)}
									isTargetSelected={selectedTargets.includes(card.id)}
									hasActivatedAbilities={card.availableActions?.some(a => a.actionType === CardActionType.CARD_ACTION_ACTIVATE_ABILITY || String(a.actionType) === 'CARD_ACTION_ACTIVATE_ABILITY' || a.actionType === 3)}
									summoningSickness={card.summoningSickness}
								/>
							{/each}
						</div>
					</div>
				{/if}

				<!-- My Battlefield -->
				<div class="my-battlefield">
					<span class="zone-label">Your Battlefield</span>
					<div class="battlefield-cards">
						{#each getPlayerBattlefieldCards(localPlayerId) as card (card.id)}
							<Card
								cardId={card.id}
								cardName={card.name}
								manaCost={card.manaCost}
								cardType={card.type}
								power={card.power}
								toughness={card.toughness}
								imageUrl=""
								isTapped={card.tapped}
								isSelected={gameState.selectedCardIds.includes(card.id)}
								size="normal"
								onclick={() => handleBattlefieldCardClick(card.id)}
								oncontextmenu={(e) => handleCardContextMenu(e, card)}
								isTargetingActive={isTargeting}
								isValidTarget={validTargets.has(card.id)}
								isTargetSelected={selectedTargets.includes(card.id)}
								hasActivatedAbilities={card.availableActions?.some(a => a.actionType === CardActionType.CARD_ACTION_ACTIVATE_ABILITY || String(a.actionType) === 'CARD_ACTION_ACTIVATE_ABILITY' || a.actionType === 3)}
								canAttack={canAttackIds.has(card.id)}
								isAttacking={attackingIds.has(card.id)}
								canBlock={canBlockIds.has(card.id)}
								isBlocking={blockingIds.has(card.id)}
								summoningSickness={card.summoningSickness}
							/>
						{/each}
						{#if getPlayerBattlefieldCards(localPlayerId).length === 0}
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

				<!-- Drop zone overlay indicator -->
				{#if isDragging}
					<div class="drop-zone-overlay" class:valid={isOverValidDrop && dropZone === 'battlefield'}>
						<span class="drop-label">
							{#if isOverValidDrop && dropZone === 'battlefield'}
								✓ Release to play
							{:else}
								Drag card here
							{/if}
						</span>
					</div>
				{/if}
			</div>

			<!-- Player Info & Zones Row -->
			<div class="player-info-row">
				{#if me}
					<div class="player-identity">
						<span class="player-name" class:has-priority={havePriority}>
							You
							{#if havePriority}
								<span class="priority-dot"></span>
							{/if}
						</span>
					</div>
					<div class="player-zones">
						<Graveyard
							cards={myGrave.map(toGameCard)}
							playerName="You"
							isOpponent={false}
							onCardClick={handleCardClick}
						/>
						<ManaPool
							mana={myMana}
							showEmpty={false}
							size="normal"
							onManaClick={() => {}}
						/>
					</div>
				{/if}
			</div>

			<!-- Player Hand -->
			<div class="hand-area">
				<PlayerHand
					onCardClick={handleCardClick}
					size="normal"
					currentPhase={step || phase}
					canDrag={havePriority && !isTargeting}
				/>
			</div>
		</main>

		<!-- Stack Overlay -->
		{#if showStackOverlay && stackCards.length > 0}
			<div class="stack-overlay" onclick={toggleStack} role="button" tabindex="0" onkeydown={(e) => e.key === 'Escape' && toggleStack()}>
				<div class="stack-panel" onclick={(e) => e.stopPropagation()} role="presentation">
					<Stack
						stackObjects={stackCards.map((c) => ({
							id: c.id,
							type: 'SPELL' as const,
							name: c.name,
							controllerId: c.controllerId,
							sourceCardId: c.id
						}))}
						{playerNames}
						onStackObjectClick={(stackId) => {
							gameStore.toggleCardSelection(stackId);
						}}
					/>
				</div>
			</div>
		{/if}

		<!-- Overlay Panels -->
		<ActionLogOverlay bind:this={actionLogRef} bind:open={showActionLog} />
		<GameChatOverlay bind:this={gameChatRef} gameId={gameId || ''} bind:open={showChat} />


		<!-- Priority Action Bar (Docked at bottom) -->
		<PriorityActionBar
			{gameId}
			hasPriority={havePriority}
			activePlayerId={gameState.gameView?.activePlayerId || ''}
			{localPlayerId}
			activePlayerName={activePlayerName}
			currentPhase={phase}
			canPassPriority={havePriority}
			isLoading={isActionLoading}
			playerLife={me?.life ?? 20}
			playerPoison={me?.poison ?? 0}
			libraryCount={me?.libraryCount ?? 0}
			onPassPriority={handlePassPriority}
			onPassUntilEOT={handlePassUntilEOT}
			onCastSpell={handleCastSpell}
			onActivateAbility={handleActivateAbility}
			onAdvancePhase={handleAdvancePhase}
			onCreateToken={() => showTokenCreator = true}
			bind:autoPassSettings
		/>

		<!-- Floating Debug Button -->
		<button 
			class="debug-fab" 
			onclick={() => showDebugOverlay = true}
			title="Open Debug View"
		>
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
			onClose={() => showDebugOverlay = false}
		/>

		<!-- Card Context Menu -->
		{#if contextMenuCard && gameId}
			<CardContextMenu
				{gameId}
				cardId={contextMenuCard.id}
				cardName={contextMenuCard.name}
				isTapped={contextMenuCard.tapped}
				isFaceDown={contextMenuCard.faceDown ?? false}
				isToken={contextMenuCard.isToken ?? false}
				currentZone="BATTLEFIELD"
				position={contextMenuPosition}
				onClose={closeContextMenu}
			/>
		{/if}

		<!-- Token Creator Dialog -->
		{#if showTokenCreator && gameId}
			<TokenCreator
				{gameId}
				onClose={() => showTokenCreator = false}
			/>
		{/if}

		<!-- Drag Ghost - Card following the cursor during drag -->
		{#if isDragging && dragCardName}
			{@const dragImageUrl = getScryfallImageUrl(dragCardName, 'small')}
			<div
				class="drag-ghost"
				style="left: {dragPos.x}px; top: {dragPos.y}px;"
			>
				<div class="drag-ghost-card" class:valid={isOverValidDrop}>
					{#if dragImageUrl}
						<img
							src={dragImageUrl}
							alt={dragCardName}
							class="drag-ghost-image"
							draggable="false"
						/>
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
		0%, 100% { border-color: rgba(251, 191, 36, 0.25); }
		50% { border-color: rgba(251, 191, 36, 0.5); }
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
		to { transform: rotate(360deg); }
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
		backdrop-filter: blur(8px);
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
		background: rgba(251, 191, 36, 0.15);
		border-color: rgba(251, 191, 36, 0.3);
	}

	.floating-btn.stack-btn:hover {
		background: rgba(251, 191, 36, 0.25);
		border-color: rgba(251, 191, 36, 0.5);
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

	.btn-primary:hover { background: #5568d3; }

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

	.btn-yes:hover { background: #16a34a; }

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

	.btn-no:hover { background: #dc2626; }

	.btn-choice {
		background: #374151;
		color: white;
		padding: 0.5rem 1rem;
		border: none;
		border-radius: 6px;
		cursor: pointer;
		transition: background 0.2s;
	}

	.btn-choice:hover { background: #4b5563; }

	/* Main Game Layout - Full Width */
	.game-layout {
		flex: 1;
		display: flex;
		flex-direction: column;
		padding: 0.75rem;
		padding-bottom: 80px; /* Space for action bar */
		gap: 0.75rem;
		overflow: hidden;
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
		transition: border-color 0.2s, box-shadow 0.2s;
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
		0%, 100% { opacity: 0.7; }
		50% { opacity: 1; }
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

	.my-battlefield {
		flex: 1;
	}

	.empty-battlefield {
		color: #4b5563;
		font-style: italic;
		font-size: 0.875rem;
		padding: 2rem;
		text-align: center;
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

	.player-identity {
		display: flex;
		align-items: center;
		gap: 0.75rem;
	}

	.player-name {
		font-weight: 700;
		font-size: 0.9375rem;
		display: flex;
		align-items: center;
		gap: 0.375rem;
		white-space: nowrap;
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

	@keyframes pulse {
		0%, 100% { opacity: 1; transform: scale(1); }
		50% { opacity: 0.6; transform: scale(1.2); }
	}

	.player-stats {
		display: flex;
		gap: 0.75rem;
		font-size: 0.8125rem;
	}

	.player-stats .life { color: #ef4444; font-weight: 700; }
	.player-stats .poison { color: #a855f7; }
	.player-stats .library { color: #3b82f6; }

	.player-zones {
		display: flex;
		gap: 0.5rem;
		align-items: center;
	}

	/* Hand Area */
	.hand-area {
		flex-shrink: 0;
	}

	/* Stack Overlay */
	.stack-overlay {
		position: fixed;
		inset: 0;
		background: rgba(0, 0, 0, 0.8);
		display: flex;
		align-items: center;
		justify-content: center;
		z-index: 100;
	}

	.stack-panel {
		width: 90%;
		max-width: 500px;
		max-height: 80vh;
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
		transition: border-color 0.15s, transform 0.15s, box-shadow 0.15s;
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
