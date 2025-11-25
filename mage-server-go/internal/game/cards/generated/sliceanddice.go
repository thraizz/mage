package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Slice And Dice", NewSliceAndDice)
}

// NewSliceAndDice creates a Slice And Dice
// {4}{R}{R} - SORCERY
func NewSliceAndDice(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Slice And Dice")
	card.ManaCost = "{4}{R}{R}"
	card.Types = []string{"SORCERY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - DamageAllEffect(4, new FilterCreaturePermanent())
	// card.AddAbility(ability0)
	// TODO: Implement spell ability with unmapped effects
	//   - DamageAllEffect(1, StaticFilters.FILTER_PERMANENT_CREATURE)
	// card.AddAbility(ability1)
	return card, nil
}
