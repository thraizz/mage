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
			...headers
		},
		body: JSON.stringify(request)
	});

	if (!response.ok) {
		const errorText = await response.text();
		throw new Error(`RPC ${method} failed: ${response.statusText} - ${errorText}`);
	}

	return await response.json();
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
	}

	/**
	 * Set session ID for authenticated requests
	 */
	setSessionId(sessionId: string) {
		this.sessionId = sessionId;
	}

	/**
	 * Get current session ID
	 */
	getSessionId(): string | null {
		return this.sessionId;
	}

	/**
	 * Clear session ID (logout)
	 */
	clearSession() {
		this.sessionId = null;
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

		if (response.success && response.sessionId) {
			this.setSessionId(response.sessionId);
		}

		return response;
	}

	/**
	 * Register a new user
	 */
	async register(userName: string, password: string, email: string): Promise<AuthRegisterResponse> {
		const request: AuthRegisterRequest = {
			userName,
			password,
			email
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
