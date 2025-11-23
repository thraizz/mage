package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/counters"
)

func init() {
	cards.Register("Lagonna Band Trailblazer", NewLagonnaBandTrailblazer)
}

// NewLagonnaBandTrailblazer creates a Lagonna Band Trailblazer
// {W} - CREATURE
func NewLagonnaBandTrailblazer(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Lagonna Band Trailblazer")
	card.ManaCost = "{W}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"CENTAUR", "SCOUT"}
	card.Power = "0"
	card.Toughness = "4"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewAddCountersSourceEffect(counters.CounterTypeP1P1.CreateInstance(1))).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	return card, nil
}
