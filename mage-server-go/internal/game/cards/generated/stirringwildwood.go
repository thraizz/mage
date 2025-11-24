package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Stirring Wildwood", NewStirringWildwood)
}

// NewStirringWildwood creates a Stirring Wildwood
//   - LAND
//
// Reach
func NewStirringWildwood(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Stirring Wildwood")
	card.ManaCost = ""
	card.Types = []string{"LAND"}
	card.Subtypes = []string{"ELEMENTAL"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordReach)
	card.AddAbility(ability0)
	ability1 := abilities.BuildSimpleManaAbility(card.ID, "G")
	card.AddAbility(ability1)
	ability2 := abilities.BuildSimpleManaAbility(card.ID, "W")
	card.AddAbility(ability2)
	return card, nil
}
