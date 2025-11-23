package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Cascading Cataracts", NewCascadingCataracts)
}

// NewCascadingCataracts creates a Cascading Cataracts
//   - LAND
//
// Indestructible
func NewCascadingCataracts(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Cascading Cataracts")
	card.ManaCost = ""
	card.Types = []string{"LAND"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordIndestructible)
	card.AddAbility(ability0)
	ability1 := abilities.BuildSimpleManaAbility(card.ID, "C")
	card.AddAbility(ability1)
	return card, nil
}
