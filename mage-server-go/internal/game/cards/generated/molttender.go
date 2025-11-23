package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Molt Tender", NewMoltTender)
}

// NewMoltTender creates a Molt Tender
// {G} - CREATURE
func NewMoltTender(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Molt Tender")
	card.ManaCost = "{G}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"INSECT", "DRUID"}
	card.Power = "1"
	card.Toughness = "1"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
