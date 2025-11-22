package abilities

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game/counters"
	"github.com/magefree/mage-server-go/internal/game/token"
)

// ========================================
// Extended Mock for Token Testing
// ========================================

// mockTokenPermanent extends mockPermanent with token-specific fields
type mockTokenPermanent struct {
	*mockPermanent
	token     *token.Token
	tapped    bool
	attacking bool
}

type mockTokenGameContext struct {
	*mockCounterGameContext
	createdTokens []*mockTokenPermanent
}

func newMockTokenGameContext() *mockTokenGameContext {
	return &mockTokenGameContext{
		mockCounterGameContext: newMockCounterGameContext(),
		createdTokens:          make([]*mockTokenPermanent, 0),
	}
}

func (m *mockTokenGameContext) CreateTokens(tok *token.Token, amount int, source uuid.UUID, tapped, attacking bool) ([]uuid.UUID, error) {
	created := make([]uuid.UUID, 0, amount)

	for i := 0; i < amount; i++ {
		// Create base permanent
		basePerm := &mockPermanent{
			id:       uuid.New(),
			counters: counters.NewCounters(),
		}

		// Create token permanent
		tokenPerm := &mockTokenPermanent{
			mockPermanent: basePerm,
			token:         tok.Copy(),
			tapped:        tapped,
			attacking:     attacking,
		}

		m.permanents[basePerm.id] = basePerm
		m.createdTokens = append(m.createdTokens, tokenPerm)
		created = append(created, basePerm.id)
	}

	return created, nil
}

// ========================================
// Token Creation Tests
// ========================================

func TestCreateTokenEffect_SingleToken(t *testing.T) {
	ctx := context.Background()
	game := newMockTokenGameContext()
	source := uuid.New()

	saproling := token.NewSaprolingToken()
	effect := NewCreateTokenEffect(saproling)

	err := effect.Apply(ctx, game, source, nil)
	if err != nil {
		t.Fatalf("Failed to create token: %v", err)
	}

	if len(game.createdTokens) != 1 {
		t.Errorf("Expected 1 token, got %d", len(game.createdTokens))
	}

	created := game.createdTokens[0]
	if created.token.Name != "Saproling Token" {
		t.Errorf("Expected Saproling Token, got %s", created.token.Name)
	}
	if created.token.Power != 1 || created.token.Toughness != 1 {
		t.Errorf("Expected 1/1, got %d/%d", created.token.Power, created.token.Toughness)
	}
	if !created.token.Color.Green {
		t.Error("Expected green token")
	}
}

func TestCreateTokenEffect_MultipleTokens(t *testing.T) {
	ctx := context.Background()
	game := newMockTokenGameContext()
	source := uuid.New()

	squirrel := token.NewSquirrelToken()
	effect := NewCreateTokenEffectAmount(squirrel, 5)

	err := effect.Apply(ctx, game, source, nil)
	if err != nil {
		t.Fatalf("Failed to create tokens: %v", err)
	}

	if len(game.createdTokens) != 5 {
		t.Errorf("Expected 5 tokens, got %d", len(game.createdTokens))
	}

	for _, created := range game.createdTokens {
		if created.token.Name != "Squirrel Token" {
			t.Errorf("Expected Squirrel Token, got %s", created.token.Name)
		}
	}
}

func TestCreateTokenEffect_TappedToken(t *testing.T) {
	ctx := context.Background()
	game := newMockTokenGameContext()
	source := uuid.New()

	soldier := token.NewSoldierToken()
	effect := NewCreateTokenEffectTapped(soldier, 1, true)

	err := effect.Apply(ctx, game, source, nil)
	if err != nil {
		t.Fatalf("Failed to create tapped token: %v", err)
	}

	if len(game.createdTokens) != 1 {
		t.Fatalf("Expected 1 token, got %d", len(game.createdTokens))
	}

	created := game.createdTokens[0]
	if !created.tapped {
		t.Error("Expected token to be tapped")
	}
}

func TestCreateTokenEffect_AttackingToken(t *testing.T) {
	ctx := context.Background()
	game := newMockTokenGameContext()
	source := uuid.New()

	goblin := token.NewGoblinToken()
	effect := NewCreateTokenEffectAttacking(goblin, 3, true, true)

	err := effect.Apply(ctx, game, source, nil)
	if err != nil {
		t.Fatalf("Failed to create attacking tokens: %v", err)
	}

	if len(game.createdTokens) != 3 {
		t.Fatalf("Expected 3 tokens, got %d", len(game.createdTokens))
	}

	for _, created := range game.createdTokens {
		if !created.tapped {
			t.Error("Expected attacking token to be tapped")
		}
		if !created.attacking {
			t.Error("Expected token to be attacking")
		}
	}
}

func TestCreateTokenEffect_WithCounters(t *testing.T) {
	ctx := context.Background()
	game := newMockTokenGameContext()
	source := uuid.New()

	zombie := token.NewZombieToken()
	effect := NewCreateTokenEffectAmount(zombie, 2).
		WithCounters(counters.CounterTypeP1P1, 2)

	err := effect.Apply(ctx, game, source, nil)
	if err != nil {
		t.Fatalf("Failed to create tokens with counters: %v", err)
	}

	if len(game.createdTokens) != 2 {
		t.Fatalf("Expected 2 tokens, got %d", len(game.createdTokens))
	}

	for _, created := range game.createdTokens {
		count := created.counters.GetCount("+1/+1")
		if count != 2 {
			t.Errorf("Expected 2 +1/+1 counters on token, got %d", count)
		}
	}
}

func TestCreateTokenEffect_GetLastAddedTokenIDs(t *testing.T) {
	ctx := context.Background()
	game := newMockTokenGameContext()
	source := uuid.New()

	treasure := token.NewTreasureToken()
	effect := NewCreateTokenEffectAmount(treasure, 3)

	err := effect.Apply(ctx, game, source, nil)
	if err != nil {
		t.Fatalf("Failed to create tokens: %v", err)
	}

	lastAdded := effect.GetLastAddedTokenIDs()
	if len(lastAdded) != 3 {
		t.Errorf("Expected 3 token IDs, got %d", len(lastAdded))
	}

	// Verify IDs match created tokens
	for i, id := range lastAdded {
		if game.createdTokens[i].id != id {
			t.Errorf("Token ID mismatch at index %d", i)
		}
	}
}

// ========================================
// Token Type Tests
// ========================================

func TestTreasureToken(t *testing.T) {
	tok := token.NewTreasureToken()

	if tok.Name != "Treasure Token" {
		t.Errorf("Expected Treasure Token, got %s", tok.Name)
	}

	hasArtifact := false
	for _, ct := range tok.CardTypes {
		if ct == token.CardTypeArtifact {
			hasArtifact = true
		}
	}
	if !hasArtifact {
		t.Error("Expected Treasure to be an artifact")
	}

	hasSubtype := false
	for _, st := range tok.Subtypes {
		if st == "Treasure" {
			hasSubtype = true
		}
	}
	if !hasSubtype {
		t.Error("Expected Treasure subtype")
	}
}

func TestSoldierTokenWithVigilance(t *testing.T) {
	tok := token.NewSoldierTokenVigilance()

	if tok.Power != 1 || tok.Toughness != 1 {
		t.Errorf("Expected 1/1, got %d/%d", tok.Power, tok.Toughness)
	}

	if !tok.Color.White {
		t.Error("Expected white token")
	}

	hasVigilance := false
	for _, ability := range tok.Abilities {
		if ability == "vigilance" {
			hasVigilance = true
		}
	}
	if !hasVigilance {
		t.Error("Expected vigilance ability")
	}
}

func TestAngelToken(t *testing.T) {
	tok := token.NewAngelToken()

	if tok.Power != 4 || tok.Toughness != 4 {
		t.Errorf("Expected 4/4, got %d/%d", tok.Power, tok.Toughness)
	}

	if !tok.Color.White {
		t.Error("Expected white token")
	}

	hasFlying := false
	for _, ability := range tok.Abilities {
		if ability == "flying" {
			hasFlying = true
		}
	}
	if !hasFlying {
		t.Error("Expected flying ability")
	}
}

func TestMerfolkTokenWithHexproof(t *testing.T) {
	tok := token.NewMerfolkTokenHexproof()

	if tok.Power != 1 || tok.Toughness != 1 {
		t.Errorf("Expected 1/1, got %d/%d", tok.Power, tok.Toughness)
	}

	if !tok.Color.Blue {
		t.Error("Expected blue token")
	}

	hasHexproof := false
	for _, ability := range tok.Abilities {
		if ability == "hexproof" {
			hasHexproof = true
		}
	}
	if !hasHexproof {
		t.Error("Expected hexproof ability")
	}
}

func TestDragonToken(t *testing.T) {
	tok := token.NewDragonToken()

	if tok.Power != 5 || tok.Toughness != 5 {
		t.Errorf("Expected 5/5, got %d/%d", tok.Power, tok.Toughness)
	}

	if !tok.Color.Red {
		t.Error("Expected red token")
	}

	hasFlying := false
	for _, ability := range tok.Abilities {
		if ability == "flying" {
			hasFlying = true
		}
	}
	if !hasFlying {
		t.Error("Expected flying ability")
	}
}

func TestThopterToken(t *testing.T) {
	tok := token.NewThopterToken()

	if tok.Power != 1 || tok.Toughness != 1 {
		t.Errorf("Expected 1/1, got %d/%d", tok.Power, tok.Toughness)
	}

	hasArtifact := false
	hasCreature := false
	for _, ct := range tok.CardTypes {
		if ct == token.CardTypeArtifact {
			hasArtifact = true
		}
		if ct == token.CardTypeCreature {
			hasCreature = true
		}
	}
	if !hasArtifact || !hasCreature {
		t.Error("Expected Thopter to be an artifact creature")
	}

	hasFlying := false
	for _, ability := range tok.Abilities {
		if ability == "flying" {
			hasFlying = true
		}
	}
	if !hasFlying {
		t.Error("Expected flying ability")
	}
}

func TestZombieTokenWithDecayed(t *testing.T) {
	tok := token.NewZombieTokenDecayed()

	if tok.Power != 2 || tok.Toughness != 2 {
		t.Errorf("Expected 2/2, got %d/%d", tok.Power, tok.Toughness)
	}

	if !tok.Color.Black {
		t.Error("Expected black token")
	}

	hasDecayed := false
	for _, ability := range tok.Abilities {
		if ability == "decayed" {
			hasDecayed = true
		}
	}
	if !hasDecayed {
		t.Error("Expected decayed ability")
	}
}

func TestClueToken(t *testing.T) {
	tok := token.NewClueToken()

	hasArtifact := false
	for _, ct := range tok.CardTypes {
		if ct == token.CardTypeArtifact {
			hasArtifact = true
		}
	}
	if !hasArtifact {
		t.Error("Expected Clue to be an artifact")
	}

	hasSubtype := false
	for _, st := range tok.Subtypes {
		if st == "Clue" {
			hasSubtype = true
		}
	}
	if !hasSubtype {
		t.Error("Expected Clue subtype")
	}
}

func TestFoodToken(t *testing.T) {
	tok := token.NewFoodToken()

	hasArtifact := false
	for _, ct := range tok.CardTypes {
		if ct == token.CardTypeArtifact {
			hasArtifact = true
		}
	}
	if !hasArtifact {
		t.Error("Expected Food to be an artifact")
	}

	hasSubtype := false
	for _, st := range tok.Subtypes {
		if st == "Food" {
			hasSubtype = true
		}
	}
	if !hasSubtype {
		t.Error("Expected Food subtype")
	}
}

// ========================================
// Effect Description Tests
// ========================================

func TestCreateTokenEffect_GetDescription_Single(t *testing.T) {
	saproling := token.NewSaprolingToken()
	effect := NewCreateTokenEffect(saproling)

	desc := effect.GetDescription()
	expected := "create a 1/1 green Saproling creature token"
	if desc != expected {
		t.Errorf("Expected description '%s', got '%s'", expected, desc)
	}
}

func TestCreateTokenEffect_GetDescription_Multiple(t *testing.T) {
	squirrel := token.NewSquirrelToken()
	effect := NewCreateTokenEffectAmount(squirrel, 3)

	desc := effect.GetDescription()
	expected := "create 3 1/1 green Squirrel creature token"
	if desc != expected {
		t.Errorf("Expected description '%s', got '%s'", expected, desc)
	}
}

func TestCreateTokenEffect_GetDescription_Tapped(t *testing.T) {
	soldier := token.NewSoldierToken()
	effect := NewCreateTokenEffectTapped(soldier, 2, true)

	desc := effect.GetDescription()
	expected := "create 2 1/1 white Soldier creature token tapped"
	if desc != expected {
		t.Errorf("Expected description '%s', got '%s'", expected, desc)
	}
}

func TestCreateTokenEffect_GetDescription_Attacking(t *testing.T) {
	goblin := token.NewGoblinToken()
	effect := NewCreateTokenEffectAttacking(goblin, 1, true, true)

	desc := effect.GetDescription()
	expected := "create a 1/1 red Goblin creature token tapped and attacking"
	if desc != expected {
		t.Errorf("Expected description '%s', got '%s'", expected, desc)
	}
}

func TestCreateTokenEffect_GetDescription_WithCounters(t *testing.T) {
	zombie := token.NewZombieToken()
	effect := NewCreateTokenEffectAmount(zombie, 1).
		WithCounters(counters.CounterTypeP1P1, 2)

	desc := effect.GetDescription()
	expected := "create a 2/2 black Zombie creature token with 2 +1/+1 counters on it"
	if desc != expected {
		t.Errorf("Expected description '%s', got '%s'", expected, desc)
	}
}

// ========================================
// Integration Tests
// ========================================

func TestCreateTokenEffect_ComplexScenario(t *testing.T) {
	ctx := context.Background()
	game := newMockTokenGameContext()
	source := uuid.New()

	// Create 3 soldier tokens with vigilance, tapped and attacking, each with 1 +1/+1 counter
	soldier := token.NewSoldierTokenVigilance()
	effect := NewCreateTokenEffectAttacking(soldier, 3, true, true).
		WithCounters(counters.CounterTypeP1P1, 1)

	err := effect.Apply(ctx, game, source, nil)
	if err != nil {
		t.Fatalf("Failed to create complex tokens: %v", err)
	}

	if len(game.createdTokens) != 3 {
		t.Fatalf("Expected 3 tokens, got %d", len(game.createdTokens))
	}

	for i, created := range game.createdTokens {
		// Verify token properties
		if created.token.Name != "Soldier Token" {
			t.Errorf("Token %d: Expected Soldier Token, got %s", i, created.token.Name)
		}
		if created.token.Power != 1 || created.token.Toughness != 1 {
			t.Errorf("Token %d: Expected 1/1, got %d/%d", i, created.token.Power, created.token.Toughness)
		}
		if !created.token.Color.White {
			t.Errorf("Token %d: Expected white token", i)
		}

		// Verify vigilance
		hasVigilance := false
		for _, ability := range created.token.Abilities {
			if ability == "vigilance" {
				hasVigilance = true
			}
		}
		if !hasVigilance {
			t.Errorf("Token %d: Expected vigilance ability", i)
		}

		// Verify tapped and attacking
		if !created.tapped {
			t.Errorf("Token %d: Expected token to be tapped", i)
		}
		if !created.attacking {
			t.Errorf("Token %d: Expected token to be attacking", i)
		}

		// Verify counters
		count := created.counters.GetCount("+1/+1")
		if count != 1 {
			t.Errorf("Token %d: Expected 1 +1/+1 counter, got %d", i, count)
		}
	}

	// Verify messages
	if len(game.messages) != 1 {
		t.Fatalf("Expected 1 message, got %d", len(game.messages))
	}
	expectedMsg := "Created 3 1/1 white Soldier creature token with vigilance"
	if game.messages[0] != expectedMsg {
		t.Errorf("Expected message '%s', got '%s'", expectedMsg, game.messages[0])
	}
}

func TestTokenCopy(t *testing.T) {
	original := token.NewAngelToken()
	copy := original.Copy()

	// Verify copy has same properties
	if copy.Name != original.Name {
		t.Error("Copy has different name")
	}
	if copy.Power != original.Power || copy.Toughness != original.Toughness {
		t.Error("Copy has different power/toughness")
	}
	if copy.Color != original.Color {
		t.Error("Copy has different color")
	}

	// Verify copy has different UUID
	if copy.ID == original.ID {
		t.Error("Copy has same UUID as original")
	}

	// Verify modifying copy doesn't affect original
	copy.Power = 10
	if original.Power == 10 {
		t.Error("Modifying copy affected original")
	}
}
