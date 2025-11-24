import type { Table, CreateTableRequest, GameFormat } from '$lib/types/table';
import type { OnlinePlayer } from '$lib/types/player';
import { getMageClient } from '$lib/grpc/client';
import type { TableView } from '$lib/generated/mage/v1/models';
import type {
	RoomCreateTableRequest,
	RoomCreateTableResponse,
	RoomJoinTableRequest,
	RoomJoinTableResponse,
	RoomLeaveTableOrTournamentRequest,
	RoomLeaveTableOrTournamentResponse
} from '$lib/generated/mage/v1/table';
import type { RoomGetUsersResponse } from '$lib/generated/mage/v1/room';

/**
 * Convert TableView from proto to our Table type
 */
function convertTableViewToTable(view: TableView): Table {
	// Parse table status from tableStateText
	let status: 'waiting' | 'ready' | 'playing' | 'finished' = 'waiting';
	const stateText = view.tableStateText.toLowerCase();
	if (stateText.includes('playing') || stateText.includes('started')) {
		status = 'playing';
	} else if (stateText.includes('ready') || stateText.includes('starting')) {
		status = 'ready';
	} else if (stateText.includes('finish') || stateText.includes('complete')) {
		status = 'finished';
	}

	// Helper to convert createTime to timestamp
	// Handles both Date objects and ISO string timestamps
	const getCreateTime = (): number => {
		if (!view.createTime) {
			return Date.now();
		}
		if (view.createTime instanceof Date) {
			return view.createTime.getTime();
		}
		if (typeof view.createTime === 'string') {
			return new Date(view.createTime).getTime();
		}
		return Date.now();
	};

	const createTime = getCreateTime();

	// Convert seats to players
	const players = view.seats
		.filter((seat) => seat.playerName) // Only include occupied seats
		.map((seat, index) => ({
			id: `${view.tableId}-${seat.seatNumber}`,
			username: seat.playerName,
			isHost: index === 0, // First player is typically the host
			isReady: status !== 'waiting', // Assume ready if game is not waiting
			joinedAt: createTime
		}));

	return {
		id: view.tableId,
		name: view.tableName || view.matchOptions?.name || 'Unnamed Table',
		format: (view.matchOptions?.gameType || view.gameType || 'Unknown') as GameFormat,
		hostUsername: view.controllerName,
		players,
		maxPlayers: view.numSeats,
		status,
		hasPassword: !!view.password,
		createdAt: createTime,
		startedAt: status === 'playing' ? Date.now() : undefined
	};
}

/**
 * Fetch all tables from the lobby
 */
export async function fetchTables(): Promise<Table[]> {
	const client = getMageClient();

	// Ensure sessionId is available (will restore from token if needed)
	const sessionId = await client.ensureSessionId();
	if (!sessionId) {
		throw new Error('No active session - please login first');
	}

	// Get main room ID first
	const roomResponse = await client.getMainRoomId();
	if (!roomResponse.roomId) {
		throw new Error('Failed to get main room ID');
	}

	// Get all tables in the room
	const response = await client.getAllTables(roomResponse.roomId);

	// Convert TableView[] to Table[]
	return response.tables.map(convertTableViewToTable);
}

/**
 * Create a new table
 */
export async function createTable(request: CreateTableRequest): Promise<Table> {
	const client = getMageClient();
	const sessionId = await client.ensureSessionId();

	if (!sessionId) {
		throw new Error('No active session - please login first');
	}

	// Get main room ID
	const roomResponse = await client.getMainRoomId();
	if (!roomResponse.roomId) {
		throw new Error('Failed to get main room ID');
	}

	// Create match options
	const matchOptions = {
		name: request.name || `${request.format} Game`,
		gameType: request.format,
		deckType: 'Constructed',
		limited: false,
		winsNeeded: 2,
		freeMulligans: 0,
		priorityTime: 0,
		rated: false,
		banlist: [],
		skillLevel: 'Casual',
		rangeOfInfluence: false,
		planeChase: false,
		rollbackTurnsAllowed: true,
		embedDeckInSavedGame: 0
	};

	const createRequest: RoomCreateTableRequest = {
		sessionId,
		roomId: roomResponse.roomId,
		matchOptions
	};

	const response = await client.call<RoomCreateTableRequest, RoomCreateTableResponse>(
		'RoomCreateTable',
		createRequest
	);

	if (!response.success) {
		throw new Error(response.error || 'Failed to create table');
	}

	// Return a basic table object - the client should refetch to get full details
	return {
		id: response.tableId,
		name: matchOptions.name,
		format: request.format,
		hostUsername: 'You', // Current user
		players: [],
		maxPlayers: request.maxPlayers,
		status: 'waiting',
		hasPassword: !!request.password,
		createdAt: Date.now()
	};
}

/**
 * Join an existing table
 */
export async function joinTable(tableId: string, password?: string): Promise<void> {
	const client = getMageClient();
	const sessionId = await client.ensureSessionId();

	if (!sessionId) {
		throw new Error('No active session - please login first');
	}

	// Get main room ID
	const roomResponse = await client.getMainRoomId();
	if (!roomResponse.roomId) {
		throw new Error('Failed to get main room ID');
	}

	const joinRequest: RoomJoinTableRequest = {
		sessionId,
		roomId: roomResponse.roomId,
		tableId,
		name: 'Player', // Player name - could be from auth store
		playerType: 'Human',
		skill: 1,
		deckType: 'Constructed',
		deck: '', // Deck will be selected later
		password: password || ''
	};

	const response = await client.call<RoomJoinTableRequest, RoomJoinTableResponse>(
		'RoomJoinTable',
		joinRequest
	);

	if (!response.success) {
		throw new Error(response.error || 'Failed to join table');
	}
}

/**
 * Leave a table
 */
export async function leaveTable(tableId: string): Promise<void> {
	const client = getMageClient();
	const sessionId = await client.ensureSessionId();

	if (!sessionId) {
		throw new Error('No active session - please login first');
	}

	// Get main room ID
	const roomResponse = await client.getMainRoomId();
	if (!roomResponse.roomId) {
		throw new Error('Failed to get main room ID');
	}

	const leaveRequest: RoomLeaveTableOrTournamentRequest = {
		sessionId,
		roomId: roomResponse.roomId,
		tableId
	};

	const response = await client.call<
		RoomLeaveTableOrTournamentRequest,
		RoomLeaveTableOrTournamentResponse
	>('RoomLeaveTableOrTournament', leaveRequest);

	if (!response.success) {
		throw new Error(response.error || 'Failed to leave table');
	}
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

/**
 * Fetch online players in the lobby
 */
export async function fetchOnlinePlayers(currentUsername?: string): Promise<OnlinePlayer[]> {
	const client = getMageClient();
	const sessionId = await client.ensureSessionId();

	if (!sessionId) {
		throw new Error('No active session - please login first');
	}

	// Get main room ID
	const roomResponse = await client.getMainRoomId();
	if (!roomResponse.roomId) {
		throw new Error('Failed to get main room ID');
	}

	// Get users in the room
	const response = await client.call<{ sessionId: string; roomId: string }, RoomGetUsersResponse>(
		'RoomGetUsers',
		{
			sessionId,
			roomId: roomResponse.roomId
		}
	);

	// Convert UserView[] to OnlinePlayer[]
	return response.users.map((user) => {
		// Helper to convert connectedAt to timestamp
		// Handles both Date objects and ISO string timestamps
		const getConnectedAt = (): number => {
			if (!user.connectedAt) {
				return Date.now();
			}
			if (user.connectedAt instanceof Date) {
				return user.connectedAt.getTime();
			}
			if (typeof user.connectedAt === 'string') {
				return new Date(user.connectedAt).getTime();
			}
			return Date.now();
		};

		return {
			id: user.userName, // Use username as ID since we don't have user ID
			username: user.userName,
			isCurrentUser: user.userName === currentUsername,
			joinedAt: getConnectedAt()
		};
	});
}
