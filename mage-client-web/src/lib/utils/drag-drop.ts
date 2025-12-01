/**
 * Drag and Drop Utilities
 * Core utilities for card drag-and-drop functionality with position tracking and drop zone detection
 */

import { writable, derived, get } from 'svelte/store';

/**
 * Valid zones for card drops
 */
export type DropZone = 'battlefield' | 'graveyard' | 'exile' | 'hand' | 'none';

/**
 * Source zones where cards can be dragged from
 */
export type SourceZone = 'hand' | 'battlefield';

/**
 * Drag state interface
 */
export interface DragState {
	isDragging: boolean;
	cardId: string | null;
	cardName: string | null;
	sourceZone: SourceZone | null;
	position: { x: number; y: number };
	startPosition: { x: number; y: number };
	validDropZones: DropZone[];
	currentDropZone: DropZone;
	isOverValidZone: boolean;
}

/**
 * Drop zone registration
 */
export interface DropZoneConfig {
	id: string;
	type: DropZone;
	element: HTMLElement;
	accepts: (cardId: string, sourceZone: SourceZone) => boolean;
	onDrop: (cardId: string) => void;
}

/**
 * Initial drag state
 */
const initialDragState: DragState = {
	isDragging: false,
	cardId: null,
	cardName: null,
	sourceZone: null,
	position: { x: 0, y: 0 },
	startPosition: { x: 0, y: 0 },
	validDropZones: [],
	currentDropZone: 'none',
	isOverValidZone: false
};

/**
 * Create the drag-drop store
 */
function createDragDropStore() {
	const { subscribe, set, update } = writable<DragState>(initialDragState);

	// Registry of drop zones
	let dropZones: Map<string, DropZoneConfig> = new Map();

	/**
	 * Register a drop zone
	 */
	function registerDropZone(config: DropZoneConfig): () => void {
		dropZones.set(config.id, config);
		// Return unregister function
		return () => {
			dropZones.delete(config.id);
		};
	}

	/**
	 * Get all registered drop zones
	 */
	function getDropZones(): DropZoneConfig[] {
		return Array.from(dropZones.values());
	}

	/**
	 * Start dragging a card
	 */
	function startDrag(
		cardId: string,
		cardName: string,
		sourceZone: SourceZone,
		startX: number,
		startY: number,
		validZones: DropZone[] = ['battlefield']
	): void {
		update((state) => ({
			...state,
			isDragging: true,
			cardId,
			cardName,
			sourceZone,
			position: { x: startX, y: startY },
			startPosition: { x: startX, y: startY },
			validDropZones: validZones,
			currentDropZone: 'none',
			isOverValidZone: false
		}));

		// Add document-level event listeners for drag tracking
		document.addEventListener('mousemove', handleMouseMove);
		document.addEventListener('mouseup', handleMouseUp);
		document.addEventListener('touchmove', handleTouchMove, { passive: false });
		document.addEventListener('touchend', handleTouchEnd);

		// Prevent text selection during drag
		document.body.style.userSelect = 'none';
	}

	/**
	 * Update drag position
	 */
	function updatePosition(x: number, y: number): void {
		const state = get({ subscribe });
		if (!state.isDragging) return;

		// Check which drop zone we're over
		const currentZone = detectDropZone(x, y, state.validDropZones);
		const isValid = currentZone !== 'none' && state.validDropZones.includes(currentZone);

		update((s) => ({
			...s,
			position: { x, y },
			currentDropZone: currentZone,
			isOverValidZone: isValid
		}));
	}

	/**
	 * Detect which drop zone contains the given coordinates
	 */
	function detectDropZone(x: number, y: number, validZones: DropZone[]): DropZone {
		const zones = Array.from(dropZones.values());
		for (const zone of zones) {
			if (!validZones.includes(zone.type)) continue;

			const rect = zone.element.getBoundingClientRect();
			if (x >= rect.left && x <= rect.right && y >= rect.top && y <= rect.bottom) {
				const state = get({ subscribe });
				if (state.cardId && state.sourceZone && zone.accepts(state.cardId, state.sourceZone)) {
					return zone.type;
				}
			}
		}
		return 'none';
	}

	/**
	 * End drag operation
	 */
	function endDrag(): { cardId: string; dropZone: DropZone; wasValid: boolean } | null {
		const state = get({ subscribe });

		// Remove event listeners
		document.removeEventListener('mousemove', handleMouseMove);
		document.removeEventListener('mouseup', handleMouseUp);
		document.removeEventListener('touchmove', handleTouchMove);
		document.removeEventListener('touchend', handleTouchEnd);

		// Restore text selection
		document.body.style.userSelect = '';

		if (!state.isDragging || !state.cardId) {
			set(initialDragState);
			return null;
		}

		const result = {
			cardId: state.cardId,
			dropZone: state.currentDropZone,
			wasValid: state.isOverValidZone
		};

		// Trigger drop handler if over valid zone
		if (state.isOverValidZone && state.currentDropZone !== 'none') {
			const zones = Array.from(dropZones.values());
			for (const zone of zones) {
				if (zone.type === state.currentDropZone) {
					const rect = zone.element.getBoundingClientRect();
					const { x, y } = state.position;
					if (
						x >= rect.left &&
						x <= rect.right &&
						y >= rect.top &&
						y <= rect.bottom &&
						state.sourceZone &&
						zone.accepts(state.cardId, state.sourceZone)
					) {
						zone.onDrop(state.cardId);
						break;
					}
				}
			}
		}

		// Reset state
		set(initialDragState);

		return result;
	}

	/**
	 * Cancel drag without dropping
	 */
	function cancelDrag(): void {
		// Remove event listeners
		document.removeEventListener('mousemove', handleMouseMove);
		document.removeEventListener('mouseup', handleMouseUp);
		document.removeEventListener('touchmove', handleTouchMove);
		document.removeEventListener('touchend', handleTouchEnd);

		// Restore text selection
		document.body.style.userSelect = '';

		set(initialDragState);
	}

	// Internal event handlers
	function handleMouseMove(e: MouseEvent): void {
		updatePosition(e.clientX, e.clientY);
	}

	function handleMouseUp(): void {
		endDrag();
	}

	function handleTouchMove(e: TouchEvent): void {
		e.preventDefault(); // Prevent scrolling while dragging
		if (e.touches.length > 0) {
			updatePosition(e.touches[0].clientX, e.touches[0].clientY);
		}
	}

	function handleTouchEnd(): void {
		endDrag();
	}

	return {
		subscribe,
		startDrag,
		updatePosition,
		endDrag,
		cancelDrag,
		registerDropZone,
		getDropZones
	};
}

/**
 * Global drag-drop store instance
 */
export const dragDropStore = createDragDropStore();

// Derived stores for convenient access
export const isDragging = derived(dragDropStore, ($state) => $state.isDragging);
export const draggedCardId = derived(dragDropStore, ($state) => $state.cardId);
export const draggedCardName = derived(dragDropStore, ($state) => $state.cardName);
export const dragPosition = derived(dragDropStore, ($state) => $state.position);
export const isOverValidDropZone = derived(dragDropStore, ($state) => $state.isOverValidZone);
export const currentDropZone = derived(dragDropStore, ($state) => $state.currentDropZone);

/**
 * Check if a card can be played from hand
 * This is a utility function to determine valid drop zones based on card type
 */
export function getValidDropZonesForCard(
	cardType: string,
	phase: string,
	hasPriority: boolean
): DropZone[] {
	if (!hasPriority) return [];

	const isLand = cardType.toLowerCase().includes('land');
	const isMainPhase = phase.includes('MAIN') || phase === 'PRECOMBAT_MAIN' || phase === 'POSTCOMBAT_MAIN';

	// Lands can only be played during main phase
	if (isLand) {
		return isMainPhase ? ['battlefield'] : [];
	}

	// Spells can be cast whenever player has priority (with mana)
	// The actual validation happens server-side
	return ['battlefield'];
}

/**
 * Calculate distance between two points
 */
export function getDistance(
	p1: { x: number; y: number },
	p2: { x: number; y: number }
): number {
	return Math.sqrt(Math.pow(p2.x - p1.x, 2) + Math.pow(p2.y - p1.y, 2));
}

/**
 * Check if drag has moved enough to be considered a real drag (not a click)
 */
export function isDragThresholdMet(
	start: { x: number; y: number },
	current: { x: number; y: number },
	threshold: number = 5
): boolean {
	return getDistance(start, current) >= threshold;
}

