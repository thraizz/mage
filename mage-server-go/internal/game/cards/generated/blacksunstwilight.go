package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Black Suns Twilight", NewBlackSunsTwilight)
}

// NewBlackSunsTwilight creates a Black Suns Twilight
// {X}{B} - INSTANT
func NewBlackSunsTwilight(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Black Suns Twilight")
	card.ManaCost = "{X}{B}"
	card.Types = []string{"INSTANT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
