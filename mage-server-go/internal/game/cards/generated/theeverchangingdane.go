package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("The Ever Changing Dane", NewTheEverChangingDane)
}

// NewTheEverChangingDane creates a The Ever Changing Dane
// {W}{U}{B} - CREATURE
func NewTheEverChangingDane(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "The Ever Changing Dane")
	card.ManaCost = "{W}{U}{B}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"SHAPESHIFTER"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "3"
	card.Toughness = "3"
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement activated ability with unmapped effects
	//   - TheEverChangingDaneEffect()
	//
	// Costs:
	//   - AddManaCost("{1}")
	// card.AddAbility(ability0)
	return card, nil
}
