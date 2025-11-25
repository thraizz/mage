package plugin

// GameRules defines the configuration for a game format.
// This struct contains all the numeric/boolean settings that determine
// how the game is set up (starting life, deck size, etc.)
type GameRules struct {
	// StartingLife is the amount of life each player starts with
	StartingLife int
	// MinimumDeckSize is the minimum number of cards required in the main deck
	MinimumDeckSize int
	// StartingHandSize is the number of cards drawn at game start
	StartingHandSize int
	// StartingPlayerSkipsDraw determines if the player who goes first skips their first draw step
	StartingPlayerSkipsDraw bool
}

// DefaultGameRules returns the standard rules for a typical MTG game
func DefaultGameRules() GameRules {
	return GameRules{
		StartingLife:            20,
		MinimumDeckSize:         60,
		StartingHandSize:        7,
		StartingPlayerSkipsDraw: true,
	}
}

// GameState is an interface representing the game state that behaviors can interact with.
// This is defined here to avoid circular dependencies - the actual implementation
// is in the game package.
type GameState interface {
	// GetPlayerIDs returns all player IDs in the game
	GetPlayerIDs() []string
	// GetPlayerLife returns a player's current life total
	GetPlayerLife(playerID string) int
	// SetPlayerLost marks a player as having lost the game
	SetPlayerLost(playerID string, reason string)
	// GetCommanderDamage returns the total commander damage dealt to a player by a specific commander
	GetCommanderDamage(playerID string, commanderID string) int
	// GetAllCommanderDamage returns all commander damage for a player (commanderID -> damage)
	GetAllCommanderDamage(playerID string) map[string]int
	// GetCommanders returns the commander card IDs for a player
	GetCommanders(playerID string) []string
	// MoveCommanderToCommandZone moves a commander from sideboard to command zone
	MoveCommanderToCommandZone(playerID string, commanderID string) error
	// IsPlayerInGame returns true if the player hasn't lost/left
	IsPlayerInGame(playerID string) bool
	// GetCardName returns the name of a card by ID
	GetCardName(cardID string) string
}

// GameBehavior defines format-specific behavior that can be plugged into the game engine.
// Different formats (Commander, Oathbreaker, etc.) can implement their own behaviors
// that get called at key points during the game.
type GameBehavior interface {
	// Name returns a unique identifier for this behavior
	Name() string

	// Init is called when the game starts, after players are set up but before
	// the first turn. This is where format-specific initialization happens
	// (e.g., moving commanders to command zone).
	Init(state GameState) error

	// CheckStateBasedActions is called during state-based action checks.
	// Returns true if any actions were performed, which will cause another
	// round of SBA checks.
	CheckStateBasedActions(state GameState) bool
}

// RulesProvider is optionally implemented by GameTypes to provide game rules.
// If a GameType doesn't implement this, DefaultGameRules() is used.
type RulesProvider interface {
	Rules() GameRules
}

// BehaviorProvider is optionally implemented by GameTypes to provide
// format-specific behaviors. If not implemented, no special behaviors are used.
type BehaviorProvider interface {
	Behaviors() []GameBehavior
}

// GetRulesForGameType returns the rules for a game type.
// If the game type implements RulesProvider, those rules are returned.
// Otherwise, DefaultGameRules() is returned.
func GetRulesForGameType(gt GameType) GameRules {
	if rp, ok := gt.(RulesProvider); ok {
		return rp.Rules()
	}
	return DefaultGameRules()
}

// GetBehaviorsForGameType returns the behaviors for a game type.
// If the game type implements BehaviorProvider, those behaviors are returned.
// Otherwise, an empty slice is returned.
func GetBehaviorsForGameType(gt GameType) []GameBehavior {
	if bp, ok := gt.(BehaviorProvider); ok {
		return bp.Behaviors()
	}
	return nil
}
