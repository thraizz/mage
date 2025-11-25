package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Demogorgons Clutches", NewDemogorgonsClutches)
}

// NewDemogorgonsClutches creates a Demogorgons Clutches
// {2}{B} - SORCERY
func NewDemogorgonsClutches(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Demogorgons Clutches")
	card.ManaCost = "{2}{B}"
	card.Types = []string{"SORCERY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - DiscardTargetEffect(2)
	//
	// Targets:
	//   - abilities.NewTargetRequirement(1, 1, abilities.NewOpponentTargetFilter())
	// card.AddAbility(ability0)
	return card, nil
}
