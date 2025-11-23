package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Basalt Golem", NewBasaltGolem)
}

// NewBasaltGolem creates a Basalt Golem
// {5} - ARTIFACT CREATURE
func NewBasaltGolem(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Basalt Golem")
	card.ManaCost = "{5}"
	card.Types = []string{"ARTIFACT", "CREATURE"}
	card.Subtypes = []string{"GOLEM"}
	card.Power = "2"
	card.Toughness = "4"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
