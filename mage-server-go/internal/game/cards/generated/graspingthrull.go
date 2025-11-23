package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Grasping Thrull", NewGraspingThrull)
}

// NewGraspingThrull creates a Grasping Thrull
// {3}{W}{B} - CREATURE
// Flying
func NewGraspingThrull(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Grasping Thrull")
	card.ManaCost = "{3}{W}{B}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"THRULL"}
	card.Power = "3"
	card.Toughness = "3"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordFlying)
	card.AddAbility(ability0)
	return card, nil
}
