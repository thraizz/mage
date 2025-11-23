package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Stromkirk Condemned", NewStromkirkCondemned)
}

// NewStromkirkCondemned creates a Stromkirk Condemned
// {B}{B} - CREATURE
func NewStromkirkCondemned(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Stromkirk Condemned")
	card.ManaCost = "{B}{B}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"VAMPIRE", "HORROR"}
	card.Power = "2"
	card.Toughness = "2"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewBoostEffect(1, 1, filter, false)).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	return card, nil
}
