package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Witchstalker", NewWitchstalker)
}

// NewWitchstalker creates a Witchstalker
// {1}{G}{G} - CREATURE
// Hexproof
func NewWitchstalker(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Witchstalker")
	card.ManaCost = "{1}{G}{G}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"WOLF"}
	card.Power = "3"
	card.Toughness = "3"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordHexproof)
	card.AddAbility(ability0)
	return card, nil
}
