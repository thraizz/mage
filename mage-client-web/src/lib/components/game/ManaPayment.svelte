<script lang="ts">
	import { sendPlayerManaType, sendPlayerBoolean } from '$lib/api/game';
	import ManaSymbol from '$lib/components/mtg/ManaSymbol.svelte';
	import type { GamePlayManaData, ManaOption } from '$lib/generated/mage/v1/websocket';

	// Props
	let {
		gameId,
		manaData,
		onComplete = () => {},
		onCancel = () => {}
	}: {
		gameId: string;
		manaData: GamePlayManaData;
		onComplete?: () => void;
		onCancel?: () => void;
	} = $props();

	// State
	let isLoading = $state(false);
	let selectedMana = $state<string | null>(null);
	let error = $state<string | null>(null);

	// Parse the message for cost display
	const costMessage = $derived(manaData.message || 'Pay mana cost');

	// Available mana options from server
	const availableOptions = $derived(manaData.manaOptions || []);

	// Check if any mana is available
	const hasAnyMana = $derived(availableOptions.some((opt) => opt.amount > 0));

	// Mana color configuration for display
	const manaColorConfig: Record<
		string,
		{ label: string; symbol: string; bgColor: string; textColor: string }
	> = {
		W: { label: 'White', symbol: 'W', bgColor: '#f8f6e3', textColor: '#000' },
		U: { label: 'Blue', symbol: 'U', bgColor: '#0e68ab', textColor: '#fff' },
		B: { label: 'Black', symbol: 'B', bgColor: '#150b00', textColor: '#fff' },
		R: { label: 'Red', symbol: 'R', bgColor: '#d3202a', textColor: '#fff' },
		G: { label: 'Green', symbol: 'G', bgColor: '#00733e', textColor: '#fff' },
		C: { label: 'Colorless', symbol: 'C', bgColor: '#9ca3af', textColor: '#000' }
	};

	/**
	 * Get display config for a mana color
	 */
	function getColorConfig(color: string) {
		return manaColorConfig[color.toUpperCase()] || manaColorConfig['C'];
	}

	/**
	 * Handle mana option selection
	 */
	async function handleManaSelect(option: ManaOption) {
		if (option.amount <= 0 || isLoading) return;

		selectedMana = option.color;
		isLoading = true;
		error = null;

		try {
			// Send the selected mana type to the server
			await sendPlayerManaType(gameId, option.color, option.color);
			onComplete();
		} catch (err) {
			error = err instanceof Error ? err.message : 'Failed to pay mana';
			selectedMana = null;
		} finally {
			isLoading = false;
		}
	}

	/**
	 * Handle cancel - send false to server to cancel mana payment
	 */
	async function handleCancel() {
		if (isLoading) return;

		isLoading = true;
		error = null;

		try {
			// Send false to indicate cancellation
			await sendPlayerBoolean(gameId, false);
			onCancel();
		} catch (err) {
			error = err instanceof Error ? err.message : 'Failed to cancel';
		} finally {
			isLoading = false;
		}
	}

	/**
	 * Handle keyboard shortcuts
	 */
	function handleKeydown(event: KeyboardEvent) {
		if (event.key === 'Escape') {
			handleCancel();
		}
	}
</script>

<svelte:window onkeydown={handleKeydown} />

<div class="mana-payment-overlay" role="dialog" aria-labelledby="mana-title" aria-modal="true">
	<div class="mana-payment-modal">
		<div class="modal-header">
			<h2 id="mana-title">Mana Payment</h2>
		</div>

		<div class="modal-content">
			<p class="cost-message">{costMessage}</p>

			{#if error}
				<div class="error-message" role="alert">
					<span class="error-icon">⚠️</span>
					{error}
				</div>
			{/if}

			{#if hasAnyMana}
				<div class="mana-options">
					<p class="instructions">Click a mana type to pay:</p>
					<div class="mana-grid">
						{#each availableOptions as option}
							{@const config = getColorConfig(option.color)}
							{@const isAvailable = option.amount > 0}
							{@const isSelected = selectedMana === option.color}
							<button
								class="mana-option"
								class:available={isAvailable}
								class:unavailable={!isAvailable}
								class:selected={isSelected}
								disabled={!isAvailable || isLoading}
								onclick={() => handleManaSelect(option)}
								title="{config.label} mana ({option.amount} available)"
								style="--mana-bg: {config.bgColor}; --mana-text: {config.textColor}"
							>
								<div class="mana-symbol-wrapper">
									<ManaSymbol symbol={config.symbol} size="lg" />
								</div>
								<span class="mana-count">{option.amount}</span>
								<span class="mana-label">{config.label}</span>
							</button>
						{/each}
					</div>
				</div>
			{:else}
				<div class="no-mana-warning">
					<span class="warning-icon">💧</span>
					<p>No mana available to pay this cost.</p>
					<p class="hint">Tap lands or other mana sources first.</p>
				</div>
			{/if}
		</div>

		<div class="modal-footer">
			<button class="btn-cancel" onclick={handleCancel} disabled={isLoading}>
				{isLoading ? 'Canceling...' : 'Cancel'}
			</button>
		</div>
	</div>
</div>

<style>
	.mana-payment-overlay {
		position: fixed;
		inset: 0;
		background: rgba(0, 0, 0, 0.85);
		display: flex;
		align-items: center;
		justify-content: center;
		z-index: 100;
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

	.mana-payment-modal {
		background: #1a1f2e;
		border: 2px solid #667eea;
		border-radius: 12px;
		width: 90%;
		max-width: 500px;
		box-shadow:
			0 25px 50px -12px rgba(0, 0, 0, 0.5),
			0 0 40px rgba(102, 126, 234, 0.15);
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

	.modal-header {
		padding: 1.25rem 1.5rem;
		background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
		border-radius: 10px 10px 0 0;
	}

	.modal-header h2 {
		margin: 0;
		font-size: 1.25rem;
		font-weight: 600;
		color: white;
		text-transform: uppercase;
		letter-spacing: 0.5px;
	}

	.modal-content {
		padding: 1.5rem;
	}

	.cost-message {
		margin: 0 0 1.5rem;
		font-size: 1.125rem;
		color: #e2e8f0;
		text-align: center;
		line-height: 1.5;
	}

	.error-message {
		display: flex;
		align-items: center;
		gap: 0.5rem;
		padding: 0.75rem 1rem;
		background: rgba(239, 68, 68, 0.15);
		border: 1px solid rgba(239, 68, 68, 0.3);
		border-radius: 8px;
		color: #f87171;
		margin-bottom: 1rem;
	}

	.error-icon {
		font-size: 1rem;
	}

	.instructions {
		margin: 0 0 1rem;
		font-size: 0.875rem;
		color: #9ca3af;
		text-align: center;
	}

	.mana-grid {
		display: grid;
		grid-template-columns: repeat(auto-fit, minmax(100px, 1fr));
		gap: 0.75rem;
	}

	.mana-option {
		display: flex;
		flex-direction: column;
		align-items: center;
		gap: 0.5rem;
		padding: 1rem 0.75rem;
		background: #141821;
		border: 2px solid #2a3441;
		border-radius: 10px;
		cursor: pointer;
		transition: all 0.2s ease;
	}

	.mana-option.available:hover {
		border-color: var(--mana-bg, #667eea);
		background: rgba(255, 255, 255, 0.05);
		transform: translateY(-2px);
		box-shadow: 0 4px 12px rgba(0, 0, 0, 0.3);
	}

	.mana-option.available:active {
		transform: translateY(0);
	}

	.mana-option.selected {
		border-color: var(--mana-bg, #667eea);
		background: rgba(var(--mana-bg), 0.1);
		box-shadow:
			0 0 0 3px rgba(102, 126, 234, 0.3),
			0 4px 12px rgba(0, 0, 0, 0.3);
	}

	.mana-option.unavailable {
		opacity: 0.4;
		cursor: not-allowed;
	}

	.mana-option:disabled {
		cursor: not-allowed;
	}

	.mana-option:focus-visible {
		outline: none;
		box-shadow:
			0 0 0 3px rgba(102, 126, 234, 0.5),
			0 4px 12px rgba(0, 0, 0, 0.3);
	}

	.mana-symbol-wrapper {
		display: flex;
		align-items: center;
		justify-content: center;
		width: 48px;
		height: 48px;
		border-radius: 50%;
		background: rgba(255, 255, 255, 0.1);
	}

	.mana-option.available:hover .mana-symbol-wrapper {
		background: rgba(255, 255, 255, 0.15);
		transform: scale(1.05);
		transition: all 0.2s ease;
	}

	.mana-count {
		font-size: 1.5rem;
		font-weight: 700;
		color: #ffffff;
	}

	.mana-option.unavailable .mana-count {
		color: #6b7280;
	}

	.mana-label {
		font-size: 0.75rem;
		color: #9ca3af;
		text-transform: uppercase;
		letter-spacing: 0.5px;
	}

	.no-mana-warning {
		text-align: center;
		padding: 2rem 1rem;
	}

	.warning-icon {
		font-size: 3rem;
		display: block;
		margin-bottom: 1rem;
		opacity: 0.6;
	}

	.no-mana-warning p {
		margin: 0 0 0.5rem;
		color: #9ca3af;
	}

	.no-mana-warning .hint {
		font-size: 0.875rem;
		color: #6b7280;
		font-style: italic;
	}

	.modal-footer {
		padding: 1rem 1.5rem;
		background: #141821;
		border-top: 1px solid #2a3441;
		border-radius: 0 0 10px 10px;
		display: flex;
		justify-content: center;
		gap: 1rem;
	}

	.btn-cancel {
		padding: 0.75rem 2rem;
		background: transparent;
		color: #9ca3af;
		border: 1px solid #374151;
		border-radius: 8px;
		font-size: 0.875rem;
		font-weight: 500;
		cursor: pointer;
		transition: all 0.2s ease;
	}

	.btn-cancel:hover:not(:disabled) {
		background: #374151;
		color: #e2e8f0;
		border-color: #4b5563;
	}

	.btn-cancel:disabled {
		opacity: 0.5;
		cursor: not-allowed;
	}

	/* Responsive */
	@media (max-width: 480px) {
		.mana-payment-modal {
			width: 95%;
			max-width: none;
			margin: 1rem;
		}

		.mana-grid {
			grid-template-columns: repeat(3, 1fr);
			gap: 0.5rem;
		}

		.mana-option {
			padding: 0.75rem 0.5rem;
		}

		.mana-symbol-wrapper {
			width: 40px;
			height: 40px;
		}

		.mana-count {
			font-size: 1.25rem;
		}
	}
</style>

