/**
 * Protected routes layout - Auth guard
 * Ensures users are authenticated before accessing protected pages
 */

import { browser } from '$app/environment';
import { redirect } from '@sveltejs/kit';
import { isAuthenticated, clearInvalidToken } from '$lib/utils/auth-guard';
import type { LayoutLoad } from './$types';

/**
 * Layout load function - runs before rendering protected pages
 * Redirects to login if not authenticated
 */
export const load: LayoutLoad = async ({ url }) => {
  // On server-side, allow through (client will handle redirect)
  if (!browser) {
    return {};
  }

  // Clear any invalid tokens
  clearInvalidToken();

  // Check authentication
  if (!isAuthenticated()) {
    // Store the original URL to redirect back after login
    const returnUrl = url.pathname + url.search;

    // Redirect to login with return URL
    throw redirect(303, `/login?returnUrl=${encodeURIComponent(returnUrl)}`);
  }

  // User is authenticated, allow access
  return {};
};
