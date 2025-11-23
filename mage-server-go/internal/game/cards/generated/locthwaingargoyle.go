package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/effects"
)

func init() {
	cards.Register("Locthwain Gargoyle", NewLocthwainGargoyle)
}

// NewLocthwainGargoyle creates a Locthwain Gargoyle
// {1} - ARTIFACT CREATURE
func NewLocthwainGargoyle(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Locthwain Gargoyle")
	card.ManaCost = "{1}"
	card.Types = []string{"ARTIFACT", "CREATURE"}
	card.Subtypes = []string{"GARGOYLE"}
	card.Power = "0"
	card.Toughness = "3"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewActivatedAbilityBuilder(card.ID).
		AddManaCost("{4}").
		AddEffect(abilities.NewBoostEffect(2, 0)).
		AddEffect(abilities.NewGrantAbilityEffect("FlyingAbility", effects.DurationEndOfTurn)).
		Build()
	card.AddAbility(ability0)
	return card, nil
}
