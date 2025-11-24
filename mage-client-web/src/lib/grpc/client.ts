/**
 * gRPC client for MAGE server
 * Provides type-safe access to all 70+ RPC methods
 * Generated from the same proto files as the Go server
 */

import type {
	AuthRegisterRequest,
	AuthRegisterResponse,
	ConnectUserRequest,
	ConnectUserResponse,
	PingRequest,
	PingResponse,
	GetServerStateRequest,
	GetServerStateResponse
} from '../generated/mage/v1/auth';

import type {
	ServerGetMainRoomIdRequest,
	ServerGetMainRoomIdResponse,
	RoomGetAllTablesRequest,
	RoomGetAllTablesResponse
} from '../generated/mage/v1/room';

import type { GameGetViewRequest, GameGetViewResponse } from '../generated/mage/v1/game';

import type { ChatSendMessageRequest, ChatSendMessageResponse } from '../generated/mage/v1/chat';

/**
 * Configuration for gRPC client connections
 */
export interface GrpcClientConfig {
	serverUrl: string;
	timeout?: number;
	headers?: Record<string, string>;
}

/**
 * Default gRPC client configuration
 */
const defaultConfig: GrpcClientConfig = {
	serverUrl: import.meta.env.VITE_GRPC_SERVER_URL || 'http://localhost:17171',
	timeout: 30000 // 30 seconds
};

/**
 * Storage key for persisting session ID
 */
const SESSION_STORAGE_KEY = 'mage_session_id';

/**
 * Converts a snake_case string to camelCase
 */
function snakeToCamel(str: string): string {
	return str.replace(/_([a-z])/g, (_, letter) => letter.toUpperCase());
}

/**
 * Recursively converts all snake_case keys in an object to camelCase
 * Handles nested objects and arrays
 */
function convertSnakeToCamel<T>(obj: unknown): T {
	if (obj === null || obj === undefined) {
		return obj as T;
	}

	// Handle arrays
	if (Array.isArray(obj)) {
		return obj.map(convertSnakeToCamel) as T;
	}

	// Handle primitive types
	if (typeof obj !== 'object') {
		return obj as T;
	}

	// Handle Date objects
	if (obj instanceof Date) {
		return obj as T;
	}

	// Convert object keys
	const converted: Record<string, unknown> = {};
	for (const [key, value] of Object.entries(obj)) {
		const camelKey = snakeToCamel(key);
		converted[camelKey] = convertSnakeToCamel(value);
	}
	return converted as T;
}

/**
 * Generic RPC method caller
 * Uses fetch() to make HTTP POST requests to gRPC-Web endpoints
 */
async function callRpc<TRequest, TResponse>(
	serverUrl: string,
	method: string,
	request: TRequest,
	headers?: Record<string, string>
): Promise<TResponse> {
	const response = await fetch(`${serverUrl}/mage.v1.MageServer/${method}`, {
		method: 'POST',
		headers: {
			'Content-Type': 'application/json',
			'X-Grpc-Web': '1', // Indicate this is a gRPC-Web request
			...headers
		},
		body: JSON.stringify(request)
	});

	if (!response.ok) {
		const errorText = await response.text();
		throw new Error(`RPC ${method} failed: ${response.statusText} - ${errorText}`);
	}

	const jsonResponse = await response.json();
	
	// Convert snake_case keys to camelCase
	// The server uses UseProtoNames: true which returns snake_case JSON
	// but TypeScript interfaces expect camelCase
	return convertSnakeToCamel<TResponse>(jsonResponse);
}

/**
 * MAGE Server Client
 * Provides type-safe access to all server RPC methods
 */
export class MageClient {
	private config: GrpcClientConfig;
	private sessionId: string | null = null;

	constructor(config?: Partial<GrpcClientConfig>) {
		this.config = { ...defaultConfig, ...config };
		// Restore session ID from localStorage if available
		this.loadSessionFromStorage();
	}

	/**
	 * Set session ID for authenticated requests
	 * Also persists to localStorage
	 */
	setSessionId(sessionId: string) {
		this.sessionId = sessionId;
		if (typeof window !== 'undefined') {
			localStorage.setItem(SESSION_STORAGE_KEY, sessionId);
		}
	}

	/**
	 * Get current session ID
	 */
	getSessionId(): string | null {
		return this.sessionId;
	}

	/**
	 * Ensure session ID is available, trying to restore from token if needed
	 * This is useful when sessionId might not be set yet (e.g., after page refresh)
	 */
	async ensureSessionId(): Promise<string | null> {
		// If we already have a sessionId, return it
		if (this.sessionId) {
			return this.sessionId;
		}

		// Try to restore from localStorage token
		if (typeof window !== 'undefined') {
			try {
				const { getSessionIdFromToken } = await import('$lib/utils/jwt');
				const AUTH_STORAGE_KEY = 'mage_auth_token';
				const token = localStorage.getItem(AUTH_STORAGE_KEY);
				if (token) {
					const sessionId = getSessionIdFromToken(token);
					if (sessionId) {
						this.setSessionId(sessionId);
						return sessionId;
					}
				}
			} catch (error) {
				console.error('Failed to restore sessionId from token:', error);
			}
		}

		return null;
	}

	/**
	 * Clear session ID (logout)
	 * Also removes from localStorage
	 */
	clearSession() {
		this.sessionId = null;
		if (typeof window !== 'undefined') {
			localStorage.removeItem(SESSION_STORAGE_KEY);
		}
	}

	/**
	 * Load session ID from localStorage
	 * Called automatically on client creation
	 * Note: Session ID is also restored from auth token by auth store when auth is loaded
	 */
	private loadSessionFromStorage() {
		if (typeof window !== 'undefined') {
			const storedSessionId = localStorage.getItem(SESSION_STORAGE_KEY);
			if (storedSessionId) {
				this.sessionId = storedSessionId;
				console.log('Session ID restored from localStorage');
			}
		}
	}

	// ==================== Authentication & Connection ====================

	/**
	 * Connect a user (login)
	 */
	async connectUser(
		userName: string,
		password: string,
		clientVersion = '1.0.0'
	): Promise<ConnectUserResponse> {
		const request: ConnectUserRequest = {
			userName,
			password,
			sessionId: '',
			clientVersion,
			userIdStr: '',
			restoreSessionId: ''
		};

		const response = await callRpc<ConnectUserRequest, ConnectUserResponse>(
			this.config.serverUrl,
			'ConnectUser',
			request,
			this.config.headers
		);

		// Debug logging in dev mode
		if (import.meta.env.DEV) {
			console.log('[MageClient] ConnectUser response:', {
				success: response.success,
				sessionId: response.sessionId,
				userId: response.userId,
				error: response.error
			});
		}

		// Set sessionId if we got one (even if it's just an empty string, we'll handle that)
		if (response.success) {
			if (response.sessionId && response.sessionId.trim() !== '') {
				this.setSessionId(response.sessionId);
			} else {
				console.warn('[MageClient] ConnectUser succeeded but sessionId is empty or missing');
			}
		}

		return response;
	}

	/**
	 * Register a new user
	 * Email is optional
	 */
	async register(userName: string, password: string, email?: string): Promise<AuthRegisterResponse> {
		const request: AuthRegisterRequest = {
			userName,
			password,
			email: email || ''
		};

		return await callRpc<AuthRegisterRequest, AuthRegisterResponse>(
			this.config.serverUrl,
			'AuthRegister',
			request,
			this.config.headers
		);
	}

	/**
	 * Keep session alive
	 */
	async ping(): Promise<PingResponse> {
		if (!this.sessionId) {
			throw new Error('No active session');
		}

		const request: PingRequest = {
			sessionId: this.sessionId
		};

		return await callRpc<PingRequest, PingResponse>(
			this.config.serverUrl,
			'Ping',
			request,
			this.config.headers
		);
	}

	/**
	 * Get server state
	 */
	async getServerState(): Promise<GetServerStateResponse> {
		if (!this.sessionId) {
			throw new Error('No active session');
		}

		const request: GetServerStateRequest = {
			sessionId: this.sessionId
		};

		return await callRpc<GetServerStateRequest, GetServerStateResponse>(
			this.config.serverUrl,
			'GetServerState',
			request,
			this.config.headers
		);
	}

	// ==================== Room/Lobby Operations ====================

	/**
	 * Get main room ID
	 */
	async getMainRoomId(): Promise<ServerGetMainRoomIdResponse> {
		if (!this.sessionId) {
			throw new Error('No active session');
		}

		const request: ServerGetMainRoomIdRequest = {
			sessionId: this.sessionId
		};

		return await callRpc<ServerGetMainRoomIdRequest, ServerGetMainRoomIdResponse>(
			this.config.serverUrl,
			'ServerGetMainRoomId',
			request,
			this.config.headers
		);
	}

	/**
	 * Get all tables in a room
	 */
	async getAllTables(roomId: string): Promise<RoomGetAllTablesResponse> {
		if (!this.sessionId) {
			throw new Error('No active session');
		}

		const request: RoomGetAllTablesRequest = {
			sessionId: this.sessionId,
			roomId
		};

		return await callRpc<RoomGetAllTablesRequest, RoomGetAllTablesResponse>(
			this.config.serverUrl,
			'RoomGetAllTables',
			request,
			this.config.headers
		);
	}

	// ==================== Game Operations ====================

	/**
	 * Get game view
	 */
	async getGameView(gameId: string, playerId: string): Promise<GameGetViewResponse> {
		if (!this.sessionId) {
			throw new Error('No active session');
		}

		const request: GameGetViewRequest = {
			sessionId: this.sessionId,
			gameId,
			playerId
		};

		return await callRpc<GameGetViewRequest, GameGetViewResponse>(
			this.config.serverUrl,
			'GameGetView',
			request,
			this.config.headers
		);
	}

	// ==================== Chat Operations ====================

	/**
	 * Send chat message
	 */
	async sendChatMessage(chatId: string, message: string): Promise<ChatSendMessageResponse> {
		if (!this.sessionId) {
			throw new Error('No active session');
		}

		const request: ChatSendMessageRequest = {
			sessionId: this.sessionId,
			chatId,
			message
		};

		return await callRpc<ChatSendMessageRequest, ChatSendMessageResponse>(
			this.config.serverUrl,
			'ChatSendMessage',
			request,
			this.config.headers
		);
	}

	// ==================== Generic RPC Method ====================

	/**
	 * Call any RPC method by name
	 * Useful for methods not wrapped in convenience functions
	 */
	async call<TRequest, TResponse>(method: string, request: TRequest): Promise<TResponse> {
		return await callRpc<TRequest, TResponse>(
			this.config.serverUrl,
			method,
			request,
			this.config.headers
		);
	}
}

/**
 * Singleton instance
 */
let clientInstance: MageClient | null = null;

/**
 * Get or create the singleton MAGE client instance
 */
export function getMageClient(config?: Partial<GrpcClientConfig>): MageClient {
	if (!clientInstance) {
		clientInstance = new MageClient(config);
	}
	return clientInstance;
}

/**
 * Reset the singleton instance
 */
export function resetMageClient() {
	clientInstance = null;
}

/**
 * Example usage:
 *
 * ```typescript
 * import { getMageClient } from '$lib/grpc/client';
 *
 * const client = getMageClient();
 *
 * // Login
 * const response = await client.connectUser('username', 'password');
 * if (response.success) {
 *   console.log('Session ID:', response.sessionId);
 * }
 *
 * // Get server state
 * const state = await client.getServerState();
 * console.log('Active players:', state.serverState?.activePlayers);
 *
 * // Get main lobby
 * const lobby = await client.getMainRoomId();
 * console.log('Lobby ID:', lobby.roomId);
 *
 * // List tables
 * const tables = await client.getAllTables(lobby.roomId);
 * console.log('Tables:', tables.tables);
 *
 * // Send chat message
 * await client.sendChatMessage('lobby', 'Hello everyone!');
 *
 * // Generic RPC call for advanced usage
 * import type { RoomCreateTableRequest, RoomCreateTableResponse } from '$lib/generated/mage/v1/table';
 *
 * const createTableReq: RoomCreateTableRequest = {
 *   sessionId: client.getSessionId()!,
 *   roomId: lobby.roomId,
 *   matchOptions: {
 *     name: 'My Game',
 *     gameType: 'TwoPlayerDuel',
 *     deckType: 'Constructed',
 *     winsNeeded: 2,
 *     // ... other options
 *   }
 * };
 *
 * const createTableRes = await client.call<RoomCreateTableRequest, RoomCreateTableResponse>(
 *   'RoomCreateTable',
 *   createTableReq
 * );
 * ```
 */
