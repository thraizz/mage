package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Frostwalk Bastion", NewFrostwalkBastion)
}

// NewFrostwalkBastion creates a Frostwalk Bastion
//  - LAND
func NewFrostwalkBastion(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Frostwalk Bastion")
	card.ManaCost = ""
	card.Types = []string{"LAND"}
	card.Subtypes = []string{"CONSTRUCT"}
	card.Supertypes = []string{"SNOW"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.BuildSimpleManaAbility(card.ID, "C")
	card.AddAbility(ability0)
	return card, nil
}