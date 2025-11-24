/**
 * Authentication related types
 */

/**
 * User information stored in the auth state
 */
export interface User {
	id: string;
	username: string;
}

/**
 * Complete authentication state
 */
export interface AuthState {
	isAuthenticated: boolean;
	token: string | null;
	user: User | null;
}

/**
 * Decoded JWT token payload
 */
export interface JwtPayload {
	sub: string; // Subject (user ID)
	username: string;
	email?: string; // Optional email
	sessionId?: string; // Optional session ID (for session-based tokens)
	exp: number; // Expiration timestamp
	iat: number; // Issued at timestamp
}

/**
 * Login credentials
 */
export interface LoginCredentials {
	username: string;
	password: string;
	rememberMe?: boolean;
}

/**
 * Registration data
 */
export interface RegisterData {
	username: string;
	password: string;
}
