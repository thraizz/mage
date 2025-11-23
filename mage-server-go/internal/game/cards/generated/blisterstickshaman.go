package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Blisterstick Shaman", NewBlisterstickShaman)
}

// NewBlisterstickShaman creates a Blisterstick Shaman
// {2}{R} - CREATURE
func NewBlisterstickShaman(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Blisterstick Shaman")
	card.ManaCost = "{2}{R}"
	card.Types = []string{"CREATURE"}
	card.Power = "2"
	card.Toughness = "1"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
