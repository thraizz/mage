<script lang="ts">
	import type { Snippet } from 'svelte';

	interface Props {
		title: string;
		description?: string;
		icon?: 'search' | 'table' | 'deck' | 'player' | 'game';
		children?: Snippet; // For action buttons
	}

	let { title, description = '', icon = 'search', children }: Props = $props();

	const icons: Record<string, string> = {
		search: 'M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z',
		table: 'M4 6h16M4 12h16M4 18h16',
		deck: 'M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10',
		player: 'M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z',
		game: 'M14.752 11.168l-3.197-2.132A1 1 0 0010 9.87v4.263a1 1 0 001.555.832l3.197-2.132a1 1 0 000-1.664z M21 12a9 9 0 11-18 0 9 9 0 0118 0z'
	};
</script>

<div class="empty-state">
	<div class="empty-icon">
		<svg
			viewBox="0 0 24 24"
			fill="none"
			stroke="currentColor"
			stroke-width="1.5"
			stroke-linecap="round"
			stroke-linejoin="round"
		>
			<path d={icons[icon]} />
		</svg>
	</div>
	<h3 class="empty-title">{title}</h3>
	{#if description}
		<p class="empty-description">{description}</p>
	{/if}
	{#if children}
		<div class="empty-actions">
			{@render children()}
		</div>
	{/if}
</div>

<style>
	.empty-state {
		display: flex;
		flex-direction: column;
		align-items: center;
		justify-content: center;
		padding: var(--space-12) var(--space-4);
		text-align: center;
	}

	.empty-icon {
		width: 4rem;
		height: 4rem;
		margin-bottom: var(--space-4);
		color: var(--text-ghost);
	}

	.empty-icon svg {
		width: 100%;
		height: 100%;
	}

	.empty-title {
		font-family: var(--font-display);
		font-size: var(--text-xl);
		font-weight: var(--weight-semibold);
		color: var(--text-muted);
		margin: 0 0 var(--space-2);
	}

	.empty-description {
		font-size: var(--text-base);
		color: var(--text-dim);
		margin: 0 0 var(--space-6);
		max-width: 24rem;
	}

	.empty-actions {
		display: flex;
		gap: var(--space-3);
	}
</style>
