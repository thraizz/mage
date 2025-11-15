package rules

import (
	"testing"
)

// mockGameStateAccessor implements GameStateAccessor for testing
type mockGameStateAccessor struct {
	cards   map[string]CardInfo
	players map[string]PlayerInfo
	zones   map[string]int
}

func newMockGameStateAccessor() *mockGameStateAccessor {
	return &mockGameStateAccessor{
		cards:   make(map[string]CardInfo),
		players: make(map[string]PlayerInfo),
		zones:   make(map[string]int),
	}
}

func (m *mockGameStateAccessor) FindCard(cardID string) (CardInfo, bool) {
	card, ok := m.cards[cardID]
	return card, ok
}

func (m *mockGameStateAccessor) FindPlayer(playerID string) (PlayerInfo, bool) {
	player, ok := m.players[playerID]
	return player, ok
}

func (m *mockGameStateAccessor) IsCardInZone(cardID string, zone int) bool {
	cardZone, ok := m.zones[cardID]
	return ok && cardZone == zone
}

func (m *mockGameStateAccessor) GetCardZone(cardID string) (int, bool) {
	zone, ok := m.zones[cardID]
	return zone, ok
}

func TestLegalityChecker_ControllerValidation(t *testing.T) {
	mockState := newMockGameStateAccessor()
	mockState.players["player1"] = PlayerInfo{
		PlayerID: "player1",
		Name:     "Player 1",
		Life:     20,
		Lost:     false,
		Left:     false,
	}

	// Set up source card on stack for spell
	mockState.cards["card1"] = CardInfo{
		ID:           "card1",
		Name:         "Test Spell",
		Zone:         zoneStack,
		ControllerID: "player1",
	}
	mockState.zones["card1"] = zoneStack

	checker := NewLegalityChecker(mockState)

	// Valid item with controller in game
	item := StackItem{
		ID:          "spell1",
		Controller:  "player1",
		Description: "Test Spell",
		Kind:        StackItemKindSpell,
		SourceID:    "card1",
	}

	result := checker.CheckStackItemLegality(item)
	if !result.Legal {
		t.Errorf("Expected legal item, got illegal: %s", result.Reason)
	}

	// Invalid: controller not found
	item.Controller = "nonexistent"
	result = checker.CheckStackItemLegality(item)
	if result.Legal {
		t.Error("Expected illegal item (controller not found), got legal")
	}
	if result.Reason != "Controller not found" {
		t.Errorf("Expected reason 'Controller not found', got '%s'", result.Reason)
	}

	// Invalid: controller lost
	mockState.players["player2"] = PlayerInfo{
		PlayerID: "player2",
		Name:     "Player 2",
		Life:     0,
		Lost:     true,
		Left:     false,
	}
	item.Controller = "player2"
	result = checker.CheckStackItemLegality(item)
	if result.Legal {
		t.Error("Expected illegal item (controller lost), got legal")
	}
}

func TestLegalityChecker_SourceValidation(t *testing.T) {
	mockState := newMockGameStateAccessor()
	mockState.players["player1"] = PlayerInfo{
		PlayerID: "player1",
		Life:     20,
		Lost:     false,
		Left:     false,
	}

	checker := NewLegalityChecker(mockState)

	// Valid spell with source on stack
	mockState.cards["card1"] = CardInfo{
		ID:           "card1",
		Name:         "Lightning Bolt",
		Zone:         zoneStack,
		ControllerID: "player1",
	}
	mockState.zones["card1"] = zoneStack

	item := StackItem{
		ID:          "spell1",
		Controller:  "player1",
		Description: "Lightning Bolt",
		Kind:        StackItemKindSpell,
		SourceID:    "card1",
	}

	result := checker.CheckStackItemLegality(item)
	if !result.Legal {
		t.Errorf("Expected legal spell, got illegal: %s", result.Reason)
	}

	// Invalid: spell source not on stack
	mockState.zones["card1"] = zoneGraveyard
	mockState.cards["card1"] = CardInfo{
		ID:   "card1",
		Zone: zoneGraveyard,
	}
	result = checker.CheckStackItemLegality(item)
	if result.Legal {
		t.Error("Expected illegal spell (source not on stack), got legal")
	}

	// Valid: triggered ability can resolve even if source moved
	item.Kind = StackItemKindTriggered
	result = checker.CheckStackItemLegality(item)
	if !result.Legal {
		t.Errorf("Expected legal triggered ability (source can move), got illegal: %s", result.Reason)
	}
}

func TestLegalityChecker_TargetValidation(t *testing.T) {
	mockState := newMockGameStateAccessor()
	mockState.players["player1"] = PlayerInfo{
		PlayerID: "player1",
		Life:     20,
		Lost:     false,
		Left:     false,
	}

	// Set up source card on stack for spell
	mockState.cards["card1"] = CardInfo{
		ID:           "card1",
		Name:         "Target Spell",
		Zone:         zoneStack,
		ControllerID: "player1",
	}
	mockState.zones["card1"] = zoneStack

	checker := NewLegalityChecker(mockState)

	// Valid target on battlefield
	mockState.cards["target1"] = CardInfo{
		ID:   "target1",
		Name: "Creature",
		Zone: zoneBattlefield,
	}
	mockState.zones["target1"] = zoneBattlefield

	item := StackItem{
		ID:          "spell1",
		Controller:  "player1",
		Description: "Target Spell",
		Kind:        StackItemKindSpell,
		SourceID:    "card1",
		Metadata: map[string]string{
			"target": "target1",
		},
	}

	result := checker.CheckStackItemLegality(item)
	if !result.Legal {
		t.Errorf("Expected legal spell with valid target, got illegal: %s", result.Reason)
	}

	// Invalid: target moved to graveyard
	mockState.zones["target1"] = zoneGraveyard
	mockState.cards["target1"] = CardInfo{
		ID:   "target1",
		Zone: zoneGraveyard,
	}
	result = checker.CheckStackItemLegality(item)
	if result.Legal {
		t.Error("Expected illegal spell (target moved), got legal")
	}

	// Invalid: target not found
	item.Metadata["target"] = "nonexistent"
	result = checker.CheckStackItemLegality(item)
	if result.Legal {
		t.Error("Expected illegal spell (target not found), got legal")
	}
}

func TestLegalityChecker_CostValidation(t *testing.T) {
	mockState := newMockGameStateAccessor()
	checker := NewLegalityChecker(mockState)

	// Valid: costs paid
	item := StackItem{
		ID:          "spell1",
		Controller:  "player1",
		Description: "Test Spell",
		Kind:        StackItemKindSpell,
		Metadata: map[string]string{
			"costs_paid": "true",
		},
	}

	result := checker.CheckCostsPaid(item)
	if !result.Legal {
		t.Errorf("Expected legal (costs paid), got illegal: %s", result.Reason)
	}

	// Invalid: costs not paid
	item.Metadata["costs_paid"] = "false"
	result = checker.CheckCostsPaid(item)
	if result.Legal {
		t.Error("Expected illegal (costs not paid), got legal")
	}
}
