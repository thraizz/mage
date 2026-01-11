<script lang="ts">
	import type { Deck } from '$lib/types/deck';
	import FormatBadge from '$lib/components/mtg/FormatBadge.svelte';
	import Badge from '$lib/components/ui/Badge.svelte';
	import Panel from '$lib/components/ui/Panel.svelte';
	import ChevronRight from '@lucide/svelte/icons/chevron-right';

	interface Props {
		deck: Deck;
		onclick?: (e: MouseEvent) => void;
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
</script>

<button class="deck-card" type="button" {onclick} aria-label={`Open deck: ${deck.name}`}>
	<Panel variant="bordered" padding="md">
		<div class="deck-card-inner" data-invalid={!deck.isValid}>
			<div class="deck-card-shine" aria-hidden="true"></div>

			<div class="deck-card-content">
				<div class="deck-header">
					<h3 class="deck-name" title={deck.name}>{deck.name}</h3>
					<ChevronRight class="chevron" size={18} aria-hidden="true" />
				</div>

				<div class="deck-badges">
					<FormatBadge format={deck.format} size="md" />
					<Badge variant={deck.isValid ? 'success' : 'error'} size="sm">
						{deck.isValid ? 'Valid' : 'Fix required'}
					</Badge>
					{#if deck.commanders && deck.commanders.length > 0}
						<Badge variant="info" size="sm">Commander</Badge>
					{/if}
				</div>

				<div class="deck-stats" aria-label="Deck stats">
					<div class="stat">
						<div class="stat-value">{deck.cardCount}</div>
						<div class="stat-label">Total</div>
					</div>
					<div class="stat">
						<div class="stat-value">{deck.mainDeck.length}</div>
						<div class="stat-label">Main</div>
					</div>
					{#if deck.sideboard && deck.sideboard.length > 0}
						<div class="stat">
							<div class="stat-value">{deck.sideboard.length}</div>
							<div class="stat-label">Side</div>
						</div>
					{/if}
					{#if deck.commanders && deck.commanders.length > 0}
						<div class="stat">
							<div class="stat-value">{deck.commanders.length}</div>
							<div class="stat-label">Cmd</div>
						</div>
					{/if}
				</div>

				<div class="deck-footer">
					<span class="modified-date">Updated {formatDate(deck.updatedAt)}</span>
					{#if !deck.isValid}
						<span class="invalid-hint">Invalid cards</span>
					{/if}
				</div>
			</div>
		</div>
	</Panel>
</button>

<style>
	.deck-card {
		width: 100%;
		display: block;
		text-align: left;
		background: transparent;
		border: none;
		padding: 0;
		cursor: pointer;
	}

	.deck-card :global(.panel) {
		transition: transform var(--transition-base), box-shadow var(--transition-base);
	}

	.deck-card:hover :global(.panel) {
		transform: translateY(-2px);
		box-shadow: var(--shadow-glow);
	}

	.deck-card:active :global(.panel) {
		transform: translateY(-1px);
	}

	.deck-card:focus-visible {
		outline: none;
	}

	.deck-card:focus-visible :global(.panel) {
		box-shadow: 0 0 0 3px var(--accent-gold-glow), var(--shadow-glow);
	}

	.deck-card-inner {
		position: relative;
	}

	.deck-card-shine {
		position: absolute;
		inset: 0;
		pointer-events: none;
		background: linear-gradient(
			135deg,
			color-mix(in srgb, var(--accent-gold) 10%, transparent) 0%,
			transparent 45%,
			color-mix(in srgb, var(--accent-gold) 6%, transparent) 100%
		);
		opacity: 0;
		transition: opacity var(--transition-fast);
	}

	.deck-card:hover .deck-card-shine {
		opacity: 1;
	}

	.deck-card-content {
		display: flex;
		flex-direction: column;
		gap: var(--space-3);
	}

	.deck-header {
		display: flex;
		align-items: flex-start;
		justify-content: space-between;
		gap: var(--space-3);
	}

	.deck-name {
		font-family: var(--font-display);
		font-size: var(--text-xl);
		font-weight: var(--weight-semibold);
		color: var(--text-bright);
		margin: 0;
		overflow: hidden;
		text-overflow: ellipsis;
		display: -webkit-box;
		-webkit-line-clamp: 2;
		-webkit-box-orient: vertical;
		line-clamp: 2;
	}

	:global(svg.chevron) {
		color: var(--text-ghost);
		flex: 0 0 auto;
		margin-top: 2px;
		transition: transform var(--transition-fast), color var(--transition-fast);
	}

	.deck-card:hover :global(svg.chevron) {
		transform: translateX(2px);
		color: var(--text-dim);
	}

	.deck-badges {
		display: flex;
		align-items: center;
		flex-wrap: wrap;
		gap: var(--space-2);
	}

	.deck-stats {
		display: grid;
		grid-template-columns: repeat(4, minmax(0, 1fr));
		gap: var(--space-3);
		padding: var(--space-3);
		border-radius: var(--radius-md);
		background: color-mix(in srgb, var(--bg-slate) 65%, transparent);
		border: 1px solid var(--border-subtle);
	}

	.stat {
		display: flex;
		flex-direction: column;
		align-items: flex-start;
		gap: 2px;
	}

	.stat-value {
		font-size: var(--text-lg);
		font-weight: var(--weight-semibold);
		color: var(--text-bright);
		letter-spacing: -0.01em;
	}

	.stat-label {
		font-size: var(--text-xs);
		color: var(--text-dim);
		text-transform: uppercase;
		letter-spacing: 0.05em;
	}

	.deck-footer {
		display: flex;
		align-items: center;
		justify-content: space-between;
		gap: var(--space-3);
	}

	.modified-date {
		font-size: var(--text-xs);
		color: var(--text-ghost);
	}

	.invalid-hint {
		font-size: var(--text-xs);
		font-weight: var(--weight-medium);
		color: var(--status-error);
		background: var(--status-error-dim);
		padding: 0.125rem var(--space-2);
		border-radius: var(--radius-full);
		white-space: nowrap;
	}

	@media (max-width: 480px) {
		.deck-stats {
			grid-template-columns: repeat(2, minmax(0, 1fr));
		}
	}
</style>
