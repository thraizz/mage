package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/counters"
)

func init() {
	cards.Register("Meathook Massacre I I", NewMeathookMassacreII)
}

// NewMeathookMassacreII creates a Meathook Massacre I I
// {X}{X}{B}{B}{B}{B} - ENCHANTMENT
func NewMeathookMassacreII(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Meathook Massacre I I")
	card.ManaCost = "{X}{X}{B}{B}{B}{B}"
	card.Types = []string{"ENCHANTMENT"}
	card.Supertypes = []string{"LEGENDARY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - SacrificeAllEffect(GetXValue.instance, StaticFilters.FILTER_PERMANENT...)
	// card.AddAbility(ability0)
	// TODO: Implement spell ability with unmapped effects
	//   - DoIfCostPaid(                         new ReturnFromGraveyardTo...)
	// card.AddAbility(ability1)
	return card, nil
}
