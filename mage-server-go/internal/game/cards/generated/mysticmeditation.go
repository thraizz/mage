package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Mystic Meditation", NewMysticMeditation)
}

// NewMysticMeditation creates a Mystic Meditation
// {3}{U} - SORCERY
func NewMysticMeditation(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Mystic Meditation")
	card.ManaCost = "{3}{U}"
	card.Types = []string{"SORCERY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - DiscardControllerEffect(2)
	//   - DoIfCostPaid(                 null, new DiscardControllerEffect...)
	// card.AddAbility(ability0)
	return card, nil
}
