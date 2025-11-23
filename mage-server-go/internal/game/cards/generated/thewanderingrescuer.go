package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/effects"
)

func init() {
	cards.Register("The Wandering Rescuer", NewTheWanderingRescuer)
}

// NewTheWanderingRescuer creates a The Wandering Rescuer
// {3}{W}{W} - CREATURE
// Flash, DoubleStrike
func NewTheWanderingRescuer(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "The Wandering Rescuer")
	card.ManaCost = "{3}{W}{W}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"HUMAN", "SAMURAI", "NOBLE"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "3"
	card.Toughness = "4"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordFlash)
	card.AddAbility(ability0)
	ability1 := abilities.NewKeywordAbility(card.ID, abilities.KeywordDoubleStrike)
	card.AddAbility(ability1)
	ability2, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewGrantAbilityEffect("HexproofAbility", effects.DurationWhileOnBattlefield)).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability2)
	return card, nil
}
