package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Kickoff Celebrations", NewKickoffCelebrations)
}

// NewKickoffCelebrations creates a Kickoff Celebrations
// {1}{R} - ENCHANTMENT
func NewKickoffCelebrations(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Kickoff Celebrations")
	card.ManaCost = "{1}{R}"
	card.Types = []string{"ENCHANTMENT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - DoIfCostPaid(new DrawCardSourceControllerEffect(2), new Discard...)
	// card.AddAbility(ability0)
	return card, nil
}
