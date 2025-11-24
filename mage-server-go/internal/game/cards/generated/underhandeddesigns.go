package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Underhanded Designs", NewUnderhandedDesigns)
}

// NewUnderhandedDesigns creates a Underhanded Designs
// {1}{B} - ENCHANTMENT
func NewUnderhandedDesigns(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Underhanded Designs")
	card.ManaCost = "{1}{B}"
	card.Types = []string{"ENCHANTMENT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - DoIfCostPaid(new LoseLifeOpponentsEffect(1), new GenericManaCos...)
	// card.AddAbility(ability0)
	return card, nil
}
