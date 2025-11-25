package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Phyrexian War Beast", NewPhyrexianWarBeast)
}

// NewPhyrexianWarBeast creates a Phyrexian War Beast
// {3} - ARTIFACT CREATURE
func NewPhyrexianWarBeast(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Phyrexian War Beast")
	card.ManaCost = "{3}"
	card.Types = []string{"ARTIFACT", "CREATURE"}
	card.Subtypes = []string{"PHYREXIAN", "BEAST"}
	card.Power = "3"
	card.Toughness = "4"
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement triggered ability: LeavesBattlefieldTriggeredAbility
	//   - Effect: SacrificeControllerEffect(StaticFilters.FILTER_LAND, 1, "")
	// card.AddAbility(ability0)
	// TODO: Implement spell ability with unmapped effects
	//   - SacrificeControllerEffect(StaticFilters.FILTER_LAND, 1, "")
	// card.AddAbility(ability1)
	return card, nil
}
