package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Soul Conduit", NewSoulConduit)
}

// NewSoulConduit creates a Soul Conduit
// {6} - ARTIFACT
func NewSoulConduit(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Soul Conduit")
	card.ManaCost = "{6}"
	card.Types = []string{"ARTIFACT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement activated ability with unmapped effects
	//   - ExchangeLifeTwoTargetEffect()
	//
	// Costs:
	//   - AddManaCost("{6}")
	//   - AddTapCost()
	// card.AddAbility(ability0)
	return card, nil
}
