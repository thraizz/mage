package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Talisman Of Conviction", NewTalismanOfConviction)
}

// NewTalismanOfConviction creates a Talisman Of Conviction
// {2} - ARTIFACT
func NewTalismanOfConviction(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Talisman Of Conviction")
	card.ManaCost = "{2}"
	card.Types = []string{"ARTIFACT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.BuildSimpleManaAbility(card.ID, "C")
	card.AddAbility(ability0)
	ability1 := abilities.BuildSimpleManaAbility(card.ID, "R")
	card.AddAbility(ability1)
	ability2 := abilities.BuildSimpleManaAbility(card.ID, "W")
	card.AddAbility(ability2)
	return card, nil
}
