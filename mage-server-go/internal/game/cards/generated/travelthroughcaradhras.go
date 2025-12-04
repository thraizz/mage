package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Travel Through Caradhras", NewTravelThroughCaradhras)
}

// NewTravelThroughCaradhras creates a Travel Through Caradhras
// {5}{G} - SORCERY
func NewTravelThroughCaradhras(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Travel Through Caradhras")
	card.ManaCost = "{5}{G}"
	card.Types = []string{"SORCERY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
