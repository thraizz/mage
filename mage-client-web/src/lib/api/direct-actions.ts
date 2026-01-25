/**
 * Direct Actions API for Rules-Light Game Engine
 *
 * These functions allow direct manipulation of game state without rule enforcement.
 * They send string commands through the existing SendPlayerString RPC endpoint.
 */

import { sendPlayerString } from './game';

/**
 * Tap or untap a card
 * @param gameId - The game ID
 * @param cardId - The card to tap/untap
 * @param tapped - Whether the card should be tapped (true) or untapped (false)
 */
export async function tapUntap(gameId: string, cardId: string, tapped: boolean): Promise<void> {
  const command = tapped ? `TAP:${cardId}` : `UNTAP:${cardId}`;
  return sendPlayerString(gameId, command);
}

/**
 * Untap all permanents controlled by the current player
 * @param gameId - The game ID
 */
export async function untapAll(gameId: string): Promise<void> {
  return sendPlayerString(gameId, 'UNTAP_ALL');
}

/**
 * Flip a card face-up or face-down
 * @param gameId - The game ID
 * @param cardId - The card to flip
 * @param faceDown - Whether the card should be face-down
 */
export async function flipCard(gameId: string, cardId: string, faceDown: boolean): Promise<void> {
  return sendPlayerString(gameId, `FLIP:${cardId}:${faceDown}`);
}

/**
 * Transform a double-faced card
 * @param gameId - The game ID
 * @param cardId - The card to transform
 */
export async function transformCard(gameId: string, cardId: string): Promise<void> {
  return sendPlayerString(gameId, `TRANSFORM:${cardId}`);
}

/**
 * Move a card to a different zone
 * @param gameId - The game ID
 * @param cardId - The card to move
 * @param targetZone - The target zone (HAND, BATTLEFIELD, GRAVEYARD, EXILE, LIBRARY, COMMAND)
 */
export async function moveCard(gameId: string, cardId: string, targetZone: string): Promise<void> {
  return sendPlayerString(gameId, `MOVE:${cardId}:${targetZone}`);
}

/**
 * Set the counter value on a card
 * @param gameId - The game ID
 * @param cardId - The card to modify
 * @param counterType - The type of counter (e.g., "+1/+1", "loyalty", "charge")
 * @param amount - The amount to set
 */
export async function setCardCounter(
  gameId: string,
  cardId: string,
  counterType: string,
  amount: number
): Promise<void> {
  return sendPlayerString(gameId, `SET_COUNTER:${cardId}:${counterType}:${amount}`);
}

/**
 * Add or remove counters from a card
 * @param gameId - The game ID
 * @param cardId - The card to modify
 * @param counterType - The type of counter
 * @param delta - The amount to add (positive) or remove (negative)
 */
export async function modifyCardCounter(
  gameId: string,
  cardId: string,
  counterType: string,
  delta: number
): Promise<void> {
  return sendPlayerString(gameId, `MODIFY_COUNTER:${cardId}:${counterType}:${delta}`);
}

/**
 * Create a token on the battlefield
 * @param gameId - The game ID
 * @param name - Token name (e.g., "Soldier", "Zombie")
 * @param types - Token types (e.g., "Creature — Soldier", "Artifact — Treasure")
 * @param power - Token power (can be "*" for variable)
 * @param toughness - Token toughness (can be "*" for variable)
 * @param color - Token color (white, blue, black, red, green, colorless, multicolor)
 * @param abilities - Array of ability text
 * @param _count - Number of tokens to create (not yet implemented server-side)
 */
export async function createToken(
  gameId: string,
  name: string,
  types: string,
  power: string,
  toughness: string,
  color: string,
  abilities: string[],
  _count: number = 1
): Promise<{ tokenId: string }> {
  const abilitiesStr = abilities.join(',');
  await sendPlayerString(
    gameId,
    `CREATE_TOKEN:${name}:${types}:${power}:${toughness}:${color}:${abilitiesStr}`
  );
  // Note: Server doesn't return token ID yet, but UI can refresh state
  return { tokenId: 'pending' };
}

/**
 * Destroy a token (remove it from the game)
 * @param gameId - The game ID
 * @param cardId - The token to destroy
 */
export async function destroyToken(gameId: string, cardId: string): Promise<void> {
  return sendPlayerString(gameId, `DESTROY_TOKEN:${cardId}`);
}

/**
 * Set a player's life total directly
 * @param gameId - The game ID
 * @param playerId - The player whose life to set
 * @param amount - The new life total
 */
export async function setPlayerLife(
  gameId: string,
  playerId: string,
  amount: number
): Promise<void> {
  return sendPlayerString(gameId, `SET_LIFE:${playerId}:${amount}`);
}

/**
 * Modify a player's life total
 * @param gameId - The game ID
 * @param playerId - The player whose life to modify
 * @param delta - The amount to add (positive) or remove (negative)
 */
export async function modifyPlayerLife(
  gameId: string,
  playerId: string,
  delta: number
): Promise<void> {
  return sendPlayerString(gameId, `MODIFY_LIFE:${playerId}:${delta}`);
}

/**
 * Draw cards for a player
 * @param gameId - The game ID
 * @param playerId - The player to draw cards (optional, defaults to current player)
 * @param count - Number of cards to draw
 */
export async function drawCards(gameId: string, playerId: string, count: number): Promise<void> {
  return sendPlayerString(gameId, `DRAW:${playerId}:${count}`);
}

/**
 * Mill cards (move top N cards from library to graveyard)
 * @param gameId - The game ID
 * @param playerId - The player whose library to mill
 * @param count - Number of cards to mill
 */
export async function millCards(gameId: string, playerId: string, count: number): Promise<void> {
  return sendPlayerString(gameId, `MILL:${playerId}:${count}`);
}

/**
 * Scry N cards (look at top N cards and rearrange them)
 * Note: This is a simplified version that just initiates a scry.
 * Full scry UI with card selection would require additional commands.
 * @param gameId - The game ID
 * @param playerId - The player who is scrying
 * @param count - Number of cards to scry
 */
export async function scryCards(gameId: string, playerId: string, count: number): Promise<void> {
  return sendPlayerString(gameId, `SCRY:${playerId}:${count}`);
}

/**
 * Set whether the top card of library is revealed to all players
 * @param gameId - The game ID
 * @param playerId - The player whose library top card to reveal/hide
 * @param revealed - Whether the top card should be revealed
 */
export async function setRevealedTop(
  gameId: string,
  playerId: string,
  revealed: boolean
): Promise<void> {
  return sendPlayerString(gameId, `REVEAL_TOP:${playerId}:${revealed}`);
}

/**
 * Mulligan (shuffle hand into library and draw N-1 cards)
 * @param gameId - The game ID
 * @param playerId - The player who is mulliganing
 */
export async function mulligan(gameId: string, playerId: string): Promise<void> {
  return sendPlayerString(gameId, `MULLIGAN:${playerId}`);
}

/**
 * Keep hand (end mulligan phase for this player)
 * @param gameId - The game ID
 * @param playerId - The player who is keeping their hand
 */
export async function keepHand(gameId: string, playerId: string): Promise<void> {
  return sendPlayerString(gameId, `KEEP_HAND:${playerId}`);
}

/**
 * Alias for modifyPlayerLife for convenience
 * @param gameId - The game ID
 * @param playerId - The player whose life to modify
 * @param delta - The amount to add (positive) or remove (negative)
 */
export async function modifyLife(gameId: string, playerId: string, delta: number): Promise<void> {
  return modifyPlayerLife(gameId, playerId, delta);
}

/**
 * Set a player counter (e.g., poison, energy, experience)
 * @param gameId - The game ID
 * @param playerId - The player whose counter to set
 * @param counterType - The type of counter (e.g., "poison", "energy")
 * @param amount - The amount to set
 */
export async function setPlayerCounter(
  gameId: string,
  playerId: string,
  counterType: string,
  amount: number
): Promise<void> {
  return sendPlayerString(gameId, `SET_PLAYER_COUNTER:${playerId}:${counterType}:${amount}`);
}

/**
 * Shuffle a player's library
 * @param gameId - The game ID
 * @param playerId - The player whose library to shuffle (optional, defaults to current player)
 */
export async function shuffleLibrary(gameId: string, playerId?: string): Promise<void> {
  const command = playerId ? `SHUFFLE:${playerId}` : 'SHUFFLE';
  return sendPlayerString(gameId, command);
}

/**
 * Advance to the next turn
 * @param gameId - The game ID
 */
export async function nextTurn(gameId: string): Promise<void> {
  return sendPlayerString(gameId, 'NEXT_TURN');
}

/**
 * Clear combat (remove all attackers and blockers)
 * @param gameId - The game ID
 */
export async function clearCombat(gameId: string): Promise<void> {
  return sendPlayerString(gameId, 'CLEAR_COMBAT');
}

/**
 * Search your library for a card
 * @param gameId - The game ID
 * @param destination - Where to put the found card: "hand", "battlefield", "top", "graveyard"
 * @param shuffle - Whether to shuffle the library after searching (default: true)
 * @param message - Optional message describing what to search for
 */
export async function searchLibrary(
  gameId: string,
  destination: 'hand' | 'battlefield' | 'top' | 'graveyard' = 'hand',
  shuffle: boolean = true,
  message?: string
): Promise<void> {
  const msgPart = message ? `:${message}` : '';
  return sendPlayerString(gameId, `SEARCH_LIBRARY:${destination}:${shuffle}${msgPart}`);
}

/**
 * Select a card from a library search
 * @param gameId - The game ID
 * @param cardId - The ID of the card to select (or "CANCEL" to cancel the search)
 */
export async function selectLibraryCard(gameId: string, cardId: string): Promise<void> {
  return sendPlayerString(gameId, cardId);
}

/**
 * Add a card to the visual stack without moving it from its current zone.
 * This is for rules-light manual tracking - the card stays where it is
 * but appears in the stack for all players to see.
 * @param gameId - The game ID
 * @param cardId - The card to add to the stack
 */
export async function addToStack(gameId: string, cardId: string): Promise<void> {
  return sendPlayerString(gameId, `STACK_ADD:${cardId}`);
}

/**
 * Remove an item from the stack (for manual resolution tracking)
 * @param gameId - The game ID
 * @param itemId - The stack item ID to remove
 */
export async function removeFromStack(gameId: string, itemId: string): Promise<void> {
  return sendPlayerString(gameId, `STACK_REMOVE:${itemId}`);
}
