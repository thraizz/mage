/**
 * Authentication store for managing user session state
 */

import { writable } from 'svelte/store';
import type { AuthState, User } from '$lib/types/auth';
import { isTokenExpired, getUserFromToken } from '$lib/utils/jwt';

// Storage key for persisting auth state
const AUTH_STORAGE_KEY = 'mage_auth_token';

/**
 * Initial/default authentication state
 */
const initialState: AuthState = {
	isAuthenticated: false,
	token: null,
	user: null
};

/**
 * Create the authentication store
 */
function createAuthStore() {
	const { subscribe, set, update } = writable<AuthState>(initialState);

	return {
		subscribe,

		/**
		 * Login with a JWT token
		 * Stores the token in localStorage and updates the store state
		 *
		 * @param token - JWT token from server
		 * @param user - User information
		 */
		login(token: string, user: User) {
			// Check if token is expired
			if (isTokenExpired(token)) {
				console.error('Cannot login with expired token');
				return;
			}

			// Store token in localStorage
			if (typeof window !== 'undefined') {
				localStorage.setItem(AUTH_STORAGE_KEY, token);
			}

			// Update store state
			set({
				isAuthenticated: true,
				token,
				user
			});

			console.log('User logged in:', user.username);
		},

		/**
		 * Logout and clear authentication state
		 * Removes token from localStorage and resets store
		 */
		logout() {
			// Remove token from localStorage
			if (typeof window !== 'undefined') {
				localStorage.removeItem(AUTH_STORAGE_KEY);
			}

			// Reset store to initial state
			set(initialState);

			console.log('User logged out');
		},

		/**
		 * Load authentication state from localStorage
		 * Restores the user session if a valid token exists
		 *
		 * @returns true if session was restored, false otherwise
		 */
		loadAuthFromStorage(): boolean {
			if (typeof window === 'undefined') {
				return false;
			}

			const token = localStorage.getItem(AUTH_STORAGE_KEY);

			if (!token) {
				console.log('No stored auth token found');
				return false;
			}

			// Check if token is expired
			if (isTokenExpired(token)) {
				console.log('Stored token is expired, removing');
				localStorage.removeItem(AUTH_STORAGE_KEY);
				return false;
			}

			// Decode token to get user info
			const user = getUserFromToken(token);

			if (!user) {
				console.error('Failed to decode user from token');
				localStorage.removeItem(AUTH_STORAGE_KEY);
				return false;
			}

			// Restore auth state
			set({
				isAuthenticated: true,
				token,
				user
			});

			console.log('Session restored for user:', user.username);
			return true;
		},

		/**
		 * Check if the current token is still valid
		 * If expired, automatically logout
		 *
		 * @returns true if authenticated and token is valid
		 */
		checkTokenValidity(): boolean {
			let valid = false;

			update((state) => {
				if (!state.isAuthenticated || !state.token) {
					valid = false;
					return state;
				}

				if (isTokenExpired(state.token)) {
					console.log('Token expired, logging out');
					valid = false;
					// Logout
					if (typeof window !== 'undefined') {
						localStorage.removeItem(AUTH_STORAGE_KEY);
					}
					return initialState;
				}

				valid = true;
				return state;
			});

			return valid;
		},

		/**
		 * Update user information
		 * Useful for updating profile data without re-authentication
		 *
		 * @param user - Updated user information
		 */
		updateUser(user: Partial<User>) {
			update((state) => {
				if (!state.isAuthenticated || !state.user) {
					return state;
				}

				return {
					...state,
					user: {
						...state.user,
						...user
					}
				};
			});
		}
	};
}

/**
 * Singleton auth store instance
 */
export const auth = createAuthStore();
