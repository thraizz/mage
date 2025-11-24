<script lang="ts">
	import { onMount } from 'svelte';
	import { page } from '$app/stores';
	import { goto } from '$app/navigation';
	import { authStore } from '$lib/stores/auth';
	import { toastStore } from '$lib/stores/toast';
	import { confirmStore } from '$lib/stores/confirm';
	import type { GameCard, GamePhase, StackObject, ManaPool } from '$lib/types/game';

	// Import game components
	import Card from '$lib/components/game/Card.svelte';
	import PlayerHand from '$lib/components/game/PlayerHand.svelte';
	import Graveyard from '$lib/components/game/Graveyard.svelte';
	import ExileZone from '$lib/components/game/ExileZone.svelte';
	import ManaPool from '$lib/components/game/ManaPool.svelte';
	import PhaseIndicator from '$lib/components/game/PhaseIndicator.svelte';
	import Stack from '$lib/components/game/Stack.svelte';

	// Game ID from route params
	const gameId = $derived($page.params.id);

	// Game state
	let loading = $state(true);
	let error = $state<string | null>(null);

	// Player state (placeholder data with proper types)
	let localPlayerId = $state('player-1');
	let playerLife = $state(20);
	let playerLibraryCount = $state(53);
	let playerHand = $state<GameCard[]>([
		{
			id: 'card-1',
			name: 'Lightning Bolt',
			manaCost: '{R}',
			cardType: 'Instant',
			imageUrl: ''
		},
		{
			id: 'card-2',
			name: 'Mountain',
			manaCost: '',
			cardType: 'Basic Land — Mountain',
			imageUrl: ''
		},
		{
			id: 'card-3',
			name: 'Goblin Guide',
			manaCost: '{R}',
			cardType: 'Creature — Goblin Scout',
			power: '2',
			toughness: '2',
			imageUrl: ''
		},
		{
			id: 'card-4',
			name: 'Lava Spike',
			manaCost: '{R}',
			cardType: 'Sorcery',
			imageUrl: ''
		},
		{
			id: 'card-5',
			name: 'Rift Bolt',
			manaCost: '{2}{R}',
			cardType: 'Sorcery',
			imageUrl: ''
		},
		{
			id: 'card-6',
			name: 'Mountain',
			manaCost: '',
			cardType: 'Basic Land — Mountain',
			imageUrl: ''
		},
		{
			id: 'card-7',
			name: 'Monastery Swiftspear',
			manaCost: '{R}',
			cardType: 'Creature — Human Monk',
			power: '1',
			toughness: '2',
			imageUrl: ''
		}
	]);
	let playerGraveyard = $state<GameCard[]>([
		{
			id: 'grave-1',
			name: 'Shock',
			manaCost: '{R}',
			cardType: 'Instant',
			imageUrl: ''
		}
	]);
	let playerExile = $state<GameCard[]>([]);
	let playerMana = $state<ManaPool>({
		white: 0,
		blue: 0,
		black: 0,
		red: 3,
		green: 0,
		colorless: 1
	});

	// Opponent state (placeholder data)
	let opponentPlayerId = $state('player-2');
	let opponentLife = $state(18);
	let opponentLibraryCount = $state(52);
	let opponentHandCount = $state(6);
	let opponentName = $state('Opponent');
	let opponentGraveyard = $state<GameCard[]>([]);
	let opponentExile = $state<GameCard[]>([]);

	// Battlefield state (placeholder)
	let battlefieldCards = $state<GameCard[]>([
		{
			id: 'bf-1',
			name: 'Mountain',
			cardType: 'Land',
			isTapped: true,
			imageUrl: ''
		},
		{
			id: 'bf-2',
			name: 'Mountain',
			cardType: 'Land',
			isTapped: false,
			imageUrl: ''
		}
	]);

	// Stack state (placeholder)
	let stackObjects = $state<StackObject[]>([]);

	// Game info state
	let currentTurn = $state(1);
	let currentPhase = $state<GamePhase>('PRECOMBAT_MAIN');
	let activePlayerId = $state(localPlayerId);
	let format = $state('Standard');

	// Selection state
	let selectedCardIds = $state<string[]>([]);
	let showStack = $state(false);

	// Game log
	let gameLog = $state<string[]>([
		'Game started',
		'You drew 7 cards',
		'Your turn begins',
		'Upkeep phase',
		'Draw step - You drew a card',
		'Main phase 1'
	]);

	/**
	 * Load game state
	 */
	async function loadGameState(): Promise<void> {
		loading = true;
		error = null;

		try {
			// TODO: Implement actual game state fetching via API
			// const response = await fetchGameState(gameId);
			// Update all state from response

			// Placeholder: simulate loading delay
			await new Promise((resolve) => setTimeout(resolve, 500));

			// Mock data loaded
			addLogEntry('Game state loaded');
		} catch (err) {
			console.error('Failed to load game state:', err);
			error = err instanceof Error ? err.message : 'Failed to load game';
		} finally {
			loading = false;
		}
	}

	/**
	 * Handle concede game
	 */
	async function handleConcede(): Promise<void> {
		const confirmed = await confirmStore.confirm(
			'Are you sure you want to concede?',
			'You will lose this game and it will count as a loss in your match history. This action cannot be undone.'
		);

		if (!confirmed) return;

		try {
			// TODO: Implement concede API call
			// await concedeGame(gameId);
			toastStore.success('Game conceded');
			goto('/lobby');
		} catch (err) {
			console.error('Failed to concede:', err);
			toastStore.error('Failed to concede game');
		}
	}

	/**
	 * Handle pass priority
	 */
	function handlePassPriority(): void {
		// TODO: Implement priority passing
		addLogEntry('You passed priority');
		toastStore.info('Priority passed');
	}

	/**
	 * Handle card click in hand
	 */
	function handleCardClick(cardId: string): void {
		console.log('Card clicked:', cardId);
		addLogEntry(`Selected card: ${playerHand.find((c) => c.id === cardId)?.name}`);
	}

	/**
	 * Handle mana click
	 */
	function handleManaClick(color: string): void {
		console.log('Mana clicked:', color);
		addLogEntry(`Tapped ${color} mana`);
		// TODO: Implement mana spending
	}

	/**
	 * Handle stack object click
	 */
	function handleStackObjectClick(stackId: string): void {
		console.log('Stack object clicked:', stackId);
		// TODO: Show stack object details
	}

	/**
	 * Handle graveyard card click
	 */
	function handleGraveyardCardClick(cardId: string): void {
		console.log('Graveyard card clicked:', cardId);
		// TODO: Implement graveyard interactions
	}

	/**
	 * Handle exile card click
	 */
	function handleExileCardClick(cardId: string): void {
		console.log('Exile card clicked:', cardId);
		// TODO: Implement exile interactions
	}

	/**
	 * Add message to game log
	 */
	function addLogEntry(message: string): void {
		gameLog = [...gameLog, message];
	}

	/**
	 * Toggle stack display
	 */
	function toggleStack(): void {
		showStack = !showStack;
	}

	// Load game on mount
	onMount(() => {
		if (!$authStore.isAuthenticated) {
			toastStore.error('Please login first');
			goto('/login');
			return;
		}

		loadGameState();
	});

	// Derived values
	const playerNames = $derived(
		new Map([
			[localPlayerId, 'You'],
			[opponentPlayerId, opponentName]
		])
	);
</script>

<svelte:head>
	<title>Game {gameId} - MAGE</title>
</svelte:head>

<div class="game-container">
	{#if loading}
		<div class="loading-overlay">
			<div class="spinner"></div>
			<p>Loading game...</p>
		</div>
	{:else if error}
		<div class="error-overlay">
			<p class="error-message">{error}</p>
			<button class="btn-primary" onclick={() => goto('/lobby')}>Return to Lobby</button>
		</div>
	{:else}
		<!-- Game Header -->
		<div class="game-header">
			<div class="game-info">
				<div class="format-badge">{format}</div>
				<div class="turn-info">
					<span class="turn-number">Turn {currentTurn}</span>
				</div>
			</div>
			<div class="header-actions">
				{#if stackObjects.length > 0}
					<button class="btn-stack" onclick={toggleStack}>
						Stack ({stackObjects.length})
					</button>
				{/if}
				<button class="btn-concede" onclick={handleConcede}>Concede</button>
			</div>
		</div>

		<!-- Phase Indicator -->
		<div class="phase-section">
			<PhaseIndicator
				{currentPhase}
				{activePlayerId}
				{localPlayerId}
				animated={true}
			/>
		</div>

		<!-- Main Game Layout -->
		<div class="game-layout">
			<!-- Opponent Area (Top) -->
			<div class="opponent-area zone">
				<div class="player-info-bar">
					<div class="player-name">{opponentName}</div>
					<div class="player-stats">
						<div class="stat life" title="Life Total">
							<span class="stat-icon">❤️</span>
							<span class="stat-value">{opponentLife}</span>
						</div>
						<div class="stat library" title="Library">
							<span class="stat-icon">📚</span>
							<span class="stat-value">{opponentLibraryCount}</span>
						</div>
						<div class="stat hand" title="Hand Size">
							<span class="stat-icon">🎴</span>
							<span class="stat-value">{opponentHandCount}</span>
						</div>
					</div>
				</div>

				<!-- Opponent Zones Row -->
				<div class="zones-row">
					<Graveyard
						cards={opponentGraveyard}
						playerName={opponentName}
						isOpponent={true}
						onCardClick={handleGraveyardCardClick}
					/>

					<ExileZone
						cards={opponentExile}
						playerName={opponentName}
						isOpponent={true}
						onCardClick={handleExileCardClick}
					/>

					<!-- Opponent Hand (Card Backs) -->
					<div class="opponent-hand">
						{#each Array(opponentHandCount) as _, i (i)}
							<div class="card-back"></div>
						{/each}
					</div>
				</div>
			</div>

			<!-- Battlefield (Middle) -->
			<div class="battlefield zone">
				<div class="zone-label">Battlefield</div>
				{#if battlefieldCards.length === 0}
					<div class="empty-state">
						<p>No permanents in play</p>
					</div>
				{:else}
					<div class="battlefield-grid">
						{#each battlefieldCards as card (card.id)}
							<div class="battlefield-card">
								<Card
									cardId={card.id}
									cardName={card.name}
									manaCost={card.manaCost || ''}
									cardType={card.cardType || ''}
									power={card.power || ''}
									toughness={card.toughness || ''}
									imageUrl={card.imageUrl || ''}
									isTapped={card.isTapped || false}
									isSelected={selectedCardIds.includes(card.id)}
									isPlaceholder={true}
									size="normal"
									onclick={() => {}}
								/>
							</div>
						{/each}
					</div>
				{/if}
			</div>

			<!-- Player Area (Bottom) -->
			<div class="player-area zone">
				<div class="player-info-bar">
					<div class="player-name">You</div>
					<div class="player-stats">
						<div class="stat life" title="Life Total">
							<span class="stat-icon">❤️</span>
							<span class="stat-value">{playerLife}</span>
						</div>
						<div class="stat library" title="Library">
							<span class="stat-icon">📚</span>
							<span class="stat-value">{playerLibraryCount}</span>
						</div>
					</div>
				</div>

				<!-- Player Zones Row -->
				<div class="zones-row">
					<Graveyard
						cards={playerGraveyard}
						playerName="You"
						isOpponent={false}
						onCardClick={handleGraveyardCardClick}
					/>

					<ExileZone
						cards={playerExile}
						playerName="You"
						isOpponent={false}
						onCardClick={handleExileCardClick}
					/>

					<ManaPool mana={playerMana} showEmpty={false} size="normal" onManaClick={handleManaClick} />
				</div>

				<!-- Player Hand -->
				<PlayerHand
					cards={playerHand}
					{selectedCardIds}
					onCardClick={handleCardClick}
					size="normal"
				/>

				<!-- Player Actions -->
				<div class="player-actions">
					<button class="btn-primary" onclick={handlePassPriority}>Pass Priority</button>
					<button class="btn-secondary">Undo</button>
				</div>
			</div>
		</div>

		<!-- Stack Overlay (if shown) -->
		{#if showStack}
			<!-- svelte-ignore a11y_click_events_have_key_events -->
			<!-- svelte-ignore a11y_no_static_element_interactions -->
			<div class="stack-overlay" onclick={toggleStack}>
				<div class="stack-panel" onclick={(e) => e.stopPropagation()}>
					<Stack
						{stackObjects}
						{playerNames}
						onStackObjectClick={handleStackObjectClick}
					/>
				</div>
			</div>
		{/if}

		<!-- Game Log Sidebar -->
		<div class="game-sidebar">
			<div class="sidebar-header">
				<h4>Game Log</h4>
			</div>
			<div class="log-entries">
				{#each gameLog as entry, i (i)}
					<div class="log-entry">{entry}</div>
				{/each}
			</div>
		</div>
	{/if}
</div>

<style>
	/* Container */
	.game-container {
		position: fixed;
		top: 0;
		left: 0;
		right: 0;
		bottom: 0;
		background: #0f1419;
		color: white;
		display: flex;
		flex-direction: column;
		overflow: hidden;
	}

	/* Loading & Error States */
	.loading-overlay,
	.error-overlay {
		position: absolute;
		top: 0;
		left: 0;
		right: 0;
		bottom: 0;
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
		to {
			transform: rotate(360deg);
		}
	}

	.loading-overlay p,
	.error-overlay p {
		color: #999;
		font-size: 1rem;
	}

	.error-message {
		color: #ef4444;
		margin-bottom: 1rem;
	}

	/* Game Header */
	.game-header {
		background: #1a1f2e;
		padding: 0.75rem 1.5rem;
		display: flex;
		justify-content: space-between;
		align-items: center;
		border-bottom: 2px solid #2a3441;
		box-shadow: 0 2px 8px rgba(0, 0, 0, 0.3);
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

	.header-actions {
		display: flex;
		gap: 1rem;
	}

	.btn-stack {
		padding: 0.625rem 1.25rem;
		background: #667eea;
		color: white;
		border: none;
		border-radius: 4px;
		font-size: 0.9375rem;
		font-weight: 600;
		cursor: pointer;
		transition: background 0.2s;
	}

	.btn-stack:hover {
		background: #5568d3;
	}

	.btn-concede {
		padding: 0.625rem 1.25rem;
		background: #ef4444;
		color: white;
		border: none;
		border-radius: 4px;
		font-size: 0.9375rem;
		font-weight: 600;
		cursor: pointer;
		transition: background 0.2s;
	}

	.btn-concede:hover {
		background: #dc2626;
	}

	/* Phase Section */
	.phase-section {
		padding: 1rem 1.5rem 0;
		margin-right: 320px;
	}

	/* Game Layout */
	.game-layout {
		flex: 1;
		display: grid;
		grid-template-rows: 200px 1fr 320px;
		padding: 1rem;
		gap: 1rem;
		overflow: hidden;
		margin-right: 320px;
	}

	.zone {
		background: #1a1f2e;
		border-radius: 8px;
		border: 1px solid #2a3441;
		overflow: hidden;
	}

	/* Player Info Bar */
	.player-info-bar {
		display: flex;
		justify-content: space-between;
		align-items: center;
		padding: 0.75rem 1rem;
		background: #141821;
		border-bottom: 1px solid #2a3441;
	}

	.player-name {
		font-size: 1rem;
		font-weight: 600;
		color: white;
	}

	.player-stats {
		display: flex;
		gap: 1.5rem;
	}

	.stat {
		display: flex;
		align-items: center;
		gap: 0.375rem;
		font-size: 1rem;
		font-weight: 600;
	}

	.stat-icon {
		font-size: 1.125rem;
	}

	.stat.life {
		color: #ef4444;
	}

	.stat.library {
		color: #3b82f6;
	}

	.stat.hand {
		color: #fbbf24;
	}

	/* Zones Row */
	.zones-row {
		display: flex;
		gap: 1rem;
		padding: 1rem;
		align-items: flex-start;
	}

	/* Opponent Area */
	.opponent-area {
		display: flex;
		flex-direction: column;
	}

	.opponent-hand {
		display: flex;
		gap: -20px;
		margin-left: auto;
	}

	.card-back {
		width: 60px;
		height: 84px;
		background: linear-gradient(135deg, #2a3441 0%, #1a1f2e 100%);
		border: 2px solid #3a4451;
		border-radius: 6px;
		margin-left: -15px;
		box-shadow: 0 2px 4px rgba(0, 0, 0, 0.3);
	}

	.card-back:first-child {
		margin-left: 0;
	}

	/* Battlefield */
	.battlefield {
		display: flex;
		flex-direction: column;
		padding: 1rem;
		background: #0d1117;
		position: relative;
	}

	.zone-label {
		font-size: 0.875rem;
		color: #6b7280;
		margin-bottom: 0.5rem;
		font-weight: 500;
		text-transform: uppercase;
		letter-spacing: 0.5px;
	}

	.empty-state {
		flex: 1;
		display: flex;
		align-items: center;
		justify-content: center;
	}

	.empty-state p {
		color: #4b5563;
		font-style: italic;
	}

	.battlefield-grid {
		display: grid;
		grid-template-columns: repeat(auto-fill, minmax(100px, 1fr));
		gap: 1rem;
		flex: 1;
		align-content: start;
	}

	/* Player Area */
	.player-area {
		display: flex;
		flex-direction: column;
	}

	.player-actions {
		padding: 0.75rem 1rem;
		border-top: 1px solid #2a3441;
		display: flex;
		justify-content: center;
		gap: 1rem;
		background: #141821;
	}

	/* Stack Overlay */
	.stack-overlay {
		position: fixed;
		top: 0;
		left: 0;
		right: 0;
		bottom: 0;
		background: rgba(0, 0, 0, 0.75);
		display: flex;
		align-items: center;
		justify-content: center;
		z-index: 100;
		animation: fadeIn 0.2s;
	}

	.stack-panel {
		width: 90%;
		max-width: 500px;
		max-height: 80vh;
	}

	/* Game Sidebar */
	.game-sidebar {
		position: fixed;
		right: 0;
		top: 0;
		bottom: 0;
		width: 320px;
		background: #1a1f2e;
		border-left: 2px solid #2a3441;
		display: flex;
		flex-direction: column;
		box-shadow: -2px 0 8px rgba(0, 0, 0, 0.3);
	}

	.sidebar-header {
		padding: 1rem;
		background: #141821;
		border-bottom: 1px solid #2a3441;
	}

	.sidebar-header h4 {
		margin: 0;
		font-size: 1rem;
		font-weight: 600;
		color: white;
	}

	.log-entries {
		flex: 1;
		overflow-y: auto;
		padding: 1rem;
		display: flex;
		flex-direction: column;
		gap: 0.5rem;
	}

	.log-entry {
		padding: 0.625rem;
		background: #0d1117;
		border-radius: 4px;
		font-size: 0.875rem;
		color: #9ca3af;
		border-left: 2px solid #3a4451;
	}

	/* Buttons */
	.btn-primary {
		padding: 0.75rem 1.5rem;
		background: #667eea;
		color: white;
		border: none;
		border-radius: 4px;
		font-size: 1rem;
		font-weight: 600;
		cursor: pointer;
		transition: background 0.2s;
	}

	.btn-primary:hover {
		background: #5568d3;
	}

	.btn-secondary {
		padding: 0.75rem 1.5rem;
		background: #2a3441;
		color: white;
		border: 1px solid #3a4451;
		border-radius: 4px;
		font-size: 1rem;
		font-weight: 600;
		cursor: pointer;
		transition: all 0.2s;
	}

	.btn-secondary:hover {
		background: #3a4451;
		border-color: #4a5461;
	}

	/* Scrollbar Styling */
	.log-entries::-webkit-scrollbar {
		width: 8px;
	}

	.log-entries::-webkit-scrollbar-track {
		background: #0d1117;
	}

	.log-entries::-webkit-scrollbar-thumb {
		background: #3a4451;
		border-radius: 4px;
	}

	.log-entries::-webkit-scrollbar-thumb:hover {
		background: #4a5461;
	}

	/* Responsive */
	@media (max-width: 1400px) {
		.game-sidebar {
			width: 280px;
		}

		.game-layout,
		.phase-section {
			margin-right: 280px;
		}
	}

	@media (max-width: 1024px) {
		.game-sidebar {
			display: none;
		}

		.game-layout,
		.phase-section {
			margin-right: 0;
		}

		.zones-row {
			flex-wrap: wrap;
		}
	}

	@media (max-width: 768px) {
		.game-layout {
			grid-template-rows: 180px 1fr 300px;
			padding: 0.5rem;
			gap: 0.5rem;
		}

		.phase-section {
			padding: 0.5rem;
		}

		.format-badge {
			display: none;
		}
	}
</style>
