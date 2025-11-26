/**
 * Game API functions for game execution operations
 * Real gRPC implementation for MTG game interactions
 */

import { getMageClient } from '$lib/grpc/client';
import type {
	GameJoinRequest,
	GameJoinResponse,
	GameGetViewRequest,
	GameGetViewResponse,
	SendPlayerActionRequest,
	SendPlayerActionResponse,
	SendPlayerUUIDRequest,
	SendPlayerUUIDResponse,
	SendPlayerBooleanRequest,
	SendPlayerBooleanResponse,
	SendPlayerIntegerRequest,
	SendPlayerIntegerResponse,
	SendPlayerStringRequest,
	SendPlayerStringResponse,
	SendPlayerManaTypeRequest,
	SendPlayerManaTypeResponse,
	SendSpecialActionRequest,
	SendSpecialActionResponse,
	MatchQuitRequest,
	MatchQuitResponse,
	PlayerAction
} from '$lib/generated/mage/v1/game';
import { SpecialActionType } from '$lib/generated/mage/v1/game';
import type { GameView } from '$lib/generated/mage/v1/models';

// Re-export enums for convenience
export { PlayerAction, SpecialActionType } from '$lib/generated/mage/v1/game';

/**
 * Join a game as a player
 * Must be called before receiving game updates
 */
export async function joinGame(gameId: string): Promise<void> {
	const client = getMageClient();
	const sessionId = await client.ensureSessionId();

	if (!sessionId) {
		throw new Error('No active session - please login first');
	}

	const request: GameJoinRequest = {
		sessionId,
		gameId
	};

	const response = await client.call<GameJoinRequest, GameJoinResponse>('GameJoin', request);

	if (!response.success) {
		throw new Error(response.error || 'Failed to join game');
	}
}

/**
 * Get current game view/state
 * Returns the full game state from the player's perspective
 */
export async function fetchGameView(gameId: string, playerId?: string): Promise<GameView> {
	const client = getMageClient();
	const sessionId = await client.ensureSessionId();

	if (!sessionId) {
		throw new Error('No active session - please login first');
	}

	const request: GameGetViewRequest = {
		sessionId,
		gameId,
		playerId: playerId || ''
	};

	const response = await client.call<GameGetViewRequest, GameGetViewResponse>(
		'GameGetView',
		request
	);

	if (!response.game) {
		throw new Error('Failed to get game view - no game data returned');
	}

	return response.game;
}

/**
 * Send a player action (pass priority, concede, etc.)
 */
export async function sendPlayerAction(gameId: string, action: PlayerAction): Promise<void> {
	const client = getMageClient();
	const sessionId = await client.ensureSessionId();

	if (!sessionId) {
		throw new Error('No active session - please login first');
	}

	const request: SendPlayerActionRequest = {
		sessionId,
		gameId,
		action
	};

	const response = await client.call<SendPlayerActionRequest, SendPlayerActionResponse>(
		'SendPlayerAction',
		request
	);

	if (!response.success) {
		throw new Error(response.error || 'Failed to send action');
	}
}

/**
 * Send a UUID selection (card, permanent, ability target)
 */
export async function sendPlayerUUID(gameId: string, uuid: string): Promise<void> {
	const client = getMageClient();
	const sessionId = await client.ensureSessionId();

	if (!sessionId) {
		throw new Error('No active session - please login first');
	}

	const request: SendPlayerUUIDRequest = {
		sessionId,
		gameId,
		uuid
	};

	const response = await client.call<SendPlayerUUIDRequest, SendPlayerUUIDResponse>(
		'SendPlayerUUID',
		request
	);

	if (!response.success) {
		throw new Error(response.error || 'Failed to send UUID selection');
	}
}

/**
 * Send a boolean response (yes/no choices)
 */
export async function sendPlayerBoolean(gameId: string, value: boolean): Promise<void> {
	const client = getMageClient();
	const sessionId = await client.ensureSessionId();

	if (!sessionId) {
		throw new Error('No active session - please login first');
	}

	const request: SendPlayerBooleanRequest = {
		sessionId,
		gameId,
		data: value
	};

	const response = await client.call<SendPlayerBooleanRequest, SendPlayerBooleanResponse>(
		'SendPlayerBoolean',
		request
	);

	if (!response.success) {
		throw new Error(response.error || 'Failed to send boolean response');
	}
}

/**
 * Send an integer response (amount selections)
 */
export async function sendPlayerInteger(gameId: string, value: number): Promise<void> {
	const client = getMageClient();
	const sessionId = await client.ensureSessionId();

	if (!sessionId) {
		throw new Error('No active session - please login first');
	}

	const request: SendPlayerIntegerRequest = {
		sessionId,
		gameId,
		data: value
	};

	const response = await client.call<SendPlayerIntegerRequest, SendPlayerIntegerResponse>(
		'SendPlayerInteger',
		request
	);

	if (!response.success) {
		throw new Error(response.error || 'Failed to send integer response');
	}
}

/**
 * Send a string response (mode/choice selections)
 */
export async function sendPlayerString(gameId: string, value: string): Promise<void> {
	const client = getMageClient();
	const sessionId = await client.ensureSessionId();

	if (!sessionId) {
		throw new Error('No active session - please login first');
	}

	const request: SendPlayerStringRequest = {
		sessionId,
		gameId,
		data: value
	};

	const response = await client.call<SendPlayerStringRequest, SendPlayerStringResponse>(
		'SendPlayerString',
		request
	);

	if (!response.success) {
		throw new Error(response.error || 'Failed to send string response');
	}
}

/**
 * Send a mana type selection
 */
export async function sendPlayerManaType(
	gameId: string,
	manaType: string,
	manaTypeStr?: string
): Promise<void> {
	const client = getMageClient();
	const sessionId = await client.ensureSessionId();

	if (!sessionId) {
		throw new Error('No active session - please login first');
	}

	const request: SendPlayerManaTypeRequest = {
		sessionId,
		gameId,
		manaType,
		manaTypeStr: manaTypeStr || manaType
	};

	const response = await client.call<SendPlayerManaTypeRequest, SendPlayerManaTypeResponse>(
		'SendPlayerManaType',
		request
	);

	if (!response.success) {
		throw new Error(response.error || 'Failed to send mana type selection');
	}
}

/**
 * Concede the game
 */
export async function concedeGame(gameId: string): Promise<void> {
	return sendPlayerAction(gameId, 4 as PlayerAction); // CONCEDE = 4
}

/**
 * Pass priority
 */
export async function passPriority(gameId: string): Promise<void> {
	return sendPlayerAction(gameId, 1 as PlayerAction); // PASS = 1
}

/**
 * Pass until end of turn (F6 in XMage)
 */
export async function passUntilEndOfTurn(gameId: string): Promise<void> {
	return sendPlayerAction(gameId, 5 as PlayerAction); // PASS_UNTIL_END_OF_TURN = 5
}

/**
 * Pass until next turn
 */
export async function passUntilNextTurn(gameId: string): Promise<void> {
	return sendPlayerAction(gameId, 6 as PlayerAction); // PASS_UNTIL_NEXT_TURN = 6
}

/**
 * Pass until stack is resolved
 */
export async function passUntilStackResolved(gameId: string): Promise<void> {
	return sendPlayerAction(gameId, 7 as PlayerAction); // PASS_UNTIL_STACK_RESOLVED = 7
}

/**
 * Pass until my next turn
 */
export async function passUntilMyNextTurn(gameId: string): Promise<void> {
	return sendPlayerAction(gameId, 8 as PlayerAction); // PASS_UNTIL_MY_NEXT_TURN = 8
}

/**
 * Keep current hand during mulligan phase
 */
export async function keepHand(gameId: string): Promise<void> {
	return sendPlayerString(gameId, 'KEEP');
}

/**
 * Mulligan (redraw hand with one fewer card)
 */
export async function mulligan(gameId: string): Promise<void> {
	return sendPlayerString(gameId, 'MULLIGAN');
}

/**
 * Quit the match entirely
 */
export async function quitMatch(gameId: string): Promise<void> {
	const client = getMageClient();
	const sessionId = await client.ensureSessionId();

	if (!sessionId) {
		throw new Error('No active session - please login first');
	}

	const request: MatchQuitRequest = {
		sessionId,
		gameId
	};

	const response = await client.call<MatchQuitRequest, MatchQuitResponse>('MatchQuit', request);

	if (!response.success) {
		throw new Error(response.error || 'Failed to quit match');
	}
}

/**
 * Send a special action (play land, foretell, etc.)
 * Per MTG Rule 116.2a: Special actions don't use the stack
 */
export async function sendSpecialAction(
	gameId: string,
	actionType: SpecialActionType,
	sourceId: string
): Promise<void> {
	const client = getMageClient();
	const sessionId = await client.ensureSessionId();

	if (!sessionId) {
		throw new Error('No active session - please login first');
	}

	const request: SendSpecialActionRequest = {
		sessionId,
		gameId,
		actionType,
		sourceId
	};

	const response = await client.call<SendSpecialActionRequest, SendSpecialActionResponse>(
		'SendSpecialAction',
		request
	);

	if (!response.success) {
		throw new Error(response.error || 'Failed to execute special action');
	}
}

/**
 * Play a land from hand
 * Convenience wrapper for sendSpecialAction with PLAY_LAND
 */
export async function playLand(gameId: string, cardId: string): Promise<void> {
	return sendSpecialAction(gameId, SpecialActionType.PLAY_LAND, cardId);
}
