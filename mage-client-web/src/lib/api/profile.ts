import type {
	UserProfile,
	UserStats,
	MatchHistory,
	ChangePasswordRequest
} from '$lib/types/profile';
import { getMageClient } from '$lib/grpc/client';
import type {
	UserGetProfileRequest,
	UserGetProfileResponse,
	UserGetStatsRequest,
	UserGetStatsResponse,
	UserGetMatchHistoryRequest,
	UserGetMatchHistoryResponse,
	UserChangePasswordRequest,
	UserChangePasswordResponse
} from '$lib/generated/mage/v1/user';

/**
 * Fetch user profile information
 */
export async function fetchUserProfile(): Promise<UserProfile> {
	const client = getMageClient();
	const sessionId = await client.ensureSessionId();

	if (!sessionId) {
		throw new Error('No active session - please login first');
	}

	const request: UserGetProfileRequest = {
		sessionId
	};

	const response = await client.call<UserGetProfileRequest, UserGetProfileResponse>(
		'UserGetProfile',
		request
	);

	if (!response.success) {
		throw new Error(response.error || 'Failed to fetch user profile');
	}

	if (!response.profile) {
		throw new Error('No profile data returned');
	}

	return {
		id: response.profile.userId,
		username: response.profile.username,
		email: response.profile.email || '',
		createdAt: response.profile.createdAt * 1000, // Convert seconds to milliseconds
		lastLogin: response.profile.lastLogin ? response.profile.lastLogin * 1000 : undefined
	};
}

/**
 * Fetch user statistics
 */
export async function fetchUserStats(): Promise<UserStats> {
	const client = getMageClient();
	const sessionId = await client.ensureSessionId();

	if (!sessionId) {
		throw new Error('No active session - please login first');
	}

	const request: UserGetStatsRequest = {
		sessionId
	};

	const response = await client.call<UserGetStatsRequest, UserGetStatsResponse>(
		'UserGetStats',
		request
	);

	if (!response.success) {
		throw new Error(response.error || 'Failed to fetch user stats');
	}

	if (!response.stats) {
		throw new Error('No stats data returned');
	}

	const stats = response.stats;
	const totalGames = stats.wins + stats.losses + stats.draws;
	const winRate = totalGames > 0 ? (stats.wins / totalGames) * 100 : 0;
	const quitRate = totalGames > 0 ? (stats.quits / totalGames) * 100 : 0;

	return {
		gamesPlayed: totalGames,
		wins: stats.wins,
		losses: stats.losses,
		draws: stats.draws,
		winRate: Math.round(winRate * 100) / 100, // Round to 2 decimal places
		quitRate: Math.round(quitRate * 100) / 100,
		totalPlayTime: stats.totalPlayTime || undefined
	};
}

/**
 * Fetch user match history
 */
export async function fetchMatchHistory(limit: number = 20): Promise<MatchHistory[]> {
	const client = getMageClient();
	const sessionId = await client.ensureSessionId();

	if (!sessionId) {
		throw new Error('No active session - please login first');
	}

	const request: UserGetMatchHistoryRequest = {
		sessionId,
		limit
	};

	const response = await client.call<UserGetMatchHistoryRequest, UserGetMatchHistoryResponse>(
		'UserGetMatchHistory',
		request
	);

	if (!response.success) {
		throw new Error(response.error || 'Failed to fetch match history');
	}

	if (!response.matches) {
		return [];
	}

	return response.matches.map((match) => ({
		id: match.matchId,
		opponent: match.opponentName || 'Unknown',
		format: match.format || 'Unknown',
		result: match.result as 'win' | 'loss' | 'draw',
		timestamp: match.timestamp * 1000, // Convert seconds to milliseconds
		duration: match.duration || undefined
	}));
}

/**
 * Change user password
 */
export async function changePassword(request: ChangePasswordRequest): Promise<void> {
	const client = getMageClient();
	const sessionId = await client.ensureSessionId();

	if (!sessionId) {
		throw new Error('No active session - please login first');
	}

	const changePasswordRequest: UserChangePasswordRequest = {
		sessionId,
		currentPassword: request.currentPassword,
		newPassword: request.newPassword
	};

	const response = await client.call<UserChangePasswordRequest, UserChangePasswordResponse>(
		'UserChangePassword',
		changePasswordRequest
	);

	if (!response.success) {
		throw new Error(response.error || 'Failed to change password');
	}
}
