package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Honden Of Seeing Winds", NewHondenOfSeeingWinds)
}

// NewHondenOfSeeingWinds creates a Honden Of Seeing Winds
// {4}{U} - ENCHANTMENT
func NewHondenOfSeeingWinds(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Honden Of Seeing Winds")
	card.ManaCost = "{4}{U}"
	card.Types = []string{"ENCHANTMENT"}
	card.Subtypes = []string{"SHRINE"}
	card.Supertypes = []string{"LEGENDARY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewDrawCardsEffect(xValue)).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	return card, nil
}
