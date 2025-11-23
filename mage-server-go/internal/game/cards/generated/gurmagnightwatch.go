package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Gurmag Nightwatch", NewGurmagNightwatch)
}

// NewGurmagNightwatch creates a Gurmag Nightwatch
// {2/B}{2/G}{2/U} - CREATURE
func NewGurmagNightwatch(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Gurmag Nightwatch")
	card.ManaCost = "{2/B}{2/G}{2/U}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"HUMAN", "RANGER"}
	card.Power = "3"
	card.Toughness = "3"
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - LookLibraryAndPickControllerEffect(                 3, 1, PutCards.TOP_ANY, PutCards....)
	// card.AddAbility(ability0)
	return card, nil
}
