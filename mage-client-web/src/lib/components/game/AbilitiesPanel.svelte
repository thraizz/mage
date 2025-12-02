<script lang="ts">
	import { onMount, onDestroy } from 'svelte';
	import type { CardAction } from '$lib/generated/mage/v1/models';
	import AbilityItem from './AbilityItem.svelte';

	// Props
	let {
		cardId,
		cardName,
		abilities,
		onActivate = () => {},
		onClose = () => {}
	}: {
		cardId: string;
		cardName: string;
		abilities: CardAction[];
		onActivate?: (abilityId: string) => void;
		onClose?: () => void;
	} = $props();

	// State
	let selectedIndex = $state(0);
	let panelRef: HTMLDivElement | undefined = $state();

	// Derived
	const enabledAbilities = $derived(abilities.filter((a) => a.isEnabled));
	const hasEnabledAbilities = $derived(enabledAbilities.length > 0);

	/**
	 * Handle keyboard navigation
	 */
	function handleKeydown(event: KeyboardEvent) {
		switch (event.key) {
			case 'Escape':
				event.preventDefault();
				onClose();
				break;
			case 'ArrowUp':
				event.preventDefault();
				selectedIndex = Math.max(0, selectedIndex - 1);
				break;
			case 'ArrowDown':
				event.preventDefault();
				selectedIndex = Math.min(abilities.length - 1, selectedIndex + 1);
				break;
			case 'Enter':
			case ' ':
				event.preventDefault();
				if (abilities[selectedIndex]?.isEnabled) {
					onActivate(abilities[selectedIndex].actionId);
				}
				break;
		}
	}

	/**
	 * Handle click outside to close
	 */
	function handleClickOutside(event: MouseEvent) {
		if (panelRef && !panelRef.contains(event.target as Node)) {
			onClose();
		}
	}

	/**
	 * Handle ability activation
	 */
	function handleAbilityActivate(abilityId: string) {
		onActivate(abilityId);
	}

	// Setup event listeners
	onMount(() => {
		console.log('[AbilitiesPanel] Mounted with:', { cardId, cardName, abilities: abilities.length });
		// Small delay to prevent immediate close from the same click that opened the panel
		setTimeout(() => {
			document.addEventListener('click', handleClickOutside);
		}, 100);
		document.addEventListener('keydown', handleKeydown);
	});

	onDestroy(() => {
		document.removeEventListener('click', handleClickOutside);
		document.removeEventListener('keydown', handleKeydown);
	});
</script>

<div class="abilities-panel-overlay">
	<div
		class="abilities-panel"
		bind:this={panelRef}
		role="dialog"
		aria-labelledby="abilities-title"
		aria-modal="true"
	>
		<div class="panel-header">
			<div class="header-icon">⚡</div>
			<div class="header-content">
				<h3 id="abilities-title">Activated Abilities</h3>
				<p class="card-name">{cardName}</p>
			</div>
			<button class="close-btn" onclick={onClose} aria-label="Close abilities panel">
				×
			</button>
		</div>

		<div class="panel-body">
			{#if abilities.length === 0}
				<div class="empty-state">
					<p>This permanent has no activated abilities.</p>
				</div>
			{:else}
				<div class="abilities-list" role="listbox" aria-label="Available abilities">
					{#each abilities as ability, index}
						<AbilityItem
							{ability}
							isSelected={index === selectedIndex}
							onActivate={handleAbilityActivate}
						/>
					{/each}
				</div>
			{/if}
		</div>

		<div class="panel-footer">
			<div class="keyboard-hints">
				<span class="hint"><kbd>↑</kbd><kbd>↓</kbd> navigate</span>
				<span class="hint"><kbd>Enter</kbd> activate</span>
				<span class="hint"><kbd>Esc</kbd> close</span>
			</div>
			{#if !hasEnabledAbilities && abilities.length > 0}
				<p class="all-disabled-hint">All abilities are currently unavailable</p>
			{/if}
		</div>
	</div>
</div>

<style>
	.abilities-panel-overlay {
		position: fixed;
		inset: 0;
		background: rgba(0, 0, 0, 0.6);
		display: flex;
		align-items: center;
		justify-content: center;
		z-index: 100;
		animation: fadeIn 0.15s ease-out;
	}

	@keyframes fadeIn {
		from {
			opacity: 0;
		}
		to {
			opacity: 1;
		}
	}

	.abilities-panel {
		background: #141821;
		border: 2px solid #667eea;
		border-radius: 12px;
		width: 90%;
		max-width: 420px;
		max-height: 70vh;
		display: flex;
		flex-direction: column;
		box-shadow:
			0 25px 50px -12px rgba(0, 0, 0, 0.6),
			0 0 40px rgba(102, 126, 234, 0.2);
		animation: slideUp 0.2s ease-out;
	}

	@keyframes slideUp {
		from {
			transform: translateY(20px);
			opacity: 0;
		}
		to {
			transform: translateY(0);
			opacity: 1;
		}
	}

	.panel-header {
		display: flex;
		align-items: center;
		gap: 0.75rem;
		padding: 1rem 1.25rem;
		background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
		border-radius: 10px 10px 0 0;
	}

	.header-icon {
		font-size: 1.5rem;
	}

	.header-content {
		flex: 1;
	}

	.header-content h3 {
		margin: 0;
		font-size: 1rem;
		font-weight: 600;
		color: white;
		text-transform: uppercase;
		letter-spacing: 0.5px;
	}

	.card-name {
		margin: 0.25rem 0 0 0;
		font-size: 0.875rem;
		color: rgba(255, 255, 255, 0.8);
	}

	.close-btn {
		width: 32px;
		height: 32px;
		display: flex;
		align-items: center;
		justify-content: center;
		background: rgba(255, 255, 255, 0.15);
		border: none;
		border-radius: 6px;
		color: white;
		font-size: 1.25rem;
		cursor: pointer;
		transition: background 0.15s;
	}

	.close-btn:hover {
		background: rgba(255, 255, 255, 0.25);
	}

	.panel-body {
		flex: 1;
		overflow-y: auto;
		padding: 1rem;
	}

	.empty-state {
		text-align: center;
		padding: 2rem 1rem;
		color: #6b7280;
	}

	.abilities-list {
		display: flex;
		flex-direction: column;
		gap: 0.625rem;
	}

	.panel-footer {
		padding: 0.75rem 1rem;
		background: #1a1f2e;
		border-top: 1px solid #2a3441;
		border-radius: 0 0 10px 10px;
	}

	.keyboard-hints {
		display: flex;
		justify-content: center;
		gap: 1rem;
		flex-wrap: wrap;
	}

	.hint {
		display: flex;
		align-items: center;
		gap: 0.25rem;
		font-size: 0.6875rem;
		color: #6b7280;
	}

	kbd {
		display: inline-block;
		padding: 0.125rem 0.375rem;
		background: #2a3441;
		border: 1px solid #374151;
		border-radius: 3px;
		font-family: monospace;
		font-size: 0.625rem;
		color: #9ca3af;
	}

	.all-disabled-hint {
		margin: 0.5rem 0 0 0;
		font-size: 0.75rem;
		color: #f87171;
		text-align: center;
		font-style: italic;
	}

	/* Responsive */
	@media (max-width: 480px) {
		.abilities-panel {
			width: 95%;
			max-height: 80vh;
		}

		.panel-header {
			padding: 0.875rem 1rem;
		}

		.header-icon {
			font-size: 1.25rem;
		}

		.keyboard-hints {
			gap: 0.75rem;
		}
	}
</style>

