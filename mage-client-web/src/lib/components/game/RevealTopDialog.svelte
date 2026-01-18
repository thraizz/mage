<script lang="ts">
	import type { CardView } from '$lib/generated/mage/v1/models';
	import Card from './Card.svelte';

	interface Props {
		cards: CardView[];
		onClose: () => void;
	}

	let { cards, onClose }: Props = $props();

	function handleKeydown(e: KeyboardEvent) {
		if (e.key === 'Escape') {
			onClose();
		}
	}
</script>

<svelte:window onkeydown={handleKeydown} />

<div class="overlay" role="dialog" aria-labelledby="reveal-dialog-title" onclick={onClose}>
	<div class="dialog" onclick={(e) => e.stopPropagation()}>
		<div class="dialog-header">
			<h2 id="reveal-dialog-title">Revealed Cards ({cards.length})</h2>
			<button class="close-button" onclick={onClose} aria-label="Close dialog">×</button>
		</div>

		<div class="cards-section">
			{#if cards.length === 0}
				<p class="no-cards">No cards to reveal</p>
			{:else}
				<div class="card-grid">
					{#each cards as card (card.id)}
						<div class="card-wrapper">
							<Card {card} interactive={false} />
						</div>
					{/each}
				</div>
			{/if}
		</div>

		<div class="dialog-footer">
			<button class="btn-primary" onclick={onClose}>Close</button>
		</div>
	</div>
</div>

<style>
	.overlay {
		position: fixed;
		inset: 0;
		background: rgba(0, 0, 0, 0.7);
		display: flex;
		align-items: center;
		justify-content: center;
		z-index: 1000;
		background: rgba(0, 0, 0, 0.85);
		animation: fadeIn 0.2s ease-out;
	}

	@keyframes fadeIn {
		from {
			opacity: 0;
		}
		to {
			opacity: 1;
		}
	}

	.dialog {
		background: #1a1f2e;
		border: 2px solid #3a4451;
		border-radius: 12px;
		padding: 1.5rem;
		max-width: 900px;
		width: 90%;
		max-height: 85vh;
		overflow-y: auto;
		box-shadow: 0 20px 25px -5px rgba(0, 0, 0, 0.5);
		animation: slideUp 0.3s ease-out;
	}

	@keyframes slideUp {
		from {
			transform: translateY(20px);
			opacity: 0;
		}
		to {
			transform: translateY(0);
			opacity: 1;
		}
	}

	.dialog-header {
		display: flex;
		justify-content: space-between;
		align-items: center;
		margin-bottom: 1.5rem;
		padding-bottom: 1rem;
		border-bottom: 1px solid #3a4451;
	}

	.dialog-header h2 {
		margin: 0;
		color: #fff;
		font-size: 1.25rem;
		font-weight: 600;
	}

	.close-button {
		background: transparent;
		border: none;
		color: #9ca3af;
		font-size: 2rem;
		line-height: 1;
		cursor: pointer;
		padding: 0;
		width: 2rem;
		height: 2rem;
		display: flex;
		align-items: center;
		justify-content: center;
		transition: color 0.2s;
	}

	.close-button:hover {
		color: #fff;
	}

	.cards-section {
		margin-bottom: 1.5rem;
	}

	.no-cards {
		color: #6b7280;
		font-style: italic;
		padding: 2rem;
		text-align: center;
		background: rgba(255, 255, 255, 0.02);
		border-radius: 6px;
	}

	.card-grid {
		display: grid;
		grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
		gap: 1rem;
		padding: 0.5rem;
	}

	.card-wrapper {
		display: flex;
		justify-content: center;
	}

	.dialog-footer {
		display: flex;
		justify-content: flex-end;
		gap: 0.75rem;
		padding-top: 1rem;
		border-top: 1px solid #3a4451;
	}

	.btn-primary {
		padding: 0.5rem 1rem;
		border-radius: 6px;
		font-weight: 600;
		font-size: 0.875rem;
		cursor: pointer;
		transition: all 0.2s;
		border: none;
		background: #667eea;
		color: white;
	}

	.btn-primary:hover {
		background: #5568d3;
		transform: translateY(-1px);
	}
</style>
