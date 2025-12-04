package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Vicious Rumors", NewViciousRumors)
}

// NewViciousRumors creates a Vicious Rumors
// {B} - SORCERY
func NewViciousRumors(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Vicious Rumors")
	card.ManaCost = "{B}"
	card.Types = []string{"SORCERY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - DiscardEachPlayerEffect(                 StaticValue.get(1), false,       ...)
	// card.AddAbility(ability0)
	return card, nil
}
