/**
 * Game page load function
 * Ensures game ID is properly passed to the page component
 */

import type { PageLoad } from './$types';

export const load: PageLoad = async ({ params }) => {
	return {
		gameId: params.id
	};
};
