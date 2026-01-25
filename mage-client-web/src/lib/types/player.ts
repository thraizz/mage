/**
 * Player type definitions for online player tracking
 */

export interface OnlinePlayer {
  id: string;
  username: string;
  isCurrentUser: boolean;
  joinedAt: number;
}

export interface PlayerListState {
  players: OnlinePlayer[];
  count: number;
  isLoading: boolean;
  error: string | null;
}
