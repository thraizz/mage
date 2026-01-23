package game

import (
	"fmt"
	"testing"

	"go.uber.org/zap/zaptest"
)

// CombatTestHarness provides utilities for setting up and running combat tests
type CombatTestHarness struct {
	t       *testing.T
	engine  *MageEngine
	gameID  string
	players []string
}

// NewCombatTestHarness creates a new test harness for combat scenarios
func NewCombatTestHarness(t *testing.T, gameID string, players []string) *CombatTestHarness {
	logger := zaptest.NewLogger(t)
	engine := NewMageEngine(logger)

	if err := engine.StartGame(gameID, players, "Duel"); err != nil {
		t.Fatalf("failed to start game: %v", err)
	}

	return &CombatTestHarness{
		t:       t,
		engine:  engine,
		gameID:  gameID,
		players: players,
	}
}

// GetGameState returns the internal game state for direct manipulation
func (h *CombatTestHarness) GetGameState() *engineGameState {
	h.engine.mu.RLock()
	gameState := h.engine.games[h.gameID]
	h.engine.mu.RUnlock()
	return gameState
}

// CreatureSpec defines the properties of a test creature
type CreatureSpec struct {
	ID         string
	Name       string
	Power      string
	Toughness  string
	Controller string
	Abilities  []string // e.g., "FlyingAbility", "VigilanceAbility"
	Tapped     bool
}

// CreateCreature adds a creature to the battlefield
func (h *CombatTestHarness) CreateCreature(spec CreatureSpec) string {
	gameState := h.GetGameState()

	gameState.mu.Lock()
	defer gameState.mu.Unlock()

	abilities := make([]EngineAbilityView, 0, len(spec.Abilities))
	for _, abilityID := range spec.Abilities {
		abilities = append(abilities, EngineAbilityView{ID: abilityID})
	}

	card := &internalCard{
		ID:           spec.ID,
		Name:         spec.Name,
		Type:         "Creature",
		Zone:         zoneBattlefield,
		OwnerID:      spec.Controller,
		ControllerID: spec.Controller,
		Power:        spec.Power,
		Toughness:    spec.Toughness,
		Tapped:       spec.Tapped,
		Abilities:    abilities,
	}

	gameState.cards[spec.ID] = card

	return spec.ID
}

// CreateAttacker creates a simple attacker creature
func (h *CombatTestHarness) CreateAttacker(id, name, controller, power, toughness string) string {
	return h.CreateCreature(CreatureSpec{
		ID:         id,
		Name:       name,
		Power:      power,
		Toughness:  toughness,
		Controller: controller,
	})
}

// CreateBlocker creates a simple blocker creature
func (h *CombatTestHarness) CreateBlocker(id, name, controller, power, toughness string) string {
	return h.CreateCreature(CreatureSpec{
		ID:         id,
		Name:       name,
		Power:      power,
		Toughness:  toughness,
		Controller: controller,
	})
}

// SetupCombat initializes combat for the given attacking player
func (h *CombatTestHarness) SetupCombat(attackingPlayer string) {
	if err := h.engine.ResetCombat(h.gameID); err != nil {
		h.t.Fatalf("failed to reset combat: %v", err)
	}

	if err := h.engine.SetAttacker(h.gameID, attackingPlayer); err != nil {
		h.t.Fatalf("failed to set attacker: %v", err)
	}

	if err := h.engine.SetDefenders(h.gameID); err != nil {
		h.t.Fatalf("failed to set defenders: %v", err)
	}
}

// DeclareAttacker declares a single creature as an attacker
func (h *CombatTestHarness) DeclareAttacker(creatureID, defenderID, controllerID string) {
	if err := h.engine.DeclareAttacker(h.gameID, creatureID, defenderID, controllerID); err != nil {
		h.t.Fatalf("failed to declare attacker %s: %v", creatureID, err)
	}
}

// DeclareBlocker declares a single creature as a blocker
func (h *CombatTestHarness) DeclareBlocker(blockerID, attackerID, controllerID string) {
	if err := h.engine.DeclareBlocker(h.gameID, blockerID, attackerID, controllerID); err != nil {
		h.t.Fatalf("failed to declare blocker %s: %v", blockerID, err)
	}
}

// AcceptBlockers finalizes blocker declarations
func (h *CombatTestHarness) AcceptBlockers() {
	if err := h.engine.AcceptBlockers(h.gameID); err != nil {
		h.t.Fatalf("failed to accept blockers: %v", err)
	}
}

// AssignDamage assigns combat damage (firstStrike true for first strike damage)
func (h *CombatTestHarness) AssignDamage(firstStrike bool) {
	if err := h.engine.AssignCombatDamage(h.gameID, firstStrike); err != nil {
		h.t.Fatalf("failed to assign combat damage: %v", err)
	}
}

// ApplyDamage applies assigned combat damage
func (h *CombatTestHarness) ApplyDamage() {
	if err := h.engine.ApplyCombatDamage(h.gameID); err != nil {
		h.t.Fatalf("failed to apply combat damage: %v", err)
	}
}

// EndCombat ends the combat phase
func (h *CombatTestHarness) EndCombat() {
	if err := h.engine.EndCombat(h.gameID); err != nil {
		h.t.Fatalf("failed to end combat: %v", err)
	}
}

// GetPlayerLife returns the current life total for a player
func (h *CombatTestHarness) GetPlayerLife(playerID string) int {
	gameState := h.GetGameState()
	gameState.mu.RLock()
	defer gameState.mu.RUnlock()

	player, exists := gameState.players[playerID]
	if !exists {
		h.t.Fatalf("player %s not found", playerID)
	}
	return player.Life
}

// GetCreatureDamage returns the damage marked on a creature
func (h *CombatTestHarness) GetCreatureDamage(creatureID string) int {
	gameState := h.GetGameState()
	gameState.mu.RLock()
	defer gameState.mu.RUnlock()

	card, exists := gameState.cards[creatureID]
	if !exists {
		h.t.Fatalf("creature %s not found", creatureID)
	}
	return card.Damage
}

// IsCreatureDead checks if a creature is in the graveyard
func (h *CombatTestHarness) IsCreatureDead(creatureID string) bool {
	gameState := h.GetGameState()
	gameState.mu.RLock()
	defer gameState.mu.RUnlock()

	card, exists := gameState.cards[creatureID]
	if !exists {
		return true // Card doesn't exist, consider it dead
	}
	return card.Zone == zoneGraveyard
}

// IsCreatureTapped checks if a creature is tapped
func (h *CombatTestHarness) IsCreatureTapped(creatureID string) bool {
	gameState := h.GetGameState()
	gameState.mu.RLock()
	defer gameState.mu.RUnlock()

	card, exists := gameState.cards[creatureID]
	if !exists {
		h.t.Fatalf("creature %s not found", creatureID)
	}
	return card.Tapped
}

// IsCreatureAttacking checks if a creature is currently attacking
func (h *CombatTestHarness) IsCreatureAttacking(creatureID string) bool {
	gameState := h.GetGameState()
	gameState.mu.RLock()
	defer gameState.mu.RUnlock()

	card, exists := gameState.cards[creatureID]
	if !exists {
		h.t.Fatalf("creature %s not found", creatureID)
	}
	return card.Attacking
}

// IsCreatureBlocking checks if a creature is currently blocking
func (h *CombatTestHarness) IsCreatureBlocking(creatureID string) bool {
	gameState := h.GetGameState()
	gameState.mu.RLock()
	defer gameState.mu.RUnlock()

	card, exists := gameState.cards[creatureID]
	if !exists {
		h.t.Fatalf("creature %s not found", creatureID)
	}
	return card.Blocking
}

// AssertPlayerLife asserts that a player has the expected life total
func (h *CombatTestHarness) AssertPlayerLife(playerID string, expectedLife int) {
	actualLife := h.GetPlayerLife(playerID)
	if actualLife != expectedLife {
		h.t.Errorf("expected %s life to be %d, got %d", playerID, expectedLife, actualLife)
	}
}

// AssertCreatureDamage asserts that a creature has the expected damage
func (h *CombatTestHarness) AssertCreatureDamage(creatureID string, expectedDamage int) {
	actualDamage := h.GetCreatureDamage(creatureID)
	if actualDamage != expectedDamage {
		h.t.Errorf("expected %s damage to be %d, got %d", creatureID, expectedDamage, actualDamage)
	}
}

// AssertCreatureDead asserts that a creature is in the graveyard
func (h *CombatTestHarness) AssertCreatureDead(creatureID string) {
	if !h.IsCreatureDead(creatureID) {
		h.t.Errorf("expected %s to be dead (in graveyard)", creatureID)
	}
}

// AssertCreatureAlive asserts that a creature is still on the battlefield
func (h *CombatTestHarness) AssertCreatureAlive(creatureID string) {
	if h.IsCreatureDead(creatureID) {
		h.t.Errorf("expected %s to be alive (on battlefield)", creatureID)
	}
}

// AssertCreatureTapped asserts that a creature is tapped
func (h *CombatTestHarness) AssertCreatureTapped(creatureID string, shouldBeTapped bool) {
	isTapped := h.IsCreatureTapped(creatureID)
	if isTapped != shouldBeTapped {
		if shouldBeTapped {
			h.t.Errorf("expected %s to be tapped", creatureID)
		} else {
			h.t.Errorf("expected %s to be untapped", creatureID)
		}
	}
}

// RunFullCombat executes a complete combat sequence with given attackers and blockers
// attackers is a map of creatureID -> defenderID
// blockers is a map of blockerID -> attackerID
func (h *CombatTestHarness) RunFullCombat(attackingPlayer string, attackers map[string]string, blockers map[string]string) {
	// Setup combat
	h.SetupCombat(attackingPlayer)

	// Declare attackers
	for creatureID, defenderID := range attackers {
		h.DeclareAttacker(creatureID, defenderID, attackingPlayer)
	}

	// Declare blockers
	for blockerID, attackerID := range blockers {
		// Find the controller of the blocker
		gameState := h.GetGameState()
		gameState.mu.RLock()
		blocker := gameState.cards[blockerID]
		controllerID := blocker.ControllerID
		gameState.mu.RUnlock()

		h.DeclareBlocker(blockerID, attackerID, controllerID)
	}

	// Accept blockers if any were declared
	if len(blockers) > 0 {
		h.AcceptBlockers()
	}

	// Check for first strike
	hasFirstStrike, _ := h.engine.HasFirstOrDoubleStrike(h.gameID)

	// Assign and apply first strike damage if needed
	if hasFirstStrike {
		h.AssignDamage(true)
		h.ApplyDamage()
	}

	// Assign and apply normal damage
	h.AssignDamage(false)
	h.ApplyDamage()

	// End combat
	h.EndCombat()
}

// Debug prints the current state of a creature (for debugging tests)
func (h *CombatTestHarness) Debug(creatureID string) {
	gameState := h.GetGameState()
	gameState.mu.RLock()
	defer gameState.mu.RUnlock()

	card, exists := gameState.cards[creatureID]
	if !exists {
		fmt.Printf("Creature %s not found\n", creatureID)
		return
	}

	fmt.Printf("Creature %s (%s):\n", creatureID, card.Name)
	fmt.Printf("  P/T: %s/%s\n", card.Power, card.Toughness)
	fmt.Printf("  Damage: %d\n", card.Damage)
	fmt.Printf("  Zone: %d\n", card.Zone)
	fmt.Printf("  Tapped: %v\n", card.Tapped)
	fmt.Printf("  Attacking: %v\n", card.Attacking)
	fmt.Printf("  Blocking: %v\n", card.Blocking)
	fmt.Printf("  Abilities: %v\n", card.Abilities)
}
