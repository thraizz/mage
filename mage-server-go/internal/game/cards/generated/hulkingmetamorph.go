package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Hulking Metamorph", NewHulkingMetamorph)
}

// NewHulkingMetamorph creates a Hulking Metamorph
// {9} - ARTIFACT CREATURE
func NewHulkingMetamorph(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Hulking Metamorph")
	card.ManaCost = "{9}"
	card.Types = []string{"ARTIFACT", "CREATURE"}
	card.Subtypes = []string{"SHAPESHIFTER"}
	card.Power = "7"
	card.Toughness = "7"
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - CopyPermanentEffect(                         StaticFilters.FILTER_CONT...)
	// card.AddAbility(ability0)
	return card, nil
}
