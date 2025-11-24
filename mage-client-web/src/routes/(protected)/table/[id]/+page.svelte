<script lang="ts">
	import { page } from '$app/stores';
	import { goto } from '$app/navigation';
	import { onMount } from 'svelte';
	import { auth } from '$lib/stores/auth';
	import { confirm } from '$lib/stores/confirm';
	import { fetchTable, toggleReady, leaveTable, startGame, kickPlayer } from '$lib/api/table';
	import type { Table } from '$lib/types/table';
	import LoadingSpinner from '$lib/components/LoadingSpinner.svelte';

	// Get table ID from URL
	const tableId = $derived($page.params.id);

	// State
	let table = $state<Table | null>(null);
	let loading = $state(true);
	let error = $state<string | null>(null);
	let togglingReady = $state(false);
	let startingGame = $state(false);

	// Derived state
	const currentPlayer = $derived(
		table?.players.find((p) => p.username === $auth.user?.username)
	);
	const isHost = $derived(currentPlayer?.isHost ?? false);
	const allPlayersReady = $derived(
		table?.players.every((p) => p.isReady) && (table?.players.length ?? 0) >= 2
	);
	const canStartGame = $derived(isHost && allPlayersReady);

	/**
	 * Load table data
	 */
	async function loadTable(): Promise<void> {
		if (!tableId) return;

		loading = true;
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
	 * Handle ready toggle
	 */
	async function handleToggleReady(): Promise<void> {
		if (!currentPlayer || togglingReady || !tableId) return;

		const newReadyState = !currentPlayer.isReady;
		togglingReady = true;

		try {
			await toggleReady(tableId, newReadyState);

			// Update local state
			if (table) {
				const player = table.players.find((p) => p.username === currentPlayer.username);
				if (player) {
					player.isReady = newReadyState;
				}
			}
		} catch (err) {
			error = err instanceof Error ? err.message : 'Failed to toggle ready status';
			console.error('Failed to toggle ready:', err);
			setTimeout(() => (error = null), 3000);
		} finally {
			togglingReady = false;
		}
	}

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
		}
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

	/**
	 * Load table on mount
	 */
	onMount(() => {
		loadTable();
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
			<div class="players-grid" style="grid-template-columns: repeat({table.maxPlayers > 4 ? 3 : 2}, 1fr);">
				{#each table.players as player (player.id)}
					<div class="player-slot occupied" class:is-current={player.username === currentPlayer?.username}>
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
								{#if player.isReady}
									<span class="ready-badge">✓ Ready</span>
								{:else}
									<span class="not-ready-badge">Not Ready</span>
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
				{#each Array(table.maxPlayers - table.players.length).fill(0).map((_, i) => i) as i (i)}
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

		<!-- Actions -->
		<div class="actions">
			{#if !isHost && currentPlayer}
				<button
					class="btn-primary btn-large"
					class:btn-success={currentPlayer.isReady}
					disabled={togglingReady}
					onclick={handleToggleReady}
				>
					{#if togglingReady}
						<LoadingSpinner size="small" color="white" />
					{:else if currentPlayer.isReady}
						✓ Ready
					{:else}
						Ready Up
					{/if}
				</button>
			{/if}

			{#if isHost}
				<button
					class="btn-primary btn-large btn-start-game"
					disabled={!canStartGame || startingGame}
					onclick={handleStartGame}
					title={!allPlayersReady ? 'All players must be ready' : 'Start the game'}
				>
					{#if startingGame}
						<LoadingSpinner size="small" color="white" />
						Starting...
					{:else}
						Start Game
					{/if}
				</button>
				{#if !allPlayersReady}
					<p class="hint">Waiting for all players to ready up...</p>
				{/if}
			{/if}
		</div>
	{/if}
</div>

<style>
	.container {
		max-width: 1200px;
		margin: 0 auto;
		padding: 2rem;
	}

	/* Loading/Error States */
	.loading-container,
	.error-container {
		display: flex;
		flex-direction: column;
		align-items: center;
		justify-content: center;
		gap: 1rem;
		min-height: 400px;
		text-align: center;
	}

	.loading-container p,
	.error-container p {
		color: #6b7280;
		font-size: 1.125rem;
		margin: 0;
	}

	.error-container {
		color: #dc2626;
	}

	.error-container svg {
		color: #dc2626;
	}

	/* Header */
	header {
		margin-bottom: 2rem;
	}

	.header-content {
		display: flex;
		justify-content: space-between;
		align-items: flex-start;
		gap: 2rem;
	}

	h1 {
		margin: 0 0 0.75rem 0;
		font-size: 2rem;
		color: #111827;
	}

	.header-meta {
		display: flex;
		flex-wrap: wrap;
		gap: 0.75rem;
		align-items: center;
	}

	.format-badge {
		display: inline-flex;
		align-items: center;
		padding: 0.375rem 0.875rem;
		background: #667eea;
		color: white;
		border-radius: 1rem;
		font-size: 0.875rem;
		font-weight: 600;
	}

	.password-badge {
		display: inline-flex;
		align-items: center;
		gap: 0.375rem;
		padding: 0.375rem 0.875rem;
		background: #f59e0b;
		color: white;
		border-radius: 1rem;
		font-size: 0.875rem;
		font-weight: 600;
	}

	.player-count {
		color: #6b7280;
		font-size: 0.875rem;
	}

	/* Error Banner */
	.error-banner {
		display: flex;
		align-items: center;
		gap: 0.75rem;
		padding: 1rem;
		background: #fef2f2;
		border: 1px solid #fecaca;
		border-radius: 0.5rem;
		color: #dc2626;
		margin-bottom: 1.5rem;
	}

	.error-banner svg {
		flex-shrink: 0;
	}

	/* Table Info */
	.table-info {
		display: flex;
		gap: 2rem;
		padding: 1.5rem;
		background: white;
		border: 1px solid #e5e7eb;
		border-radius: 0.5rem;
		margin-bottom: 2rem;
	}

	.info-item {
		display: flex;
		align-items: center;
		gap: 0.5rem;
	}

	.info-label {
		font-weight: 600;
		color: #374151;
	}

	.info-value {
		color: #6b7280;
	}

	.status-badge {
		padding: 0.25rem 0.75rem;
		border-radius: 1rem;
		font-size: 0.875rem;
		font-weight: 600;
		text-transform: capitalize;
	}

	.status-waiting {
		background: #fef3c7;
		color: #92400e;
	}

	.status-ready {
		background: #d1fae5;
		color: #065f46;
	}

	.status-playing {
		background: #dbeafe;
		color: #1e40af;
	}

	/* Players Section */
	.players-section {
		margin-bottom: 2rem;
	}

	.players-section h2 {
		margin: 0 0 1rem 0;
		font-size: 1.5rem;
		color: #111827;
	}

	.players-grid {
		display: grid;
		gap: 1rem;
	}

	.player-slot {
		background: white;
		border: 2px solid #e5e7eb;
		border-radius: 0.75rem;
		padding: 1.5rem;
		display: flex;
		flex-direction: column;
		align-items: center;
		gap: 1rem;
		min-height: 180px;
		transition: all 0.2s;
	}

	.player-slot.occupied {
		border-color: #667eea;
	}

	.player-slot.is-current {
		border-color: #10b981;
		background: #f0fdf4;
	}

	.player-slot.empty {
		background: #f9fafb;
		border-style: dashed;
		color: #9ca3af;
	}

	.player-avatar {
		width: 64px;
		height: 64px;
		border-radius: 50%;
		background: #667eea;
		display: flex;
		align-items: center;
		justify-content: center;
		color: white;
	}

	.player-slot.is-current .player-avatar {
		background: #10b981;
	}

	.empty-icon {
		opacity: 0.3;
	}

	.player-info {
		text-align: center;
		width: 100%;
	}

	.player-info h3 {
		margin: 0 0 0.5rem 0;
		font-size: 1.125rem;
		color: #111827;
	}

	.you-badge {
		font-size: 0.875rem;
		color: #10b981;
		font-weight: 600;
	}

	.player-badges {
		display: flex;
		align-items: center;
		justify-content: center;
		gap: 0.5rem;
		margin-bottom: 0.5rem;
	}

	.host-badge {
		display: inline-block;
		padding: 0.25rem 0.625rem;
		background: #fbbf24;
		color: #78350f;
		border-radius: 0.75rem;
		font-size: 0.75rem;
		font-weight: 600;
	}

	.ready-badge {
		display: inline-block;
		padding: 0.25rem 0.625rem;
		background: #10b981;
		color: white;
		border-radius: 0.75rem;
		font-size: 0.75rem;
		font-weight: 600;
	}

	.not-ready-badge {
		display: inline-block;
		padding: 0.25rem 0.625rem;
		background: #9ca3af;
		color: white;
		border-radius: 0.75rem;
		font-size: 0.75rem;
		font-weight: 600;
	}

	.btn-kick {
		margin-top: 0.5rem;
		padding: 0.375rem 0.75rem;
		background: white;
		color: #ef4444;
		border: 1px solid #ef4444;
		border-radius: 0.375rem;
		font-size: 0.75rem;
		font-weight: 600;
		cursor: pointer;
		transition: all 0.2s;
	}

	.btn-kick:hover {
		background: #ef4444;
		color: white;
	}

	/* Actions */
	.actions {
		display: flex;
		flex-direction: column;
		align-items: center;
		gap: 1rem;
	}

	.hint {
		margin: 0;
		color: #6b7280;
		font-size: 0.875rem;
		font-style: italic;
	}

	/* Buttons */
	.btn-primary,
	.btn-danger {
		display: inline-flex;
		align-items: center;
		justify-content: center;
		gap: 0.5rem;
		padding: 0.75rem 1.5rem;
		border: none;
		border-radius: 0.5rem;
		font-size: 1rem;
		font-weight: 600;
		cursor: pointer;
		transition: all 0.2s;
	}

	.btn-primary {
		background: #667eea;
		color: white;
	}

	.btn-primary:hover:not(:disabled) {
		background: #5568d3;
	}

	.btn-primary:disabled {
		background: #9ca3af;
		cursor: not-allowed;
		opacity: 0.6;
	}

	.btn-success {
		background: #10b981;
	}

	.btn-success:hover:not(:disabled) {
		background: #059669;
	}

	.btn-large {
		padding: 1rem 2rem;
		font-size: 1.125rem;
	}

	.btn-start-game {
		min-width: 200px;
	}

	.btn-danger {
		background: #ef4444;
		color: white;
	}

	.btn-danger:hover {
		background: #dc2626;
	}

	/* Responsive */
	@media (max-width: 768px) {
		.container {
			padding: 1rem;
		}

		.header-content {
			flex-direction: column;
			align-items: stretch;
		}

		h1 {
			font-size: 1.5rem;
		}

		.players-grid {
			grid-template-columns: 1fr !important;
		}

		.table-info {
			flex-direction: column;
			gap: 1rem;
		}
	}
</style>
