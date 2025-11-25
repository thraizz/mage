package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Xantcha Sleeper Agent", NewXantchaSleeperAgent)
}

// NewXantchaSleeperAgent creates a Xantcha Sleeper Agent
// {1}{B}{R} - CREATURE
func NewXantchaSleeperAgent(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Xantcha Sleeper Agent")
	card.ManaCost = "{1}{B}{R}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"PHYREXIAN", "MINION"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "5"
	card.Toughness = "5"
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement activated ability with unmapped effects
	//   - LoseLifePermanentControllerEffect()
	//
	// Costs:
	//   - AddManaCost("{3}")
	// card.AddAbility(ability0)
	return card, nil
}
