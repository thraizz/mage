package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Lorthos The Tidemaker", NewLorthosTheTidemaker)
}

// NewLorthosTheTidemaker creates a Lorthos The Tidemaker
// {5}{U}{U}{U} - CREATURE
func NewLorthosTheTidemaker(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Lorthos The Tidemaker")
	card.ManaCost = "{5}{U}{U}{U}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"OCTOPUS"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "8"
	card.Toughness = "8"
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - DoIfCostPaid(new TapTargetEffect(), new GenericManaCost(8), "Pa...)
	// card.AddAbility(ability0)
	return card, nil
}