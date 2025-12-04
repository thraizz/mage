package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Yawgmoths Vile Offering", NewYawgmothsVileOffering)
}

// NewYawgmothsVileOffering creates a Yawgmoths Vile Offering
// {4}{B} - SORCERY
func NewYawgmothsVileOffering(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Yawgmoths Vile Offering")
	card.ManaCost = "{4}{B}"
	card.Types = []string{"SORCERY"}
	card.Supertypes = []string{"LEGENDARY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
