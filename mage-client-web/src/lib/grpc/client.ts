/**
 * Basic gRPC client setup for MAGE multiplayer system
 * This module provides the foundation for connecting to gRPC services
 */

import { credentials, type ChannelCredentials } from '@grpc/grpc-js';
import { GameServiceClient } from '$lib/generated/game';
import { LobbyServiceClient } from '$lib/generated/lobby';

/**
 * Configuration for gRPC client connections
 */
export interface GrpcClientConfig {
	serverUrl: string;
	credentials?: ChannelCredentials;
	timeout?: number;
}

/**
 * Default gRPC client configuration
 * Can be overridden via environment variables
 */
const defaultConfig: GrpcClientConfig = {
	serverUrl: import.meta.env.VITE_GRPC_SERVER_URL || 'localhost:50051',
	credentials: credentials.createInsecure(),
	timeout: 30000 // 30 seconds default timeout
};

/**
 * Create a GameService client instance
 * @param config Optional configuration override
 * @returns GameServiceClient instance
 */
export function createGameServiceClient(config?: Partial<GrpcClientConfig>): GameServiceClient {
	const finalConfig = { ...defaultConfig, ...config };
	return new GameServiceClient(finalConfig.serverUrl, finalConfig.credentials!);
}

/**
 * Create a LobbyService client instance
 * @param config Optional configuration override
 * @returns LobbyServiceClient instance
 */
export function createLobbyServiceClient(config?: Partial<GrpcClientConfig>): LobbyServiceClient {
	const finalConfig = { ...defaultConfig, ...config };
	return new LobbyServiceClient(finalConfig.serverUrl, finalConfig.credentials!);
}

/**
 * Test gRPC connection
 * @param client Any gRPC client
 * @returns Promise that resolves if connection is successful
 */
export async function testConnection(client: GameServiceClient | LobbyServiceClient): Promise<boolean> {
	return new Promise((resolve) => {
		// Use waitForReady to test connection
		const deadline = new Date();
		deadline.setSeconds(deadline.getSeconds() + 5);

		client.waitForReady(deadline, (error) => {
			if (error) {
				console.error('gRPC connection test failed:', error);
				resolve(false);
			} else {
				console.log('gRPC connection successful');
				resolve(true);
			}
		});
	});
}

/**
 * Example usage:
 *
 * import { createGameServiceClient, createLobbyServiceClient } from '$lib/grpc/client';
 *
 * const gameClient = createGameServiceClient();
 * const lobbyClient = createLobbyServiceClient();
 *
 * // Use the clients to make RPC calls
 * lobbyClient.listTables({ formatFilter: '', openOnly: false }, (err, response) => {
 *   if (err) {
 *     console.error('Error:', err);
 *   } else {
 *     console.log('Tables:', response.tables);
 *   }
 * });
 */
