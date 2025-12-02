<script lang="ts">
	import type { CardAction } from '$lib/generated/mage/v1/models';

	// Props
	let {
		ability,
		isSelected = false,
		onActivate = () => {}
	}: {
		ability: CardAction;
		isSelected?: boolean;
		onActivate?: (abilityId: string) => void;
	} = $props();

	// Parse cost from display text (e.g., "{2}, {T}: Draw a card" -> cost: "{2}, {T}", effect: "Draw a card")
	const parsed = $derived.by(() => {
		const text = ability.displayText || 'Activate ability';
		const colonIndex = text.indexOf(':');
		if (colonIndex > -1) {
			return {
				cost: text.substring(0, colonIndex).trim(),
				effect: text.substring(colonIndex + 1).trim()
			};
		}
		return {
			cost: '',
			effect: text
		};
	});

	/**
	 * Handle click on ability
	 */
	function handleClick() {
		if (ability.isEnabled) {
			onActivate(ability.actionId);
		}
	}

	/**
	 * Handle keyboard activation
	 */
	function handleKeydown(event: KeyboardEvent) {
		if ((event.key === 'Enter' || event.key === ' ') && ability.isEnabled) {
			event.preventDefault();
			onActivate(ability.actionId);
		}
	}
</script>

<button
	class="ability-item"
	class:enabled={ability.isEnabled}
	class:disabled={!ability.isEnabled}
	class:selected={isSelected}
	onclick={handleClick}
	onkeydown={handleKeydown}
	disabled={!ability.isEnabled}
	title={!ability.isEnabled ? ability.disabledReason : 'Click to activate'}
	aria-disabled={!ability.isEnabled}
>
	<div class="ability-content">
		{#if parsed.cost}
			<span class="ability-cost">{parsed.cost}:</span>
		{/if}
		<span class="ability-effect">{parsed.effect}</span>
	</div>

	{#if !ability.isEnabled && ability.disabledReason}
		<div class="disabled-reason">
			<span class="reason-icon">⚠</span>
			<span class="reason-text">{ability.disabledReason}</span>
		</div>
	{/if}
</button>

<style>
	.ability-item {
		display: flex;
		flex-direction: column;
		gap: 0.375rem;
		width: 100%;
		padding: 0.875rem 1rem;
		background: #1a1f2e;
		border: 2px solid #2a3441;
		border-radius: 8px;
		text-align: left;
		cursor: pointer;
		transition: all 0.15s ease;
	}

	.ability-item.enabled:hover {
		border-color: #667eea;
		background: #1e2438;
		transform: translateX(4px);
	}

	.ability-item.enabled:focus-visible {
		outline: none;
		border-color: #667eea;
		box-shadow: 0 0 0 3px rgba(102, 126, 234, 0.3);
	}

	.ability-item.selected {
		border-color: #667eea;
		background: rgba(102, 126, 234, 0.15);
	}

	.ability-item.disabled {
		opacity: 0.6;
		cursor: not-allowed;
	}

	.ability-content {
		display: flex;
		flex-wrap: wrap;
		gap: 0.375rem;
		font-size: 0.9375rem;
		line-height: 1.5;
	}

	.ability-cost {
		color: #fbbf24;
		font-weight: 600;
		font-family: 'Courier New', monospace;
	}

	.ability-effect {
		color: #e2e8f0;
	}

	.disabled-reason {
		display: flex;
		align-items: center;
		gap: 0.375rem;
		font-size: 0.75rem;
		color: #f87171;
	}

	.reason-icon {
		font-size: 0.75rem;
	}

	.reason-text {
		font-style: italic;
	}
</style>

