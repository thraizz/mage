package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Mudbrawler Raiders", NewMudbrawlerRaiders)
}

// NewMudbrawlerRaiders creates a Mudbrawler Raiders
// {2}{R/G}{R/G} - CREATURE
func NewMudbrawlerRaiders(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Mudbrawler Raiders")
	card.ManaCost = "{2}{R/G}{R/G}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"GOBLIN", "WARRIOR"}
	card.Power = "3"
	card.Toughness = "3"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
