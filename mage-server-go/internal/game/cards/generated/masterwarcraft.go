package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Master Warcraft", NewMasterWarcraft)
}

// NewMasterWarcraft creates a Master Warcraft
// {2}{R/W}{R/W} - INSTANT
func NewMasterWarcraft(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Master Warcraft")
	card.ManaCost = "{2}{R/W}{R/W}"
	card.Types = []string{"INSTANT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}