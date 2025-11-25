<script lang="ts">
	interface Props {
		format: string;
		size?: 'sm' | 'md' | 'lg';
	}

	let { format, size = 'md' }: Props = $props();

	const formatColors: Record<string, string> = {
		standard: '#f59e0b',
		pioneer: '#8b5cf6',
		modern: '#ef4444',
		legacy: '#3b82f6',
		vintage: '#10b981',
		commander: '#ec4899',
		edh: '#ec4899',
		pauper: '#71717a',
		limited: '#06b6d4',
		draft: '#06b6d4',
		sealed: '#14b8a6',
		historic: '#f97316',
		alchemy: '#a855f7',
		brawl: '#f472b6',
		freeform: '#6b7280',
		twoplayerduel: '#3b82f6'
	};

	const normalizedFormat = $derived(format.toLowerCase().replace(/[^a-z]/g, ''));
	const color = $derived(formatColors[normalizedFormat] || '#71717a');
	const displayName = $derived(format.charAt(0).toUpperCase() + format.slice(1));
</script>

<span class="format-badge format-badge-{size}" style="--format-color: {color}">
	{displayName}
</span>

<style>
	.format-badge {
		display: inline-flex;
		align-items: center;
		font-family: var(--font-body);
		font-weight: var(--weight-semibold);
		text-transform: uppercase;
		letter-spacing: 0.05em;
		border-radius: var(--radius-sm);
		background: color-mix(in srgb, var(--format-color) 15%, transparent);
		color: var(--format-color);
		border: 1px solid color-mix(in srgb, var(--format-color) 30%, transparent);
		white-space: nowrap;
	}

	/* Sizes */
	.format-badge-sm {
		padding: 0.125rem var(--space-2);
		font-size: 0.625rem;
	}

	.format-badge-md {
		padding: var(--space-1) var(--space-2);
		font-size: var(--text-xs);
	}

	.format-badge-lg {
		padding: var(--space-1) var(--space-3);
		font-size: var(--text-sm);
	}
</style>
