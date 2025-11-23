package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Reap The Past", NewReapThePast)
}

// NewReapThePast creates a Reap The Past
// {X}{R}{G} - SORCERY
func NewReapThePast(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Reap The Past")
	card.ManaCost = "{X}{R}{G}"
	card.Types = []string{"SORCERY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
