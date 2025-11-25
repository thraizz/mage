<script lang="ts">
	import type { CardView } from '$lib/generated/mage/v1/models';
	import Card from './Card.svelte';

	interface Props {
		cards: CardView[];
		mulliganCount: number;
		onKeep: () => void;
		onMulligan: () => void;
		isLoading?: boolean;
	}

	let { cards, mulliganCount, onKeep, onMulligan, isLoading = false }: Props = $props();

	// Calculate next hand size after mulligan
	const nextHandSize = $derived(Math.max(0, 7 - (mulliganCount + 1)));

	// Determine if this is the first draw or subsequent mulligan
	const isFirstDraw = $derived(mulliganCount === 0);

	// Show warning when hand will be very small
	const showWarning = $derived(nextHandSize <= 3);
</script>

<div class="mulligan-overlay">
	<div class="mulligan-dialog">
		<div class="dialog-header">
			<h2>
				{#if isFirstDraw}
					Opening Hand
				{:else}
					Mulligan #{mulliganCount}
				{/if}
			</h2>
			<p class="hand-info">
				{cards.length} cards in hand
				{#if mulliganCount > 0}
					<span class="mulligan-note">({mulliganCount} put on bottom)</span>
				{/if}
			</p>
		</div>

		<div class="hand-display">
			{#each cards as card (card.id)}
				<div class="mulligan-card">
					<Card
						cardId={card.id}
						cardName={card.name}
						manaCost={card.manaCost}
						cardType={card.type}
						power={card.power}
						toughness={card.toughness}
						imageUrl=""
						isTapped={false}
						isSelected={false}
						size="large"
					/>
				</div>
			{/each}
		</div>

		<div class="mulligan-info">
			{#if nextHandSize > 0}
				<p class="next-hand-info">
					If you mulligan, you'll draw 7 cards and put <strong>{mulliganCount + 1}</strong> on the bottom.
				</p>
			{:else}
				<p class="warning">
					⚠️ If you mulligan again, you'll have 0 cards!
				</p>
			{/if}
		</div>

		<div class="dialog-actions">
			<button
				class="btn-keep"
				onclick={onKeep}
				disabled={isLoading}
			>
				{#if isLoading}
					<span class="spinner-small"></span>
				{:else}
					Keep Hand
				{/if}
			</button>
			<button
				class="btn-mulligan"
				onclick={onMulligan}
				disabled={isLoading || nextHandSize === 0}
				class:warning={showWarning}
			>
				{#if isLoading}
					<span class="spinner-small"></span>
				{:else if nextHandSize === 0}
					No Mulligan Available
				{:else}
					Mulligan to {nextHandSize} cards
				{/if}
			</button>
		</div>
	</div>
</div>

<style>
	.mulligan-overlay {
		position: fixed;
		inset: 0;
		background: rgba(0, 0, 0, 0.9);
		display: flex;
		align-items: center;
		justify-content: center;
		z-index: 200;
		padding: 2rem;
	}

	.mulligan-dialog {
		background: linear-gradient(135deg, #1a1f2e 0%, #0f1419 100%);
		border-radius: 16px;
		border: 2px solid #667eea;
		max-width: 1000px;
		width: 100%;
		padding: 2rem;
		box-shadow: 0 25px 50px -12px rgba(0, 0, 0, 0.5),
			0 0 40px rgba(102, 126, 234, 0.3);
	}

	.dialog-header {
		text-align: center;
		margin-bottom: 1.5rem;
	}

	.dialog-header h2 {
		font-size: 1.75rem;
		font-weight: 700;
		color: #fff;
		margin: 0 0 0.5rem;
		text-shadow: 0 2px 4px rgba(0, 0, 0, 0.3);
	}

	.hand-info {
		color: #94a3b8;
		font-size: 1rem;
		margin: 0;
	}

	.mulligan-note {
		color: #fbbf24;
		font-style: italic;
	}

	.hand-display {
		display: flex;
		justify-content: center;
		gap: 0.75rem;
		flex-wrap: wrap;
		padding: 1.5rem;
		background: rgba(0, 0, 0, 0.3);
		border-radius: 12px;
		margin-bottom: 1.5rem;
		min-height: 280px;
		align-items: center;
	}

	.mulligan-card {
		transition: transform 0.2s ease, box-shadow 0.2s ease;
	}

	.mulligan-card:hover {
		transform: translateY(-12px) scale(1.05);
		z-index: 10;
	}

	.mulligan-info {
		text-align: center;
		margin-bottom: 1.5rem;
		padding: 1rem;
		background: rgba(102, 126, 234, 0.1);
		border-radius: 8px;
		border: 1px solid rgba(102, 126, 234, 0.2);
	}

	.next-hand-info {
		color: #94a3b8;
		margin: 0;
	}

	.next-hand-info strong {
		color: #fbbf24;
	}

	.warning {
		color: #fbbf24;
		margin: 0;
		font-weight: 600;
	}

	.dialog-actions {
		display: flex;
		justify-content: center;
		gap: 1.5rem;
	}

	.btn-keep,
	.btn-mulligan {
		padding: 1rem 2.5rem;
		font-size: 1.125rem;
		font-weight: 700;
		border: none;
		border-radius: 8px;
		cursor: pointer;
		transition: all 0.2s ease;
		display: flex;
		align-items: center;
		gap: 0.5rem;
	}

	.btn-keep {
		background: linear-gradient(135deg, #22c55e 0%, #16a34a 100%);
		color: white;
		box-shadow: 0 4px 15px rgba(34, 197, 94, 0.3);
	}

	.btn-keep:hover:not(:disabled) {
		transform: translateY(-2px);
		box-shadow: 0 6px 20px rgba(34, 197, 94, 0.4);
	}

	.btn-keep:disabled {
		opacity: 0.6;
		cursor: not-allowed;
	}

	.btn-mulligan {
		background: linear-gradient(135deg, #667eea 0%, #5568d3 100%);
		color: white;
		box-shadow: 0 4px 15px rgba(102, 126, 234, 0.3);
	}

	.btn-mulligan:hover:not(:disabled) {
		transform: translateY(-2px);
		box-shadow: 0 6px 20px rgba(102, 126, 234, 0.4);
	}

	.btn-mulligan.warning {
		background: linear-gradient(135deg, #f59e0b 0%, #d97706 100%);
		box-shadow: 0 4px 15px rgba(245, 158, 11, 0.3);
	}

	.btn-mulligan.warning:hover:not(:disabled) {
		box-shadow: 0 6px 20px rgba(245, 158, 11, 0.4);
	}

	.btn-mulligan:disabled {
		opacity: 0.5;
		cursor: not-allowed;
		background: #374151;
		box-shadow: none;
	}

	.spinner-small {
		width: 20px;
		height: 20px;
		border: 2px solid transparent;
		border-top-color: currentColor;
		border-radius: 50%;
		animation: spin 0.6s linear infinite;
	}

	@keyframes spin {
		to {
			transform: rotate(360deg);
		}
	}

	/* Responsive */
	@media (max-width: 768px) {
		.mulligan-overlay {
			padding: 1rem;
		}

		.mulligan-dialog {
			padding: 1rem;
		}

		.dialog-header h2 {
			font-size: 1.25rem;
		}

		.hand-display {
			gap: 0.5rem;
			padding: 1rem;
		}

		.dialog-actions {
			flex-direction: column;
			gap: 0.75rem;
		}

		.btn-keep,
		.btn-mulligan {
			width: 100%;
			justify-content: center;
		}
	}
</style>

