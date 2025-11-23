package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Valgavoths Lair", NewValgavothsLair)
}

// NewValgavothsLair creates a Valgavoths Lair
//   - ENCHANTMENT LAND
//
// Hexproof
func NewValgavothsLair(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Valgavoths Lair")
	card.ManaCost = ""
	card.Types = []string{"ENCHANTMENT", "LAND"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordHexproof)
	card.AddAbility(ability0)
	return card, nil
}
