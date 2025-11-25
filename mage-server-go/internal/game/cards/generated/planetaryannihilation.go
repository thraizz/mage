package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Planetary Annihilation", NewPlanetaryAnnihilation)
}

// NewPlanetaryAnnihilation creates a Planetary Annihilation
// {3}{R}{R} - SORCERY
func NewPlanetaryAnnihilation(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Planetary Annihilation")
	card.ManaCost = "{3}{R}{R}"
	card.Types = []string{"SORCERY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - DamageAllEffect(6, StaticFilters.FILTER_PERMANENT_CREATURE)
	// card.AddAbility(ability0)
	return card, nil
}
