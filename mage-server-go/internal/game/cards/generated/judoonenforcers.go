package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Judoon Enforcers", NewJudoonEnforcers)
}

// NewJudoonEnforcers creates a Judoon Enforcers
// {5}{R}{W} - CREATURE
// Trample
func NewJudoonEnforcers(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Judoon Enforcers")
	card.ManaCost = "{5}{R}{W}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"ALIEN", "RHINO", "SOLDIER"}
	card.Power = "8"
	card.Toughness = "8"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordTrample)
	card.AddAbility(ability0)
	return card, nil
}