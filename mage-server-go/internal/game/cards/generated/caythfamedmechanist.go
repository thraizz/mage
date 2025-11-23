package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/effects"
)

func init() {
	cards.Register("Cayth Famed Mechanist", NewCaythFamedMechanist)
}

// NewCaythFamedMechanist creates a Cayth Famed Mechanist
// {1}{U}{R}{W} - CREATURE
func NewCaythFamedMechanist(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Cayth Famed Mechanist")
	card.ManaCost = "{1}{U}{R}{W}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"DWARF", "ARTIFICER"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "3"
	card.Toughness = "3"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewGrantAbilityEffect(new FabricateAbility(1), true)).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	return card, nil
}