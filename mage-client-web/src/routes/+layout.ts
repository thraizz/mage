/**
 * Root layout load function
 * Initializes authentication state and connection from localStorage on app startup
 */

import { browser } from '$app/environment';
import { auth } from '$lib/stores/auth';
import { connection } from '$lib/stores/connection';
import type { LayoutLoad } from './$types';

/**
 * Load authentication state before rendering any page
 * This ensures the auth store and connection are initialized before route guards run
 */
export const load: LayoutLoad = async () => {
	// Only run on client-side (browser)
	if (browser) {
		// Load authentication from localStorage if available
		auth.loadAuthFromStorage();
		
		// Initialize connection with session keep-alive
		// Server lease period is 120 seconds, so ping every 60 seconds to keep session alive
		connection.initialize({
			enableHealthCheck: true,
			healthCheckInterval: 60000, // Ping every 60 seconds (server lease is 120 seconds)
			autoReconnect: true
		});
	}

	return {};
};

