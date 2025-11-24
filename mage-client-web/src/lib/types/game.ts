/**
 * Game state and card types
 */

export interface GameCard {
	id: string;
	name: string;
	manaCost?: string;
	cardType?: string;
	types?: string[];
	colors?: string[];
	power?: string;
	toughness?: string;
	imageUrl?: string;
	isTapped?: boolean;
	isSelected?: boolean;
	counters?: CardCounter[];
	zone?: CardZone;
	ownerId?: string;
	controllerId?: string;
}

export interface CardCounter {
	type: CounterType;
	count: number;
}

export type CounterType =
	| 'P1P1' // +1/+1 counter
	| 'M1M1' // -1/-1 counter
	| 'LOYALTY' // Planeswalker loyalty
	| 'POISON' // Poison counter
	| 'ENERGY' // Energy counter
	| 'CHARGE' // Charge counter
	| 'TIME' // Time counter
	| 'FADE' // Fade counter
	| 'AGE' // Age counter
	| 'QUEST' // Quest counter
	| 'OTHER'; // Generic counter

export type CardZone =
	| 'HAND'
	| 'BATTLEFIELD'
	| 'GRAVEYARD'
	| 'EXILE'
	| 'LIBRARY'
	| 'STACK'
	| 'COMMAND';

export interface GameState {
	id: string;
	format: string;
	turn: number;
	phase: GamePhase;
	activePlayerId: string;
	priorityPlayerId: string;
	players: GamePlayer[];
	battlefield: GameCard[];
	stack: StackObject[];
}

export interface GamePlayer {
	id: string;
	username: string;
	life: number;
	libraryCount: number;
	handCount: number;
	hand?: GameCard[]; // Only visible for local player
	graveyard: GameCard[];
	exile: GameCard[];
	commandZone: GameCard[];
	manaPool: ManaPool;
}

export interface ManaPool {
	white: number;
	blue: number;
	black: number;
	red: number;
	green: number;
	colorless: number;
}

export type GamePhase =
	| 'BEGINNING'
	| 'UNTAP'
	| 'UPKEEP'
	| 'DRAW'
	| 'PRECOMBAT_MAIN'
	| 'COMBAT'
	| 'DECLARE_ATTACKERS'
	| 'DECLARE_BLOCKERS'
	| 'COMBAT_DAMAGE'
	| 'END_OF_COMBAT'
	| 'POSTCOMBAT_MAIN'
	| 'END'
	| 'END_OF_TURN'
	| 'CLEANUP';

export interface StackObject {
	id: string;
	type: 'SPELL' | 'ABILITY';
	name: string;
	controllerId: string;
	sourceCardId?: string;
	targets?: string[];
}

export interface GameAction {
	type: GameActionType;
	playerId: string;
	cardId?: string;
	targetIds?: string[];
	amount?: number;
	manaPayment?: ManaPool;
}

export type GameActionType =
	| 'PLAY_LAND'
	| 'CAST_SPELL'
	| 'ACTIVATE_ABILITY'
	| 'PASS_PRIORITY'
	| 'ATTACK'
	| 'BLOCK'
	| 'CONCEDE'
	| 'MULLIGAN';

export interface GameEvent {
	type: GameEventType;
	playerId?: string;
	cardId?: string;
	message: string;
	timestamp: number;
}

export type GameEventType =
	| 'GAME_START'
	| 'TURN_START'
	| 'PHASE_CHANGE'
	| 'CARD_DRAWN'
	| 'CARD_PLAYED'
	| 'CARD_TAPPED'
	| 'CARD_UNTAPPED'
	| 'DAMAGE_DEALT'
	| 'LIFE_CHANGED'
	| 'COUNTER_ADDED'
	| 'COUNTER_REMOVED'
	| 'CARD_MOVED'
	| 'PRIORITY_PASSED'
	| 'GAME_END';
