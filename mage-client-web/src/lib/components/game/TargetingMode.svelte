<script lang="ts">
	import { onMount, onDestroy } from 'svelte';
	import {
		targetingStore,
		isTargetingActive,
		targetingMessage,
		selectedTargetIds,
		canConfirmTargets,
		isTargetingRequired
	} from '$lib/stores/game-targeting';

	// Props
	let {
		onConfirm = () => {},
		onCancel = () => {}
	}: {
		// eslint-disable-next-line no-unused-vars
		onConfirm?: (targetIds: string[]) => void;
		onCancel?: () => void;
	} = $props();

	// Derived state from stores
	const isActive = $derived($isTargetingActive);
	const message = $derived($targetingMessage);
	const selectedIds = $derived($selectedTargetIds);
	const canConfirm = $derived($canConfirmTargets);
	const required = $derived($isTargetingRequired);

	// Get targeting state for display
	const targetingState = $derived.by(() => {
		let state: { isActive: boolean; minTargets: number; maxTargets: number } = {
			isActive: false,
			minTargets: 1,
			maxTargets: 1
		};
		targetingStore.subscribe((s) => {
			state = { isActive: s.isActive, minTargets: s.minTargets, maxTargets: s.maxTargets };
		})();
		return state;
	});

	/**
	 * Handle keyboard shortcuts
	 */
	function handleKeydown(event: KeyboardEvent) {
		if (!isActive) return;

		switch (event.key) {
			case 'Escape':
				if (!required) {
					event.preventDefault();
					handleCancel();
				}
				break;
			case 'Enter':
				if (canConfirm) {
					event.preventDefault();
					handleConfirm();
				}
				break;
		}
	}

	/**
	 * Handle confirm button click
	 */
	function handleConfirm() {
		if (!canConfirm) return;
		const targets = targetingStore.getSelectedTargets();
		onConfirm(targets);
		targetingStore.exitTargetingMode();
	}

	/**
	 * Handle cancel button click
	 */
	function handleCancel() {
		if (required) return;
		onCancel();
		targetingStore.exitTargetingMode();
	}

	/**
	 * Handle backdrop click - only cancel if not required
	 */
	function handleBackdropClick(event: MouseEvent) {
		// Only trigger if clicking directly on the backdrop, not children
		if (event.target === event.currentTarget && !required) {
			handleCancel();
		}
	}

	/**
	 * Handle right-click to cancel
	 */
	function handleContextMenu(event: MouseEvent) {
		if (isActive && !required) {
			event.preventDefault();
			handleCancel();
		}
	}

	// Add global keyboard listener
	onMount(() => {
		window.addEventListener('keydown', handleKeydown);
		window.addEventListener('contextmenu', handleContextMenu);
	});

	onDestroy(() => {
		window.removeEventListener('keydown', handleKeydown);
		window.removeEventListener('contextmenu', handleContextMenu);
	});
</script>

{#if isActive}
	<!-- Targeting Mode Overlay -->
	<div
		class="targeting-overlay"
		onclick={handleBackdropClick}
		role="dialog"
		aria-modal="true"
		aria-labelledby="targeting-title"
	>
		<!-- Top Banner -->
		<div class="targeting-banner">
			<div class="targeting-icon">🎯</div>
			<div class="targeting-content">
				<h3 id="targeting-title" class="targeting-title">Select Target</h3>
				<p class="targeting-message">{message}</p>
				<div class="targeting-status">
					{#if targetingState.maxTargets > 1}
						<span class="target-count">
							{selectedIds.length} / {targetingState.maxTargets} targets selected
						</span>
					{:else}
						<span class="target-hint">Click a highlighted card to select it as target</span>
					{/if}
				</div>
			</div>
		</div>

		<!-- Bottom Action Bar -->
		<div class="targeting-actions">
			<div class="action-hints">
				{#if !required}
					<span class="hint">
						<kbd>ESC</kbd> or <kbd>Right-click</kbd> to cancel
					</span>
				{/if}
				{#if canConfirm}
					<span class="hint">
						<kbd>Enter</kbd> to confirm
					</span>
				{/if}
			</div>
			<div class="action-buttons">
				{#if !required}
					<button class="btn-cancel" onclick={handleCancel} type="button">
						Cancel
					</button>
				{/if}
				<button
					class="btn-confirm"
					onclick={handleConfirm}
					disabled={!canConfirm}
					type="button"
				>
					{#if targetingState.maxTargets > 1}
						Confirm ({selectedIds.length})
					{:else}
						Confirm
					{/if}
				</button>
			</div>
		</div>
	</div>
{/if}

<style>
	/* Overlay - semi-transparent with crosshair cursor */
	.targeting-overlay {
		position: fixed;
		inset: 0;
		background: rgba(0, 0, 0, 0.4);
		z-index: 100;
		cursor: crosshair;
		display: flex;
		flex-direction: column;
		justify-content: space-between;
		pointer-events: auto;
	}

	/* Top Banner */
	.targeting-banner {
		display: flex;
		align-items: center;
		gap: 1rem;
		padding: 1rem 1.5rem;
		background: linear-gradient(180deg, rgba(251, 191, 36, 0.15) 0%, transparent 100%);
		border-bottom: 2px solid rgba(251, 191, 36, 0.5);
		pointer-events: auto;
		animation: banner-slide-in 0.2s ease-out;
	}

	@keyframes banner-slide-in {
		from {
			opacity: 0;
			transform: translateY(-20px);
		}
		to {
			opacity: 1;
			transform: translateY(0);
		}
	}

	.targeting-icon {
		font-size: 2rem;
		animation: pulse-icon 1.5s ease-in-out infinite;
	}

	@keyframes pulse-icon {
		0%,
		100% {
			transform: scale(1);
		}
		50% {
			transform: scale(1.1);
		}
	}

	.targeting-content {
		flex: 1;
	}

	.targeting-title {
		margin: 0 0 0.25rem 0;
		font-size: 1.125rem;
		font-weight: 700;
		color: #fbbf24;
		text-transform: uppercase;
		letter-spacing: 1px;
	}

	.targeting-message {
		margin: 0 0 0.5rem 0;
		font-size: 1rem;
		color: #ffffff;
	}

	.targeting-status {
		font-size: 0.875rem;
	}

	.target-count {
		color: #22c55e;
		font-weight: 600;
	}

	.target-hint {
		color: #94a3b8;
		font-style: italic;
	}

	/* Bottom Action Bar */
	.targeting-actions {
		display: flex;
		justify-content: space-between;
		align-items: center;
		padding: 1rem 1.5rem;
		background: linear-gradient(0deg, rgba(0, 0, 0, 0.8) 0%, transparent 100%);
		pointer-events: auto;
	}

	.action-hints {
		display: flex;
		gap: 1.5rem;
	}

	.hint {
		color: #6b7280;
		font-size: 0.75rem;
	}

	kbd {
		display: inline-block;
		padding: 0.125rem 0.375rem;
		background: #374151;
		border: 1px solid #4b5563;
		border-radius: 4px;
		font-family: monospace;
		font-size: 0.6875rem;
		color: #e5e7eb;
	}

	.action-buttons {
		display: flex;
		gap: 0.75rem;
	}

	.btn-cancel,
	.btn-confirm {
		padding: 0.625rem 1.25rem;
		border-radius: 6px;
		font-weight: 600;
		font-size: 0.875rem;
		cursor: pointer;
		transition: all 0.2s;
		border: none;
	}

	.btn-cancel {
		background: #374151;
		color: #e5e7eb;
	}

	.btn-cancel:hover {
		background: #4b5563;
	}

	.btn-confirm {
		background: #fbbf24;
		color: #0a0d12;
	}

	.btn-confirm:hover:not(:disabled) {
		background: #f59e0b;
	}

	.btn-confirm:disabled {
		background: #6b7280;
		color: #9ca3af;
		cursor: not-allowed;
	}

	/* Responsive */
	@media (max-width: 600px) {
		.targeting-banner {
			padding: 0.75rem 1rem;
		}

		.targeting-icon {
			font-size: 1.5rem;
		}

		.targeting-title {
			font-size: 1rem;
		}

		.targeting-message {
			font-size: 0.875rem;
		}

		.targeting-actions {
			flex-direction: column;
			gap: 0.75rem;
		}

		.action-hints {
			flex-wrap: wrap;
			justify-content: center;
		}
	}
</style>

