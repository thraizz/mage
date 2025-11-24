package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Thawing Glaciers", NewThawingGlaciers)
}

// NewThawingGlaciers creates a Thawing Glaciers
//   - LAND
func NewThawingGlaciers(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Thawing Glaciers")
	card.ManaCost = ""
	card.Types = []string{"LAND"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewActivatedAbilityBuilder(card.ID).
		AddManaCost("{1}").
		AddTapCost().
		// TODO: SearchLibraryPutInPlayEffect with complex parameters
		// TODO: ReturnToHandSourceEffect with complex parameters
		Build()
	card.AddAbility(ability0)
	return card, nil
}
