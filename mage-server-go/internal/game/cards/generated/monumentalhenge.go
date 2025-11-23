package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Monumental Henge", NewMonumentalHenge)
}

// NewMonumentalHenge creates a Monumental Henge
//   - LAND
func NewMonumentalHenge(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Monumental Henge")
	card.ManaCost = ""
	card.Types = []string{"LAND"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.BuildSimpleManaAbility(card.ID, "W")
	card.AddAbility(ability0)
	// TODO: Implement activated ability with unmapped effects
	//   - LookLibraryAndPickControllerEffect(5, 1, filterCard, PutCards.HAND, PutCards.BOTTOM_R...)
	//
	// Costs:
	//   - AddTapCost()
	// card.AddAbility(ability1)
	return card, nil
}
