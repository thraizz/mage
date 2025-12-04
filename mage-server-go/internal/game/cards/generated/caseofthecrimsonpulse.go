package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Case Of The Crimson Pulse", NewCaseOfTheCrimsonPulse)
}

// NewCaseOfTheCrimsonPulse creates a Case Of The Crimson Pulse
// {2}{R} - ENCHANTMENT
func NewCaseOfTheCrimsonPulse(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Case Of The Crimson Pulse")
	card.ManaCost = "{2}{R}"
	card.Types = []string{"ENCHANTMENT"}
	card.Subtypes = []string{"CASE"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - DiscardControllerEffect(1)
	//   - DiscardHandControllerEffect()
	// card.AddAbility(ability0)
	return card, nil
}
