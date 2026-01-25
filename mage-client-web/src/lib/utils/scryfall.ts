/**
 * Scryfall API utilities for fetching card images
 *
 * Scryfall provides free card images for Magic: The Gathering cards.
 * API docs: https://scryfall.com/docs/api
 */

import type { CardView } from '$lib/generated/mage/v1/models';

/**
 * Generate a Scryfall image URL for a card.
 * Prefers direct CDN URLs from the database if available, otherwise falls back to the redirect API.
 *
 * @param cardName - The exact name of the card
 * @param version - Image version: 'small' | 'normal' | 'large' | 'png' | 'art_crop' | 'border_crop'
 * @param card - Optional CardView with direct image URIs from the database
 * @returns The Scryfall image URL
 */
export function getScryfallImageUrl(
  cardName: string,
  version: 'small' | 'normal' | 'large' | 'png' | 'art_crop' | 'border_crop' = 'normal',
  card?: CardView
): string {
  if (!cardName) return '';

  // Use direct CDN URLs if available (much faster, no redirect needed)
  if (card) {
    switch (version) {
      case 'small':
        if (card.imageUriSmall) return card.imageUriSmall;
        break;
      case 'normal':
        if (card.imageUriNormal) return card.imageUriNormal;
        break;
      case 'large':
        if (card.imageUriLarge) return card.imageUriLarge;
        break;
      case 'png':
        if (card.imageUriPng) return card.imageUriPng;
        break;
    }
  }

  // Fall back to redirect API if direct URLs not available
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
export function getScryfallVersionForSize(
  size: 'small' | 'normal' | 'large'
): 'small' | 'normal' | 'large' {
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

/**
 * Generate a Scryfall search URL for a token with specific characteristics.
 * This helps find the correct token variant when multiple versions exist.
 *
 * @param cardName - The token name (e.g., "Bird", "Soldier")
 * @param power - The power value (for creatures)
 * @param toughness - The toughness value (for creatures)
 * @param color - The color (e.g., "white", "blue", "red", "green", "black", "colorless")
 * @returns A Scryfall search URL that returns JSON with matching token cards
 */
export function getScryfallTokenSearchUrl(
  cardName: string,
  power?: string,
  toughness?: string,
  color?: string
): string {
  if (!cardName) return '';

  // Build a Scryfall search query for tokens
  const parts = ['t:token'];

  // Add name search (exact match preferred)
  parts.push(`name:"${cardName}"`);

  // Add power/toughness if provided
  if (power) {
    parts.push(`pow=${power}`);
  }
  if (toughness) {
    parts.push(`tou=${toughness}`);
  }

  // Add color identity if provided
  if (color && color !== 'colorless' && color !== 'multicolor') {
    const colorMap: Record<string, string> = {
      white: 'w',
      blue: 'u',
      black: 'b',
      red: 'r',
      green: 'g'
    };
    const colorCode = colorMap[color.toLowerCase()];
    if (colorCode) {
      parts.push(`c:${colorCode}`);
    }
  } else if (color === 'colorless') {
    parts.push('c:c');
  }

  // Combine query parts
  const query = parts.join(' ');
  const encodedQuery = encodeURIComponent(query);

  return `https://api.scryfall.com/cards/search?q=${encodedQuery}&order=released`;
}
