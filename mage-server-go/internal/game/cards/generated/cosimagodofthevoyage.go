package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Cosima God Of The Voyage", NewCosimaGodOfTheVoyage)
}

// NewCosimaGodOfTheVoyage creates a Cosima God Of The Voyage
//   - CREATURE
func NewCosimaGodOfTheVoyage(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Cosima God Of The Voyage")
	card.ManaCost = ""
	card.Types = []string{"CREATURE"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
