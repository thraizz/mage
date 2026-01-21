<script lang="ts">
	import type { OnlinePlayer } from '$lib/types/player';
	import Users from '@lucide/svelte/icons/users';
	import ChevronDown from '@lucide/svelte/icons/chevron-down';

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
				<Users size={18} aria-hidden="true" />
			</div>
			<div class="header-text">
				<span class="header-title">Online Players</span>
				<span class="player-count">{onlineCount}</span>
			</div>
		</div>
		<div class="toggle-icon" class:open={isOpen}>
			<ChevronDown size={16} aria-hidden="true" />
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
		background: var(--bg-obsidian);
		border: 1px solid var(--border-subtle);
		border-radius: var(--radius-lg);
		overflow: hidden;
	}

	/* Header */
	.players-header {
		width: 100%;
		padding: var(--space-4);
		display: flex;
		justify-content: space-between;
		align-items: center;
		background: var(--bg-obsidian);
		border: none;
		cursor: pointer;
		transition: background-color var(--transition-fast);
	}

	.players-header:hover {
		background-color: var(--bg-slate);
	}

	.header-content {
		display: flex;
		align-items: center;
		gap: var(--space-3);
	}

	.header-icon {
		display: flex;
		align-items: center;
		justify-content: center;
		width: 2.5rem;
		height: 2.5rem;
		background: linear-gradient(135deg, var(--ci-jace-cloak) 0%, #2563eb 100%);
		color: var(--ci-scroll-parchment);
		border-radius: var(--radius-md);
		box-shadow: 0 2px 8px rgba(59, 130, 246, 0.3);
	}

	.header-text {
		display: flex;
		flex-direction: column;
		align-items: flex-start;
		gap: 0.125rem;
	}

	.header-title {
		font-size: var(--text-sm);
		font-weight: var(--weight-semibold);
		color: var(--ci-scroll-parchment);
		text-transform: uppercase;
		letter-spacing: 0.05em;
	}

	.player-count {
		font-size: var(--text-xl);
		font-weight: var(--weight-bold);
		color: var(--ci-jace-cloak);
	}

	.toggle-icon {
		display: flex;
		align-items: center;
		justify-content: center;
		color: var(--ci-swamp-obsidian);
		transition: transform var(--transition-base);
	}

	.toggle-icon.open {
		transform: rotate(180deg);
	}

	/* Player List */
	.players-list {
		max-height: 400px;
		overflow-y: auto;
		border-top: 1px solid var(--border-subtle);
	}

	/* Player Item */
	.player-item {
		display: flex;
		align-items: center;
		gap: var(--space-3);
		padding: var(--space-3) var(--space-4);
		transition: background-color var(--transition-fast);
	}

	.player-item:hover {
		background-color: var(--bg-slate);
	}

	.player-item:not(:last-child) {
		border-bottom: 1px solid var(--border-subtle);
	}

	.player-status {
		flex-shrink: 0;
	}

	.status-dot {
		width: 0.5rem;
		height: 0.5rem;
		background-color: var(--ci-forest-emerald);
		border-radius: var(--radius-full);
		box-shadow: 0 0 8px rgba(46, 204, 113, 0.5);
		animation: pulse-glow 2s ease-in-out infinite;
	}

	@keyframes pulse-glow {
		0%,
		100% {
			box-shadow: 0 0 4px rgba(46, 204, 113, 0.4);
		}
		50% {
			box-shadow: 0 0 12px rgba(46, 204, 113, 0.6);
		}
	}

	.player-info {
		flex: 1;
		min-width: 0;
	}

	.player-name {
		font-size: var(--text-sm);
		font-weight: var(--weight-medium);
		color: var(--ci-scroll-parchment);
		display: flex;
		align-items: center;
		gap: var(--space-2);
	}

	.you-badge {
		display: inline-flex;
		align-items: center;
		padding: 0.125rem var(--space-2);
		background: linear-gradient(135deg, var(--ci-jace-cloak) 0%, #2563eb 100%);
		color: var(--ci-scroll-parchment);
		font-size: var(--text-xs);
		font-weight: var(--weight-semibold);
		border-radius: var(--radius-sm);
		text-transform: uppercase;
		letter-spacing: 0.05em;
	}

	/* Empty State */
	.empty-state {
		padding: var(--space-8) var(--space-4);
		text-align: center;
	}

	.empty-state p {
		color: var(--ci-swamp-obsidian);
		font-size: var(--text-sm);
		font-style: italic;
		margin: 0;
	}

	/* Responsive */
	@media (max-width: 640px) {
		.players-header {
			padding: var(--space-3);
		}

		.player-item {
			padding: var(--space-2) var(--space-3);
		}

		.players-list {
			max-height: 300px;
		}
	}
</style>
