<script lang="ts">
	type ManaColor = 'W' | 'U' | 'B' | 'R' | 'G' | 'C' | 'X' | 'T' | 'Q';

	interface Props {
		symbol: ManaColor | number | string;
		size?: 'sm' | 'md' | 'lg' | 'xl' | '2x' | '3x';
		interactive?: boolean;
		onclick?: () => void;
	}

	let { symbol, size = 'md', interactive = false, onclick }: Props = $props();

	// Convert symbol to mana-font class
	function getManaClass(sym: ManaColor | number | string): string {
		if (typeof sym === 'number') {
			return `ms-${sym}`;
		}
		// Handle string symbols - normalize to lowercase for mana-font
		const s = String(sym).toLowerCase();
		switch (s) {
			case 'w':
				return 'ms-w';
			case 'u':
				return 'ms-u';
			case 'b':
				return 'ms-b';
			case 'r':
				return 'ms-r';
			case 'g':
				return 'ms-g';
			case 'c':
				return 'ms-c';
			case 'x':
				return 'ms-x';
			case 't':
				return 'ms-tap';
			case 'q':
				return 'ms-untap';
			case 's':
				return 'ms-s'; // snow
			case 'e':
				return 'ms-e'; // energy
			default:
				// Try parsing as number
				const num = parseInt(s, 10);
				if (!isNaN(num)) {
					return `ms-${num}`;
				}
				return 'ms-c'; // fallback to colorless
		}
	}

	const manaClass = $derived(getManaClass(symbol));

	// Size classes - mana-font uses ms-cost for circle backgrounds, ms-shadow for shadows
	const sizeClass = $derived(
		{
			sm: 'mana-size-sm',
			md: 'mana-size-md',
			lg: 'mana-size-lg',
			xl: 'mana-size-xl',
			'2x': 'ms-2x',
			'3x': 'ms-3x'
		}[size] || 'mana-size-md'
	);
</script>

{#if interactive}
	<button type="button" class="mana-symbol-wrapper mana-interactive {sizeClass}" {onclick}>
		<i class="ms ms-cost {manaClass}"></i>
		<span class="sr-only">Mana symbol: {symbol}</span>
	</button>
{:else}
	<span class="mana-symbol-wrapper {sizeClass}">
		<i class="ms ms-cost {manaClass}"></i>
	</span>
{/if}

<style>
	.mana-symbol-wrapper {
		display: inline-flex;
		align-items: center;
		justify-content: center;
		vertical-align: middle;
	}

	/* Size overrides for our custom sizes */
	.mana-size-sm :global(.ms) {
		font-size: 0.875rem;
	}

	.mana-size-md :global(.ms) {
		font-size: 1.125rem;
	}

	.mana-size-lg :global(.ms) {
		font-size: 1.5rem;
	}

	.mana-size-xl :global(.ms) {
		font-size: 2rem;
	}

	/* Interactive state */
	.mana-interactive {
		cursor: pointer;
		transition: transform var(--transition-fast, 0.15s);
		background: none;
		border: none;
		padding: 0;
	}

	.mana-interactive:hover {
		transform: scale(1.15);
	}

	.mana-interactive:active {
		transform: scale(0.95);
	}

	.mana-interactive:focus-visible {
		outline: 2px solid var(--accent-gold, #d4af37);
		outline-offset: 2px;
		border-radius: 50%;
	}
</style>
