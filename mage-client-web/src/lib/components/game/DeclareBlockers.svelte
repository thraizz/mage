<!--
  DeclareBlockers.svelte
  
  Combat component for declaring blockers during the Declare Blockers step.
  Uses two-step selection: click a blocker, then click an attacker to assign.
  
  Features:
  - Two-step blocker assignment
  - Support multiple blockers per attacker
  - Visual feedback for selection state
  - Keyboard shortcuts: Enter to confirm, Escape to decline
-->
<script lang="ts">
	import { onMount, onDestroy } from 'svelte';
	import { declareBlockers, declineToBlock } from '$lib/api/game';
	import {
		combatStore,
		assignedBlockerIds,
		selectedBlockerId as selectedBlockerIdStore
	} from '$lib/stores/combat';
	import type { DeclaredAttacker, ParsedCombatOptions } from '$lib/types/combat';
	import type { CardView } from '$lib/generated/mage/v1/models';

	// Props
	let {
		gameId,
		options,
		battlefieldCards,
		attackingCreatures,
		onComplete = () => {}
	}: {
		gameId: string;
		options: ParsedCombatOptions;
		battlefieldCards: CardView[];
		attackingCreatures: DeclaredAttacker[];
		onComplete?: () => void;
	} = $props();

	// State
	let isSubmitting = $state(false);
	let error = $state<string | null>(null);

	// Derived state from combat store
	const assignedIds = $derived($assignedBlockerIds);
	const selectedBlockerId = $derived($selectedBlockerIdStore);

	// Build card name maps
	const cardNames = $derived.by(() => {
		const map = new Map<string, string>();
		for (const card of battlefieldCards) {
			map.set(card.id, card.name);
		}
		return map;
	});

	// Get attackers by ID for quick lookup
	const attackersById = $derived.by(() => {
		const map = new Map<string, DeclaredAttacker>();
		for (const attacker of attackingCreatures) {
			map.set(attacker.cardId, attacker);
		}
		return map;
	});

	// Get available blockers (cards that can block)
	const availableBlockerIds = $derived.by(() => {
		const ids = new Set<string>();
		for (const opt of options.blockOptions) {
			ids.add(opt.blockerId);
		}
		return ids;
	});

	// Get attackers that a specific blocker can block
	function getBlockableAttackers(blockerId: string): string[] {
		return options.blockOptions
			.filter((opt) => opt.blockerId === blockerId)
			.map((opt) => opt.attackerId);
	}

	// Initialize combat store with parsed options
	$effect(() => {
		combatStore.enterDeclareBlockersPhase(options, cardNames(), attackingCreatures);
	});

	/**
	 * Handle clicking on a blocker
	 */
	function handleBlockerClick(blockerId: string) {
		if (isSubmitting) return;
		if (!availableBlockerIds().has(blockerId)) return;

		// If already blocking, clicking again removes the assignment
		if (assignedIds.has(blockerId)) {
			combatStore.removeBlockAssignment(blockerId);
			error = null;
			return;
		}

		// Select this blocker
		combatStore.selectBlocker(blockerId);
		error = null;
	}

	/**
	 * Handle clicking on an attacker (to assign selected blocker)
	 */
	function handleAttackerClick(attackerId: string) {
		if (isSubmitting) return;
		if (!selectedBlockerId) return;

		// Check if selected blocker can block this attacker
		const canBlock = getBlockableAttackers(selectedBlockerId).includes(attackerId);
		if (!canBlock) {
			error = `${cardNames().get(selectedBlockerId) || 'Selected creature'} cannot block ${cardNames().get(attackerId) || 'that attacker'}`;
			return;
		}

		// Assign the block
		combatStore.assignBlocker(attackerId);
		error = null;
	}

	/**
	 * Cancel selection
	 */
	function handleCancelSelection() {
		combatStore.selectBlocker(null);
	}

	/**
	 * Submit all declared blockers
	 */
	async function handleConfirm() {
		if (isSubmitting) return;

		isSubmitting = true;
		error = null;

		try {
			const blockers = combatStore.getBlockAssignments();
			await declareBlockers(gameId, blockers);
			combatStore.reset();
			onComplete();
		} catch (err) {
			error = err instanceof Error ? err.message : 'Failed to declare blockers';
			console.error('Failed to declare blockers:', err);
		} finally {
			isSubmitting = false;
		}
	}

	/**
	 * Decline to block (no blockers)
	 */
	async function handleDecline() {
		if (isSubmitting) return;

		isSubmitting = true;
		error = null;

		try {
			await declineToBlock(gameId);
			combatStore.reset();
			onComplete();
		} catch (err) {
			error = err instanceof Error ? err.message : 'Failed to decline blocking';
			console.error('Failed to decline blocking:', err);
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
				if (selectedBlockerId) {
					handleCancelSelection();
				} else {
					handleDecline();
				}
				break;
		}
	}

	// Get block assignments grouped by attacker
	const blocksByAttacker = $derived.by(() => {
		const assignments = combatStore.getBlockAssignments();
		const grouped = new Map<string, string[]>();

		for (const { blockerId, attackerId } of assignments) {
			const blockers = grouped.get(attackerId) || [];
			blockers.push(blockerId);
			grouped.set(attackerId, blockers);
		}

		return grouped;
	});

	// Get the attacker that a blocker is assigned to
	function getAssignedAttacker(blockerId: string): string | undefined {
		const assignments = combatStore.getBlockAssignments();
		const assignment = assignments.find((a) => a.blockerId === blockerId);
		return assignment?.attackerId;
	}

	// Add global keyboard listener
	onMount(() => {
		window.addEventListener('keydown', handleKeydown);
	});

	onDestroy(() => {
		window.removeEventListener('keydown', handleKeydown);
	});
</script>

<!-- Overlay -->
<div
	class="declare-blockers-overlay"
	role="dialog"
	aria-modal="true"
	aria-labelledby="blockers-title"
>
	<!-- Top Banner -->
	<div class="blockers-banner">
		<div class="banner-icon">🛡️</div>
		<div class="banner-content">
			<h3 id="blockers-title" class="banner-title">Declare Blockers</h3>
			{#if selectedBlockerId}
				<p class="banner-description select-attacker">
					Now click an <strong>attacker</strong> to block with {cardNames().get(
						selectedBlockerId
					) || 'selected creature'}
				</p>
			{:else}
				<p class="banner-description">Click a creature you control to select it as a blocker</p>
			{/if}
			<div class="banner-status">
				{#if assignedIds.size > 0}
					<span class="blocker-count"
						>{assignedIds.size} blocker{assignedIds.size !== 1 ? 's' : ''} assigned</span
					>
				{:else}
					<span class="blocker-hint">No blockers assigned yet</span>
				{/if}
			</div>
		</div>
		{#if selectedBlockerId}
			<button class="cancel-selection" onclick={handleCancelSelection}> Cancel Selection </button>
		{/if}
	</div>

	<!-- Combat Layout -->
	<div class="combat-layout">
		<!-- Attackers Section -->
		<div class="attackers-section">
			<h4 class="section-title">Attacking Creatures</h4>
			<div class="attacker-cards">
				{#each attackingCreatures as attacker}
					{@const isBlockable =
						selectedBlockerId && getBlockableAttackers(selectedBlockerId).includes(attacker.cardId)}
					{@const blockers = blocksByAttacker().get(attacker.cardId) || []}
					<button
						class="attacker-card"
						class:blockable={isBlockable}
						class:has-blockers={blockers.length > 0}
						class:selecting={selectedBlockerId !== null}
						onclick={() => handleAttackerClick(attacker.cardId)}
						disabled={isSubmitting || !selectedBlockerId}
					>
						<div class="card-name">{attacker.cardName}</div>
						<div class="attack-target">→ {attacker.defenderName}</div>
						{#if blockers.length > 0}
							<div class="blockers-assigned">
								<span class="blocked-label">Blocked by:</span>
								{#each blockers as blockerId}
									<span class="blocker-chip">{cardNames().get(blockerId) || 'Unknown'}</span>
								{/each}
							</div>
						{:else if isBlockable}
							<div class="block-prompt">Click to block</div>
						{:else if selectedBlockerId}
							<div class="cannot-block">Cannot be blocked</div>
						{/if}
					</button>
				{/each}
			</div>
		</div>

		<!-- Blockers Section -->
		<div class="blockers-section">
			<h4 class="section-title">Your Creatures</h4>
			<div class="blocker-cards">
				{#each battlefieldCards.filter((c) => availableBlockerIds().has(c.id)) as card}
					{@const isSelected = selectedBlockerId === card.id}
					{@const isAssigned = assignedIds.has(card.id)}
					{@const assignedTo = getAssignedAttacker(card.id)}
					<button
						class="blocker-card"
						class:selected={isSelected}
						class:assigned={isAssigned}
						onclick={() => handleBlockerClick(card.id)}
						disabled={isSubmitting}
					>
						<div class="card-name">{card.name}</div>
						{#if card.power && card.toughness}
							<div class="card-stats">{card.power}/{card.toughness}</div>
						{/if}
						{#if isAssigned && assignedTo}
							<div class="assigned-to">
								Blocking: {cardNames().get(assignedTo) || 'Unknown'}
							</div>
						{:else if isSelected}
							<div class="select-hint">Select an attacker to block</div>
						{/if}
					</button>
				{/each}
				{#if availableBlockerIds().size === 0}
					<div class="no-blockers">No creatures can block</div>
				{/if}
			</div>
		</div>
	</div>

	<!-- Error Message -->
	{#if error}
		<div class="error-message" role="alert">
			{error}
		</div>
	{/if}

	<!-- Bottom Action Bar -->
	<div class="blockers-actions">
		<div class="action-hints">
			<span class="hint">
				<kbd>ESC</kbd>
				{selectedBlockerId ? 'cancel selection' : 'no blocks'}
			</span>
			<span class="hint">
				<kbd>Enter</kbd> to confirm
			</span>
		</div>
		<div class="action-buttons">
			<button class="btn-decline" onclick={handleDecline} disabled={isSubmitting} type="button">
				{#if isSubmitting}
					<span class="spinner"></span>
				{/if}
				No Blocks
			</button>
			<button class="btn-confirm" onclick={handleConfirm} disabled={isSubmitting} type="button">
				{#if isSubmitting}
					<span class="spinner"></span>
				{/if}
				{#if assignedIds.size > 0}
					Confirm Blockers ({assignedIds.size})
				{:else}
					Confirm (No Blockers)
				{/if}
			</button>
		</div>
	</div>
</div>

<style>
	/* Overlay */
	.declare-blockers-overlay {
		position: fixed;
		inset: 0;
		background: rgba(0, 0, 0, 0.6);
		z-index: 100;
		display: flex;
		flex-direction: column;
		pointer-events: auto;
	}

	/* Top Banner */
	.blockers-banner {
		display: flex;
		align-items: center;
		gap: 1rem;
		padding: 1rem 1.5rem;
		background: linear-gradient(180deg, rgba(59, 130, 246, 0.2) 0%, transparent 100%);
		border-bottom: 2px solid rgba(59, 130, 246, 0.5);
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
		animation: shield-pulse 1.5s ease-in-out infinite;
	}

	@keyframes shield-pulse {
		0%,
		100% {
			transform: scale(1);
		}
		50% {
			transform: scale(1.1);
		}
	}

	.banner-content {
		flex: 1;
	}

	.banner-title {
		margin: 0 0 0.25rem 0;
		font-size: 1.25rem;
		font-weight: 700;
		color: #3b82f6;
		text-transform: uppercase;
		letter-spacing: 1px;
	}

	.banner-description {
		margin: 0 0 0.5rem 0;
		font-size: 0.9375rem;
		color: #ffffff;
	}

	.banner-description.select-attacker {
		color: #fbbf24;
	}

	.banner-status {
		font-size: 0.875rem;
	}

	.blocker-count {
		color: #3b82f6;
		font-weight: 600;
	}

	.blocker-hint {
		color: #9ca3af;
		font-style: italic;
	}

	.cancel-selection {
		padding: 0.5rem 1rem;
		background: rgba(239, 68, 68, 0.2);
		border: 1px solid rgba(239, 68, 68, 0.5);
		border-radius: 6px;
		color: #fca5a5;
		font-size: 0.875rem;
		cursor: pointer;
		transition: all 0.2s;
	}

	.cancel-selection:hover {
		background: rgba(239, 68, 68, 0.3);
		border-color: #ef4444;
	}

	/* Combat Layout */
	.combat-layout {
		flex: 1;
		display: grid;
		grid-template-rows: 1fr 1fr;
		gap: 1rem;
		padding: 1rem;
		overflow: auto;
	}

	.section-title {
		margin: 0 0 0.75rem 0;
		font-size: 0.875rem;
		font-weight: 600;
		text-transform: uppercase;
		letter-spacing: 0.5px;
	}

	.attackers-section .section-title {
		color: #ef4444;
	}

	.blockers-section .section-title {
		color: #3b82f6;
	}

	.attacker-cards,
	.blocker-cards {
		display: flex;
		flex-wrap: wrap;
		gap: 0.75rem;
	}

	/* Attacker Cards */
	.attacker-card {
		min-width: 150px;
		padding: 0.75rem;
		background: rgba(239, 68, 68, 0.15);
		border: 2px solid rgba(239, 68, 68, 0.3);
		border-radius: 8px;
		text-align: left;
		cursor: default;
		transition: all 0.2s;
	}

	.attacker-card.selecting {
		cursor: pointer;
	}

	.attacker-card.blockable {
		border-color: #22c55e;
		background: rgba(34, 197, 94, 0.15);
		box-shadow: 0 0 15px rgba(34, 197, 94, 0.3);
		animation: blockable-pulse 1s ease-in-out infinite;
	}

	@keyframes blockable-pulse {
		0%,
		100% {
			box-shadow: 0 0 15px rgba(34, 197, 94, 0.3);
		}
		50% {
			box-shadow: 0 0 25px rgba(34, 197, 94, 0.5);
		}
	}

	.attacker-card.blockable:hover:not(:disabled) {
		transform: translateY(-2px);
		box-shadow: 0 0 25px rgba(34, 197, 94, 0.5);
	}

	.attacker-card.has-blockers {
		border-color: #3b82f6;
		background: rgba(59, 130, 246, 0.15);
	}

	.attacker-card .card-name {
		font-weight: 600;
		color: #ffffff;
		margin-bottom: 0.25rem;
	}

	.attack-target {
		font-size: 0.75rem;
		color: #f87171;
	}

	.blockers-assigned {
		margin-top: 0.5rem;
		padding-top: 0.5rem;
		border-top: 1px solid rgba(59, 130, 246, 0.3);
	}

	.blocked-label {
		font-size: 0.6875rem;
		color: #93c5fd;
		display: block;
		margin-bottom: 0.25rem;
	}

	.blocker-chip {
		display: inline-block;
		padding: 0.125rem 0.5rem;
		background: rgba(59, 130, 246, 0.3);
		border-radius: 9999px;
		font-size: 0.6875rem;
		color: #93c5fd;
		margin-right: 0.25rem;
	}

	.block-prompt {
		margin-top: 0.5rem;
		font-size: 0.75rem;
		color: #22c55e;
		font-style: italic;
	}

	.cannot-block {
		margin-top: 0.5rem;
		font-size: 0.75rem;
		color: #6b7280;
		font-style: italic;
	}

	/* Blocker Cards */
	.blocker-card {
		min-width: 130px;
		padding: 0.75rem;
		background: rgba(59, 130, 246, 0.1);
		border: 2px solid rgba(59, 130, 246, 0.3);
		border-radius: 8px;
		text-align: left;
		cursor: pointer;
		transition: all 0.2s;
	}

	.blocker-card:hover:not(:disabled) {
		border-color: #3b82f6;
		background: rgba(59, 130, 246, 0.2);
		transform: translateY(-2px);
	}

	.blocker-card.selected {
		border-color: #fbbf24;
		background: rgba(251, 191, 36, 0.15);
		box-shadow: 0 0 20px rgba(251, 191, 36, 0.4);
	}

	.blocker-card.assigned {
		border-color: #22c55e;
		background: rgba(34, 197, 94, 0.15);
	}

	.blocker-card .card-name {
		font-weight: 600;
		color: #ffffff;
	}

	.blocker-card .card-stats {
		font-size: 0.875rem;
		color: #fbbf24;
		font-weight: 600;
	}

	.assigned-to {
		margin-top: 0.5rem;
		padding-top: 0.5rem;
		border-top: 1px solid rgba(34, 197, 94, 0.3);
		font-size: 0.75rem;
		color: #86efac;
	}

	.select-hint {
		margin-top: 0.5rem;
		font-size: 0.75rem;
		color: #fbbf24;
		font-style: italic;
	}

	.no-blockers {
		padding: 1rem;
		color: #6b7280;
		font-style: italic;
	}

	/* Error Message */
	.error-message {
		margin: 0.5rem 1rem;
		padding: 0.75rem 1rem;
		background: rgba(239, 68, 68, 0.2);
		border: 1px solid #ef4444;
		border-radius: 6px;
		color: #fca5a5;
		font-size: 0.875rem;
	}

	/* Bottom Action Bar */
	.blockers-actions {
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

	.btn-decline,
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

	.btn-decline {
		background: #374151;
		color: #e5e7eb;
	}

	.btn-decline:hover:not(:disabled) {
		background: #4b5563;
	}

	.btn-confirm {
		background: linear-gradient(135deg, #3b82f6 0%, #2563eb 100%);
		color: white;
	}

	.btn-confirm:hover:not(:disabled) {
		background: linear-gradient(135deg, #60a5fa 0%, #3b82f6 100%);
		transform: translateY(-1px);
		box-shadow: 0 4px 12px rgba(59, 130, 246, 0.4);
	}

	.btn-decline:disabled,
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
	@media (max-width: 768px) {
		.combat-layout {
			grid-template-rows: auto auto;
		}

		.blockers-actions {
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
