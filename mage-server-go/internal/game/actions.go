package game

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

// actions.go implements all operations from playtest-game.ts lines 492-1151
// Each function is a direct translation of the TypeScript playtest operations

// DrawCards draws N cards from a player's library to their hand
// From playtest-game.ts lines 492-527
func (e *GameEngine) DrawCards(gameID, playerID string, count int) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	state := e.games[gameID]
	if state == nil {
		return fmt.Errorf("game not found: %s", gameID)
	}

	player, ok := state.Players[playerID]
	if !ok {
		return fmt.Errorf("player not found: %s", playerID)
	}

	before := len(player.Library)
	actualCount := min(count, len(player.Library))

	// Draw cards from top of library
	drawn := make([]*Card, actualCount)
	copy(drawn, player.Library[:actualCount])
	player.Library = player.Library[actualCount:]

	// Update zone and make cards visible
	for _, card := range drawn {
		card.Zone = ZoneHandStr
		card.FaceDown = false
	}

	// Add to hand
	player.Hand = append(player.Hand, drawn...)
	player.HandCount = len(player.Hand)
	player.LibraryCount = len(player.Library)

	msg := fmt.Sprintf("%s draws %d", player.Name, len(drawn))
	if count != len(drawn) {
		msg += fmt.Sprintf(" (requested %d)", count)
	}
	msg += fmt.Sprintf(". Library: %d → %d.", before, len(player.Library))

	e.appendLog(state, "draw", msg)
	e.broadcast(gameID)
	return nil
}

// PlayCard moves a card from hand to battlefield
// From playtest-game.ts lines 532-568
func (e *GameEngine) PlayCard(gameID, playerID, cardID string, tapped bool) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	state := e.games[gameID]
	if state == nil {
		return fmt.Errorf("game not found: %s", gameID)
	}

	player, ok := state.Players[playerID]
	if !ok {
		return fmt.Errorf("player not found: %s", playerID)
	}

	// Find card in hand
	cardIndex := -1
	var card *Card
	for i, c := range player.Hand {
		if c.ID == cardID {
			cardIndex = i
			card = c
			break
		}
	}

	if cardIndex == -1 {
		return fmt.Errorf("card not found in hand: %s", cardID)
	}

	// Remove from hand
	player.Hand = append(player.Hand[:cardIndex], player.Hand[cardIndex+1:]...)
	player.HandCount = len(player.Hand)

	// Update card properties
	card.Zone = ZoneBattlefieldStr
	card.ControllerID = playerID
	card.Tapped = tapped
	card.FaceDown = false

	// Add to battlefield
	state.Battlefield = append(state.Battlefield, card)

	msg := fmt.Sprintf("%s plays %s", player.Name, card.Name)
	if tapped {
		msg += " (tapped)"
	}
	msg += fmt.Sprintf(". Hand: %d → %d.", len(player.Hand)+1, len(player.Hand))

	e.appendLog(state, "play", msg)
	e.broadcast(gameID)
	return nil
}

// MoveCard moves a card to a different zone
// From playtest-game.ts lines 574-607
func (e *GameEngine) MoveCard(gameID, playerID, cardID, targetZone string) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	state := e.games[gameID]
	if state == nil {
		return fmt.Errorf("game not found: %s", gameID)
	}

	// Find the card in any zone
	card, sourceZone := e.findCardInState(state, cardID)
	if card == nil {
		return fmt.Errorf("card not found: %s", cardID)
	}

	// MTG Rule: Tokens cease to exist when they leave the battlefield
	// From playtest-game.ts lines 590-596
	if strings.HasPrefix(cardID, "token-") && sourceZone == ZoneBattlefieldStr {
		e.removeCardFromZone(state, cardID, sourceZone)
		msg := fmt.Sprintf("%s moves %s to %s. Token ceases to exist.",
			e.getPlayerName(state, playerID), card.Name, targetZone)
		e.appendLog(state, "move", msg)
		e.broadcast(gameID)
		return nil
	}

	// Remove from source zone
	e.removeCardFromZone(state, cardID, sourceZone)

	// Add to target zone
	e.addCardToZone(state, card, targetZone, playerID)

	msg := fmt.Sprintf("%s moves %s from %s to %s.",
		e.getPlayerName(state, playerID), card.Name, sourceZone, targetZone)
	e.appendLog(state, "move", msg)
	e.broadcast(gameID)
	return nil
}

// TapCard taps or untaps a card
// From playtest-game.ts lines 612-630
func (e *GameEngine) TapCard(gameID, playerID, cardID string, tapped bool) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	state := e.games[gameID]
	if state == nil {
		return fmt.Errorf("game not found: %s", gameID)
	}

	// Find card on battlefield
	var card *Card
	for _, c := range state.Battlefield {
		if c.ID == cardID {
			card = c
			break
		}
	}

	if card == nil {
		return fmt.Errorf("card not found on battlefield: %s", cardID)
	}

	card.Tapped = tapped

	action := "taps"
	if !tapped {
		action = "untaps"
	}
	msg := fmt.Sprintf("%s %s %s.", e.getPlayerName(state, playerID), action, card.Name)
	e.appendLog(state, "tap", msg)
	e.broadcast(gameID)
	return nil
}

// UntapAll untaps all permanents controlled by a player
// From playtest-game.ts lines 635-646
func (e *GameEngine) UntapAll(gameID, playerID string) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	state := e.games[gameID]
	if state == nil {
		return fmt.Errorf("game not found: %s", gameID)
	}

	for _, card := range state.Battlefield {
		if card.ControllerID == playerID {
			card.Tapped = false
		}
	}

	msg := fmt.Sprintf("%s untaps all permanents.", e.getPlayerName(state, playerID))
	e.appendLog(state, "untap", msg)
	e.broadcast(gameID)
	return nil
}

// FlipCard flips a card face up/down
// From playtest-game.ts lines 651-663
func (e *GameEngine) FlipCard(gameID, playerID, cardID string, faceDown bool) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	state := e.games[gameID]
	if state == nil {
		return fmt.Errorf("game not found: %s", gameID)
	}

	var card *Card
	for _, c := range state.Battlefield {
		if c.ID == cardID {
			card = c
			break
		}
	}

	if card == nil {
		return fmt.Errorf("card not found: %s", cardID)
	}

	card.FaceDown = faceDown

	direction := "up"
	if faceDown {
		direction = "down"
	}
	msg := fmt.Sprintf("%s flips %s face %s.", e.getPlayerName(state, playerID), card.Name, direction)
	e.appendLog(state, "flip", msg)
	e.broadcast(gameID)
	return nil
}

// ModifyLife modifies a player's life total
// From playtest-game.ts lines 668-679
func (e *GameEngine) ModifyLife(gameID, playerID string, delta int) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	state := e.games[gameID]
	if state == nil {
		return fmt.Errorf("game not found: %s", gameID)
	}

	player, ok := state.Players[playerID]
	if !ok {
		return fmt.Errorf("player not found: %s", playerID)
	}

	player.Life = max(0, player.Life+delta)

	msg := fmt.Sprintf("%s modifies life by %d.", player.Name, delta)
	e.appendLog(state, "life", msg)
	e.broadcast(gameID)
	return nil
}

// SetPlayerCounter sets a player counter (poison, energy, etc.)
// From playtest-game.ts lines 684-700
func (e *GameEngine) SetPlayerCounter(gameID, playerID, counterType string, value int) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	state := e.games[gameID]
	if state == nil {
		return fmt.Errorf("game not found: %s", gameID)
	}

	player, ok := state.Players[playerID]
	if !ok {
		return fmt.Errorf("player not found: %s", playerID)
	}

	switch counterType {
	case "poison":
		player.Poison = max(0, value)
	case "energy":
		player.Energy = max(0, value)
	default:
		return fmt.Errorf("unknown counter type: %s", counterType)
	}

	msg := fmt.Sprintf("%s sets %s to %d.", player.Name, counterType, value)
	e.appendLog(state, "counter", msg)
	e.broadcast(gameID)
	return nil
}

// ShuffleLibrary shuffles a player's library using Fisher-Yates algorithm
// From playtest-game.ts lines 705-717
func (e *GameEngine) ShuffleLibrary(gameID, playerID string) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	state := e.games[gameID]
	if state == nil {
		return fmt.Errorf("game not found: %s", gameID)
	}

	player, ok := state.Players[playerID]
	if !ok {
		return fmt.Errorf("player not found: %s", playerID)
	}

	// Fisher-Yates shuffle
	library := player.Library
	for i := len(library) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		library[i], library[j] = library[j], library[i]
	}

	player.RevealedTopCard = false // Clear revealed top when shuffling

	msg := fmt.Sprintf("%s shuffles their library.", player.Name)
	e.appendLog(state, "shuffle", msg)
	e.broadcast(gameID)
	return nil
}

// CreateToken creates a token on the battlefield
// From playtest-game.ts lines 759-805
func (e *GameEngine) CreateToken(gameID, playerID, name, types, power, toughness, color string) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	state := e.games[gameID]
	if state == nil {
		return fmt.Errorf("game not found: %s", gameID)
	}

	// Generate unique token ID
	tokenID := fmt.Sprintf("token-%d-%d", time.Now().UnixNano(), rand.Intn(1000000))

	token := &Card{
		ID:                tokenID,
		Name:              name,
		DisplayName:       name,
		SubTypes:          "",
		SuperTypes:        "",
		Color:             color,
		Type:              types,
		Power:             power,
		Toughness:         toughness,
		Loyalty:           "",
		ManaCost:          "",
		CardNumber:        0,
		ExpansionSet:      "",
		Rarity:            "",
		RulesText:         "",
		Zone:              ZoneBattlefieldStr,
		OwnerID:           playerID,
		ControllerID:      playerID,
		Tapped:            false,
		Flipped:           false,
		Transformed:       false,
		FaceDown:          false,
		Counters:          make([]Counter, 0),
		AttachedTo:        make([]string, 0),
		SummoningSickness: true,
	}

	state.Battlefield = append(state.Battlefield, token)

	msg := fmt.Sprintf("%s creates %s token.", e.getPlayerName(state, playerID), token.Name)
	e.appendLog(state, "create", msg)
	e.broadcast(gameID)
	return nil
}

// AddCounter adds counters to a card
// From playtest-game.ts lines 810-840
func (e *GameEngine) AddCounter(gameID, playerID, cardID, counterName string, amount int) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	state := e.games[gameID]
	if state == nil {
		return fmt.Errorf("game not found: %s", gameID)
	}

	card, sourceZone := e.findCardInState(state, cardID)
	if card == nil {
		return fmt.Errorf("card not found: %s", cardID)
	}

	// Find existing counter or create new one
	found := false
	for i := range card.Counters {
		if card.Counters[i].Name == counterName {
			card.Counters[i].Count += amount
			found = true
			break
		}
	}

	if !found {
		card.Counters = append(card.Counters, Counter{
			Name:  counterName,
			Count: amount,
		})
	}

	msg := fmt.Sprintf("Added %d %s counter(s) to %s (zone: %s).", amount, counterName, card.Name, sourceZone)
	e.appendLog(state, "counter", msg)
	e.broadcast(gameID)
	return nil
}

// RemoveCounter removes counters from a card
// From playtest-game.ts lines 845-882
func (e *GameEngine) RemoveCounter(gameID, playerID, cardID, counterName string, amount int) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	state := e.games[gameID]
	if state == nil {
		return fmt.Errorf("game not found: %s", gameID)
	}

	card, _ := e.findCardInState(state, cardID)
	if card == nil {
		return fmt.Errorf("card not found: %s", cardID)
	}

	// Find and update counter
	for i := range card.Counters {
		if card.Counters[i].Name == counterName {
			newCount := max(0, card.Counters[i].Count-amount)
			if newCount == 0 {
				// Remove counter if count reaches 0
				card.Counters = append(card.Counters[:i], card.Counters[i+1:]...)
			} else {
				card.Counters[i].Count = newCount
			}
			break
		}
	}

	msg := fmt.Sprintf("Removed %d %s counter(s) from %s.", amount, counterName, card.Name)
	e.appendLog(state, "counter", msg)
	e.broadcast(gameID)
	return nil
}

// SetCounter sets a counter to a specific value
// From playtest-game.ts lines 887-918
func (e *GameEngine) SetCounter(gameID, playerID, cardID, counterName string, amount int) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	state := e.games[gameID]
	if state == nil {
		return fmt.Errorf("game not found: %s", gameID)
	}

	card, _ := e.findCardInState(state, cardID)
	if card == nil {
		return fmt.Errorf("card not found: %s", cardID)
	}

	if amount <= 0 {
		// Remove counter if setting to 0 or less
		for i := range card.Counters {
			if card.Counters[i].Name == counterName {
				card.Counters = append(card.Counters[:i], card.Counters[i+1:]...)
				break
			}
		}
	} else {
		// Set or update counter
		found := false
		for i := range card.Counters {
			if card.Counters[i].Name == counterName {
				card.Counters[i].Count = amount
				found = true
				break
			}
		}
		if !found {
			card.Counters = append(card.Counters, Counter{
				Name:  counterName,
				Count: amount,
			})
		}
	}

	msg := fmt.Sprintf("Set %s counters on %s to %d.", counterName, card.Name, amount)
	e.appendLog(state, "counter", msg)
	e.broadcast(gameID)
	return nil
}

// MillCards moves top N cards from library to graveyard
// From playtest-game.ts lines 923-957
func (e *GameEngine) MillCards(gameID, playerID string, count int) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	state := e.games[gameID]
	if state == nil {
		return fmt.Errorf("game not found: %s", gameID)
	}

	player, ok := state.Players[playerID]
	if !ok {
		return fmt.Errorf("player not found: %s", playerID)
	}

	before := len(player.Library)
	actualCount := min(count, len(player.Library))

	// Mill cards from top of library
	milled := make([]*Card, actualCount)
	copy(milled, player.Library[:actualCount])
	player.Library = player.Library[actualCount:]

	// Update zone for milled cards
	for _, card := range milled {
		card.Zone = ZoneGraveyardStr
		card.FaceDown = false
	}

	// Add to graveyard
	player.Graveyard = append(player.Graveyard, milled...)
	player.LibraryCount = len(player.Library)

	msg := fmt.Sprintf("%s mills %d card(s). Library: %d → %d.",
		player.Name, actualCount, before, len(player.Library))
	e.appendLog(state, "mill", msg)
	e.broadcast(gameID)
	return nil
}

// ScryCards implements scry by reordering library cards
// From playtest-game.ts lines 1016-1053
func (e *GameEngine) ScryCards(gameID, playerID string, scryCount int, keepOnTop, putToBottom []string) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	state := e.games[gameID]
	if state == nil {
		return fmt.Errorf("game not found: %s", gameID)
	}

	player, ok := state.Players[playerID]
	if !ok {
		return fmt.Errorf("player not found: %s", playerID)
	}

	// Build sets of card IDs being scried
	scryCardIDs := make(map[string]bool)
	for _, id := range keepOnTop {
		scryCardIDs[id] = true
	}
	for _, id := range putToBottom {
		scryCardIDs[id] = true
	}

	// Separate scried cards from remaining library
	keepCards := make([]*Card, 0)
	bottomCards := make([]*Card, 0)
	remaining := make([]*Card, 0)

	for _, card := range player.Library {
		if scryCardIDs[card.ID] {
			// Determine if this card goes to top or bottom
			isKept := false
			for _, id := range keepOnTop {
				if id == card.ID {
					keepCards = append(keepCards, card)
					isKept = true
					break
				}
			}
			if !isKept {
				for _, id := range putToBottom {
					if id == card.ID {
						bottomCards = append(bottomCards, card)
						break
					}
				}
			}
		} else {
			remaining = append(remaining, card)
		}
	}

	// Rebuild library: keep on top, remaining cards, put to bottom
	player.Library = make([]*Card, 0, len(keepCards)+len(remaining)+len(bottomCards))
	player.Library = append(player.Library, keepCards...)
	player.Library = append(player.Library, remaining...)
	player.Library = append(player.Library, bottomCards...)

	msg := fmt.Sprintf("%s completes scry %d (%d kept on top, %d to bottom).",
		player.Name, scryCount, len(keepCards), len(bottomCards))
	e.appendLog(state, "scry", msg)
	e.broadcast(gameID)
	return nil
}

// SetRevealedTop sets whether the top card of library is revealed
// From playtest-game.ts lines 1058-1071
func (e *GameEngine) SetRevealedTop(gameID, playerID string, revealed bool) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	state := e.games[gameID]
	if state == nil {
		return fmt.Errorf("game not found: %s", gameID)
	}

	player, ok := state.Players[playerID]
	if !ok {
		return fmt.Errorf("player not found: %s", playerID)
	}

	player.RevealedTopCard = revealed

	action := "reveals"
	if !revealed {
		action = "hides"
	}
	msg := fmt.Sprintf("%s %s the top card of their library.", player.Name, action)
	e.appendLog(state, "reveal", msg)
	e.broadcast(gameID)
	return nil
}

// NextTurn advances to the next player's turn
// From playtest-game.ts lines 1076-1090
func (e *GameEngine) NextTurn(gameID, playerID string) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	state := e.games[gameID]
	if state == nil {
		return fmt.Errorf("game not found: %s", gameID)
	}

	// Get player order.
	//
	// This MUST come from state.TurnOrder, not from ranging over state.Players:
	// Players is a map and Go randomizes map iteration, so the previous version
	// of this loop produced a different turn order on every single call.
	playerIDs := state.turnOrder()
	if len(playerIDs) == 0 {
		return fmt.Errorf("game has no players: %s", gameID)
	}

	// Find next player
	currentIdx := -1
	for i, pid := range playerIDs {
		if pid == state.ActivePlayerID {
			currentIdx = i
			break
		}
	}

	nextIdx := (currentIdx + 1) % len(playerIDs)
	nextPlayerID := playerIDs[nextIdx]

	state.Turn++
	state.ActivePlayerID = nextPlayerID
	state.ActiveControlSeat = nextPlayerID

	msg := fmt.Sprintf("%s ends their turn.", e.getPlayerName(state, playerID))
	e.appendLog(state, "endTurn", msg)
	e.broadcast(gameID)
	return nil
}

// Mulligan performs a mulligan for a player
// From playtest-game.ts lines 1095-1146
func (e *GameEngine) Mulligan(gameID, playerID string) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	state := e.games[gameID]
	if state == nil {
		return fmt.Errorf("game not found: %s", gameID)
	}

	player, ok := state.Players[playerID]
	if !ok {
		return fmt.Errorf("player not found: %s", playerID)
	}

	newMulliganCount := player.MulliganCount + 1

	// Return hand to library
	for _, card := range player.Hand {
		card.Zone = ZoneLibraryStr
		card.FaceDown = true
		player.Library = append(player.Library, card)
	}

	// Shuffle library
	library := player.Library
	for i := len(library) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		library[i], library[j] = library[j], library[i]
	}

	// Calculate new hand size based on free mulligans
	var newHandSize int
	if newMulliganCount <= state.FreeMulligans {
		newHandSize = 7 // Free mulligan - draw full 7
	} else {
		penaltyMulligans := newMulliganCount - state.FreeMulligans
		newHandSize = max(0, 7-penaltyMulligans)
	}

	// Draw new hand
	actualHandSize := min(newHandSize, len(player.Library))
	newHand := make([]*Card, actualHandSize)
	copy(newHand, player.Library[:actualHandSize])
	player.Library = player.Library[actualHandSize:]

	for _, card := range newHand {
		card.Zone = ZoneHandStr
		card.FaceDown = false
	}

	player.Hand = newHand
	player.HandCount = len(newHand)
	player.LibraryCount = len(player.Library)
	player.KeptHand = false
	player.MulliganCount = newMulliganCount

	msg := fmt.Sprintf("%s mulligans their hand.", player.Name)
	e.appendLog(state, "mulligan", msg)
	e.broadcast(gameID)
	return nil
}

// KeepHand marks that a player is keeping their hand
// From playtest-game.ts lines 1151-1163
func (e *GameEngine) KeepHand(gameID, playerID string) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	state := e.games[gameID]
	if state == nil {
		return fmt.Errorf("game not found: %s", gameID)
	}

	player, ok := state.Players[playerID]
	if !ok {
		return fmt.Errorf("player not found: %s", playerID)
	}

	player.KeptHand = true

	msg := fmt.Sprintf("%s keeps their hand.", player.Name)
	e.appendLog(state, "keep", msg)
	e.broadcast(gameID)
	return nil
}

// Helper functions

func (e *GameEngine) appendLog(state *GameState, kind, message string) {
	entry := LogEntry{
		Kind:      kind,
		Message:   message,
		Timestamp: time.Now(),
	}
	state.Log = append(state.Log, entry)
}

func (e *GameEngine) getPlayerName(state *GameState, playerID string) string {
	if player, ok := state.Players[playerID]; ok {
		return player.Name
	}
	return playerID
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
