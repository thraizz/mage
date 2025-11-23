package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/effects"
)

func init() {
	cards.Register("Cephalid Inkshrouder", NewCephalidInkshrouder)
}

// NewCephalidInkshrouder creates a Cephalid Inkshrouder
// {2}{U} - CREATURE
func NewCephalidInkshrouder(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Cephalid Inkshrouder")
	card.ManaCost = "{2}{U}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"OCTOPUS"}
	card.Power = "2"
	card.Toughness = "1"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewActivatedAbilityBuilder(card.ID).
		AddEffect(abilities.NewGrantAbilityEffect("ShroudAbility", effects.DurationEndOfTurn)).
		Build()
	card.AddAbility(ability0)
	return card, nil
}
