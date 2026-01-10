package rules

import (
	"github.com/google/uuid"
)

// StateBasedActions implements MTG Rule 704 - State-Based Actions
// These are checks performed automatically whenever a player would receive priority
// or when a spell/ability finishes resolving.
type StateBasedActions struct{}

// NewStateBasedActions creates a new SBA checker
func NewStateBasedActions() *StateBasedActions {
	return &StateBasedActions{}
}

// Action represents a state-based action that needs to be performed
type Action interface {
	// Execute performs the action
	Execute(state GameStateReader) error

	// GetDescription returns a human-readable description
	GetDescription() string
}

// GameStateReader provides read access to game state for SBA checks
type GameStateReader interface {
	// Player queries
	GetPlayer(id uuid.UUID) (Player, bool)
	GetAllPlayers() []Player

	// Permanent queries
	GetPermanent(id uuid.UUID) (Permanent, bool)
	GetAllPermanents() []Permanent
	GetPermanentsControlledBy(playerID uuid.UUID) []Permanent

	// Stack queries
	IsOnStack(id uuid.UUID) bool
}

// Player represents player state for SBA checks
type Player struct {
	ID          uuid.UUID
	Life        int
	Poison      int
	LibrarySize int
}

// Permanent represents permanent state for SBA checks
type Permanent struct {
	ID            uuid.UUID
	Name          string
	Types         []string
	Subtypes      []string
	Power         int
	Toughness     int
	Damage        int
	Counters      map[string]int
	AttachedTo    uuid.UUID
	ControllerID  uuid.UUID
	IsToken       bool
	DamageSources map[uuid.UUID]bool // Sources that dealt damage to this permanent
}

// Check performs all state-based action checks and returns actions to execute
// RULES-LIGHT: Returns empty - players handle SBAs manually
// The original checks are kept below for reference/UI hints
func (sba *StateBasedActions) Check(state GameStateReader) []Action {
	// RULES-LIGHT: No automatic state-based actions
	// Players are trusted to handle lethal damage, 0-toughness creatures, etc. manually
	_ = state // Keep parameter for interface compatibility
	return []Action{}
}

// CheckAdvisory performs SBA checks and returns advisory hints (not enforced)
// This can be used by UI to show warnings like "Creature has lethal damage"
func (sba *StateBasedActions) CheckAdvisory(state GameStateReader) []Action {
	actions := []Action{}

	// 704.5a: Player with 0 or less life loses
	actions = append(actions, sba.checkPlayerLife(state)...)

	// 704.5b: Player with 10+ poison counters loses
	actions = append(actions, sba.checkPoisonCounters(state)...)

	// 704.5c: Player tried to draw from empty library loses
	// (This is handled during draw action, not here)

	// 704.5d: Token in non-battlefield zone ceases to exist
	// (Handled during zone change)

	// 704.5e: Copy of spell not on stack ceases to exist
	actions = append(actions, sba.checkSpellCopies(state)...)

	// 704.5f: Creature with toughness ≤ 0 is put into graveyard
	actions = append(actions, sba.checkCreatureToughness(state)...)

	// 704.5g: Creature with lethal damage is destroyed
	actions = append(actions, sba.checkLethalDamage(state)...)

	// 704.5h: Creature with deathtouch damage is destroyed
	actions = append(actions, sba.checkDeathtouchDamage(state)...)

	// 704.5i: Planeswalker with 0 loyalty is put into graveyard
	actions = append(actions, sba.checkPlaneswalkerLoyalty(state)...)

	// 704.5j: Two+ legendary permanents with same name
	actions = append(actions, sba.checkLegendRule(state)...)

	// 704.5k: Aura attached to illegal object/player
	actions = append(actions, sba.checkAuraAttachment(state)...)

	// 704.5m: Equipment/Fortification attached to illegal permanent
	actions = append(actions, sba.checkEquipmentAttachment(state)...)

	// 704.5n: Creature/planeswalker that's also Equipment/Fortification becomes unattached
	actions = append(actions, sba.checkCreatureEquipment(state)...)

	// 704.5q: +1/+1 and -1/-1 counters on same permanent
	actions = append(actions, sba.checkCounterAnnihilation(state)...)

	return actions
}

// checkPlayerLife checks Rule 704.5a

// checkPlayerLife checks Rule 704.5a
func (sba *StateBasedActions) checkPlayerLife(state GameStateReader) []Action {
	actions := []Action{}
	for _, player := range state.GetAllPlayers() {
		if player.Life <= 0 {
			actions = append(actions, &LoseGameAction{
				PlayerID: player.ID,
				Reason:   "life total of 0 or less",
			})
		}
	}
	return actions
}

// checkPoisonCounters checks Rule 704.5b
func (sba *StateBasedActions) checkPoisonCounters(state GameStateReader) []Action {
	actions := []Action{}
	for _, player := range state.GetAllPlayers() {
		if player.Poison >= 10 {
			actions = append(actions, &LoseGameAction{
				PlayerID: player.ID,
				Reason:   "10 or more poison counters",
			})
		}
	}
	return actions
}

// checkSpellCopies checks Rule 704.5e
func (sba *StateBasedActions) checkSpellCopies(state GameStateReader) []Action {
	// TODO: Implement when spell copy system exists
	return []Action{}
}

// checkCreatureToughness checks Rule 704.5f
func (sba *StateBasedActions) checkCreatureToughness(state GameStateReader) []Action {
	actions := []Action{}
	for _, permanent := range state.GetAllPermanents() {
		if sba.hasType(permanent, "CREATURE") {
			// Toughness includes all modifiers from continuous effects
			if permanent.Toughness <= 0 {
				actions = append(actions, &PutIntoGraveyardAction{
					PermanentID: permanent.ID,
					Reason:      "toughness 0 or less",
				})
			}
		}
	}
	return actions
}

// checkLethalDamage checks Rule 704.5g
func (sba *StateBasedActions) checkLethalDamage(state GameStateReader) []Action {
	actions := []Action{}
	for _, permanent := range state.GetAllPermanents() {
		if sba.hasType(permanent, "CREATURE") {
			// Rule 702.12: Indestructible permanents can't be destroyed
			if sba.hasAbility(permanent, "indestructible") {
				continue
			}

			if permanent.Damage >= permanent.Toughness && permanent.Toughness > 0 {
				actions = append(actions, &DestroyAction{
					PermanentID: permanent.ID,
					Reason:      "lethal damage",
				})
			}
		}
	}
	return actions
}

// checkDeathtouchDamage checks Rule 704.5h
func (sba *StateBasedActions) checkDeathtouchDamage(state GameStateReader) []Action {
	actions := []Action{}
	for _, permanent := range state.GetAllPermanents() {
		if sba.hasType(permanent, "CREATURE") && permanent.Damage > 0 {
			// Rule 702.12: Indestructible permanents can't be destroyed
			if sba.hasAbility(permanent, "indestructible") {
				continue
			}

			// Check if any damage source had deathtouch
			for sourceID := range permanent.DamageSources {
				if source, ok := state.GetPermanent(sourceID); ok {
					if sba.hasAbility(source, "deathtouch") {
						actions = append(actions, &DestroyAction{
							PermanentID: permanent.ID,
							Reason:      "deathtouch damage",
						})
						break
					}
				}
			}
		}
	}
	return actions
}

// checkPlaneswalkerLoyalty checks Rule 704.5i
func (sba *StateBasedActions) checkPlaneswalkerLoyalty(state GameStateReader) []Action {
	actions := []Action{}
	for _, permanent := range state.GetAllPermanents() {
		if sba.hasType(permanent, "PLANESWALKER") {
			loyalty := permanent.Counters["LOYALTY"]
			if loyalty <= 0 {
				actions = append(actions, &PutIntoGraveyardAction{
					PermanentID: permanent.ID,
					Reason:      "0 loyalty",
				})
			}
		}
	}
	return actions
}

// checkLegendRule checks Rule 704.5j
func (sba *StateBasedActions) checkLegendRule(state GameStateReader) []Action {
	actions := []Action{}

	// Group legendary permanents by controller and name
	legendsByController := make(map[uuid.UUID]map[string][]Permanent)

	for _, permanent := range state.GetAllPermanents() {
		if sba.hasType(permanent, "LEGENDARY") {
			controller := permanent.ControllerID
			if legendsByController[controller] == nil {
				legendsByController[controller] = make(map[string][]Permanent)
			}
			legendsByController[controller][permanent.Name] = append(
				legendsByController[controller][permanent.Name],
				permanent,
			)
		}
	}

	// For each controller, check for duplicate legendary names
	for controllerID, legends := range legendsByController {
		for name, permanents := range legends {
			if len(permanents) > 1 {
				actions = append(actions, &ChooseLegendaryAction{
					PlayerID:   controllerID,
					Name:       name,
					Permanents: permanents,
				})
			}
		}
	}

	return actions
}

// checkAuraAttachment checks Rule 704.5k
func (sba *StateBasedActions) checkAuraAttachment(state GameStateReader) []Action {
	actions := []Action{}
	for _, permanent := range state.GetAllPermanents() {
		if sba.hasSubtype(permanent, "Aura") && permanent.AttachedTo != uuid.Nil {
			// Check if attachment is still legal
			// TODO: Implement full legality check based on enchant ability
			if _, ok := state.GetPermanent(permanent.AttachedTo); !ok {
				// Attached object no longer exists
				actions = append(actions, &PutIntoGraveyardAction{
					PermanentID: permanent.ID,
					Reason:      "illegal attachment",
				})
			}
		}
	}
	return actions
}

// checkEquipmentAttachment checks Rule 704.5m
func (sba *StateBasedActions) checkEquipmentAttachment(state GameStateReader) []Action {
	actions := []Action{}
	for _, permanent := range state.GetAllPermanents() {
		if (sba.hasSubtype(permanent, "Equipment") || sba.hasSubtype(permanent, "Fortification")) &&
			permanent.AttachedTo != uuid.Nil {

			attached, ok := state.GetPermanent(permanent.AttachedTo)
			if !ok {
				// Attached permanent no longer exists
				actions = append(actions, &UnattachAction{
					PermanentID: permanent.ID,
					Reason:      "attached permanent no longer exists",
				})
			} else if sba.hasSubtype(permanent, "Equipment") && !sba.hasType(attached, "CREATURE") {
				// Equipment attached to non-creature
				actions = append(actions, &UnattachAction{
					PermanentID: permanent.ID,
					Reason:      "attached to non-creature",
				})
			}
		}
	}
	return actions
}

// checkCreatureEquipment checks Rule 704.5n
func (sba *StateBasedActions) checkCreatureEquipment(state GameStateReader) []Action {
	actions := []Action{}
	for _, permanent := range state.GetAllPermanents() {
		if (sba.hasType(permanent, "CREATURE") || sba.hasType(permanent, "PLANESWALKER")) &&
			(sba.hasSubtype(permanent, "Equipment") || sba.hasSubtype(permanent, "Fortification")) &&
			permanent.AttachedTo != uuid.Nil {

			actions = append(actions, &UnattachAction{
				PermanentID: permanent.ID,
				Reason:      "creature/planeswalker can't be attached",
			})
		}
	}
	return actions
}

// checkCounterAnnihilation checks Rule 704.5q
func (sba *StateBasedActions) checkCounterAnnihilation(state GameStateReader) []Action {
	actions := []Action{}
	for _, permanent := range state.GetAllPermanents() {
		plus := permanent.Counters["P1P1"]
		minus := permanent.Counters["M1M1"]

		if plus > 0 && minus > 0 {
			remove := plus
			if minus < plus {
				remove = minus
			}
			actions = append(actions, &RemoveCountersAction{
				PermanentID: permanent.ID,
				P1P1:        remove,
				M1M1:        remove,
			})
		}
	}
	return actions
}

// Helper functions

func (sba *StateBasedActions) hasType(permanent Permanent, cardType string) bool {
	for _, t := range permanent.Types {
		if t == cardType {
			return true
		}
	}
	return false
}

func (sba *StateBasedActions) hasSubtype(permanent Permanent, subtype string) bool {
	for _, st := range permanent.Subtypes {
		if st == subtype {
			return true
		}
	}
	return false
}

func (sba *StateBasedActions) hasAbility(permanent Permanent, ability string) bool {
	// TODO: Integrate with abilities system to check for keyword abilities
	// For now, this is a placeholder
	return false
}

// ========================================
// Action Types
// ========================================

// LoseGameAction represents a player losing the game
type LoseGameAction struct {
	PlayerID uuid.UUID
	Reason   string
}

func (a *LoseGameAction) Execute(state GameStateReader) error {
	// TODO: Implement player loss
	return nil
}

func (a *LoseGameAction) GetDescription() string {
	return "player loses the game: " + a.Reason
}

// DestroyAction represents destroying a permanent
type DestroyAction struct {
	PermanentID uuid.UUID
	Reason      string
}

func (a *DestroyAction) Execute(state GameStateReader) error {
	// TODO: Implement destroy
	return nil
}

func (a *DestroyAction) GetDescription() string {
	return "destroy permanent: " + a.Reason
}

// PutIntoGraveyardAction represents putting a permanent into graveyard
type PutIntoGraveyardAction struct {
	PermanentID uuid.UUID
	Reason      string
}

func (a *PutIntoGraveyardAction) Execute(state GameStateReader) error {
	// TODO: Implement zone change to graveyard
	return nil
}

func (a *PutIntoGraveyardAction) GetDescription() string {
	return "put into graveyard: " + a.Reason
}

// UnattachAction represents unattaching an Equipment/Aura/Fortification
type UnattachAction struct {
	PermanentID uuid.UUID
	Reason      string
}

func (a *UnattachAction) Execute(state GameStateReader) error {
	// TODO: Implement unattach
	return nil
}

func (a *UnattachAction) GetDescription() string {
	return "unattach: " + a.Reason
}

// RemoveCountersAction represents removing counters
type RemoveCountersAction struct {
	PermanentID uuid.UUID
	P1P1        int
	M1M1        int
}

func (a *RemoveCountersAction) Execute(state GameStateReader) error {
	// TODO: Implement counter removal
	return nil
}

func (a *RemoveCountersAction) GetDescription() string {
	return "remove +1/+1 and -1/-1 counters"
}

// ChooseLegendaryAction represents choosing which legendary to keep
type ChooseLegendaryAction struct {
	PlayerID   uuid.UUID
	Name       string
	Permanents []Permanent
}

func (a *ChooseLegendaryAction) Execute(state GameStateReader) error {
	// TODO: Implement legendary choice
	// Player chooses one to keep, others go to graveyard
	return nil
}

func (a *ChooseLegendaryAction) GetDescription() string {
	return "choose legendary permanent to keep: " + a.Name
}
