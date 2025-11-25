package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Call Forth The Tempest", NewCallForthTheTempest)
}

// NewCallForthTheTempest creates a Call Forth The Tempest
// {5}{R}{R}{R} - SORCERY
func NewCallForthTheTempest(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Call Forth The Tempest")
	card.ManaCost = "{5}{R}{R}{R}"
	card.Types = []string{"SORCERY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - DamageAllEffect(                 CallForthTheTempestDynamicValue.i...)
	// card.AddAbility(ability0)
	return card, nil
}
