/**
 * Visual Stack Store
 * UI for displaying the MTG stack. Syncs with server state
 * and provides reordering/notes as visual aids.
 */

import { writable, derived } from 'svelte/store';
import type { SourceZone } from '$lib/utils/drag-drop';

/**
 * Represents a single item on the visual stack
 */
export interface VisualStackItem {
  localId: string; // Unique local ID (allows same card multiple times)
  cardId: string; // Original card ID from the game
  cardName: string;
  imageUrl?: string;
  sourceZone: SourceZone;
  addedAt: number; // Timestamp for ordering
  note?: string; // Optional player note (e.g., "ETB trigger")
  controllerId?: string; // Player who controls this item
}

/**
 * Internal state for the visual stack
 */
interface VisualStackState {
  items: VisualStackItem[];
  isOpen: boolean;
}

const initialState: VisualStackState = {
  items: [],
  isOpen: false
};

/**
 * Generate a unique local ID for stack items
 */
function generateLocalId(): string {
  return `vs-${Date.now()}-${Math.random().toString(36).substr(2, 9)}`;
}

/**
 * Create the visual stack store
 */
function createVisualStackStore() {
  const { subscribe, set, update } = writable<VisualStackState>(initialState);

  return {
    subscribe,

    /**
     * Add an item to the top of the stack
     * @param cardId - The card ID from the game
     * @param cardName - Display name of the card
     * @param sourceZone - Zone the card came from
     * @param options - Additional options including localId from server
     * @returns The localId of the added item
     */
    addItem(
      cardId: string,
      cardName: string,
      sourceZone: SourceZone,
      options?: {
        imageUrl?: string;
        note?: string;
        controllerId?: string;
        /** Server-provided stack item ID - use this as localId for server sync */
        localId?: string;
      }
    ): string {
      // Use server-provided ID if available, otherwise generate one
      const localId = options?.localId || generateLocalId();
      update((state) => ({
        ...state,
        items: [
          ...state.items,
          {
            localId,
            cardId,
            cardName,
            sourceZone,
            addedAt: Date.now(),
            imageUrl: options?.imageUrl,
            note: options?.note,
            controllerId: options?.controllerId
          }
        ]
      }));
      return localId;
    },

    /**
     * Remove an item from the stack by its local ID
     */
    removeItem(localId: string): void {
      update((state) => ({
        ...state,
        items: state.items.filter((item) => item.localId !== localId)
      }));
    },

    /**
     * Remove the top item from the stack (resolve)
     */
    resolveTop(): VisualStackItem | null {
      let removedItem: VisualStackItem | null = null;
      update((state) => {
        if (state.items.length === 0) return state;
        removedItem = state.items[state.items.length - 1];
        return {
          ...state,
          items: state.items.slice(0, -1)
        };
      });
      return removedItem;
    },

    /**
     * Reorder items by moving an item from one index to another
     */
    reorderItems(fromIndex: number, toIndex: number): void {
      update((state) => {
        if (fromIndex < 0 || fromIndex >= state.items.length) return state;
        if (toIndex < 0 || toIndex >= state.items.length) return state;
        if (fromIndex === toIndex) return state;

        const newItems = [...state.items];
        const [movedItem] = newItems.splice(fromIndex, 1);
        newItems.splice(toIndex, 0, movedItem);

        return {
          ...state,
          items: newItems
        };
      });
    },

    /**
     * Update the note on a stack item
     */
    updateNote(localId: string, note: string): void {
      update((state) => ({
        ...state,
        items: state.items.map((item) => (item.localId === localId ? { ...item, note } : item))
      }));
    },

    /**
     * Clear all items from the stack
     */
    clearStack(): void {
      update((state) => ({
        ...state,
        items: []
      }));
    },

    /**
     * Toggle the stack panel visibility
     */
    toggleOpen(): void {
      update((state) => ({
        ...state,
        isOpen: !state.isOpen
      }));
    },

    /**
     * Set the stack panel visibility
     */
    setOpen(isOpen: boolean): void {
      update((state) => ({
        ...state,
        isOpen
      }));
    },

    /**
     * Reset to initial state
     */
    reset(): void {
      set(initialState);
    }
  };
}

/**
 * Global visual stack store instance
 */
export const visualStackStore = createVisualStackStore();

// Derived stores for convenient access
export const visualStackItems = derived(visualStackStore, ($state) => $state.items);

export const visualStackIsOpen = derived(visualStackStore, ($state) => $state.isOpen);

export const visualStackCount = derived(visualStackStore, ($state) => $state.items.length);

export const visualStackIsEmpty = derived(visualStackStore, ($state) => $state.items.length === 0);

/**
 * Get the top item (the one that would resolve next)
 */
export const visualStackTop = derived(visualStackStore, ($state) =>
  $state.items.length > 0 ? $state.items[$state.items.length - 1] : null
);
