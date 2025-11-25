package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Omnath Locus Of All", NewOmnathLocusOfAll)
}

// NewOmnathLocusOfAll creates a Omnath Locus Of All
// {W}{U}{B/P}{R}{G} - CREATURE
func NewOmnathLocusOfAll(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Omnath Locus Of All")
	card.ManaCost = "{W}{U}{B/P}{R}{G}"
	card.Types = []string{"CREATURE"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "4"
	card.Toughness = "4"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
