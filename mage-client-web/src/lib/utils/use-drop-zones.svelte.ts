/**
 * Shared drop zone setup for game pages
 * Eliminates duplication of $effect blocks across playtest and multiplayer pages
 */

import { dragDropStore, type SourceZone } from './drag-drop';

export interface DropZoneRefs {
	battlefield: HTMLDivElement | null;
	graveyard: HTMLElement | null;
	exile: HTMLElement | null;
	hand: HTMLElement | null;
	library: HTMLElement | null;
	command: HTMLElement | null;
}

export interface DropZoneHandlers {
	onBattlefieldDrop: (cardId: string) => void;
	onGraveyardDrop: (cardId: string) => void;
	onExileDrop: (cardId: string) => void;
	onHandDrop: (cardId: string) => void;
	onLibraryDrop: (cardId: string) => void;
	onCommandDrop: (cardId: string) => void;
}

/**
 * Sets up all 6 drop zones with automatic cleanup
 * Call once in component with refs getter and handlers
 */
export function useDropZones(
	getRefs: () => DropZoneRefs,
	handlers: DropZoneHandlers
): void {
	// Battlefield drop zone
	let battlefieldUnregister: (() => void) | null = null;
	$effect(() => {
		const refs = getRefs();
		if (refs.battlefield && !battlefieldUnregister) {
			battlefieldUnregister = dragDropStore.registerDropZone({
				id: 'battlefield',
				type: 'battlefield',
				element: refs.battlefield,
				accepts: (_cardId: string, sourceZone: SourceZone) => sourceZone !== 'battlefield',
				onDrop: handlers.onBattlefieldDrop
			});
		}
		return () => {
			if (battlefieldUnregister) {
				battlefieldUnregister();
				battlefieldUnregister = null;
			}
		};
	});

	// Graveyard drop zone
	let graveyardUnregister: (() => void) | null = null;
	$effect(() => {
		const refs = getRefs();
		if (refs.graveyard && !graveyardUnregister) {
			graveyardUnregister = dragDropStore.registerDropZone({
				id: 'graveyard',
				type: 'graveyard',
				element: refs.graveyard,
				accepts: (_cardId: string, sourceZone: SourceZone) => sourceZone !== 'graveyard',
				onDrop: handlers.onGraveyardDrop
			});
		}
		return () => {
			if (graveyardUnregister) {
				graveyardUnregister();
				graveyardUnregister = null;
			}
		};
	});

	// Exile drop zone
	let exileUnregister: (() => void) | null = null;
	$effect(() => {
		const refs = getRefs();
		if (refs.exile && !exileUnregister) {
			exileUnregister = dragDropStore.registerDropZone({
				id: 'exile',
				type: 'exile',
				element: refs.exile,
				accepts: (_cardId: string, sourceZone: SourceZone) => sourceZone !== 'exile',
				onDrop: handlers.onExileDrop
			});
		}
		return () => {
			if (exileUnregister) {
				exileUnregister();
				exileUnregister = null;
			}
		};
	});

	// Hand drop zone
	let handUnregister: (() => void) | null = null;
	$effect(() => {
		const refs = getRefs();
		if (refs.hand && !handUnregister) {
			handUnregister = dragDropStore.registerDropZone({
				id: 'hand',
				type: 'hand',
				element: refs.hand,
				accepts: (_cardId: string, sourceZone: SourceZone) => sourceZone !== 'hand',
				onDrop: handlers.onHandDrop
			});
		}
		return () => {
			if (handUnregister) {
				handUnregister();
				handUnregister = null;
			}
		};
	});

	// Library drop zone
	let libraryUnregister: (() => void) | null = null;
	$effect(() => {
		const refs = getRefs();
		if (refs.library && !libraryUnregister) {
			libraryUnregister = dragDropStore.registerDropZone({
				id: 'library',
				type: 'library',
				element: refs.library,
				accepts: (_cardId: string, sourceZone: SourceZone) => sourceZone !== 'library',
				onDrop: handlers.onLibraryDrop
			});
		}
		return () => {
			if (libraryUnregister) {
				libraryUnregister();
				libraryUnregister = null;
			}
		};
	});

	// Command drop zone
	let commandUnregister: (() => void) | null = null;
	$effect(() => {
		const refs = getRefs();
		if (refs.command && !commandUnregister) {
			commandUnregister = dragDropStore.registerDropZone({
				id: 'command',
				type: 'command',
				element: refs.command,
				accepts: (_cardId: string, sourceZone: SourceZone) => sourceZone !== 'command',
				onDrop: handlers.onCommandDrop
			});
		}
		return () => {
			if (commandUnregister) {
				commandUnregister();
				commandUnregister = null;
			}
		};
	});
}
