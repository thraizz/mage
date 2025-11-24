package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Nyleas Huntmaster", NewNyleasHuntmaster)
}

// NewNyleasHuntmaster creates a Nyleas Huntmaster
// {3}{G} - CREATURE
func NewNyleasHuntmaster(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Nyleas Huntmaster")
	card.ManaCost = "{3}{G}"
	card.Types = []string{"CREATURE"}
	card.Power = "4"
	card.Toughness = "3"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}