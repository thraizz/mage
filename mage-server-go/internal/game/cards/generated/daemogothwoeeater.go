package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Daemogoth Woe Eater", NewDaemogothWoeEater)
}

// NewDaemogothWoeEater creates a Daemogoth Woe Eater
// {1}{B}{B/G}{G} - CREATURE
func NewDaemogothWoeEater(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Daemogoth Woe Eater")
	card.ManaCost = "{1}{B}{B/G}{G}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"DEMON"}
	card.Power = "7"
	card.Toughness = "6"
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement triggered ability: SacrificeSourceTriggeredAbility
	//   - Effect: DiscardEachPlayerEffect(TargetController.OPPONENT)
	// card.AddAbility(ability0)
	// TODO: Implement spell ability with unmapped effects
	//   - SacrificeControllerEffect(                 StaticFilters.FILTER_PERMANENT_CR...)
	// card.AddAbility(ability1)
	return card, nil
}
