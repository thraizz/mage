package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Yoda Jedi Master", NewYodaJediMaster)
}

// NewYodaJediMaster creates a Yoda Jedi Master
// {1}{G}{U} - PLANESWALKER
func NewYodaJediMaster(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Yoda Jedi Master")
	card.ManaCost = "{1}{G}{U}"
	card.Types = []string{"PLANESWALKER"}
	card.Subtypes = []string{"YODA"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Loyalty = "3"
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement triggered ability: LoyaltyAbility
	//   - Effect: YodaJediMasterEffect()
	//
	// Targets:
	//   - abilities.NewTargetRequirement(1, 1, abilities.NewPermanentTargetFilter())
	// card.AddAbility(ability0)
	// TODO: Implement spell ability with unmapped effects
	//   - LookLibraryAndPickControllerEffect(2, 1, PutCards.BOTTOM_ANY, PutCards.TOP_ANY)
	// card.AddAbility(ability1)
	return card, nil
}
