import { getMageClient } from '$lib/grpc/client';
import type {
	AdminGetAllActiveGamesRequest,
	AdminGetAllActiveGamesResponse,
	AdminGetServerDebugStateRequest,
	AdminGetServerDebugStateResponse,
	DebugActiveGameInfo
} from '$lib/generated/mage/v1/admin';

export interface ActiveGameDebug {
	gameId: string;
	tableId: string;
	gameType: string;
	players: string[];
	turnNumber: number;
	state: string;
	createdAt: string;
	updatedAt: string;
	inMemory: boolean;
	inDatabase: boolean;
}

export interface ServerDebugState {
	memoryActiveGames: number;
	memoryActiveTables: number;
	memoryActiveSessions: number;
	dbActiveGames: number;
	dbMatchHistory: number;
	gamesInMemoryOnly: string[];
	gamesInDbOnly: string[];
	serverUptime: string;
	lastGameSave: string;
}

export interface AllActiveGamesResult {
	games: ActiveGameDebug[];
	totalInMemory: number;
	totalInDatabase: number;
}

/**
 * Fetch all active games from both server memory and database
 * This is an admin/debug endpoint
 */
export async function fetchAllActiveGames(): Promise<AllActiveGamesResult> {
	const client = getMageClient();
	const sessionId = await client.ensureSessionId();

	if (!sessionId) {
		throw new Error('No active session - please login first');
	}

	const request: AdminGetAllActiveGamesRequest = {
		sessionId
	};

	const response = await client.call<AdminGetAllActiveGamesRequest, AdminGetAllActiveGamesResponse>(
		'AdminGetAllActiveGames',
		request
	);

	return {
		games: (response.games || []).map((g: DebugActiveGameInfo) => ({
			gameId: g.gameId,
			tableId: g.tableId,
			gameType: g.gameType,
			players: g.players || [],
			turnNumber: g.turnNumber,
			state: g.state,
			createdAt: g.createdAt,
			updatedAt: g.updatedAt,
			inMemory: g.inMemory,
			inDatabase: g.inDatabase
		})),
		totalInMemory: response.totalInMemory,
		totalInDatabase: response.totalInDatabase
	};
}

/**
 * Fetch server debug state including memory vs database comparison
 * This is an admin/debug endpoint
 */
export async function fetchServerDebugState(): Promise<ServerDebugState> {
	const client = getMageClient();
	const sessionId = await client.ensureSessionId();

	if (!sessionId) {
		throw new Error('No active session - please login first');
	}

	const request: AdminGetServerDebugStateRequest = {
		sessionId
	};

	const response = await client.call<
		AdminGetServerDebugStateRequest,
		AdminGetServerDebugStateResponse
	>('AdminGetServerDebugState', request);

	return {
		memoryActiveGames: response.memoryActiveGames,
		memoryActiveTables: response.memoryActiveTables,
		memoryActiveSessions: response.memoryActiveSessions,
		dbActiveGames: response.dbActiveGames,
		dbMatchHistory: response.dbMatchHistory,
		gamesInMemoryOnly: response.gamesInMemoryOnly || [],
		gamesInDbOnly: response.gamesInDbOnly || [],
		serverUptime: response.serverUptime,
		lastGameSave: response.lastGameSave
	};
}
