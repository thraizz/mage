package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Loch Dragon", NewLochDragon)
}

// NewLochDragon creates a Loch Dragon
// {U/R}{U/R}{U/R}{U/R} - CREATURE
// Flying
func NewLochDragon(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Loch Dragon")
	card.ManaCost = "{U/R}{U/R}{U/R}{U/R}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"DRAGON"}
	card.Power = "3"
	card.Toughness = "2"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordFlying)
	card.AddAbility(ability0)
	// TODO: Implement spell ability with unmapped effects
	//   - DoIfCostPaid(new DrawCardSourceControllerEffect(1), new Discard...)
	// card.AddAbility(ability1)
	return card, nil
}
