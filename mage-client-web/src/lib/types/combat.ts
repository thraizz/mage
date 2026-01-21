/**
 * Combat System Types
 * Type definitions for MTG combat phases: declare attackers, blockers, and damage assignment
 */

/**
 * Combat phase states
 */
export type CombatPhase = 'idle' | 'declare-attackers' | 'declare-blockers' | 'assign-damage';

/**
 * Defender target (player or planeswalker)
 */
export interface DefenderTarget {
	id: string;
	name: string;
	type: 'player' | 'planeswalker';
	life?: number; // For players
	loyalty?: number; // For planeswalkers
}

/**
 * Available attacker info parsed from server prompt
 */
export interface AvailableAttacker {
	cardId: string;
	cardName: string;
	validDefenders: string[]; // Defender IDs this creature can attack
}

/**
 * Declared attacker with its target
 */
export interface DeclaredAttacker {
	cardId: string;
	cardName: string;
	defenderId: string;
	defenderName: string;
}

/**
 * Available blocker info parsed from server prompt
 */
export interface AvailableBlocker {
	cardId: string;
	cardName: string;
	canBlock: string[]; // Attacker IDs this creature can block
}

/**
 * Block assignment
 */
export interface BlockAssignment {
	blockerId: string;
	blockerName: string;
	attackerId: string;
	attackerName: string;
}

/**
 * Blocker in damage assignment order
 */
export interface OrderedBlocker {
	id: string;
	name: string;
	toughness: number;
	damage: number; // Already marked damage
	order: number; // Position in damage order (0 = first to receive damage)
}

/**
 * Damage assignment prompt from server
 */
export interface DamageAssignmentPrompt {
	attackerId: string;
	attackerName: string;
	attackerPower: number;
	blockers: OrderedBlocker[];
	hasTrample: boolean;
	defendingPlayerId: string;
	defendingPlayerName: string;
}

/**
 * Single damage assignment entry
 */
export interface DamageAssignment {
	targetId: string;
	targetType: 'creature' | 'player';
	damage: number;
}

/**
 * Combat store state
 */
export interface CombatStoreState {
	// Current combat phase
	phase: CombatPhase;

	// Declare Attackers phase
	availableAttackers: Map<string, AvailableAttacker>; // cardId -> attacker info
	declaredAttackers: Map<string, string>; // cardId -> defenderId
	defenders: Map<string, DefenderTarget>; // defenderId -> defender info

	// Declare Blockers phase
	attackingCreatures: Map<string, DeclaredAttacker>; // cardId -> attacker info (from opponent)
	availableBlockers: Map<string, AvailableBlocker>; // cardId -> blocker info
	blockAssignments: Map<string, string>; // blockerId -> attackerId

	// Damage Assignment phase
	pendingDamageAssignment: DamageAssignmentPrompt | null;
	damageAssignments: DamageAssignment[];

	// UI state
	selectedBlockerId: string | null; // Currently selected blocker (for assignment)
	isSubmitting: boolean;
}

/**
 * Parsed combat options from server prompt
 */
export interface ParsedCombatOptions {
	type: 'attack' | 'block' | 'damage' | 'none';
	attackOptions: Array<{ cardId: string; defenderId: string }>;
	blockOptions: Array<{ blockerId: string; attackerId: string }>;
	isDone: boolean; // DONE_ATTACKING or DONE_BLOCKING present
}

/**
 * Parse combat options from prompt string array
 */
export function parseCombatOptions(options: string[]): ParsedCombatOptions {
	const attackOptions: Array<{ cardId: string; defenderId: string }> = [];
	const blockOptions: Array<{ blockerId: string; attackerId: string }> = [];
	let isDone = false;
	let type: 'attack' | 'block' | 'damage' | 'none' = 'none';

	for (const opt of options) {
		if (opt === 'DONE_ATTACKING' || opt === 'DONE_BLOCKING') {
			isDone = true;
			continue;
		}

		if (opt.startsWith('ATTACK:')) {
			const parts = opt.split(':');
			if (parts.length >= 3) {
				attackOptions.push({
					cardId: parts[1],
					defenderId: parts[2]
				});
				type = 'attack';
			}
		} else if (opt.startsWith('BLOCK:')) {
			const parts = opt.split(':');
			if (parts.length >= 3) {
				blockOptions.push({
					blockerId: parts[1],
					attackerId: parts[2]
				});
				type = 'block';
			}
		}
	}

	return {
		type,
		attackOptions,
		blockOptions,
		isDone
	};
}

/**
 * Group attack options by card ID to get all valid defenders for each creature
 */
export function groupAttackOptionsByCard(
	options: Array<{ cardId: string; defenderId: string }>
): Map<string, string[]> {
	const grouped = new Map<string, string[]>();

	for (const opt of options) {
		const existing = grouped.get(opt.cardId) || [];
		existing.push(opt.defenderId);
		grouped.set(opt.cardId, existing);
	}

	return grouped;
}

/**
 * Group block options by blocker ID to get all attackers each creature can block
 */
export function groupBlockOptionsByBlocker(
	options: Array<{ blockerId: string; attackerId: string }>
): Map<string, string[]> {
	const grouped = new Map<string, string[]>();

	for (const opt of options) {
		const existing = grouped.get(opt.blockerId) || [];
		existing.push(opt.attackerId);
		grouped.set(opt.blockerId, existing);
	}

	return grouped;
}
