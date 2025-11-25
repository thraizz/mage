<script lang="ts">
	type ManaColor = 'W' | 'U' | 'B' | 'R' | 'G' | 'C' | 'X';

	interface Props {
		symbol: ManaColor | number;
		size?: 'sm' | 'md' | 'lg' | 'xl';
		interactive?: boolean;
		onclick?: () => void;
	}

	let { symbol, size = 'md', interactive = false, onclick }: Props = $props();

	const colorStyles: Record<string, { bg: string; text: string; border: string }> = {
		W: { bg: '#f8f6e3', text: '#1a1a1a', border: '#c9c5a8' },
		U: { bg: '#0e68ab', text: '#ffffff', border: '#0a4d80' },
		B: { bg: '#150b00', text: '#a8a8a8', border: '#3d2d1a' },
		R: { bg: '#d3202a', text: '#ffffff', border: '#a01920' },
		G: { bg: '#00733e', text: '#ffffff', border: '#005830' },
		C: { bg: '#9ca3af', text: '#1a1a1a', border: '#6b7280' },
		X: { bg: '#6b7280', text: '#ffffff', border: '#4b5563' }
	};

	const isNumeric = $derived(typeof symbol === 'number');
	const displaySymbol = $derived(isNumeric ? symbol.toString() : symbol);
	const style = $derived(isNumeric ? colorStyles['C'] : (colorStyles[symbol as string] || colorStyles['C']));
</script>

{#if interactive}
	<button
		type="button"
		class="mana-symbol mana-symbol-{size} mana-interactive"
		style="--mana-bg: {style.bg}; --mana-text: {style.text}; --mana-border: {style.border}"
		{onclick}
	>
		{displaySymbol}
	</button>
{:else}
	<span
		class="mana-symbol mana-symbol-{size}"
		style="--mana-bg: {style.bg}; --mana-text: {style.text}; --mana-border: {style.border}"
	>
		{displaySymbol}
	</span>
{/if}

<style>
	.mana-symbol {
		display: inline-flex;
		align-items: center;
		justify-content: center;
		font-family: var(--font-body);
		font-weight: var(--weight-bold);
		border-radius: var(--radius-full);
		background: var(--mana-bg);
		color: var(--mana-text);
		border: 2px solid var(--mana-border);
		box-shadow:
			inset 0 2px 4px rgba(255, 255, 255, 0.2),
			inset 0 -2px 4px rgba(0, 0, 0, 0.3),
			0 1px 2px rgba(0, 0, 0, 0.3);
		user-select: none;
		flex-shrink: 0;
	}

	/* Sizes */
	.mana-symbol-sm {
		width: 1rem;
		height: 1rem;
		font-size: 0.625rem;
		border-width: 1px;
	}

	.mana-symbol-md {
		width: 1.25rem;
		height: 1.25rem;
		font-size: 0.75rem;
	}

	.mana-symbol-lg {
		width: 1.5rem;
		height: 1.5rem;
		font-size: 0.875rem;
	}

	.mana-symbol-xl {
		width: 2rem;
		height: 2rem;
		font-size: 1rem;
		border-width: 3px;
	}

	/* Interactive state */
	.mana-interactive {
		cursor: pointer;
		transition: all var(--transition-fast);
	}

	.mana-interactive:hover {
		transform: scale(1.1);
		box-shadow:
			inset 0 2px 4px rgba(255, 255, 255, 0.2),
			inset 0 -2px 4px rgba(0, 0, 0, 0.3),
			0 2px 8px rgba(0, 0, 0, 0.4);
	}

	.mana-interactive:active {
		transform: scale(0.95);
	}

	.mana-interactive:focus-visible {
		outline: 2px solid var(--accent-gold);
		outline-offset: 2px;
	}
</style>
