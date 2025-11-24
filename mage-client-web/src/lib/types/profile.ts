/**
 * User profile and statistics types
 */

export interface UserProfile {
	id: string;
	username: string;
	email: string;
	createdAt: number; // timestamp
	lastLogin?: number; // timestamp
}

export interface UserStats {
	gamesPlayed: number;
	wins: number;
	losses: number;
	draws: number;
	winRate: number; // percentage (0-100)
	quitRate: number; // percentage (0-100)
	totalPlayTime?: number; // seconds
}

export interface MatchHistory {
	id: string;
	opponent: string;
	format: string;
	result: 'win' | 'loss' | 'draw';
	timestamp: number;
	duration?: number; // seconds
}

export interface ChangePasswordRequest {
	currentPassword: string;
	newPassword: string;
}
