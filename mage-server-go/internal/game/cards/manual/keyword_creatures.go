package manual

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

// Register creatures with keyword abilities on package import
func init() {
	cards.Register("Serra Angel", NewSerraAngel)
	cards.Register("Serra's Guardian", NewSerrasGuardian)
	cards.Register("Vampire Nighthawk", NewVampireNighthawk)
}

// NewSerraAngel creates a Serra Angel
// {3}{W}{W} - Creature — Angel - 4/4
// Flying, Vigilance
func NewSerraAngel(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Serra Angel")
	card.ManaCost = "{3}{W}{W}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"ANGEL"}
	card.Power = "4"
	card.Toughness = "4"
	card.SetCode = "M15"
	card.Rarity = "uncommon"
	card.RulesText = "Flying, Vigilance"

	// Add flying keyword ability
	flyingAbility := abilities.NewKeywordAbility(card.ID, abilities.KeywordFlying)
	card.AddAbility(flyingAbility)

	// Add vigilance keyword ability
	vigilanceAbility := abilities.NewKeywordAbility(card.ID, abilities.KeywordVigilance)
	card.AddAbility(vigilanceAbility)

	return card, nil
}

// NewSerrasGuardian creates a Serra's Guardian
// {4}{W} - Creature — Angel - 3/4
// Flying, Vigilance
func NewSerrasGuardian(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Serra's Guardian")
	card.ManaCost = "{4}{W}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"ANGEL"}
	card.Power = "3"
	card.Toughness = "4"
	card.SetCode = "M20"
	card.Rarity = "common"
	card.RulesText = "Flying, Vigilance"

	// Add flying keyword ability
	flyingAbility := abilities.NewKeywordAbility(card.ID, abilities.KeywordFlying)
	card.AddAbility(flyingAbility)

	// Add vigilance keyword ability
	vigilanceAbility := abilities.NewKeywordAbility(card.ID, abilities.KeywordVigilance)
	card.AddAbility(vigilanceAbility)

	return card, nil
}

// NewVampireNighthawk creates a Vampire Nighthawk
// {1}{B}{B} - Creature — Vampire Shaman - 2/3
// Flying, Deathtouch, Lifelink
func NewVampireNighthawk(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Vampire Nighthawk")
	card.ManaCost = "{1}{B}{B}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"VAMPIRE", "SHAMAN"}
	card.Power = "2"
	card.Toughness = "3"
	card.SetCode = "M13"
	card.Rarity = "uncommon"
	card.RulesText = "Flying, Deathtouch, Lifelink"

	// Add flying keyword ability
	flyingAbility := abilities.NewKeywordAbility(card.ID, abilities.KeywordFlying)
	card.AddAbility(flyingAbility)

	// Add deathtouch keyword ability
	deathtouchAbility := abilities.NewKeywordAbility(card.ID, abilities.KeywordDeathtouch)
	card.AddAbility(deathtouchAbility)

	// Add lifelink keyword ability
	lifelinkAbility := abilities.NewKeywordAbility(card.ID, abilities.KeywordLifelink)
	card.AddAbility(lifelinkAbility)

	return card, nil
}
