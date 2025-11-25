/**
 * Lobby Updates Service
 * Handles real-time table updates via WebSocket
 */

import { websocketStore } from '$lib/stores/websocket';
import { CallbackMethod } from '$lib/generated/mage/v1/websocket';
import type { TableView } from '$lib/generated/mage/v1/models';
import type { Table } from '$lib/types/table';

/**
 * Table update event types
 */
export type TableUpdateType = 'created' | 'updated' | 'deleted';

/**
 * Table update event
 */
export interface TableUpdateEvent {
	type: TableUpdateType;
	table: Table;
}

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
	const getCreateTime = (): number => {
		if (!view.createTime) return Date.now();
		if (view.createTime instanceof Date) return view.createTime.getTime();
		if (typeof view.createTime === 'string') return new Date(view.createTime).getTime();
		return Date.now();
	};

	const createTime = getCreateTime();

	// Convert seats to players
	const players = view.seats
		.filter((seat) => seat.playerName)
		.map((seat, index) => ({
			id: `${view.tableId}-${seat.seatNumber}`,
			username: seat.playerName,
			isHost: index === 0,
			joinedAt: createTime
		}));

	return {
		id: view.tableId,
		name: view.tableName || view.matchOptions?.name || 'Unnamed Table',
		format: (view.matchOptions?.gameType || view.gameType || 'Unknown') as Table['format'],
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
 * Subscribe to lobby table updates
 * Returns unsubscribe function
 */
export function subscribeLobbyUpdates(onUpdate: (event: TableUpdateEvent) => void): () => void {
	const unsubscribers: (() => void)[] = [];

	// Subscribe to TABLE_WAITING (table created/updated)
	unsubscribers.push(
		websocketStore.on(CallbackMethod.TABLE_WAITING, (data) => {
			try {
				// Parse the table data
				// The server sends TableView in the data field
				const tableView = data as TableView;
				const table = convertTableViewToTable(tableView);

				// Determine if this is a new table or an update
				// You might want to track table IDs to differentiate
				// For now, we'll treat all as updates
				onUpdate({
					type: 'updated',
					table
				});
			} catch (err) {
				console.error('[Lobby Updates] Failed to process TABLE_WAITING:', err);
			}
		})
	);

	// Subscribe to JOINED_TABLE (player joined a table)
	unsubscribers.push(
		websocketStore.on(CallbackMethod.JOINED_TABLE, (data) => {
			try {
				// When a player joins, we need to refresh that table's data
				// The data contains tableId and roomId
				const joinData = data as { tableId: string; roomId: string };

				// Trigger a refresh request or mark the table as needing update
				// For simplicity, we'll just log it - the TABLE_WAITING event should follow
				console.log('[Lobby Updates] Player joined table:', joinData.tableId);
			} catch (err) {
				console.error('[Lobby Updates] Failed to process JOINED_TABLE:', err);
			}
		})
	);

	// Subscribe to chat messages (for lobby chat)
	unsubscribers.push(
		websocketStore.on(CallbackMethod.CHATMESSAGE, (data) => {
			try {
				// Chat messages are handled by the chat component
				// We just log them here for debugging
				console.log('[Lobby Updates] Chat message received');
			} catch (err) {
				console.error('[Lobby Updates] Failed to process CHATMESSAGE:', err);
			}
		})
	);

	// Subscribe to server messages
	unsubscribers.push(
		websocketStore.on(CallbackMethod.SERVER_MESSAGE, (data) => {
			try {
				const serverMsg = data as { message: string; isError: boolean };
				console.log(
					`[Lobby Updates] Server message: ${serverMsg.message}`,
					serverMsg.isError ? '(error)' : ''
				);
			} catch (err) {
				console.error('[Lobby Updates] Failed to process SERVER_MESSAGE:', err);
			}
		})
	);

	// Return combined unsubscribe function
	return () => {
		unsubscribers.forEach((unsub) => unsub());
	};
}

/**
 * Connect to WebSocket with session ID
 */
export async function connectLobbyUpdates(sessionId: string): Promise<void> {
	await websocketStore.connect(sessionId);
}

/**
 * Disconnect from WebSocket
 */
export function disconnectLobbyUpdates(): void {
	websocketStore.disconnect();
}

/**
 * Check if WebSocket is connected
 */
export function isLobbyUpdatesConnected(): boolean {
	return websocketStore.isConnected();
}
