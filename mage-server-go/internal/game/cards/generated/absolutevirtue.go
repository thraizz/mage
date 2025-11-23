package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Absolute Virtue", NewAbsoluteVirtue)
}

// NewAbsoluteVirtue creates a Absolute Virtue
// {6}{W}{U} - CREATURE
// Flying
func NewAbsoluteVirtue(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Absolute Virtue")
	card.ManaCost = "{6}{W}{U}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"AVATAR", "WARRIOR"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "8"
	card.Toughness = "8"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordFlying)
	card.AddAbility(ability0)
	return card, nil
}
