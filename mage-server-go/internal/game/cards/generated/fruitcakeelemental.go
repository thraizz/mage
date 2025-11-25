package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Fruitcake Elemental", NewFruitcakeElemental)
}

// NewFruitcakeElemental creates a Fruitcake Elemental
// {1}{G}{G} - CREATURE
// Indestructible
func NewFruitcakeElemental(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Fruitcake Elemental")
	card.ManaCost = "{1}{G}{G}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"ELEMENTAL"}
	card.Power = "7"
	card.Toughness = "7"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordIndestructible)
	card.AddAbility(ability0)
	// TODO: Implement activated ability with unmapped effects
	//   - FruitcakeElementalEffect()
	// card.AddAbility(ability1)
	return card, nil
}
