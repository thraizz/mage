package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Contaminated Aquifer", NewContaminatedAquifer)
}

// NewContaminatedAquifer creates a Contaminated Aquifer
//   - LAND
func NewContaminatedAquifer(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Contaminated Aquifer")
	card.ManaCost = ""
	card.Types = []string{"LAND"}
	card.Subtypes = []string{"ISLAND", "SWAMP"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.BuildSimpleManaAbility(card.ID, "U")
	card.AddAbility(ability0)
	ability1 := abilities.BuildSimpleManaAbility(card.ID, "B")
	card.AddAbility(ability1)
	return card, nil
}
