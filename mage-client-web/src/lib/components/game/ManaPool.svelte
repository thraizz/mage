<script lang="ts">
	import type { ManaPool } from '$lib/types/game';

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
	function handleManaClick(color: string): void {
		onManaClick(color);
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
	const sizeClass = $derived(() => {
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

<div class="mana-pool {sizeClass()}">
	<div class="mana-pool-header">
		<span class="label">Mana Pool</span>
		{#if hasAnyMana}
			<span class="total">({totalMana})</span>
		{/if}
	</div>

	<div class="mana-orbs">
		{#each manaColors as manaColor}
			{@const count = getManaCount(manaColor.key)}
			{#if count > 0 || showEmpty}
				<button
					class="mana-orb"
					class:has-mana={count > 0}
					class:empty={count === 0}
					style="--orb-color: {manaColor.color}; --orb-text-color: {manaColor.textColor}"
					onclick={() => handleManaClick(manaColor.key)}
					disabled={count === 0}
					title="{count} {manaColor.label} mana"
				>
					<div class="orb-inner">
						<span class="mana-symbol">{manaColor.symbol}</span>
						{#if count > 0}
							<span class="mana-count">{count}</span>
						{/if}
					</div>
				</button>
			{/if}
		{/each}
	</div>

	{#if !hasAnyMana && !showEmpty}
		<div class="empty-state">
			<span>No mana available</span>
		</div>
	{/if}
</div>

<style>
	.mana-pool {
		background: #1a1f2e;
		border: 1px solid #2a3441;
		border-radius: 8px;
		padding: 0.75rem;
		display: flex;
		flex-direction: column;
		gap: 0.75rem;
	}

	.mana-pool-header {
		display: flex;
		align-items: center;
		gap: 0.5rem;
	}

	.label {
		font-size: 0.875rem;
		color: #9ca3af;
		font-weight: 500;
		text-transform: uppercase;
		letter-spacing: 0.5px;
	}

	.total {
		font-size: 0.875rem;
		color: #6b7280;
	}

	.mana-orbs {
		display: flex;
		gap: 0.5rem;
		flex-wrap: wrap;
	}

	.mana-orb {
		position: relative;
		width: 48px;
		height: 48px;
		border-radius: 50%;
		background: var(--orb-color);
		border: 2px solid transparent;
		cursor: pointer;
		transition:
			transform 0.2s,
			border-color 0.2s,
			box-shadow 0.2s;
		display: flex;
		align-items: center;
		justify-content: center;
		box-shadow: 0 2px 4px rgba(0, 0, 0, 0.3);
	}

	.mana-orb.has-mana:hover {
		transform: scale(1.1);
		border-color: rgba(255, 255, 255, 0.5);
		box-shadow: 0 4px 8px rgba(0, 0, 0, 0.4);
	}

	.mana-orb.has-mana:active {
		transform: scale(0.95);
	}

	.mana-orb.empty {
		opacity: 0.3;
		cursor: not-allowed;
	}

	.mana-orb:focus {
		outline: none;
		border-color: rgba(255, 255, 255, 0.7);
	}

	.orb-inner {
		display: flex;
		flex-direction: column;
		align-items: center;
		justify-content: center;
		gap: 0.125rem;
	}

	.mana-symbol {
		font-size: 1rem;
		font-weight: 700;
		color: var(--orb-text-color);
	}

	.mana-count {
		font-size: 0.75rem;
		font-weight: 700;
		color: var(--orb-text-color);
		background: rgba(0, 0, 0, 0.3);
		padding: 0.125rem 0.375rem;
		border-radius: 10px;
		min-width: 1.25rem;
		text-align: center;
	}

	.empty-state {
		text-align: center;
		padding: 1rem;
		color: #6b7280;
		font-size: 0.875rem;
		font-style: italic;
	}

	/* Size Variants */
	.mana-pool-small {
		padding: 0.5rem;
		gap: 0.5rem;
	}

	.mana-pool-small .mana-orb {
		width: 36px;
		height: 36px;
	}

	.mana-pool-small .mana-symbol {
		font-size: 0.875rem;
	}

	.mana-pool-small .mana-count {
		font-size: 0.625rem;
	}

	.mana-pool-large {
		padding: 1rem;
		gap: 1rem;
	}

	.mana-pool-large .mana-orb {
		width: 64px;
		height: 64px;
	}

	.mana-pool-large .mana-symbol {
		font-size: 1.5rem;
	}

	.mana-pool-large .mana-count {
		font-size: 1rem;
	}

	/* Responsive */
	@media (max-width: 768px) {
		.mana-pool {
			padding: 0.5rem;
		}

		.mana-orbs {
			gap: 0.375rem;
		}

		.mana-orb {
			width: 40px;
			height: 40px;
		}

		.mana-symbol {
			font-size: 0.875rem;
		}

		.mana-count {
			font-size: 0.625rem;
		}
	}
</style>
