<script lang="ts">
	import { onMount, onDestroy } from 'svelte';
	import { goto } from '$app/navigation';
	import { auth } from '$lib/stores/auth';
	import { fetchTables } from '$lib/api/lobby';
	import type { Table } from '$lib/types/table';

	// State
	let tables = $state<Table[]>([]);
	let loading = $state(true);
	let error = $state<string | null>(null);
	let lastRefresh = $state<Date | null>(null);
	let autoRefresh = $state(true);
	let refreshInterval: ReturnType<typeof setInterval> | null = null;

	// Computed
	const activeGames = $derived(
		tables.filter((t) => t.state === 'PLAYING' || t.state === 'SIDEBOARDING')
	);
	const waitingTables = $derived(
		tables.filter((t) => t.state === 'WAITING' || t.state === 'READY')
	);
	const finishedTables = $derived(tables.filter((t) => t.state === 'FINISHED'));

	/**
	 * Fetch all tables from the server
	 */
	async function loadTables() {
		try {
			loading = true;
			error = null;
			tables = await fetchTables();
			lastRefresh = new Date();
		} catch (err) {
			console.error('[AdminGames] Failed to load tables:', err);
			error = err instanceof Error ? err.message : 'Failed to load tables';
		} finally {
			loading = false;
		}
	}

	/**
	 * Toggle auto-refresh
	 */
	function toggleAutoRefresh() {
		autoRefresh = !autoRefresh;
		if (autoRefresh) {
			startAutoRefresh();
		} else {
			stopAutoRefresh();
		}
	}

	/**
	 * Start auto-refresh interval
	 */
	function startAutoRefresh() {
		if (refreshInterval) return;
		refreshInterval = setInterval(loadTables, 3000);
	}

	/**
	 * Stop auto-refresh interval
	 */
	function stopAutoRefresh() {
		if (refreshInterval) {
			clearInterval(refreshInterval);
			refreshInterval = null;
		}
	}

	/**
	 * Navigate to game debug view
	 */
	function openGameDebug(gameId: string) {
		goto(`/game/${gameId}/debug`);
	}

	/**
	 * Navigate to game view
	 */
	function openGame(gameId: string) {
		goto(`/game/${gameId}`);
	}

	/**
	 * Format table state for display
	 */
	function formatState(state: string): { label: string; class: string } {
		switch (state) {
			case 'PLAYING':
				return { label: '🎮 Playing', class: 'state-playing' };
			case 'SIDEBOARDING':
				return { label: '📋 Sideboarding', class: 'state-sideboard' };
			case 'WAITING':
				return { label: '⏳ Waiting', class: 'state-waiting' };
			case 'READY':
				return { label: '✅ Ready', class: 'state-ready' };
			case 'FINISHED':
				return { label: '🏁 Finished', class: 'state-finished' };
			default:
				return { label: state, class: 'state-unknown' };
		}
	}

	/**
	 * Format time ago
	 */
	function timeAgo(date: Date): string {
		const seconds = Math.floor((Date.now() - date.getTime()) / 1000);
		if (seconds < 60) return `${seconds}s ago`;
		const minutes = Math.floor(seconds / 60);
		if (minutes < 60) return `${minutes}m ago`;
		const hours = Math.floor(minutes / 60);
		return `${hours}h ${minutes % 60}m ago`;
	}

	onMount(() => {
		if (!$auth.isAuthenticated) {
			goto('/login');
			return;
		}

		loadTables();
		if (autoRefresh) {
			startAutoRefresh();
		}
	});

	onDestroy(() => {
		stopAutoRefresh();
	});
</script>

<svelte:head>
	<title>Admin - All Games - MAGE</title>
</svelte:head>

<div class="admin-container">
	<header class="admin-header">
		<div class="header-left">
			<a href="/lobby" class="back-link">← Back to Lobby</a>
			<h1>🛠️ Admin: Game Monitor</h1>
		</div>
		<div class="header-right">
			<div class="stats-bar">
				<div class="stat">
					<span class="stat-value">{activeGames.length}</span>
					<span class="stat-label">Active</span>
				</div>
				<div class="stat">
					<span class="stat-value">{waitingTables.length}</span>
					<span class="stat-label">Waiting</span>
				</div>
				<div class="stat">
					<span class="stat-value">{tables.length}</span>
					<span class="stat-label">Total</span>
				</div>
			</div>
			<div class="refresh-controls">
				<button class="btn-refresh" onclick={loadTables} disabled={loading}>
					{loading ? '⟳' : '↻'} Refresh
				</button>
				<label class="auto-refresh-toggle">
					<input type="checkbox" checked={autoRefresh} onchange={toggleAutoRefresh} />
					Auto (3s)
				</label>
			</div>
			{#if lastRefresh}
				<span class="last-refresh">Updated: {timeAgo(lastRefresh)}</span>
			{/if}
		</div>
	</header>

	{#if error}
		<div class="error-banner">
			<span>⚠️ {error}</span>
			<button onclick={loadTables}>Retry</button>
		</div>
	{/if}

	<main class="admin-content">
		<!-- Active Games Section -->
		<section class="game-section">
			<h2 class="section-title">
				<span class="section-icon">🎮</span>
				Active Games ({activeGames.length})
			</h2>
			{#if activeGames.length === 0}
				<div class="empty-state">No active games</div>
			{:else}
				<div class="games-grid">
					{#each activeGames as table (table.id)}
						{@const stateInfo = formatState(table.state)}
						<div class="game-card active">
							<div class="game-header">
								<span class="game-name">{table.name}</span>
								<span class="game-state {stateInfo.class}">{stateInfo.label}</span>
							</div>
							<div class="game-info">
								<div class="info-row">
									<span class="label">Table ID:</span>
									<code class="value">{table.id.slice(0, 8)}...</code>
								</div>
								{#if table.gameId}
									<div class="info-row">
										<span class="label">Game ID:</span>
										<code class="value">{table.gameId.slice(0, 8)}...</code>
									</div>
								{/if}
								<div class="info-row">
									<span class="label">Format:</span>
									<span class="value">{table.deckType || 'Unknown'}</span>
								</div>
								<div class="info-row">
									<span class="label">Players:</span>
									<span class="value"
										>{table.seats?.filter((s) => s?.player).length || 0}/{table.seats?.length ||
											0}</span
									>
								</div>
							</div>
							<div class="game-players">
								{#each table.seats || [] as seat}
									{#if seat?.player}
										<span class="player-badge">{seat.player.name}</span>
									{/if}
								{/each}
							</div>
							<div class="game-actions">
								{#if table.gameId}
									<button class="btn-debug" onclick={() => openGameDebug(table.gameId!)}>
										🔧 Debug
									</button>
									<button class="btn-spectate" onclick={() => openGame(table.gameId!)}>
										👁️ View
									</button>
								{/if}
							</div>
						</div>
					{/each}
				</div>
			{/if}
		</section>

		<!-- Waiting Tables Section -->
		<section class="game-section">
			<h2 class="section-title">
				<span class="section-icon">⏳</span>
				Waiting Tables ({waitingTables.length})
			</h2>
			{#if waitingTables.length === 0}
				<div class="empty-state">No waiting tables</div>
			{:else}
				<div class="tables-list">
					{#each waitingTables as table (table.id)}
						{@const stateInfo = formatState(table.state)}
						<div class="table-row">
							<span class="table-name">{table.name}</span>
							<span class="table-state {stateInfo.class}">{stateInfo.label}</span>
							<span class="table-format">{table.deckType || 'Unknown'}</span>
							<span class="table-players">
								{table.seats?.filter((s) => s?.player).length || 0}/{table.seats?.length || 0} players
							</span>
							<code class="table-id">{table.id.slice(0, 8)}...</code>
						</div>
					{/each}
				</div>
			{/if}
		</section>

		<!-- Recently Finished Section -->
		{#if finishedTables.length > 0}
			<section class="game-section collapsed">
				<h2 class="section-title">
					<span class="section-icon">🏁</span>
					Recently Finished ({finishedTables.length})
				</h2>
				<div class="tables-list">
					{#each finishedTables.slice(0, 10) as table (table.id)}
						<div class="table-row finished">
							<span class="table-name">{table.name}</span>
							<span class="table-format">{table.deckType || 'Unknown'}</span>
							<code class="table-id">{table.id.slice(0, 8)}...</code>
						</div>
					{/each}
				</div>
			</section>
		{/if}
	</main>
</div>

<style>
	.admin-container {
		min-height: 100vh;
		background: #1a1a2e;
		color: #e0e0e0;
		font-family: 'JetBrains Mono', 'Fira Code', monospace;
	}

	.admin-header {
		display: flex;
		justify-content: space-between;
		align-items: center;
		padding: 1rem 1.5rem;
		background: linear-gradient(135deg, #16213e, #1a1a2e);
		border-bottom: 2px solid #0f3460;
		flex-wrap: wrap;
		gap: 1rem;
	}

	.header-left {
		display: flex;
		align-items: center;
		gap: 1.5rem;
	}

	.back-link {
		color: #888;
		text-decoration: none;
		font-size: 0.875rem;
	}

	.back-link:hover {
		color: #00d9ff;
	}

	.admin-header h1 {
		font-size: 1.25rem;
		font-weight: 600;
		margin: 0;
		color: #00d9ff;
	}

	.header-right {
		display: flex;
		align-items: center;
		gap: 1.5rem;
	}

	.stats-bar {
		display: flex;
		gap: 1.5rem;
	}

	.stat {
		display: flex;
		flex-direction: column;
		align-items: center;
		padding: 0.5rem 1rem;
		background: rgba(0, 217, 255, 0.1);
		border: 1px solid rgba(0, 217, 255, 0.3);
		border-radius: 8px;
	}

	.stat-value {
		font-size: 1.5rem;
		font-weight: 700;
		color: #00d9ff;
	}

	.stat-label {
		font-size: 0.625rem;
		text-transform: uppercase;
		color: #888;
	}

	.refresh-controls {
		display: flex;
		align-items: center;
		gap: 0.75rem;
	}

	.btn-refresh {
		padding: 0.5rem 1rem;
		background: #0f3460;
		border: 1px solid #00d9ff;
		border-radius: 6px;
		color: #00d9ff;
		font-family: inherit;
		font-size: 0.875rem;
		cursor: pointer;
		transition: all 0.2s;
	}

	.btn-refresh:hover:not(:disabled) {
		background: #00d9ff;
		color: #1a1a2e;
	}

	.btn-refresh:disabled {
		opacity: 0.5;
		cursor: not-allowed;
	}

	.auto-refresh-toggle {
		display: flex;
		align-items: center;
		gap: 0.5rem;
		font-size: 0.75rem;
		color: #888;
		cursor: pointer;
	}

	.auto-refresh-toggle input {
		accent-color: #00d9ff;
	}

	.last-refresh {
		font-size: 0.75rem;
		color: #666;
	}

	.error-banner {
		display: flex;
		justify-content: space-between;
		align-items: center;
		padding: 0.75rem 1.5rem;
		background: rgba(255, 82, 82, 0.15);
		border-bottom: 1px solid rgba(255, 82, 82, 0.3);
		color: #ff5252;
	}

	.error-banner button {
		padding: 0.375rem 0.75rem;
		background: transparent;
		border: 1px solid #ff5252;
		border-radius: 4px;
		color: #ff5252;
		cursor: pointer;
	}

	.admin-content {
		padding: 1.5rem;
		max-width: 1600px;
		margin: 0 auto;
	}

	.game-section {
		margin-bottom: 2rem;
	}

	.section-title {
		display: flex;
		align-items: center;
		gap: 0.75rem;
		font-size: 1rem;
		font-weight: 600;
		margin: 0 0 1rem 0;
		color: #fff;
	}

	.section-icon {
		font-size: 1.25rem;
	}

	.empty-state {
		padding: 2rem;
		text-align: center;
		color: #666;
		background: rgba(255, 255, 255, 0.02);
		border: 1px dashed #333;
		border-radius: 8px;
	}

	.games-grid {
		display: grid;
		grid-template-columns: repeat(auto-fill, minmax(320px, 1fr));
		gap: 1rem;
	}

	.game-card {
		background: #16213e;
		border: 1px solid #0f3460;
		border-radius: 12px;
		overflow: hidden;
		transition: all 0.2s;
	}

	.game-card.active {
		border-color: #00d9ff;
		box-shadow: 0 0 20px rgba(0, 217, 255, 0.1);
	}

	.game-card:hover {
		transform: translateY(-2px);
		box-shadow: 0 4px 20px rgba(0, 0, 0, 0.3);
	}

	.game-header {
		display: flex;
		justify-content: space-between;
		align-items: center;
		padding: 0.875rem 1rem;
		background: #0f3460;
	}

	.game-name {
		font-weight: 600;
		color: #fff;
	}

	.game-state {
		font-size: 0.75rem;
		padding: 0.25rem 0.625rem;
		border-radius: 999px;
	}

	.state-playing {
		background: rgba(0, 255, 136, 0.2);
		color: #00ff88;
	}

	.state-sideboard {
		background: rgba(255, 193, 7, 0.2);
		color: #ffc107;
	}

	.state-waiting {
		background: rgba(255, 255, 255, 0.1);
		color: #888;
	}

	.state-ready {
		background: rgba(0, 217, 255, 0.2);
		color: #00d9ff;
	}

	.state-finished {
		background: rgba(255, 255, 255, 0.05);
		color: #666;
	}

	.game-info {
		padding: 1rem;
		display: flex;
		flex-direction: column;
		gap: 0.5rem;
	}

	.info-row {
		display: flex;
		justify-content: space-between;
		font-size: 0.8125rem;
	}

	.info-row .label {
		color: #888;
	}

	.info-row .value {
		color: #e0e0e0;
	}

	.info-row code {
		font-family: inherit;
		color: #00d9ff;
		background: rgba(0, 217, 255, 0.1);
		padding: 0.125rem 0.375rem;
		border-radius: 4px;
	}

	.game-players {
		display: flex;
		flex-wrap: wrap;
		gap: 0.5rem;
		padding: 0 1rem 1rem;
	}

	.player-badge {
		font-size: 0.75rem;
		padding: 0.25rem 0.625rem;
		background: rgba(255, 255, 255, 0.1);
		border-radius: 4px;
		color: #fff;
	}

	.game-actions {
		display: flex;
		gap: 0.5rem;
		padding: 1rem;
		background: rgba(0, 0, 0, 0.2);
		border-top: 1px solid #0f3460;
	}

	.btn-debug,
	.btn-spectate {
		flex: 1;
		padding: 0.625rem;
		border: none;
		border-radius: 6px;
		font-family: inherit;
		font-size: 0.8125rem;
		font-weight: 600;
		cursor: pointer;
		transition: all 0.2s;
	}

	.btn-debug {
		background: #ff6b35;
		color: #fff;
	}

	.btn-debug:hover {
		background: #ff8c5a;
	}

	.btn-spectate {
		background: #0f3460;
		color: #00d9ff;
		border: 1px solid #00d9ff;
	}

	.btn-spectate:hover {
		background: #00d9ff;
		color: #1a1a2e;
	}

	.tables-list {
		display: flex;
		flex-direction: column;
		gap: 0.5rem;
	}

	.table-row {
		display: flex;
		align-items: center;
		gap: 1rem;
		padding: 0.75rem 1rem;
		background: rgba(255, 255, 255, 0.02);
		border: 1px solid #333;
		border-radius: 6px;
	}

	.table-row.finished {
		opacity: 0.5;
	}

	.table-name {
		flex: 1;
		font-weight: 600;
	}

	.table-state {
		font-size: 0.75rem;
		padding: 0.25rem 0.5rem;
		border-radius: 4px;
	}

	.table-format {
		font-size: 0.75rem;
		color: #888;
	}

	.table-players {
		font-size: 0.75rem;
		color: #888;
	}

	.table-id {
		font-size: 0.6875rem;
		color: #00d9ff;
		background: rgba(0, 217, 255, 0.1);
		padding: 0.25rem 0.5rem;
		border-radius: 4px;
	}

	@media (max-width: 768px) {
		.admin-header {
			flex-direction: column;
			align-items: flex-start;
		}

		.header-right {
			flex-wrap: wrap;
		}

		.games-grid {
			grid-template-columns: 1fr;
		}

		.table-row {
			flex-wrap: wrap;
		}
	}
</style>
