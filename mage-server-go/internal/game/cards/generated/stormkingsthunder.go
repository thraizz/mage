package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Storm Kings Thunder", NewStormKingsThunder)
}

// NewStormKingsThunder creates a Storm Kings Thunder
// {X}{R}{R}{R} - INSTANT
func NewStormKingsThunder(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Storm Kings Thunder")
	card.ManaCost = "{X}{R}{R}{R}"
	card.Types = []string{"INSTANT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
