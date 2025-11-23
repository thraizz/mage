package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Wonderscape Sage", NewWonderscapeSage)
}

// NewWonderscapeSage creates a Wonderscape Sage
// {1}{U} - CREATURE
// Flying
func NewWonderscapeSage(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Wonderscape Sage")
	card.ManaCost = "{1}{U}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"MOONFOLK", "WIZARD"}
	card.Power = "1"
	card.Toughness = "3"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordFlying)
	card.AddAbility(ability0)
	return card, nil
}
