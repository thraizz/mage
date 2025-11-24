package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Lurking Evil", NewLurkingEvil)
}

// NewLurkingEvil creates a Lurking Evil
// {B}{B}{B} - ENCHANTMENT
// Flying
func NewLurkingEvil(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Lurking Evil")
	card.ManaCost = "{B}{B}{B}"
	card.Types = []string{"ENCHANTMENT"}
	card.Subtypes = []string{"PHYREXIAN", "HORROR"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordFlying)
	card.AddAbility(ability0)
	return card, nil
}
