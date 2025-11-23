package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Desert Of The Mindful", NewDesertOfTheMindful)
}

// NewDesertOfTheMindful creates a Desert Of The Mindful
//   - LAND
func NewDesertOfTheMindful(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Desert Of The Mindful")
	card.ManaCost = ""
	card.Types = []string{"LAND"}
	card.Subtypes = []string{"DESERT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.BuildSimpleManaAbility(card.ID, "U")
	card.AddAbility(ability0)
	return card, nil
}
