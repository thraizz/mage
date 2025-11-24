package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Ultimate Green Goblin", NewUltimateGreenGoblin)
}

// NewUltimateGreenGoblin creates a Ultimate Green Goblin
// {1}{B/R}{B/R} - CREATURE
func NewUltimateGreenGoblin(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Ultimate Green Goblin")
	card.ManaCost = "{1}{B/R}{B/R}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"GOBLIN", "VILLAIN"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "5"
	card.Toughness = "4"
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - DiscardControllerEffect(1)
	// card.AddAbility(ability0)
	return card, nil
}