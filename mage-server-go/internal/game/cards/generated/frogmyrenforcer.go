package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Frogmyr Enforcer", NewFrogmyrEnforcer)
}

// NewFrogmyrEnforcer creates a Frogmyr Enforcer
// {7} - ARTIFACT CREATURE
func NewFrogmyrEnforcer(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Frogmyr Enforcer")
	card.ManaCost = "{7}"
	card.Types = []string{"ARTIFACT", "CREATURE"}
	card.Subtypes = []string{"FROG", "MYR"}
	card.Power = "4"
	card.Toughness = "4"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
