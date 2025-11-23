package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Arterial Flow", NewArterialFlow)
}

// NewArterialFlow creates a Arterial Flow
// {1}{B}{B} - SORCERY
func NewArterialFlow(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Arterial Flow")
	card.ManaCost = "{1}{B}{B}"
	card.Types = []string{"SORCERY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - DiscardEachPlayerEffect(StaticValue.get(2), false, TargetController.OPPONE...)
	// card.AddAbility(ability0)
	return card, nil
}
