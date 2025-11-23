package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Phantasmal Image", NewPhantasmalImage)
}

// NewPhantasmalImage creates a Phantasmal Image
// {1}{U} - CREATURE
func NewPhantasmalImage(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Phantasmal Image")
	card.ManaCost = "{1}{U}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"ILLUSION"}
	card.Power = "0"
	card.Toughness = "0"
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - CopyPermanentEffect(StaticFilters.FILTER_PERMANENT_CREATURE, phantasma...)
	// card.AddAbility(ability0)
	return card, nil
}
