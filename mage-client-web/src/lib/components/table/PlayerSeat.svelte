<script lang="ts">
	import Badge from '$lib/components/ui/Badge.svelte';
	import Button from '$lib/components/ui/Button.svelte';

	interface Player {
		id: string;
		username: string;
		isHost: boolean;
		isReady?: boolean;
	}

	interface Props {
		player?: Player;
		seatNumber: number;
		isCurrentUser?: boolean;
		canKick?: boolean;
		onkick?: (playerId: string) => void;
	}

	let { player, seatNumber, isCurrentUser = false, canKick = false, onkick }: Props = $props();

	const isEmpty = $derived(!player);
</script>

<div class="player-seat" class:empty={isEmpty} class:current={isCurrentUser}>
	<div class="seat-number">Seat {seatNumber}</div>

	{#if player}
		<div class="avatar">
			<svg width="32" height="32" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
				<path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"></path>
				<circle cx="12" cy="7" r="4"></circle>
			</svg>
		</div>

		<div class="player-info">
			<span class="player-name">
				{player.username}
				{#if isCurrentUser}
					<Badge variant="info" size="sm">You</Badge>
				{/if}
			</span>

			<div class="player-badges">
				{#if player.isHost}
					<Badge variant="warning" size="sm">Host</Badge>
				{/if}
			</div>
		</div>

		{#if canKick && !player.isHost && onkick}
			<Button variant="danger" size="sm" onclick={() => onkick(player.id)}>Kick</Button>
		{/if}
	{:else}
		<div class="empty-avatar">
			<svg width="32" height="32" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" opacity="0.3">
				<path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"></path>
				<circle cx="12" cy="7" r="4"></circle>
			</svg>
		</div>
		<span class="waiting-text">Waiting...</span>
	{/if}
</div>

<style>
	.player-seat {
		display: flex;
		flex-direction: column;
		align-items: center;
		gap: var(--space-3);
		padding: var(--space-5);
		background: var(--bg-obsidian);
		border: 2px solid var(--border-subtle);
		border-radius: var(--radius-lg);
		min-height: 180px;
		transition: all var(--transition-base);
	}

	.player-seat:not(.empty) {
		border-color: var(--border-default);
	}

	.player-seat.current {
		border-color: var(--accent-gold);
		box-shadow: 0 0 15px var(--accent-gold-glow);
	}

	.player-seat.empty {
		border-style: dashed;
		background: var(--bg-slate);
	}

	.seat-number {
		font-size: var(--text-xs);
		font-weight: var(--weight-semibold);
		color: var(--text-ghost);
		text-transform: uppercase;
		letter-spacing: 0.05em;
	}

	.avatar {
		width: 4rem;
		height: 4rem;
		border-radius: var(--radius-full);
		background: var(--accent-gold-dim);
		display: flex;
		align-items: center;
		justify-content: center;
		color: var(--bg-void);
	}

	.current .avatar {
		background: var(--accent-gold);
	}

	.empty-avatar {
		width: 4rem;
		height: 4rem;
		display: flex;
		align-items: center;
		justify-content: center;
		color: var(--text-ghost);
	}

	.player-info {
		display: flex;
		flex-direction: column;
		align-items: center;
		gap: var(--space-2);
		text-align: center;
	}

	.player-name {
		display: flex;
		align-items: center;
		gap: var(--space-2);
		font-size: var(--text-base);
		font-weight: var(--weight-semibold);
		color: var(--text-bright);
	}

	.player-badges {
		display: flex;
		gap: var(--space-1);
	}

	.waiting-text {
		font-size: var(--text-sm);
		color: var(--text-ghost);
		font-style: italic;
	}
</style>
