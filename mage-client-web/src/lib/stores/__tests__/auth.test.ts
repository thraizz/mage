/**
 * Auth store tests
 */

import { describe, it, expect, beforeEach } from 'vitest';
import { get } from 'svelte/store';
import { auth } from '../auth';

// Mock localStorage
const localStorageMock = (() => {
	let store: Record<string, string> = {};
	return {
		getItem: (key: string) => store[key] || null,
		setItem: (key: string, value: string) => {
			store[key] = value;
		},
		removeItem: (key: string) => {
			delete store[key];
		},
		clear: () => {
			store = {};
		}
	};
})();

// @ts-expect-error - Mocking localStorage for tests
global.localStorage = localStorageMock;

describe('Auth Store', () => {
	beforeEach(() => {
		localStorageMock.clear();
		auth.logout(); // Reset to initial state
	});

	it('should have initial unauthenticated state', () => {
		const state = get(auth);
		expect(state.isAuthenticated).toBe(false);
		expect(state.token).toBeNull();
		expect(state.user).toBeNull();
	});

	it('should login with valid token', () => {
		// Create a valid JWT token (expires in 1 hour)
		const futureExp = Math.floor(Date.now() / 1000) + 3600;
		const payload = {
			sub: 'user123',
			username: 'testuser',
			email: 'test@example.com',
			exp: futureExp,
			iat: Math.floor(Date.now() / 1000)
		};
		const encodedPayload = btoa(JSON.stringify(payload));
		const token = `header.${encodedPayload}.signature`;

		const user = {
			id: 'user123',
			username: 'testuser',
			email: 'test@example.com'
		};

		auth.login(token, user);

		const state = get(auth);
		expect(state.isAuthenticated).toBe(true);
		expect(state.token).toBe(token);
		expect(state.user).toEqual(user);
		expect(localStorageMock.getItem('mage_auth_token')).toBe(token);
	});

	it('should not login with expired token', () => {
		// Create an expired JWT token
		const pastExp = Math.floor(Date.now() / 1000) - 3600; // 1 hour ago
		const payload = {
			sub: 'user123',
			username: 'testuser',
			email: 'test@example.com',
			exp: pastExp,
			iat: Math.floor(Date.now() / 1000) - 7200
		};
		const encodedPayload = btoa(JSON.stringify(payload));
		const token = `header.${encodedPayload}.signature`;

		const user = {
			id: 'user123',
			username: 'testuser',
			email: 'test@example.com'
		};

		auth.login(token, user);

		const state = get(auth);
		expect(state.isAuthenticated).toBe(false);
		expect(state.token).toBeNull();
		expect(localStorageMock.getItem('mage_auth_token')).toBeNull();
	});

	it('should logout and clear state', () => {
		// First login
		const futureExp = Math.floor(Date.now() / 1000) + 3600;
		const payload = {
			sub: 'user123',
			username: 'testuser',
			email: 'test@example.com',
			exp: futureExp,
			iat: Math.floor(Date.now() / 1000)
		};
		const encodedPayload = btoa(JSON.stringify(payload));
		const token = `header.${encodedPayload}.signature`;

		auth.login(token, {
			id: 'user123',
			username: 'testuser',
			email: 'test@example.com'
		});

		// Then logout
		auth.logout();

		const state = get(auth);
		expect(state.isAuthenticated).toBe(false);
		expect(state.token).toBeNull();
		expect(state.user).toBeNull();
		expect(localStorageMock.getItem('mage_auth_token')).toBeNull();
	});

	it('should load auth from storage with valid token', () => {
		// Create a valid JWT token and store it
		const futureExp = Math.floor(Date.now() / 1000) + 3600;
		const payload = {
			sub: 'user123',
			username: 'testuser',
			email: 'test@example.com',
			exp: futureExp,
			iat: Math.floor(Date.now() / 1000)
		};
		const encodedPayload = btoa(JSON.stringify(payload));
		const token = `header.${encodedPayload}.signature`;

		localStorageMock.setItem('mage_auth_token', token);

		const result = auth.loadAuthFromStorage();

		expect(result).toBe(true);
		const state = get(auth);
		expect(state.isAuthenticated).toBe(true);
		expect(state.token).toBe(token);
		expect(state.user?.username).toBe('testuser');
	});

	it('should not load auth from storage with expired token', () => {
		// Create an expired token and store it
		const pastExp = Math.floor(Date.now() / 1000) - 3600;
		const payload = {
			sub: 'user123',
			username: 'testuser',
			email: 'test@example.com',
			exp: pastExp,
			iat: Math.floor(Date.now() / 1000) - 7200
		};
		const encodedPayload = btoa(JSON.stringify(payload));
		const token = `header.${encodedPayload}.signature`;

		localStorageMock.setItem('mage_auth_token', token);

		const result = auth.loadAuthFromStorage();

		expect(result).toBe(false);
		const state = get(auth);
		expect(state.isAuthenticated).toBe(false);
		expect(localStorageMock.getItem('mage_auth_token')).toBeNull();
	});

	it('should check token validity', () => {
		// Login with valid token
		const futureExp = Math.floor(Date.now() / 1000) + 3600;
		const payload = {
			sub: 'user123',
			username: 'testuser',
			email: 'test@example.com',
			exp: futureExp,
			iat: Math.floor(Date.now() / 1000)
		};
		const encodedPayload = btoa(JSON.stringify(payload));
		const token = `header.${encodedPayload}.signature`;

		auth.login(token, {
			id: 'user123',
			username: 'testuser',
			email: 'test@example.com'
		});

		const isValid = auth.checkTokenValidity();
		expect(isValid).toBe(true);
	});

	it('should update user information', () => {
		// Login first
		const futureExp = Math.floor(Date.now() / 1000) + 3600;
		const payload = {
			sub: 'user123',
			username: 'testuser',
			email: 'test@example.com',
			exp: futureExp,
			iat: Math.floor(Date.now() / 1000)
		};
		const encodedPayload = btoa(JSON.stringify(payload));
		const token = `header.${encodedPayload}.signature`;

		auth.login(token, {
			id: 'user123',
			username: 'testuser',
			email: 'test@example.com'
		});

		// Update user info
		auth.updateUser({ email: 'newemail@example.com' });

		const state = get(auth);
		expect(state.user?.email).toBe('newemail@example.com');
		expect(state.user?.username).toBe('testuser'); // Should keep other fields
	});
});
