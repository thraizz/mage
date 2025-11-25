package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Polygraph Orb", NewPolygraphOrb)
}

// NewPolygraphOrb creates a Polygraph Orb
// {4}{B} - ARTIFACT
func NewPolygraphOrb(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Polygraph Orb")
	card.ManaCost = "{4}{B}"
	card.Types = []string{"ARTIFACT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement triggered ability: EntersBattlefieldTriggeredAbility
	//   - Effect: LookLibraryAndPickControllerEffect(                         4, 2, PutCards.HAND, PutC...)
	// card.AddAbility(ability0)
	// TODO: Implement activated ability with unmapped effects
	//   - PolygraphOrbEffect()
	//
	// Costs:
	//   - AddManaCost("{2}")
	//   - AddTapCost()
	// card.AddAbility(ability1)
	return card, nil
}
