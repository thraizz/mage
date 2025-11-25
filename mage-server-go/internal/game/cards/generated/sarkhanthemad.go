package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Sarkhan The Mad", NewSarkhanTheMad)
}

// NewSarkhanTheMad creates a Sarkhan The Mad
// {3}{B}{R} - PLANESWALKER
func NewSarkhanTheMad(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Sarkhan The Mad")
	card.ManaCost = "{3}{B}{R}"
	card.Types = []string{"PLANESWALKER"}
	card.Subtypes = []string{"SARKHAN"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Loyalty = "7"
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement triggered ability: LoyaltyAbility
	//   - Effect: SarkhanTheMadSacEffect()
	// card.AddAbility(ability0)
	// TODO: Implement triggered ability: LoyaltyAbility
	//   - Effect: SarkhanTheMadDragonDamageEffect()
	//
	// Targets:
	//   - abilities.NewTargetRequirement(1, 1, abilities.NewPlayerTargetFilter())
	// card.AddAbility(ability1)
	return card, nil
}
