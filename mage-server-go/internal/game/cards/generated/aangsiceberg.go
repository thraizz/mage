package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Aangs Iceberg", NewAangsIceberg)
}

// NewAangsIceberg creates a Aangs Iceberg
// {2}{W} - ENCHANTMENT
// Flash
func NewAangsIceberg(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Aangs Iceberg")
	card.ManaCost = "{2}{W}"
	card.Types = []string{"ENCHANTMENT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordFlash)
	card.AddAbility(ability0)
	return card, nil
}
