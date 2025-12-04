package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Karn The Great Creator", NewKarnTheGreatCreator)
}

// NewKarnTheGreatCreator creates a Karn The Great Creator
// {4} - PLANESWALKER
func NewKarnTheGreatCreator(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Karn The Great Creator")
	card.ManaCost = "{4}"
	card.Types = []string{"PLANESWALKER"}
	card.Subtypes = []string{"KARN"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Loyalty = "5"
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement triggered ability: LoyaltyAbility
	//   - Effect: KarnTheGreatCreatorAnimateEffect()
	//
	// Targets:
	//   - abilities.NewTargetRequirement(1, 1, abilities.NewPermanentTargetFilter())
	// card.AddAbility(ability0)
	return card, nil
}
