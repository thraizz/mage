package manual

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

// Register creatures with activated abilities on package import
func init() {
	cards.Register("Llanowar Elves", NewLlanowarElves)
	cards.Register("Prodigal Pyromancer", NewProdigalPyromancer)
	cards.Register("Prodigal Sorcerer", NewProdigalSorcerer)
}

// NewLlanowarElves creates a Llanowar Elves
// {G} - Creature — Elf Druid - 1/1
// {T}: Add {G}
func NewLlanowarElves(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Llanowar Elves")
	card.ManaCost = "{G}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"ELF", "DRUID"}
	card.Power = "1"
	card.Toughness = "1"
	card.SetCode = "M10"
	card.Rarity = "common"
	card.RulesText = "{T}: Add {G}."

	// Add mana ability: {T}: Add {G}
	ability := abilities.BuildSimpleManaAbility(card.ID, "G")
	card.AddAbility(ability)

	return card, nil
}

// NewProdigalPyromancer creates a Prodigal Pyromancer
// {2}{R} - Creature — Human Wizard - 1/1
// {T}: Prodigal Pyromancer deals 1 damage to any target
func NewProdigalPyromancer(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Prodigal Pyromancer")
	card.ManaCost = "{2}{R}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"HUMAN", "WIZARD"}
	card.Power = "1"
	card.Toughness = "1"
	card.SetCode = "M10"
	card.Rarity = "uncommon"
	card.RulesText = "{T}: Prodigal Pyromancer deals 1 damage to any target."

	// Add damage ability: {T}: Deal 1 damage to any target
	ability := abilities.BuildSimpleDamageAbility(card.ID, 1)
	card.AddAbility(ability)

	return card, nil
}

// NewProdigalSorcerer creates a Prodigal Sorcerer (Tim)
// {2}{U} - Creature — Human Wizard - 1/1
// {T}: Prodigal Sorcerer deals 1 damage to any target
func NewProdigalSorcerer(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Prodigal Sorcerer")
	card.ManaCost = "{2}{U}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"HUMAN", "WIZARD"}
	card.Power = "1"
	card.Toughness = "1"
	card.SetCode = "M10"
	card.Rarity = "common"
	card.RulesText = "{T}: Prodigal Sorcerer deals 1 damage to any target."

	// Add damage ability: {T}: Deal 1 damage to any target
	ability := abilities.BuildSimpleDamageAbility(card.ID, 1)
	card.AddAbility(ability)

	return card, nil
}
