package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Celestial Sword", NewCelestialSword)
}

// NewCelestialSword creates a Celestial Sword
// {6} - ARTIFACT
func NewCelestialSword(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Celestial Sword")
	card.ManaCost = "{6}"
	card.Types = []string{"ARTIFACT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement activated ability with unmapped effects
	//   - SacrificeTargetEffect("its controller sacrifices it")
	//
	// Costs:
	//   - AddManaCost("{3}")
	//   - AddTapCost()
	// card.AddAbility(ability0)
	return card, nil
}