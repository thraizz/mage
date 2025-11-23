package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Elemental Uprising", NewElementalUprising)
}

// NewElementalUprising creates a Elemental Uprising
// {1}{G} - INSTANT
// Haste
func NewElementalUprising(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Elemental Uprising")
	card.ManaCost = "{1}{G}"
	card.Types = []string{"INSTANT"}
	card.Subtypes = []string{"ELEMENTAL"}
	card.Power = "4"
	card.Toughness = "4"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordHaste)
	card.AddAbility(ability0)
	return card, nil
}
