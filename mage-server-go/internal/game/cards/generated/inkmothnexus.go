package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Inkmoth Nexus", NewInkmothNexus)
}

// NewInkmothNexus creates a Inkmoth Nexus
//  - LAND
// Flying
func NewInkmothNexus(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Inkmoth Nexus")
	card.ManaCost = ""
	card.Types = []string{"LAND"}
	card.Subtypes = []string{"PHYREXIAN", "BLINKMOTH"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordFlying)
	card.AddAbility(ability0)
	ability1 := abilities.BuildSimpleManaAbility(card.ID, "C")
	card.AddAbility(ability1)
	return card, nil
}