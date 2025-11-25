<script lang="ts">
	interface Props {
		variant?: 'text' | 'circular' | 'rectangular';
		width?: string;
		height?: string;
		lines?: number;
	}

	let {
		variant = 'text',
		width = '100%',
		height = '',
		lines = 1
	}: Props = $props();

	const defaultHeight = $derived(
		variant === 'text' ? '1rem' :
		variant === 'circular' ? width :
		'3rem'
	);

	const actualHeight = $derived(height || defaultHeight);
</script>

{#if variant === 'text' && lines > 1}
	<div class="skeleton-lines">
		{#each Array(lines) as _, i}
			<div
				class="skeleton skeleton-text"
				style="width: {i === lines - 1 ? '80%' : width}; height: {actualHeight}"
			></div>
		{/each}
	</div>
{:else}
	<div
		class="skeleton skeleton-{variant}"
		style="width: {width}; height: {actualHeight}"
	></div>
{/if}

<style>
	.skeleton {
		background: linear-gradient(
			90deg,
			var(--bg-slate) 25%,
			var(--bg-iron) 50%,
			var(--bg-slate) 75%
		);
		background-size: 200% 100%;
		animation: shimmer 1.5s infinite;
	}

	.skeleton-text {
		border-radius: var(--radius-sm);
	}

	.skeleton-circular {
		border-radius: var(--radius-full);
	}

	.skeleton-rectangular {
		border-radius: var(--radius-md);
	}

	.skeleton-lines {
		display: flex;
		flex-direction: column;
		gap: var(--space-2);
	}

	@keyframes shimmer {
		0% {
			background-position: 200% 0;
		}
		100% {
			background-position: -200% 0;
		}
	}
</style>
