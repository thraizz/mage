package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("T45 Power Armor", NewT45PowerArmor)
}

// NewT45PowerArmor creates a T45 Power Armor
// {2} - ARTIFACT
func NewT45PowerArmor(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "T45 Power Armor")
	card.ManaCost = "{2}"
	card.Types = []string{"ARTIFACT"}
	card.Subtypes = []string{"EQUIPMENT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0, err := abilities.NewEquipAbility(card.ID, "{3}", false)
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	// TODO: Implement spell ability with unmapped effects
	//   - DoIfCostPaid(                         new T45PowerArmorEffect()...)
	// card.AddAbility(ability1)
	return card, nil
}
