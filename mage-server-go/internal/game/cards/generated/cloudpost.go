package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Cloudpost", NewCloudpost)
}

// NewCloudpost creates a Cloudpost
//   - LAND
func NewCloudpost(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Cloudpost")
	card.ManaCost = ""
	card.Types = []string{"LAND"}
	card.Subtypes = []string{"LOCUS"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
