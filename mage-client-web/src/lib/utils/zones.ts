/**
 * Zone Constants and Utilities
 *
 * Centralized constants and helpers for MTG card zones.
 * Maps between zone names (strings) and numeric zone IDs used in the game state.
 */

import type { CardZone } from '$lib/types/game';

/**
 * Numeric zone IDs used in CardView.zone property
 */
export enum ZoneId {
	LIBRARY = 0,
	HAND = 1,
	BATTLEFIELD = 2,
	GRAVEYARD = 3,
	EXILE = 4,
	COMMAND = 5,
	STACK = 6
}

/**
 * Map zone names to numeric IDs
 */
export const ZONE_NAME_TO_ID: Record<CardZone, ZoneId> = {
	LIBRARY: ZoneId.LIBRARY,
	HAND: ZoneId.HAND,
	BATTLEFIELD: ZoneId.BATTLEFIELD,
	GRAVEYARD: ZoneId.GRAVEYARD,
	EXILE: ZoneId.EXILE,
	STACK: ZoneId.STACK,
	COMMAND: ZoneId.COMMAND
};

/**
 * Map numeric IDs to zone names
 */
export const ZONE_ID_TO_NAME: Record<ZoneId, CardZone> = {
	[ZoneId.LIBRARY]: 'LIBRARY',
	[ZoneId.HAND]: 'HAND',
	[ZoneId.BATTLEFIELD]: 'BATTLEFIELD',
	[ZoneId.GRAVEYARD]: 'GRAVEYARD',
	[ZoneId.EXILE]: 'EXILE',
	[ZoneId.STACK]: 'STACK',
	[ZoneId.COMMAND]: 'COMMAND'
};

/**
 * Zones that are player-specific (each player has their own)
 */
export const PLAYER_ZONES: ReadonlySet<CardZone> = new Set([
	'HAND',
	'LIBRARY',
	'GRAVEYARD'
] as const);

/**
 * Zones that are shared (global game zones)
 */
export const SHARED_ZONES: ReadonlySet<CardZone> = new Set([
	'BATTLEFIELD',
	'EXILE',
	'STACK',
	'COMMAND'
] as const);

/**
 * Convert zone name to numeric ID
 */
export function zoneNameToId(name: string): ZoneId {
	const upperName = name.toUpperCase();
	const zoneId = ZONE_NAME_TO_ID[upperName as CardZone];

	if (zoneId === undefined) {
		throw new Error(`Invalid zone name: ${name}`);
	}

	return zoneId;
}

/**
 * Convert numeric zone ID to zone name
 */
export function zoneIdToName(id: number): CardZone {
	const zoneName = ZONE_ID_TO_NAME[id as ZoneId];

	if (!zoneName) {
		throw new Error(`Invalid zone ID: ${id}`);
	}

	return zoneName;
}

/**
 * Check if a zone name is valid
 */
export function isValidZoneName(zone: unknown): zone is CardZone {
	if (typeof zone !== 'string') return false;
	const upperZone = zone.toUpperCase();
	return upperZone in ZONE_NAME_TO_ID;
}

/**
 * Check if a zone is player-specific
 */
export function isPlayerZone(zone: CardZone | ZoneId): boolean {
	const zoneName = typeof zone === 'number' ? zoneIdToName(zone) : zone;
	return PLAYER_ZONES.has(zoneName);
}

/**
 * Check if a zone is shared/global
 */
export function isSharedZone(zone: CardZone | ZoneId): boolean {
	const zoneName = typeof zone === 'number' ? zoneIdToName(zone) : zone;
	return SHARED_ZONES.has(zoneName);
}

/**
 * Normalize zone name (handles LIBRARY_TOP, LIBRARY_BOTTOM variants)
 */
export function normalizeZoneName(zone: string): CardZone {
	const upperZone = zone.toUpperCase();

	// Handle library variants
	if (upperZone === 'LIBRARY_TOP' || upperZone === 'LIBRARY_BOTTOM') {
		return 'LIBRARY';
	}

	if (!isValidZoneName(upperZone)) {
		throw new Error(`Invalid zone name: ${zone}`);
	}

	return upperZone as CardZone;
}

/**
 * Check if zone string indicates library top placement
 */
export function isLibraryTop(zone: string): boolean {
	return zone.toUpperCase() === 'LIBRARY_TOP' || zone.toUpperCase() === 'LIBRARY';
}

/**
 * Check if zone string indicates library bottom placement
 */
export function isLibraryBottom(zone: string): boolean {
	return zone.toUpperCase() === 'LIBRARY_BOTTOM';
}
