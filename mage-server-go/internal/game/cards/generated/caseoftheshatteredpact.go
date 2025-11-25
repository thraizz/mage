package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/effects"
)

func init() {
	cards.Register("Case Of The Shattered Pact", NewCaseOfTheShatteredPact)
}

// NewCaseOfTheShatteredPact creates a Case Of The Shattered Pact
// {2} - ENCHANTMENT
func NewCaseOfTheShatteredPact(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Case Of The Shattered Pact")
	card.ManaCost = "{2}"
	card.Types = []string{"ENCHANTMENT"}
	card.Subtypes = []string{"CASE"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewGrantAbilityEffect("FlyingAbility", effects.DurationEndOfTurn)).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	return card, nil
}
