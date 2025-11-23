package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Sacred Foundry", NewSacredFoundry)
}

// NewSacredFoundry creates a Sacred Foundry
//   - LAND
func NewSacredFoundry(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Sacred Foundry")
	card.ManaCost = ""
	card.Types = []string{"LAND"}
	card.Subtypes = []string{"MOUNTAIN", "PLAINS"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.BuildSimpleManaAbility(card.ID, "R")
	card.AddAbility(ability0)
	ability1 := abilities.BuildSimpleManaAbility(card.ID, "W")
	card.AddAbility(ability1)
	return card, nil
}
