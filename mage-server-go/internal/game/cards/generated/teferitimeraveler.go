package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Teferi Time Raveler", NewTeferiTimeRaveler)
}

// NewTeferiTimeRaveler creates a Teferi Time Raveler
// {1}{W}{U} - PLANESWALKER
func NewTeferiTimeRaveler(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Teferi Time Raveler")
	card.ManaCost = "{1}{W}{U}"
	card.Types = []string{"PLANESWALKER"}
	card.Subtypes = []string{"TEFERI"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Loyalty = "4"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		// TODO: ReturnToHandTargetEffect with complex parameters
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	return card, nil
}
