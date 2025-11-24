package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Youve Been Caught Stealing", NewYouveBeenCaughtStealing)
}

// NewYouveBeenCaughtStealing creates a Youve Been Caught Stealing
// {1}{R} - SORCERY
func NewYouveBeenCaughtStealing(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Youve Been Caught Stealing")
	card.ManaCost = "{1}{R}"
	card.Types = []string{"SORCERY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
