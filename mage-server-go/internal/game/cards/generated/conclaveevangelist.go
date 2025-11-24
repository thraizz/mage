package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Conclave Evangelist", NewConclaveEvangelist)
}

// NewConclaveEvangelist creates a Conclave Evangelist
// {3}{G/W}{G/W} - CREATURE
func NewConclaveEvangelist(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Conclave Evangelist")
	card.ManaCost = "{3}{G/W}{G/W}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"ELEPHANT", "CLERIC"}
	card.Power = "4"
	card.Toughness = "4"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
