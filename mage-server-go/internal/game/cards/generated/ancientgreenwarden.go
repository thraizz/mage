package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Ancient Greenwarden", NewAncientGreenwarden)
}

// NewAncientGreenwarden creates a Ancient Greenwarden
// {4}{G}{G} - CREATURE
// Reach
func NewAncientGreenwarden(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Ancient Greenwarden")
	card.ManaCost = "{4}{G}{G}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"ELEMENTAL"}
	card.Power = "5"
	card.Toughness = "7"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordReach)
	card.AddAbility(ability0)
	return card, nil
}
