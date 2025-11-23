package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Mammoth Spider", NewMammothSpider)
}

// NewMammothSpider creates a Mammoth Spider
// {4}{G} - CREATURE
// Reach
func NewMammothSpider(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Mammoth Spider")
	card.ManaCost = "{4}{G}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"SPIDER"}
	card.Power = "3"
	card.Toughness = "5"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordReach)
	card.AddAbility(ability0)
	return card, nil
}
