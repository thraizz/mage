package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Smoldering Marsh", NewSmolderingMarsh)
}

// NewSmolderingMarsh creates a Smoldering Marsh
//  - LAND
func NewSmolderingMarsh(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Smoldering Marsh")
	card.ManaCost = ""
	card.Types = []string{"LAND"}
	card.Subtypes = []string{"SWAMP", "MOUNTAIN"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.BuildSimpleManaAbility(card.ID, "B")
	card.AddAbility(ability0)
	ability1 := abilities.BuildSimpleManaAbility(card.ID, "R")
	card.AddAbility(ability1)
	return card, nil
}