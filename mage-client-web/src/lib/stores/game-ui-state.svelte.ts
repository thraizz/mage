import type { CardView } from '$lib/generated/mage/v1/models';
import type { ScrySession } from './multiplayer-game';

/**
 * Shared UI state for game pages (playtest and multiplayer)
 *
 * This store manages all UI toggles, dialogs, and overlay states that are common
 * between the playtest and multiplayer game pages. By centralizing this state,
 * we avoid duplication and ensure consistent behavior across both game modes.
 */
export function createGameUIState() {
  // Token creation
  let showTokenCreator = $state(false);
  let showCreateTokenDialog = $state(false);

  // Counter management
  let showCounterDialog = $state(false);
  let selectedCardForCounters = $state<{ id: string; name: string } | null>(null);

  // Modals and overlays
  let showKeyboardShortcuts = $state(false);
  let showAllHands = $state(false);
  let showMenu = $state(false);
  let showDebugOverlay = $state(false);
  let showGameStateLog = $state(false);

  // Card interaction
  let hoveredCardId = $state<string | null>(null);

  // Life menu state (for current player)
  let showLifeMenu = $state(false);

  // Opponent state
  let selectedOpponentId = $state<string | null>(null);
  let showOpponentLifeMenu = $state(false);

  // Library search
  let showDeckSearch = $state(false);

  // Deck context menu
  let showDeckContextMenu = $state(false);
  let deckContextMenuPosition = $state<{ x: number; y: number }>({ x: 0, y: 0 });

  // Number input dialog
  let showNumberInputDialog = $state(false);
  let numberInputDialogConfig = $state<{
    title: string;
    defaultValue: number;
    min: number;
    max: number;
    onConfirm: (value: number) => void;
  } | null>(null);

  // Scry dialog
  let showScryDialog = $state(false);
  let currentScrySession = $state<ScrySession | null>(null);

  // Reveal top dialog
  let showRevealTopDialog = $state(false);
  let revealedCards = $state<CardView[]>([]);

  return {
    // Token creation
    get showTokenCreator() {
      return showTokenCreator;
    },
    set showTokenCreator(value: boolean) {
      showTokenCreator = value;
    },

    get showCreateTokenDialog() {
      return showCreateTokenDialog;
    },
    set showCreateTokenDialog(value: boolean) {
      showCreateTokenDialog = value;
    },

    // Counter management
    get showCounterDialog() {
      return showCounterDialog;
    },
    set showCounterDialog(value: boolean) {
      showCounterDialog = value;
    },

    get selectedCardForCounters() {
      return selectedCardForCounters;
    },
    set selectedCardForCounters(value: { id: string; name: string } | null) {
      selectedCardForCounters = value;
    },

    // Modals and overlays
    get showKeyboardShortcuts() {
      return showKeyboardShortcuts;
    },
    set showKeyboardShortcuts(value: boolean) {
      showKeyboardShortcuts = value;
    },

    get showAllHands() {
      return showAllHands;
    },
    set showAllHands(value: boolean) {
      showAllHands = value;
    },

    get showMenu() {
      return showMenu;
    },
    set showMenu(value: boolean) {
      showMenu = value;
    },

    get showDebugOverlay() {
      return showDebugOverlay;
    },
    set showDebugOverlay(value: boolean) {
      showDebugOverlay = value;
    },

    get showGameStateLog() {
      return showGameStateLog;
    },
    set showGameStateLog(value: boolean) {
      showGameStateLog = value;
    },

    // Card interaction
    get hoveredCardId() {
      return hoveredCardId;
    },
    set hoveredCardId(value: string | null) {
      hoveredCardId = value;
    },

    // Life menu
    get showLifeMenu() {
      return showLifeMenu;
    },
    set showLifeMenu(value: boolean) {
      showLifeMenu = value;
    },

    // Opponent state
    get selectedOpponentId() {
      return selectedOpponentId;
    },
    set selectedOpponentId(value: string | null) {
      selectedOpponentId = value;
    },

    get showOpponentLifeMenu() {
      return showOpponentLifeMenu;
    },
    set showOpponentLifeMenu(value: boolean) {
      showOpponentLifeMenu = value;
    },

    // Library search
    get showDeckSearch() {
      return showDeckSearch;
    },
    set showDeckSearch(value: boolean) {
      showDeckSearch = value;
    },

    // Deck context menu
    get showDeckContextMenu() {
      return showDeckContextMenu;
    },
    set showDeckContextMenu(value: boolean) {
      showDeckContextMenu = value;
    },

    get deckContextMenuPosition() {
      return deckContextMenuPosition;
    },
    set deckContextMenuPosition(value: { x: number; y: number }) {
      deckContextMenuPosition = value;
    },

    // Number input dialog
    get showNumberInputDialog() {
      return showNumberInputDialog;
    },
    set showNumberInputDialog(value: boolean) {
      showNumberInputDialog = value;
    },

    get numberInputDialogConfig() {
      return numberInputDialogConfig;
    },
    set numberInputDialogConfig(
      value: {
        title: string;
        defaultValue: number;
        min: number;
        max: number;
        onConfirm: (value: number) => void;
      } | null
    ) {
      numberInputDialogConfig = value;
    },

    // Scry dialog
    get showScryDialog() {
      return showScryDialog;
    },
    set showScryDialog(value: boolean) {
      showScryDialog = value;
    },

    get currentScrySession() {
      return currentScrySession;
    },
    set currentScrySession(value: ScrySession | null) {
      currentScrySession = value;
    },

    // Reveal top dialog
    get showRevealTopDialog() {
      return showRevealTopDialog;
    },
    set showRevealTopDialog(value: boolean) {
      showRevealTopDialog = value;
    },

    get revealedCards() {
      return revealedCards;
    },
    set revealedCards(value: CardView[]) {
      revealedCards = value;
    }
  };
}

export type GameUIState = ReturnType<typeof createGameUIState>;
