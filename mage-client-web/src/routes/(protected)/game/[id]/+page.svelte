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
		passUntilEndOfTurn,
		concedeGame,
		sendPlayerUUID,
		sendPlayerBoolean,
		sendPlayerString,
		keepHand,
		mulligan
	} from '$lib/api/game';
	import type { CardView } from '$lib/generated/mage/v1/models';
	import type { GameCard, GamePhase } from '$lib/types/game';

	// Game components
	import Card from '$lib/components/game/Card.svelte';
	import PlayerHand from '$lib/components/game/PlayerHand.svelte';
	import Graveyard from '$lib/components/game/Graveyard.svelte';
	import ManaPool from '$lib/components/game/ManaPool.svelte';
	import PhaseIndicator from '$lib/components/game/PhaseIndicator.svelte';
	import Stack from '$lib/components/game/Stack.svelte';
	import PriorityIndicator from '$lib/components/game/PriorityIndicator.svelte';
	import GameActionsPanel from '$lib/components/game/GameActionsPanel.svelte';
	import GameChat from '$lib/components/game/GameChat.svelte';
	import ActionLog from '$lib/components/game/ActionLog.svelte';
	import MulliganDialog from '$lib/components/game/MulliganDialog.svelte';

	// Game ID from route params
	const gameId = $derived($page.params.id);

	// UI state
	let chatCollapsed = $state(false);
	let actionLogCollapsed = $state(false);
	let gameChatRef: GameChat | undefined;
	let actionLogRef: ActionLog | undefined;
	let isActionLoading = $state(false);
	let showStackOverlay = $state(false);
	let initialized = $state(false);

	// Mulligan state
	let mulliganCount = $state(0);
	let isMulliganLoading = $state(false);

	// Get local player ID from auth (server uses usernames as player IDs)
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

	// Mulligan phase detection
	const isMulliganPhase = $derived(
		gameState.gameView?.state?.toLowerCase() === 'mulligan'
	);

	// Get active player name
	const activePlayerName = $derived(() => {
		const gv = gameState.gameView;
		if (!gv) return '';
		const active = allPlayers.find((p) => p.playerId === gv.activePlayerId);
		return active?.name || 'Unknown';
	});

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
	 * Convert phase string to GamePhase type
	 */
	function toGamePhase(phase: string): GamePhase {
		const phases: Record<string, GamePhase> = {
			'BEGINNING': 'BEGINNING',
			'UNTAP': 'UNTAP',
			'UPKEEP': 'UPKEEP',
			'DRAW': 'DRAW',
			'PRECOMBAT_MAIN': 'PRECOMBAT_MAIN',
			'COMBAT': 'COMBAT',
			'DECLARE_ATTACKERS': 'DECLARE_ATTACKERS',
			'DECLARE_BLOCKERS': 'DECLARE_BLOCKERS',
			'COMBAT_DAMAGE': 'COMBAT_DAMAGE',
			'END_OF_COMBAT': 'END_OF_COMBAT',
			'POSTCOMBAT_MAIN': 'POSTCOMBAT_MAIN',
			'END': 'END',
			'END_OF_TURN': 'END_OF_TURN',
			'CLEANUP': 'CLEANUP'
		};
		return phases[phase] || 'PRECOMBAT_MAIN';
	}

	/**
	 * Initialize game connection
	 */
	async function initializeGame() {
		if (!localPlayerId || !gameId) {
			console.error('Missing player ID or game ID');
			return;
		}

		try {
			console.log('[GamePage] Starting game initialization...', { gameId, localPlayerId });

			// Step 1: Connect to WebSocket FIRST to ensure we receive events
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
			} else {
				console.log('[GamePage] WebSocket already connected');
			}

			// Step 2: Initialize game store and subscribe to events AFTER WebSocket is connected
			console.log('[GamePage] Initializing game store and subscribing to events...');
			gameStore.initGame(gameId, localPlayerId);

			// Step 3: Join the game (server will send updates to our connected WebSocket)
			console.log('[GamePage] Joining game...');
			await joinGame(gameId);
			console.log('[GamePage] Joined game successfully');

			// Step 4: Fetch initial game state as fallback (in case we missed GAME_INIT)
			console.log('[GamePage] Fetching initial game state...');
			const gameView = await fetchGameView(gameId, localPlayerId);
			console.log('[GamePage] Got game state:', {
				players: gameView.players?.length,
				turn: gameView.turn,
				phase: gameView.phase,
				priorityPlayerId: gameView.priorityPlayerId
			});
			gameStore.setGameView(gameView);

			initialized = true;
			console.log('[GamePage] Game initialization complete');
		} catch (err) {
			console.error('[GamePage] Failed to initialize game:', err);
			gameStore.setError(err instanceof Error ? err.message : 'Failed to load game');
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
	async function handlePassUntilEOT() {
		if (!havePriority || isActionLoading || !gameId) return;

		isActionLoading = true;
		try {
			await passUntilEndOfTurn(gameId);
			addLogEntry('You passed until end of turn');
		} catch (err) {
			console.error('Failed to pass until EOT:', err);
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
	 */
	function handleCardClick(cardId: string) {
		// If we have a target prompt, send the selection
		if (prompt?.type === 'target' && gameId) {
			sendPlayerUUID(gameId, cardId).catch(console.error);
			gameStore.clearPrompt();
			return;
		}

		// Otherwise toggle selection
		gameStore.toggleCardSelection(cardId);

		const card = myCards.find((c) => c.id === cardId);
		if (card) {
			addLogEntry(`Selected: ${card.name}`);
		}
	}

	/**
	 * Handle battlefield card click
	 */
	function handleBattlefieldCardClick(cardId: string) {
		// If we have a target prompt, send the selection
		if (prompt?.type === 'target' && gameId) {
			sendPlayerUUID(gameId, cardId).catch(console.error);
			gameStore.clearPrompt();
			return;
		}

		// Toggle selection
		gameStore.toggleCardSelection(cardId);
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
	 * Handle mulligan during mulligan phase
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
	 * Format life total for display (Commander starts at 40)
	 */
	function formatLife(life: number): string {
		return life.toString();
	}

	/**
	 * Get position class for opponent based on index (for 4-player layout)
	 */
	function getOpponentPosition(index: number, total: number): string {
		if (total === 1) return 'top';
		if (total === 2) return index === 0 ? 'left' : 'right';
		if (total === 3) {
			if (index === 0) return 'left';
			if (index === 1) return 'top';
			return 'right';
		}
		return 'top';
	}

	// Initialize on mount
	onMount(() => {
		if (!$auth.isAuthenticated) {
			goto('/login');
			return;
		}

		initializeGame();
	});

	// Cleanup on destroy
	onDestroy(() => {
		gameStore.reset();
	});
</script>

<svelte:head>
	<title>Game {gameId} - MAGE</title>
</svelte:head>

<div class="game-container">
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
		<!-- Mulligan Phase -->
		<MulliganDialog
			cards={myCards}
			mulliganCount={mulliganCount}
			onKeep={handleKeepHand}
			onMulligan={handleMulligan}
			isLoading={isMulliganLoading}
		/>
	{:else}
		<!-- Game Header -->
		<div class="game-header">
			<div class="game-info">
				<div class="format-badge">Commander</div>
				<div class="turn-info">
					<span class="turn-number">Turn {turn}</span>
					<span class="active-player">{activePlayerName()}'s turn</span>
				</div>
			</div>
			<div class="header-actions">
				{#if stackCards.length > 0}
					<button class="btn-stack" onclick={toggleStack}>
						Stack ({stackCards.length})
					</button>
				{/if}
				<button class="btn-concede" onclick={handleConcede}>Concede</button>
			</div>
		</div>

		<!-- Phase & Priority Row -->
		<div class="phase-section">
			<div class="phase-priority-row">
				<PhaseIndicator
					currentPhase={toGamePhase(phase)}
					activePlayerId={gameState.gameView?.activePlayerId || ''}
					{localPlayerId}
					animated={true}
				/>
				<PriorityIndicator
					hasPriority={havePriority}
					activePlayerId={gameState.gameView?.activePlayerId || ''}
					{localPlayerId}
					playerName={activePlayerName()}
					animated={true}
				/>
			</div>
		</div>

		<!-- Prompt Overlay -->
		{#if prompt}
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
					{:else if prompt.type === 'target'}
						<p class="prompt-hint">Click a valid target to select it</p>
					{/if}
				</div>
			</div>
		{/if}

		<!-- Main 4-Player Layout -->
		<div class="game-layout commander-layout">
			<!-- Opponent Areas (up to 3 opponents) -->
			{#each otherPlayers as opponent, idx (opponent.playerId)}
				{@const position = getOpponentPosition(idx, otherPlayers.length)}
				<div class="opponent-area opponent-{position}">
					<div class="player-panel">
						<div class="player-header">
							<span class="player-name" class:has-priority={opponent.hasPriority}>
								{opponent.name}
								{#if opponent.hasPriority}
									<span class="priority-dot"></span>
								{/if}
							</span>
							<div class="player-stats">
								<span class="life" title="Life">{formatLife(opponent.life)}</span>
								<span class="poison" title="Poison" class:active={opponent.poison > 0}>
									{opponent.poison}
								</span>
								<span class="library" title="Library">{opponent.libraryCount}</span>
								<span class="hand-count" title="Hand">{opponent.handCount}</span>
							</div>
						</div>
						<!-- Opponent's battlefield section -->
						<div class="opponent-battlefield">
							{#each getPlayerBattlefieldCards(opponent.playerId) as card (card.id)}
								<div class="battlefield-card">
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
									/>
								</div>
							{/each}
						</div>
						<!-- Opponent zones -->
						<div class="opponent-zones">
							<Graveyard
								cards={opponent.graveyard.map(toGameCard)}
								playerName={opponent.name}
								isOpponent={true}
								onCardClick={() => {}}
							/>
						</div>
					</div>
				</div>
			{/each}

			<!-- Central Battlefield (shared zone) -->
			<div class="central-battlefield">
				<div class="zone-label">Battlefield</div>
				<!-- Command Zone -->
				{#if commandCards.length > 0}
					<div class="command-zone">
						<span class="zone-sublabel">Command Zone</span>
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
								/>
							{/each}
						</div>
					</div>
				{/if}
			</div>

			<!-- Local Player Area -->
			<div class="player-area">
				{#if me}
					<div class="player-panel local-player">
						<div class="player-header">
							<span class="player-name" class:has-priority={havePriority}>
								You
								{#if havePriority}
									<span class="priority-dot"></span>
								{/if}
							</span>
							<div class="player-stats">
								<span class="life" title="Life">{formatLife(me.life)}</span>
								<span class="poison" title="Poison" class:active={me.poison > 0}>
									{me.poison}
								</span>
								<span class="library" title="Library">{me.libraryCount}</span>
							</div>
						</div>

						<!-- Local player's battlefield -->
						<div class="my-battlefield">
							{#each getPlayerBattlefieldCards(localPlayerId) as card (card.id)}
								<div class="battlefield-card">
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
									/>
								</div>
							{/each}
						</div>

						<!-- Player zones row -->
						<div class="zones-row">
							<Graveyard
								cards={myGrave.map(toGameCard)}
								playerName="You"
								isOpponent={false}
								onCardClick={() => {}}
							/>
							<ManaPool
								mana={myMana}
								showEmpty={false}
								size="normal"
								onManaClick={() => {}}
							/>
						</div>

						<!-- Player hand -->
						<PlayerHand
							cards={myCards.map(toGameCard)}
							selectedCardIds={gameState.selectedCardIds}
							onCardClick={handleCardClick}
							size="normal"
						/>
					</div>
				{/if}
			</div>
		</div>

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
						onStackObjectClick={() => {}}
					/>
				</div>
			</div>
		{/if}

		<!-- Sidebars -->
		<ActionLog bind:this={actionLogRef} bind:collapsed={actionLogCollapsed} />
		<div class="game-sidebar" class:collapsed={chatCollapsed}>
			<GameChat bind:this={gameChatRef} gameId={gameId || ''} bind:collapsed={chatCollapsed} />
			<!-- Action buttons -->
			<div class="sidebar-actions">
				<GameActionsPanel
					hasPriority={havePriority}
					canPassPriority={havePriority}
					isLoading={isActionLoading}
					onPassPriority={handlePassPriority}
					onCastSpell={() => {}}
					onActivateAbility={() => {}}
				/>
				<button class="btn-f6" onclick={handlePassUntilEOT} disabled={!havePriority}>
					F6 (Pass Turn)
				</button>
			</div>
		</div>
	{/if}
</div>

<style>
	.game-container {
		position: fixed;
		inset: 0;
		background: #0f1419;
		color: white;
		display: flex;
		flex-direction: column;
		overflow: hidden;
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
		background: #0f1419;
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

	/* Game Header */
	.game-header {
		background: #1a1f2e;
		padding: 0.75rem 1.5rem;
		display: flex;
		justify-content: space-between;
		align-items: center;
		border-bottom: 2px solid #2a3441;
		z-index: 10;
	}

	.game-info {
		display: flex;
		align-items: center;
		gap: 2rem;
	}

	.format-badge {
		padding: 0.375rem 0.75rem;
		background: #667eea;
		border-radius: 4px;
		font-size: 0.875rem;
		font-weight: 600;
	}

	.turn-info {
		display: flex;
		gap: 1rem;
		font-size: 0.9375rem;
	}

	.turn-number {
		color: #fbbf24;
		font-weight: 600;
	}

	.active-player {
		color: #94a3b8;
	}

	.header-actions {
		display: flex;
		gap: 1rem;
	}

	.btn-stack,
	.btn-concede,
	.btn-primary,
	.btn-f6 {
		padding: 0.5rem 1rem;
		border: none;
		border-radius: 4px;
		font-size: 0.875rem;
		font-weight: 600;
		cursor: pointer;
		transition: background 0.2s;
	}

	.btn-stack {
		background: #667eea;
		color: white;
	}

	.btn-stack:hover { background: #5568d3; }

	.btn-concede {
		background: #ef4444;
		color: white;
	}

	.btn-concede:hover { background: #dc2626; }

	.btn-primary {
		background: #667eea;
		color: white;
	}

	.btn-primary:hover { background: #5568d3; }

	.btn-f6 {
		background: #374151;
		color: white;
	}

	.btn-f6:hover:not(:disabled) { background: #4b5563; }
	.btn-f6:disabled { opacity: 0.5; cursor: not-allowed; }

	/* Phase Section */
	.phase-section {
		padding: 0.5rem 1rem;
		margin: 0 320px;
	}

	.phase-priority-row {
		display: flex;
		gap: 1rem;
		align-items: center;
	}

	/* Prompt Overlay */
	.prompt-overlay {
		position: fixed;
		inset: 0;
		background: rgba(0, 0, 0, 0.8);
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
		border-radius: 4px;
		font-weight: 600;
		cursor: pointer;
	}

	.btn-yes:hover { background: #16a34a; }

	.btn-no {
		background: #ef4444;
		color: white;
		padding: 0.75rem 2rem;
		border: none;
		border-radius: 4px;
		font-weight: 600;
		cursor: pointer;
	}

	.btn-no:hover { background: #dc2626; }

	.btn-choice {
		background: #374151;
		color: white;
		padding: 0.5rem 1rem;
		border: none;
		border-radius: 4px;
		cursor: pointer;
	}

	.btn-choice:hover { background: #4b5563; }

	/* 4-Player Commander Layout */
	.game-layout.commander-layout {
		flex: 1;
		display: grid;
		grid-template-areas:
			"left top right"
			"left center right"
			"bottom bottom bottom";
		grid-template-columns: 250px 1fr 250px;
		grid-template-rows: 1fr 200px 400px;
		gap: 0.5rem;
		padding: 0.5rem;
		margin: 0 320px;
		overflow: hidden;
	}

	/* Opponent areas */
	.opponent-area {
		background: #1a1f2e;
		border-radius: 8px;
		border: 1px solid #2a3441;
		overflow: hidden;
	}

	.opponent-left { grid-area: left; }
	.opponent-top { grid-area: top; }
	.opponent-right { grid-area: right; }

	.player-panel {
		height: 100%;
		display: flex;
		flex-direction: column;
		padding: 0.5rem;
	}

	.player-header {
		display: flex;
		justify-content: space-between;
		align-items: center;
		padding: 0.5rem;
		background: #141821;
		border-radius: 4px;
		margin-bottom: 0.5rem;
	}

	.player-name {
		font-weight: 600;
		display: flex;
		align-items: center;
		gap: 0.5rem;
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
		0%, 100% { opacity: 1; }
		50% { opacity: 0.5; }
	}

	.player-stats {
		display: flex;
		gap: 0.75rem;
		font-size: 0.875rem;
	}

	.player-stats span {
		display: flex;
		align-items: center;
		gap: 0.25rem;
	}

	.player-stats .life { color: #ef4444; font-weight: 700; }
	.player-stats .poison { color: #6b7280; }
	.player-stats .poison.active { color: #a855f7; }
	.player-stats .library { color: #3b82f6; }
	.player-stats .hand-count { color: #fbbf24; }

	.opponent-battlefield,
	.my-battlefield {
		flex: 1;
		display: flex;
		flex-wrap: wrap;
		gap: 0.25rem;
		align-content: flex-start;
		overflow-y: auto;
	}

	.opponent-zones {
		display: flex;
		gap: 0.5rem;
		padding-top: 0.5rem;
		border-top: 1px solid #2a3441;
	}

	/* Central battlefield */
	.central-battlefield {
		grid-area: center;
		background: #0d1117;
		border-radius: 8px;
		border: 1px solid #2a3441;
		padding: 1rem;
		display: flex;
		flex-direction: column;
	}

	.zone-label {
		font-size: 0.75rem;
		color: #6b7280;
		text-transform: uppercase;
		letter-spacing: 0.5px;
		margin-bottom: 0.5rem;
	}

	.zone-sublabel {
		font-size: 0.625rem;
		color: #4b5563;
	}

	.command-zone {
		margin-bottom: 1rem;
	}

	.command-cards {
		display: flex;
		gap: 0.5rem;
		flex-wrap: wrap;
	}

	/* Local player area */
	.player-area {
		grid-area: bottom;
		background: #1a1f2e;
		border-radius: 8px;
		border: 1px solid #2a3441;
	}

	.local-player {
		padding: 1rem;
		height: 100%;
		display: flex;
		flex-direction: column;
		overflow: hidden;
	}

	.zones-row {
		display: flex;
		gap: 1rem;
		padding: 0.5rem 0;
		align-items: flex-start;
	}

	.battlefield-card {
		flex-shrink: 0;
	}

	/* Stack overlay */
	.stack-overlay {
		position: fixed;
		inset: 0;
		background: rgba(0, 0, 0, 0.75);
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

	/* Sidebars */
	.game-sidebar {
		position: fixed;
		right: 0;
		top: 0;
		bottom: 0;
		width: 320px;
		z-index: 20;
		display: flex;
		flex-direction: column;
		padding: 0.5rem;
		gap: 0.5rem;
		overflow-y: auto;
		background: #0f1419;
	}

	.game-sidebar.collapsed {
		width: auto;
	}

	.game-sidebar > :global(.game-chat) {
		flex: 1 1 auto;
		min-height: 0;
		display: flex;
		flex-direction: column;
	}

	.sidebar-actions {
		display: flex;
		flex-direction: column;
		gap: 0.75rem;
		flex-shrink: 0;
		padding-top: 0.5rem;
		border-top: 1px solid #2a3441;
		margin-top: 0.5rem;
	}

	.sidebar-actions .btn-f6 {
		width: 100%;
		padding: 0.75rem 1rem;
		background: #374151;
		color: white;
		border: none;
		border-radius: 6px;
		font-size: 0.875rem;
		font-weight: 600;
		cursor: pointer;
		transition: background 0.2s;
	}

	.sidebar-actions .btn-f6:hover:not(:disabled) {
		background: #4b5563;
	}

	.sidebar-actions .btn-f6:disabled {
		opacity: 0.5;
		cursor: not-allowed;
	}

	/* Responsive */
	@media (max-width: 1400px) {
		.phase-section,
		.game-layout.commander-layout {
			margin: 0 280px;
		}

		.game-sidebar:not(.collapsed) {
			width: 280px;
		}
	}

	@media (max-width: 1024px) {
		.phase-section {
			margin: 0 0 0 48px;
		}

		.game-layout.commander-layout {
			margin: 0 0 0 48px;
			grid-template-columns: 200px 1fr 200px;
		}

		.game-sidebar {
			width: 280px;
		}
	}

	@media (max-width: 768px) {
		.phase-section,
		.game-layout.commander-layout {
			margin: 0;
		}

		.game-layout.commander-layout {
			grid-template-areas:
				"top"
				"center"
				"bottom";
			grid-template-columns: 1fr;
			grid-template-rows: 150px 1fr 300px;
		}

		.opponent-left,
		.opponent-right {
			display: none;
		}
	}
</style>
