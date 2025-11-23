package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Prahv Spires Of Order", NewPrahvSpiresOfOrder)
}

// NewPrahvSpiresOfOrder creates a Prahv Spires Of Order
//   - LAND
func NewPrahvSpiresOfOrder(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Prahv Spires Of Order")
	card.ManaCost = ""
	card.Types = []string{"LAND"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.BuildSimpleManaAbility(card.ID, "C")
	card.AddAbility(ability0)
	return card, nil
}
