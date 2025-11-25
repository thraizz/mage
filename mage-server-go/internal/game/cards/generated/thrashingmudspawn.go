package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Thrashing Mudspawn", NewThrashingMudspawn)
}

// NewThrashingMudspawn creates a Thrashing Mudspawn
// {3}{B}{B} - CREATURE
func NewThrashingMudspawn(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Thrashing Mudspawn")
	card.ManaCost = "{3}{B}{B}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"BEAST"}
	card.Power = "4"
	card.Toughness = "4"
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement triggered ability: DealtDamageToSourceTriggeredAbility
	//   - Effect: ThrashingMudspawnEffect()
	// card.AddAbility(ability0)
	return card, nil
}
