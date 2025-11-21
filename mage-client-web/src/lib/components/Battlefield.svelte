<script lang="ts">
	import type { Card as CardType } from '../types';
	import Card from './Card.svelte';
	import { game } from '../stores/game';

	interface Props {
		cards: CardType[];
		title: string;
		onCardClick?: (card: CardType) => void;
	}

	let { cards, title, onCardClick }: Props = $props();
</script>

<div class="battlefield">
	<h2>{title} ({cards.length} creatures)</h2>
	
	<div class="cards-grid">
		{#each cards as card (card.id)}
			<Card {card} onclick={() => onCardClick?.(card)} />
		{/each}
	</div>
	
	{#if cards.length === 0}
		<div class="empty">No creatures on battlefield</div>
	{/if}
</div>

<style>
	.battlefield {
		padding: 16px;
		background: rgba(0, 0, 0, 0.1);
		border-radius: 8px;
		min-height: 250px;
	}

	h2 {
		margin: 0 0 16px 0;
		font-size: 18px;
		color: #333;
	}

	.cards-grid {
		display: grid;
		grid-template-columns: repeat(auto-fill, minmax(150px, 1fr));
		gap: 16px;
	}

	.empty {
		text-align: center;
		padding: 64px;
		color: #999;
		font-style: italic;
	}
</style>
