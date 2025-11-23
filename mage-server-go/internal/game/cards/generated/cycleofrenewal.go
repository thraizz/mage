package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Cycle Of Renewal", NewCycleOfRenewal)
}

// NewCycleOfRenewal creates a Cycle Of Renewal
// {2}{G} - INSTANT
func NewCycleOfRenewal(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Cycle Of Renewal")
	card.ManaCost = "{2}{G}"
	card.Types = []string{"INSTANT"}
	card.Subtypes = []string{"LESSON"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - SacrificeControllerEffect(StaticFilters.FILTER_LAND, 1, "")
	// card.AddAbility(ability0)
	return card, nil
}
