/**
 * JWT utility functions for token handling
 */

import type { JwtPayload } from '$lib/types/auth';

/**
 * Decode a JWT token without verification
 * NOTE: This is client-side decoding only - tokens should be verified on the server!
 *
 * @param token - JWT token string
 * @returns Decoded payload or null if invalid
 */
export function decodeJwt(token: string): JwtPayload | null {
	try {
		// JWT format: header.payload.signature
		const parts = token.split('.');
		if (parts.length !== 3) {
			console.error('Invalid JWT format: expected 3 parts');
			return null;
		}

		// Decode the payload (second part)
		const payload = parts[1];

		// Base64url decode
		const base64 = payload.replace(/-/g, '+').replace(/_/g, '/');
		const jsonPayload = decodeURIComponent(
			atob(base64)
				.split('')
				.map((c) => '%' + ('00' + c.charCodeAt(0).toString(16)).slice(-2))
				.join('')
		);

		return JSON.parse(jsonPayload) as JwtPayload;
	} catch (error) {
		console.error('Failed to decode JWT:', error);
		return null;
	}
}

/**
 * Check if a JWT token is expired
 *
 * @param token - JWT token string
 * @returns true if expired, false if valid
 */
export function isTokenExpired(token: string): boolean {
	const payload = decodeJwt(token);
	if (!payload || !payload.exp) {
		return true;
	}

	// exp is in seconds, Date.now() is in milliseconds
	const currentTime = Date.now() / 1000;
	return payload.exp < currentTime;
}

/**
 * Get time remaining until token expiration in seconds
 *
 * @param token - JWT token string
 * @returns seconds until expiration, or 0 if expired/invalid
 */
export function getTokenTimeRemaining(token: string): number {
	const payload = decodeJwt(token);
	if (!payload || !payload.exp) {
		return 0;
	}

	const currentTime = Date.now() / 1000;
	const remaining = payload.exp - currentTime;
	return Math.max(0, remaining);
}

/**
 * Extract user information from JWT token
 *
 * @param token - JWT token string
 * @returns User info object or null if invalid
 */
export function getUserFromToken(
	token: string
): { id: string; username: string; email: string } | null {
	const payload = decodeJwt(token);
	if (!payload) {
		return null;
	}

	return {
		id: payload.sub,
		username: payload.username,
		email: payload.email || `${payload.username}@example.com`
	};
}

/**
 * Extract session ID from JWT token
 *
 * @param token - JWT token string
 * @returns Session ID or null if not found
 */
export function getSessionIdFromToken(token: string): string | null {
	const payload = decodeJwt(token);
	if (!payload) {
		return null;
	}
	return payload.sessionId || null;
}

/**
 * Create a session-based token from server response
 * This creates a JWT-like token containing session information
 *
 * @param sessionId - Session ID from server
 * @param userId - User ID from server
 * @param username - Username
 * @param email - Optional email address
 * @param expiresIn - Expiration time in seconds (default: 24 hours)
 * @returns JWT-formatted token string
 */
export function createSessionToken(
	sessionId: string,
	userId: string,
	username: string,
	email?: string,
	expiresIn: number = 86400 // 24 hours
): string {
	const now = Math.floor(Date.now() / 1000);
	const exp = now + expiresIn;

	const payload = {
		sub: userId,
		sessionId: sessionId, // Store sessionId in payload
		username: username,
		email: email || `${username}@example.com`,
		exp,
		iat: now
	};

	const payloadStr = JSON.stringify(payload);
	const encodedPayload = btoa(payloadStr);

	// Create a JWT-like token: header.payload.signature
	// For session tokens, we use "session" as the header type
	const header = btoa(JSON.stringify({ typ: 'JWT', alg: 'session' }));

	// Use sessionId as part of the signature for validation
	// In a real implementation, this would be signed by the server
	const signature = btoa(sessionId).slice(0, 16); // Truncate for consistency

	return `${header}.${encodedPayload}.${signature}`;
}
