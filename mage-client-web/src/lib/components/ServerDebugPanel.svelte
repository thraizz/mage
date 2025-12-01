<script lang="ts">
	import { onMount, onDestroy } from 'svelte';
	import {
		fetchAllActiveGames,
		fetchServerDebugState,
		type ActiveGameDebug,
		type ServerDebugState,
		type AllActiveGamesResult
	} from '$lib/api/debug';
	import { fetchMatchHistory, type Match } from '$lib/api/match_history';

	// Props
	let { open = $bindable(false) } = $props<{ open?: boolean }>();

	// State
	let loading = $state(true);
	let error = $state<string | null>(null);
	let serverState = $state<ServerDebugState | null>(null);
	let activeGames = $state<AllActiveGamesResult | null>(null);
	let matchHistory = $state<Match[]>([]);
	let matchHistoryCount = $state(0);
	let autoRefresh = $state(false);
	let refreshInterval: ReturnType<typeof setInterval> | null = null;
	let lastRefresh = $state<Date | null>(null);

	// Active tab
	let activeTab = $state<'state' | 'games' | 'history'>('state');

	async function loadData() {
		loading = true;
		error = null;
		try {
			const [stateResult, gamesResult, historyResult] = await Promise.all([
				fetchServerDebugState(),
				fetchAllActiveGames(),
				fetchMatchHistory(10, 0)
			]);
			serverState = stateResult;
			activeGames = gamesResult;
			matchHistory = historyResult.matches;
			matchHistoryCount = historyResult.totalCount;
			lastRefresh = new Date();
		} catch (err) {
			error = err instanceof Error ? err.message : 'Failed to load debug data';
			console.error('[ServerDebugPanel] Error:', err);
		} finally {
			loading = false;
		}
	}

	function toggleAutoRefresh() {
		autoRefresh = !autoRefresh;
		if (autoRefresh) {
			refreshInterval = setInterval(loadData, 3000);
		} else if (refreshInterval) {
			clearInterval(refreshInterval);
			refreshInterval = null;
		}
	}

	function timeAgo(date: Date): string {
		const seconds = Math.floor((Date.now() - date.getTime()) / 1000);
		if (seconds < 60) return `${seconds}s ago`;
		const minutes = Math.floor(seconds / 60);
		if (minutes < 60) return `${minutes}m ago`;
		const hours = Math.floor(minutes / 60);
		return `${hours}h ${minutes % 60}m ago`;
	}

	function formatDuration(seconds: number): string {
		const mins = Math.floor(seconds / 60);
		const secs = seconds % 60;
		return `${mins}m ${secs}s`;
	}

	onMount(() => {
		if (open) {
			loadData();
		}
	});

	onDestroy(() => {
		if (refreshInterval) {
			clearInterval(refreshInterval);
		}
	});

	// Load data when panel opens
	$effect(() => {
		if (open && !lastRefresh) {
			loadData();
		}
	});
</script>

{#if open}
	<div class="debug-panel">
		<div class="debug-header">
			<div class="header-left">
				<span class="debug-icon">🔧</span>
				<h2>Server Debug</h2>
			</div>
			<div class="header-right">
				<div class="refresh-controls">
					<button class="btn-refresh" onclick={loadData} disabled={loading}>
						{loading ? '⟳' : '↻'} Refresh
					</button>
					<label class="auto-toggle">
						<input type="checkbox" checked={autoRefresh} onchange={toggleAutoRefresh} />
						Auto (3s)
					</label>
				</div>
				{#if lastRefresh}
					<span class="last-update">Updated: {timeAgo(lastRefresh)}</span>
				{/if}
				<button class="btn-close" onclick={() => (open = false)}>✕</button>
			</div>
		</div>

		{#if error}
			<div class="error-banner">
				<span>⚠️ {error}</span>
				<button onclick={loadData}>Retry</button>
			</div>
		{/if}

		<div class="tabs">
			<button class="tab" class:active={activeTab === 'state'} onclick={() => (activeTab = 'state')}>
				📊 Server State
			</button>
			<button class="tab" class:active={activeTab === 'games'} onclick={() => (activeTab = 'games')}>
				🎮 Active Games ({activeGames?.games.length ?? 0})
			</button>
			<button class="tab" class:active={activeTab === 'history'} onclick={() => (activeTab = 'history')}>
				📜 Match History ({matchHistoryCount})
			</button>
		</div>

		<div class="tab-content">
			{#if loading && !lastRefresh}
				<div class="loading">Loading debug data...</div>
			{:else if activeTab === 'state' && serverState}
				<!-- Server State Tab -->
				<div class="state-grid">
					<div class="stat-card memory">
						<h3>🧠 Memory State</h3>
						<div class="stat-row">
							<span class="label">Active Games:</span>
							<span class="value">{serverState.memoryActiveGames}</span>
						</div>
						<div class="stat-row">
							<span class="label">Active Tables:</span>
							<span class="value">{serverState.memoryActiveTables}</span>
						</div>
						<div class="stat-row">
							<span class="label">Active Sessions:</span>
							<span class="value">{serverState.memoryActiveSessions}</span>
						</div>
					</div>

					<div class="stat-card database">
						<h3>💾 Database State</h3>
						<div class="stat-row">
							<span class="label">Active Games (DB):</span>
							<span class="value">{serverState.dbActiveGames}</span>
						</div>
						<div class="stat-row">
							<span class="label">Match History:</span>
							<span class="value">{serverState.dbMatchHistory}</span>
						</div>
					</div>

					{#if serverState.gamesInMemoryOnly.length > 0 || serverState.gamesInDbOnly.length > 0}
						<div class="stat-card discrepancies warning">
							<h3>⚠️ Discrepancies</h3>
							{#if serverState.gamesInMemoryOnly.length > 0}
								<div class="discrepancy-section">
									<h4>In Memory Only ({serverState.gamesInMemoryOnly.length})</h4>
									{#each serverState.gamesInMemoryOnly as gameId}
										<code class="game-id">{gameId.slice(0, 8)}...</code>
									{/each}
								</div>
							{/if}
							{#if serverState.gamesInDbOnly.length > 0}
								<div class="discrepancy-section">
									<h4>In Database Only ({serverState.gamesInDbOnly.length})</h4>
									{#each serverState.gamesInDbOnly as gameId}
										<code class="game-id">{gameId.slice(0, 8)}...</code>
									{/each}
								</div>
							{/if}
						</div>
					{:else}
						<div class="stat-card success">
							<h3>✅ Sync Status</h3>
							<p>Memory and database are in sync!</p>
						</div>
					{/if}
				</div>
			{:else if activeTab === 'games' && activeGames}
				<!-- Active Games Tab -->
				<div class="games-summary">
					<span class="summary-item">
						<strong>{activeGames.totalInMemory}</strong> in memory
					</span>
					<span class="summary-item">
						<strong>{activeGames.totalInDatabase}</strong> in database
					</span>
				</div>

				{#if activeGames.games.length === 0}
					<div class="empty-state">No active games</div>
				{:else}
					<div class="games-table">
						<table>
							<thead>
								<tr>
									<th>Game ID</th>
									<th>Type</th>
									<th>Players</th>
									<th>Turn</th>
									<th>State</th>
									<th>Memory</th>
									<th>DB</th>
								</tr>
							</thead>
							<tbody>
								{#each activeGames.games as game (game.gameId)}
									<tr class:warning={!game.inMemory || !game.inDatabase}>
										<td>
											<code class="game-id">{game.gameId.slice(0, 8)}...</code>
										</td>
										<td>{game.gameType}</td>
										<td>{game.players.join(', ')}</td>
										<td>{game.turnNumber}</td>
										<td>
											<span class="state-badge state-{game.state.toLowerCase()}">
												{game.state}
											</span>
										</td>
										<td class:yes={game.inMemory} class:no={!game.inMemory}>
											{game.inMemory ? '✓' : '✗'}
										</td>
										<td class:yes={game.inDatabase} class:no={!game.inDatabase}>
											{game.inDatabase ? '✓' : '✗'}
										</td>
									</tr>
								{/each}
							</tbody>
						</table>
					</div>
				{/if}
			{:else if activeTab === 'history'}
				<!-- Match History Tab -->
				{#if matchHistory.length === 0}
					<div class="empty-state">No match history</div>
				{:else}
					<div class="history-list">
						{#each matchHistory as match (match.id)}
							<div class="history-card">
								<div class="history-header">
									<span class="game-type">{match.gameType}</span>
									<span class="duration">{formatDuration(match.durationSeconds)}</span>
								</div>
								<div class="history-players">
									{#each match.players as player}
										<span
											class="player-badge"
											class:winner={player.result === 'win'}
											class:loser={player.result === 'loss'}
										>
											{player.username}
											{#if player.result === 'win'}🏆{/if}
										</span>
									{/each}
								</div>
								{#if match.winnerName}
									<div class="winner-info">Winner: {match.winnerName}</div>
								{/if}
							</div>
						{/each}
					</div>
				{/if}
			{/if}
		</div>
	</div>
{/if}

<style>
	.debug-panel {
		background: linear-gradient(135deg, #1a1a2e 0%, #16213e 100%);
		border: 1px solid #0f3460;
		border-radius: 12px;
		overflow: hidden;
		font-family: 'JetBrains Mono', 'Fira Code', monospace;
		color: #e0e0e0;
		font-size: 0.875rem;
	}

	.debug-header {
		display: flex;
		justify-content: space-between;
		align-items: center;
		padding: 0.875rem 1rem;
		background: linear-gradient(135deg, #0f3460, #16213e);
		border-bottom: 1px solid #00d9ff33;
	}

	.header-left {
		display: flex;
		align-items: center;
		gap: 0.625rem;
	}

	.debug-icon {
		font-size: 1.25rem;
	}

	.debug-header h2 {
		margin: 0;
		font-size: 1rem;
		font-weight: 600;
		color: #00d9ff;
	}

	.header-right {
		display: flex;
		align-items: center;
		gap: 1rem;
	}

	.refresh-controls {
		display: flex;
		align-items: center;
		gap: 0.5rem;
	}

	.btn-refresh {
		padding: 0.375rem 0.75rem;
		background: #0f3460;
		border: 1px solid #00d9ff;
		border-radius: 6px;
		color: #00d9ff;
		font-family: inherit;
		font-size: 0.75rem;
		cursor: pointer;
		transition: all 0.2s;
	}

	.btn-refresh:hover:not(:disabled) {
		background: #00d9ff;
		color: #1a1a2e;
	}

	.btn-refresh:disabled {
		opacity: 0.5;
	}

	.auto-toggle {
		display: flex;
		align-items: center;
		gap: 0.375rem;
		font-size: 0.6875rem;
		color: #888;
		cursor: pointer;
	}

	.auto-toggle input {
		accent-color: #00d9ff;
	}

	.last-update {
		font-size: 0.6875rem;
		color: #666;
	}

	.btn-close {
		background: none;
		border: none;
		color: #666;
		font-size: 1rem;
		cursor: pointer;
		padding: 0.25rem;
		transition: color 0.2s;
	}

	.btn-close:hover {
		color: #ff5252;
	}

	.error-banner {
		display: flex;
		justify-content: space-between;
		align-items: center;
		padding: 0.625rem 1rem;
		background: rgba(255, 82, 82, 0.15);
		border-bottom: 1px solid rgba(255, 82, 82, 0.3);
		color: #ff5252;
		font-size: 0.8125rem;
	}

	.error-banner button {
		padding: 0.25rem 0.625rem;
		background: transparent;
		border: 1px solid #ff5252;
		border-radius: 4px;
		color: #ff5252;
		cursor: pointer;
		font-size: 0.75rem;
	}

	.tabs {
		display: flex;
		border-bottom: 1px solid #333;
	}

	.tab {
		flex: 1;
		padding: 0.625rem;
		background: transparent;
		border: none;
		color: #888;
		font-family: inherit;
		font-size: 0.75rem;
		cursor: pointer;
		transition: all 0.2s;
	}

	.tab:hover {
		color: #e0e0e0;
		background: rgba(255, 255, 255, 0.03);
	}

	.tab.active {
		color: #00d9ff;
		background: rgba(0, 217, 255, 0.1);
		border-bottom: 2px solid #00d9ff;
	}

	.tab-content {
		padding: 1rem;
		max-height: 400px;
		overflow-y: auto;
	}

	.loading,
	.empty-state {
		text-align: center;
		padding: 2rem;
		color: #666;
	}

	/* State Tab */
	.state-grid {
		display: grid;
		grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
		gap: 1rem;
	}

	.stat-card {
		background: rgba(255, 255, 255, 0.03);
		border: 1px solid #333;
		border-radius: 8px;
		padding: 1rem;
	}

	.stat-card h3 {
		margin: 0 0 0.75rem 0;
		font-size: 0.875rem;
		font-weight: 600;
		color: #fff;
	}

	.stat-card.memory {
		border-color: #00d9ff33;
	}

	.stat-card.database {
		border-color: #00ff8833;
	}

	.stat-card.warning {
		border-color: rgba(245, 158, 11, 0.5);
		background: rgba(245, 158, 11, 0.1);
	}

	.stat-card.success {
		border-color: rgba(46, 204, 113, 0.5);
		background: rgba(46, 204, 113, 0.1);
	}

	.stat-card.success p {
		margin: 0;
		color: #2ecc71;
	}

	.stat-row {
		display: flex;
		justify-content: space-between;
		margin-bottom: 0.5rem;
	}

	.stat-row .label {
		color: #888;
	}

	.stat-row .value {
		font-weight: 600;
		color: #00d9ff;
	}

	.discrepancy-section {
		margin-top: 0.5rem;
	}

	.discrepancy-section h4 {
		margin: 0 0 0.375rem 0;
		font-size: 0.75rem;
		color: #f59e0b;
	}

	/* Games Tab */
	.games-summary {
		display: flex;
		gap: 1.5rem;
		margin-bottom: 1rem;
		padding-bottom: 0.75rem;
		border-bottom: 1px solid #333;
	}

	.summary-item {
		font-size: 0.8125rem;
		color: #888;
	}

	.summary-item strong {
		color: #00d9ff;
	}

	.games-table {
		overflow-x: auto;
	}

	.games-table table {
		width: 100%;
		border-collapse: collapse;
	}

	.games-table th,
	.games-table td {
		padding: 0.5rem;
		text-align: left;
		border-bottom: 1px solid #333;
	}

	.games-table th {
		font-weight: 600;
		color: #888;
		font-size: 0.6875rem;
		text-transform: uppercase;
	}

	.games-table tr.warning {
		background: rgba(245, 158, 11, 0.1);
	}

	.game-id {
		font-size: 0.75rem;
		color: #00d9ff;
		background: rgba(0, 217, 255, 0.1);
		padding: 0.125rem 0.375rem;
		border-radius: 4px;
	}

	.state-badge {
		padding: 0.125rem 0.375rem;
		border-radius: 4px;
		font-size: 0.6875rem;
		text-transform: uppercase;
	}

	.state-in_progress {
		background: rgba(0, 255, 136, 0.2);
		color: #00ff88;
	}

	.state-mulligan {
		background: rgba(245, 158, 11, 0.2);
		color: #f59e0b;
	}

	.state-starting {
		background: rgba(0, 217, 255, 0.2);
		color: #00d9ff;
	}

	.state-finished {
		background: rgba(255, 255, 255, 0.1);
		color: #888;
	}

	.yes {
		color: #2ecc71;
	}

	.no {
		color: #ff5252;
	}

	/* History Tab */
	.history-list {
		display: flex;
		flex-direction: column;
		gap: 0.75rem;
	}

	.history-card {
		background: rgba(255, 255, 255, 0.03);
		border: 1px solid #333;
		border-radius: 8px;
		padding: 0.75rem;
	}

	.history-header {
		display: flex;
		justify-content: space-between;
		margin-bottom: 0.5rem;
	}

	.game-type {
		font-weight: 600;
		color: #fff;
	}

	.duration {
		font-size: 0.75rem;
		color: #888;
	}

	.history-players {
		display: flex;
		gap: 0.5rem;
		flex-wrap: wrap;
	}

	.player-badge {
		padding: 0.25rem 0.5rem;
		background: rgba(255, 255, 255, 0.1);
		border-radius: 4px;
		font-size: 0.75rem;
	}

	.player-badge.winner {
		background: rgba(46, 204, 113, 0.2);
		color: #2ecc71;
	}

	.player-badge.loser {
		background: rgba(255, 82, 82, 0.1);
		color: #ff5252;
	}

	.winner-info {
		margin-top: 0.5rem;
		font-size: 0.75rem;
		color: #f59e0b;
	}
</style>


