<script lang="ts">
	import { fade, scale } from 'svelte/transition';
	import { onMount, onDestroy } from 'svelte';
	import type { ModalSize } from '$lib/types/modal';

	// Props
	export let open = false;
	export let title: string | undefined = undefined;
	export let size: ModalSize = 'medium';
	export let closeOnBackdrop = true;
	export let closeOnEsc = true;
	export let showCloseButton = true;
	export let onClose: (() => void) | undefined = undefined;

	let modalElement: HTMLDivElement;
	let previouslyFocusedElement: HTMLElement | null = null;

	// Handle ESC key press
	function handleKeydown(event: KeyboardEvent) {
		if (event.key === 'Escape' && closeOnEsc && open) {
			handleClose();
		}
	}

	// Handle backdrop click
	function handleBackdropClick(event: MouseEvent) {
		if (closeOnBackdrop && event.target === event.currentTarget) {
			handleClose();
		}
	}

	// Close modal
	function handleClose() {
		open = false;
		if (onClose) {
			onClose();
		}
	}

	// Prevent body scroll when modal is open
	function updateBodyScroll() {
		if (typeof window === 'undefined') return;

		if (open) {
			// Save currently focused element
			previouslyFocusedElement = document.activeElement as HTMLElement;

			// Prevent body scroll
			document.body.style.overflow = 'hidden';

			// Focus the modal
			if (modalElement) {
				modalElement.focus();
			}
		} else {
			// Restore body scroll
			document.body.style.overflow = '';

			// Restore focus to previously focused element
			if (previouslyFocusedElement) {
				previouslyFocusedElement.focus();
			}
		}
	}

	// Focus trap: keep focus within modal
	function handleFocusTrap(event: FocusEvent) {
		if (!open || !modalElement) return;

		const focusableElements = modalElement.querySelectorAll(
			'button, [href], input, select, textarea, [tabindex]:not([tabindex="-1"])'
		);
		const firstElement = focusableElements[0] as HTMLElement;

		// If focus moves outside modal, bring it back
		if (!modalElement.contains(event.target as Node)) {
			firstElement?.focus();
		}
	}

	// Watch for open state changes
	$: if (typeof window !== 'undefined') {
		updateBodyScroll();
	}

	onMount(() => {
		document.addEventListener('keydown', handleKeydown);
		document.addEventListener('focusin', handleFocusTrap);
	});

	onDestroy(() => {
		if (typeof window !== 'undefined') {
			document.removeEventListener('keydown', handleKeydown);
			document.removeEventListener('focusin', handleFocusTrap);
			// Ensure body scroll is restored
			document.body.style.overflow = '';
		}
	});

	function getSizeClass(): string {
		switch (size) {
			case 'small':
				return 'modal-small';
			case 'large':
				return 'modal-large';
			case 'medium':
			default:
				return 'modal-medium';
		}
	}
</script>

{#if open}
	<!-- Backdrop -->
	<div
		class="modal-backdrop"
		on:click={handleBackdropClick}
		transition:fade={{ duration: 200 }}
		role="presentation"
	>
		<!-- Modal Dialog -->
		<div
			bind:this={modalElement}
			class="modal-dialog {getSizeClass()}"
			role="dialog"
			aria-modal="true"
			aria-labelledby={title ? 'modal-title' : undefined}
			tabindex="-1"
			transition:scale={{ duration: 200, start: 0.95 }}
		>
			<!-- Header -->
			{#if title || showCloseButton}
				<div class="modal-header">
					{#if title}
						<h2 id="modal-title" class="modal-title">{title}</h2>
					{/if}

					{#if showCloseButton}
						<button
							class="modal-close-button"
							on:click={handleClose}
							aria-label="Close modal"
							type="button"
						>
							<svg
								xmlns="http://www.w3.org/2000/svg"
								width="20"
								height="20"
								viewBox="0 0 24 24"
								fill="none"
								stroke="currentColor"
								stroke-width="2"
								stroke-linecap="round"
								stroke-linejoin="round"
							>
								<line x1="18" y1="6" x2="6" y2="18"></line>
								<line x1="6" y1="6" x2="18" y2="18"></line>
							</svg>
						</button>
					{/if}
				</div>
			{/if}

			<!-- Content -->
			<div class="modal-content">
				<slot />
			</div>

			<!-- Footer (optional slot) -->
			{#if $$slots.footer}
				<div class="modal-footer">
					<slot name="footer" />
				</div>
			{/if}
		</div>
	</div>
{/if}

<style>
	.modal-backdrop {
		position: fixed;
		top: 0;
		left: 0;
		right: 0;
		bottom: 0;
		background-color: rgba(0, 0, 0, 0.5);
		display: flex;
		align-items: center;
		justify-content: center;
		z-index: 1000;
		padding: 1rem;
	}

	.modal-dialog {
		background-color: white;
		border-radius: 0.75rem;
		box-shadow:
			0 20px 25px -5px rgba(0, 0, 0, 0.1),
			0 10px 10px -5px rgba(0, 0, 0, 0.04);
		display: flex;
		flex-direction: column;
		max-height: calc(100vh - 2rem);
		width: 100%;
		position: relative;
	}

	.modal-dialog:focus {
		outline: none;
	}

	/* Size variants */
	.modal-small {
		max-width: 400px;
	}

	.modal-medium {
		max-width: 600px;
	}

	.modal-large {
		max-width: 800px;
	}

	/* Header */
	.modal-header {
		display: flex;
		align-items: center;
		justify-content: space-between;
		padding: 1.5rem;
		border-bottom: 1px solid #e5e7eb;
	}

	.modal-title {
		margin: 0;
		font-size: 1.25rem;
		font-weight: 600;
		color: #111827;
		flex: 1;
	}

	.modal-close-button {
		background: none;
		border: none;
		cursor: pointer;
		padding: 0.5rem;
		color: #6b7280;
		transition: all 0.2s;
		border-radius: 0.375rem;
		display: flex;
		align-items: center;
		justify-content: center;
		margin-left: 1rem;
	}

	.modal-close-button:hover {
		color: #111827;
		background-color: #f3f4f6;
	}

	/* Content */
	.modal-content {
		flex: 1;
		padding: 1.5rem;
		overflow-y: auto;
	}

	/* Footer */
	.modal-footer {
		padding: 1.5rem;
		border-top: 1px solid #e5e7eb;
		display: flex;
		gap: 0.75rem;
		justify-content: flex-end;
	}

	/* Responsive */
	@media (max-width: 640px) {
		.modal-dialog {
			max-width: 100%;
			margin: 0;
			max-height: 100vh;
			border-radius: 0;
		}

		.modal-backdrop {
			padding: 0;
		}
	}
</style>
