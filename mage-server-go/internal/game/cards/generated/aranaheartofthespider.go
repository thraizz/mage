package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Arana Heart Of The Spider", NewAranaHeartOfTheSpider)
}

// NewAranaHeartOfTheSpider creates a Arana Heart Of The Spider
// {1}{R}{W} - CREATURE
func NewAranaHeartOfTheSpider(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Arana Heart Of The Spider")
	card.ManaCost = "{1}{R}{W}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"SPIDER", "HUMAN", "HERO"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "3"
	card.Toughness = "3"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
