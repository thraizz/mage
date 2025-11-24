package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Krosan Cloudscraper", NewKrosanCloudscraper)
}

// NewKrosanCloudscraper creates a Krosan Cloudscraper
// {7}{G}{G}{G} - CREATURE
func NewKrosanCloudscraper(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Krosan Cloudscraper")
	card.ManaCost = "{7}{G}{G}{G}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"BEAST", "MUTANT"}
	card.Power = "13"
	card.Toughness = "13"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
