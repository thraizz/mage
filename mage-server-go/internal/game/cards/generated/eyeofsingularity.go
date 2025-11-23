package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Eye Of Singularity", NewEyeOfSingularity)
}

// NewEyeOfSingularity creates a Eye Of Singularity
// {3}{W} - ENCHANTMENT
func NewEyeOfSingularity(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Eye Of Singularity")
	card.ManaCost = "{3}{W}"
	card.Types = []string{"ENCHANTMENT"}
	card.Supertypes = []string{"WORLD"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
