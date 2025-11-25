<script lang="ts">
	import type { Deck } from '$lib/types/deck';
	import FormatBadge from '$lib/components/mtg/FormatBadge.svelte';
	import Badge from '$lib/components/ui/Badge.svelte';

	interface Props {
		deck: Deck;
		onclick?: () => void;
	}

	let { deck, onclick }: Props = $props();

	function formatDate(timestamp: number): string {
		const date = new Date(timestamp);
		const now = new Date();
		const diffMs = now.getTime() - date.getTime();
		const diffDays = Math.floor(diffMs / (1000 * 60 * 60 * 24));

		if (diffDays === 0) return 'Today';
		if (diffDays === 1) return 'Yesterday';
		if (diffDays < 7) return `${diffDays}d ago`;
		if (diffDays < 30) return `${Math.floor(diffDays / 7)}w ago`;
		return date.toLocaleDateString();
	}

	function handleKeypress(e: KeyboardEvent) {
		if (e.key === 'Enter' || e.key === ' ') {
			e.preventDefault();
			onclick?.();
		}
	}
</script>

<div
	class="deck-card"
	onclick={onclick}
	onkeypress={handleKeypress}
	role="button"
	tabindex="0"
>
	<header class="card-header">
		<FormatBadge format={deck.format} size="md" />
		<Badge variant={deck.isValid ? 'success' : 'error'} size="sm">
			{deck.isValid ? 'Valid' : 'Invalid'}
		</Badge>
	</header>

	<h3 class="deck-name" title={deck.name}>{deck.name}</h3>

	<div class="deck-stats">
		<div class="stat">
			<span class="stat-value">{deck.cardCount}</span>
			<span class="stat-label">cards</span>
		</div>
		<div class="stat">
			<span class="stat-value">{deck.mainDeck.length}</span>
			<span class="stat-label">main</span>
		</div>
		{#if deck.sideboard && deck.sideboard.length > 0}
			<div class="stat">
				<span class="stat-value">{deck.sideboard.length}</span>
				<span class="stat-label">side</span>
			</div>
		{/if}
	</div>

	<footer class="card-footer">
		<span class="modified-date">{formatDate(deck.updatedAt)}</span>
	</footer>
</div>

<style>
	.deck-card {
		background: var(--bg-obsidian);
		border: 1px solid var(--border-subtle);
		border-radius: var(--radius-lg);
		padding: var(--space-4);
		cursor: pointer;
		transition: all var(--transition-base);
	}

	.deck-card:hover {
		border-color: var(--accent-gold-dim);
		box-shadow: var(--shadow-glow);
		transform: translateY(-2px);
	}

	.deck-card:focus-visible {
		outline: 2px solid var(--accent-gold);
		outline-offset: 2px;
	}

	/* Header */
	.card-header {
		display: flex;
		justify-content: space-between;
		align-items: center;
		margin-bottom: var(--space-3);
	}

	/* Deck Name */
	.deck-name {
		font-family: var(--font-display);
		font-size: var(--text-lg);
		font-weight: var(--weight-semibold);
		color: var(--text-bright);
		margin: 0 0 var(--space-4);
		overflow: hidden;
		text-overflow: ellipsis;
		white-space: nowrap;
	}

	/* Stats */
	.deck-stats {
		display: flex;
		gap: var(--space-4);
		margin-bottom: var(--space-4);
	}

	.stat {
		display: flex;
		flex-direction: column;
		align-items: center;
	}

	.stat-value {
		font-size: var(--text-xl);
		font-weight: var(--weight-bold);
		color: var(--accent-gold);
	}

	.stat-label {
		font-size: var(--text-xs);
		color: var(--text-dim);
		text-transform: uppercase;
		letter-spacing: 0.05em;
	}

	/* Footer */
	.card-footer {
		padding-top: var(--space-3);
		border-top: 1px solid var(--border-subtle);
	}

	.modified-date {
		font-size: var(--text-xs);
		color: var(--text-ghost);
	}
</style>
