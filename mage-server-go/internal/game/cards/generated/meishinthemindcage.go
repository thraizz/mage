package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Meishin The Mind Cage", NewMeishinTheMindCage)
}

// NewMeishinTheMindCage creates a Meishin The Mind Cage
// {4}{U}{U}{U} - ENCHANTMENT
func NewMeishinTheMindCage(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Meishin The Mind Cage")
	card.ManaCost = "{4}{U}{U}{U}"
	card.Types = []string{"ENCHANTMENT"}
	card.Supertypes = []string{"LEGENDARY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewBoostEffect(new SignInversionDynamicValue(CardsInControllerHandCount.ANY), StaticValue.get(0), false)).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	return card, nil
}