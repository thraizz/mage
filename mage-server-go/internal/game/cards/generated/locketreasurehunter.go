package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Locke Treasure Hunter", NewLockeTreasureHunter)
}

// NewLockeTreasureHunter creates a Locke Treasure Hunter
// {1}{B}{R} - CREATURE
func NewLockeTreasureHunter(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Locke Treasure Hunter")
	card.ManaCost = "{1}{B}{R}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"HUMAN", "ROGUE"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "2"
	card.Toughness = "3"
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement activated ability with unmapped effects
	//   - LockeTreasureHunterEffect()
	//
	// Costs:
	//   - AddManaCost("{0}")
	// card.AddAbility(ability0)
	return card, nil
}
