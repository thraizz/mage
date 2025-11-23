package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Mirrormade", NewMirrormade)
}

// NewMirrormade creates a Mirrormade
// {1}{U}{U} - ENCHANTMENT
func NewMirrormade(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Mirrormade")
	card.ManaCost = "{1}{U}{U}"
	card.Types = []string{"ENCHANTMENT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - CopyPermanentEffect(StaticFilters.FILTER_PERMANENT_ARTIFACT_OR_ENCHANT...)
	// card.AddAbility(ability0)
	return card, nil
}
