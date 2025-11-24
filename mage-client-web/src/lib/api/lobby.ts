import type { Table, CreateTableRequest, GameFormat } from '$lib/types/table';

/**
 * Mock table data for development
 */
const MOCK_TABLES: Table[] = [
	{
		id: '1',
		name: 'Commander Night',
		format: 'Commander',
		hostUsername: 'alice',
		players: [
			{ id: 'p1', username: 'alice', isHost: true, isReady: true, joinedAt: Date.now() - 300000 },
			{ id: 'p2', username: 'bob', isHost: false, isReady: true, joinedAt: Date.now() - 240000 },
			{
				id: 'p3',
				username: 'charlie',
				isHost: false,
				isReady: false,
				joinedAt: Date.now() - 120000
			}
		],
		maxPlayers: 4,
		status: 'waiting',
		hasPassword: false,
		createdAt: Date.now() - 300000
	},
	{
		id: '2',
		name: 'Standard Ranked',
		format: 'Standard',
		hostUsername: 'dave',
		players: [
			{ id: 'p4', username: 'dave', isHost: true, isReady: true, joinedAt: Date.now() - 180000 },
			{ id: 'p5', username: 'eve', isHost: false, isReady: true, joinedAt: Date.now() - 90000 }
		],
		maxPlayers: 2,
		status: 'ready',
		hasPassword: false,
		createdAt: Date.now() - 180000
	},
	{
		id: '3',
		name: 'Modern Tournament',
		format: 'Modern',
		hostUsername: 'frank',
		players: [
			{ id: 'p6', username: 'frank', isHost: true, isReady: true, joinedAt: Date.now() - 600000 },
			{ id: 'p7', username: 'grace', isHost: false, isReady: true, joinedAt: Date.now() - 540000 }
		],
		maxPlayers: 2,
		status: 'playing',
		hasPassword: true,
		createdAt: Date.now() - 600000,
		startedAt: Date.now() - 60000
	},
	{
		id: '4',
		name: 'Casual Legacy',
		format: 'Legacy',
		hostUsername: 'henry',
		players: [
			{ id: 'p8', username: 'henry', isHost: true, isReady: false, joinedAt: Date.now() - 60000 }
		],
		maxPlayers: 2,
		status: 'waiting',
		hasPassword: false,
		createdAt: Date.now() - 60000
	},
	{
		id: '5',
		name: 'Pioneer League',
		format: 'Pioneer',
		hostUsername: 'iris',
		players: [
			{ id: 'p9', username: 'iris', isHost: true, isReady: true, joinedAt: Date.now() - 420000 }
		],
		maxPlayers: 4,
		status: 'waiting',
		hasPassword: false,
		createdAt: Date.now() - 420000
	},
	{
		id: '6',
		name: 'Pauper Fun',
		format: 'Pauper',
		hostUsername: 'jack',
		players: [
			{ id: 'p10', username: 'jack', isHost: true, isReady: true, joinedAt: Date.now() - 150000 },
			{ id: 'p11', username: 'kate', isHost: false, isReady: true, joinedAt: Date.now() - 100000 },
			{ id: 'p12', username: 'liam', isHost: false, isReady: true, joinedAt: Date.now() - 50000 },
			{ id: 'p13', username: 'mia', isHost: false, isReady: false, joinedAt: Date.now() - 25000 }
		],
		maxPlayers: 4,
		status: 'waiting',
		hasPassword: false,
		createdAt: Date.now() - 150000
	}
];

/**
 * Fetch all tables from the lobby
 */
export async function fetchTables(): Promise<Table[]> {
	// Simulate network delay
	await new Promise((resolve) => setTimeout(resolve, 500));

	// In production, this would be:
	// const response = await grpcCall(lobbyService.listTables, {}, 'LobbyService.listTables');
	// return response.tables;

	return MOCK_TABLES;
}

/**
 * Create a new table
 */
export async function createTable(request: CreateTableRequest): Promise<Table> {
	// Simulate network delay
	await new Promise((resolve) => setTimeout(resolve, 800));

	// In production, this would be:
	// const response = await grpcCall(lobbyService.createTable, request, 'LobbyService.createTable');
	// return response.table;

	const newTable: Table = {
		id: `table-${Date.now()}`,
		name: request.name || `${request.format} Game`,
		format: request.format,
		hostUsername: 'currentuser', // This would come from auth store
		players: [
			{
				id: 'current',
				username: 'currentuser',
				isHost: true,
				isReady: false,
				joinedAt: Date.now()
			}
		],
		maxPlayers: request.maxPlayers,
		status: 'waiting',
		hasPassword: !!request.password,
		createdAt: Date.now()
	};

	return newTable;
}

/**
 * Join an existing table
 */
export async function joinTable(tableId: string, password?: string): Promise<void> {
	// Simulate network delay
	await new Promise((resolve) => setTimeout(resolve, 600));

	// In production, this would be:
	// await grpcCall(lobbyService.joinTable, { tableId, password }, 'LobbyService.joinTable');

	console.log(`Joining table ${tableId}`, { password });
}

/**
 * Leave a table
 */
export async function leaveTable(tableId: string): Promise<void> {
	// Simulate network delay
	await new Promise((resolve) => setTimeout(resolve, 400));

	// In production, this would be:
	// await grpcCall(lobbyService.leaveTable, { tableId }, 'LobbyService.leaveTable');

	console.log(`Leaving table ${tableId}`);
}

/**
 * Get all available game formats
 */
export function getGameFormats(): GameFormat[] {
	return [
		'Standard',
		'Modern',
		'Legacy',
		'Vintage',
		'Commander',
		'Pauper',
		'Pioneer',
		'Historic',
		'Alchemy',
		'Brawl',
		'Limited',
		'Draft',
		'Sealed'
	];
}
