package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Defend The Celestus", NewDefendTheCelestus)
}

// NewDefendTheCelestus creates a Defend The Celestus
// {2}{G}{G} - INSTANT
func NewDefendTheCelestus(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Defend The Celestus")
	card.ManaCost = "{2}{G}{G}"
	card.Types = []string{"INSTANT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
