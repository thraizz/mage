<script lang="ts">
	import Card from './Card.svelte';
	import type { GameCard } from '$lib/types/game';

	// Props
	let {
		cards = [],
		selectedCardIds = [],
		onCardClick = (cardId: string) => {},
		onCardHover = (cardId: string) => {},
		size = 'normal'
	}: {
		cards?: GameCard[];
		selectedCardIds?: string[];
		onCardClick?: (cardId: string) => void;
		onCardHover?: (cardId: string) => void;
		size?: 'small' | 'normal' | 'large';
	} = $props();

	// State
	let selectedCardId = $state<string | null>(null);
	let multiSelectMode = $state(false);

	/**
	 * Handle card click
	 */
	function handleCardClick(cardId: string, event: MouseEvent): void {
		// Check for multi-select (Shift key)
		if (event.shiftKey) {
			multiSelectMode = true;
			// Toggle selection
			if (selectedCardIds.includes(cardId)) {
				selectedCardIds = selectedCardIds.filter((id) => id !== cardId);
			} else {
				selectedCardIds = [...selectedCardIds, cardId];
			}
		} else {
			// Single select
			multiSelectMode = false;
			selectedCardId = selectedCardId === cardId ? null : cardId;
			selectedCardIds = selectedCardId ? [selectedCardId] : [];
		}

		onCardClick(cardId);
	}

	/**
	 * Handle card hover
	 */
	function handleCardHover(cardId: string): void {
		onCardHover(cardId);
	}

	/**
	 * Check if card is selected
	 */
	function isCardSelected(cardId: string): boolean {
		return selectedCardIds.includes(cardId);
	}

	// Derived values
	const handCount = $derived(cards.length);
	const isEmpty = $derived(handCount === 0);
</script>

<div class="player-hand">
	<div class="hand-header">
		<span class="hand-label">Your Hand</span>
		<span class="hand-count">({handCount})</span>
		{#if selectedCardIds.length > 0}
			<span class="selected-count">{selectedCardIds.length} selected</span>
		{/if}
	</div>

	{#if isEmpty}
		<div class="empty-state">
			<p>No cards in hand</p>
		</div>
	{:else}
		<div class="hand-cards">
			{#each cards as card (card.id)}
				<div class="card-wrapper">
					<Card
						cardId={card.id}
						cardName={card.name}
						manaCost={card.manaCost || ''}
						cardType={card.cardType || ''}
						power={card.power || ''}
						toughness={card.toughness || ''}
						imageUrl={card.imageUrl || ''}
						isTapped={card.isTapped || false}
						isSelected={isCardSelected(card.id)}
						counters={card.counters || []}
						{size}
						onclick={() => handleCardClick(card.id, event as MouseEvent)}
						onhover={() => handleCardHover(card.id)}
					/>
				</div>
			{/each}
		</div>
	{/if}

	{#if multiSelectMode}
		<div class="multi-select-hint">
			<span>💡 Shift+Click to select multiple cards</span>
		</div>
	{/if}
</div>

<style>
	.player-hand {
		display: flex;
		flex-direction: column;
		height: 100%;
		overflow: hidden;
	}

	.hand-header {
		display: flex;
		align-items: center;
		gap: 0.5rem;
		padding: 0.75rem 1rem;
		background: #141821;
		border-bottom: 1px solid #2a3441;
	}

	.hand-label {
		font-size: 0.875rem;
		color: #6b7280;
		font-weight: 500;
		text-transform: uppercase;
		letter-spacing: 0.5px;
	}

	.hand-count {
		font-size: 0.875rem;
		color: #9ca3af;
	}

	.selected-count {
		margin-left: auto;
		padding: 0.25rem 0.5rem;
		background: #667eea;
		border-radius: 4px;
		font-size: 0.75rem;
		font-weight: 600;
		color: white;
	}

	.empty-state {
		flex: 1;
		display: flex;
		align-items: center;
		justify-content: center;
	}

	.empty-state p {
		color: #4b5563;
		font-style: italic;
		margin: 0;
	}

	.hand-cards {
		flex: 1;
		display: flex;
		gap: 0.5rem;
		padding: 1rem;
		overflow-x: auto;
		overflow-y: hidden;
		align-items: center;
	}

	.card-wrapper {
		flex-shrink: 0;
	}

	.multi-select-hint {
		padding: 0.5rem 1rem;
		background: #141821;
		border-top: 1px solid #2a3441;
		text-align: center;
	}

	.multi-select-hint span {
		font-size: 0.75rem;
		color: #6b7280;
	}

	/* Scrollbar Styling */
	.hand-cards::-webkit-scrollbar {
		height: 8px;
	}

	.hand-cards::-webkit-scrollbar-track {
		background: #0d1117;
	}

	.hand-cards::-webkit-scrollbar-thumb {
		background: #3a4451;
		border-radius: 4px;
	}

	.hand-cards::-webkit-scrollbar-thumb:hover {
		background: #4a5461;
	}

	/* Responsive */
	@media (max-width: 768px) {
		.hand-cards {
			padding: 0.5rem;
			gap: 0.25rem;
		}

		.hand-header {
			padding: 0.5rem;
		}

		.multi-select-hint {
			display: none;
		}
	}
</style>
