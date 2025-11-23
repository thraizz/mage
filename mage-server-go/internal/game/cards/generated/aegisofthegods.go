package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Aegis Of The Gods", NewAegisOfTheGods)
}

// NewAegisOfTheGods creates a Aegis Of The Gods
// {1}{W} - ENCHANTMENT CREATURE
// Hexproof
func NewAegisOfTheGods(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Aegis Of The Gods")
	card.ManaCost = "{1}{W}"
	card.Types = []string{"ENCHANTMENT", "CREATURE"}
	card.Subtypes = []string{"HUMAN", "SOLDIER"}
	card.Power = "2"
	card.Toughness = "1"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordHexproof)
	card.AddAbility(ability0)
	return card, nil
}
