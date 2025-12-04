package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Phyrexian Metamorph", NewPhyrexianMetamorph)
}

// NewPhyrexianMetamorph creates a Phyrexian Metamorph
// {3}{U/P} - ARTIFACT CREATURE
func NewPhyrexianMetamorph(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Phyrexian Metamorph")
	card.ManaCost = "{3}{U/P}"
	card.Types = []string{"ARTIFACT", "CREATURE"}
	card.Subtypes = []string{"PHYREXIAN", "SHAPESHIFTER"}
	card.Power = "0"
	card.Toughness = "0"
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - CopyPermanentEffect(                 StaticFilters.FILTER_PERMANENT_AR...)
	// card.AddAbility(ability0)
	return card, nil
}
