package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Gravebreaker Lamia", NewGravebreakerLamia)
}

// NewGravebreakerLamia creates a Gravebreaker Lamia
// {4}{B} - ENCHANTMENT CREATURE
// Lifelink
func NewGravebreakerLamia(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Gravebreaker Lamia")
	card.ManaCost = "{4}{B}"
	card.Types = []string{"ENCHANTMENT", "CREATURE"}
	card.Subtypes = []string{"SNAKE", "LAMIA"}
	card.Power = "4"
	card.Toughness = "4"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordLifelink)
	card.AddAbility(ability0)
	return card, nil
}
