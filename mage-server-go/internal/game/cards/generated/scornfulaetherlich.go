package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/effects"
)

func init() {
	cards.Register("Scornful Aether Lich", NewScornfulAetherLich)
}

// NewScornfulAetherLich creates a Scornful Aether Lich
// {3}{U} - ARTIFACT CREATURE
func NewScornfulAetherLich(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Scornful Aether Lich")
	card.ManaCost = "{3}{U}"
	card.Types = []string{"ARTIFACT", "CREATURE"}
	card.Subtypes = []string{"ZOMBIE", "WIZARD"}
	card.Power = "2"
	card.Toughness = "4"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewActivatedAbilityBuilder(card.ID).
		AddEffect(abilities.NewGrantAbilityEffect("FearAbility", effects.DurationEndOfTurn)).
		AddEffect(abilities.NewGrantAbilityEffect("VigilanceAbility", effects.DurationEndOfTurn)).
		Build()
	card.AddAbility(ability0)
	return card, nil
}
