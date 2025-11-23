package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Soulblade Corrupter", NewSoulbladeCorrupter)
}

// NewSoulbladeCorrupter creates a Soulblade Corrupter
// {4}{B} - CREATURE
// Deathtouch
func NewSoulbladeCorrupter(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Soulblade Corrupter")
	card.ManaCost = "{4}{B}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"HUMAN", "WARRIOR"}
	card.Power = "3"
	card.Toughness = "3"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordDeathtouch)
	card.AddAbility(ability0)
	return card, nil
}
