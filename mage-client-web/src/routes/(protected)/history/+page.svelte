<script lang="ts">
	import { onMount } from 'svelte';
	import { fetchMatchHistory, formatDuration, formatRelativeTime } from '$lib/api/match_history';
	import type { Match } from '$lib/api/match_history';

	let matches: Match[] = [];
	let loading = true;
	let error: string | null = null;
	let totalCount = 0;

	// Pagination
	let currentPage = 1;
	let pageSize = 20;

	$: totalPages = Math.ceil(totalCount / pageSize);
	$: offset = (currentPage - 1) * pageSize;

	// Load matches on mount
	onMount(async () => {
		await loadMatches();
	});

	async function loadMatches() {
		loading = true;
		error = null;
		try {
			const result = await fetchMatchHistory(pageSize, offset);
			matches = result.matches;
			totalCount = result.totalCount;
		} catch (err) {
			error = err instanceof Error ? err.message : 'Failed to load match history';
			console.error('Failed to load match history:', err);
		} finally {
			loading = false;
		}
	}

	async function goToPage(page: number) {
		if (page < 1 || page > totalPages) return;
		currentPage = page;
		await loadMatches();
	}

	function getResultBadgeClass(result: string): string {
		switch (result) {
			case 'win':
				return 'badge-win';
			case 'loss':
				return 'badge-loss';
			case 'draw':
				return 'badge-draw';
			case 'concede':
				return 'badge-concede';
			default:
				return 'badge-default';
		}
	}

	function getResultEmoji(result: string): string {
		switch (result) {
			case 'win':
				return '🏆';
			case 'loss':
				return '💀';
			case 'draw':
				return '🤝';
			case 'concede':
				return '🏳️';
			default:
				return '❓';
		}
	}
</script>

<svelte:head>
	<title>Match History - MAGE</title>
</svelte:head>

<div class="container">
	<header>
		<h1>Match History</h1>
		<div class="stats">
			<span class="stat">Total Matches: {totalCount}</span>
		</div>
	</header>

	{#if loading}
		<div class="loading">
			<p>Loading match history...</p>
		</div>
	{:else if error}
		<div class="error">
			<p>❌ {error}</p>
			<button class="btn-primary" onclick={loadMatches}>Retry</button>
		</div>
	{:else if matches.length === 0}
		<div class="empty-state">
			<p>🎮 No matches found</p>
			<p class="hint">Play your first game to see match history here!</p>
		</div>
	{:else}
		<div class="matches-list">
			{#each matches as match (match.id)}
				<div class="match-card">
					<div class="match-header">
						<div class="match-info">
							<h3>{match.gameType}</h3>
							<span class="match-time">{formatRelativeTime(match.endTime)}</span>
						</div>
						<div class="match-duration">
							<span class="duration-badge">{formatDuration(match.durationSeconds)}</span>
						</div>
					</div>

					<div class="players-section">
						{#each match.players as player}
							<div class="player-row">
								<div class="player-info">
									<span class="player-name">{player.username}</span>
									<span class="player-deck">{player.deck || 'Unknown Deck'}</span>
								</div>
								<div class="player-result">
									<span class={`result-badge ${getResultBadgeClass(player.result)}`}>
										{getResultEmoji(player.result)}
										{player.result.toUpperCase()}
									</span>
								</div>
							</div>
						{/each}
					</div>

					{#if match.winnerName}
						<div class="match-winner">
							<span>Winner: <strong>{match.winnerName}</strong></span>
						</div>
					{/if}

					<div class="match-meta">
						{#if match.tableId}
							<span class="meta-item">Table: {match.tableId}</span>
						{/if}
						{#if match.tournamentId}
							<span class="meta-item">Tournament: {match.tournamentId}</span>
						{/if}
					</div>
				</div>
			{/each}
		</div>

		{#if totalPages > 1}
			<div class="pagination">
				<button
					class="btn-pagination"
					disabled={currentPage === 1}
					onclick={() => goToPage(currentPage - 1)}
				>
					Previous
				</button>

				<div class="page-numbers">
					{#if currentPage > 2}
						<button class="btn-page" onclick={() => goToPage(1)}>1</button>
						{#if currentPage > 3}
							<span class="ellipsis">...</span>
						{/if}
					{/if}

					{#if currentPage > 1}
						<button class="btn-page" onclick={() => goToPage(currentPage - 1)}>
							{currentPage - 1}
						</button>
					{/if}

					<button class="btn-page active">{currentPage}</button>

					{#if currentPage < totalPages}
						<button class="btn-page" onclick={() => goToPage(currentPage + 1)}>
							{currentPage + 1}
						</button>
					{/if}

					{#if currentPage < totalPages - 1}
						{#if currentPage < totalPages - 2}
							<span class="ellipsis">...</span>
						{/if}
						<button class="btn-page" onclick={() => goToPage(totalPages)}>
							{totalPages}
						</button>
					{/if}
				</div>

				<button
					class="btn-pagination"
					disabled={currentPage === totalPages}
					onclick={() => goToPage(currentPage + 1)}
				>
					Next
				</button>
			</div>
		{/if}
	{/if}
</div>

<style>
	.container {
		max-width: 1200px;
		margin: 0 auto;
		padding: 2rem;
	}

	header {
		display: flex;
		justify-content: space-between;
		align-items: center;
		margin-bottom: 2rem;
	}

	h1 {
		margin: 0;
		font-size: 2.5rem;
		color: #333;
	}

	.stats {
		display: flex;
		gap: 2rem;
	}

	.stat {
		font-size: 1rem;
		color: #666;
		font-weight: 500;
	}

	.loading,
	.error {
		text-align: center;
		padding: 3rem;
		background: #f9fafb;
		border-radius: 8px;
	}

	.error {
		border: 2px solid #ef4444;
	}

	.error p {
		color: #ef4444;
		font-size: 1.125rem;
		margin: 0 0 1rem 0;
	}

	.loading p {
		color: #666;
		font-size: 1.125rem;
		margin: 0;
	}

	.matches-list {
		display: flex;
		flex-direction: column;
		gap: 1rem;
	}

	.match-card {
		background: white;
		border: 1px solid #ddd;
		border-radius: 8px;
		padding: 1.5rem;
		transition: box-shadow 0.2s;
	}

	.match-card:hover {
		box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
	}

	.match-header {
		display: flex;
		justify-content: space-between;
		align-items: center;
		margin-bottom: 1rem;
		padding-bottom: 1rem;
		border-bottom: 1px solid #eee;
	}

	.match-info h3 {
		margin: 0 0 0.25rem 0;
		font-size: 1.25rem;
		color: #333;
	}

	.match-time {
		font-size: 0.875rem;
		color: #888;
	}

	.duration-badge {
		padding: 0.25rem 0.75rem;
		background: #f3f4f6;
		color: #555;
		border-radius: 12px;
		font-size: 0.875rem;
		font-weight: 500;
	}

	.players-section {
		display: flex;
		flex-direction: column;
		gap: 0.75rem;
		margin-bottom: 1rem;
	}

	.player-row {
		display: flex;
		justify-content: space-between;
		align-items: center;
		padding: 0.75rem;
		background: #f9fafb;
		border-radius: 6px;
	}

	.player-info {
		display: flex;
		flex-direction: column;
		gap: 0.25rem;
	}

	.player-name {
		font-weight: 600;
		color: #333;
	}

	.player-deck {
		font-size: 0.875rem;
		color: #666;
	}

	.result-badge {
		padding: 0.375rem 0.75rem;
		border-radius: 12px;
		font-size: 0.75rem;
		font-weight: 600;
	}

	.badge-win {
		background: #dcfce7;
		color: #16a34a;
	}

	.badge-loss {
		background: #fee2e2;
		color: #dc2626;
	}

	.badge-draw {
		background: #e0e7ff;
		color: #4f46e5;
	}

	.badge-concede {
		background: #f3f4f6;
		color: #6b7280;
	}

	.match-winner {
		padding: 0.75rem;
		background: #fef3c7;
		border-radius: 6px;
		text-align: center;
		margin-bottom: 0.75rem;
		color: #92400e;
	}

	.match-winner strong {
		font-weight: 600;
	}

	.match-meta {
		display: flex;
		gap: 1rem;
		font-size: 0.875rem;
		color: #888;
	}

	.meta-item {
		padding: 0.25rem 0.5rem;
		background: #f3f4f6;
		border-radius: 4px;
	}

	.pagination {
		display: flex;
		justify-content: center;
		align-items: center;
		gap: 0.5rem;
		margin-top: 2rem;
		padding: 1rem;
	}

	.btn-pagination {
		padding: 0.5rem 1rem;
		background: white;
		color: #667eea;
		border: 1px solid #667eea;
		border-radius: 4px;
		font-size: 0.875rem;
		font-weight: 500;
		cursor: pointer;
		transition: all 0.2s;
	}

	.btn-pagination:hover:not(:disabled) {
		background: #667eea;
		color: white;
	}

	.btn-pagination:disabled {
		opacity: 0.5;
		cursor: not-allowed;
	}

	.page-numbers {
		display: flex;
		gap: 0.25rem;
		align-items: center;
	}

	.btn-page {
		padding: 0.5rem 0.75rem;
		background: white;
		color: #667eea;
		border: 1px solid #ddd;
		border-radius: 4px;
		font-size: 0.875rem;
		font-weight: 500;
		cursor: pointer;
		transition: all 0.2s;
		min-width: 2.5rem;
	}

	.btn-page:hover {
		background: #f3f4f6;
	}

	.btn-page.active {
		background: #667eea;
		color: white;
		border-color: #667eea;
	}

	.ellipsis {
		padding: 0 0.5rem;
		color: #999;
	}

	.btn-primary {
		padding: 0.75rem 1.5rem;
		background: #667eea;
		color: white;
		border: none;
		border-radius: 4px;
		font-size: 1rem;
		font-weight: 500;
		cursor: pointer;
		transition: background 0.2s;
	}

	.btn-primary:hover {
		background: #5568d3;
	}

	.empty-state {
		text-align: center;
		padding: 3rem;
		background: #f9fafb;
		border: 2px dashed #ddd;
		border-radius: 8px;
	}

	.empty-state p {
		margin: 0 0 1rem 0;
		font-size: 1.125rem;
		color: #666;
	}

	.empty-state .hint {
		font-size: 0.875rem;
		color: #888;
	}
</style>
