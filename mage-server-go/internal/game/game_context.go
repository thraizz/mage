package game

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/counters"
	manapool "github.com/magefree/mage-server-go/internal/game/mana"
	"github.com/magefree/mage-server-go/internal/game/token"
	"go.uber.org/zap"
)

// GameContext implements the abilities.TokenGameContext interface.
// This provides the game engine implementation that token and counter effects need.
type GameContext struct {
	gameID uuid.UUID
	engine *MageEngine
	logger *zap.Logger
}

// NewGameContext creates a new game context for the given game.
func NewGameContext(gameID uuid.UUID, engine *MageEngine, logger *zap.Logger) *GameContext {
	return &GameContext{
		gameID: gameID,
		engine: engine,
		logger: logger,
	}
}

// GetGameID returns the game ID for this context.
func (gc *GameContext) GetGameID() uuid.UUID {
	return gc.gameID
}

// ==============================================================================
// abilities.GameContext interface implementation
// ==============================================================================

// GetCard retrieves a card by ID.
func (gc *GameContext) GetCard(id uuid.UUID) (interface{}, error) {
	gc.engine.mu.RLock()
	gameState, ok := gc.engine.games[gc.gameID.String()]
	gc.engine.mu.RUnlock()

	if !ok {
		return nil, fmt.Errorf("game %s not found", gc.gameID)
	}

	gameState.mu.RLock()
	defer gameState.mu.RUnlock()

	// Check battlefield (shared zone)
	for _, permanent := range gameState.battlefield {
		if permanent.ID == id.String() {
			return permanent, nil
		}
	}

	// Search through player zones
	for _, player := range gameState.players {
		// Check hand
		for _, card := range player.Hand {
			if card.ID == id.String() {
				return card, nil
			}
		}

		// Check graveyard
		for _, card := range player.Graveyard {
			if card.ID == id.String() {
				return card, nil
			}
		}

		// Check library
		for _, card := range player.Library {
			if card.ID == id.String() {
				return card, nil
			}
		}
	}

	// Check exile (shared zone)
	for _, card := range gameState.exile {
		if card.ID == id.String() {
			return card, nil
		}
	}

	// Check command zone (shared zone)
	for _, card := range gameState.command {
		if card.ID == id.String() {
			return card, nil
		}
	}

	return nil, fmt.Errorf("card %s not found", id)
}

// GetPlayer retrieves a player by ID.
func (gc *GameContext) GetPlayer(id uuid.UUID) (interface{}, error) {
	gc.engine.mu.RLock()
	gameState, ok := gc.engine.games[gc.gameID.String()]
	gc.engine.mu.RUnlock()

	if !ok {
		return nil, fmt.Errorf("game %s not found", gc.gameID)
	}

	gameState.mu.RLock()
	defer gameState.mu.RUnlock()

	player, ok := gameState.players[id.String()]
	if !ok {
		return nil, fmt.Errorf("player %s not found", id)
	}

	return player, nil
}

// DealDamage deals damage from a source to a target.
func (gc *GameContext) DealDamage(sourceID, targetID uuid.UUID, amount int) error {
	gc.engine.mu.RLock()
	gameState, ok := gc.engine.games[gc.gameID.String()]
	gc.engine.mu.RUnlock()

	if !ok {
		return fmt.Errorf("game %s not found", gc.gameID)
	}

	gameState.mu.Lock()
	defer gameState.mu.Unlock()

	// Try to find target as a player
	if player, ok := gameState.players[targetID.String()]; ok {
		player.Life -= amount
		gc.logger.Info("dealt damage to player",
			zap.String("source", sourceID.String()),
			zap.String("target", targetID.String()),
			zap.Int("amount", amount),
			zap.Int("new_life", player.Life))
		return nil
	}

	// Try to find target as a permanent
	for _, permanent := range gameState.battlefield {
		if permanent.ID == targetID.String() {
			permanent.Damage += amount
			if permanent.DamageSources == nil {
				permanent.DamageSources = make(map[string]int)
			}
			permanent.DamageSources[sourceID.String()] += amount

			gc.logger.Info("dealt damage to permanent",
				zap.String("source", sourceID.String()),
				zap.String("target", targetID.String()),
				zap.String("permanent", permanent.Name),
				zap.Int("amount", amount),
				zap.Int("total_damage", permanent.Damage))
			return nil
		}
	}

	return fmt.Errorf("target %s not found", targetID)
}

// DrawCards draws the specified number of cards for a player.
func (gc *GameContext) DrawCards(playerID uuid.UUID, amount int) error {
	gc.engine.mu.RLock()
	gameState, ok := gc.engine.games[gc.gameID.String()]
	gc.engine.mu.RUnlock()

	if !ok {
		return fmt.Errorf("game %s not found", gc.gameID)
	}

	gameState.mu.Lock()
	defer gameState.mu.Unlock()

	player, ok := gameState.players[playerID.String()]
	if !ok {
		return fmt.Errorf("player %s not found", playerID)
	}

	for i := 0; i < amount; i++ {
		if len(player.Library) == 0 {
			// Player loses if they can't draw from empty library
			player.Lost = true
			gc.logger.Info("player loses by drawing from empty library",
				zap.String("player", playerID.String()))
			return fmt.Errorf("cannot draw from empty library")
		}

		// Draw top card of library
		card := player.Library[len(player.Library)-1]
		player.Library = player.Library[:len(player.Library)-1]
		player.Hand = append(player.Hand, card)

		gc.logger.Debug("player drew card",
			zap.String("player", playerID.String()),
			zap.String("card", card.Name))
	}

	gc.logger.Info("player drew cards",
		zap.String("player", playerID.String()),
		zap.Int("amount", amount))

	return nil
}

// DestroyPermanent destroys a permanent.
func (gc *GameContext) DestroyPermanent(permanentID uuid.UUID) error {
	gc.engine.mu.RLock()
	gameState, ok := gc.engine.games[gc.gameID.String()]
	gc.engine.mu.RUnlock()

	if !ok {
		return fmt.Errorf("game %s not found", gc.gameID)
	}

	gameState.mu.Lock()
	defer gameState.mu.Unlock()

	// Find and remove the permanent from the battlefield
	for i, permanent := range gameState.battlefield {
		if permanent.ID == permanentID.String() {
			// Store LKI BEFORE removing from battlefield
			// This captures the permanent's state (including counters) at the moment of leaving
			gameState.storeLKI(permanent)

			// Remove from battlefield
			gameState.battlefield = append(gameState.battlefield[:i], gameState.battlefield[i+1:]...)

			// Add to graveyard (tokens don't go to graveyard)
			if !permanent.IsToken {
				// Find the owner and add to their graveyard
				if owner, ok := gameState.players[permanent.OwnerID]; ok {
					owner.Graveyard = append(owner.Graveyard, permanent)
					gc.logger.Info("permanent destroyed and sent to graveyard",
						zap.String("permanent", permanent.Name),
						zap.String("owner", permanent.OwnerID))
				}
			} else {
				gc.logger.Info("token destroyed and removed from game",
					zap.String("token", permanent.Name),
					zap.String("owner", permanent.OwnerID))
			}

			return nil
		}
	}

	return fmt.Errorf("permanent %s not found", permanentID)
}

// AddMana adds mana to a player's mana pool.
func (gc *GameContext) AddMana(playerID uuid.UUID, mana *abilities.Mana) error {
	gc.engine.mu.RLock()
	gameState, ok := gc.engine.games[gc.gameID.String()]
	gc.engine.mu.RUnlock()

	if !ok {
		return fmt.Errorf("game %s not found", gc.gameID)
	}

	gameState.mu.Lock()
	defer gameState.mu.Unlock()

	player, ok := gameState.players[playerID.String()]
	if !ok {
		return fmt.Errorf("player %s not found", playerID)
	}

	// Add mana to pool - need to import mana package at top
	if mana.White > 0 {
		player.ManaPool.Add(manapool.ManaWhite, mana.White)
	}
	if mana.Blue > 0 {
		player.ManaPool.Add(manapool.ManaBlue, mana.Blue)
	}
	if mana.Black > 0 {
		player.ManaPool.Add(manapool.ManaBlack, mana.Black)
	}
	if mana.Red > 0 {
		player.ManaPool.Add(manapool.ManaRed, mana.Red)
	}
	if mana.Green > 0 {
		player.ManaPool.Add(manapool.ManaGreen, mana.Green)
	}
	if mana.Colorless > 0 {
		player.ManaPool.Add(manapool.ManaColorless, mana.Colorless)
	}

	gc.logger.Debug("added mana to pool",
		zap.String("player", playerID.String()),
		zap.Int("white", mana.White),
		zap.Int("blue", mana.Blue),
		zap.Int("black", mana.Black),
		zap.Int("red", mana.Red),
		zap.Int("green", mana.Green),
		zap.Int("colorless", mana.Colorless))

	return nil
}

// GetManaPool returns a player's mana pool for cost payment.
func (gc *GameContext) GetManaPool(playerID uuid.UUID) abilities.ManaPoolInterface {
	gc.engine.mu.RLock()
	gameState, ok := gc.engine.games[gc.gameID.String()]
	gc.engine.mu.RUnlock()

	if !ok {
		gc.logger.Error("game not found when getting mana pool",
			zap.String("game_id", gc.gameID.String()))
		return nil
	}

	gameState.mu.RLock()
	defer gameState.mu.RUnlock()

	player, ok := gameState.players[playerID.String()]
	if !ok {
		gc.logger.Error("player not found when getting mana pool",
			zap.String("player_id", playerID.String()))
		return nil
	}

	// Return the player's mana pool which implements ManaPoolInterface
	return player.ManaPool
}

// TapPermanent taps a permanent.
func (gc *GameContext) TapPermanent(permanentID uuid.UUID) error {
	gc.engine.mu.RLock()
	gameState, ok := gc.engine.games[gc.gameID.String()]
	gc.engine.mu.RUnlock()

	if !ok {
		return fmt.Errorf("game %s not found", gc.gameID)
	}

	gameState.mu.Lock()
	defer gameState.mu.Unlock()

	// Find the permanent
	for _, permanent := range gameState.battlefield {
		if permanent.ID == permanentID.String() {
			permanent.Tapped = true
			gc.logger.Debug("permanent tapped",
				zap.String("permanent", permanent.Name))
			return nil
		}
	}

	return fmt.Errorf("permanent %s not found", permanentID)
}

// UntapPermanent untaps a permanent.
func (gc *GameContext) UntapPermanent(permanentID uuid.UUID) error {
	gc.engine.mu.RLock()
	gameState, ok := gc.engine.games[gc.gameID.String()]
	gc.engine.mu.RUnlock()

	if !ok {
		return fmt.Errorf("game %s not found", gc.gameID)
	}

	gameState.mu.Lock()
	defer gameState.mu.Unlock()

	// Find the permanent
	for _, permanent := range gameState.battlefield {
		if permanent.ID == permanentID.String() {
			permanent.Tapped = false
			gc.logger.Debug("permanent untapped",
				zap.String("permanent", permanent.Name))
			return nil
		}
	}

	return fmt.Errorf("permanent %s not found", permanentID)
}

// IsPermanentTapped checks if a permanent is tapped.
func (gc *GameContext) IsPermanentTapped(permanentID uuid.UUID) bool {
	gc.engine.mu.RLock()
	gameState, ok := gc.engine.games[gc.gameID.String()]
	gc.engine.mu.RUnlock()

	if !ok {
		return false
	}

	gameState.mu.RLock()
	defer gameState.mu.RUnlock()

	// Find the permanent
	for _, permanent := range gameState.battlefield {
		if permanent.ID == permanentID.String() {
			return permanent.Tapped
		}
	}

	return false
}

// SacrificePermanent sacrifices a permanent (moves to graveyard, triggers dies events).
func (gc *GameContext) SacrificePermanent(permanentID uuid.UUID) error {
	gc.engine.mu.RLock()
	gameState, ok := gc.engine.games[gc.gameID.String()]
	gc.engine.mu.RUnlock()

	if !ok {
		return fmt.Errorf("game %s not found", gc.gameID)
	}

	gameState.mu.Lock()
	defer gameState.mu.Unlock()

	// Find and remove the permanent from the battlefield
	for i, permanent := range gameState.battlefield {
		if permanent.ID == permanentID.String() {
			// Store LKI BEFORE removing from battlefield
			// This captures the permanent's state (including counters) at the moment of leaving
			gameState.storeLKI(permanent)

			// Remove from battlefield
			gameState.battlefield = append(gameState.battlefield[:i], gameState.battlefield[i+1:]...)

			// Add to graveyard (tokens don't go to graveyard)
			if !permanent.IsToken {
				// Find the owner and add to their graveyard
				if owner, ok := gameState.players[permanent.OwnerID]; ok {
					owner.Graveyard = append(owner.Graveyard, permanent)
					gc.logger.Info("permanent sacrificed and sent to graveyard",
						zap.String("permanent", permanent.Name),
						zap.String("owner", permanent.OwnerID))
				}
			} else {
				gc.logger.Info("token sacrificed and removed from game",
					zap.String("token", permanent.Name),
					zap.String("owner", permanent.OwnerID))
			}

			// TODO: Trigger dies events for creatures/planeswalkers
			return nil
		}
	}

	return fmt.Errorf("permanent %s not found", permanentID)
}

// DiscardCard discards a card from a player's hand.
func (gc *GameContext) DiscardCard(playerID uuid.UUID, cardID uuid.UUID) error {
	gc.engine.mu.RLock()
	gameState, ok := gc.engine.games[gc.gameID.String()]
	gc.engine.mu.RUnlock()

	if !ok {
		return fmt.Errorf("game %s not found", gc.gameID)
	}

	gameState.mu.Lock()
	defer gameState.mu.Unlock()

	player, ok := gameState.players[playerID.String()]
	if !ok {
		return fmt.Errorf("player %s not found", playerID)
	}

	// Find and remove the card from hand
	for i, card := range player.Hand {
		if card.ID == cardID.String() {
			// Remove from hand
			player.Hand = append(player.Hand[:i], player.Hand[i+1:]...)

			// Add to graveyard
			player.Graveyard = append(player.Graveyard, card)

			gc.logger.Info("card discarded",
				zap.String("card", card.Name),
				zap.String("player", playerID.String()))

			// TODO: Trigger discard events
			return nil
		}
	}

	return fmt.Errorf("card %s not found in player's hand", cardID)
}

// GetPlayerHand returns the cards in a player's hand.
func (gc *GameContext) GetPlayerHand(playerID uuid.UUID) ([]interface{}, error) {
	gc.engine.mu.RLock()
	gameState, ok := gc.engine.games[gc.gameID.String()]
	gc.engine.mu.RUnlock()

	if !ok {
		return nil, fmt.Errorf("game %s not found", gc.gameID)
	}

	gameState.mu.RLock()
	defer gameState.mu.RUnlock()

	player, ok := gameState.players[playerID.String()]
	if !ok {
		return nil, fmt.Errorf("player %s not found", playerID)
	}

	// Convert to []interface{}
	hand := make([]interface{}, len(player.Hand))
	for i, card := range player.Hand {
		hand[i] = card
	}

	return hand, nil
}

// GetPermanentsControlledByPlayer returns all permanents controlled by a player.
func (gc *GameContext) GetPermanentsControlledByPlayer(playerID uuid.UUID) ([]interface{}, error) {
	gc.engine.mu.RLock()
	gameState, ok := gc.engine.games[gc.gameID.String()]
	gc.engine.mu.RUnlock()

	if !ok {
		return nil, fmt.Errorf("game %s not found", gc.gameID)
	}

	gameState.mu.RLock()
	defer gameState.mu.RUnlock()

	permanents := make([]interface{}, 0)
	for _, permanent := range gameState.battlefield {
		if permanent.ControllerID == playerID.String() {
			permanents = append(permanents, permanent)
		}
	}

	return permanents, nil
}

// ==============================================================================
// abilities.CounterGameContext interface implementation
// ==============================================================================

// GetPermanent retrieves a permanent by ID.
func (gc *GameContext) GetPermanent(id uuid.UUID) (interface{}, error) {
	gc.engine.mu.RLock()
	gameState, ok := gc.engine.games[gc.gameID.String()]
	gc.engine.mu.RUnlock()

	if !ok {
		return nil, fmt.Errorf("game %s not found", gc.gameID)
	}

	gameState.mu.RLock()
	defer gameState.mu.RUnlock()

	for _, permanent := range gameState.battlefield {
		if permanent.ID == id.String() {
			return permanent, nil
		}
	}

	return nil, fmt.Errorf("permanent %s not found", id)
}

// GetPermanentOrLKI retrieves a permanent by ID, falling back to Last Known Information
// if the permanent is no longer on the battlefield.
// Java: Game.getPermanentOrLKIBattlefield()
// MTG Rules: 113.7a (Objects track their last known information when they leave zones)
//
// This is critical for triggered abilities like Resourceful Defense that need to
// know the state of a permanent (especially its counters) at the moment it left
// the battlefield.
func (gc *GameContext) GetPermanentOrLKI(id uuid.UUID) (*LastKnownInfo, error) {
	gc.engine.mu.RLock()
	gameState, ok := gc.engine.games[gc.gameID.String()]
	gc.engine.mu.RUnlock()

	if !ok {
		return nil, fmt.Errorf("game %s not found", gc.gameID)
	}

	gameState.mu.RLock()
	defer gameState.mu.RUnlock()

	// First try to find the permanent on the battlefield
	for _, permanent := range gameState.battlefield {
		if permanent.ID == id.String() {
			// Permanent still exists - create LKI from current state
			zoneCounter := gameState.lkiZoneCounter[permanent.ID]
			return createLKIFromCard(permanent, zoneCounter), nil
		}
	}

	// Permanent not on battlefield - check LKI
	if lki := gameState.getLKI(id.String()); lki != nil {
		return lki, nil
	}

	return nil, fmt.Errorf("permanent %s not found (no LKI available)", id)
}

// AddCountersToPermanent adds counters to a permanent.
func (gc *GameContext) AddCountersToPermanent(permanent interface{}, counter *counters.Counter) error {
	perm, ok := permanent.(*internalCard)
	if !ok {
		return fmt.Errorf("invalid permanent type")
	}

	gc.engine.mu.RLock()
	gameState, ok := gc.engine.games[gc.gameID.String()]
	gc.engine.mu.RUnlock()

	if !ok {
		return fmt.Errorf("game %s not found", gc.gameID)
	}

	gameState.mu.Lock()
	defer gameState.mu.Unlock()

	// Find the permanent and add counters
	for _, p := range gameState.battlefield {
		if p.ID == perm.ID {
			p.Counters.AddCounter(counter)

			gc.logger.Info("added counters to permanent",
				zap.String("permanent", p.Name),
				zap.String("counter", counter.Name),
				zap.Int("amount", counter.Count),
				zap.Int("total", p.Counters.GetCount(counter.Name)))
			return nil
		}
	}

	return fmt.Errorf("permanent %s not found", perm.ID)
}

// AddCountersToPlayer adds counters to a player.
func (gc *GameContext) AddCountersToPlayer(player interface{}, counter *counters.Counter) error {
	p, ok := player.(*internalPlayer)
	if !ok {
		return fmt.Errorf("invalid player type")
	}

	gc.engine.mu.RLock()
	gameState, ok := gc.engine.games[gc.gameID.String()]
	gc.engine.mu.RUnlock()

	if !ok {
		return fmt.Errorf("game %s not found", gc.gameID)
	}

	gameState.mu.Lock()
	defer gameState.mu.Unlock()

	playerState, ok := gameState.players[p.PlayerID]
	if !ok {
		return fmt.Errorf("player %s not found", p.PlayerID)
	}

	// Handle specific counter types that players track
	switch counter.Name {
	case string(counters.CounterTypePoison):
		playerState.Poison += counter.Count
		gc.logger.Info("added poison counters to player",
			zap.String("player", p.PlayerID),
			zap.Int("amount", counter.Count),
			zap.Int("total", playerState.Poison))
	case string(counters.CounterTypeEnergy):
		playerState.Energy += counter.Count
		gc.logger.Info("added energy counters to player",
			zap.String("player", p.PlayerID),
			zap.Int("amount", counter.Count),
			zap.Int("total", playerState.Energy))
	default:
		// For other counter types, we'd need a map. For now, log a warning
		gc.logger.Warn("unsupported player counter type",
			zap.String("counter", counter.Name),
			zap.String("player", p.PlayerID))
	}

	return nil
}

// AddCountersToCard adds counters to a card (in any zone).
func (gc *GameContext) AddCountersToCard(card interface{}, counter *counters.Counter) error {
	c, ok := card.(*internalCard)
	if !ok {
		return fmt.Errorf("invalid card type")
	}

	gc.engine.mu.RLock()
	gameState, ok := gc.engine.games[gc.gameID.String()]
	gc.engine.mu.RUnlock()

	if !ok {
		return fmt.Errorf("game %s not found", gc.gameID)
	}

	gameState.mu.Lock()
	defer gameState.mu.Unlock()

	// Check battlefield first
	for _, zoneCard := range gameState.battlefield {
		if zoneCard.ID == c.ID {
			zoneCard.Counters.AddCounter(counter)

			gc.logger.Info("added counters to card on battlefield",
				zap.String("card", zoneCard.Name),
				zap.String("counter", counter.Name),
				zap.Int("amount", counter.Count),
				zap.Int("total", zoneCard.Counters.GetCount(counter.Name)))
			return nil
		}
	}

	// Search player zones
	for _, player := range gameState.players {
		// Check all zones
		allCards := [][]*internalCard{
			player.Hand,
			player.Graveyard,
			player.Library,
		}

		for _, zone := range allCards {
			for _, zoneCard := range zone {
				if zoneCard.ID == c.ID {
					zoneCard.Counters.AddCounter(counter)

					gc.logger.Info("added counters to card",
						zap.String("card", zoneCard.Name),
						zap.String("counter", counter.Name),
						zap.Int("amount", counter.Count),
						zap.Int("total", zoneCard.Counters.GetCount(counter.Name)))
					return nil
				}
			}
		}
	}

	// Check shared zones
	for _, zoneCard := range gameState.exile {
		if zoneCard.ID == c.ID {
			zoneCard.Counters.AddCounter(counter)
			gc.logger.Info("added counters to exiled card",
				zap.String("card", zoneCard.Name),
				zap.String("counter", counter.Name),
				zap.Int("amount", counter.Count),
				zap.Int("total", zoneCard.Counters.GetCount(counter.Name)))
			return nil
		}
	}

	for _, zoneCard := range gameState.command {
		if zoneCard.ID == c.ID {
			zoneCard.Counters.AddCounter(counter)
			gc.logger.Info("added counters to card in command zone",
				zap.String("card", zoneCard.Name),
				zap.String("counter", counter.Name),
				zap.Int("amount", counter.Count),
				zap.Int("total", zoneCard.Counters.GetCount(counter.Name)))
			return nil
		}
	}

	return fmt.Errorf("card %s not found", c.ID)
}

// GetAllPermanents returns all permanents on the battlefield.
func (gc *GameContext) GetAllPermanents() ([]interface{}, error) {
	gc.engine.mu.RLock()
	gameState, ok := gc.engine.games[gc.gameID.String()]
	gc.engine.mu.RUnlock()

	if !ok {
		return nil, fmt.Errorf("game %s not found", gc.gameID)
	}

	gameState.mu.RLock()
	defer gameState.mu.RUnlock()

	permanents := make([]interface{}, 0, len(gameState.battlefield))
	for _, permanent := range gameState.battlefield {
		permanents = append(permanents, permanent)
	}

	return permanents, nil
}

// InformPlayers sends an informational message to all players.
func (gc *GameContext) InformPlayers(message string) {
	gc.logger.Info("game message", zap.String("message", message))

	gc.engine.mu.RLock()
	gameState, ok := gc.engine.games[gc.gameID.String()]
	gc.engine.mu.RUnlock()

	if !ok {
		return
	}

	gameState.mu.Lock()
	defer gameState.mu.Unlock()

	// Add message to game messages
	gameState.messages = append(gameState.messages, EngineMessage{
		Text:      message,
		Color:     "",
		Timestamp: gameState.startedAt,
	})
}

// ==============================================================================
// abilities.TokenGameContext interface implementation
// ==============================================================================

// CreateTokens creates the specified number of tokens on the battlefield.
// Returns the UUIDs of the created permanents.
func (gc *GameContext) CreateTokens(tok *token.Token, amount int, source uuid.UUID, tapped, attacking bool) ([]uuid.UUID, error) {
	gc.engine.mu.RLock()
	gameState, ok := gc.engine.games[gc.gameID.String()]
	gc.engine.mu.RUnlock()

	if !ok {
		return nil, fmt.Errorf("game %s not found", gc.gameID)
	}

	gameState.mu.Lock()
	defer gameState.mu.Unlock()

	// Find the controller of the source card
	var controllerID string
	sourceFound := false

	// Check battlefield first
	for _, permanent := range gameState.battlefield {
		if permanent.ID == source.String() {
			controllerID = permanent.ControllerID
			sourceFound = true
			break
		}
	}

	// If not found on battlefield, check player zones
	if !sourceFound {
		for playerID, player := range gameState.players {
			// Also check other zones
			allCards := [][]*internalCard{player.Hand, player.Graveyard}
			for _, zone := range allCards {
				for _, card := range zone {
					if card.ID == source.String() {
						controllerID = playerID
						sourceFound = true
						break
					}
				}
				if sourceFound {
					break
				}
			}
			if sourceFound {
				break
			}
		}
	}

	// Check command zone
	if !sourceFound {
		for _, card := range gameState.command {
			if card.ID == source.String() {
				controllerID = card.ControllerID
				sourceFound = true
				break
			}
		}
	}

	// If source not found, use first player as controller
	if !sourceFound {
		if len(gameState.playerOrder) > 0 {
			controllerID = gameState.playerOrder[0]
		} else {
			return nil, fmt.Errorf("no players in game")
		}
	}

	if _, ok = gameState.players[controllerID]; !ok {
		return nil, fmt.Errorf("controller %s not found", controllerID)
	}

	createdIDs := make([]uuid.UUID, 0, amount)

	for i := 0; i < amount; i++ {
		// Create a copy of the token for each instance
		tokenCopy := tok.Copy()
		tokenID := uuid.New()

		// Convert token to internal card
		permanent := tokenToInternalCard(tokenCopy, tokenID, controllerID)
		permanent.IsToken = true
		permanent.Tapped = tapped
		permanent.Attacking = attacking

		// Handle summoning sickness (tokens have summoning sickness unless they have haste)
		hasHaste := false
		for _, ability := range tokenCopy.Abilities {
			if ability == "haste" {
				hasHaste = true
				break
			}
		}
		if !hasHaste {
			permanent.SummoningSickness = true
		}

		// Add to battlefield
		gameState.battlefield = append(gameState.battlefield, permanent)
		createdIDs = append(createdIDs, tokenID)

		gc.logger.Info("created token",
			zap.String("token", tokenCopy.Name),
			zap.String("controller", controllerID),
			zap.String("id", tokenID.String()),
			zap.Bool("tapped", tapped),
			zap.Bool("attacking", attacking))
	}

	// Update the token's LastAddedIDs
	tok.LastAddedIDs = createdIDs

	// Inform players
	gc.InformPlayers(fmt.Sprintf("Created %d %s token(s)", amount, tok.Name))

	return createdIDs, nil
}

// tokenToInternalCard converts a token.Token to an internalCard.
func tokenToInternalCard(tok *token.Token, id uuid.UUID, controllerID string) *internalCard {
	// Convert card types
	cardTypes := make([]string, 0, len(tok.CardTypes))
	for _, ct := range tok.CardTypes {
		cardTypes = append(cardTypes, cardTypeToString(ct))
	}

	// Convert color
	colorStr := ""
	if tok.Color.White {
		colorStr += "W"
	}
	if tok.Color.Blue {
		colorStr += "U"
	}
	if tok.Color.Black {
		colorStr += "B"
	}
	if tok.Color.Red {
		colorStr += "R"
	}
	if tok.Color.Green {
		colorStr += "G"
	}
	if colorStr == "" {
		colorStr = "C" // Colorless
	}

	// Convert abilities to EngineAbilityView
	abilities := make([]EngineAbilityView, 0, len(tok.Abilities))
	for _, ability := range tok.Abilities {
		abilities = append(abilities, EngineAbilityView{
			ID:   uuid.New().String(),
			Text: ability,
			Rule: ability,
		})
	}

	return &internalCard{
		ID:            id.String(),
		Name:          tok.Name,
		DisplayName:   tok.Name,
		Type:          joinStrings(cardTypes),
		SubTypes:      tok.Subtypes,
		SuperTypes:    []string{},
		Color:         colorStr,
		Power:         fmt.Sprintf("%d", tok.Power),
		Toughness:     fmt.Sprintf("%d", tok.Toughness),
		OwnerID:       controllerID,
		ControllerID:  controllerID,
		Abilities:     abilities,
		Counters:      counters.NewCounters(),
		IsToken:       false, // Will be set to true by caller
		Damage:        0,
		DamageSources: make(map[string]int),
	}
}

// cardTypeToString converts a token.CardType to a string.
func cardTypeToString(ct token.CardType) string {
	switch ct {
	case token.CardTypeArtifact:
		return "ARTIFACT"
	case token.CardTypeCreature:
		return "CREATURE"
	case token.CardTypeEnchantment:
		return "ENCHANTMENT"
	case token.CardTypeLand:
		return "LAND"
	case token.CardTypePlaneswalker:
		return "PLANESWALKER"
	case token.CardTypeInstant:
		return "INSTANT"
	case token.CardTypeSorcery:
		return "SORCERY"
	default:
		return ""
	}
}

// joinStrings joins a slice of strings with spaces.
func joinStrings(strs []string) string {
	result := ""
	for i, s := range strs {
		if i > 0 {
			result += " "
		}
		result += s
	}
	return result
}

// ==============================================================================
// CDA (Characteristic-Defining Ability) Support - Rule 604
// ==============================================================================

// GetAllCardsInZone returns all cards in a specific zone
// Used by CDAs like Tarmogoyf (counts card types in graveyards)
func (gc *GameContext) GetAllCardsInZone(ctx context.Context, zone int) []abilities.CardInfo {
	gc.engine.mu.RLock()
	gameState, ok := gc.engine.games[gc.gameID.String()]
	gc.engine.mu.RUnlock()

	if !ok {
		return []abilities.CardInfo{}
	}

	gameState.mu.RLock()
	defer gameState.mu.RUnlock()

	result := []abilities.CardInfo{}

	// Zone constants from abilities package: 0=Library, 1=Hand, 2=Battlefield, 3=Graveyard, 4=Stack, 5=Exile, 6=Command
	switch zone {
	case 0: // ZoneLibrary
		// Get all cards from all players' libraries
		for _, player := range gameState.players {
			for _, card := range player.Library {
				result = append(result, &cardInfoAdapter{card: card})
			}
		}

	case 1: // ZoneHand
		// Get all cards from all players' hands
		for _, player := range gameState.players {
			for _, card := range player.Hand {
				result = append(result, &cardInfoAdapter{card: card})
			}
		}

	case 2: // ZoneBattlefield
		// Get all permanents on battlefield
		for _, card := range gameState.battlefield {
			result = append(result, &cardInfoAdapter{card: card})
		}

	case 3: // ZoneGraveyard
		// Get all cards from all players' graveyards
		for _, player := range gameState.players {
			for _, card := range player.Graveyard {
				result = append(result, &cardInfoAdapter{card: card})
			}
		}

	case 4: // ZoneStack
		// Get all cards on the stack
		// TODO: Implement when stack access method is available
		// gameState.stack is a *rules.StackManager, not a slice

	case 5: // ZoneExile
		// Get all cards in exile
		for _, card := range gameState.exile {
			result = append(result, &cardInfoAdapter{card: card})
		}

	case 6: // ZoneCommand
		// Get all cards in command zone
		for _, card := range gameState.command {
			result = append(result, &cardInfoAdapter{card: card})
		}
	}

	return result
}

// GetCreaturesControlledBy returns all creatures controlled by a player
// Used by CDAs like "power/toughness equal to creatures you control"
func (gc *GameContext) GetCreaturesControlledBy(ctx context.Context, playerID uuid.UUID) []abilities.CardInfo {
	gc.engine.mu.RLock()
	gameState, ok := gc.engine.games[gc.gameID.String()]
	gc.engine.mu.RUnlock()

	if !ok {
		return []abilities.CardInfo{}
	}

	gameState.mu.RLock()
	defer gameState.mu.RUnlock()

	result := []abilities.CardInfo{}
	playerIDStr := playerID.String()

	for _, card := range gameState.battlefield {
		// Check if controlled by player and is a creature
		if card.ControllerID == playerIDStr && isCreature(card) {
			result = append(result, &cardInfoAdapter{card: card})
		}
	}

	return result
}

// GetPlayerHandForCDA returns cards in a player's hand for CDA calculations
// Used by CDAs like Maro (power/toughness equal to hand size)
func (gc *GameContext) GetPlayerHandForCDA(ctx context.Context, playerID uuid.UUID) []abilities.CardInfo {
	gc.engine.mu.RLock()
	gameState, ok := gc.engine.games[gc.gameID.String()]
	gc.engine.mu.RUnlock()

	if !ok {
		return []abilities.CardInfo{}
	}

	gameState.mu.RLock()
	defer gameState.mu.RUnlock()

	player, exists := gameState.players[playerID.String()]
	if !exists {
		return []abilities.CardInfo{}
	}

	result := make([]abilities.CardInfo, 0, len(player.Hand))
	for _, card := range player.Hand {
		result = append(result, &cardInfoAdapter{card: card})
	}

	return result
}

// GetCountersOnPermanent returns the number of a specific counter type on a permanent
// Used by CDAs where power/toughness equals counter count
func (gc *GameContext) GetCountersOnPermanent(ctx context.Context, permanentID uuid.UUID, counterType string) int {
	gc.engine.mu.RLock()
	gameState, ok := gc.engine.games[gc.gameID.String()]
	gc.engine.mu.RUnlock()

	if !ok {
		return 0
	}

	gameState.mu.RLock()
	defer gameState.mu.RUnlock()

	permanentIDStr := permanentID.String()

	// Search battlefield for the permanent
	for _, card := range gameState.battlefield {
		if card.ID == permanentIDStr {
			if card.Counters == nil {
				return 0
			}
			return card.Counters.GetCount(counterType)
		}
	}

	return 0
}

// GetAllCountersOnPermanent returns all counters on a permanent as a map (name -> count)
// Java: permanent.getCounters(game).values()
// Used by effects like Resourceful Defense that need to access all counter types
func (gc *GameContext) GetAllCountersOnPermanent(ctx context.Context, permanentID uuid.UUID) map[string]int {
	gc.engine.mu.RLock()
	gameState, ok := gc.engine.games[gc.gameID.String()]
	gc.engine.mu.RUnlock()

	if !ok {
		return make(map[string]int)
	}

	gameState.mu.RLock()
	defer gameState.mu.RUnlock()

	permanentIDStr := permanentID.String()

	// Search battlefield for the permanent
	for _, card := range gameState.battlefield {
		if card.ID == permanentIDStr {
			if card.Counters == nil {
				return make(map[string]int)
			}
			// Convert counters to simple map
			result := make(map[string]int)
			for name, counter := range card.Counters.GetAll() {
				result[name] = counter.Count
			}
			return result
		}
	}

	return make(map[string]int)
}

// RemoveCountersFromPermanent removes counters from a permanent
// Java: permanent.removeCounters(counterName, amount, source, game)
// Returns error if the permanent doesn't exist
func (gc *GameContext) RemoveCountersFromPermanent(ctx context.Context, permanentID uuid.UUID, counterName string, amount int) error {
	gc.engine.mu.RLock()
	gameState, ok := gc.engine.games[gc.gameID.String()]
	gc.engine.mu.RUnlock()

	if !ok {
		return fmt.Errorf("game %s not found", gc.gameID)
	}

	gameState.mu.Lock()
	defer gameState.mu.Unlock()

	permanentIDStr := permanentID.String()

	// Search battlefield for the permanent
	for _, card := range gameState.battlefield {
		if card.ID == permanentIDStr {
			if card.Counters == nil {
				return nil // No counters to remove
			}
			removed := card.Counters.RemoveCounter(counterName, amount)
			if removed {
				gc.logger.Info("removed counters from permanent",
					zap.String("permanent", card.Name),
					zap.String("counter", counterName),
					zap.Int("amount", amount),
					zap.Int("remaining", card.Counters.GetCount(counterName)))
			}
			return nil
		}
	}

	return fmt.Errorf("permanent %s not found", permanentID)
}

// GetMultiAmountChoice asks the player to distribute amounts among multiple options
// Java: player.getMultiAmountWithIndividualConstraints()
// This is a stub implementation that returns default values for testing/AI
// A real implementation would prompt the player through the UI
func (gc *GameContext) GetMultiAmountChoice(
	ctx context.Context,
	playerID uuid.UUID,
	choices []abilities.MultiAmountChoice,
	totalMin, totalMax int,
	choiceType abilities.MultiAmountType,
) ([]int, error) {
	if len(choices) == 0 {
		return []int{}, nil
	}

	// Simple default strategy: distribute maximum amounts
	// In a real implementation, this would prompt the player
	result := make([]int, len(choices))
	remaining := totalMax

	for i := range choices {
		// Take as much as possible from each choice, up to its max
		take := choices[i].Max
		if take > remaining {
			take = remaining
		}
		if take < choices[i].Min {
			take = choices[i].Min
		}
		result[i] = take
		remaining -= take

		if remaining <= 0 {
			break
		}
	}

	gc.logger.Info("multi-amount choice made (default strategy)",
		zap.String("player", playerID.String()),
		zap.Int("total", totalMax-remaining),
		zap.Int("choiceType", int(choiceType)))

	return result, nil
}

// isCreature checks if a card is a creature
func isCreature(card *internalCard) bool {
	if card == nil {
		return false
	}
	// Simple check - look for "Creature" in Type string
	cardType := card.Type
	if len(cardType) < 8 {
		return false
	}
	for i := 0; i <= len(cardType)-8; i++ {
		if (cardType[i] == 'C' || cardType[i] == 'c') &&
			(cardType[i+1] == 'R' || cardType[i+1] == 'r') &&
			(cardType[i+2] == 'E' || cardType[i+2] == 'e') &&
			(cardType[i+3] == 'A' || cardType[i+3] == 'a') &&
			(cardType[i+4] == 'T' || cardType[i+4] == 't') &&
			(cardType[i+5] == 'U' || cardType[i+5] == 'u') &&
			(cardType[i+6] == 'R' || cardType[i+6] == 'r') &&
			(cardType[i+7] == 'E' || cardType[i+7] == 'e') {
			return true
		}
	}
	return false
}

// cardInfoAdapter adapts internalCard to abilities.CardInfo interface
// This prevents circular dependencies between abilities and game packages
type cardInfoAdapter struct {
	card *internalCard
}

// GetID returns the card's unique identifier
func (c *cardInfoAdapter) GetID() uuid.UUID {
	id, _ := uuid.Parse(c.card.ID)
	return id
}

// GetName returns the card's name
func (c *cardInfoAdapter) GetName() string {
	return c.card.Name
}

// GetTypes returns the card's types (Creature, Artifact, etc.)
func (c *cardInfoAdapter) GetTypes() []string {
	// Parse Type string into slice
	// Format is like "Creature — Human Warrior" or "Instant"
	types := []string{}
	parts := splitOnDash(c.card.Type)
	if len(parts) > 0 {
		typePart := trimSpace(parts[0])
		types = splitOnSpace(typePart)
	}
	return types
}

// GetSubtypes returns the card's subtypes (Human, Warrior, etc.)
func (c *cardInfoAdapter) GetSubtypes() []string {
	// Use the SubTypes field directly - it's already a []string
	return c.card.SubTypes
}

// GetPower returns the card's power (for creatures)
func (c *cardInfoAdapter) GetPower() int {
	power, _ := parseIntOrZero(c.card.Power)
	return power
}

// GetToughness returns the card's toughness (for creatures)
func (c *cardInfoAdapter) GetToughness() int {
	toughness, _ := parseIntOrZero(c.card.Toughness)
	return toughness
}

// Helper functions for string parsing
func splitOnDash(s string) []string {
	result := []string{}
	current := ""
	for _, ch := range s {
		if ch == '—' || ch == '-' {
			if current != "" {
				result = append(result, current)
				current = ""
			}
		} else {
			current += string(ch)
		}
	}
	if current != "" {
		result = append(result, current)
	}
	return result
}

func trimSpace(s string) string {
	start := 0
	end := len(s)

	for start < len(s) && (s[start] == ' ' || s[start] == '\t') {
		start++
	}

	for end > start && (s[end-1] == ' ' || s[end-1] == '\t') {
		end--
	}

	return s[start:end]
}

func splitOnSpace(s string) []string {
	result := []string{}
	current := ""
	for _, ch := range s {
		if ch == ' ' || ch == '\t' {
			if current != "" {
				result = append(result, current)
				current = ""
			}
		} else {
			current += string(ch)
		}
	}
	if current != "" {
		result = append(result, current)
	}
	return result
}

func parseIntOrZero(s string) (int, error) {
	if s == "" || s == "*" || s == "X" {
		return 0, nil
	}

	// Handle cases like "1+*" for Tarmogoyf
	if len(s) > 1 && (s[len(s)-1] == '*' || s[len(s)-1] == 'X') {
		// Just parse the numeric part
		numPart := ""
		for _, ch := range s {
			if ch >= '0' && ch <= '9' {
				numPart += string(ch)
			}
		}
		if numPart != "" {
			result := 0
			for _, ch := range numPart {
				result = result*10 + int(ch-'0')
			}
			return result, nil
		}
		return 0, nil
	}

	// Simple integer parsing
	result := 0
	negative := false
	start := 0

	if len(s) > 0 && s[0] == '-' {
		negative = true
		start = 1
	}

	for i := start; i < len(s); i++ {
		if s[i] < '0' || s[i] > '9' {
			return 0, fmt.Errorf("invalid integer: %s", s)
		}
		result = result*10 + int(s[i]-'0')
	}

	if negative {
		result = -result
	}

	return result, nil
}

// Ensure GameContext implements all required interfaces at compile time
var _ abilities.GameContext = (*GameContext)(nil)
var _ abilities.CounterGameContext = (*GameContext)(nil)
var _ abilities.TokenGameContext = (*GameContext)(nil)
