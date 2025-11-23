package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("God Eternal Bontu", NewGodEternalBontu)
}

// NewGodEternalBontu creates a God Eternal Bontu
// {3}{B}{B} - CREATURE
func NewGodEternalBontu(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "God Eternal Bontu")
	card.ManaCost = "{3}{B}{B}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"ZOMBIE", "GOD"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "5"
	card.Toughness = "6"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
