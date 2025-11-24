package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Angrath The Flame Chained", NewAngrathTheFlameChained)
}

// NewAngrathTheFlameChained creates a Angrath The Flame Chained
// {3}{B}{R} - PLANESWALKER
func NewAngrathTheFlameChained(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Angrath The Flame Chained")
	card.ManaCost = "{3}{B}{R}"
	card.Types = []string{"PLANESWALKER"}
	card.Subtypes = []string{"ANGRATH"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Loyalty = "4"
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - DiscardEachPlayerEffect(TargetController.OPPONENT)
	//   - SacrificeTargetEffect("sacrifice this", source.getControllerId())
	// card.AddAbility(ability0)
	return card, nil
}