<!--
  AssignDamage.svelte
  
  Combat component for assigning combat damage when manual assignment is required.
  Used when an attacker has trample or is blocked by multiple creatures.
  
  Features:
  - Show attacker with power value
  - List blockers in damage order (first must receive lethal before second)
  - Number inputs with +/- buttons for each blocker
  - Trample row for excess damage to defending player
  - Validation: total must equal attacker power, lethal damage ordering enforced
  - Keyboard shortcuts: Enter to confirm
-->
<script lang="ts">
	import { onMount, onDestroy } from 'svelte';
	import { assignCombatDamage } from '$lib/api/game';
	import { combatStore } from '$lib/stores/combat';
	import type { DamageAssignmentPrompt, DamageAssignment, OrderedBlocker } from '$lib/types/combat';

	// Props
	let {
		gameId,
		prompt,
		onComplete = () => {}
	}: {
		gameId: string;
		prompt: DamageAssignmentPrompt;
		onComplete?: () => void;
	} = $props();

	// State - damage assignments per target
	let assignments = $state<Map<string, number>>(new Map());
	let isSubmitting = $state(false);
	let validationError = $state<string | null>(null);

	// Initialize assignments with lethal damage to first blocker(s)
	$effect(() => {
		const initial = new Map<string, number>();
		let remainingPower = prompt.attackerPower;

		// Assign lethal to each blocker in order
		for (const blocker of prompt.blockers.sort((a, b) => a.order - b.order)) {
			const lethalNeeded = Math.max(0, blocker.toughness - blocker.damage);
			const toAssign = Math.min(remainingPower, lethalNeeded);
			initial.set(blocker.id, toAssign);
			remainingPower -= toAssign;
		}

		// If trample and remaining damage, assign to defending player
		if (prompt.hasTrample && remainingPower > 0) {
			initial.set(prompt.defendingPlayerId, remainingPower);
		}

		assignments = initial;
	});

	// Computed values
	const totalAssigned = $derived(() => {
		let total = 0;
		for (const damage of assignments.values()) {
			total += damage;
		}
		return total;
	});

	const remainingDamage = $derived(() => prompt.attackerPower - totalAssigned());

	// Sorted blockers by damage order
	const sortedBlockers = $derived(() => [...prompt.blockers].sort((a, b) => a.order - b.order));

	/**
	 * Get damage assigned to a target
	 */
	function getDamage(targetId: string): number {
		return assignments.get(targetId) || 0;
	}

	/**
	 * Set damage for a target
	 */
	function setDamage(targetId: string, damage: number) {
		const newAssignments = new Map(assignments);
		newAssignments.set(targetId, Math.max(0, damage));
		assignments = newAssignments;
		validationError = null;
	}

	/**
	 * Adjust damage by delta (+1 or -1)
	 */
	function adjustDamage(targetId: string, delta: number) {
		const current = getDamage(targetId);
		setDamage(targetId, current + delta);
	}

	/**
	 * Calculate lethal damage needed for a blocker
	 */
	function getLethalDamage(blocker: OrderedBlocker): number {
		return Math.max(0, blocker.toughness - blocker.damage);
	}

	/**
	 * Check if a blocker has received lethal damage
	 */
	function hasLethalDamage(blocker: OrderedBlocker): boolean {
		return getDamage(blocker.id) >= getLethalDamage(blocker);
	}

	/**
	 * Validate damage assignment per MTG rules
	 */
	function validateAssignment(): { valid: boolean; error?: string } {
		// Total must equal attacker power
		if (totalAssigned() !== prompt.attackerPower) {
			return {
				valid: false,
				error: `Must assign exactly ${prompt.attackerPower} damage (currently ${totalAssigned()})`
			};
		}

		// Check lethal damage ordering
		let mustAssignLethal = true;
		for (const blocker of sortedBlockers()) {
			const assigned = getDamage(blocker.id);
			const lethal = getLethalDamage(blocker);

			if (mustAssignLethal && assigned < lethal) {
				// Check if later targets have damage
				const laterTargets = [...sortedBlockers().filter((b) => b.order > blocker.order)];
				if (prompt.hasTrample) {
					// Also check defending player
					const playerDamage = getDamage(prompt.defendingPlayerId);
					if (playerDamage > 0) {
						return {
							valid: false,
							error: `Must assign lethal damage (${lethal}) to ${blocker.name} before assigning trample damage`
						};
					}
				}

				for (const later of laterTargets) {
					if (getDamage(later.id) > 0) {
						return {
							valid: false,
							error: `Must assign lethal damage (${lethal}) to ${blocker.name} before assigning to ${later.name}`
						};
					}
				}
			}

			if (assigned >= lethal) {
				// This blocker has lethal, can proceed to next
			} else {
				// Not lethal, can't assign to later targets
				mustAssignLethal = false;
			}
		}

		return { valid: true };
	}

	/**
	 * Submit damage assignment
	 */
	async function handleSubmit() {
		if (isSubmitting) return;

		const validation = validateAssignment();
		if (!validation.valid) {
			validationError = validation.error || 'Invalid damage assignment';
			return;
		}

		isSubmitting = true;
		validationError = null;

		try {
			const damageList: Array<{ targetId: string; damage: number }> = [];

			// Add blocker damage
			for (const blocker of prompt.blockers) {
				const damage = getDamage(blocker.id);
				if (damage > 0) {
					damageList.push({ targetId: blocker.id, damage });
				}
			}

			// Add trample damage if any
			if (prompt.hasTrample) {
				const trampleDamage = getDamage(prompt.defendingPlayerId);
				if (trampleDamage > 0) {
					damageList.push({ targetId: prompt.defendingPlayerId, damage: trampleDamage });
				}
			}

			await assignCombatDamage(gameId, damageList);
			combatStore.reset();
			onComplete();
		} catch (err) {
			validationError = err instanceof Error ? err.message : 'Failed to assign damage';
			console.error('Failed to assign damage:', err);
		} finally {
			isSubmitting = false;
		}
	}

	/**
	 * Handle keyboard shortcuts
	 */
	function handleKeydown(event: KeyboardEvent) {
		if (isSubmitting) return;

		if (event.key === 'Enter') {
			event.preventDefault();
			handleSubmit();
		}
	}

	onMount(() => {
		window.addEventListener('keydown', handleKeydown);
	});

	onDestroy(() => {
		window.removeEventListener('keydown', handleKeydown);
	});
</script>

<!-- Overlay -->
<div class="assign-damage-overlay" role="dialog" aria-modal="true" aria-labelledby="damage-title">
	<div class="damage-panel">
		<!-- Header -->
		<div class="damage-header">
			<div class="header-icon">⚔️</div>
			<div class="header-content">
				<h3 id="damage-title" class="header-title">Assign Combat Damage</h3>
				<p class="header-subtitle">{prompt.message || 'Distribute damage among blockers'}</p>
			</div>
		</div>

		<!-- Attacker Info -->
		<div class="attacker-info">
			<div class="attacker-card">
				<span class="attacker-name">{prompt.attackerName}</span>
				<span class="attacker-power">{prompt.attackerPower} Power</span>
			</div>
			<div class="damage-counter">
				<span
					class="remaining"
					class:over={remainingDamage() < 0}
					class:exact={remainingDamage() === 0}
				>
					{remainingDamage() >= 0 ? remainingDamage() : Math.abs(remainingDamage())}
				</span>
				<span class="remaining-label">
					{#if remainingDamage() > 0}
						remaining
					{:else if remainingDamage() < 0}
						over!
					{:else}
						✓ all assigned
					{/if}
				</span>
			</div>
		</div>

		<!-- Blockers List -->
		<div class="blockers-list">
			<h4 class="list-title">Blockers (in damage order)</h4>

			{#each sortedBlockers() as blocker, index}
				{@const damage = getDamage(blocker.id)}
				{@const lethal = getLethalDamage(blocker)}
				{@const isLethal = damage >= lethal}
				<div class="blocker-row" class:has-lethal={isLethal}>
					<div class="blocker-order">{index + 1}</div>
					<div class="blocker-info">
						<span class="blocker-name">{blocker.name}</span>
						<span class="blocker-stats">
							Toughness: {blocker.toughness}
							{#if blocker.damage > 0}
								<span class="marked-damage">(already marked: {blocker.damage})</span>
							{/if}
						</span>
						<span class="lethal-info">
							Lethal: {lethal}
							{#if isLethal}
								<span class="lethal-check">✓</span>
							{/if}
						</span>
					</div>
					<div class="damage-controls">
						<button
							class="btn-adjust"
							onclick={() => adjustDamage(blocker.id, -1)}
							disabled={isSubmitting || damage <= 0}
							aria-label="Decrease damage"
						>
							−
						</button>
						<input
							type="number"
							class="damage-input"
							value={damage}
							min="0"
							onchange={(e) => setDamage(blocker.id, parseInt(e.currentTarget.value) || 0)}
							disabled={isSubmitting}
						/>
						<button
							class="btn-adjust"
							onclick={() => adjustDamage(blocker.id, 1)}
							disabled={isSubmitting}
							aria-label="Increase damage"
						>
							+
						</button>
					</div>
				</div>
			{/each}

			<!-- Trample row -->
			{#if prompt.hasTrample}
				{@const trampleDamage = getDamage(prompt.defendingPlayerId)}
				<div class="blocker-row trample-row">
					<div class="blocker-order">⚡</div>
					<div class="blocker-info">
						<span class="blocker-name">{prompt.defendingPlayerName} (Trample)</span>
						<span class="blocker-stats">Excess damage to defending player</span>
					</div>
					<div class="damage-controls">
						<button
							class="btn-adjust"
							onclick={() => adjustDamage(prompt.defendingPlayerId, -1)}
							disabled={isSubmitting || trampleDamage <= 0}
							aria-label="Decrease trample damage"
						>
							−
						</button>
						<input
							type="number"
							class="damage-input trample"
							value={trampleDamage}
							min="0"
							onchange={(e) =>
								setDamage(prompt.defendingPlayerId, parseInt(e.currentTarget.value) || 0)}
							disabled={isSubmitting}
						/>
						<button
							class="btn-adjust"
							onclick={() => adjustDamage(prompt.defendingPlayerId, 1)}
							disabled={isSubmitting}
							aria-label="Increase trample damage"
						>
							+
						</button>
					</div>
				</div>
			{/if}
		</div>

		<!-- Validation Error -->
		{#if validationError}
			<div class="error-message" role="alert">
				{validationError}
			</div>
		{/if}

		<!-- Actions -->
		<div class="damage-actions">
			<div class="action-hint">
				<kbd>Enter</kbd> to confirm
			</div>
			<button
				class="btn-confirm"
				onclick={handleSubmit}
				disabled={isSubmitting || remainingDamage() !== 0}
				type="button"
			>
				{#if isSubmitting}
					<span class="spinner"></span>
				{/if}
				Confirm Damage
			</button>
		</div>
	</div>
</div>

<style>
	/* Overlay */
	.assign-damage-overlay {
		position: fixed;
		inset: 0;
		background: rgba(0, 0, 0, 0.75);
		z-index: 100;
		display: flex;
		align-items: center;
		justify-content: center;
		padding: 1rem;
	}

	/* Panel */
	.damage-panel {
		background: #1f2937;
		border: 2px solid #ef4444;
		border-radius: 12px;
		width: 100%;
		max-width: 500px;
		max-height: 90vh;
		overflow: auto;
		animation: panel-appear 0.2s ease-out;
	}

	@keyframes panel-appear {
		from {
			opacity: 0;
			transform: scale(0.95);
		}
		to {
			opacity: 1;
			transform: scale(1);
		}
	}

	/* Header */
	.damage-header {
		display: flex;
		align-items: center;
		gap: 1rem;
		padding: 1rem 1.5rem;
		background: linear-gradient(180deg, rgba(239, 68, 68, 0.2) 0%, transparent 100%);
		border-bottom: 1px solid rgba(239, 68, 68, 0.3);
	}

	.header-icon {
		font-size: 2rem;
	}

	.header-content {
		flex: 1;
	}

	.header-title {
		margin: 0 0 0.25rem 0;
		font-size: 1.125rem;
		font-weight: 700;
		color: #ef4444;
	}

	.header-subtitle {
		margin: 0;
		font-size: 0.875rem;
		color: #9ca3af;
	}

	/* Attacker Info */
	.attacker-info {
		display: flex;
		justify-content: space-between;
		align-items: center;
		padding: 1rem 1.5rem;
		background: rgba(239, 68, 68, 0.1);
		border-bottom: 1px solid rgba(239, 68, 68, 0.2);
	}

	.attacker-card {
		display: flex;
		flex-direction: column;
		gap: 0.25rem;
	}

	.attacker-name {
		font-weight: 600;
		color: #ffffff;
		font-size: 1rem;
	}

	.attacker-power {
		color: #f87171;
		font-weight: 700;
		font-size: 1.25rem;
	}

	.damage-counter {
		display: flex;
		flex-direction: column;
		align-items: center;
		padding: 0.5rem 1rem;
		background: #111827;
		border-radius: 8px;
		border: 1px solid #374151;
	}

	.remaining {
		font-size: 1.5rem;
		font-weight: 700;
		color: #fbbf24;
	}

	.remaining.exact {
		color: #22c55e;
	}

	.remaining.over {
		color: #ef4444;
	}

	.remaining-label {
		font-size: 0.75rem;
		color: #6b7280;
	}

	/* Blockers List */
	.blockers-list {
		padding: 1rem 1.5rem;
	}

	.list-title {
		margin: 0 0 1rem 0;
		font-size: 0.75rem;
		font-weight: 600;
		color: #9ca3af;
		text-transform: uppercase;
		letter-spacing: 0.5px;
	}

	.blocker-row {
		display: flex;
		align-items: center;
		gap: 0.75rem;
		padding: 0.75rem;
		background: #111827;
		border: 1px solid #374151;
		border-radius: 8px;
		margin-bottom: 0.5rem;
		transition: all 0.2s;
	}

	.blocker-row.has-lethal {
		border-color: #22c55e;
		background: rgba(34, 197, 94, 0.1);
	}

	.blocker-row.trample-row {
		border-color: #f59e0b;
		background: rgba(245, 158, 11, 0.1);
	}

	.blocker-order {
		width: 2rem;
		height: 2rem;
		display: flex;
		align-items: center;
		justify-content: center;
		background: #374151;
		border-radius: 50%;
		font-weight: 700;
		font-size: 0.875rem;
		color: #e5e7eb;
	}

	.blocker-info {
		flex: 1;
		display: flex;
		flex-direction: column;
		gap: 0.125rem;
	}

	.blocker-name {
		font-weight: 600;
		color: #ffffff;
	}

	.blocker-stats {
		font-size: 0.75rem;
		color: #9ca3af;
	}

	.marked-damage {
		color: #f87171;
	}

	.lethal-info {
		font-size: 0.75rem;
		color: #fbbf24;
	}

	.lethal-check {
		color: #22c55e;
		margin-left: 0.25rem;
	}

	/* Damage Controls */
	.damage-controls {
		display: flex;
		align-items: center;
		gap: 0.25rem;
	}

	.btn-adjust {
		width: 2rem;
		height: 2rem;
		display: flex;
		align-items: center;
		justify-content: center;
		background: #374151;
		border: 1px solid #4b5563;
		border-radius: 4px;
		color: #e5e7eb;
		font-size: 1.25rem;
		font-weight: 700;
		cursor: pointer;
		transition: all 0.15s;
	}

	.btn-adjust:hover:not(:disabled) {
		background: #4b5563;
		border-color: #6b7280;
	}

	.btn-adjust:disabled {
		opacity: 0.4;
		cursor: not-allowed;
	}

	.damage-input {
		width: 3rem;
		height: 2rem;
		padding: 0 0.25rem;
		background: #1f2937;
		border: 2px solid #4b5563;
		border-radius: 4px;
		color: #ffffff;
		font-size: 1rem;
		font-weight: 700;
		text-align: center;
		-moz-appearance: textfield;
	}

	.damage-input::-webkit-outer-spin-button,
	.damage-input::-webkit-inner-spin-button {
		-webkit-appearance: none;
		margin: 0;
	}

	.damage-input:focus {
		outline: none;
		border-color: #ef4444;
	}

	.damage-input.trample {
		border-color: #f59e0b;
	}

	.damage-input.trample:focus {
		border-color: #fbbf24;
	}

	/* Error Message */
	.error-message {
		margin: 0 1.5rem 1rem;
		padding: 0.75rem;
		background: rgba(239, 68, 68, 0.2);
		border: 1px solid #ef4444;
		border-radius: 6px;
		color: #fca5a5;
		font-size: 0.875rem;
	}

	/* Actions */
	.damage-actions {
		display: flex;
		justify-content: space-between;
		align-items: center;
		padding: 1rem 1.5rem;
		border-top: 1px solid #374151;
	}

	.action-hint {
		font-size: 0.75rem;
		color: #6b7280;
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

	.btn-confirm {
		display: flex;
		align-items: center;
		gap: 0.5rem;
		padding: 0.625rem 1.5rem;
		background: linear-gradient(135deg, #ef4444 0%, #dc2626 100%);
		border: none;
		border-radius: 6px;
		color: white;
		font-weight: 600;
		font-size: 0.875rem;
		cursor: pointer;
		transition: all 0.2s;
	}

	.btn-confirm:hover:not(:disabled) {
		background: linear-gradient(135deg, #f87171 0%, #ef4444 100%);
		transform: translateY(-1px);
		box-shadow: 0 4px 12px rgba(239, 68, 68, 0.4);
	}

	.btn-confirm:disabled {
		opacity: 0.5;
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
	@media (max-width: 480px) {
		.damage-panel {
			max-height: 100vh;
			border-radius: 0;
		}

		.attacker-info {
			flex-direction: column;
			gap: 0.75rem;
			text-align: center;
		}

		.blocker-row {
			flex-wrap: wrap;
		}

		.blocker-info {
			flex: 1 1 100%;
			order: 2;
		}

		.damage-controls {
			order: 1;
			margin-left: auto;
		}
	}
</style>
