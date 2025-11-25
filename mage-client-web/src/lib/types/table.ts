/**
 * Table game format types
 */
export type GameFormat =
	| 'Standard'
	| 'Modern'
	| 'Legacy'
	| 'Vintage'
	| 'Commander'
	| 'Pauper'
	| 'Pioneer'
	| 'Historic'
	| 'Alchemy'
	| 'Brawl'
	| 'Limited'
	| 'Draft'
	| 'Sealed';

/**
 * Table status
 */
export type TableStatus = 'waiting' | 'ready' | 'playing' | 'finished';

/**
 * Player in a table
 */
export interface TablePlayer {
	id: string;
	username: string;
	isHost: boolean;
	joinedAt: number;
}

/**
 * Table information
 */
export interface Table {
	id: string;
	name: string;
	format: GameFormat;
	hostUsername: string;
	players: TablePlayer[];
	maxPlayers: number;
	status: TableStatus;
	hasPassword: boolean;
	createdAt: number;
	startedAt?: number;
}

/**
 * Table creation request
 */
export interface CreateTableRequest {
	name?: string;
	format: GameFormat;
	maxPlayers: number;
	password?: string;
}

/**
 * Table filter options
 */
export interface TableFilters {
	format?: GameFormat;
	openOnly?: boolean;
	searchQuery?: string;
}
