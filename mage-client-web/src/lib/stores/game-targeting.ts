/**
 * Game Targeting Store
 * Manages target selection mode for spells and abilities
 */

import { writable, derived, get } from 'svelte/store';
import { pendingPrompt } from './game';
import type { CardView } from '$lib/generated/mage/v1/models';
import type { GameTargetData } from '$lib/generated/mage/v1/websocket';

/**
 * Targeting mode state
 */
export interface TargetingState {
	// Whether targeting mode is active
	isActive: boolean;
	
	// The prompt message from the server
	message: string;
	
	// Valid target card IDs (from server)
	validTargetIds: Set<string>;
	
	// Currently selected target IDs (user selections)
	selectedTargetIds: string[];
	
	// Minimum number of targets required
	minTargets: number;
	
	// Maximum number of targets allowed
	maxTargets: number;
	
	// Whether targeting is required (can't cancel)
	required: boolean;
	
	// Full target card data for display
	validTargets: CardView[];
}

const initialState: TargetingState = {
	isActive: false,
	message: '',
	validTargetIds: new Set(),
	selectedTargetIds: [],
	minTargets: 1,
	maxTargets: 1,
	required: false,
	validTargets: []
};

/**
 * Create the targeting store
 */
function createTargetingStore() {
	const { subscribe, set, update } = writable<TargetingState>(initialState);

	/**
	 * Enter targeting mode with server-provided target data
	 */
	function enterTargetingMode(data: GameTargetData) {
		const validIds = new Set<string>(data.targets.map((t) => t.id));
		
		// Parse min/max from options if provided
		const minTargets = parseInt(data.options['minTargets'] || '1', 10);
		const maxTargets = parseInt(data.options['maxTargets'] || '1', 10);

		set({
			isActive: true,
			message: data.message,
			validTargetIds: validIds,
			selectedTargetIds: [],
			minTargets,
			maxTargets,
			required: data.required,
			validTargets: data.targets
		});
	}

	/**
	 * Exit targeting mode and reset state
	 */
	function exitTargetingMode() {
		set(initialState);
	}

	/**
	 * Toggle a target selection
	 * Returns true if the target was toggled successfully
	 */
	function toggleTarget(targetId: string): boolean {
		let success = false;
		
		update((state) => {
			// Must be a valid target
			if (!state.validTargetIds.has(targetId)) {
				return state;
			}

			const isSelected = state.selectedTargetIds.includes(targetId);

			if (isSelected) {
				// Remove from selection
				success = true;
				return {
					...state,
					selectedTargetIds: state.selectedTargetIds.filter((id) => id !== targetId)
				};
			} else {
				// Add to selection if under max
				if (state.selectedTargetIds.length < state.maxTargets) {
					success = true;
					return {
						...state,
						selectedTargetIds: [...state.selectedTargetIds, targetId]
					};
				} else if (state.maxTargets === 1) {
					// Single target mode - replace selection
					success = true;
					return {
						...state,
						selectedTargetIds: [targetId]
					};
				}
			}

			return state;
		});

		return success;
	}

	/**
	 * Select a single target (clears previous selections)
	 */
	function selectTarget(targetId: string): boolean {
		const state = get({ subscribe });
		
		if (!state.validTargetIds.has(targetId)) {
			return false;
		}

		update((s) => ({
			...s,
			selectedTargetIds: [targetId]
		}));

		return true;
	}

	/**
	 * Clear all selected targets
	 */
	function clearSelection() {
		update((state) => ({
			...state,
			selectedTargetIds: []
		}));
	}

	/**
	 * Check if the current selection is valid for confirmation
	 */
	function canConfirm(): boolean {
		const state = get({ subscribe });
		return (
			state.selectedTargetIds.length >= state.minTargets &&
			state.selectedTargetIds.length <= state.maxTargets
		);
	}

	/**
	 * Get the selected target IDs
	 */
	function getSelectedTargets(): string[] {
		return get({ subscribe }).selectedTargetIds;
	}

	return {
		subscribe,
		enterTargetingMode,
		exitTargetingMode,
		toggleTarget,
		selectTarget,
		clearSelection,
		canConfirm,
		getSelectedTargets
	};
}

/**
 * Global targeting store instance
 */
export const targetingStore = createTargetingStore();

// Derived stores for convenient access

/**
 * Whether targeting mode is currently active
 */
export const isTargetingActive = derived(targetingStore, ($targeting) => $targeting.isActive);

/**
 * Current targeting message
 */
export const targetingMessage = derived(targetingStore, ($targeting) => $targeting.message);

/**
 * Set of valid target IDs
 */
export const validTargetIds = derived(targetingStore, ($targeting) => $targeting.validTargetIds);

/**
 * Array of selected target IDs
 */
export const selectedTargetIds = derived(
	targetingStore,
	($targeting) => $targeting.selectedTargetIds
);

/**
 * Whether the current selection can be confirmed
 */
export const canConfirmTargets = derived(targetingStore, ($targeting) => {
	return (
		$targeting.selectedTargetIds.length >= $targeting.minTargets &&
		$targeting.selectedTargetIds.length <= $targeting.maxTargets
	);
});

/**
 * Whether targeting is required (can't cancel)
 */
export const isTargetingRequired = derived(targetingStore, ($targeting) => $targeting.required);

/**
 * Check if a specific card is a valid target
 */
export function isValidTarget(cardId: string): boolean {
	const state = get(targetingStore);
	return state.isActive && state.validTargetIds.has(cardId);
}

/**
 * Check if a specific card is currently selected as a target
 */
export function isTargetSelected(cardId: string): boolean {
	const state = get(targetingStore);
	return state.selectedTargetIds.includes(cardId);
}

/**
 * Auto-sync with game store's pending prompt
 * When GAME_TARGET prompt comes in, activate targeting mode
 */
export function syncWithGamePrompt() {
	return pendingPrompt.subscribe((prompt) => {
		if (prompt?.type === 'target') {
			const targetData = prompt.data as GameTargetData;
			targetingStore.enterTargetingMode(targetData);
		} else {
			// Clear targeting when prompt changes to something else
			const state = get(targetingStore);
			if (state.isActive) {
				targetingStore.exitTargetingMode();
			}
		}
	});
}

