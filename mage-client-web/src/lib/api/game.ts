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
	ActivateAbilityRequest,
	ActivateAbilityResponse,
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
	console.log('[fetchGameView] Starting', { gameId, playerId });
	
	const client = getMageClient();
	const sessionId = await client.ensureSessionId();
	console.log('[fetchGameView] Got sessionId', { sessionId });

	if (!sessionId) {
		throw new Error('No active session - please login first');
	}

	const request: GameGetViewRequest = {
		sessionId,
		gameId,
		playerId: playerId || ''
	};

	console.log('[fetchGameView] Calling GameGetView RPC...');
	const response = await client.call<GameGetViewRequest, GameGetViewResponse>(
		'GameGetView',
		request
	);
	console.log('[fetchGameView] Got response', { hasGame: !!response.game });

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

/**
 * Advance to the next phase/step
 * Debug/development feature for manual turn progression
 */
export async function advancePhase(gameId: string): Promise<void> {
	return sendSpecialAction(gameId, SpecialActionType.ADVANCE_PHASE, '');
}

/**
 * Activate a mana ability on a permanent (tap for mana)
 * Per MTG Rule 605: Mana abilities don't use the stack and resolve immediately
 */
export async function activateManaAbility(gameId: string, permanentId: string): Promise<void> {
	return sendSpecialAction(gameId, SpecialActionType.ACTIVATE_MANA_ABILITY, permanentId);
}

/**
 * Activate an activated ability on a permanent (non-mana abilities)
 * Per MTG Rule 602: Activated abilities use the stack
 */
export async function activateAbility(
	gameId: string,
	cardId: string,
	abilityId: string,
	targets: string[] = []
): Promise<void> {
	const client = getMageClient();
	const sessionId = await client.ensureSessionId();

	if (!sessionId) {
		throw new Error('No active session - please login first');
	}

	const request: ActivateAbilityRequest = {
		sessionId,
		gameId,
		cardId,
		abilityId,
		targets
	};

	const response = await client.call<ActivateAbilityRequest, ActivateAbilityResponse>(
		'ActivateAbility',
		request
	);

	if (!response.success) {
		throw new Error(response.error || 'Failed to activate ability');
	}
}

// ============================================================================
// Combat Phase API Functions
// These functions handle declare attackers, blockers, and damage assignment
// ============================================================================

/**
 * Declare a creature as an attacker targeting a specific defender
 * Per MTG Rule 508: Declaring attackers
 * @param gameId - The game ID
 * @param cardId - The attacking creature's ID
 * @param defenderId - The defender (player or planeswalker) being attacked
 */
export async function declareAttacker(
	gameId: string,
	cardId: string,
	defenderId: string
): Promise<void> {
	return sendPlayerString(gameId, `ATTACK:${cardId}:${defenderId}`);
}

/**
 * Finish declaring attackers and move to the next step
 * Per MTG Rule 508.1: All attackers are declared simultaneously
 */
export async function finishDeclaringAttackers(gameId: string): Promise<void> {
	return sendPlayerString(gameId, 'DONE_ATTACKING');
}

/**
 * Declare multiple attackers at once
 * Convenience function to declare all attackers in sequence
 * @param gameId - The game ID
 * @param attackers - Array of attacker declarations
 */
export async function declareAttackers(
	gameId: string,
	attackers: Array<{ cardId: string; defenderId: string }>
): Promise<void> {
	for (const attacker of attackers) {
		await declareAttacker(gameId, attacker.cardId, attacker.defenderId);
	}
	await finishDeclaringAttackers(gameId);
}

/**
 * Declare a creature as a blocker of a specific attacker
 * Per MTG Rule 509: Declaring blockers
 * @param gameId - The game ID
 * @param blockerId - The blocking creature's ID
 * @param attackerId - The attacking creature being blocked
 */
export async function declareBlocker(
	gameId: string,
	blockerId: string,
	attackerId: string
): Promise<void> {
	return sendPlayerString(gameId, `BLOCK:${blockerId}:${attackerId}`);
}

/**
 * Finish declaring blockers and move to the next step
 * Per MTG Rule 509.1: All blockers are declared simultaneously
 */
export async function finishDeclaringBlockers(gameId: string): Promise<void> {
	return sendPlayerString(gameId, 'DONE_BLOCKING');
}

/**
 * Declare multiple blockers at once
 * Convenience function to declare all blockers in sequence
 * @param gameId - The game ID
 * @param blockers - Array of blocker assignments
 */
export async function declareBlockers(
	gameId: string,
	blockers: Array<{ blockerId: string; attackerId: string }>
): Promise<void> {
	for (const blocker of blockers) {
		await declareBlocker(gameId, blocker.blockerId, blocker.attackerId);
	}
	await finishDeclaringBlockers(gameId);
}

/**
 * Assign combat damage from an attacker to its blockers and/or defending player
 * Per MTG Rule 510: Combat damage step
 * @param gameId - The game ID
 * @param assignments - Array of damage assignments (targetId + damage amount)
 */
export async function assignCombatDamage(
	gameId: string,
	assignments: Array<{ targetId: string; damage: number }>
): Promise<void> {
	// Format: DAMAGE:targetId:amount for each assignment
	for (const assignment of assignments) {
		await sendPlayerString(gameId, `DAMAGE:${assignment.targetId}:${assignment.damage}`);
	}
	await sendPlayerString(gameId, 'DONE_DAMAGE');
}

/**
 * Skip combat - declare no attackers
 * Convenience function to immediately finish declaring attackers with none selected
 */
export async function skipCombat(gameId: string): Promise<void> {
	return finishDeclaringAttackers(gameId);
}

/**
 * Decline to block - finish declaring blockers with none assigned
 * Convenience function to immediately finish declaring blockers with none selected
 */
export async function declineToBlock(gameId: string): Promise<void> {
	return finishDeclaringBlockers(gameId);
}

// ============================================================================
// Rollback API Functions
// These functions handle game state rollback requests and responses
// ============================================================================

import type {
	RequestRollbackRequest,
	RequestRollbackResponse,
	RespondToRollbackRequest,
	RespondToRollbackResponse,
	CancelRollbackRequest,
	CancelRollbackResponse
} from '$lib/generated/mage/v1/game';

/**
 * Request a rollback to a specific game log message
 * In multiplayer games, this sends a request to opponents for consent
 * @param gameId - The game ID
 * @param messageId - The message ID to rollback to
 * @returns Object containing success status, requestId (for multiplayer), and whether consent is required
 */
export async function requestRollback(
	gameId: string,
	messageId: number
): Promise<{ success: boolean; requestId: string; requiresConsent: boolean; error?: string }> {
	const client = getMageClient();
	const sessionId = await client.ensureSessionId();

	if (!sessionId) {
		throw new Error('No active session - please login first');
	}

	const request: RequestRollbackRequest = {
		sessionId,
		gameId,
		messageId
	};

	const response = await client.call<RequestRollbackRequest, RequestRollbackResponse>(
		'RequestRollback',
		request
	);

	return {
		success: response.success,
		requestId: response.requestId,
		requiresConsent: response.requiresConsent,
		error: response.error
	};
}

/**
 * Respond to a rollback request from another player
 * @param gameId - The game ID
 * @param requestId - The rollback request ID
 * @param approved - Whether to approve or deny the rollback
 */
export async function respondToRollback(
	gameId: string,
	requestId: string,
	approved: boolean
): Promise<{ success: boolean; error?: string }> {
	const client = getMageClient();
	const sessionId = await client.ensureSessionId();

	if (!sessionId) {
		throw new Error('No active session - please login first');
	}

	const request: RespondToRollbackRequest = {
		sessionId,
		gameId,
		requestId,
		approved
	};

	const response = await client.call<RespondToRollbackRequest, RespondToRollbackResponse>(
		'RespondToRollback',
		request
	);

	return {
		success: response.success,
		error: response.error
	};
}

/**
 * Cancel a pending rollback request
 * @param gameId - The game ID
 * @param requestId - The rollback request ID to cancel
 */
export async function cancelRollback(
	gameId: string,
	requestId: string
): Promise<{ success: boolean; error?: string }> {
	const client = getMageClient();
	const sessionId = await client.ensureSessionId();

	if (!sessionId) {
		throw new Error('No active session - please login first');
	}

	const request: CancelRollbackRequest = {
		sessionId,
		gameId,
		requestId
	};

	const response = await client.call<CancelRollbackRequest, CancelRollbackResponse>(
		'CancelRollback',
		request
	);

	return {
		success: response.success,
		error: response.error
	};
}
