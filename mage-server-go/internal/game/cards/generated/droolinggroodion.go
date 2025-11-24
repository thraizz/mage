package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Drooling Groodion", NewDroolingGroodion)
}

// NewDroolingGroodion creates a Drooling Groodion
// {3}{B}{B}{G} - CREATURE
func NewDroolingGroodion(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Drooling Groodion")
	card.ManaCost = "{3}{B}{B}{G}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"BEAST"}
	card.Power = "4"
	card.Toughness = "3"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewActivatedAbilityBuilder(card.ID).
		AddEffect(abilities.NewBoostEffect(2, 2)).
		AddEffect(abilities.NewBoostEffect(-2, -2)).
		Build()
	card.AddAbility(ability0)
	return card, nil
}
