<script lang="ts">
	import type { ManaPool } from '$lib/types/game';
	import ManaSymbol from '$lib/components/mtg/ManaSymbol.svelte';

	// Props
	let {
		mana = { white: 0, blue: 0, black: 0, red: 0, green: 0, colorless: 0 },
		showEmpty = false,
		size = 'normal',

		onManaClick = (color: string) => {}
	}: {
		mana?: ManaPool;
		showEmpty?: boolean;
		size?: 'small' | 'normal' | 'large';

		onManaClick?: (color: string) => void;
	} = $props();

	// Mana colors configuration
	const manaColors = [
		{ key: 'white', symbol: 'W', color: '#f0f0d8', textColor: '#000', label: 'White' },
		{ key: 'blue', symbol: 'U', color: '#0e68ab', textColor: '#fff', label: 'Blue' },
		{ key: 'black', symbol: 'B', color: '#150b00', textColor: '#fff', label: 'Black' },
		{ key: 'red', symbol: 'R', color: '#d3202a', textColor: '#fff', label: 'Red' },
		{ key: 'green', symbol: 'G', color: '#00733e', textColor: '#fff', label: 'Green' },
		{ key: 'colorless', symbol: 'C', color: '#ccc', textColor: '#000', label: 'Colorless' }
	] as const;

	/**
	 * Handle mana orb click
	 */
	function handleManaClick(_color: string): void {
		onManaClick(_color);
	}

	/**
	 * Get mana count for a color
	 */
	function getManaCount(key: string): number {
		return mana[key as keyof ManaPool] || 0;
	}

	// Derived values
	const totalMana = $derived(
		mana.white + mana.blue + mana.black + mana.red + mana.green + mana.colorless
	);
	const hasAnyMana = $derived(totalMana > 0);
	const sizeClass = $derived.by(() => {
		switch (size) {
			case 'small':
				return 'mana-pool-small';
			case 'large':
				return 'mana-pool-large';
			default:
				return 'mana-pool-normal';
		}
	});
</script>

<div class="mana-pool {sizeClass}" class:has-mana={hasAnyMana}>
	<span class="label">Mana</span>

	{#if hasAnyMana}
		<div class="mana-orbs">
			{#each manaColors as manaColor}
				{@const count = getManaCount(manaColor.key)}
				{#if count > 0 || showEmpty}
					<button
						class="mana-orb"
						class:has-mana={count > 0}
						class:empty={count === 0}
						onclick={() => handleManaClick(manaColor.key)}
						disabled={count === 0}
						title="{count} {manaColor.label} mana"
					>
						<ManaSymbol symbol={manaColor.symbol} size="sm" />
						{#if count > 1}
							<span class="mana-count">{count}</span>
						{/if}
					</button>
				{/if}
			{/each}
		</div>
	{:else}
		<span class="empty-indicator">0</span>
	{/if}
</div>

<style>
	.mana-pool {
		display: flex;
		align-items: center;
		gap: 0.5rem;
		padding: 0.375rem 0.625rem;
		background: rgba(26, 31, 46, 0.6);
		border: 1px solid #2a3441;
		border-radius: 6px;
		min-height: 32px;
	}

	.mana-pool.has-mana {
		background: rgba(26, 31, 46, 0.9);
		border-color: rgba(251, 191, 36, 0.3);
	}

	.label {
		font-size: 0.6875rem;
		color: #6b7280;
		font-weight: 600;
		text-transform: uppercase;
		letter-spacing: 0.5px;
		white-space: nowrap;
	}

	.mana-orbs {
		display: flex;
		gap: 0.25rem;
		flex-wrap: wrap;
		align-items: center;
	}

	.mana-orb {
		position: relative;
		display: flex;
		align-items: center;
		justify-content: center;
		padding: 0.125rem;
		border-radius: 4px;
		background: transparent;
		border: none;
		cursor: pointer;
		transition: all 0.15s;
	}

	.mana-orb.has-mana:hover {
		background: rgba(255, 255, 255, 0.1);
		transform: scale(1.1);
	}

	.mana-orb.has-mana:active {
		transform: scale(0.95);
	}

	.mana-orb.empty {
		opacity: 0.25;
		cursor: not-allowed;
	}

	.mana-orb:focus {
		outline: none;
	}

	.mana-count {
		position: absolute;
		bottom: -2px;
		right: -4px;
		font-size: 0.5625rem;
		font-weight: 700;
		color: #fff;
		background: rgba(0, 0, 0, 0.75);
		padding: 0 0.25rem;
		border-radius: 4px;
		min-width: 0.875rem;
		text-align: center;
		line-height: 1.2;
	}

	.empty-indicator {
		font-size: 0.75rem;
		color: #4b5563;
		font-weight: 500;
	}

	/* Size Variants */
	.mana-pool-small {
		padding: 0.25rem 0.5rem;
		gap: 0.375rem;
		min-height: 28px;
	}

	.mana-pool-small .label {
		font-size: 0.625rem;
	}

	.mana-pool-large {
		padding: 0.5rem 0.75rem;
		gap: 0.625rem;
		min-height: 40px;
	}

	.mana-pool-large .mana-orb {
		padding: 0.25rem;
	}

	.mana-pool-large .mana-count {
		font-size: 0.6875rem;
	}
</style>
