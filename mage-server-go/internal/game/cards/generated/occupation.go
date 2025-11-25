package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Occupation", NewOccupation)
}

// NewOccupation creates a Occupation
// {W}{B} - ENCHANTMENT
func NewOccupation(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Occupation")
	card.ManaCost = "{W}{B}"
	card.Types = []string{"ENCHANTMENT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement activated ability with unmapped effects
	//   - OccupationOneShotEffect()
	// card.AddAbility(ability0)
	return card, nil
}
