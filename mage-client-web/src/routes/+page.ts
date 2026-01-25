/**
 * Root page load function
 * Redirects to login if not authenticated, or to lobby if authenticated
 */

import { browser } from '$app/environment';
import { redirect } from '@sveltejs/kit';
import { isAuthenticated } from '$lib/utils/auth-guard';
import type { PageLoad } from './$types';

/**
 * Page load function - runs before rendering the root page
 * Automatically redirects based on authentication status
 */
export const load: PageLoad = async () => {
  // Only run on client-side (browser)
  if (!browser) {
    return {};
  }

  // Check authentication status
  if (isAuthenticated()) {
    // User is logged in, redirect to lobby
    throw redirect(303, '/lobby');
  } else {
    // User is not logged in, redirect to login page
    throw redirect(303, '/login');
  }
};
