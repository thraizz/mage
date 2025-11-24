package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Olorins Searing Light", NewOlorinsSearingLight)
}

// NewOlorinsSearingLight creates a Olorins Searing Light
// {2}{R}{W} - INSTANT
func NewOlorinsSearingLight(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Olorins Searing Light")
	card.ManaCost = "{2}{R}{W}"
	card.Types = []string{"INSTANT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}