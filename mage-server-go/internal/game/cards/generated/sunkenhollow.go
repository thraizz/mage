package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Sunken Hollow", NewSunkenHollow)
}

// NewSunkenHollow creates a Sunken Hollow
//   - LAND
func NewSunkenHollow(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Sunken Hollow")
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
