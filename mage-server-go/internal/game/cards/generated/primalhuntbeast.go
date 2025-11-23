package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Primal Huntbeast", NewPrimalHuntbeast)
}

// NewPrimalHuntbeast creates a Primal Huntbeast
// {3}{G} - CREATURE
// Hexproof
func NewPrimalHuntbeast(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Primal Huntbeast")
	card.ManaCost = "{3}{G}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"BEAST"}
	card.Power = "3"
	card.Toughness = "3"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordHexproof)
	card.AddAbility(ability0)
	return card, nil
}
