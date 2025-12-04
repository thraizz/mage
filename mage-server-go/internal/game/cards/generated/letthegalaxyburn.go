package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Let The Galaxy Burn", NewLetTheGalaxyBurn)
}

// NewLetTheGalaxyBurn creates a Let The Galaxy Burn
// {X}{5}{R} - SORCERY
func NewLetTheGalaxyBurn(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Let The Galaxy Burn")
	card.ManaCost = "{X}{5}{R}"
	card.Types = []string{"SORCERY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - DamageAllEffect(xValue, filter)
	// card.AddAbility(ability0)
	return card, nil
}
