<!--
  DeclareAttackers.svelte
  
  Combat component for declaring attackers during the Declare Attackers step.
  Shows available creatures that can attack and allows player to select them.
  
  Features:
  - Toggle attackers by clicking creatures
  - Multi-defender support (for planeswalkers)
  - Visual indicators for attacking creatures
  - Keyboard shortcuts: Enter to confirm, Escape to skip
-->
<script lang="ts">
	import { onMount, onDestroy } from 'svelte';
	import { declareAttackers, skipCombat } from '$lib/api/game';
	import { combatStore, declaredAttackerIds, declaredAttackerCount } from '$lib/stores/combat';
	import type { DefenderTarget, ParsedCombatOptions } from '$lib/types/combat';
	import type { CardView } from '$lib/generated/mage/v1/models';

	// Props
	let {
		gameId,
		options,
		battlefieldCards,
		defenders,
		onComplete = () => {}
	}: {
		gameId: string;
		options: ParsedCombatOptions;
		battlefieldCards: CardView[];
		defenders: DefenderTarget[];
		onComplete?: () => void;
	} = $props();

	// State
	let isSubmitting = $state(false);
	let error = $state<string | null>(null);

	// Derived state from combat store
	const declaredIds = $derived($declaredAttackerIds);
	const attackerCount = $derived($declaredAttackerCount);

	// Build a map of card names for display
	const cardNames = $derived(() => {
		const map = new Map<string, string>();
		for (const card of battlefieldCards) {
			map.set(card.id, card.name);
		}
		return map;
	});

	// Initialize combat store with parsed options
	$effect(() => {
		combatStore.enterDeclareAttackersPhase(options, cardNames(), defenders);
	});

	/**
	 * Toggle an attacker when clicking a card
	 */
	function handleCardClick(cardId: string) {
		if (isSubmitting) return;
		
		// Check if this card can attack (is in available attackers)
		const canAttackDefenders = options.attackOptions
			.filter((opt) => opt.cardId === cardId)
			.map((opt) => opt.defenderId);
		
		if (canAttackDefenders.length === 0) {
			// Card cannot attack
			return;
		}

		// Toggle the attacker
		combatStore.toggleAttacker(cardId, canAttackDefenders[0]);
		error = null;
	}

	/**
	 * Change defender for an attacker (when multiple defenders available)
	 */
	function handleDefenderChange(cardId: string, defenderId: string) {
		if (isSubmitting) return;
		combatStore.changeAttackTarget(cardId, defenderId);
	}

	/**
	 * Submit all declared attackers
	 */
	async function handleConfirm() {
		if (isSubmitting) return;

		isSubmitting = true;
		error = null;

		try {
			const attackers = combatStore.getDeclaredAttackers();
			await declareAttackers(gameId, attackers);
			combatStore.reset();
			onComplete();
		} catch (err) {
			error = err instanceof Error ? err.message : 'Failed to declare attackers';
			console.error('Failed to declare attackers:', err);
		} finally {
			isSubmitting = false;
		}
	}

	/**
	 * Skip combat (declare no attackers)
	 */
	async function handleSkip() {
		if (isSubmitting) return;

		isSubmitting = true;
		error = null;

		try {
			await skipCombat(gameId);
			combatStore.reset();
			onComplete();
		} catch (err) {
			error = err instanceof Error ? err.message : 'Failed to skip combat';
			console.error('Failed to skip combat:', err);
		} finally {
			isSubmitting = false;
		}
	}

	/**
	 * Handle keyboard shortcuts
	 */
	function handleKeydown(event: KeyboardEvent) {
		if (isSubmitting) return;

		switch (event.key) {
			case 'Enter':
				event.preventDefault();
				handleConfirm();
				break;
			case 'Escape':
				event.preventDefault();
				handleSkip();
				break;
		}
	}

	// Get valid defenders for a specific attacker
	function getValidDefendersForAttacker(cardId: string): string[] {
		return options.attackOptions
			.filter((opt) => opt.cardId === cardId)
			.map((opt) => opt.defenderId);
	}

	// Get defender name by ID
	function getDefenderName(defenderId: string): string {
		const defender = defenders.find((d) => d.id === defenderId);
		return defender?.name || 'Unknown';
	}

	// Get available attackers (cards that can attack)
	const availableAttackerIds = $derived(() => {
		const ids = new Set<string>();
		for (const opt of options.attackOptions) {
			ids.add(opt.cardId);
		}
		return ids;
	});

	// Add global keyboard listener
	onMount(() => {
		window.addEventListener('keydown', handleKeydown);
	});

	onDestroy(() => {
		window.removeEventListener('keydown', handleKeydown);
		// Don't reset combat store on destroy - let the parent handle it
	});
</script>

<!-- Overlay -->
<div class="declare-attackers-overlay" role="dialog" aria-modal="true" aria-labelledby="attackers-title">
	<!-- Top Banner -->
	<div class="attackers-banner">
		<div class="banner-icon">⚔️</div>
		<div class="banner-content">
			<h3 id="attackers-title" class="banner-title">Declare Attackers</h3>
			<p class="banner-description">
				Click creatures on the battlefield to declare them as attackers
			</p>
			<div class="banner-status">
				{#if attackerCount > 0}
					<span class="attacker-count">{attackerCount} creature{attackerCount !== 1 ? 's' : ''} attacking</span>
				{:else}
					<span class="attacker-hint">No attackers declared yet</span>
				{/if}
			</div>
		</div>
	</div>

	<!-- Attacker List (shows declared attackers with defender selection) -->
	{#if attackerCount > 0}
		<div class="attackers-list">
			<h4 class="list-title">Declared Attackers</h4>
			<div class="attacker-items">
				{#each combatStore.getDeclaredAttackers() as attacker}
					{@const validDefenders = getValidDefendersForAttacker(attacker.cardId)}
					<div class="attacker-item">
						<span class="attacker-name">{cardNames().get(attacker.cardId) || 'Unknown'}</span>
						<span class="attack-arrow">→</span>
						{#if validDefenders.length > 1}
							<select
								class="defender-select"
								value={attacker.defenderId}
								onchange={(e) => handleDefenderChange(attacker.cardId, e.currentTarget.value)}
								disabled={isSubmitting}
							>
								{#each validDefenders as defenderId}
									<option value={defenderId}>{getDefenderName(defenderId)}</option>
								{/each}
							</select>
						{:else}
							<span class="defender-name">{getDefenderName(attacker.defenderId)}</span>
						{/if}
						<button
							class="remove-attacker"
							onclick={() => handleCardClick(attacker.cardId)}
							disabled={isSubmitting}
							title="Remove attacker"
						>
							✕
						</button>
					</div>
				{/each}
			</div>
		</div>
	{/if}

	<!-- Error Message -->
	{#if error}
		<div class="error-message" role="alert">
			{error}
		</div>
	{/if}

	<!-- Bottom Action Bar -->
	<div class="attackers-actions">
		<div class="action-hints">
			<span class="hint">
				<kbd>ESC</kbd> to skip combat
			</span>
			<span class="hint">
				<kbd>Enter</kbd> to confirm
			</span>
		</div>
		<div class="action-buttons">
			<button
				class="btn-skip"
				onclick={handleSkip}
				disabled={isSubmitting}
				type="button"
			>
				{#if isSubmitting}
					<span class="spinner"></span>
				{/if}
				Skip Combat
			</button>
			<button
				class="btn-confirm"
				onclick={handleConfirm}
				disabled={isSubmitting}
				type="button"
			>
				{#if isSubmitting}
					<span class="spinner"></span>
				{/if}
				{#if attackerCount > 0}
					Attack with {attackerCount}
				{:else}
					Confirm (No Attackers)
				{/if}
			</button>
		</div>
	</div>
</div>

<style>
	/* Overlay */
	.declare-attackers-overlay {
		position: fixed;
		inset: 0;
		background: rgba(0, 0, 0, 0.6);
		z-index: 100;
		display: flex;
		flex-direction: column;
		justify-content: space-between;
		pointer-events: auto;
	}

	/* Top Banner */
	.attackers-banner {
		display: flex;
		align-items: center;
		gap: 1rem;
		padding: 1rem 1.5rem;
		background: linear-gradient(180deg, rgba(239, 68, 68, 0.2) 0%, transparent 100%);
		border-bottom: 2px solid rgba(239, 68, 68, 0.5);
		animation: banner-slide-in 0.2s ease-out;
	}

	@keyframes banner-slide-in {
		from {
			opacity: 0;
			transform: translateY(-20px);
		}
		to {
			opacity: 1;
			transform: translateY(0);
		}
	}

	.banner-icon {
		font-size: 2.5rem;
		animation: sword-bounce 1s ease-in-out infinite;
	}

	@keyframes sword-bounce {
		0%, 100% {
			transform: translateY(0) rotate(-5deg);
		}
		50% {
			transform: translateY(-5px) rotate(5deg);
		}
	}

	.banner-content {
		flex: 1;
	}

	.banner-title {
		margin: 0 0 0.25rem 0;
		font-size: 1.25rem;
		font-weight: 700;
		color: #ef4444;
		text-transform: uppercase;
		letter-spacing: 1px;
	}

	.banner-description {
		margin: 0 0 0.5rem 0;
		font-size: 0.9375rem;
		color: #ffffff;
	}

	.banner-status {
		font-size: 0.875rem;
	}

	.attacker-count {
		color: #ef4444;
		font-weight: 600;
	}

	.attacker-hint {
		color: #9ca3af;
		font-style: italic;
	}

	/* Attackers List */
	.attackers-list {
		padding: 1rem 1.5rem;
		background: rgba(0, 0, 0, 0.4);
		max-height: 200px;
		overflow-y: auto;
	}

	.list-title {
		margin: 0 0 0.75rem 0;
		font-size: 0.875rem;
		font-weight: 600;
		color: #ef4444;
		text-transform: uppercase;
		letter-spacing: 0.5px;
	}

	.attacker-items {
		display: flex;
		flex-direction: column;
		gap: 0.5rem;
	}

	.attacker-item {
		display: flex;
		align-items: center;
		gap: 0.75rem;
		padding: 0.5rem 0.75rem;
		background: rgba(239, 68, 68, 0.15);
		border: 1px solid rgba(239, 68, 68, 0.3);
		border-radius: 6px;
	}

	.attacker-name {
		font-weight: 600;
		color: #ffffff;
		flex: 1;
	}

	.attack-arrow {
		color: #ef4444;
		font-size: 1.25rem;
	}

	.defender-name {
		color: #f87171;
		font-weight: 500;
	}

	.defender-select {
		padding: 0.25rem 0.5rem;
		background: #1f2937;
		border: 1px solid #374151;
		border-radius: 4px;
		color: #f87171;
		font-size: 0.875rem;
		cursor: pointer;
	}

	.defender-select:hover {
		border-color: #ef4444;
	}

	.remove-attacker {
		padding: 0.25rem 0.5rem;
		background: transparent;
		border: 1px solid rgba(239, 68, 68, 0.5);
		border-radius: 4px;
		color: #f87171;
		font-size: 0.75rem;
		cursor: pointer;
		transition: all 0.2s;
	}

	.remove-attacker:hover:not(:disabled) {
		background: rgba(239, 68, 68, 0.2);
		border-color: #ef4444;
	}

	.remove-attacker:disabled {
		opacity: 0.5;
		cursor: not-allowed;
	}

	/* Error Message */
	.error-message {
		margin: 0.5rem 1.5rem;
		padding: 0.75rem 1rem;
		background: rgba(239, 68, 68, 0.2);
		border: 1px solid #ef4444;
		border-radius: 6px;
		color: #fca5a5;
		font-size: 0.875rem;
	}

	/* Bottom Action Bar */
	.attackers-actions {
		display: flex;
		justify-content: space-between;
		align-items: center;
		padding: 1rem 1.5rem;
		background: linear-gradient(0deg, rgba(0, 0, 0, 0.8) 0%, transparent 100%);
	}

	.action-hints {
		display: flex;
		gap: 1.5rem;
	}

	.hint {
		color: #6b7280;
		font-size: 0.75rem;
	}

	kbd {
		display: inline-block;
		padding: 0.125rem 0.375rem;
		background: #374151;
		border: 1px solid #4b5563;
		border-radius: 4px;
		font-family: monospace;
		font-size: 0.6875rem;
		color: #e5e7eb;
	}

	.action-buttons {
		display: flex;
		gap: 0.75rem;
	}

	.btn-skip,
	.btn-confirm {
		display: flex;
		align-items: center;
		gap: 0.5rem;
		padding: 0.625rem 1.25rem;
		border-radius: 6px;
		font-weight: 600;
		font-size: 0.875rem;
		cursor: pointer;
		transition: all 0.2s;
		border: none;
	}

	.btn-skip {
		background: #374151;
		color: #e5e7eb;
	}

	.btn-skip:hover:not(:disabled) {
		background: #4b5563;
	}

	.btn-confirm {
		background: linear-gradient(135deg, #ef4444 0%, #dc2626 100%);
		color: white;
	}

	.btn-confirm:hover:not(:disabled) {
		background: linear-gradient(135deg, #f87171 0%, #ef4444 100%);
		transform: translateY(-1px);
		box-shadow: 0 4px 12px rgba(239, 68, 68, 0.4);
	}

	.btn-skip:disabled,
	.btn-confirm:disabled {
		opacity: 0.6;
		cursor: not-allowed;
		transform: none;
	}

	.spinner {
		width: 14px;
		height: 14px;
		border: 2px solid rgba(255, 255, 255, 0.3);
		border-top-color: white;
		border-radius: 50%;
		animation: spin 0.6s linear infinite;
	}

	@keyframes spin {
		to {
			transform: rotate(360deg);
		}
	}

	/* Responsive */
	@media (max-width: 600px) {
		.attackers-banner {
			padding: 0.75rem 1rem;
		}

		.banner-icon {
			font-size: 2rem;
		}

		.banner-title {
			font-size: 1rem;
		}

		.banner-description {
			font-size: 0.875rem;
		}

		.attackers-actions {
			flex-direction: column;
			gap: 0.75rem;
		}

		.action-hints {
			flex-wrap: wrap;
			justify-content: center;
		}

		.action-buttons {
			width: 100%;
			justify-content: center;
		}
	}
</style>

