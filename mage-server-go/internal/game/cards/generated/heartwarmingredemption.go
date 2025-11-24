package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Heartwarming Redemption", NewHeartwarmingRedemption)
}

// NewHeartwarmingRedemption creates a Heartwarming Redemption
// {2}{R}{W} - INSTANT
func NewHeartwarmingRedemption(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Heartwarming Redemption")
	card.ManaCost = "{2}{R}{W}"
	card.Types = []string{"INSTANT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}