<script lang="ts">
	import type { CardView } from '$lib/generated/mage/v1/models';
	import Card from './Card.svelte';

	interface Props {
		cards: CardView[];
		onComplete: (keepOnTop: CardView[], putToBottom: CardView[]) => void;
		onCancel: () => void;
	}

	let { cards, onComplete, onCancel }: Props = $props();

	let keepOnTop = $state<CardView[]>([]);
	let putToBottom = $state<CardView[]>([]);

	let remaining = $derived(
		cards.filter(
			(card) =>
				!keepOnTop.some((c) => c.id === card.id) && !putToBottom.some((c) => c.id === card.id)
		)
	);

	function moveToTop(card: CardView) {
		keepOnTop = [...keepOnTop, card];
	}

	function moveToBottom(card: CardView) {
		putToBottom = [...putToBottom, card];
	}

	function handleComplete() {
		if (remaining.length > 0) {
			alert('Please assign all cards to either keep on top or put to bottom.');
			return;
		}
		onComplete(keepOnTop, putToBottom);
	}

	function handleKeydown(e: KeyboardEvent) {
		if (e.key === 'Escape') {
			onCancel();
		} else if (e.key === 'Enter' && remaining.length === 0) {
			handleComplete();
		}
	}
</script>

<svelte:window onkeydown={handleKeydown} />

<div class="overlay" role="dialog" aria-labelledby="scry-dialog-title">
	<div class="dialog">
		<div class="dialog-header">
			<h2 id="scry-dialog-title">Scry {cards.length}</h2>
			<button class="close-button" onclick={onCancel} aria-label="Close dialog">×</button>
		</div>

		<div class="instructions">
			<p>Choose where to put each card - keep on top or put to bottom.</p>
			{#if remaining.length > 0}
				<p class="warning">⚠️ {remaining.length} card(s) not yet assigned</p>
			{/if}
		</div>

		<div class="card-grid">
			{#each remaining as card (card.id)}
				<div class="card-item">
					<Card
						cardId={card.id}
						cardName={card.name}
						manaCost={card.manaCost}
						cardType={card.type}
						power={card.power}
						toughness={card.toughness}
						isSelected={false}
						size="normal"
						onclick={() => {}}
						onhover={() => {}}
					/>
					<div class="card-actions">
						<button class="btn-top" onclick={() => moveToTop(card)}>↑ Top</button>
						<button class="btn-bottom" onclick={() => moveToBottom(card)}>↓ Bottom</button>
					</div>
				</div>
			{/each}
			{#if remaining.length === 0}
				<p class="all-assigned">All cards assigned!</p>
			{/if}
		</div>

		<div class="dialog-footer">
			<div class="assignment-summary">
				<span>Top: {keepOnTop.length}</span>
				<span>Bottom: {putToBottom.length}</span>
			</div>
			<div class="footer-actions">
				<button class="btn-secondary" onclick={onCancel}>Cancel</button>
				<button class="btn-primary" onclick={handleComplete}> Complete </button>
			</div>
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
		width: 95%;
		max-height: 90vh;
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
		margin-bottom: 1rem;
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

	.instructions {
		margin-bottom: 1.5rem;
		padding: 1rem;
		background: rgba(102, 126, 234, 0.1);
		border: 1px solid #667eea;
		border-radius: 8px;
	}

	.instructions p {
		margin: 0;
		color: #ddd;
		font-size: 0.875rem;
	}

	.instructions .warning {
		margin-top: 0.5rem;
		color: #fbbf24;
		font-weight: 600;
	}

	.card-grid {
		display: grid;
		grid-template-columns: repeat(auto-fill, minmax(180px, 1fr));
		gap: 1rem;
		margin-bottom: 1.5rem;
		min-height: 200px;
	}

	.card-item {
		background: rgba(255, 255, 255, 0.05);
		border: 1px solid #3a4451;
		border-radius: 8px;
		padding: 0.75rem;
		display: flex;
		flex-direction: column;
		gap: 0.75rem;
		transition: all 0.2s;
	}

	.card-item:hover {
		background: rgba(255, 255, 255, 0.08);
		transform: translateY(-2px);
	}

	.card-actions {
		display: flex;
		gap: 0.5rem;
	}

	.btn-top,
	.btn-bottom {
		flex: 1;
		padding: 0.5rem;
		border-radius: 6px;
		font-size: 0.875rem;
		font-weight: 600;
		cursor: pointer;
		transition: all 0.2s;
		border: none;
	}

	.btn-top {
		background: rgba(16, 185, 129, 0.2);
		color: #10b981;
		border: 1px solid #10b981;
	}

	.btn-top:hover {
		background: rgba(16, 185, 129, 0.3);
		transform: translateY(-1px);
	}

	.btn-bottom {
		background: rgba(239, 68, 68, 0.2);
		color: #ef4444;
		border: 1px solid #ef4444;
	}

	.btn-bottom:hover {
		background: rgba(239, 68, 68, 0.3);
		transform: translateY(-1px);
	}

	.all-assigned {
		grid-column: 1 / -1;
		text-align: center;
		color: #10b981;
		font-weight: 600;
		font-size: 1rem;
		padding: 2rem;
	}

	.dialog-footer {
		display: flex;
		justify-content: space-between;
		align-items: center;
		padding-top: 1rem;
		border-top: 1px solid #3a4451;
	}

	.assignment-summary {
		display: flex;
		gap: 1.5rem;
		color: #9ca3af;
		font-size: 0.875rem;
		font-weight: 600;
	}

	.footer-actions {
		display: flex;
		gap: 0.75rem;
	}
</style>
