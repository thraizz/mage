<script lang="ts">
	import type { Snippet } from 'svelte';

	interface Props {
		variant?: 'primary' | 'secondary' | 'ghost' | 'danger';
		size?: 'sm' | 'md' | 'lg';
		disabled?: boolean;
		loading?: boolean;
		fullWidth?: boolean;
		type?: 'button' | 'submit' | 'reset';
		onclick?: (e: MouseEvent) => void;
		children?: Snippet;
	}

	let {
		variant = 'primary',
		size = 'md',
		disabled = false,
		loading = false,
		fullWidth = false,
		type = 'button',
		onclick,
		children
	}: Props = $props();
</script>

<button
	{type}
	class="btn btn-{variant} btn-{size}"
	class:btn-full={fullWidth}
	class:btn-loading={loading}
	disabled={disabled || loading}
	{onclick}
>
	{#if loading}
		<span class="btn-spinner"></span>
	{/if}
	<span class="btn-content" class:invisible={loading}>
		{#if children}
			{@render children()}
		{/if}
	</span>
</button>

<style>
	.btn {
		display: inline-flex;
		align-items: center;
		justify-content: center;
		gap: var(--space-2);
		font-family: var(--font-body);
		font-weight: var(--weight-medium);
		border: 1px solid transparent;
		border-radius: var(--radius-md);
		cursor: pointer;
		transition: all var(--transition-fast);
		position: relative;
		text-decoration: none;
	}

	.btn:focus-visible {
		outline: 2px solid var(--accent-gold);
		outline-offset: 2px;
	}

	.btn:disabled {
		opacity: 0.5;
		cursor: not-allowed;
	}

	/* Primary variant - gold accent */
	.btn-primary {
		background: var(--accent-gold);
		color: var(--bg-void);
		border-color: var(--accent-gold);
	}

	.btn-primary:hover:not(:disabled) {
		background: var(--accent-gold-bright);
		box-shadow: var(--shadow-glow);
	}

	.btn-primary:active:not(:disabled) {
		background: var(--accent-gold-dim);
	}

	/* Secondary variant - subtle */
	.btn-secondary {
		background: var(--bg-iron);
		color: var(--text-bright);
		border-color: var(--border-default);
	}

	.btn-secondary:hover:not(:disabled) {
		background: var(--bg-steel);
		border-color: var(--border-strong);
	}

	.btn-secondary:active:not(:disabled) {
		background: var(--bg-iron);
	}

	/* Ghost variant - minimal */
	.btn-ghost {
		background: transparent;
		color: var(--text-muted);
		border-color: transparent;
	}

	.btn-ghost:hover:not(:disabled) {
		background: var(--bg-iron);
		color: var(--text-bright);
	}

	.btn-ghost:active:not(:disabled) {
		background: var(--bg-slate);
	}

	/* Danger variant - destructive actions */
	.btn-danger {
		background: var(--status-error);
		color: white;
		border-color: var(--status-error);
	}

	.btn-danger:hover:not(:disabled) {
		background: #dc2626;
		border-color: #dc2626;
	}

	.btn-danger:active:not(:disabled) {
		background: #b91c1c;
	}

	/* Size variants */
	.btn-sm {
		padding: var(--space-1) var(--space-3);
		font-size: var(--text-sm);
		min-height: 32px;
	}

	.btn-md {
		padding: var(--space-2) var(--space-4);
		font-size: var(--text-base);
		min-height: 40px;
	}

	.btn-lg {
		padding: var(--space-3) var(--space-6);
		font-size: var(--text-lg);
		min-height: 48px;
	}

	.btn-full {
		width: 100%;
	}

	/* Loading state */
	.btn-loading {
		pointer-events: none;
	}

	.btn-spinner {
		position: absolute;
		width: 1.25em;
		height: 1.25em;
		border: 2px solid currentColor;
		border-top-color: transparent;
		border-radius: 50%;
		animation: spin 0.6s linear infinite;
	}

	.btn-content {
		display: inline-flex;
		align-items: center;
		gap: var(--space-2);
	}

	.invisible {
		visibility: hidden;
	}

	@keyframes spin {
		to {
			transform: rotate(360deg);
		}
	}
</style>
