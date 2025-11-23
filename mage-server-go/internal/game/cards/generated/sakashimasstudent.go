package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Sakashimas Student", NewSakashimasStudent)
}

// NewSakashimasStudent creates a Sakashimas Student
// {2}{U}{U} - CREATURE
func NewSakashimasStudent(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Sakashimas Student")
	card.ManaCost = "{2}{U}{U}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"HUMAN", "NINJA"}
	card.Power = "0"
	card.Toughness = "0"
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - CopyPermanentEffect(StaticFilters.FILTER_PERMANENT_CREATURE, new AddSu...)
	// card.AddAbility(ability0)
	return card, nil
}
