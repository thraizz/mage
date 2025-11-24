package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Sharlayan Nation Of Scholars", NewSharlayanNationOfScholars)
}

// NewSharlayanNationOfScholars creates a Sharlayan Nation Of Scholars
//  - LAND
func NewSharlayanNationOfScholars(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Sharlayan Nation Of Scholars")
	card.ManaCost = ""
	card.Types = []string{"LAND"}
	card.Subtypes = []string{"TOWN"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.BuildSimpleManaAbility(card.ID, "W")
	card.AddAbility(ability0)
	ability1 := abilities.BuildSimpleManaAbility(card.ID, "U")
	card.AddAbility(ability1)
	return card, nil
}