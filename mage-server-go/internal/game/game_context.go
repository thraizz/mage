package game

import (
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

// Ensure GameContext implements all required interfaces at compile time
var _ abilities.GameContext = (*GameContext)(nil)
var _ abilities.CounterGameContext = (*GameContext)(nil)
var _ abilities.TokenGameContext = (*GameContext)(nil)
