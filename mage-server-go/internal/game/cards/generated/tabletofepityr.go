package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Tablet Of Epityr", NewTabletOfEpityr)
}

// NewTabletOfEpityr creates a Tablet Of Epityr
// {1} - ARTIFACT
func NewTabletOfEpityr(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Tablet Of Epityr")
	card.ManaCost = "{1}"
	card.Types = []string{"ARTIFACT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - DoIfCostPaid(new GainLifeEffect(1), new GenericManaCost(1))
	// card.AddAbility(ability0)
	return card, nil
}