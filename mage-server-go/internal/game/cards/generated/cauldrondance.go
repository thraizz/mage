package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Cauldron Dance", NewCauldronDance)
}

// NewCauldronDance creates a Cauldron Dance
// {4}{B}{R} - INSTANT
func NewCauldronDance(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Cauldron Dance")
	card.ManaCost = "{4}{B}{R}"
	card.Types = []string{"INSTANT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - SacrificeTargetEffect("sacrifice " + card.getName(), source.getControlle...)
	// card.AddAbility(ability0)
	return card, nil
}
