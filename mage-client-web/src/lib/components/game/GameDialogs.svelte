<script lang="ts">
	import type { CardView, CounterView } from '$lib/generated/mage/v1/models';
	import type { MenuAction } from './DeckContextMenu.svelte';

	// Dialog components
	import TokenCreator from './TokenCreator.svelte';
	import CreateTokenDialog from './CreateTokenDialog.svelte';
	import CounterDialog from './CounterDialog.svelte';
	import LibrarySearch from './LibrarySearch.svelte';
	import PlaytestLibrarySearch from './PlaytestLibrarySearch.svelte';
	import DeckContextMenu from './DeckContextMenu.svelte';
	import NumberInputDialog from './NumberInputDialog.svelte';
	import ScryDialog from './ScryDialog.svelte';
	import RevealTopDialog from './RevealTopDialog.svelte';
	import KeyboardShortcutsModal from './KeyboardShortcutsModal.svelte';
	import GameChatOverlay from './GameChatOverlay.svelte';

	// Types
	interface Player {
		playerId: string;
		name: string;
		library: CardView[];
		libraryCount: number;
		hand: CardView[];
	}

	interface GameUIState {
		showTokenCreator: boolean;
		showCreateTokenDialog: boolean;
		showCounterDialog: boolean;
		selectedCardForCounters: { id: string; name: string } | null;
		showDeckSearch: boolean;
		showDeckContextMenu: boolean;
		deckContextMenuPosition: { x: number; y: number };
		showNumberInputDialog: boolean;
		numberInputDialogConfig: {
			title: string;
			defaultValue: number;
			min: number;
			max: number;
			onConfirm: (value: number) => void;
		} | null;
		showScryDialog: boolean;
		currentScrySession: {
			cards: CardView[];
		} | null;
		showRevealTopDialog: boolean;
		revealedCards: CardView[];
		showKeyboardShortcuts: boolean;
	}

	interface Props {
		// UI State
		uiState: GameUIState;

		// Game data
		gameId: string;
		me: Player | null;
		selectedCardForCountersData: CardView | null;
		deckContextMenuActions: MenuAction[];

		// Event handlers
		onCreateToken: (name: string, types: string, power: string, toughness: string, color: string) => void;
		onAddCounter: (cardId: string, counterName: string, amount: number) => void;
		onRemoveCounter: (cardId: string, counterName: string, amount: number) => void;
		onSetCounter: (cardId: string, counterName: string, amount: number) => void;
		onScryComplete: (keepOnTop: CardView[], putToBottom: CardView[]) => void;
		onNumberConfirm: (value: number) => void;

		// Variants
		keyboardShortcutsMode: 'game' | 'playtest';
		librarySearchVariant: 'multiplayer' | 'playtest';

		// Multiplayer-specific (optional)
		showGameChatOverlay?: boolean;

		// Playtest-specific library search handlers (optional)
		onLibraryMove?: (cardId: string, zone: 'HAND' | 'BATTLEFIELD' | 'GRAVEYARD' | 'EXILE') => void;
		onLibraryShuffle?: () => void;
		onLibraryClose?: () => void;

		// Multiplayer-specific library search handlers (optional)
		onLibraryComplete?: () => void;
		onLibraryCancel?: () => void;
	}

	let {
		uiState,
		gameId,
		me,
		selectedCardForCountersData,
		deckContextMenuActions,
		onCreateToken,
		onAddCounter,
		onRemoveCounter,
		onSetCounter,
		onScryComplete,
		onNumberConfirm,
		keyboardShortcutsMode,
		librarySearchVariant,
		showGameChatOverlay = false,
		onLibraryMove,
		onLibraryShuffle,
		onLibraryClose,
		onLibraryComplete,
		onLibraryCancel
	}: Props = $props();
</script>

<!-- Token Creator -->
{#if uiState.showTokenCreator}
	<TokenCreator {gameId} onClose={() => (uiState.showTokenCreator = false)} />
{/if}

<!-- Create Token Dialog -->
{#if uiState.showCreateTokenDialog}
	<CreateTokenDialog
		{onCreateToken}
		onClose={() => (uiState.showCreateTokenDialog = false)}
	/>
{/if}

<!-- Counter Dialog -->
{#if uiState.showCounterDialog && uiState.selectedCardForCounters && selectedCardForCountersData}
	<CounterDialog
		cardName={selectedCardForCountersData.name}
		cardId={selectedCardForCountersData.id}
		currentCounters={selectedCardForCountersData.counters}
		onAddCounter={(counterName, amount) => {
			const card = selectedCardForCountersData;
			onAddCounter(card.id, counterName, amount);
		}}
		onRemoveCounter={(counterName, amount) => {
			const card = selectedCardForCountersData;
			onRemoveCounter(card.id, counterName, amount);
		}}
		onSetCounter={(counterName, amount) => {
			const card = selectedCardForCountersData;
			onSetCounter(card.id, counterName, amount);
		}}
		onClose={() => {
			uiState.showCounterDialog = false;
			uiState.selectedCardForCounters = null;
		}}
	/>
{/if}

<!-- Deck Search - Multiplayer variant -->
{#if librarySearchVariant === 'multiplayer' && uiState.showDeckSearch && me}
	<LibrarySearch
		{gameId}
		librarySearchData={{
			playerId: me.playerId,
			message: 'Search your library',
			destination: 'hand',
			cards: me.library,
			canCancel: true
		}}
		onComplete={onLibraryComplete || (() => (uiState.showDeckSearch = false))}
		onCancel={onLibraryCancel || (() => (uiState.showDeckSearch = false))}
	/>
{/if}

<!-- Deck Search - Playtest variant -->
{#if librarySearchVariant === 'playtest' && uiState.showDeckSearch && me}
	<PlaytestLibrarySearch
		cards={me.library}
		playerName="You"
		onMove={onLibraryMove || (() => {})}
		onShuffle={onLibraryShuffle}
		onClose={onLibraryClose || (() => (uiState.showDeckSearch = false))}
	/>
{/if}

<!-- Deck Context Menu -->
{#if uiState.showDeckContextMenu}
	<DeckContextMenu
		position={uiState.deckContextMenuPosition}
		deckCount={me?.libraryCount || 0}
		playerName={me?.name || 'You'}
		onClose={() => (uiState.showDeckContextMenu = false)}
		actions={deckContextMenuActions}
	/>
{/if}

<!-- Number Input Dialog -->
{#if uiState.showNumberInputDialog && uiState.numberInputDialogConfig}
	<NumberInputDialog
		title={uiState.numberInputDialogConfig.title}
		defaultValue={uiState.numberInputDialogConfig.defaultValue}
		min={uiState.numberInputDialogConfig.min}
		max={uiState.numberInputDialogConfig.max}
		onConfirm={onNumberConfirm}
		onCancel={() => {
			uiState.showNumberInputDialog = false;
			uiState.numberInputDialogConfig = null;
		}}
	/>
{/if}

<!-- Scry Dialog -->
{#if uiState.showScryDialog && uiState.currentScrySession}
	<ScryDialog
		cards={uiState.currentScrySession.cards}
		onComplete={onScryComplete}
		onCancel={() => {
			uiState.showScryDialog = false;
			uiState.currentScrySession = null;
		}}
	/>
{/if}

<!-- Reveal Top Dialog -->
{#if uiState.showRevealTopDialog}
	<RevealTopDialog
		cards={uiState.revealedCards}
		onClose={() => {
			uiState.showRevealTopDialog = false;
			uiState.revealedCards = [];
		}}
	/>
{/if}

<!-- Keyboard Shortcuts Modal -->
<KeyboardShortcutsModal bind:open={uiState.showKeyboardShortcuts} mode={keyboardShortcutsMode} />

<!-- Game Chat Overlay (Multiplayer only) -->
{#if showGameChatOverlay}
	<GameChatOverlay {gameId} />
{/if}
