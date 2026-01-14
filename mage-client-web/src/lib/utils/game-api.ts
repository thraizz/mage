/**
 * Game API Utilities
 *
 * Shared utilities for making gRPC game API calls.
 * Reduces boilerplate for session management and error handling.
 */

import { getMageClient } from '$lib/grpc/client';

/**
 * Execute a game API call with automatic session management
 *
 * @param method - The gRPC method name (e.g., 'GameJoin', 'SendPlayerAction')
 * @param buildRequest - Function that builds the request object given a sessionId
 * @returns The response from the gRPC call
 * @throws Error if no active session or if the call fails
 */
export async function withSession<TRequest, TResponse>(
	method: string,
	buildRequest: (sessionId: string) => TRequest
): Promise<TResponse> {
	const client = getMageClient();
	const sessionId = await client.ensureSessionId();

	if (!sessionId) {
		throw new Error('No active session - please login first');
	}

	return client.call<TRequest, TResponse>(method, buildRequest(sessionId));
}

/**
 * Execute a game API call with automatic session management and success validation
 *
 * @param method - The gRPC method name
 * @param buildRequest - Function that builds the request object given a sessionId
 * @param errorMessage - Custom error message prefix (default: 'API call failed')
 * @throws Error if no active session, if the call fails, or if response.success is false
 */
export async function withSessionValidated<
	TRequest,
	TResponse extends { success: boolean; error?: string }
>(
	method: string,
	buildRequest: (sessionId: string) => TRequest,
	errorMessage: string = 'API call failed'
): Promise<TResponse> {
	const response = await withSession<TRequest, TResponse>(method, buildRequest);

	if (!response.success) {
		throw new Error(response.error || errorMessage);
	}

	return response;
}

/**
 * Build a standard game request with session and game ID
 */
export function gameRequest(
	sessionId: string,
	gameId: string
): { sessionId: string; gameId: string } {
	return { sessionId, gameId };
}

/**
 * Build a game request with additional data
 */
export function gameRequestWithData<T>(
	sessionId: string,
	gameId: string,
	data: T
): { sessionId: string; gameId: string } & T {
	return { sessionId, gameId, ...data };
}
