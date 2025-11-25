package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Words Of War", NewWordsOfWar)
}

// NewWordsOfWar creates a Words Of War
// {2}{R} - ENCHANTMENT
func NewWordsOfWar(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Words Of War")
	card.ManaCost = "{2}{R}"
	card.Types = []string{"ENCHANTMENT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement activated ability with unmapped effects
	//   - WordsOfWarEffect()
	//
	// Costs:
	//   - AddManaCost("{1}")
	// card.AddAbility(ability0)
	return card, nil
}
