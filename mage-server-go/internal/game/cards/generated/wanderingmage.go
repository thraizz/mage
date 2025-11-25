package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Wandering Mage", NewWanderingMage)
}

// NewWanderingMage creates a Wandering Mage
// {W}{U}{B} - CREATURE
func NewWanderingMage(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Wandering Mage")
	card.ManaCost = "{W}{U}{B}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"HUMAN", "CLERIC", "WIZARD"}
	card.Power = "0"
	card.Toughness = "3"
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement activated ability with unmapped effects
	//   - PreventDamageToTargetEffect()
	// card.AddAbility(ability0)
	// TODO: Implement activated ability with unmapped effects
	//   - PreventDamageToTargetEffect()
	// card.AddAbility(ability1)
	// TODO: Implement activated ability with unmapped effects
	//   - PreventDamageToTargetEffect()
	// card.AddAbility(ability2)
	return card, nil
}
