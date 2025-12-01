<script lang="ts">
	/**
	 * Component that renders text with inline mana symbols
	 * Converts {W}, {U}, {B}, {R}, {G}, {C}, {T}, {Q}, {1}, {2}, etc. to mana font icons
	 */

	interface Props {
		text: string;
		size?: 'sm' | 'md' | 'lg';
	}

	let { text, size = 'sm' }: Props = $props();

	type ParsedPart = { type: 'text'; content: string } | { type: 'mana'; symbol: string };

	// Parse text into parts (text and mana symbols)
	function parseText(input: string): ParsedPart[] {
		const parts: ParsedPart[] = [];
		// Match mana symbols like {W}, {U}, {B}, {R}, {G}, {C}, {T}, {Q}, {S}, {E}, {X}, {0}-{20}
		const regex = /\{([WUBRGCTSQEX]|\d+)\}/gi;
		let lastIndex = 0;
		let match;

		while ((match = regex.exec(input)) !== null) {
			// Add text before the match
			if (match.index > lastIndex) {
				parts.push({ type: 'text', content: input.slice(lastIndex, match.index) });
			}
			// Add the mana symbol
			parts.push({ type: 'mana', symbol: match[1].toUpperCase() });
			lastIndex = regex.lastIndex;
		}

		// Add remaining text
		if (lastIndex < input.length) {
			parts.push({ type: 'text', content: input.slice(lastIndex) });
		}

		return parts;
	}

	// Convert symbol to mana-font class
	function getManaClass(symbol: string): string {
		switch (symbol) {
			case 'W': return 'ms-w';
			case 'U': return 'ms-u';
			case 'B': return 'ms-b';
			case 'R': return 'ms-r';
			case 'G': return 'ms-g';
			case 'C': return 'ms-c';
			case 'X': return 'ms-x';
			case 'T': return 'ms-tap';
			case 'Q': return 'ms-untap';
			case 'S': return 'ms-s';
			case 'E': return 'ms-e';
			default:
				// Numeric values
				const num = parseInt(symbol, 10);
				if (!isNaN(num)) {
					return `ms-${num}`;
				}
				return 'ms-c';
		}
	}

	const parsedParts = $derived(parseText(text));
	
	const sizeClass = $derived({
		sm: 'mana-text-sm',
		md: 'mana-text-md',
		lg: 'mana-text-lg'
	}[size] || 'mana-text-sm');
</script>

<span class="mana-text {sizeClass}">
	{#each parsedParts as part}
		{#if part.type === 'text'}
			{part.content}
		{:else}
			<i class="ms ms-cost {getManaClass(part.symbol)}"></i>
		{/if}
	{/each}
</span>

<style>
	.mana-text {
		display: inline;
		vertical-align: baseline;
	}

	.mana-text :global(.ms) {
		vertical-align: middle;
		margin: 0 1px;
	}

	.mana-text-sm :global(.ms) {
		font-size: 0.9em;
	}

	.mana-text-md :global(.ms) {
		font-size: 1.1em;
	}

	.mana-text-lg :global(.ms) {
		font-size: 1.3em;
	}
</style>


