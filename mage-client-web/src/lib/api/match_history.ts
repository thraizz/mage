import { getMageClient } from '$lib/grpc/client';
import type {
  GetMatchHistoryRequest,
  GetMatchHistoryResponse,
  GetMatchByIdRequest,
  GetMatchByIdResponse,
  MatchHistoryEntry,
  MatchPlayerInfo
} from '$lib/generated/mage/v1/room';

/**
 * Match history type for use in components
 */
export interface Match {
  id: number;
  gameId: string;
  tableId: string;
  tournamentId: string;
  players: Player[];
  gameType: string;
  startTime: Date;
  endTime: Date;
  durationSeconds: number;
  winnerId: number;
  winnerName: string;
}

/**
 * Player info in a match
 */
export interface Player {
  userId: number;
  username: string;
  deck: string;
  result: 'win' | 'loss' | 'draw' | 'concede';
}

/**
 * Match history list with pagination
 */
export interface MatchHistoryList {
  matches: Match[];
  totalCount: number;
}

/**
 * Match details including optional replay data
 */
export interface MatchDetails {
  match: Match;
  replayData?: string;
}

/**
 * Convert proto MatchPlayerInfo to our Player type
 */
function convertPlayerInfo(playerInfo: MatchPlayerInfo): Player {
  return {
    userId: playerInfo.userId,
    username: playerInfo.username,
    deck: playerInfo.deck,
    result: playerInfo.result as 'win' | 'loss' | 'draw' | 'concede'
  };
}

/**
 * Convert proto MatchHistoryEntry to our Match type
 */
function convertMatchEntry(entry: MatchHistoryEntry): Match {
  return {
    id: entry.id,
    gameId: entry.gameId,
    tableId: entry.tableId,
    tournamentId: entry.tournamentId,
    players: entry.players.map(convertPlayerInfo),
    gameType: entry.gameType,
    startTime: entry.startTime || new Date(),
    endTime: entry.endTime || new Date(),
    durationSeconds: entry.durationSeconds,
    winnerId: entry.winnerId,
    winnerName: entry.winnerName
  };
}

/**
 * Fetch user's match history with pagination
 */
export async function fetchMatchHistory(
  limit: number = 50,
  offset: number = 0
): Promise<MatchHistoryList> {
  const client = getMageClient();
  const sessionId = await client.ensureSessionId();

  if (!sessionId) {
    throw new Error('No active session - please login first');
  }

  // Clamp limit to max 100
  const clampedLimit = Math.min(Math.max(limit, 1), 100);

  const request: GetMatchHistoryRequest = {
    sessionId,
    limit: clampedLimit,
    offset: Math.max(offset, 0)
  };

  const response = await client.call<GetMatchHistoryRequest, GetMatchHistoryResponse>(
    'GetMatchHistory',
    request
  );

  if (!response.success) {
    throw new Error(response.error || 'Failed to fetch match history');
  }

  return {
    matches: response.matches.map(convertMatchEntry),
    totalCount: response.totalCount
  };
}

/**
 * Get full details of a specific match including replay data
 */
export async function getMatchDetails(matchId: number): Promise<MatchDetails> {
  const client = getMageClient();
  const sessionId = await client.ensureSessionId();

  if (!sessionId) {
    throw new Error('No active session - please login first');
  }

  const request: GetMatchByIdRequest = {
    sessionId,
    matchId
  };

  const response = await client.call<GetMatchByIdRequest, GetMatchByIdResponse>(
    'GetMatchById',
    request
  );

  if (!response.success) {
    throw new Error(response.error || 'Failed to fetch match details');
  }

  if (!response.match) {
    throw new Error('Match not found');
  }

  return {
    match: convertMatchEntry(response.match),
    replayData: response.replayData || undefined
  };
}

/**
 * Format duration in seconds to human-readable string
 */
export function formatDuration(seconds: number): string {
  const hours = Math.floor(seconds / 3600);
  const minutes = Math.floor((seconds % 3600) / 60);
  const secs = seconds % 60;

  if (hours > 0) {
    return `${hours}h ${minutes}m`;
  } else if (minutes > 0) {
    return `${minutes}m ${secs}s`;
  } else {
    return `${secs}s`;
  }
}

/**
 * Format date to relative time (e.g., "2 hours ago", "3 days ago")
 */
export function formatRelativeTime(date: Date): string {
  const now = new Date();
  const diffMs = now.getTime() - date.getTime();
  const diffSecs = Math.floor(diffMs / 1000);
  const diffMins = Math.floor(diffSecs / 60);
  const diffHours = Math.floor(diffMins / 60);
  const diffDays = Math.floor(diffHours / 24);

  if (diffDays > 30) {
    return date.toLocaleDateString();
  } else if (diffDays > 0) {
    return `${diffDays} day${diffDays > 1 ? 's' : ''} ago`;
  } else if (diffHours > 0) {
    return `${diffHours} hour${diffHours > 1 ? 's' : ''} ago`;
  } else if (diffMins > 0) {
    return `${diffMins} minute${diffMins > 1 ? 's' : ''} ago`;
  } else {
    return 'just now';
  }
}
