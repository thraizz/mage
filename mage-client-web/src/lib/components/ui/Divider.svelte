<script lang="ts">
	interface Props {
		orientation?: 'horizontal' | 'vertical';
		spacing?: 'none' | 'sm' | 'md' | 'lg';
		label?: string;
	}

	let { orientation = 'horizontal', spacing = 'md', label = '' }: Props = $props();
</script>

{#if orientation === 'horizontal'}
	<div class="divider divider-horizontal divider-spacing-{spacing}" role="separator">
		{#if label}
			<span class="divider-line"></span>
			<span class="divider-label">{label}</span>
			<span class="divider-line"></span>
		{:else}
			<span class="divider-line divider-line-full"></span>
		{/if}
	</div>
{:else}
	<div class="divider divider-vertical divider-spacing-{spacing}" role="separator"></div>
{/if}

<style>
	.divider {
		display: flex;
		align-items: center;
	}

	/* Horizontal */
	.divider-horizontal {
		width: 100%;
	}

	.divider-horizontal .divider-line {
		flex: 1;
		height: 1px;
		background: var(--border-subtle);
	}

	.divider-horizontal .divider-line-full {
		flex: 1;
	}

	.divider-label {
		padding: 0 var(--space-3);
		font-size: var(--text-sm);
		color: var(--text-dim);
		white-space: nowrap;
	}

	/* Vertical */
	.divider-vertical {
		width: 1px;
		background: var(--border-subtle);
		align-self: stretch;
	}

	/* Spacing */
	.divider-spacing-none {
		margin: 0;
	}

	.divider-spacing-sm {
		margin: var(--space-2) 0;
	}

	.divider-horizontal.divider-spacing-md {
		margin: var(--space-4) 0;
	}

	.divider-horizontal.divider-spacing-lg {
		margin: var(--space-6) 0;
	}

	.divider-vertical.divider-spacing-sm {
		margin: 0 var(--space-2);
	}

	.divider-vertical.divider-spacing-md {
		margin: 0 var(--space-4);
	}

	.divider-vertical.divider-spacing-lg {
		margin: 0 var(--space-6);
	}
</style>
