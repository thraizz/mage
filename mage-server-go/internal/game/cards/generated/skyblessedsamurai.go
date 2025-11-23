package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Sky Blessed Samurai", NewSkyBlessedSamurai)
}

// NewSkyBlessedSamurai creates a Sky Blessed Samurai
// {6}{W} - ENCHANTMENT CREATURE
// Flying
func NewSkyBlessedSamurai(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Sky Blessed Samurai")
	card.ManaCost = "{6}{W}"
	card.Types = []string{"ENCHANTMENT", "CREATURE"}
	card.Subtypes = []string{"HUMAN", "SAMURAI"}
	card.Power = "4"
	card.Toughness = "4"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordFlying)
	card.AddAbility(ability0)
	return card, nil
}
