package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Ultimate Magic Meteor", NewUltimateMagicMeteor)
}

// NewUltimateMagicMeteor creates a Ultimate Magic Meteor
// {5}{R} - SORCERY
func NewUltimateMagicMeteor(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Ultimate Magic Meteor")
	card.ManaCost = "{5}{R}"
	card.Types = []string{"SORCERY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - DamageAllEffect(7, StaticFilters.FILTER_PERMANENT_CREATURE)
	// card.AddAbility(ability0)
	return card, nil
}
