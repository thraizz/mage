package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Slab Hammer", NewSlabHammer)
}

// NewSlabHammer creates a Slab Hammer
// {2} - ARTIFACT
func NewSlabHammer(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Slab Hammer")
	card.ManaCost = "{2}"
	card.Types = []string{"ARTIFACT"}
	card.Subtypes = []string{"EQUIPMENT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement triggered ability: AttacksAttachedTriggeredAbility
	//   - Effect: DoIfCostPaid(new BoostEquippedEffect(2, 2, Duration.EndOfTurn),...)
	// card.AddAbility(ability0)
	return card, nil
}
