package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Goldbug Scrappy Scout", NewGoldbugScrappyScout)
}

// NewGoldbugScrappyScout creates a Goldbug Scrappy Scout
//   - ARTIFACT
func NewGoldbugScrappyScout(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Goldbug Scrappy Scout")
	card.ManaCost = ""
	card.Types = []string{"ARTIFACT"}
	card.Subtypes = []string{"VEHICLE"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "1"
	card.Toughness = "3"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
