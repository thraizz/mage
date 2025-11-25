package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Kaya Ghost Assassin", NewKayaGhostAssassin)
}

// NewKayaGhostAssassin creates a Kaya Ghost Assassin
// {2}{W}{B} - PLANESWALKER
func NewKayaGhostAssassin(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Kaya Ghost Assassin")
	card.ManaCost = "{2}{W}{B}"
	card.Types = []string{"PLANESWALKER"}
	card.Subtypes = []string{"KAYA"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Loyalty = "5"
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement triggered ability: LoyaltyAbility
	//   - Effect: KayaGhostAssassinEffect()
	//
	// Targets:
	//   - abilities.NewTargetRequirement(1, 1, abilities.NewCreatureTargetFilter())
	// card.AddAbility(ability0)
	// TODO: Implement spell ability with unmapped effects
	//   - DiscardEachPlayerEffect(TargetController.OPPONENT)
	// card.AddAbility(ability1)
	return card, nil
}
