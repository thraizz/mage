/**
 * Utility functions for handling session errors and redirecting to login
 */

import { browser } from '$app/environment';
import { goto } from '$app/navigation';
import { page } from '$app/stores';
import { get } from 'svelte/store';
import { auth } from '$lib/stores/auth';
import { toast } from '$lib/stores/toast';

/**
 * Handle session errors by logging out and redirecting to login
 * This should be called when we detect that the session is invalid
 */
export function handleSessionError(message = 'Session expired. Please log in again.'): void {
	// Only handle in browser context
	if (!browser) {
		return;
	}

	// Logout user
	auth.logout();

	// Show error toast
	toast.error(message);

	// Get current page URL for return redirect
	const currentPage = get(page);
	const returnUrl = currentPage.url.pathname + currentPage.url.search;

	// Redirect to login with return URL
	goto(`/login?returnUrl=${encodeURIComponent(returnUrl)}`);
}

/**
 * Check if a response has a session error
 * Returns true if the response indicates a session error
 */
export function isSessionErrorResponse(response: { error?: string; success?: boolean }): boolean {
	if (!response) {
		return false;
	}

	// Check if response has an error field indicating session issues
	if (response.error) {
		const errorMsg = response.error.toLowerCase();
		return (
			errorMsg.includes('session not found') ||
			errorMsg.includes('invalid or expired session') ||
			errorMsg.includes('missing session') ||
			errorMsg.includes('session expired') ||
			errorMsg.includes('unauthorized') ||
			errorMsg.includes('unauthenticated')
		);
	}

	// Check if success is false and we're in an authenticated context
	// (this might indicate session issues, but we need to be careful not to
	// trigger on all failed requests)
	if (response.success === false) {
		const authState = get(auth);
		if (authState.isAuthenticated) {
			// If we're authenticated but got a failure, it might be a session issue
			// We'll let the caller decide based on the error message
			return false; // Don't auto-detect based on success=false alone
		}
	}

	return false;
}
