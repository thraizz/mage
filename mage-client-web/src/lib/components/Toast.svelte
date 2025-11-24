<script lang="ts">
	import { fly } from 'svelte/transition';
	import { toast } from '$lib/stores/toast';
	import type { Toast as ToastType } from '$lib/types/toast';

	export let data: ToastType;

	function handleDismiss() {
		toast.dismiss(data.id);
	}

	function getIcon(type: ToastType['type']): string {
		switch (type) {
			case 'success':
				return '✓';
			case 'error':
				return '✕';
			case 'warning':
				return '⚠';
			case 'info':
			default:
				return 'ℹ';
		}
	}

	function getColorClass(type: ToastType['type']): string {
		switch (type) {
			case 'success':
				return 'toast-success';
			case 'error':
				return 'toast-error';
			case 'warning':
				return 'toast-warning';
			case 'info':
			default:
				return 'toast-info';
		}
	}
</script>

<div
	class="toast {getColorClass(data.type)}"
	role="alert"
	aria-live="polite"
	transition:fly={{ x: 300, duration: 300 }}
>
	<div class="toast-icon">
		{getIcon(data.type)}
	</div>

	<div class="toast-message">
		{data.message}
	</div>

	{#if data.dismissible}
		<button class="toast-dismiss" on:click={handleDismiss} aria-label="Dismiss notification">
			<svg
				xmlns="http://www.w3.org/2000/svg"
				width="16"
				height="16"
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

<style>
	.toast {
		display: flex;
		align-items: center;
		gap: 0.75rem;
		padding: 1rem;
		margin-bottom: 0.75rem;
		border-radius: 0.5rem;
		box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
		min-width: 300px;
		max-width: 500px;
		background-color: white;
		border-left: 4px solid;
		position: relative;
	}

	.toast-icon {
		flex-shrink: 0;
		width: 24px;
		height: 24px;
		display: flex;
		align-items: center;
		justify-content: center;
		font-weight: 700;
		font-size: 1.25rem;
		border-radius: 50%;
		color: white;
	}

	.toast-message {
		flex: 1;
		font-size: 0.875rem;
		line-height: 1.5;
		color: #374151;
		word-wrap: break-word;
	}

	.toast-dismiss {
		flex-shrink: 0;
		background: none;
		border: none;
		cursor: pointer;
		padding: 0.25rem;
		color: #9ca3af;
		transition: color 0.2s;
		display: flex;
		align-items: center;
		justify-content: center;
		border-radius: 0.25rem;
	}

	.toast-dismiss:hover {
		color: #4b5563;
		background-color: #f3f4f6;
	}

	/* Success Toast */
	.toast-success {
		border-left-color: #10b981;
	}

	.toast-success .toast-icon {
		background-color: #10b981;
	}

	/* Error Toast */
	.toast-error {
		border-left-color: #ef4444;
	}

	.toast-error .toast-icon {
		background-color: #ef4444;
	}

	/* Warning Toast */
	.toast-warning {
		border-left-color: #f59e0b;
	}

	.toast-warning .toast-icon {
		background-color: #f59e0b;
	}

	/* Info Toast */
	.toast-info {
		border-left-color: #3b82f6;
	}

	.toast-info .toast-icon {
		background-color: #3b82f6;
	}

	/* Responsive */
	@media (max-width: 640px) {
		.toast {
			min-width: 280px;
			max-width: calc(100vw - 2rem);
		}
	}
</style>
