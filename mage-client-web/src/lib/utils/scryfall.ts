/**
 * Scryfall API utilities for fetching card images
 * 
 * Scryfall provides free card images for Magic: The Gathering cards.
 * API docs: https://scryfall.com/docs/api
 */

/**
 * Generate a Scryfall image URL for a card by exact name.
 * Uses the redirect endpoint which returns a 302 to the actual image.
 * 
 * @param cardName - The exact name of the card
 * @param version - Image version: 'small' | 'normal' | 'large' | 'png' | 'art_crop' | 'border_crop'
 * @returns The Scryfall image URL
 */
export function getScryfallImageUrl(
	cardName: string,
	version: 'small' | 'normal' | 'large' | 'png' | 'art_crop' | 'border_crop' = 'normal'
): string {
	if (!cardName) return '';

	// Use the named card image redirect endpoint
	// This is faster than making an API call to get the image URI
	const encodedName = encodeURIComponent(cardName);
	return `https://api.scryfall.com/cards/named?format=image&version=${version}&exact=${encodedName}`;
}

/**
 * Generate a Scryfall image URL with fuzzy matching (more forgiving of typos)
 * 
 * @param cardName - The card name (fuzzy matched)
 * @param version - Image version
 * @returns The Scryfall image URL
 */
export function getScryfallImageUrlFuzzy(
	cardName: string,
	version: 'small' | 'normal' | 'large' | 'png' | 'art_crop' | 'border_crop' = 'normal'
): string {
	if (!cardName) return '';

	const encodedName = encodeURIComponent(cardName);
	return `https://api.scryfall.com/cards/named?format=image&version=${version}&fuzzy=${encodedName}`;
}

/**
 * Get the appropriate image size for card display
 */
export function getScryfallVersionForSize(size: 'small' | 'normal' | 'large'): 'small' | 'normal' | 'large' {
	switch (size) {
		case 'small':
			return 'small'; // 146 × 204
		case 'normal':
			return 'normal'; // 488 × 680
		case 'large':
			return 'large'; // 672 × 936
		default:
			return 'normal';
	}
}

