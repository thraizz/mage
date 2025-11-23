package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Chandra The Firebrand", NewChandraTheFirebrand)
}

// NewChandraTheFirebrand creates a Chandra The Firebrand
// {3}{R} - PLANESWALKER
func NewChandraTheFirebrand(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Chandra The Firebrand")
	card.ManaCost = "{3}{R}"
	card.Types = []string{"PLANESWALKER"}
	card.Subtypes = []string{"CHANDRA"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Loyalty = "3"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewDamageEffect(1)).
		AddEffect(abilities.NewDamageEffect(6, true)).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	return card, nil
}
