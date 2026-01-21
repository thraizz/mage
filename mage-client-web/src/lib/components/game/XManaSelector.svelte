<script lang="ts">
	import { sendPlayerInteger, sendPlayerBoolean } from '$lib/api/game';
	import type { GamePlayXManaData } from '$lib/generated/mage/v1/websocket';

	// Props
	let {
		gameId,
		xManaData,
		onComplete = () => {},
		onCancel = () => {}
	}: {
		gameId: string;
		xManaData: GamePlayXManaData;
		onComplete?: () => void;
		onCancel?: () => void;
	} = $props();

	// State
	let selectedAmount = $state(0);
	let isLoading = $state(false);
	let error = $state<string | null>(null);

	// Maximum available mana from server
	const maxAmount = $derived(xManaData.available || 0);
	const message = $derived(xManaData.message || 'Choose X value');

	// Predefined quick-select amounts
	const quickAmounts = $derived.by(() => {
		const amounts: number[] = [];
		if (maxAmount >= 1) amounts.push(1);
		if (maxAmount >= 2) amounts.push(2);
		if (maxAmount >= 3) amounts.push(3);
		if (maxAmount >= 5) amounts.push(5);
		if (maxAmount >= 10) amounts.push(10);
		if (maxAmount > 10 && !amounts.includes(maxAmount)) amounts.push(maxAmount);
		return amounts;
	});

	/**
	 * Handle increment
	 */
	function increment() {
		if (selectedAmount < maxAmount) {
			selectedAmount++;
		}
	}

	/**
	 * Handle decrement
	 */
	function decrement() {
		if (selectedAmount > 0) {
			selectedAmount--;
		}
	}

	/**
	 * Handle quick amount selection
	 */
	function selectQuickAmount(amount: number) {
		selectedAmount = Math.min(amount, maxAmount);
	}

	/**
	 * Handle max selection
	 */
	function selectMax() {
		selectedAmount = maxAmount;
	}

	/**
	 * Handle confirm
	 */
	async function handleConfirm() {
		if (isLoading) return;

		isLoading = true;
		error = null;

		try {
			await sendPlayerInteger(gameId, selectedAmount);
			onComplete();
		} catch (err) {
			error = err instanceof Error ? err.message : 'Failed to confirm selection';
		} finally {
			isLoading = false;
		}
	}

	/**
	 * Handle cancel
	 */
	async function handleCancel() {
		if (isLoading) return;

		isLoading = true;
		error = null;

		try {
			// Send false to indicate cancellation (server interprets as X=0 or cancel)
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
		} else if (event.key === 'Enter' && selectedAmount > 0) {
			handleConfirm();
		} else if (event.key === 'ArrowUp') {
			event.preventDefault();
			increment();
		} else if (event.key === 'ArrowDown') {
			event.preventDefault();
			decrement();
		}
	}
</script>

<svelte:window onkeydown={handleKeydown} />

<div class="xmana-overlay" role="dialog" aria-labelledby="xmana-title" aria-modal="true">
	<div class="xmana-modal">
		<div class="modal-header">
			<h2 id="xmana-title">Choose X Value</h2>
		</div>

		<div class="modal-content">
			<p class="message">{message}</p>

			{#if error}
				<div class="error-message" role="alert">
					<span class="error-icon">⚠️</span>
					{error}
				</div>
			{/if}

			<div class="mana-available">
				<span class="label">Available mana:</span>
				<span class="amount">{maxAmount}</span>
			</div>

			<!-- Number selector -->
			<div class="number-selector">
				<button
					class="btn-adjust btn-decrement"
					onclick={decrement}
					disabled={selectedAmount <= 0 || isLoading}
					aria-label="Decrease X"
				>
					−
				</button>

				<div class="number-display">
					<span class="x-label">X =</span>
					<span class="number-value">{selectedAmount}</span>
				</div>

				<button
					class="btn-adjust btn-increment"
					onclick={increment}
					disabled={selectedAmount >= maxAmount || isLoading}
					aria-label="Increase X"
				>
					+
				</button>
			</div>

			<!-- Quick select buttons -->
			{#if quickAmounts().length > 0}
				<div class="quick-select">
					{#each quickAmounts() as amount}
						<button
							class="btn-quick"
							class:active={selectedAmount === amount}
							onclick={() => selectQuickAmount(amount)}
							disabled={isLoading}
						>
							{amount}
						</button>
					{/each}
					<button class="btn-quick btn-max" onclick={selectMax} disabled={isLoading}> MAX </button>
				</div>
			{/if}

			<!-- Slider for larger amounts -->
			{#if maxAmount > 5}
				<div class="slider-container">
					<input
						type="range"
						min="0"
						max={maxAmount}
						bind:value={selectedAmount}
						disabled={isLoading}
						class="mana-slider"
					/>
					<div class="slider-labels">
						<span>0</span>
						<span>{maxAmount}</span>
					</div>
				</div>
			{/if}
		</div>

		<div class="modal-footer">
			<button class="btn-cancel" onclick={handleCancel} disabled={isLoading}> Cancel </button>
			<button class="btn-confirm" onclick={handleConfirm} disabled={isLoading}>
				{isLoading ? 'Confirming...' : `Pay X = ${selectedAmount}`}
			</button>
		</div>
	</div>
</div>

<style>
	.xmana-overlay {
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

	.xmana-modal {
		background: #1a1f2e;
		border: 2px solid #667eea;
		border-radius: 12px;
		width: 90%;
		max-width: 400px;
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

	.message {
		margin: 0 0 1rem;
		font-size: 1rem;
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

	.mana-available {
		display: flex;
		justify-content: center;
		align-items: center;
		gap: 0.5rem;
		padding: 0.75rem 1rem;
		background: rgba(102, 126, 234, 0.1);
		border: 1px solid rgba(102, 126, 234, 0.3);
		border-radius: 8px;
		margin-bottom: 1.5rem;
	}

	.mana-available .label {
		color: #9ca3af;
		font-size: 0.875rem;
	}

	.mana-available .amount {
		font-size: 1.25rem;
		font-weight: 700;
		color: #667eea;
	}

	/* Number selector */
	.number-selector {
		display: flex;
		align-items: center;
		justify-content: center;
		gap: 1rem;
		margin-bottom: 1.5rem;
	}

	.btn-adjust {
		width: 48px;
		height: 48px;
		border-radius: 50%;
		background: #374151;
		border: 2px solid #4b5563;
		color: white;
		font-size: 1.5rem;
		font-weight: 700;
		cursor: pointer;
		transition: all 0.2s ease;
		display: flex;
		align-items: center;
		justify-content: center;
	}

	.btn-adjust:hover:not(:disabled) {
		background: #4b5563;
		border-color: #667eea;
		transform: scale(1.05);
	}

	.btn-adjust:disabled {
		opacity: 0.3;
		cursor: not-allowed;
	}

	.btn-decrement:active:not(:disabled) {
		background: #ef4444;
		border-color: #ef4444;
	}

	.btn-increment:active:not(:disabled) {
		background: #22c55e;
		border-color: #22c55e;
	}

	.number-display {
		display: flex;
		align-items: baseline;
		gap: 0.5rem;
		min-width: 100px;
		justify-content: center;
	}

	.x-label {
		font-size: 1.25rem;
		color: #9ca3af;
		font-weight: 500;
	}

	.number-value {
		font-size: 3rem;
		font-weight: 700;
		color: #fbbf24;
		line-height: 1;
	}

	/* Quick select */
	.quick-select {
		display: flex;
		flex-wrap: wrap;
		justify-content: center;
		gap: 0.5rem;
		margin-bottom: 1rem;
	}

	.btn-quick {
		min-width: 48px;
		padding: 0.5rem 0.75rem;
		background: #141821;
		border: 2px solid #2a3441;
		border-radius: 8px;
		color: #9ca3af;
		font-size: 0.875rem;
		font-weight: 600;
		cursor: pointer;
		transition: all 0.2s ease;
	}

	.btn-quick:hover:not(:disabled) {
		background: #2a3441;
		border-color: #667eea;
		color: white;
	}

	.btn-quick.active {
		background: #667eea;
		border-color: #667eea;
		color: white;
	}

	.btn-quick.btn-max {
		background: rgba(251, 191, 36, 0.15);
		border-color: rgba(251, 191, 36, 0.3);
		color: #fbbf24;
	}

	.btn-quick.btn-max:hover:not(:disabled) {
		background: rgba(251, 191, 36, 0.25);
		border-color: #fbbf24;
	}

	.btn-quick:disabled {
		opacity: 0.5;
		cursor: not-allowed;
	}

	/* Slider */
	.slider-container {
		margin-top: 0.5rem;
	}

	.mana-slider {
		width: 100%;
		height: 8px;
		border-radius: 4px;
		background: #141821;
		outline: none;
		-webkit-appearance: none;
		appearance: none;
	}

	.mana-slider::-webkit-slider-thumb {
		-webkit-appearance: none;
		width: 24px;
		height: 24px;
		border-radius: 50%;
		background: #667eea;
		border: 3px solid #fff;
		cursor: pointer;
		box-shadow: 0 2px 4px rgba(0, 0, 0, 0.3);
		transition: transform 0.2s ease;
	}

	.mana-slider::-webkit-slider-thumb:hover {
		transform: scale(1.1);
	}

	.mana-slider::-moz-range-thumb {
		width: 24px;
		height: 24px;
		border-radius: 50%;
		background: #667eea;
		border: 3px solid #fff;
		cursor: pointer;
		box-shadow: 0 2px 4px rgba(0, 0, 0, 0.3);
	}

	.slider-labels {
		display: flex;
		justify-content: space-between;
		margin-top: 0.25rem;
		font-size: 0.75rem;
		color: #6b7280;
	}

	/* Footer */
	.modal-footer {
		padding: 1rem 1.5rem;
		background: #141821;
		border-top: 1px solid #2a3441;
		border-radius: 0 0 10px 10px;
		display: flex;
		justify-content: space-between;
		gap: 1rem;
	}

	.btn-cancel {
		padding: 0.75rem 1.5rem;
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

	.btn-confirm {
		padding: 0.75rem 1.5rem;
		background: linear-gradient(135deg, #22c55e 0%, #16a34a 100%);
		color: white;
		border: none;
		border-radius: 8px;
		font-size: 0.875rem;
		font-weight: 600;
		cursor: pointer;
		transition: all 0.2s ease;
	}

	.btn-confirm:hover:not(:disabled) {
		transform: translateY(-1px);
		box-shadow: 0 4px 12px rgba(34, 197, 94, 0.3);
	}

	.btn-confirm:disabled,
	.btn-cancel:disabled {
		opacity: 0.5;
		cursor: not-allowed;
	}

	/* Responsive */
	@media (max-width: 480px) {
		.xmana-modal {
			width: 95%;
			margin: 1rem;
		}

		.number-value {
			font-size: 2.5rem;
		}

		.btn-adjust {
			width: 40px;
			height: 40px;
			font-size: 1.25rem;
		}

		.quick-select {
			gap: 0.375rem;
		}

		.btn-quick {
			min-width: 40px;
			padding: 0.375rem 0.5rem;
			font-size: 0.75rem;
		}
	}
</style>
