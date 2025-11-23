package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Crystal Seer", NewCrystalSeer)
}

// NewCrystalSeer creates a Crystal Seer
// {4}{U} - CREATURE
func NewCrystalSeer(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Crystal Seer")
	card.ManaCost = "{4}{U}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"VEDALKEN", "WIZARD"}
	card.Power = "2"
	card.Toughness = "2"
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - LookLibraryControllerEffect(4)
	// card.AddAbility(ability0)
	return card, nil
}
