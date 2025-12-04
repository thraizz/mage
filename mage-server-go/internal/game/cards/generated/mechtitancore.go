package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Mechtitan Core", NewMechtitanCore)
}

// NewMechtitanCore creates a Mechtitan Core
// {2} - ARTIFACT
func NewMechtitanCore(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Mechtitan Core")
	card.ManaCost = "{2}"
	card.Types = []string{"ARTIFACT"}
	card.Subtypes = []string{"VEHICLE"}
	card.Power = "2"
	card.Toughness = "4"
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement activated ability with unmapped effects
	//   - MechtitanCoreTokenEffect()
	//
	// Costs:
	//   - AddManaCost("{5}")
	// card.AddAbility(ability0)
	return card, nil
}
