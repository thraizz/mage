package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Nyleas Presence", NewNyleasPresence)
}

// NewNyleasPresence creates a Nyleas Presence
// {1}{G} - ENCHANTMENT
func NewNyleasPresence(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Nyleas Presence")
	card.ManaCost = "{1}{G}"
	card.Types = []string{"ENCHANTMENT"}
	card.Subtypes = []string{"AURA"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewEnchantAbility(card.ID, abilities.NewTargetRequirement(1, 1, abilities.NewCreatureTargetFilter()))
	card.AddAbility(ability0)
	ability1 := abilities.BuildSimpleManaAbility(card.ID, "W")
	card.AddAbility(ability1)
	ability2 := abilities.BuildSimpleManaAbility(card.ID, "U")
	card.AddAbility(ability2)
	ability3 := abilities.BuildSimpleManaAbility(card.ID, "B")
	card.AddAbility(ability3)
	ability4 := abilities.BuildSimpleManaAbility(card.ID, "R")
	card.AddAbility(ability4)
	ability5 := abilities.BuildSimpleManaAbility(card.ID, "G")
	card.AddAbility(ability5)
	ability6, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewAttachEffect(abilities.OutcomeAddAbility)).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability6)
	ability7, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewDrawCardsEffect(1)).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability7)
	return card, nil
}
