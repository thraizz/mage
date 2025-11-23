package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Needleshot Gourna", NewNeedleshotGourna)
}

// NewNeedleshotGourna creates a Needleshot Gourna
// {4}{G}{G} - CREATURE
// Reach
func NewNeedleshotGourna(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Needleshot Gourna")
	card.ManaCost = "{4}{G}{G}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"BEAST"}
	card.Power = "3"
	card.Toughness = "6"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordReach)
	card.AddAbility(ability0)
	return card, nil
}
