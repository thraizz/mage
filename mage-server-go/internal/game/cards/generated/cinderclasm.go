package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Cinderclasm", NewCinderclasm)
}

// NewCinderclasm creates a Cinderclasm
// {1}{R} - INSTANT
func NewCinderclasm(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Cinderclasm")
	card.ManaCost = "{1}{R}"
	card.Types = []string{"INSTANT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - DamageAllEffect(2, StaticFilters.FILTER_PERMANENT_CREATURE)
	//   - DamageAllEffect(1, StaticFilters.FILTER_PERMANENT_CREATURE)
	// card.AddAbility(ability0)
	return card, nil
}
