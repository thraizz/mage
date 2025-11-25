package plugin

// CommanderBehavior implements GameBehavior for Commander format games.
// This handles commander-specific rules like:
// - Moving commanders to command zone at game start
// - Checking the 21 commander damage rule (903.10a)
// - Commander replacement effects (returning to command zone from graveyard/exile)
type CommanderBehavior struct {
	// CheckDamage determines if commander damage is tracked and checked.
	// Some variants (like Duel Commander since Nov 2016) don't use this rule.
	CheckDamage bool

	// DamageThreshold is the amount of combat damage from a single commander
	// that causes a player to lose. Default is 21.
	DamageThreshold int

	// AlsoHand determines if commanders going to hand are replaced to command zone.
	// This was part of older commander rules (pre-2017).
	AlsoHand bool

	// AlsoLibrary determines if commanders going to library are replaced to command zone.
	// This was part of older commander rules (pre-2017).
	AlsoLibrary bool
}

// NewCommanderBehavior creates a CommanderBehavior with default settings.
// Default: Check damage enabled, threshold 21, no hand/library replacement.
func NewCommanderBehavior() *CommanderBehavior {
	return &CommanderBehavior{
		CheckDamage:     true,
		DamageThreshold: 21,
		AlsoHand:        false,
		AlsoLibrary:     false,
	}
}

// Name returns the identifier for this behavior.
func (b *CommanderBehavior) Name() string {
	return "Commander"
}

// Init initializes commander-specific state at game start.
// Per rule 903.3 and GameCommanderImpl.init() in Java:
// - Cards designated as commanders are moved from sideboard to command zone
// - Commander watchers are set up for tracking damage and cast count
func (b *CommanderBehavior) Init(state GameState) error {
	// Move commanders to command zone for each player
	for _, playerID := range state.GetPlayerIDs() {
		commanders := state.GetCommanders(playerID)
		for _, commanderID := range commanders {
			if err := state.MoveCommanderToCommandZone(playerID, commanderID); err != nil {
				return err
			}
		}
	}
	return nil
}

// CheckStateBasedActions checks commander-specific state-based actions.
// Per rule 903.10a: A player that's been dealt 21 or more combat damage by the
// same commander over the course of the game loses the game.
// Returns true if any state-based actions were performed.
func (b *CommanderBehavior) CheckStateBasedActions(state GameState) bool {
	if !b.CheckDamage {
		return false
	}

	somethingHappened := false

	for _, playerID := range state.GetPlayerIDs() {
		if !state.IsPlayerInGame(playerID) {
			continue
		}

		// Check all commander damage entries for this player
		commanderDamage := state.GetAllCommanderDamage(playerID)
		for commanderID, damage := range commanderDamage {
			if damage >= b.DamageThreshold {
				commanderName := state.GetCardName(commanderID)
				if commanderName == "" {
					commanderName = commanderID
				}
				state.SetPlayerLost(playerID, "dealt "+string(rune(damage))+" commander damage by "+commanderName)
				somethingHappened = true
				break // Player already lost, no need to check more commanders
			}
		}
	}

	return somethingHappened
}

// IsCommanderDamageEnabled returns whether this behavior checks commander damage.
func (b *CommanderBehavior) IsCommanderDamageEnabled() bool {
	return b.CheckDamage
}

// GetDamageThreshold returns the commander damage threshold.
func (b *CommanderBehavior) GetDamageThreshold() int {
	return b.DamageThreshold
}
