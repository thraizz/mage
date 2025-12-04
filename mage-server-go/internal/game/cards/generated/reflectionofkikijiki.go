package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Reflection Of Kiki Jiki", NewReflectionOfKikiJiki)
}

// NewReflectionOfKikiJiki creates a Reflection Of Kiki Jiki
//   - ENCHANTMENT CREATURE
func NewReflectionOfKikiJiki(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Reflection Of Kiki Jiki")
	card.ManaCost = ""
	card.Types = []string{"ENCHANTMENT", "CREATURE"}
	card.Subtypes = []string{"GOBLIN", "SHAMAN"}
	card.Power = "2"
	card.Toughness = "2"
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - CreateTokenCopyTargetEffect(null, null, true)
	// card.AddAbility(ability0)
	// TODO: Implement activated ability with unmapped effects
	//   - ReflectionOfKikiJikiEffect()
	//
	// Costs:
	//   - AddManaCost("{1}")
	//   - AddTapCost()
	// card.AddAbility(ability1)
	return card, nil
}
