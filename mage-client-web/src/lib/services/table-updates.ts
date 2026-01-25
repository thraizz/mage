/**
 * Table Updates Service
 * Handles real-time table updates via WebSocket for the table view page
 */

import { websocketStore } from '$lib/stores/websocket';
import { CallbackMethod } from '$lib/generated/mage/v1/websocket';
import type { Table } from '$lib/types/table';

export interface TableUpdateCallback {
  (table: Table): void;
}

/**
 * Subscribe to table updates for a specific table
 * Listens for TABLE_WAITING events and filters by table ID
 */
export function subscribeTableUpdates(tableId: string, callback: TableUpdateCallback): () => void {
  const unsubscribe = websocketStore.on(CallbackMethod.TABLE_WAITING, (data) => {
    try {
      // Parse table data from WebSocket event
      const tableData = data as any;

      // Check if this event is for our table
      if (tableData.tableId === tableId || tableData.id === tableId) {
        // Transform to Table type
        const table: Table = {
          id: tableData.tableId || tableData.id,
          name: tableData.name || tableData.tableName || '',
          hostUsername: tableData.hostUsername || tableData.host || '',
          format: tableData.gameType || tableData.format || 'TwoPlayerDuel',
          maxPlayers: tableData.maxPlayers || 2,
          hasPassword: tableData.hasPassword || false,
          status: mapTableStatus(tableData.state || tableData.status),
          createdAt: tableData.createTime ? new Date(tableData.createTime).getTime() : Date.now(),
          players: (tableData.seats || tableData.players || []).map((seat: any) => ({
            id: seat.playerId || seat.id,
            username: seat.playerName || seat.username || seat.name,
            isHost: seat.isHost || false,
            joinedAt: Date.now()
          }))
        };

        callback(table);
      }
    } catch (err) {
      console.error('[Table Updates] Failed to process TABLE_WAITING:', err);
    }
  });

  return unsubscribe;
}

/**
 * Map server table state to client status
 */
function mapTableStatus(state: string): 'waiting' | 'ready' | 'playing' | 'finished' {
  const lowerState = (state || '').toLowerCase();

  if (lowerState.includes('waiting') || lowerState.includes('created')) {
    return 'waiting';
  }
  if (lowerState.includes('ready') || lowerState.includes('starting')) {
    return 'ready';
  }
  if (lowerState.includes('playing') || lowerState.includes('dueling')) {
    return 'playing';
  }
  if (lowerState.includes('finished') || lowerState.includes('ended')) {
    return 'finished';
  }

  return 'waiting';
}
