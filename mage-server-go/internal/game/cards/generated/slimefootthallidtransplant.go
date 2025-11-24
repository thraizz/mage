package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Slimefoot Thallid Transplant", NewSlimefootThallidTransplant)
}

// NewSlimefootThallidTransplant creates a Slimefoot Thallid Transplant
// {B}{G} - CREATURE
func NewSlimefootThallidTransplant(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Slimefoot Thallid Transplant")
	card.ManaCost = "{B}{G}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"FUNGUS"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "2"
	card.Toughness = "2"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
