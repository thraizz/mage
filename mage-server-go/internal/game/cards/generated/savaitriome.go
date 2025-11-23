package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Savai Triome", NewSavaiTriome)
}

// NewSavaiTriome creates a Savai Triome
//   - LAND
func NewSavaiTriome(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Savai Triome")
	card.ManaCost = ""
	card.Types = []string{"LAND"}
	card.Subtypes = []string{"MOUNTAIN", "PLAINS", "SWAMP"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.BuildSimpleManaAbility(card.ID, "R")
	card.AddAbility(ability0)
	ability1 := abilities.BuildSimpleManaAbility(card.ID, "W")
	card.AddAbility(ability1)
	ability2 := abilities.BuildSimpleManaAbility(card.ID, "B")
	card.AddAbility(ability2)
	return card, nil
}
