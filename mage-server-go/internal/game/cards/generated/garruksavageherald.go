package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/effects"
)

func init() {
	cards.Register("Garruk Savage Herald", NewGarrukSavageHerald)
}

// NewGarrukSavageHerald creates a Garruk Savage Herald
// {4}{G}{G} - PLANESWALKER
func NewGarrukSavageHerald(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Garruk Savage Herald")
	card.ManaCost = "{4}{G}{G}"
	card.Types = []string{"PLANESWALKER"}
	card.Subtypes = []string{"GARRUK"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Loyalty = "5"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewGrantAbilityEffect("DamageAsThoughNotBlockedAbility", effects.DurationEndOfTurn)).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	return card, nil
}