package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Dina Soul Steeper", NewDinaSoulSteeper)
}

// NewDinaSoulSteeper creates a Dina Soul Steeper
// {B}{G} - CREATURE
func NewDinaSoulSteeper(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Dina Soul Steeper")
	card.ManaCost = "{B}{G}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"DRYAD", "DRUID"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "1"
	card.Toughness = "3"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewActivatedAbilityBuilder(card.ID).
		AddManaCost("{1}").
		AddEffect(abilities.NewBoostEffect(SacrificeCostCreaturesPower.instance, StaticValue.get(0))).
		Build()
	card.AddAbility(ability0)
	return card, nil
}
