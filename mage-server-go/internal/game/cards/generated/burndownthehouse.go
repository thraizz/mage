package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Burn Down The House", NewBurnDownTheHouse)
}

// NewBurnDownTheHouse creates a Burn Down The House
// {3}{R}{R} - SORCERY
func NewBurnDownTheHouse(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Burn Down The House")
	card.ManaCost = "{3}{R}{R}"
	card.Types = []string{"SORCERY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - DamageAllEffect(5, filter)
	// card.AddAbility(ability0)
	return card, nil
}
