package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Evil Eye Of Orms By Gore", NewEvilEyeOfOrmsByGore)
}

// NewEvilEyeOfOrmsByGore creates a Evil Eye Of Orms By Gore
// {4}{B} - CREATURE
func NewEvilEyeOfOrmsByGore(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Evil Eye Of Orms By Gore")
	card.ManaCost = "{4}{B}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"EYE"}
	card.Power = "3"
	card.Toughness = "6"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}