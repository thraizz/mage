<script lang="ts">
	import { page } from '$app/stores';
	import { goto } from '$app/navigation';
	import { onMount, onDestroy } from 'svelte';
	import { auth } from '$lib/stores/auth';
	import { confirm } from '$lib/stores/confirm';
	import { websocketStore } from '$lib/stores/websocket';
	import { fetchTable, leaveTable, startGame, kickPlayer } from '$lib/api/table';
	import type { Table } from '$lib/types/table';
	import { usePolling } from '$lib/utils/polling';
	import { subscribeTableUpdates } from '$lib/services/table-updates';
	import LoadingSpinner from '$lib/components/LoadingSpinner.svelte';
	import TableChat from '$lib/components/TableChat.svelte';
	import GameStartCountdown from '$lib/components/GameStartCountdown.svelte';
	import { getMageClient } from '$lib/grpc/client';
	import { CallbackMethod } from '$lib/generated/mage/v1/websocket';

	// Get table ID from URL
	const tableId = $derived($page.params.id ?? '');

	// State
	let table = $state<Table | null>(null);
	let loading = $state(true);
	let error = $state<string | null>(null);
	let startingGame = $state(false);
	let showCountdown = $state(false);

	// Derived state
	const currentPlayer = $derived(table?.players.find((p) => p.username === $auth.user?.username));
	const isHost = $derived(currentPlayer?.isHost ?? false);
	const hasMinPlayers = $derived((table?.players.length ?? 0) >= 2);
	const canStartGame = $derived(isHost && hasMinPlayers);

	/**
	 * Load table data
	 */
	async function loadTable(): Promise<void> {
		if (!tableId) return;

		// Only show loading spinner on initial load
		if (!table) {
			loading = true;
		}
		error = null;

		try {
			table = await fetchTable(tableId);
		} catch (err) {
			error = err instanceof Error ? err.message : 'Failed to load table';
			console.error('Failed to load table:', err);
		} finally {
			loading = false;
		}
	}

	/**
	 * Setup polling for table data
	 */
	usePolling(loadTable, {
		interval: 5000, // Poll every 5s (safety net even when WS connected)
		intervalWhenHidden: 30000, // Poll every 30s when tab hidden
		pollWhenConnected: true, // Continue polling even when WS connected
		immediate: false // Don't fetch immediately (we do it manually in onMount)
	});

	/**
	 * Handle leave table
	 */
	async function handleLeaveTable(): Promise<void> {
		if (!tableId) return;

		const confirmed = await confirm.confirm({
			title: 'Leave Table',
			message: 'Are you sure you want to leave this table?',
			confirmText: 'Leave',
			cancelText: 'Stay',
			destructive: false
		});

		if (!confirmed) return;

		try {
			await leaveTable(tableId);
			// Navigate back to lobby
			goto('/lobby');
		} catch (err) {
			error = err instanceof Error ? err.message : 'Failed to leave table';
			console.error('Failed to leave table:', err);
			setTimeout(() => (error = null), 3000);
		}
	}

	/**
	 * Handle start game (host only)
	 */
	async function handleStartGame(): Promise<void> {
		if (!canStartGame || startingGame || !tableId) return;

		// Show countdown first
		showCountdown = true;
	}

	/**
	 * Handle countdown complete - actually start the game
	 */
	async function handleCountdownComplete(): Promise<void> {
		if (!tableId) return;

		startingGame = true;

		try {
			const gameId = await startGame(tableId);
			// Navigate to game view
			goto(`/game/${gameId}`);
		} catch (err) {
			error = err instanceof Error ? err.message : 'Failed to start game';
			console.error('Failed to start game:', err);
			setTimeout(() => (error = null), 3000);
			startingGame = false;
			showCountdown = false;
		}
	}

	/**
	 * Handle countdown cancelled
	 */
	function handleCountdownCancel(): void {
		showCountdown = false;
		startingGame = false;
	}

	/**
	 * Handle kick player (host only)
	 */
	async function handleKickPlayer(playerId: string, playerName: string): Promise<void> {
		if (!isHost || !tableId) return;

		const confirmed = await confirm.confirm({
			title: 'Kick Player',
			message: `Are you sure you want to kick ${playerName} from the table?`,
			confirmText: 'Kick',
			cancelText: 'Cancel',
			destructive: true
		});

		if (!confirmed) return;

		try {
			await kickPlayer(tableId, playerId);
			// Remove player from local state
			if (table) {
				table.players = table.players.filter((p) => p.id !== playerId);
			}
		} catch (err) {
			error = err instanceof Error ? err.message : 'Failed to kick player';
			console.error('Failed to kick player:', err);
			setTimeout(() => (error = null), 3000);
		}
	}

	// WebSocket subscription cleanup functions
	let unsubscribeTable: (() => void) | null = null;
	let unsubscribeStartGame: (() => void) | null = null;

	/**
	 * Handle table updates from WebSocket
	 */
	function handleTableUpdate(updatedTable: Table): void {
		console.log('[Table] Received table update via WebSocket:', updatedTable);
		table = updatedTable;
	}

	/**
	 * Connect to WebSocket for real-time updates
	 */
	async function connectWebSocket(): Promise<void> {
		try {
			const client = getMageClient();
			const sessionId = await client.ensureSessionId();

			if (!sessionId) {
				console.warn('[Table] No session ID available for WebSocket connection');
				return;
			}

			// Connect to WebSocket
			if (!websocketStore.isConnected()) {
				await websocketStore.connect(sessionId);
				console.log('[Table] WebSocket connected');
			}

			// Subscribe to table updates
			unsubscribeTable = subscribeTableUpdates(tableId, handleTableUpdate);
			console.log('[Table] Subscribed to table updates');

			// Subscribe to START_GAME events to handle when host starts the game
			unsubscribeStartGame = websocketStore.on(CallbackMethod.START_GAME, (data: unknown) => {
				// Extract gameId from the event data
				const eventData = data as { gameId?: string; playerNames?: string[] };
				console.log('[Table] Received START_GAME event:', eventData);

				if (eventData?.gameId) {
					// Navigate to the game page
					console.log('[Table] Game started, navigating to game:', eventData.gameId);
					goto(`/game/${eventData.gameId}`);
				}
			});
			console.log('[Table] Subscribed to START_GAME events');
		} catch (err) {
			console.error('[Table] Failed to connect WebSocket:', err);
		}
	}

	/**
	 * Load table on mount and connect WebSocket
	 */
	onMount(async () => {
		loadTable();
		await connectWebSocket();
	});

	/**
	 * Cleanup on unmount
	 */
	onDestroy(() => {
		// Unsubscribe from table updates
		if (unsubscribeTable) {
			unsubscribeTable();
			unsubscribeTable = null;
		}

		// Unsubscribe from START_GAME events
		if (unsubscribeStartGame) {
			unsubscribeStartGame();
			unsubscribeStartGame = null;
		}

		// Note: We don't disconnect WebSocket here as other pages might still need it
		console.log('[Table] Page unmounted');
	});
</script>

<svelte:head>
	<title>Table {tableId} - MAGE</title>
</svelte:head>

<div class="container">
	{#if loading}
		<div class="loading-container">
			<LoadingSpinner size="large" />
			<p>Loading table...</p>
		</div>
	{:else if error && !table}
		<div class="error-container">
			<svg
				xmlns="http://www.w3.org/2000/svg"
				width="48"
				height="48"
				viewBox="0 0 24 24"
				fill="none"
				stroke="currentColor"
				stroke-width="2"
				stroke-linecap="round"
				stroke-linejoin="round"
			>
				<circle cx="12" cy="12" r="10"></circle>
				<line x1="12" y1="8" x2="12" y2="12"></line>
				<line x1="12" y1="16" x2="12.01" y2="16"></line>
			</svg>
			<p>{error}</p>
			<button class="btn-primary" onclick={loadTable}>Retry</button>
		</div>
	{:else if table}
		<!-- Header -->
		<header>
			<div class="header-content">
				<div>
					<h1>{table.name || `Table #${table.id}`}</h1>
					<div class="header-meta">
						<span class="format-badge">{table.format}</span>
						{#if table.hasPassword}
							<span class="password-badge">
								<svg
									xmlns="http://www.w3.org/2000/svg"
									width="14"
									height="14"
									viewBox="0 0 24 24"
									fill="none"
									stroke="currentColor"
									stroke-width="2"
									stroke-linecap="round"
									stroke-linejoin="round"
								>
									<rect x="3" y="11" width="18" height="11" rx="2" ry="2"></rect>
									<path d="M7 11V7a5 5 0 0 1 10 0v4"></path>
								</svg>
								Password Protected
							</span>
						{/if}
						<span class="player-count">{table.players.length} / {table.maxPlayers} players</span>
					</div>
				</div>
				<button class="btn-danger" onclick={handleLeaveTable}>Leave Table</button>
			</div>
		</header>

		<!-- Error Banner -->
		{#if error}
			<div class="error-banner">
				<svg
					xmlns="http://www.w3.org/2000/svg"
					width="20"
					height="20"
					viewBox="0 0 24 24"
					fill="none"
					stroke="currentColor"
					stroke-width="2"
					stroke-linecap="round"
					stroke-linejoin="round"
				>
					<circle cx="12" cy="12" r="10"></circle>
					<line x1="12" y1="8" x2="12" y2="12"></line>
					<line x1="12" y1="16" x2="12.01" y2="16"></line>
				</svg>
				<span>{error}</span>
			</div>
		{/if}

		<!-- Table Info -->
		<div class="table-info">
			<div class="info-item">
				<span class="info-label">Host:</span>
				<span class="info-value">{table.hostUsername}</span>
			</div>
			<div class="info-item">
				<span class="info-label">Status:</span>
				<span class="status-badge status-{table.status}">{table.status}</span>
			</div>
		</div>

		<!-- Players Grid -->
		<div class="players-section">
			<h2>Players</h2>
			<div
				class="players-grid"
				style="grid-template-columns: repeat({table.maxPlayers > 4 ? 3 : 2}, 1fr);"
			>
				{#each table.players as player (player.id)}
					<div
						class="player-slot occupied"
						class:is-current={player.username === currentPlayer?.username}
					>
						<div class="player-avatar">
							<svg
								xmlns="http://www.w3.org/2000/svg"
								width="48"
								height="48"
								viewBox="0 0 24 24"
								fill="none"
								stroke="currentColor"
								stroke-width="2"
								stroke-linecap="round"
								stroke-linejoin="round"
							>
								<path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"></path>
								<circle cx="12" cy="7" r="4"></circle>
							</svg>
						</div>
						<div class="player-info">
							<h3>
								{player.username}
								{#if player.username === currentPlayer?.username}
									<span class="you-badge">(You)</span>
								{/if}
							</h3>
							<div class="player-badges">
								{#if player.isHost}
									<span class="host-badge">Host</span>
								{/if}
							</div>
							{#if isHost && !player.isHost}
								<button
									class="btn-kick"
									onclick={() => handleKickPlayer(player.id, player.username)}
									title="Kick player"
								>
									Kick
								</button>
							{/if}
						</div>
					</div>
				{/each}

				<!-- Empty Slots -->
				{#each Array(table.maxPlayers - table.players.length)
					.fill(0)
					.map((_, i) => i) as i (i)}
					<div class="player-slot empty">
						<div class="empty-icon">
							<svg
								xmlns="http://www.w3.org/2000/svg"
								width="48"
								height="48"
								viewBox="0 0 24 24"
								fill="none"
								stroke="currentColor"
								stroke-width="2"
								stroke-linecap="round"
								stroke-linejoin="round"
								opacity="0.3"
							>
								<path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"></path>
								<circle cx="12" cy="7" r="4"></circle>
							</svg>
						</div>
						<p>Waiting for player...</p>
					</div>
				{/each}
			</div>
		</div>

		<!-- Main Content Grid: Players + Chat -->
		<div class="content-grid">
			<div class="left-column">
				<!-- Actions -->
				<div class="actions">
					{#if isHost}
						<button
							class="btn-primary btn-large btn-start-game"
							disabled={!canStartGame || startingGame}
							onclick={handleStartGame}
							title={!hasMinPlayers ? 'Need at least 2 players' : 'Start the game'}
						>
							{#if startingGame}
								<LoadingSpinner size="small" color="white" />
								Starting...
							{:else}
								Start Game
							{/if}
						</button>
						{#if !hasMinPlayers}
							<p class="hint">Waiting for more players to join...</p>
						{/if}
					{:else if currentPlayer}
						<p class="hint">Waiting for host to start the game...</p>
					{/if}
				</div>
			</div>

			<!-- Chat Column -->
			<div class="chat-column">
				<TableChat {tableId} />
			</div>
		</div>
	{/if}
</div>

<!-- Game Start Countdown -->
<GameStartCountdown
	bind:show={showCountdown}
	onComplete={handleCountdownComplete}
	onCancel={handleCountdownCancel}
/>

<style>
	.container {
		max-width: 1200px;
		margin: 0 auto;
		padding: var(--space-8);
		background: var(--bg-void);
		min-height: 100vh;
	}

	/* Loading/Error States */
	.loading-container,
	.error-container {
		display: flex;
		flex-direction: column;
		align-items: center;
		justify-content: center;
		gap: var(--space-4);
		min-height: 400px;
		text-align: center;
	}

	.loading-container p,
	.error-container p {
		color: var(--text-muted);
		font-size: var(--text-lg);
		margin: 0;
	}

	.error-container {
		color: var(--status-error);
	}

	.error-container svg {
		color: var(--status-error);
	}

	/* Header */
	header {
		margin-bottom: var(--space-8);
	}

	.header-content {
		display: flex;
		justify-content: space-between;
		align-items: flex-start;
		gap: var(--space-8);
	}

	h1 {
		margin: 0 0 var(--space-3) 0;
		font-family: var(--font-display);
		font-size: var(--text-3xl);
		font-weight: var(--weight-bold);
		color: var(--text-bright);
	}

	.header-meta {
		display: flex;
		flex-wrap: wrap;
		gap: var(--space-3);
		align-items: center;
	}

	.format-badge {
		display: inline-flex;
		align-items: center;
		padding: var(--space-2) var(--space-4);
		background: var(--accent-gold);
		color: var(--bg-void);
		border-radius: var(--radius-full);
		font-size: var(--text-sm);
		font-weight: var(--weight-semibold);
	}

	.password-badge {
		display: inline-flex;
		align-items: center;
		gap: var(--space-2);
		padding: var(--space-2) var(--space-4);
		background: var(--status-warning);
		color: var(--bg-void);
		border-radius: var(--radius-full);
		font-size: var(--text-sm);
		font-weight: var(--weight-semibold);
	}

	.player-count {
		color: var(--text-muted);
		font-size: var(--text-sm);
	}

	/* Error Banner */
	.error-banner {
		display: flex;
		align-items: center;
		gap: var(--space-3);
		padding: var(--space-4);
		background: var(--status-error-dim);
		border: 1px solid var(--status-error);
		border-radius: var(--radius-md);
		color: var(--status-error);
		margin-bottom: var(--space-6);
	}

	.error-banner svg {
		flex-shrink: 0;
	}

	/* Table Info */
	.table-info {
		display: flex;
		gap: var(--space-8);
		padding: var(--space-6);
		background: var(--bg-obsidian);
		border: 1px solid var(--border-subtle);
		border-radius: var(--radius-lg);
		margin-bottom: var(--space-8);
	}

	.info-item {
		display: flex;
		align-items: center;
		gap: var(--space-2);
	}

	.info-label {
		font-weight: var(--weight-semibold);
		color: var(--text-muted);
	}

	.info-value {
		color: var(--text-bright);
	}

	.status-badge {
		padding: var(--space-1) var(--space-3);
		border-radius: var(--radius-full);
		font-size: var(--text-sm);
		font-weight: var(--weight-semibold);
		text-transform: capitalize;
	}

	.status-waiting {
		background: var(--status-warning-dim);
		color: var(--status-warning);
	}

	.status-ready {
		background: var(--status-success-dim);
		color: var(--status-success);
	}

	.status-playing {
		background: var(--status-info-dim);
		color: var(--status-info);
	}

	/* Players Section */
	.players-section {
		margin-bottom: var(--space-8);
	}

	.players-section h2 {
		margin: 0 0 var(--space-4) 0;
		font-family: var(--font-display);
		font-size: var(--text-2xl);
		font-weight: var(--weight-bold);
		color: var(--text-bright);
	}

	.players-grid {
		display: grid;
		gap: var(--space-4);
	}

	.player-slot {
		background: var(--bg-obsidian);
		border: 2px solid var(--border-subtle);
		border-radius: var(--radius-lg);
		padding: var(--space-6);
		display: flex;
		flex-direction: column;
		align-items: center;
		gap: var(--space-4);
		min-height: 180px;
		transition: all var(--transition-fast);
	}

	.player-slot.occupied {
		border-color: var(--accent-gold-dim);
	}

	.player-slot.is-current {
		border-color: var(--status-success);
		background: var(--status-success-dim);
	}

	.player-slot.empty {
		background: var(--bg-slate);
		border-style: dashed;
		border-color: var(--border-default);
	}

	.player-avatar {
		width: 64px;
		height: 64px;
		border-radius: var(--radius-full);
		background: var(--accent-gold);
		display: flex;
		align-items: center;
		justify-content: center;
		color: var(--bg-void);
	}

	.player-slot.is-current .player-avatar {
		background: var(--status-success);
	}

	.empty-icon {
		opacity: 0.3;
		color: var(--text-ghost);
	}

	.player-info {
		text-align: center;
		width: 100%;
	}

	.player-info h3 {
		margin: 0 0 var(--space-2) 0;
		font-size: var(--text-lg);
		font-weight: var(--weight-semibold);
		color: var(--text-bright);
	}

	.you-badge {
		font-size: var(--text-sm);
		color: var(--status-success);
		font-weight: var(--weight-semibold);
	}

	.player-badges {
		display: flex;
		align-items: center;
		justify-content: center;
		gap: var(--space-2);
		margin-bottom: var(--space-2);
	}

	.host-badge {
		display: inline-block;
		padding: var(--space-1) var(--space-3);
		background: var(--accent-gold);
		color: var(--bg-void);
		border-radius: var(--radius-full);
		font-size: var(--text-xs);
		font-weight: var(--weight-semibold);
	}

	.ready-badge {
		display: inline-block;
		padding: var(--space-1) var(--space-3);
		background: var(--status-success);
		color: var(--bg-void);
		border-radius: var(--radius-full);
		font-size: var(--text-xs);
		font-weight: var(--weight-semibold);
	}

	.not-ready-badge {
		display: inline-block;
		padding: var(--space-1) var(--space-3);
		background: var(--text-ghost);
		color: var(--text-bright);
		border-radius: var(--radius-full);
		font-size: var(--text-xs);
		font-weight: var(--weight-semibold);
	}

	.btn-kick {
		margin-top: var(--space-2);
		padding: var(--space-2) var(--space-3);
		background: transparent;
		color: var(--status-error);
		border: 1px solid var(--status-error);
		border-radius: var(--radius-md);
		font-size: var(--text-xs);
		font-weight: var(--weight-semibold);
		cursor: pointer;
		transition: all var(--transition-fast);
	}

	.btn-kick:hover {
		background: var(--status-error);
		color: var(--text-bright);
	}

	/* Content Grid */
	.content-grid {
		display: grid;
		grid-template-columns: 1fr 400px;
		gap: var(--space-8);
		margin-top: var(--space-8);
	}

	.left-column {
		display: flex;
		flex-direction: column;
	}

	.chat-column {
		height: 600px;
	}

	/* Actions */
	.actions {
		display: flex;
		flex-direction: column;
		align-items: center;
		gap: var(--space-4);
	}

	.hint {
		margin: 0;
		color: var(--text-muted);
		font-size: var(--text-sm);
		font-style: italic;
	}

	/* Buttons */
	.btn-primary,
	.btn-danger {
		display: inline-flex;
		align-items: center;
		justify-content: center;
		gap: var(--space-2);
		padding: var(--space-3) var(--space-6);
		border: none;
		border-radius: var(--radius-md);
		font-size: var(--text-base);
		font-weight: var(--weight-semibold);
		cursor: pointer;
		transition: all var(--transition-fast);
	}

	.btn-primary {
		background: var(--accent-gold);
		color: var(--bg-void);
	}

	.btn-primary:hover:not(:disabled) {
		background: var(--accent-gold-bright);
	}

	.btn-primary:disabled {
		background: var(--bg-steel);
		color: var(--text-ghost);
		cursor: not-allowed;
		opacity: 0.6;
	}

	.btn-success {
		background: var(--status-success);
	}

	.btn-success:hover:not(:disabled) {
		background: #059669;
	}

	.btn-large {
		padding: var(--space-4) var(--space-8);
		font-size: var(--text-lg);
	}

	.btn-start-game {
		min-width: 200px;
	}

	.btn-danger {
		background: var(--status-error);
		color: var(--text-bright);
	}

	.btn-danger:hover {
		background: #dc2626;
	}

	/* Responsive */
	@media (max-width: 1024px) {
		.content-grid {
			grid-template-columns: 1fr;
			gap: var(--space-6);
		}

		.chat-column {
			height: 400px;
		}
	}

	@media (max-width: 768px) {
		.container {
			padding: var(--space-4);
		}

		.header-content {
			flex-direction: column;
			align-items: stretch;
		}

		h1 {
			font-size: var(--text-2xl);
		}

		.players-grid {
			grid-template-columns: 1fr !important;
		}

		.table-info {
			flex-direction: column;
			gap: var(--space-4);
		}

		.content-grid {
			gap: var(--space-4);
		}

		.chat-column {
			height: 300px;
		}
	}
</style>
