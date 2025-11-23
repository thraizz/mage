package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Rock Basilisk", NewRockBasilisk)
}

// NewRockBasilisk creates a Rock Basilisk
// {4}{R}{G} - CREATURE
func NewRockBasilisk(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Rock Basilisk")
	card.ManaCost = "{4}{R}{G}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"BASILISK"}
	card.Power = "4"
	card.Toughness = "5"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
