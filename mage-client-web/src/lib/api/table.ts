/**
 * Table API functions for table lobby operations
 * Mock implementation - will be replaced with actual gRPC calls
 */

import type { Table } from '$lib/types/table';

/**
 * Fetch table details by ID
 */
export async function fetchTable(tableId: string): Promise<Table> {
	// Simulate network delay
	await new Promise((resolve) => setTimeout(resolve, 300));

	// Mock table data
	const mockTable: Table = {
		id: tableId,
		name: 'Friday Night Magic',
		format: 'Standard',
		hostUsername: 'MageMaster',
		players: [
			{
				id: '1',
				username: 'MageMaster',
				isHost: true,
				isReady: true,
				joinedAt: Date.now() - 300000
			},
			{
				id: '2',
				username: 'CurrentUser',
				isHost: false,
				isReady: false,
				joinedAt: Date.now() - 120000
			}
		],
		maxPlayers: 4,
		status: 'waiting',
		hasPassword: false,
		createdAt: Date.now() - 300000
	};

	return mockTable;
}

/**
 * Toggle ready status for current player
 */
export async function toggleReady(tableId: string, isReady: boolean): Promise<void> {
	// Simulate network delay
	await new Promise((resolve) => setTimeout(resolve, 200));

	// Mock success
	console.log(`Toggled ready status to ${isReady} for table ${tableId}`);
}

/**
 * Leave table
 */
export async function leaveTable(tableId: string): Promise<void> {
	// Simulate network delay
	await new Promise((resolve) => setTimeout(resolve, 200));

	// Mock success
	console.log(`Left table ${tableId}`);
}

/**
 * Start game (host only)
 */
export async function startGame(tableId: string): Promise<string> {
	// Simulate network delay
	await new Promise((resolve) => setTimeout(resolve, 500));

	// Mock game ID
	const gameId = `game-${Date.now()}`;
	console.log(`Started game ${gameId} for table ${tableId}`);
	return gameId;
}

/**
 * Kick player from table (host only)
 */
export async function kickPlayer(tableId: string, playerId: string): Promise<void> {
	// Simulate network delay
	await new Promise((resolve) => setTimeout(resolve, 300));

	// Mock success
	console.log(`Kicked player ${playerId} from table ${tableId}`);
}
