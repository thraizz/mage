/**
 * Shared Game Page Controller
 *
 * Provides all event handlers for both playtest and multiplayer game pages.
 * Eliminates code duplication by accepting any game store with the required interface.
 */

import type { CardView } from '$lib/generated/mage/v1/models';
import type { ScrySession } from '$lib/types/gamestore';
import { toast } from '$lib/stores/toast';
import type { MenuAction } from '$lib/components/game/DeckContextMenu.svelte';

/**
 * Common game store interface shared by both playtestGameStore and multiplayerGameStore
 */
export interface GameStore {
	subscribe: (fn: (state: any) => void) => () => void;
	modifyLife: (playerId: string, delta: number) => void;
	setPlayerCounter: (playerId: string, counterType: string, value: number) => void;
	drawCards: (playerId: string, count: number) => void;
	shuffleLibrary: (playerId: string) => void;
	untapAll: (playerId: string) => void;
	nextTurn: () => void;
	tapCard: (cardId: string, tapped: boolean) => void;
	moveCardToZone: (cardId: string, zone: string) => void;
	millCards: (playerId: string, count: number) => void;
	scryCards: (playerId: string, count: number) => ScrySession | null;
	applyScryDecision: (
		playerId: string,
		scryCount: number,
		keepOnTop: CardView[],
		putToBottom: CardView[]
	) => void;
	revealTopCards: (playerId: string, count: number) => CardView[];
	setRevealedTop: (playerId: string, revealed: boolean) => void;
	createToken: (name: string, types: string, power: any, toughness: any, color: string) => void;
	addCounter: (cardId: string, counterName: string, amount: number) => void;
	removeCounter: (cardId: string, counterName: string, amount: number) => void;
	setCounter: (cardId: string, counterName: string, amount: number) => void;
	mulligan?: (playerId: string) => void | Promise<void>;
	keepHand?: (playerId: string) => void | Promise<void>;
	switchControlSeat?: (playerId: string) => void;
	[key: string]: any; // Allow additional methods from concrete store implementations
}

/**
 * Player interface (common to both stores)
 */
export interface GamePlayer {
	playerId: string;
	name: string;
	life: number;
	poison?: number;
	libraryCount: number;
	handCount?: number;
	hand: CardView[];
	library: CardView[];
	graveyard: CardView[];
	manaPool?: {
		white: number;
		blue: number;
		black: number;
		red: number;
		green: number;
		colorless: number;
	};
	keptHand?: boolean;
	mulliganCount?: number;
	revealedTopCard?: boolean;
}

/**
 * Game state interface (common fields from both stores)
 */
export interface GameState {
	activePlayerId: string;
	activeControlSeat: string;
	players: GamePlayer[];
	battlefield: CardView[];
	turn: number;
}

/**
 * Configuration for creating a game page controller
 */
export interface GamePageControllerConfig {
	gameStore: GameStore;
	getState: () => GameState;
	getLocalPlayer: () => GamePlayer | null;
	getPlayers: () => GamePlayer[];
	getBattlefield: () => CardView[];
}

/**
 * Callbacks for dialog state management
 */
export interface DialogCallbacks {
	setScryDialog: (show: boolean, session: ScrySession | null) => void;
	setRevealTopDialog: (show: boolean, cards: CardView[]) => void;
	setNumberInputDialog: (config: {
		show: boolean;
		title?: string;
		defaultValue?: number;
		onConfirm?: (value: number) => void;
	}) => void;
	setDeckContextMenu: (show: boolean, position?: { x: number; y: number }) => void;
}

/**
 * Creates a game page controller with all event handlers
 */
export function createGamePageController(
	config: GamePageControllerConfig,
	callbacks?: DialogCallbacks
) {
	const { gameStore, getState, getLocalPlayer, getPlayers, getBattlefield } = config;

	/**
	 * Handle life change
	 */
	function handleLifeChange(delta: number, playerId?: string): void {
		const targetPlayerId = playerId || getLocalPlayer()?.playerId;
		if (!targetPlayerId) return;
		gameStore.modifyLife(targetPlayerId, delta);
	}

	/**
	 * Handle poison counter change
	 */
	function handlePoisonChange(delta: number, playerId?: string): void {
		const targetPlayerId = playerId || getLocalPlayer()?.playerId;
		if (!targetPlayerId) return;
		const player = getPlayers().find((p) => p.playerId === targetPlayerId);
		if (!player) return;
		const newValue = Math.max(0, (player.poison || 0) + delta);
		gameStore.setPlayerCounter(targetPlayerId, 'poison', newValue);
	}

	/**
	 * Draw a card
	 */
	function handleDrawCard(): void {
		const me = getLocalPlayer();
		if (!me) return;
		gameStore.drawCards(me.playerId, 1);
		toast.success('Drew a card');
	}

	/**
	 * Shuffle library
	 */
	function handleShuffleLibrary(): void {
		const me = getLocalPlayer();
		if (!me) return;
		gameStore.shuffleLibrary(me.playerId);
		toast.success('Shuffled library');
	}

	/**
	 * Untap all permanents
	 */
	function handleUntapAll(): void {
		const me = getLocalPlayer();
		if (!me) return;
		gameStore.untapAll(me.playerId);
		toast.success('Untapped all');
	}

	/**
	 * Next turn
	 */
	function handleNextTurn(): void {
		gameStore.nextTurn();
		const state = getState();
		const newActivePlayer = getPlayers().find((p) => p.playerId === state.activePlayerId);
		if (newActivePlayer) {
			toast.info(`${newActivePlayer.name}'s turn`);
		}
	}

	/**
	 * Handle mulligan decision
	 */
	function handleMulligan(): void | Promise<void> {
		const me = getLocalPlayer();
		if (!me || !gameStore.mulligan) return;
		return gameStore.mulligan(me.playerId);
	}

	/**
	 * Handle keep hand decision
	 */
	function handleKeepHand(): void | Promise<void> {
		const me = getLocalPlayer();
		if (!me || !gameStore.keepHand) return;
		return gameStore.keepHand(me.playerId);
	}

	/**
	 * Deck context menu handlers
	 */
	function handleDeckContextMenu(event: MouseEvent): void {
		const me = getLocalPlayer();
		if (!me) return;
		callbacks?.setDeckContextMenu(true, { x: event.clientX, y: event.clientY });
	}

	function handleDrawN(count: number): void {
		const me = getLocalPlayer();
		if (!me) return;
		gameStore.drawCards(me.playerId, count);
		toast.success(`Drew ${count} card(s)`);
	}

	function handleMill(count: number): void {
		const me = getLocalPlayer();
		if (!me) return;
		gameStore.millCards(me.playerId, count);
		toast.success(`Milled ${count} card(s)`);
	}

	function handleScry(count: number): void {
		const me = getLocalPlayer();
		if (!me) return;
		const session = gameStore.scryCards(me.playerId, count);
		if (session) {
			callbacks?.setScryDialog(true, session);
		} else {
			toast.error('No cards to scry');
		}
	}

	function handleScryComplete(
		keepOnTop: CardView[],
		putToBottom: CardView[],
		currentSession: ScrySession
	): void {
		const me = getLocalPlayer();
		if (!me) return;
		const scryCount = currentSession.cards.length;
		gameStore.applyScryDecision(me.playerId, scryCount, keepOnTop, putToBottom);
		callbacks?.setScryDialog(false, null);
		toast.success(`Scry ${scryCount} complete`);
	}

	function handleRevealTop(count: number): void {
		const me = getLocalPlayer();
		if (!me) return;
		const cards = gameStore.revealTopCards(me.playerId, count);
		callbacks?.setRevealTopDialog(true, cards);
	}

	function handleToggleRevealedTop(): void {
		const me = getLocalPlayer();
		if (!me) return;
		const willReveal = !me.revealedTopCard;
		gameStore.setRevealedTop(me.playerId, willReveal);
		toast.info(willReveal ? 'Top card revealed permanently' : 'Top card hidden');
	}

	function showNumberInput(
		title: string,
		defaultValue: number,
		onConfirm: (value: number) => void
	): void {
		callbacks?.setNumberInputDialog({
			show: true,
			title,
			defaultValue,
			onConfirm: (value) => {
				onConfirm(value);
				callbacks?.setNumberInputDialog({ show: false });
			}
		});
	}

	/**
	 * Create deck context menu actions
	 */
	function createDeckContextMenuActions(onSearchLibrary: () => void): MenuAction[] {
		const me = getLocalPlayer();
		if (!me) return [];

		return [
			{
				label: 'Draw Cards',
				submenu: [
					{ label: '1 Card', onClick: () => handleDrawN(1) },
					{ label: '2 Cards', onClick: () => handleDrawN(2) },
					{ label: '3 Cards', onClick: () => handleDrawN(3) },
					{ label: '7 Cards', onClick: () => handleDrawN(7) },
					{
						label: 'Custom...',
						onClick: () => showNumberInput('Draw N Cards', 1, handleDrawN)
					}
				]
			},
			{
				label: 'Scry',
				submenu: [
					{ label: '1 Card', onClick: () => handleScry(1) },
					{ label: '2 Cards', onClick: () => handleScry(2) },
					{ label: '3 Cards', onClick: () => handleScry(3) },
					{
						label: 'Custom...',
						onClick: () => showNumberInput('Scry N Cards', 1, handleScry)
					}
				]
			},
			{
				label: 'Mill Cards',
				submenu: [
					{ label: '1 Card', onClick: () => handleMill(1) },
					{ label: '2 Cards', onClick: () => handleMill(2) },
					{ label: '3 Cards', onClick: () => handleMill(3) },
					{ label: '5 Cards', onClick: () => handleMill(5) },
					{
						label: 'Custom...',
						onClick: () => showNumberInput('Mill N Cards', 1, handleMill)
					}
				]
			},
			{ divider: true },
			{
				label: 'Reveal Top Card',
				onClick: () => handleRevealTop(1)
			},
			{
				label: me.revealedTopCard ? 'Hide Revealed Top' : 'Reveal Top Permanently',
				onClick: handleToggleRevealedTop
			},
			{ divider: true },
			{
				label: 'Search Library',
				onClick: onSearchLibrary
			},
			{
				label: 'Shuffle Library',
				onClick: handleShuffleLibrary
			}
		];
	}

	/**
	 * Handle battlefield card click
	 */
	function handleBattlefieldCardClick(cardId: string): void {
		const card = getBattlefield().find((c) => c.id === cardId);
		if (!card) return;
		gameStore.tapCard(cardId, !card.tapped);
	}

	/**
	 * Handle battlefield drop
	 */
	function handleBattlefieldDrop(cardId: string, sourceZone: string | null): void {
		if (sourceZone === 'hand') {
			gameStore.moveCardToZone(cardId, 'BATTLEFIELD');
		} else if (sourceZone && sourceZone !== 'battlefield') {
			gameStore.moveCardToZone(cardId, 'BATTLEFIELD');
		}
	}

	/**
	 * Handle zone drop (graveyard, exile, hand)
	 */
	function handleZoneDrop(cardId: string, zone: string): void {
		gameStore.moveCardToZone(cardId, zone);
	}

	/**
	 * Handle keyboard shortcuts for hovered card
	 */
	function handleHoveredCardShortcut(key: string, hoveredCardId: string | null): boolean {
		if (!hoveredCardId) return false;

		switch (key) {
			case 'd':
				gameStore.moveCardToZone(hoveredCardId, 'GRAVEYARD');
				return true;
			case 's':
				gameStore.moveCardToZone(hoveredCardId, 'EXILE');
				return true;
			case 'r':
				gameStore.moveCardToZone(hoveredCardId, 'HAND');
				return true;
			case 't':
				gameStore.moveCardToZone(hoveredCardId, 'LIBRARY');
				return true;
			default:
				return false;
		}
	}

	return {
		handleLifeChange,
		handlePoisonChange,
		handleDrawCard,
		handleShuffleLibrary,
		handleUntapAll,
		handleNextTurn,
		handleMulligan,
		handleKeepHand,
		handleDeckContextMenu,
		handleDrawN,
		handleMill,
		handleScry,
		handleScryComplete,
		handleRevealTop,
		handleToggleRevealedTop,
		showNumberInput,
		createDeckContextMenuActions,
		handleBattlefieldCardClick,
		handleBattlefieldDrop,
		handleZoneDrop,
		handleHoveredCardShortcut
	};
}
