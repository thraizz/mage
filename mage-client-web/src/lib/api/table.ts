/**
 * Table API functions for table lobby operations
 * Real gRPC implementation
 */

import type { Table } from '$lib/types/table';
import { getMageClient } from '$lib/grpc/client';
import type { TableView } from '$lib/generated/mage/v1/models';
import type {
	RoomGetTableByIdRequest,
	RoomGetTableByIdResponse
} from '$lib/generated/mage/v1/room';
import type { GameFormat } from '$lib/types/table';

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
 * Fetch table details by ID
 */
export async function fetchTable(tableId: string): Promise<Table> {
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

	const request: RoomGetTableByIdRequest = {
		sessionId,
		roomId: roomResponse.roomId,
		tableId
	};

	const response = await client.call<RoomGetTableByIdRequest, RoomGetTableByIdResponse>(
		'RoomGetTableById',
		request
	);

	if (!response.table) {
		throw new Error('Table not found');
	}

	return convertTableViewToTable(response.table);
}

/**
 * Toggle ready status for current player
 *
 * NOTE: The TableSetReady RPC does not exist in the current server implementation.
 * In XMage, ready status is controlled by deck submission via DeckSubmit.
 * This function is a placeholder for future implementation.
 */
export async function toggleReady(_tableId: string, _isReady: boolean): Promise<void> {
	// TableSetReady RPC doesn't exist - ready status is controlled by deck submission
	// For now, throw an error to make it clear this functionality isn't implemented
	throw new Error('Ready status is controlled by deck submission. Please submit your deck to mark yourself as ready.');
}

/**
 * Join table with deck
 */
export async function joinTable(
	tableId: string,
	deckList: string,
	password?: string
): Promise<void> {
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

	const request = {
		sessionId,
		roomId: roomResponse.roomId,
		tableId,
		playerName: '', // Server derives from session
		playerType: 'Human',
		skillLevel: 0, // Default skill level
		deckList,
		password: password || ''
	};

	const response = await client.call<typeof request, { success: boolean; error?: string }>(
		'RoomJoinTable',
		request
	);

	if (!response.success) {
		throw new Error(response.error || 'Failed to join table');
	}
}

/**
 * Leave table
 * Uses the same RPC as leaving from lobby
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

	const leaveRequest = {
		sessionId,
		roomId: roomResponse.roomId,
		tableId
	};

	const response = await client.call<typeof leaveRequest, { success: boolean; error?: string }>(
		'RoomLeaveTableOrTournament',
		leaveRequest
	);

	if (!response.success) {
		throw new Error(response.error || 'Failed to leave table');
	}
}

/**
 * Start game (host only)
 *
 * Note: The actual method to start a match might be different
 * This assumes MatchStart exists
 */
export async function startGame(tableId: string): Promise<string> {
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

	try {
		const request = {
			sessionId,
			roomId: roomResponse.roomId,
			tableId
		};

		const response = await client.call<
			typeof request,
			{ success: boolean; error?: string; gameId?: string }
		>('MatchStart', request);

		if (!response.success) {
			throw new Error(response.error || 'Failed to start game');
		}

		// Return game ID if provided, otherwise use table ID
		return response.gameId || tableId;
	} catch (error) {
		console.error('Failed to start game:', error);
		throw error;
	}
}

/**
 * Submit a deck to a table (for deck validation and game start)
 *
 * This allows players to submit/update their deck after joining a table.
 * The deck is sent in structured format with card names and quantities.
 */
export async function submitDeck(
	tableId: string,
	deck: {
		mainDeck: { name: string; quantity: number }[];
		sideboard: { name: string; quantity: number }[];
		commanders: { name: string; quantity: number }[];
	}
): Promise<void> {
	const client = getMageClient();
	const sessionId = await client.ensureSessionId();

	if (!sessionId) {
		throw new Error('No active session - please login first');
	}

	// Convert to the proto DeckCardLists format
	const deckCardLists = {
		mainDeck: deck.mainDeck.map((c) => ({ name: c.name, quantity: c.quantity })),
		sideboard: deck.sideboard.map((c) => ({ name: c.name, quantity: c.quantity })),
		commanders: deck.commanders.map((c) => ({ name: c.name, quantity: c.quantity }))
	};

	const request = {
		sessionId,
		tableId,
		deck: deckCardLists
	};

	const response = await client.call<typeof request, { success: boolean; error?: string }>(
		'DeckSubmit',
		request
	);

	if (!response.success) {
		throw new Error(response.error || 'Failed to submit deck');
	}
}

/**
 * Kick player from table (host only)
 *
 * Note: This functionality might not be available via RPC
 * It might require admin privileges or be handled differently
 */
export async function kickPlayer(_tableId: string, _playerId: string): Promise<void> {
	const client = getMageClient();
	const sessionId = await client.ensureSessionId();

	if (!sessionId) {
		throw new Error('No active session - please login first');
	}

	// This might require a different RPC method or admin privileges
	// Placeholder implementation
	console.warn('Kick player functionality not yet implemented via RPC');
	throw new Error('Kick player functionality not available');
}
