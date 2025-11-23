package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Needle Spires", NewNeedleSpires)
}

// NewNeedleSpires creates a Needle Spires
//   - LAND
//
// DoubleStrike
func NewNeedleSpires(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Needle Spires")
	card.ManaCost = ""
	card.Types = []string{"LAND"}
	card.Subtypes = []string{"ELEMENTAL"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordDoubleStrike)
	card.AddAbility(ability0)
	ability1 := abilities.BuildSimpleManaAbility(card.ID, "R")
	card.AddAbility(ability1)
	ability2 := abilities.BuildSimpleManaAbility(card.ID, "W")
	card.AddAbility(ability2)
	return card, nil
}
