package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Koth Of The Hammer", NewKothOfTheHammer)
}

// NewKothOfTheHammer creates a Koth Of The Hammer
// {2}{R}{R} - PLANESWALKER
func NewKothOfTheHammer(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Koth Of The Hammer")
	card.ManaCost = "{2}{R}{R}"
	card.Types = []string{"PLANESWALKER"}
	card.Subtypes = []string{"KOTH", "ELEMENTAL"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "4"
	card.Toughness = "4"
	card.Loyalty = "3"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		// TODO: UntapTargetEffect with complex parameters
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	return card, nil
}
