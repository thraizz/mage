package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Leshracs Sigil", NewLeshracsSigil)
}

// NewLeshracsSigil creates a Leshracs Sigil
// {B}{B} - ENCHANTMENT
func NewLeshracsSigil(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Leshracs Sigil")
	card.ManaCost = "{B}{B}"
	card.Types = []string{"ENCHANTMENT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - DoIfCostPaid(new LookTargetHandChooseDiscardEffect(), new ManaC...)
	// card.AddAbility(ability0)
	return card, nil
}