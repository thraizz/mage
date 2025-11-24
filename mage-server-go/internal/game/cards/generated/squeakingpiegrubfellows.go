package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Squeaking Pie Grubfellows", NewSqueakingPieGrubfellows)
}

// NewSqueakingPieGrubfellows creates a Squeaking Pie Grubfellows
// {3}{B} - CREATURE
func NewSqueakingPieGrubfellows(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Squeaking Pie Grubfellows")
	card.ManaCost = "{3}{B}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"GOBLIN", "SHAMAN"}
	card.Power = "3"
	card.Toughness = "2"
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - DiscardEachPlayerEffect(StaticValue.get(1), false, TargetController.OPPONE...)
	// card.AddAbility(ability0)
	return card, nil
}
