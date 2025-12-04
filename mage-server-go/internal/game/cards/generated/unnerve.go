package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Unnerve", NewUnnerve)
}

// NewUnnerve creates a Unnerve
// {3}{B} - SORCERY
func NewUnnerve(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Unnerve")
	card.ManaCost = "{3}{B}"
	card.Types = []string{"SORCERY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - DiscardEachPlayerEffect(StaticValue.get(2), false, TargetController.OPPONE...)
	// card.AddAbility(ability0)
	return card, nil
}
