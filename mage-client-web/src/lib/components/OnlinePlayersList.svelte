<script lang="ts">
	import type { OnlinePlayer } from '$lib/types/player';

	// Props
	let {
		players = $bindable([]),
		isOpen = $bindable(false),
		currentUsername
	}: {
		players?: OnlinePlayer[];
		isOpen?: boolean;
		currentUsername: string;
	} = $props();

	// Computed
	const onlineCount = $derived(players.length);

	function toggleOpen() {
		isOpen = !isOpen;
	}
</script>

<div class="online-players">
	<!-- Header with count and toggle -->
	<button class="players-header" onclick={toggleOpen} aria-expanded={isOpen}>
		<div class="header-content">
			<div class="header-icon">
				<svg
					xmlns="http://www.w3.org/2000/svg"
					width="18"
					height="18"
					viewBox="0 0 24 24"
					fill="none"
					stroke="currentColor"
					stroke-width="2"
					stroke-linecap="round"
					stroke-linejoin="round"
				>
					<path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2"></path>
					<circle cx="9" cy="7" r="4"></circle>
					<path d="M23 21v-2a4 4 0 0 0-3-3.87"></path>
					<path d="M16 3.13a4 4 0 0 1 0 7.75"></path>
				</svg>
			</div>
			<div class="header-text">
				<span class="header-title">Online Players</span>
				<span class="player-count">{onlineCount}</span>
			</div>
		</div>
		<div class="toggle-icon" class:open={isOpen}>
			<svg
				xmlns="http://www.w3.org/2000/svg"
				width="16"
				height="16"
				viewBox="0 0 24 24"
				fill="none"
				stroke="currentColor"
				stroke-width="2"
				stroke-linecap="round"
				stroke-linejoin="round"
			>
				<polyline points="6 9 12 15 18 9"></polyline>
			</svg>
		</div>
	</button>

	<!-- Collapsible player list -->
	{#if isOpen}
		<div class="players-list" role="list">
			{#if players.length === 0}
				<div class="empty-state">
					<p>No players online</p>
				</div>
			{:else}
				{#each players as player (player.id)}
					<div class="player-item" role="listitem">
						<div class="player-status">
							<div class="status-dot"></div>
						</div>
						<div class="player-info">
							<span class="player-name">
								{player.username}
								{#if player.username === currentUsername}
									<span class="you-badge">You</span>
								{/if}
							</span>
						</div>
					</div>
				{/each}
			{/if}
		</div>
	{/if}
</div>

<style>
	.online-players {
		background: white;
		border-radius: 0.75rem;
		box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
		overflow: hidden;
	}

	/* Header */
	.players-header {
		width: 100%;
		padding: 1rem;
		display: flex;
		justify-content: space-between;
		align-items: center;
		background: white;
		border: none;
		cursor: pointer;
		transition: background-color 0.2s;
	}

	.players-header:hover {
		background-color: #f9fafb;
	}

	.header-content {
		display: flex;
		align-items: center;
		gap: 0.75rem;
	}

	.header-icon {
		display: flex;
		align-items: center;
		justify-content: center;
		width: 2.5rem;
		height: 2.5rem;
		background-color: #667eea;
		color: white;
		border-radius: 0.5rem;
	}

	.header-text {
		display: flex;
		flex-direction: column;
		align-items: flex-start;
		gap: 0.125rem;
	}

	.header-title {
		font-size: 0.875rem;
		font-weight: 600;
		color: #374151;
	}

	.player-count {
		font-size: 1.25rem;
		font-weight: 700;
		color: #667eea;
	}

	.toggle-icon {
		display: flex;
		align-items: center;
		justify-content: center;
		color: #9ca3af;
		transition: transform 0.2s;
	}

	.toggle-icon.open {
		transform: rotate(180deg);
	}

	/* Player List */
	.players-list {
		max-height: 400px;
		overflow-y: auto;
		border-top: 1px solid #e5e7eb;
	}

	.players-list::-webkit-scrollbar {
		width: 8px;
	}

	.players-list::-webkit-scrollbar-track {
		background: #f3f4f6;
	}

	.players-list::-webkit-scrollbar-thumb {
		background: #d1d5db;
		border-radius: 4px;
	}

	.players-list::-webkit-scrollbar-thumb:hover {
		background: #9ca3af;
	}

	/* Player Item */
	.player-item {
		display: flex;
		align-items: center;
		gap: 0.75rem;
		padding: 0.75rem 1rem;
		transition: background-color 0.15s;
	}

	.player-item:hover {
		background-color: #f9fafb;
	}

	.player-item:not(:last-child) {
		border-bottom: 1px solid #f3f4f6;
	}

	.player-status {
		flex-shrink: 0;
	}

	.status-dot {
		width: 0.625rem;
		height: 0.625rem;
		background-color: #10b981;
		border-radius: 50%;
		box-shadow: 0 0 0 2px rgba(16, 185, 129, 0.2);
	}

	.player-info {
		flex: 1;
		min-width: 0;
	}

	.player-name {
		font-size: 0.875rem;
		font-weight: 500;
		color: #111827;
		display: flex;
		align-items: center;
		gap: 0.5rem;
	}

	.you-badge {
		display: inline-flex;
		align-items: center;
		padding: 0.125rem 0.5rem;
		background-color: #667eea;
		color: white;
		font-size: 0.6875rem;
		font-weight: 600;
		border-radius: 0.25rem;
		text-transform: uppercase;
	}

	/* Empty State */
	.empty-state {
		padding: 2rem 1rem;
		text-align: center;
	}

	.empty-state p {
		color: #6b7280;
		font-size: 0.875rem;
		margin: 0;
	}

	/* Responsive */
	@media (max-width: 640px) {
		.players-header {
			padding: 0.875rem;
		}

		.player-item {
			padding: 0.625rem 0.875rem;
		}

		.players-list {
			max-height: 300px;
		}
	}
</style>
