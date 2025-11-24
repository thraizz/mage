package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Tibalt The Fiend Blooded", NewTibaltTheFiendBlooded)
}

// NewTibaltTheFiendBlooded creates a Tibalt The Fiend Blooded
// {R}{R} - PLANESWALKER
func NewTibaltTheFiendBlooded(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Tibalt The Fiend Blooded")
	card.ManaCost = "{R}{R}"
	card.Types = []string{"PLANESWALKER"}
	card.Subtypes = []string{"TIBALT"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Loyalty = "2"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewDrawCardsEffect(1)).
		AddEffect(abilities.NewDamageEffect(CardsInTargetHandCount.instance)).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	return card, nil
}