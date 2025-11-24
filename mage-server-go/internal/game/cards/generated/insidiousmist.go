package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Insidious Mist", NewInsidiousMist)
}

// NewInsidiousMist creates a Insidious Mist
//  - CREATURE
// Hexproof, Indestructible
func NewInsidiousMist(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Insidious Mist")
	card.ManaCost = ""
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"ELEMENTAL"}
	card.Power = "0"
	card.Toughness = "1"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordHexproof)
	card.AddAbility(ability0)
	ability1 := abilities.NewKeywordAbility(card.ID, abilities.KeywordIndestructible)
	card.AddAbility(ability1)
	// TODO: Implement spell ability with unmapped effects
	//   - TransformSourceEffect()
	//   - DoIfCostPaid(new TransformSourceEffect()                 .setTe...)
	// card.AddAbility(ability2)
	return card, nil
}