<script lang="ts">
	import Card from './Card.svelte';
	import type { GameCard } from '$lib/types/game';

	// Props
	let {
		cards = [],
		playerName = 'Player',
		isOpponent = false,
		// eslint-disable-next-line no-unused-vars
		onCardClick = (cardId: string) => {}
	}: {
		cards?: GameCard[];
		playerName?: string;
		isOpponent?: boolean;
		// eslint-disable-next-line no-unused-vars
		onCardClick?: (cardId: string) => void;
	} = $props();

	// State
	let showModal = $state(false);
	let selectedCardId = $state<string | null>(null);

	/**
	 * Toggle graveyard modal
	 */
	function toggleModal(): void {
		showModal = !showModal;
	}

	/**
	 * Close modal
	 */
	function closeModal(): void {
		showModal = false;
	}

	/**
	 * Handle card click in modal
	 */
	function handleCardClick(_cardId: string): void {
		selectedCardId = selectedCardId === _cardId ? null : _cardId;
		onCardClick(_cardId);
	}

	/**
	 * Handle backdrop click
	 */
	function handleBackdropClick(event: MouseEvent): void {
		if (event.target === event.currentTarget) {
			closeModal();
		}
	}

	// Derived values
	const cardCount = $derived(cards.length);
	const isEmpty = $derived(cardCount === 0);
	const topCard = $derived(cards.length > 0 ? cards[cards.length - 1] : null);
</script>

<div class="graveyard-zone" class:opponent={isOpponent}>
	<button
		class="graveyard-button"
		class:empty={isEmpty}
		onclick={toggleModal}
		disabled={isEmpty}
		title="{playerName}'s Graveyard ({cardCount} cards)"
	>
		{#if topCard}
			<div class="top-card-preview">
				<Card
					cardId={topCard.id}
					cardName={topCard.name}
					manaCost={topCard.manaCost || ''}
					cardType={topCard.cardType || ''}
					power={topCard.power || ''}
					toughness={topCard.toughness || ''}
					imageUrl={topCard.imageUrl || ''}
					size="small"
					onclick={() => {}}
					onhover={() => {}}
				/>
			</div>
		{:else}
			<div class="empty-graveyard-icon">🪦</div>
		{/if}

		<div class="card-count-badge" class:zero={cardCount === 0}>
			{cardCount}
		</div>
	</button>

	<div class="graveyard-label">Graveyard</div>
</div>

<!-- Graveyard Modal -->
{#if showModal}
	<!-- svelte-ignore a11y_click_events_have_key_events -->
	<!-- svelte-ignore a11y_no_static_element_interactions -->
	<div class="graveyard-modal-backdrop" onclick={handleBackdropClick}>
		<div class="graveyard-modal">
			<div class="modal-header">
				<h3>{playerName}'s Graveyard</h3>
				<span class="card-count-text">{cardCount} card{cardCount !== 1 ? 's' : ''}</span>
				<button class="close-button" onclick={closeModal} title="Close">✕</button>
			</div>

			<div class="modal-content">
				{#if isEmpty}
					<div class="empty-state">
						<p>No cards in graveyard</p>
					</div>
				{:else}
					<div class="card-grid">
						{#each cards as card (card.id)}
							<div class="card-grid-item">
								<Card
									cardId={card.id}
									cardName={card.name}
									manaCost={card.manaCost || ''}
									cardType={card.cardType || ''}
									power={card.power || ''}
									toughness={card.toughness || ''}
									imageUrl={card.imageUrl || ''}
									isSelected={selectedCardId === card.id}
									size="normal"
									onclick={() => handleCardClick(card.id)}
								/>
							</div>
						{/each}
					</div>
				{/if}
			</div>
		</div>
	</div>
{/if}

<style>
	/* Graveyard Zone */
	.graveyard-zone {
		display: flex;
		flex-direction: column;
		align-items: center;
		gap: 0.5rem;
	}

	.graveyard-button {
		position: relative;
		width: 80px;
		height: 112px;
		border: 2px solid #4b5563;
		border-radius: 8px;
		background: #1a1f2e;
		cursor: pointer;
		transition:
			transform 0.2s,
			border-color 0.2s,
			box-shadow 0.2s;
		overflow: hidden;
	}

	.graveyard-button:not(.empty):hover {
		transform: scale(1.05);
		border-color: #9ca3af;
		box-shadow: 0 4px 8px rgba(0, 0, 0, 0.3);
	}

	.graveyard-button:disabled {
		cursor: not-allowed;
		opacity: 0.5;
	}

	.graveyard-button.empty {
		border-style: dashed;
		border-color: #374151;
	}

	.top-card-preview {
		width: 100%;
		height: 100%;
		display: flex;
		align-items: center;
		justify-content: center;
	}

	.empty-graveyard-icon {
		width: 100%;
		height: 100%;
		display: flex;
		align-items: center;
		justify-content: center;
		font-size: 2rem;
		opacity: 0.3;
	}

	.card-count-badge {
		position: absolute;
		top: 0.25rem;
		right: 0.25rem;
		padding: 0.25rem 0.5rem;
		background: #374151;
		border-radius: 4px;
		font-size: 0.75rem;
		font-weight: 700;
		color: white;
		box-shadow: 0 2px 4px rgba(0, 0, 0, 0.5);
	}

	.card-count-badge.zero {
		background: #1f2937;
		color: #6b7280;
	}

	.graveyard-label {
		font-size: 0.75rem;
		color: #6b7280;
		text-transform: uppercase;
		letter-spacing: 0.5px;
	}

	/* Modal Backdrop */
	.graveyard-modal-backdrop {
		position: fixed;
		top: 0;
		left: 0;
		right: 0;
		bottom: 0;
		background: rgba(0, 0, 0, 0.75);
		display: flex;
		align-items: center;
		justify-content: center;
		z-index: 1000;
		animation: fadeIn 0.2s;
	}

	@keyframes fadeIn {
		from {
			opacity: 0;
		}
		to {
			opacity: 1;
		}
	}

	/* Modal */
	.graveyard-modal {
		background: #1a1f2e;
		border-radius: 12px;
		border: 2px solid #2a3441;
		max-width: 90vw;
		max-height: 85vh;
		width: 800px;
		display: flex;
		flex-direction: column;
		box-shadow: 0 20px 25px -5px rgba(0, 0, 0, 0.5);
		animation: slideIn 0.3s;
	}

	@keyframes slideIn {
		from {
			transform: translateY(20px);
			opacity: 0;
		}
		to {
			transform: translateY(0);
			opacity: 1;
		}
	}

	.modal-header {
		display: flex;
		align-items: center;
		gap: 1rem;
		padding: 1.5rem;
		border-bottom: 1px solid #2a3441;
		background: #141821;
	}

	.modal-header h3 {
		margin: 0;
		font-size: 1.25rem;
		font-weight: 600;
		color: white;
	}

	.card-count-text {
		font-size: 0.875rem;
		color: #9ca3af;
	}

	.close-button {
		margin-left: auto;
		width: 32px;
		height: 32px;
		display: flex;
		align-items: center;
		justify-content: center;
		background: transparent;
		border: none;
		border-radius: 4px;
		font-size: 1.5rem;
		color: #9ca3af;
		cursor: pointer;
		transition: all 0.2s;
	}

	.close-button:hover {
		background: #2a3441;
		color: white;
	}

	.modal-content {
		flex: 1;
		overflow-y: auto;
		padding: 1.5rem;
	}

	.empty-state {
		display: flex;
		align-items: center;
		justify-content: center;
		min-height: 200px;
	}

	.empty-state p {
		color: #6b7280;
		font-style: italic;
		margin: 0;
	}

	.card-grid {
		display: grid;
		grid-template-columns: repeat(auto-fill, minmax(100px, 1fr));
		gap: 1rem;
	}

	.card-grid-item {
		display: flex;
		justify-content: center;
	}

	/* Scrollbar Styling */
	.modal-content::-webkit-scrollbar {
		width: 8px;
	}

	.modal-content::-webkit-scrollbar-track {
		background: #0d1117;
	}

	.modal-content::-webkit-scrollbar-thumb {
		background: #3a4451;
		border-radius: 4px;
	}

	.modal-content::-webkit-scrollbar-thumb:hover {
		background: #4a5461;
	}

	/* Responsive */
	@media (max-width: 768px) {
		.graveyard-modal {
			width: 95vw;
			max-height: 90vh;
		}

		.modal-header {
			padding: 1rem;
		}

		.modal-content {
			padding: 1rem;
		}

		.card-grid {
			grid-template-columns: repeat(auto-fill, minmax(80px, 1fr));
			gap: 0.75rem;
		}
	}
</style>
