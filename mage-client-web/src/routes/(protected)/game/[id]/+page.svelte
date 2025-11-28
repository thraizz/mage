<script lang="ts">
	import { onMount, onDestroy } from 'svelte';
	import { page } from '$app/stores';
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
		isLoading
	} from '$lib/stores/game';
	import {
		joinGame,
		fetchGameView,
		passPriority,
		passUntilMyNextTurn,
		concedeGame,
		sendPlayerUUID,
		sendPlayerBoolean,
		sendPlayerString,
		keepHand,
		mulligan,
		playLand,
		advancePhase
	} from '$lib/api/game';
	import { CardActionType, type CardView } from '$lib/generated/mage/v1/models';
	import type { GameCard, GamePhase } from '$lib/types/game';

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

	// Targeting store
	import {
		targetingStore,
		isTargetingActive,
		validTargetIds,
		selectedTargetIds,
		syncWithGamePrompt
	} from '$lib/stores/game-targeting';

	// Game ID from route params
	const gameId = $derived($page.params.id);

	// UI state
	let showActionLog = $state(false);
	let showChat = $state(false);
	let showDebugOverlay = $state(false);
	let actionLogRef = $state<ActionLogOverlay | undefined>(undefined);
	let gameChatRef = $state<GameChatOverlay | undefined>(undefined);
	let isActionLoading = $state(false);
	let showStackOverlay = $state(false);
	let initialized = $state(false);

	// Targeting state (from store)
	const isTargeting = $derived($isTargetingActive);
	const validTargets = $derived($validTargetIds);
	const selectedTargets = $derived($selectedTargetIds);

	// Track targeting sync unsubscriber
	let targetingSyncUnsub: (() => void) | null = null;

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

			const wsState = $websocketStore;
			if (wsState.state !== 'connected') {
				const token = $auth.token;
				const sessionId = token ? getSessionIdFromToken(token) : null;
				if (sessionId) {
					console.log('[GamePage] Connecting to WebSocket...');
					await websocketStore.connect(sessionId);
					console.log('[GamePage] WebSocket connected');
				} else {
					throw new Error('No session ID available');
				}
			}

			console.log('[GamePage] Initializing game store...');
			gameStore.initGame(gameId, playerId);

			console.log('[GamePage] Joining game...');
			await joinGame(gameId);
			console.log('[GamePage] Joined game successfully');

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
			gameStore.setError(err instanceof Error ? err.message : 'Failed to load game');
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
	 * Handle pass until player's next turn (F6)
	 * This passes priority automatically until your next upkeep step
	 */
	async function handlePassUntilEOT() {
		if (!havePriority || isActionLoading || !gameId) return;

		isActionLoading = true;
		try {
			await passUntilMyNextTurn(gameId);
			addLogEntry('You will pass until your next turn');
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
	 * Handle activate ability
	 */
	function handleActivateAbility() {
		if (!havePriority || isActionLoading) return;
		// TODO: Implement ability activation UI
		addLogEntry('Activate ability: select a permanent first');
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
	function handleBattlefieldCardClick(cardId: string) {
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

		gameStore.toggleCardSelection(cardId);
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
	<title>Game {gameId} - MAGE</title>
</svelte:head>

<div class="game-container" class:has-priority={havePriority}>
	{#if loading && !initialized}
		<div class="loading-overlay">
			<div class="spinner"></div>
			<p>Loading game...</p>
		</div>
	{:else if error && !gameState.gameView}
		<div class="error-overlay">
			<p class="error-message">{error}</p>
			<button class="btn-primary" onclick={() => goto('/lobby')}>Return to Lobby</button>
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

		<!-- Prompt Overlay (non-target prompts) -->
		{#if prompt && prompt.type !== 'target'}
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

			<!-- Central Battlefield Area -->
			<div class="battlefield-area">
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
								isTargetingActive={isTargeting}
								isValidTarget={validTargets.has(card.id)}
								isTargetSelected={selectedTargets.includes(card.id)}
							/>
						{/each}
						{#if getPlayerBattlefieldCards(localPlayerId).length === 0}
							<div class="empty-battlefield">No permanents</div>
						{/if}
					</div>
				</div>
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
						<div class="player-stats">
							<span class="life" title="Life">❤️ {formatLife(me.life)}</span>
							{#if me.poison > 0}
								<span class="poison" title="Poison">☠️ {me.poison}</span>
							{/if}
							<span class="library" title="Library">📚 {me.libraryCount}</span>
						</div>
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
				<PlayerHand onCardClick={handleCardClick} size="normal" />
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
			hasPriority={havePriority}
			activePlayerId={gameState.gameView?.activePlayerId || ''}
			{localPlayerId}
			activePlayerName={activePlayerName}
			canPassPriority={havePriority}
			isLoading={isActionLoading}
			onPassPriority={handlePassPriority}
			onPassUntilEOT={handlePassUntilEOT}
			onCastSpell={handleCastSpell}
			onActivateAbility={handleActivateAbility}
			onAdvancePhase={handleAdvancePhase}
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
			{gameId}
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

	.error-message { color: #ef4444; }

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
		gap: 1rem;
		padding: 0.5rem 1rem;
		background: #1a1f2e;
		border-radius: 8px;
		border: 1px solid #2a3441;
	}

	.player-identity {
		display: flex;
		flex-direction: column;
		gap: 0.25rem;
	}

	.player-name {
		font-weight: 700;
		font-size: 1rem;
		display: flex;
		align-items: center;
		gap: 0.5rem;
	}

	.player-name.has-priority {
		color: #22c55e;
	}

	.priority-dot {
		width: 10px;
		height: 10px;
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
		gap: 1rem;
		font-size: 0.875rem;
	}

	.player-stats .life { color: #ef4444; font-weight: 700; }
	.player-stats .poison { color: #a855f7; }
	.player-stats .library { color: #3b82f6; }

	.player-zones {
		display: flex;
		gap: 1rem;
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
</style>
