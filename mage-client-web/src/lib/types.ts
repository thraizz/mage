// Game types matching the Go server

export interface Card {
	id: string;
	name: string;
	type: string;
	power?: string;
	toughness?: string;
	zone: string;
	tapped: boolean;
	attacking: boolean;
	blocking: boolean;
	damage: number;
	controller: string;
	owner: string;
	abilities: Ability[];
}

export interface Ability {
	id: string;
	text: string;
}

export interface Player {
	id: string;
	name: string;
	life: number;
	library_count: number;
	hand_count: number;
	graveyard_count: number;
}

export interface GameState {
	game_id: string;
	current_player: string;
	active_player: string;
	priority_player: string;
	phase: string;
	step: string;
	turn: number;
	players: Player[];
	battlefield: Card[];
	hand: Card[];
	graveyard: Card[];
	exile: Card[];
	stack: any[];
	combat?: CombatState;
}

export interface CombatState {
	attacking_player: string;
	defenders: string[];
	attackers: string[];
	blockers: string[];
}

export interface WSMessage {
	type: string;
	game_id?: string;
	player_id?: string;
	data?: any;
}
