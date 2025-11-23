package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Merfolk Of The Pearl Trident", NewMerfolkOfThePearlTrident)
}

// NewMerfolkOfThePearlTrident creates a Merfolk Of The Pearl Trident
// {U} - CREATURE
func NewMerfolkOfThePearlTrident(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Merfolk Of The Pearl Trident")
	card.ManaCost = "{U}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"MERFOLK"}
	card.Power = "1"
	card.Toughness = "1"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
