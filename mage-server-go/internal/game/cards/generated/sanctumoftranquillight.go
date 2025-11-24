package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Sanctum Of Tranquil Light", NewSanctumOfTranquilLight)
}

// NewSanctumOfTranquilLight creates a Sanctum Of Tranquil Light
// {W} - ENCHANTMENT
func NewSanctumOfTranquilLight(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Sanctum Of Tranquil Light")
	card.ManaCost = "{W}"
	card.Types = []string{"ENCHANTMENT"}
	card.Subtypes = []string{"SHRINE"}
	card.Supertypes = []string{"LEGENDARY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewActivatedAbilityBuilder(card.ID).
		// TODO: TapTargetEffect with complex parameters
		Build()
	card.AddAbility(ability0)
	return card, nil
}
