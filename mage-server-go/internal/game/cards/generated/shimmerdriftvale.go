package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Shimmerdrift Vale", NewShimmerdriftVale)
}

// NewShimmerdriftVale creates a Shimmerdrift Vale
//  - LAND
func NewShimmerdriftVale(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Shimmerdrift Vale")
	card.ManaCost = ""
	card.Types = []string{"LAND"}
	card.Supertypes = []string{"SNOW"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}