package game

import (
	"time"
)

// From playtest-game.ts lines 25-59: State structure

// EngineGameState represents the complete game state for the rules-light engine
type EngineGameState struct {
	GameID            string                   `json:"gameId"`
	ActiveControlSeat string                   `json:"activeControlSeat"` // Which player perspective is controlling
	Players           map[string]*EnginePlayer `json:"players"`
	Battlefield       []*EngineCard            `json:"battlefield"`
	Exile             []*EngineCard            `json:"exile"`
	Stack             []*EngineCard            `json:"stack"`
	Command           []*EngineCard            `json:"command"`
	Turn              int                      `json:"turn"`
	ActivePlayerID    string                   `json:"activePlayerId"`
	IsInitialized     bool                     `json:"isInitialized"`
	Log               []EngineLogEntry         `json:"log"`
	MulliganType      string                   `json:"mulliganType"` // "london"
	FreeMulligans     int                      `json:"freeMulligans"`
	StartedAt         time.Time                `json:"startedAt"`
}

// EnginePlayer represents a player in the rules-light engine
// From playtest-game.ts lines 25-40
type EnginePlayer struct {
	PlayerID        string          `json:"playerId"`
	Name            string          `json:"name"`
	Life            int             `json:"life"`
	Poison          int             `json:"poison"`
	Energy          int             `json:"energy"`
	LibraryCount    int             `json:"libraryCount"`
	HandCount       int             `json:"handCount"`
	Hand            []*EngineCard   `json:"hand"`
	Library         []*EngineCard   `json:"library"`
	Graveyard       []*EngineCard   `json:"graveyard"`
	ManaPool        *EngineManaPool `json:"manaPool"`
	KeptHand        bool            `json:"keptHand"`
	MulliganCount   int             `json:"mulliganCount"`
	RevealedTopCard bool            `json:"revealedTopCard"` // When true, top card is visible
}

// EngineCard represents a card in the rules-light engine
// Simplified from Card type - no rules enforcement
type EngineCard struct {
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
	Zone         string          `json:"zone"`
	ControllerID string          `json:"controllerId"`
	Tapped       bool            `json:"tapped"`
	Flipped      bool            `json:"flipped"`
	Transformed  bool            `json:"transformed"`
	FaceDown     bool            `json:"faceDown"`
	Counters     []EngineCounter `json:"counters"`

	// Combat state (manual tracking, no enforcement)
	Attacking         bool `json:"attacking"`
	Blocking          bool `json:"blocking"`
	SummoningSickness bool `json:"summoningSickness"`

	// Attachments
	AttachedTo []string `json:"attachedTo"`
}

// EngineCounter represents a counter on a card
type EngineCounter struct {
	Name  string `json:"name"`
	Count int    `json:"count"`
}

// EngineManaPool represents a player's mana pool (client-side tracking only)
type EngineManaPool struct {
	White     int `json:"white"`
	Blue      int `json:"blue"`
	Black     int `json:"black"`
	Red       int `json:"red"`
	Green     int `json:"green"`
	Colorless int `json:"colorless"`
}

// EngineLogEntry represents a log entry in the game
type EngineLogEntry struct {
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

// NewEngineGameState creates a new empty game state
func NewEngineGameState(gameID string, playerIDs []string, playerNames map[string]string) *EngineGameState {
	state := &EngineGameState{
		GameID:        gameID,
		Players:       make(map[string]*EnginePlayer),
		Battlefield:   make([]*EngineCard, 0),
		Exile:         make([]*EngineCard, 0),
		Stack:         make([]*EngineCard, 0),
		Command:       make([]*EngineCard, 0),
		Turn:          1,
		IsInitialized: false,
		Log:           make([]EngineLogEntry, 0),
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
		state.Players[playerID] = &EnginePlayer{
			PlayerID:        playerID,
			Name:            name,
			Life:            20, // Default starting life
			Poison:          0,
			Energy:          0,
			LibraryCount:    0,
			HandCount:       0,
			Hand:            make([]*EngineCard, 0),
			Library:         make([]*EngineCard, 0),
			Graveyard:       make([]*EngineCard, 0),
			ManaPool:        &EngineManaPool{},
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

// ConvertCardToEngineCard converts a Card to an EngineCard
// From playtest-game.ts card structure
func ConvertCardToEngineCard(card *Card) *EngineCard {
	engineCounters := make([]EngineCounter, 0)
	if card.Counters != nil {
		for name, counter := range card.Counters.GetAll() {
			if counter != nil && counter.Count > 0 {
				engineCounters = append(engineCounters, EngineCounter{
					Name:  name,
					Count: counter.Count,
				})
			}
		}
	}

	return &EngineCard{
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
