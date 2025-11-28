<script lang="ts">
	import Card from './Card.svelte';
	import type { GameCard } from '$lib/types/game';
	import type { CardView } from '$lib/generated/mage/v1/models';
	import { myHand, selectedCards, gameStore } from '$lib/stores/game';
	import {
		isTargetingActive,
		validTargetIds,
		selectedTargetIds,
		targetingStore
	} from '$lib/stores/game-targeting';

	// Props (optional callbacks for additional handling like target selection)
	let {
		// eslint-disable-next-line no-unused-vars
		onCardClick = (cardId: string) => {},
		// eslint-disable-next-line no-unused-vars
		onCardHover = (cardId: string) => {},
		size = 'normal'
	}: {
		// eslint-disable-next-line no-unused-vars
		onCardClick?: (cardId: string) => void;
		// eslint-disable-next-line no-unused-vars
		onCardHover?: (cardId: string) => void;
		size?: 'small' | 'normal' | 'large';
	} = $props();

	// Targeting state from store
	const isTargeting = $derived($isTargetingActive);
	const validTargets = $derived($validTargetIds);
	const selectedTargets = $derived($selectedTargetIds);

	// State
	let multiSelectMode = $state(false);

	/**
	 * Convert CardView from proto to GameCard for components
	 */
	function toGameCard(card: CardView): GameCard {
		return {
			id: card.id,
			name: card.name,
			manaCost: card.manaCost,
			cardType: card.type,
			power: card.power,
			toughness: card.toughness,
			imageUrl: '',
			isTapped: card.tapped,
			isSelected: false,
			ownerId: card.ownerId,
			controllerId: card.controllerId
		};
	}

	// Always use game store for cards and selection
	const handCards = $derived(($myHand || []).map(toGameCard));
	const selectedCardIds = $derived($selectedCards || []);

	/**
	 * Handle card click
	 */
	function handleCardClick(cardId: string, event?: MouseEvent | KeyboardEvent): void {
		// Handle targeting mode - toggle target selection
		if (isTargeting) {
			targetingStore.toggleTarget(cardId);
			onCardClick(cardId);
			return;
		}

		// Check for multi-select (Shift key)
		if (event?.shiftKey) {
			multiSelectMode = true;
			// Toggle selection in game store
			gameStore.toggleCardSelection(cardId);
		} else {
			// Single select - clear all and select this one
			multiSelectMode = false;
			const isSelected = selectedCardIds.includes(cardId);
			gameStore.clearSelection();
			if (!isSelected) {
				gameStore.toggleCardSelection(cardId);
			}
		}

		onCardClick(cardId);
	}

	/**
	 * Handle keyboard events for accessibility
	 */
	function handleKeydown(cardId: string, event: KeyboardEvent): void {
		if (event.key === 'Enter' || event.key === ' ') {
			event.preventDefault();
			handleCardClick(cardId, event);
		}
	}

	/**
	 * Handle card hover
	 */
	function handleCardHover(_cardId: string): void {
		onCardHover(_cardId);
	}

	/**
	 * Check if card is selected
	 */
	function isCardSelected(cardId: string): boolean {
		return selectedCardIds.includes(cardId);
	}

	// Derived values
	const handCount = $derived(handCards.length);
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
			{#each handCards as card (card.id)}
				<div
					class="card-wrapper"
					role="button"
					tabindex="0"
					onclick={(e) => handleCardClick(card.id, e)}
					onkeydown={(e) => handleKeydown(card.id, e)}
					onmouseenter={() => handleCardHover(card.id)}
					aria-label={`Card: ${card.name}`}
				>
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
						onclick={() => {}}
						onhover={() => {}}
						isTargetingActive={isTargeting}
						isValidTarget={validTargets.has(card.id)}
						isTargetSelected={selectedTargets.includes(card.id)}
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
