package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Honored Heirloom", NewHonoredHeirloom)
}

// NewHonoredHeirloom creates a Honored Heirloom
// {3} - ARTIFACT
func NewHonoredHeirloom(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Honored Heirloom")
	card.ManaCost = "{3}"
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
