package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Prismari Pledgemage", NewPrismariPledgemage)
}

// NewPrismariPledgemage creates a Prismari Pledgemage
// {U/R}{U/R} - CREATURE
// Defender
func NewPrismariPledgemage(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Prismari Pledgemage")
	card.ManaCost = "{U/R}{U/R}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"ORC", "WIZARD"}
	card.Power = "3"
	card.Toughness = "3"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordDefender)
	card.AddAbility(ability0)
	return card, nil
}
