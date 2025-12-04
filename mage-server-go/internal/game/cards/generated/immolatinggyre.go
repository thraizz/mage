package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Immolating Gyre", NewImmolatingGyre)
}

// NewImmolatingGyre creates a Immolating Gyre
// {4}{R}{R} - SORCERY
func NewImmolatingGyre(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Immolating Gyre")
	card.ManaCost = "{4}{R}{R}"
	card.Types = []string{"SORCERY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - DamageAllEffect(xValue, filter)
	// card.AddAbility(ability0)
	return card, nil
}
