<script lang="ts">
	import ManaSymbol from './ManaSymbol.svelte';

	type ManaColor = 'W' | 'U' | 'B' | 'R' | 'G' | 'C' | 'X';
	type ManaSymbolType = ManaColor | number;

	interface Props {
		cost: string; // e.g., "2WU", "3BB", "WUBRG"
		size?: 'sm' | 'md' | 'lg';
	}

	let { cost, size = 'md' }: Props = $props();

	// Parse mana cost string into array of symbols
	function parseCost(costString: string): ManaSymbolType[] {
		const symbols: ManaSymbolType[] = [];
		let i = 0;

		while (i < costString.length) {
			const char = costString[i];

			// Check for numeric (could be multi-digit like "10", "15")
			if (/\d/.test(char)) {
				let num = '';
				while (i < costString.length && /\d/.test(costString[i])) {
					num += costString[i];
					i++;
				}
				symbols.push(parseInt(num, 10));
			}
			// Single letter mana symbol
			else if (/[WUBRGCX]/i.test(char)) {
				symbols.push(char.toUpperCase() as ManaColor);
				i++;
			}
			// Handle hybrid mana like {W/U} or {2/W}
			else if (char === '{') {
				const end = costString.indexOf('}', i);
				if (end !== -1) {
					const inner = costString.slice(i + 1, end);
					// For now, just take the first part of hybrid
					const parts = inner.split('/');
					if (parts.length > 0) {
						const first = parts[0];
						if (/\d+/.test(first)) {
							symbols.push(parseInt(first, 10));
						} else {
							symbols.push(first.toUpperCase() as ManaColor);
						}
					}
					i = end + 1;
				} else {
					i++;
				}
			} else {
				i++;
			}
		}

		return symbols;
	}

	const symbols = $derived(parseCost(cost));
</script>

<span class="mana-cost">
	{#each symbols as symbol}
		<ManaSymbol {symbol} {size} />
	{/each}
</span>

<style>
	.mana-cost {
		display: inline-flex;
		align-items: center;
		gap: 2px;
	}
</style>
