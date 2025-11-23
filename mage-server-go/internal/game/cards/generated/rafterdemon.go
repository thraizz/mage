package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Rafter Demon", NewRafterDemon)
}

// NewRafterDemon creates a Rafter Demon
// {2}{B}{R} - CREATURE
func NewRafterDemon(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Rafter Demon")
	card.ManaCost = "{2}{B}{R}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"DEMON"}
	card.Power = "4"
	card.Toughness = "2"
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - DiscardEachPlayerEffect(                 StaticValue.get(1), false, Target...)
	// card.AddAbility(ability0)
	return card, nil
}