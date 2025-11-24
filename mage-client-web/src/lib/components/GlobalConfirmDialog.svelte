<script lang="ts">
	import { confirm } from '$lib/stores/confirm';
	import Modal from './Modal.svelte';

	$: open = $confirm.open;
	$: title = $confirm.title || 'Confirm';
	$: message = $confirm.message || 'Are you sure?';
	$: confirmText = $confirm.confirmText || 'Confirm';
	$: cancelText = $confirm.cancelText || 'Cancel';
	$: destructive = $confirm.destructive || false;

	function handleConfirm() {
		confirm.handleConfirm();
	}

	function handleCancel() {
		confirm.handleCancel();
	}

	// Handle keyboard shortcuts
	function handleKeydown(event: KeyboardEvent) {
		if (!open) return;

		if (event.key === 'Enter') {
			event.preventDefault();
			handleConfirm();
		} else if (event.key === 'Escape') {
			event.preventDefault();
			handleCancel();
		}
	}
</script>

<svelte:window on:keydown={handleKeydown} />

<Modal {open} {title} size="small" closeOnBackdrop={false} onClose={handleCancel}>
	<div class="confirm-message">
		{message}
	</div>

	<div slot="footer" class="confirm-actions">
		<button class="btn-cancel" on:click={handleCancel}>
			{cancelText}
		</button>
		<button class="btn-confirm" class:destructive on:click={handleConfirm}>
			{confirmText}
		</button>
	</div>
</Modal>

<style>
	.confirm-message {
		color: #374151;
		font-size: 0.9375rem;
		line-height: 1.6;
		margin: 0;
	}

	.confirm-actions {
		display: flex;
		gap: 0.75rem;
		justify-content: flex-end;
		width: 100%;
	}

	.btn-cancel,
	.btn-confirm {
		padding: 0.625rem 1.25rem;
		border-radius: 0.5rem;
		font-weight: 600;
		font-size: 0.875rem;
		cursor: pointer;
		transition: all 0.2s;
		border: none;
	}

	.btn-cancel {
		background-color: #e5e7eb;
		color: #374151;
	}

	.btn-cancel:hover {
		background-color: #d1d5db;
	}

	.btn-confirm {
		background-color: #667eea;
		color: white;
	}

	.btn-confirm:hover {
		background-color: #5568d3;
	}

	.btn-confirm.destructive {
		background-color: #ef4444;
	}

	.btn-confirm.destructive:hover {
		background-color: #dc2626;
	}
</style>
