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
		email: payload.email
	};
}
