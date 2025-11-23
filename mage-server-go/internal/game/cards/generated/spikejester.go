package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Spike Jester", NewSpikeJester)
}

// NewSpikeJester creates a Spike Jester
// {B}{R} - CREATURE
// Haste
func NewSpikeJester(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Spike Jester")
	card.ManaCost = "{B}{R}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"GOBLIN", "WARRIOR"}
	card.Power = "3"
	card.Toughness = "1"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordHaste)
	card.AddAbility(ability0)
	return card, nil
}
