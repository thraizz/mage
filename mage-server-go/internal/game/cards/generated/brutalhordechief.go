package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Brutal Hordechief", NewBrutalHordechief)
}

// NewBrutalHordechief creates a Brutal Hordechief
// {3}{B} - CREATURE
func NewBrutalHordechief(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Brutal Hordechief")
	card.ManaCost = "{3}{B}"
	card.Types = []string{"CREATURE"}
	card.Power = "3"
	card.Toughness = "3"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
