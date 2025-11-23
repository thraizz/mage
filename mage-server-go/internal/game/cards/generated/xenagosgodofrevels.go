package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Xenagos God Of Revels", NewXenagosGodOfRevels)
}

// NewXenagosGodOfRevels creates a Xenagos God Of Revels
// {3}{R}{G} - ENCHANTMENT CREATURE
// Indestructible
func NewXenagosGodOfRevels(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Xenagos God Of Revels")
	card.ManaCost = "{3}{R}{G}"
	card.Types = []string{"ENCHANTMENT", "CREATURE"}
	card.Subtypes = []string{"GOD"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "6"
	card.Toughness = "5"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordIndestructible)
	card.AddAbility(ability0)
	return card, nil
}
