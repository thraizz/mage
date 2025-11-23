package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Concealing Curtains", NewConcealingCurtains)
}

// NewConcealingCurtains creates a Concealing Curtains
// {B} - CREATURE
// Defender
func NewConcealingCurtains(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Concealing Curtains")
	card.ManaCost = "{B}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"WALL"}
	card.Power = "0"
	card.Toughness = "4"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordDefender)
	card.AddAbility(ability0)
	// TODO: Implement spell ability with unmapped effects
	//   - TransformSourceEffect()
	// card.AddAbility(ability1)
	return card, nil
}
