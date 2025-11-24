package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Giant Trap Door Spider", NewGiantTrapDoorSpider)
}

// NewGiantTrapDoorSpider creates a Giant Trap Door Spider
// {1}{R}{G} - CREATURE
func NewGiantTrapDoorSpider(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Giant Trap Door Spider")
	card.ManaCost = "{1}{R}{G}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"SPIDER"}
	card.Power = "2"
	card.Toughness = "3"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewActivatedAbilityBuilder(card.ID).
		AddTapCost().
		// TODO: ExileSourceEffect with complex parameters
		// TODO: ExileTargetEffect with complex parameters
		Build()
	card.AddAbility(ability0)
	return card, nil
}
