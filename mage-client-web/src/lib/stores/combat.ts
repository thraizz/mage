/**
 * Combat Store
 * Manages combat phase UI state for declare attackers, blockers, and damage assignment
 */

import { writable, derived, get } from 'svelte/store';
import type {
	CombatStoreState,
	CombatPhase,
	AvailableAttacker,
	AvailableBlocker,
	DeclaredAttacker,
	DefenderTarget,
	DamageAssignmentPrompt,
	DamageAssignment,
	ParsedCombatOptions
} from '$lib/types/combat';
import {
	parseCombatOptions,
	groupAttackOptionsByCard,
	groupBlockOptionsByBlocker
} from '$lib/types/combat';

const initialState: CombatStoreState = {
	phase: 'idle',
	availableAttackers: new Map(),
	declaredAttackers: new Map(),
	defenders: new Map(),
	attackingCreatures: new Map(),
	availableBlockers: new Map(),
	blockAssignments: new Map(),
	pendingDamageAssignment: null,
	damageAssignments: [],
	selectedBlockerId: null,
	isSubmitting: false
};

/**
 * Create the combat store
 */
function createCombatStore() {
	const { subscribe, set, update } = writable<CombatStoreState>(initialState);

	/**
	 * Enter declare attackers phase with parsed options from server
	 */
	function enterDeclareAttackersPhase(
		options: ParsedCombatOptions,
		cardNames: Map<string, string>,
		defenderInfo: DefenderTarget[]
	) {
		const attackerGroups = groupAttackOptionsByCard(options.attackOptions);
		const availableAttackers = new Map<string, AvailableAttacker>();

		for (const [cardId, validDefenders] of attackerGroups) {
			availableAttackers.set(cardId, {
				cardId,
				cardName: cardNames.get(cardId) || 'Unknown',
				validDefenders
			});
		}

		const defenders = new Map<string, DefenderTarget>();
		for (const defender of defenderInfo) {
			defenders.set(defender.id, defender);
		}

		update((state) => ({
			...state,
			phase: 'declare-attackers',
			availableAttackers,
			declaredAttackers: new Map(),
			defenders,
			isSubmitting: false
		}));
	}

	/**
	 * Enter declare blockers phase with parsed options from server
	 */
	function enterDeclareBlockersPhase(
		options: ParsedCombatOptions,
		cardNames: Map<string, string>,
		attackerInfo: DeclaredAttacker[]
	) {
		const blockerGroups = groupBlockOptionsByBlocker(options.blockOptions);
		const availableBlockers = new Map<string, AvailableBlocker>();

		for (const [blockerId, canBlock] of blockerGroups) {
			availableBlockers.set(blockerId, {
				cardId: blockerId,
				cardName: cardNames.get(blockerId) || 'Unknown',
				canBlock
			});
		}

		const attackingCreatures = new Map<string, DeclaredAttacker>();
		for (const attacker of attackerInfo) {
			attackingCreatures.set(attacker.cardId, attacker);
		}

		update((state) => ({
			...state,
			phase: 'declare-blockers',
			attackingCreatures,
			availableBlockers,
			blockAssignments: new Map(),
			selectedBlockerId: null,
			isSubmitting: false
		}));
	}

	/**
	 * Enter damage assignment phase
	 */
	function enterDamageAssignmentPhase(prompt: DamageAssignmentPrompt) {
		// Initialize damage assignments with lethal to first blocker
		const initialAssignments: DamageAssignment[] = [];
		let remainingPower = prompt.attackerPower;

		for (const blocker of prompt.blockers) {
			const lethalDamage = Math.max(0, blocker.toughness - blocker.damage);
			const assignedDamage = Math.min(remainingPower, lethalDamage);
			initialAssignments.push({
				targetId: blocker.id,
				targetType: 'creature',
				damage: assignedDamage
			});
			remainingPower -= assignedDamage;
		}

		// Assign remaining to defending player if trample
		if (prompt.hasTrample && remainingPower > 0) {
			initialAssignments.push({
				targetId: prompt.defendingPlayerId,
				targetType: 'player',
				damage: remainingPower
			});
		}

		update((state) => ({
			...state,
			phase: 'assign-damage',
			pendingDamageAssignment: prompt,
			damageAssignments: initialAssignments,
			isSubmitting: false
		}));
	}

	/**
	 * Toggle an attacker declaration
	 */
	function toggleAttacker(cardId: string, defenderId?: string): boolean {
		let success = false;

		update((state) => {
			if (state.phase !== 'declare-attackers') return state;

			const attacker = state.availableAttackers.get(cardId);
			if (!attacker) return state;

			const newDeclared = new Map(state.declaredAttackers);

			if (newDeclared.has(cardId)) {
				// Remove attacker
				newDeclared.delete(cardId);
				success = true;
			} else {
				// Add attacker - use provided defender or first valid one
				const targetDefender = defenderId || attacker.validDefenders[0];
				if (targetDefender && attacker.validDefenders.includes(targetDefender)) {
					newDeclared.set(cardId, targetDefender);
					success = true;
				}
			}

			return {
				...state,
				declaredAttackers: newDeclared
			};
		});

		return success;
	}

	/**
	 * Change the defender target for a declared attacker
	 */
	function changeAttackTarget(cardId: string, defenderId: string): boolean {
		let success = false;

		update((state) => {
			if (state.phase !== 'declare-attackers') return state;

			const attacker = state.availableAttackers.get(cardId);
			if (!attacker || !attacker.validDefenders.includes(defenderId)) return state;

			if (!state.declaredAttackers.has(cardId)) return state;

			const newDeclared = new Map(state.declaredAttackers);
			newDeclared.set(cardId, defenderId);
			success = true;

			return {
				...state,
				declaredAttackers: newDeclared
			};
		});

		return success;
	}

	/**
	 * Select a blocker for assignment
	 */
	function selectBlocker(blockerId: string | null) {
		update((state) => {
			if (state.phase !== 'declare-blockers') return state;

			// Toggle off if same blocker selected
			if (state.selectedBlockerId === blockerId) {
				return { ...state, selectedBlockerId: null };
			}

			// Verify it's a valid blocker
			if (blockerId && !state.availableBlockers.has(blockerId)) {
				return state;
			}

			return { ...state, selectedBlockerId: blockerId };
		});
	}

	/**
	 * Assign selected blocker to an attacker
	 */
	function assignBlocker(attackerId: string): boolean {
		let success = false;

		update((state) => {
			if (state.phase !== 'declare-blockers') return state;
			if (!state.selectedBlockerId) return state;

			const blocker = state.availableBlockers.get(state.selectedBlockerId);
			if (!blocker || !blocker.canBlock.includes(attackerId)) return state;

			const newAssignments = new Map(state.blockAssignments);
			newAssignments.set(state.selectedBlockerId, attackerId);
			success = true;

			return {
				...state,
				blockAssignments: newAssignments,
				selectedBlockerId: null
			};
		});

		return success;
	}

	/**
	 * Remove a block assignment
	 */
	function removeBlockAssignment(blockerId: string): boolean {
		let success = false;

		update((state) => {
			if (state.phase !== 'declare-blockers') return state;

			if (!state.blockAssignments.has(blockerId)) return state;

			const newAssignments = new Map(state.blockAssignments);
			newAssignments.delete(blockerId);
			success = true;

			return {
				...state,
				blockAssignments: newAssignments
			};
		});

		return success;
	}

	/**
	 * Update damage assignment for a target
	 */
	function updateDamageAssignment(targetId: string, damage: number): boolean {
		let success = false;

		update((state) => {
			if (state.phase !== 'assign-damage') return state;
			if (!state.pendingDamageAssignment) return state;

			const newAssignments = state.damageAssignments.map((a) =>
				a.targetId === targetId ? { ...a, damage: Math.max(0, damage) } : a
			);
			success = true;

			return {
				...state,
				damageAssignments: newAssignments
			};
		});

		return success;
	}

	/**
	 * Validate damage assignment follows MTG rules
	 * - Total damage must equal attacker power
	 * - Must assign lethal damage to blocker before moving to next in order
	 */
	function validateDamageAssignment(): { valid: boolean; error?: string } {
		const state = get({ subscribe });
		if (!state.pendingDamageAssignment) {
			return { valid: false, error: 'No damage assignment pending' };
		}

		const prompt = state.pendingDamageAssignment;
		const assignments = state.damageAssignments;

		// Check total damage equals attacker power
		const totalDamage = assignments.reduce((sum, a) => sum + a.damage, 0);
		if (totalDamage !== prompt.attackerPower) {
			return {
				valid: false,
				error: `Total damage (${totalDamage}) must equal attacker power (${prompt.attackerPower})`
			};
		}

		// Check lethal damage ordering for blockers
		let mustAssignLethal = true;
		for (const blocker of prompt.blockers) {
			const assignment = assignments.find((a) => a.targetId === blocker.id);
			const assignedDamage = assignment?.damage || 0;
			const lethalDamage = Math.max(0, blocker.toughness - blocker.damage);

			if (mustAssignLethal && assignedDamage < lethalDamage) {
				// Check if there's damage assigned to later blockers or player
				const laterAssignments = assignments.filter((a) => {
					if (a.targetType === 'player') return true;
					const laterBlocker = prompt.blockers.find((b) => b.id === a.targetId);
					return laterBlocker && laterBlocker.order > blocker.order;
				});

				const laterDamage = laterAssignments.reduce((sum, a) => sum + a.damage, 0);
				if (laterDamage > 0) {
					return {
						valid: false,
						error: `Must assign lethal damage (${lethalDamage}) to ${blocker.name} before assigning to later blockers`
					};
				}
			}

			if (assignedDamage >= lethalDamage) {
				// This blocker has lethal, continue to next
			} else {
				// Not lethal, can't assign to later blockers
				mustAssignLethal = false;
			}
		}

		return { valid: true };
	}

	/**
	 * Set submitting state
	 */
	function setSubmitting(submitting: boolean) {
		update((state) => ({ ...state, isSubmitting: submitting }));
	}

	/**
	 * Reset combat state
	 */
	function reset() {
		set(initialState);
	}

	/**
	 * Get declared attackers as array for API submission
	 */
	function getDeclaredAttackers(): Array<{ cardId: string; defenderId: string }> {
		const state = get({ subscribe });
		return Array.from(state.declaredAttackers.entries()).map(([cardId, defenderId]) => ({
			cardId,
			defenderId
		}));
	}

	/**
	 * Get block assignments as array for API submission
	 */
	function getBlockAssignments(): Array<{ blockerId: string; attackerId: string }> {
		const state = get({ subscribe });
		return Array.from(state.blockAssignments.entries()).map(([blockerId, attackerId]) => ({
			blockerId,
			attackerId
		}));
	}

	/**
	 * Get damage assignments for API submission
	 */
	function getDamageAssignments(): DamageAssignment[] {
		return get({ subscribe }).damageAssignments;
	}

	return {
		subscribe,
		enterDeclareAttackersPhase,
		enterDeclareBlockersPhase,
		enterDamageAssignmentPhase,
		toggleAttacker,
		changeAttackTarget,
		selectBlocker,
		assignBlocker,
		removeBlockAssignment,
		updateDamageAssignment,
		validateDamageAssignment,
		setSubmitting,
		reset,
		getDeclaredAttackers,
		getBlockAssignments,
		getDamageAssignments
	};
}

/**
 * Global combat store instance
 */
export const combatStore = createCombatStore();

// Derived stores for convenient access

/**
 * Current combat phase
 */
export const combatPhase = derived(combatStore, ($combat) => $combat.phase);

/**
 * Whether in any combat selection phase
 */
export const isInCombat = derived(combatStore, ($combat) => $combat.phase !== 'idle');

/**
 * Whether in declare attackers phase
 */
export const isDeclaringAttackers = derived(
	combatStore,
	($combat) => $combat.phase === 'declare-attackers'
);

/**
 * Whether in declare blockers phase
 */
export const isDeclaringBlockers = derived(
	combatStore,
	($combat) => $combat.phase === 'declare-blockers'
);

/**
 * Whether in damage assignment phase
 */
export const isAssigningDamage = derived(
	combatStore,
	($combat) => $combat.phase === 'assign-damage'
);

/**
 * Set of card IDs that can attack
 */
export const canAttackCardIds = derived(
	combatStore,
	($combat) => new Set($combat.availableAttackers.keys())
);

/**
 * Set of card IDs that are declared as attackers
 */
export const declaredAttackerIds = derived(
	combatStore,
	($combat) => new Set($combat.declaredAttackers.keys())
);

/**
 * Number of declared attackers
 */
export const declaredAttackerCount = derived(
	combatStore,
	($combat) => $combat.declaredAttackers.size
);

/**
 * Set of card IDs that can block
 */
export const canBlockCardIds = derived(
	combatStore,
	($combat) => new Set($combat.availableBlockers.keys())
);

/**
 * Set of card IDs that have block assignments
 */
export const assignedBlockerIds = derived(
	combatStore,
	($combat) => new Set($combat.blockAssignments.keys())
);

/**
 * Currently selected blocker ID
 */
export const selectedBlockerId = derived(combatStore, ($combat) => $combat.selectedBlockerId);

/**
 * Whether combat submission is in progress
 */
export const isCombatSubmitting = derived(combatStore, ($combat) => $combat.isSubmitting);

/**
 * Check if a card is declared as an attacker
 */
export function isCardAttacking(cardId: string): boolean {
	return get(combatStore).declaredAttackers.has(cardId);
}

/**
 * Check if a card can attack
 */
export function canCardAttack(cardId: string): boolean {
	return get(combatStore).availableAttackers.has(cardId);
}

/**
 * Check if a card is assigned as a blocker
 */
export function isCardBlocking(cardId: string): boolean {
	return get(combatStore).blockAssignments.has(cardId);
}

/**
 * Check if a card can block
 */
export function canCardBlock(cardId: string): boolean {
	return get(combatStore).availableBlockers.has(cardId);
}

/**
 * Get the defender ID that a card is attacking
 */
export function getAttackTarget(cardId: string): string | undefined {
	return get(combatStore).declaredAttackers.get(cardId);
}

/**
 * Get the attacker ID that a card is blocking
 */
export function getBlockTarget(cardId: string): string | undefined {
	return get(combatStore).blockAssignments.get(cardId);
}
