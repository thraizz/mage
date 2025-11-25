package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Navigation Orb", NewNavigationOrb)
}

// NewNavigationOrb creates a Navigation Orb
// {3} - ARTIFACT
func NewNavigationOrb(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Navigation Orb")
	card.ManaCost = "{3}"
	card.Types = []string{"ARTIFACT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement activated ability with unmapped effects
	//   - SearchLibraryPutOntoBattlefieldTappedRestInHandEffect()
	//
	// Costs:
	//   - AddManaCost("{2}")
	//   - AddTapCost()
	//   - AddSacrificeSourceCost()
	// card.AddAbility(ability0)
	return card, nil
}
