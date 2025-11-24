package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Moratorium Stone", NewMoratoriumStone)
}

// NewMoratoriumStone creates a Moratorium Stone
// {1} - ARTIFACT
func NewMoratoriumStone(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Moratorium Stone")
	card.ManaCost = "{1}"
	card.Types = []string{"ARTIFACT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewActivatedAbilityBuilder(card.ID).
		AddManaCost("{2}").
		AddTapCost().
		// TODO: ExileTargetEffect with complex parameters
		Build()
	card.AddAbility(ability0)
	return card, nil
}
