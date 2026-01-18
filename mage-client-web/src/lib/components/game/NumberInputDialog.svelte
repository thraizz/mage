<script lang="ts">
	interface Props {
		title: string;
		defaultValue?: number;
		min?: number;
		max?: number;
		onConfirm: (value: number) => void;
		onCancel: () => void;
	}

	let { title, defaultValue = 1, min = 1, max = 99, onConfirm, onCancel }: Props = $props();

	let value = $state(defaultValue);

	function increment() {
		if (max === undefined || value < max) {
			value++;
		}
	}

	function decrement() {
		if (value > min) {
			value--;
		}
	}

	function handleConfirm() {
		onConfirm(value);
	}

	function handleKeydown(e: KeyboardEvent) {
		if (e.key === 'Enter') {
			handleConfirm();
		} else if (e.key === 'Escape') {
			onCancel();
		}
	}
</script>

<svelte:window onkeydown={handleKeydown} />

<div class="overlay" role="dialog" aria-labelledby="number-dialog-title">
	<div class="dialog">
		<div class="dialog-header">
			<h2 id="number-dialog-title">{title}</h2>
			<button class="close-button" onclick={onCancel} aria-label="Close dialog">×</button>
		</div>

		<div class="number-input-section">
			<button
				class="btn-adjust"
				onclick={decrement}
				disabled={value <= min}
				aria-label="Decrease value"
			>
				−
			</button>

			<input
				type="number"
				bind:value
				{min}
				{max}
				class="number-input"
				autofocus
			/>

			<button
				class="btn-adjust"
				onclick={increment}
				disabled={max !== undefined && value >= max}
				aria-label="Increase value"
			>
				+
			</button>
		</div>

		<div class="dialog-footer">
			<button class="btn-secondary" onclick={onCancel}>Cancel</button>
			<button class="btn-primary" onclick={handleConfirm}>Confirm</button>
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
		backdrop-filter: blur(4px);
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
		max-width: 400px;
		width: 90%;
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

	.number-input-section {
		display: flex;
		align-items: center;
		justify-content: center;
		gap: 1rem;
		margin-bottom: 1.5rem;
	}

	.btn-adjust {
		width: 3rem;
		height: 3rem;
		background: rgba(102, 126, 234, 0.1);
		border: 2px solid #667eea;
		border-radius: 8px;
		color: #667eea;
		font-size: 1.5rem;
		font-weight: 700;
		cursor: pointer;
		transition: all 0.2s;
		display: flex;
		align-items: center;
		justify-content: center;
	}

	.btn-adjust:hover:not(:disabled) {
		background: rgba(102, 126, 234, 0.2);
		transform: translateY(-1px);
	}

	.btn-adjust:disabled {
		opacity: 0.3;
		cursor: not-allowed;
	}

	.number-input {
		width: 120px;
		padding: 0.75rem;
		background: rgba(255, 255, 255, 0.05);
		border: 2px solid #3a4451;
		border-radius: 8px;
		color: #fff;
		font-size: 2rem;
		font-weight: 700;
		text-align: center;
		transition: border-color 0.2s;
	}

	.number-input:focus {
		outline: none;
		border-color: #667eea;
	}

	/* Hide number input spinners */
	.number-input::-webkit-outer-spin-button,
	.number-input::-webkit-inner-spin-button {
		-webkit-appearance: none;
		margin: 0;
	}

	.number-input[type='number'] {
		-moz-appearance: textfield;
	}

	.dialog-footer {
		display: flex;
		justify-content: flex-end;
		gap: 0.75rem;
		padding-top: 1rem;
		border-top: 1px solid #3a4451;
	}

	.btn-primary,
	.btn-secondary {
		padding: 0.5rem 1rem;
		border-radius: 6px;
		font-weight: 600;
		font-size: 0.875rem;
		cursor: pointer;
		transition: all 0.2s;
		border: none;
	}

	.btn-primary {
		background: #667eea;
		color: white;
	}

	.btn-primary:hover {
		background: #5568d3;
		transform: translateY(-1px);
	}

	.btn-secondary {
		background: rgba(255, 255, 255, 0.1);
		color: #fff;
		border: 1px solid #3a4451;
	}

	.btn-secondary:hover {
		background: rgba(255, 255, 255, 0.15);
	}
</style>
