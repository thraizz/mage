package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Benalish Trapper", NewBenalishTrapper)
}

// NewBenalishTrapper creates a Benalish Trapper
// {1}{W} - CREATURE
func NewBenalishTrapper(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Benalish Trapper")
	card.ManaCost = "{1}{W}"
	card.Types = []string{"CREATURE"}
	card.Power = "1"
	card.Toughness = "2"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewActivatedAbilityBuilder(card.ID).
		AddTapCost().
		// TODO: TapTargetEffect with complex parameters
		Build()
	card.AddAbility(ability0)
	return card, nil
}
