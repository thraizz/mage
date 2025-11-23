package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Dualcaster Mage", NewDualcasterMage)
}

// NewDualcasterMage creates a Dualcaster Mage
// {1}{R}{R} - CREATURE
// Flash
func NewDualcasterMage(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Dualcaster Mage")
	card.ManaCost = "{1}{R}{R}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"HUMAN", "WIZARD"}
	card.Power = "2"
	card.Toughness = "2"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordFlash)
	card.AddAbility(ability0)
	return card, nil
}
