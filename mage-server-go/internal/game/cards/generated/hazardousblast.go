package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Hazardous Blast", NewHazardousBlast)
}

// NewHazardousBlast creates a Hazardous Blast
// {3}{R} - SORCERY
func NewHazardousBlast(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Hazardous Blast")
	card.ManaCost = "{3}{R}"
	card.Types = []string{"SORCERY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - DamageAllEffect(                 1, StaticFilters.FILTER_OPPONENTS...)
	// card.AddAbility(ability0)
	return card, nil
}
