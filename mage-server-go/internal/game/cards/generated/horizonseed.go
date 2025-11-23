package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Horizon Seed", NewHorizonSeed)
}

// NewHorizonSeed creates a Horizon Seed
// {4}{W} - CREATURE
func NewHorizonSeed(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Horizon Seed")
	card.ManaCost = "{4}{W}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"SPIRIT"}
	card.Power = "2"
	card.Toughness = "1"
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - RegenerateTargetEffect()
	// card.AddAbility(ability0)
	return card, nil
}
