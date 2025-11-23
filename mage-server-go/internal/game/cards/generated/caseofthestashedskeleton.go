package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Case Of The Stashed Skeleton", NewCaseOfTheStashedSkeleton)
}

// NewCaseOfTheStashedSkeleton creates a Case Of The Stashed Skeleton
// {1}{B} - ENCHANTMENT
func NewCaseOfTheStashedSkeleton(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Case Of The Stashed Skeleton")
	card.ManaCost = "{1}{B}"
	card.Types = []string{"ENCHANTMENT"}
	card.Subtypes = []string{"CASE"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
