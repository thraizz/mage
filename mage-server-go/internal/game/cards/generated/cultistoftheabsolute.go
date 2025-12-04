package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Cultist Of The Absolute", NewCultistOfTheAbsolute)
}

// NewCultistOfTheAbsolute creates a Cultist Of The Absolute
// {B} - ENCHANTMENT
func NewCultistOfTheAbsolute(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Cultist Of The Absolute")
	card.ManaCost = "{B}"
	card.Types = []string{"ENCHANTMENT"}
	card.Subtypes = []string{"BACKGROUND"}
	card.Supertypes = []string{"LEGENDARY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - SacrificeControllerEffect(                         StaticFilters.FILTER_PERM...)
	// card.AddAbility(ability0)
	return card, nil
}
