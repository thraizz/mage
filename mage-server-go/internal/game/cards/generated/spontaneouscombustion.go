package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Spontaneous Combustion", NewSpontaneousCombustion)
}

// NewSpontaneousCombustion creates a Spontaneous Combustion
// {1}{B}{R} - INSTANT
func NewSpontaneousCombustion(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Spontaneous Combustion")
	card.ManaCost = "{1}{B}{R}"
	card.Types = []string{"INSTANT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - DamageAllEffect(3, new FilterCreaturePermanent())
	// card.AddAbility(ability0)
	return card, nil
}
