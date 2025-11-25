package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Stone Of Erech", NewStoneOfErech)
}

// NewStoneOfErech creates a Stone Of Erech
// {1} - ARTIFACT
func NewStoneOfErech(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Stone Of Erech")
	card.ManaCost = "{1}"
	card.Types = []string{"ARTIFACT"}
	card.Supertypes = []string{"LEGENDARY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement activated ability with unmapped effects
	//   - ExileGraveyardAllTargetPlayerEffect()
	//
	// Costs:
	//   - AddTapCost()
	//   - AddManaCost("{2}")
	//   - AddSacrificeSourceCost()
	// card.AddAbility(ability0)
	return card, nil
}
