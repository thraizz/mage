package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Evolutionary Leap", NewEvolutionaryLeap)
}

// NewEvolutionaryLeap creates a Evolutionary Leap
// {1}{G} - ENCHANTMENT
func NewEvolutionaryLeap(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Evolutionary Leap")
	card.ManaCost = "{1}{G}"
	card.Types = []string{"ENCHANTMENT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement activated ability with unmapped effects
	//   - RevealCardsFromLibraryUntilEffect()
	// card.AddAbility(ability0)
	return card, nil
}
