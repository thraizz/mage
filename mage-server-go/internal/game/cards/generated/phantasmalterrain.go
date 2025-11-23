package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Phantasmal Terrain", NewPhantasmalTerrain)
}

// NewPhantasmalTerrain creates a Phantasmal Terrain
// {U}{U} - ENCHANTMENT
func NewPhantasmalTerrain(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Phantasmal Terrain")
	card.ManaCost = "{U}{U}"
	card.Types = []string{"ENCHANTMENT"}
	card.Subtypes = []string{"AURA"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewEnchantAbility(card.ID, abilities.NewTargetRequirement(1, 1, abilities.NewCreatureTargetFilter()))
	card.AddAbility(ability0)
	ability1 := abilities.BuildSimpleManaAbility(card.ID, "G")
	card.AddAbility(ability1)
	ability2 := abilities.BuildSimpleManaAbility(card.ID, "W")
	card.AddAbility(ability2)
	ability3 := abilities.BuildSimpleManaAbility(card.ID, "R")
	card.AddAbility(ability3)
	ability4 := abilities.BuildSimpleManaAbility(card.ID, "U")
	card.AddAbility(ability4)
	ability5 := abilities.BuildSimpleManaAbility(card.ID, "B")
	card.AddAbility(ability5)
	ability6, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewAttachEffect(abilities.OutcomeDetriment)).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability6)
	return card, nil
}
