package game

import (
	"sort"
	"time"
)

// game_state.go defines the core game state structures for the playtest engine.
// Based on playtest-game.ts lines 25-59: State structure

// GameState represents the complete game state for the rules-light game engine.
// All zones, players, and cards in one synchronized structure.
type GameState struct {
	GameID            string             `json:"gameId"`
	ActiveControlSeat string             `json:"activeControlSeat"` // Which player perspective is controlling
	Players           map[string]*Player `json:"players"`
	// TurnOrder is the seating order of Players, captured once at game
	// creation from the players slice handed to StartGame/StartGameWithDecks.
	//
	// It exists because Players is a map and Go randomizes map iteration:
	// deriving turn order by ranging over Players produced a different order on
	// every call, so NextTurn was non-deterministic and no test that depends on
	// sequencing could be reproducible. Sorting the IDs would also be
	// deterministic, but real Magic turn order is seating order, not
	// alphabetical, and the seating order is already available here.
	TurnOrder      []string   `json:"turnOrder"`
	Battlefield    []*Card    `json:"battlefield"`
	Exile          []*Card    `json:"exile"`
	Stack          []*Card    `json:"stack"`
	Command        []*Card    `json:"command"`
	Turn           int        `json:"turn"`
	ActivePlayerID string     `json:"activePlayerId"`
	IsInitialized  bool       `json:"isInitialized"`
	Log            []LogEntry `json:"log"`
	MulliganType   string     `json:"mulliganType"` // "london"
	FreeMulligans  int        `json:"freeMulligans"`
	StartedAt      time.Time  `json:"startedAt"`
}

// Player represents a player in the game engine.
// Based on playtest-game.ts lines 25-40
type Player struct {
	PlayerID        string    `json:"playerId"`
	Name            string    `json:"name"`
	Life            int       `json:"life"`
	Poison          int       `json:"poison"`
	Energy          int       `json:"energy"`
	LibraryCount    int       `json:"libraryCount"`
	HandCount       int       `json:"handCount"`
	Hand            []*Card   `json:"hand"`
	Library         []*Card   `json:"library"`
	Graveyard       []*Card   `json:"graveyard"`
	ManaPool        *ManaPool `json:"manaPool"`
	KeptHand        bool      `json:"keptHand"`
	MulliganCount   int       `json:"mulliganCount"`
	RevealedTopCard bool      `json:"revealedTopCard"` // When true, top card is visible
}

// Card represents a card in the game engine.
// Simplified from LegacyCard type - no rules enforcement, direct player control.
type Card struct {
	// Identity
	ID          string `json:"id"`
	Name        string `json:"name"`
	DisplayName string `json:"displayName"`
	OwnerID     string `json:"ownerId"`

	// Card characteristics
	ManaCost     string `json:"manaCost"`
	Type         string `json:"type"` // Full type line
	SubTypes     string `json:"subTypes"`
	SuperTypes   string `json:"superTypes"`
	Color        string `json:"color"`
	Power        string `json:"power"`
	Toughness    string `json:"toughness"`
	Loyalty      string `json:"loyalty"`
	CardNumber   int    `json:"cardNumber"`
	ExpansionSet string `json:"expansionSetCode"`
	Rarity       string `json:"rarity"`
	RulesText    string `json:"rulesText"`

	// Game state
	Zone         string    `json:"zone"`
	ControllerID string    `json:"controllerId"`
	Tapped       bool      `json:"tapped"`
	Flipped      bool      `json:"flipped"`
	Transformed  bool      `json:"transformed"`
	FaceDown     bool      `json:"faceDown"`
	Counters     []Counter `json:"counters"`

	// Combat state (manual tracking, no enforcement)
	Attacking         bool `json:"attacking"`
	Blocking          bool `json:"blocking"`
	SummoningSickness bool `json:"summoningSickness"`

	// Attachments
	AttachedTo []string `json:"attachedTo"`
}

// Counter represents a counter on a card
type Counter struct {
	Name  string `json:"name"`
	Count int    `json:"count"`
}

// ManaPool represents a player's mana pool (client-side tracking only)
type ManaPool struct {
	White     int `json:"white"`
	Blue      int `json:"blue"`
	Black     int `json:"black"`
	Red       int `json:"red"`
	Green     int `json:"green"`
	Colorless int `json:"colorless"`
}

// LogEntry represents a log entry in the game
type LogEntry struct {
	Kind      string    `json:"kind"` // draw, play, tap, etc.
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
}

// Zone constants matching TypeScript ZoneId
const (
	ZoneLibraryStr     = "LIBRARY"
	ZoneHandStr        = "HAND"
	ZoneBattlefieldStr = "BATTLEFIELD"
	ZoneGraveyardStr   = "GRAVEYARD"
	ZoneExileStr       = "EXILE"
	ZoneStackStr       = "STACK"
	ZoneCommandStr     = "COMMAND"
)

// NewGameState creates a new empty game state
func NewGameState(gameID string, playerIDs []string, playerNames map[string]string) *GameState {
	state := &GameState{
		GameID:        gameID,
		TurnOrder:     append([]string(nil), playerIDs...),
		Players:       make(map[string]*Player),
		Battlefield:   make([]*Card, 0),
		Exile:         make([]*Card, 0),
		Stack:         make([]*Card, 0),
		Command:       make([]*Card, 0),
		Turn:          1,
		IsInitialized: false,
		Log:           make([]LogEntry, 0),
		MulliganType:  "london",
		FreeMulligans: 1, // Standard is 1 free mulligan
		StartedAt:     time.Now(),
	}

	// Initialize players
	for _, playerID := range playerIDs {
		name := playerNames[playerID]
		if name == "" {
			name = playerID
		}
		state.Players[playerID] = &Player{
			PlayerID:        playerID,
			Name:            name,
			Life:            20, // Default starting life
			Poison:          0,
			Energy:          0,
			LibraryCount:    0,
			HandCount:       0,
			Hand:            make([]*Card, 0),
			Library:         make([]*Card, 0),
			Graveyard:       make([]*Card, 0),
			ManaPool:        &ManaPool{},
			KeptHand:        false,
			MulliganCount:   0,
			RevealedTopCard: false,
		}
	}

	// Set first player as active
	if len(playerIDs) > 0 {
		state.ActivePlayerID = playerIDs[0]
		state.ActiveControlSeat = playerIDs[0]
	}

	return state
}

// ConvertLegacyCardToCard converts a LegacyCard to a Card (playtest engine card)
// From playtest-game.ts card structure
func ConvertLegacyCardToCard(card *LegacyCard) *Card {
	engineCounters := make([]Counter, 0)
	if card.Counters != nil {
		for name, counter := range card.Counters.GetAll() {
			if counter != nil && counter.Count > 0 {
				engineCounters = append(engineCounters, Counter{
					Name:  name,
					Count: counter.Count,
				})
			}
		}
	}

	return &Card{
		ID:                card.ID.String(),
		Name:              card.Name,
		DisplayName:       card.Name,
		OwnerID:           card.OwnerID.String(),
		ManaCost:          card.ManaCost,
		Type:              joinTypes(card.Types, card.Subtypes, card.Supertypes),
		SubTypes:          joinSlice(card.Subtypes),
		SuperTypes:        joinSlice(card.Supertypes),
		Color:             card.Color,
		Power:             card.Power,
		Toughness:         card.Toughness,
		Loyalty:           card.Loyalty,
		CardNumber:        card.CardNumber,
		ExpansionSet:      card.SetCode,
		Rarity:            card.Rarity,
		RulesText:         card.RulesText,
		Zone:              ZoneToString(card.Zone),
		ControllerID:      card.ControllerID.String(),
		Tapped:            card.Tapped,
		Flipped:           card.Flipped,
		Transformed:       card.Transformed,
		FaceDown:          card.FaceDown,
		Counters:          engineCounters,
		Attacking:         card.Attacking,
		Blocking:          card.Blocking,
		SummoningSickness: card.SummoningSickness,
		AttachedTo:        make([]string, 0),
	}
}

// ZoneToString converts Zone to string representation
func ZoneToString(zone Zone) string {
	switch zone {
	case ZoneLibrary:
		return ZoneLibraryStr
	case ZoneHand:
		return ZoneHandStr
	case ZoneBattlefield:
		return ZoneBattlefieldStr
	case ZoneGraveyard:
		return ZoneGraveyardStr
	case ZoneExile:
		return ZoneExileStr
	case ZoneStack:
		return ZoneStackStr
	case ZoneCommand:
		return ZoneCommandStr
	default:
		return ZoneLibraryStr
	}
}

// StringToZone converts string to Zone
func StringToZone(zoneStr string) Zone {
	switch zoneStr {
	case ZoneLibraryStr:
		return ZoneLibrary
	case ZoneHandStr:
		return ZoneHand
	case ZoneBattlefieldStr:
		return ZoneBattlefield
	case ZoneGraveyardStr:
		return ZoneGraveyard
	case ZoneExileStr:
		return ZoneExile
	case ZoneStackStr:
		return ZoneStack
	case ZoneCommandStr:
		return ZoneCommand
	default:
		return ZoneLibrary
	}
}

// Helper to join string slices
func joinSlice(slice []string) string {
	if len(slice) == 0 {
		return ""
	}
	result := ""
	for i, s := range slice {
		if i > 0 {
			result += " "
		}
		result += s
	}
	return result
}

// turnOrder returns the seating order of the game's players.
//
// It prefers the explicit TurnOrder slice captured at game creation. The sorted
// fallback exists only for GameState values that predate the field (a state
// restored from an older serialized snapshot, or built by a test that fills the
// Players map directly): it is arbitrary, but it is at least stable, which is
// the property NextTurn actually needs.
func (s *GameState) turnOrder() []string {
	if len(s.TurnOrder) > 0 {
		return s.TurnOrder
	}
	ids := make([]string, 0, len(s.Players))
	for pid := range s.Players {
		ids = append(ids, pid)
	}
	sort.Strings(ids)
	return ids
}
