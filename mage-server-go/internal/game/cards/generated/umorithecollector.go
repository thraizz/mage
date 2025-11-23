package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Umori The Collector", NewUmoriTheCollector)
}

// NewUmoriTheCollector creates a Umori The Collector
// {2}{B/G}{B/G} - CREATURE
func NewUmoriTheCollector(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Umori The Collector")
	card.ManaCost = "{2}{B/G}{B/G}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"OOZE"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "4"
	card.Toughness = "5"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
