package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Lord Windgrace", NewLordWindgrace)
}

// NewLordWindgrace creates a Lord Windgrace
// {2}{B}{R}{G} - PLANESWALKER
func NewLordWindgrace(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Lord Windgrace")
	card.ManaCost = "{2}{B}{R}{G}"
	card.Types = []string{"PLANESWALKER"}
	card.Subtypes = []string{"WINDGRACE"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Loyalty = "5"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewDestroyEffect()).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	return card, nil
}
