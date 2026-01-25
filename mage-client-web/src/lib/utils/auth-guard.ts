/**
 * Auth guard utility functions for route protection
 */

import { browser } from '$app/environment';

/**
 * Storage key for auth token
 */
const AUTH_STORAGE_KEY = 'mage_auth_token';

/**
 * Check if user has a valid authentication token
 * @returns true if authenticated with valid token, false otherwise
 */
export function isAuthenticated(): boolean {
  if (!browser) {
    // On server, can't check localStorage
    return false;
  }

  const token = localStorage.getItem(AUTH_STORAGE_KEY);
  if (!token) {
    return false;
  }

  return isTokenValid(token);
}

/**
 * Validate JWT token structure and expiry
 * @param token - JWT token to validate
 * @returns true if token is valid and not expired
 */
export function isTokenValid(token: string): boolean {
  try {
    const parts = token.split('.');
    if (parts.length !== 3) {
      return false;
    }

    const payload = JSON.parse(atob(parts[1]));
    const exp = payload.exp;

    if (!exp) {
      return false;
    }

    // Check if token is expired
    const now = Math.floor(Date.now() / 1000);
    return exp >= now;
  } catch {
    // Invalid token format
    return false;
  }
}

/**
 * Clear invalid or expired tokens from storage
 */
export function clearInvalidToken(): void {
  if (!browser) return;

  const token = localStorage.getItem(AUTH_STORAGE_KEY);
  if (token && !isTokenValid(token)) {
    localStorage.removeItem(AUTH_STORAGE_KEY);
  }
}
