package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Crimson Honor Guard", NewCrimsonHonorGuard)
}

// NewCrimsonHonorGuard creates a Crimson Honor Guard
// {3}{R}{R} - CREATURE
// Trample
func NewCrimsonHonorGuard(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Crimson Honor Guard")
	card.ManaCost = "{3}{R}{R}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"VAMPIRE", "KNIGHT"}
	card.Power = "4"
	card.Toughness = "5"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordTrample)
	card.AddAbility(ability0)
	return card, nil
}
