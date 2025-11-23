package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Kotose The Silent Spider", NewKotoseTheSilentSpider)
}

// NewKotoseTheSilentSpider creates a Kotose The Silent Spider
// {3}{U}{B} - CREATURE
func NewKotoseTheSilentSpider(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Kotose The Silent Spider")
	card.ManaCost = "{3}{U}{B}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"HUMAN", "NINJA"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "4"
	card.Toughness = "4"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
