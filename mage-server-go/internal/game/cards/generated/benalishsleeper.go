package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Benalish Sleeper", NewBenalishSleeper)
}

// NewBenalishSleeper creates a Benalish Sleeper
// {1}{W} - CREATURE
func NewBenalishSleeper(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Benalish Sleeper")
	card.ManaCost = "{1}{W}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"PHYREXIAN", "HUMAN", "SOLDIER"}
	card.Power = "3"
	card.Toughness = "1"
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - SacrificeAllEffect(StaticFilters.FILTER_PERMANENT_CREATURE)
	// card.AddAbility(ability0)
	return card, nil
}
