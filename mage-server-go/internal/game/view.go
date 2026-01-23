package game

import (
	"strings"
)

// view.go implements game view with hidden information filtering
// Opponent hands show counts only, not cards
// Opponent libraries show counts only, not cards
// Public zones (battlefield, graveyard, exile) visible to all

// PlaytestGameView represents a player-specific view of the game
type PlaytestGameView struct {
	GameID            string                  `json:"gameId"`
	ViewerID          string                  `json:"viewerId"`
	ActiveControlSeat string                  `json:"activeControlSeat"`
	Me                *PlaytestPlayerView     `json:"me"`
	Opponents         []*PlaytestOpponentView `json:"opponents"`
	Battlefield       []*Card                 `json:"battlefield"`
	Exile             []*Card                 `json:"exile"`
	Stack             []*Card                 `json:"stack"`
	Command           []*Card                 `json:"command"`
	Turn              int                     `json:"turn"`
	ActivePlayerID    string                  `json:"activePlayerId"`
	IsInitialized     bool                    `json:"isInitialized"`
	Log               []LogEntry              `json:"log"`
	MulliganType      string                  `json:"mulliganType"`
	FreeMulligans     int                     `json:"freeMulligans"`
}

// PlaytestPlayerView represents the viewing player's full state
type PlaytestPlayerView struct {
	PlayerID        string    `json:"playerId"`
	Name            string    `json:"name"`
	Life            int       `json:"life"`
	Poison          int       `json:"poison"`
	Energy          int       `json:"energy"`
	LibraryCount    int       `json:"libraryCount"`
	HandCount       int       `json:"handCount"`
	Hand            []*Card   `json:"hand"`    // Full visibility of own hand
	Library         []*Card   `json:"library"` // Full visibility of own library
	Graveyard       []*Card   `json:"graveyard"`
	ManaPool        *ManaPool `json:"manaPool"`
	KeptHand        bool      `json:"keptHand"`
	MulliganCount   int       `json:"mulliganCount"`
	RevealedTopCard bool      `json:"revealedTopCard"`
}

// PlaytestOpponentView represents an opponent with hidden information
type PlaytestOpponentView struct {
	PlayerID        string    `json:"playerId"`
	Name            string    `json:"name"`
	Life            int       `json:"life"`
	Poison          int       `json:"poison"`
	Energy          int       `json:"energy"`
	LibraryCount    int       `json:"libraryCount"` // Count only
	HandCount       int       `json:"handCount"`    // Count only
	Hand            []*Card   `json:"hand"`         // Empty - hidden
	Library         []*Card   `json:"library"`      // Empty - hidden
	TopCard         *Card     `json:"topCard"`      // Only if revealed
	Graveyard       []*Card   `json:"graveyard"`    // Public
	ManaPool        *ManaPool `json:"manaPool"`
	KeptHand        bool      `json:"keptHand"`
	MulliganCount   int       `json:"mulliganCount"`
	RevealedTopCard bool      `json:"revealedTopCard"`
}

// buildGameView creates a player-specific game view with hidden information filtering
func (e *GameEngine) buildGameView(state *GameState, viewerID string) *PlaytestGameView {
	view := &PlaytestGameView{
		GameID:            state.GameID,
		ViewerID:          viewerID,
		ActiveControlSeat: state.ActiveControlSeat,
		Battlefield:       state.Battlefield, // Public zone
		Exile:             state.Exile,       // Public zone
		Stack:             state.Stack,       // Public zone
		Command:           state.Command,     // Public zone
		Turn:              state.Turn,
		ActivePlayerID:    state.ActivePlayerID,
		IsInitialized:     state.IsInitialized,
		Log:               state.Log,
		MulliganType:      state.MulliganType,
		FreeMulligans:     state.FreeMulligans,
		Opponents:         make([]*PlaytestOpponentView, 0),
	}

	// Build player views
	for playerID, player := range state.Players {
		if playerID == viewerID {
			// Full visibility for the viewing player
			view.Me = &PlaytestPlayerView{
				PlayerID:        player.PlayerID,
				Name:            player.Name,
				Life:            player.Life,
				Poison:          player.Poison,
				Energy:          player.Energy,
				LibraryCount:    player.LibraryCount,
				HandCount:       player.HandCount,
				Hand:            player.Hand,    // Full visibility
				Library:         player.Library, // Full visibility
				Graveyard:       player.Graveyard,
				ManaPool:        player.ManaPool,
				KeptHand:        player.KeptHand,
				MulliganCount:   player.MulliganCount,
				RevealedTopCard: player.RevealedTopCard,
			}
		} else {
			// Hidden information for opponents
			opponentView := &PlaytestOpponentView{
				PlayerID:        player.PlayerID,
				Name:            player.Name,
				Life:            player.Life,
				Poison:          player.Poison,
				Energy:          player.Energy,
				LibraryCount:    player.LibraryCount, // Count only
				HandCount:       player.HandCount,    // Count only
				Hand:            make([]*Card, 0),    // Hidden
				Library:         make([]*Card, 0),    // Hidden
				Graveyard:       player.Graveyard,    // Public
				ManaPool:        player.ManaPool,
				KeptHand:        player.KeptHand,
				MulliganCount:   player.MulliganCount,
				RevealedTopCard: player.RevealedTopCard,
			}

			// If top card is revealed, include it
			if player.RevealedTopCard && len(player.Library) > 0 {
				opponentView.TopCard = player.Library[0]
			}

			view.Opponents = append(view.Opponents, opponentView)
		}
	}

	return view
}

// findCardInState finds a card in any zone of the game state
// From playtest-helpers.ts lines 30-63
func (e *GameEngine) findCardInState(state *GameState, cardID string) (*Card, string) {
	// Check battlefield
	for _, card := range state.Battlefield {
		if card.ID == cardID {
			return card, ZoneBattlefieldStr
		}
	}

	// Check exile
	for _, card := range state.Exile {
		if card.ID == cardID {
			return card, ZoneExileStr
		}
	}

	// Check command
	for _, card := range state.Command {
		if card.ID == cardID {
			return card, ZoneCommandStr
		}
	}

	// Check stack
	for _, card := range state.Stack {
		if card.ID == cardID {
			return card, ZoneStackStr
		}
	}

	// Check player zones
	for playerID, player := range state.Players {
		for _, card := range player.Hand {
			if card.ID == cardID {
				return card, "hand:" + playerID
			}
		}

		for _, card := range player.Library {
			if card.ID == cardID {
				return card, "library:" + playerID
			}
		}

		for _, card := range player.Graveyard {
			if card.ID == cardID {
				return card, "graveyard:" + playerID
			}
		}
	}

	return nil, ""
}

// removeCardFromZone removes a card from its source zone
// From playtest-helpers.ts lines 120-170
func (e *GameEngine) removeCardFromZone(state *GameState, cardID, sourceZone string) {
	if sourceZone == ZoneBattlefieldStr {
		newBattlefield := make([]*Card, 0, len(state.Battlefield))
		for _, card := range state.Battlefield {
			if card.ID != cardID {
				newBattlefield = append(newBattlefield, card)
			}
		}
		state.Battlefield = newBattlefield
		return
	}

	if sourceZone == ZoneExileStr {
		newExile := make([]*Card, 0, len(state.Exile))
		for _, card := range state.Exile {
			if card.ID != cardID {
				newExile = append(newExile, card)
			}
		}
		state.Exile = newExile
		return
	}

	if sourceZone == ZoneCommandStr {
		newCommand := make([]*Card, 0, len(state.Command))
		for _, card := range state.Command {
			if card.ID != cardID {
				newCommand = append(newCommand, card)
			}
		}
		state.Command = newCommand
		return
	}

	if sourceZone == ZoneStackStr {
		newStack := make([]*Card, 0, len(state.Stack))
		for _, card := range state.Stack {
			if card.ID != cardID {
				newStack = append(newStack, card)
			}
		}
		state.Stack = newStack
		return
	}

	// Player zones
	if len(sourceZone) > 5 && sourceZone[:5] == "hand:" {
		playerID := sourceZone[5:]
		if player, ok := state.Players[playerID]; ok {
			newHand := make([]*Card, 0, len(player.Hand))
			for _, card := range player.Hand {
				if card.ID != cardID {
					newHand = append(newHand, card)
				}
			}
			player.Hand = newHand
			player.HandCount = len(newHand)
		}
		return
	}

	if len(sourceZone) > 8 && sourceZone[:8] == "library:" {
		playerID := sourceZone[8:]
		if player, ok := state.Players[playerID]; ok {
			newLibrary := make([]*Card, 0, len(player.Library))
			for _, card := range player.Library {
				if card.ID != cardID {
					newLibrary = append(newLibrary, card)
				}
			}
			player.Library = newLibrary
			player.LibraryCount = len(newLibrary)
		}
		return
	}

	if len(sourceZone) > 10 && sourceZone[:10] == "graveyard:" {
		playerID := sourceZone[10:]
		if player, ok := state.Players[playerID]; ok {
			newGraveyard := make([]*Card, 0, len(player.Graveyard))
			for _, card := range player.Graveyard {
				if card.ID != cardID {
					newGraveyard = append(newGraveyard, card)
				}
			}
			player.Graveyard = newGraveyard
		}
		return
	}
}

// addCardToZone adds a card to a target zone
// From playtest-helpers.ts lines 175-243
func (e *GameEngine) addCardToZone(state *GameState, card *Card, targetZone, controllerID string) {
	if card == nil {
		return
	}

	// Prepare card for new zone
	card.FaceDown = false

	normalizedZone := normalizeZoneName(targetZone)

	switch normalizedZone {
	case ZoneBattlefieldStr:
		card.Zone = ZoneBattlefieldStr
		card.ControllerID = controllerID
		state.Battlefield = append(state.Battlefield, card)

	case ZoneGraveyardStr:
		card.Zone = ZoneGraveyardStr
		ownerID := card.OwnerID
		if ownerID == "" {
			ownerID = controllerID
		}
		if player, ok := state.Players[ownerID]; ok {
			player.Graveyard = append(player.Graveyard, card)
		}

	case ZoneExileStr:
		card.Zone = ZoneExileStr
		state.Exile = append(state.Exile, card)

	case ZoneHandStr:
		card.Zone = ZoneHandStr
		if player, ok := state.Players[controllerID]; ok {
			player.Hand = append(player.Hand, card)
			player.HandCount = len(player.Hand)
		}

	case ZoneCommandStr:
		card.Zone = ZoneCommandStr
		card.FaceDown = false
		state.Command = append(state.Command, card)

	case ZoneLibraryStr:
		card.Zone = ZoneLibraryStr
		card.FaceDown = true
		if player, ok := state.Players[controllerID]; ok {
			// Check for position specifiers (TOP, BOTTOM)
			if targetZone == "LIBRARY_TOP" || targetZone == "LIBRARY" {
				// Add to top of library
				player.Library = append([]*Card{card}, player.Library...)
			} else if targetZone == "LIBRARY_BOTTOM" {
				// Add to bottom of library
				player.Library = append(player.Library, card)
			} else {
				// Default to top
				player.Library = append([]*Card{card}, player.Library...)
			}
			player.LibraryCount = len(player.Library)
		}

	case ZoneStackStr:
		card.Zone = ZoneStackStr
		state.Stack = append(state.Stack, card)
	}
}

// normalizeZoneName converts zone name to standardized format
func normalizeZoneName(zone string) string {
	// Convert to uppercase and strip position specifiers
	upper := strings.ToUpper(zone)

	if strings.HasPrefix(upper, "LIBRARY") {
		return ZoneLibraryStr
	}
	if strings.HasPrefix(upper, "HAND") {
		return ZoneHandStr
	}
	if strings.HasPrefix(upper, "BATTLEFIELD") {
		return ZoneBattlefieldStr
	}
	if strings.HasPrefix(upper, "GRAVEYARD") {
		return ZoneGraveyardStr
	}
	if strings.HasPrefix(upper, "EXILE") {
		return ZoneExileStr
	}
	if strings.HasPrefix(upper, "COMMAND") {
		return ZoneCommandStr
	}
	if strings.HasPrefix(upper, "STACK") {
		return ZoneStackStr
	}

	// Abbreviations
	switch upper {
	case "H":
		return ZoneHandStr
	case "L":
		return ZoneLibraryStr
	case "B":
		return ZoneBattlefieldStr
	case "G":
		return ZoneGraveyardStr
	case "E":
		return ZoneExileStr
	case "C":
		return ZoneCommandStr
	case "S":
		return ZoneStackStr
	}

	return upper
}
