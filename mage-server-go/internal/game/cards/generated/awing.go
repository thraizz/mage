package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("A Wing", NewAWing)
}

// NewAWing creates a A Wing
// {2}{R} - ARTIFACT CREATURE
// Haste
func NewAWing(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "A Wing")
	card.ManaCost = "{2}{R}"
	card.Types = []string{"ARTIFACT", "CREATURE"}
	card.Subtypes = []string{"REBEL", "STARSHIP"}
	card.Power = "2"
	card.Toughness = "2"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordHaste)
	card.AddAbility(ability0)
	// TODO: Implement activated ability with unmapped effects
	//   - RemoveFromCombatSourceEffect()
	//
	// Costs:
	//   - AddManaCost("{1}")
	// card.AddAbility(ability1)
	return card, nil
}
