package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Capital City", NewCapitalCity)
}

// NewCapitalCity creates a Capital City
//   - LAND
func NewCapitalCity(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Capital City")
	card.ManaCost = ""
	card.Types = []string{"LAND"}
	card.Subtypes = []string{"TOWN"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.BuildSimpleManaAbility(card.ID, "C")
	card.AddAbility(ability0)
	return card, nil
}
