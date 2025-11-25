package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Nyleas Intervention", NewNyleasIntervention)
}

// NewNyleasIntervention creates a Nyleas Intervention
// {X}{G}{G} - SORCERY
func NewNyleasIntervention(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Nyleas Intervention")
	card.ManaCost = "{X}{G}{G}"
	card.Types = []string{"SORCERY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - DamageAllEffect(xValue, filter)
	// card.AddAbility(ability0)
	return card, nil
}
