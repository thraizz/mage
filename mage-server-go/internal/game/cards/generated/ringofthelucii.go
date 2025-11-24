package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Ring Of The Lucii", NewRingOfTheLucii)
}

// NewRingOfTheLucii creates a Ring Of The Lucii
// {4} - ARTIFACT
func NewRingOfTheLucii(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Ring Of The Lucii")
	card.ManaCost = "{4}"
	card.Types = []string{"ARTIFACT"}
	card.Supertypes = []string{"LEGENDARY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewActivatedAbilityBuilder(card.ID).
		AddManaCost("{2}").
		AddTapCost().
		// TODO: TapTargetEffect with complex parameters
		Build()
	card.AddAbility(ability0)
	return card, nil
}
