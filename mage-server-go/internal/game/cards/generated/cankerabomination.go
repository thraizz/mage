package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Canker Abomination", NewCankerAbomination)
}

// NewCankerAbomination creates a Canker Abomination
// {2}{B/G}{B/G} - CREATURE
func NewCankerAbomination(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Canker Abomination")
	card.ManaCost = "{2}{B/G}{B/G}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"TREEFOLK", "HORROR"}
	card.Power = "6"
	card.Toughness = "6"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}