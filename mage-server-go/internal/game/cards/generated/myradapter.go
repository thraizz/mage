package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Myr Adapter", NewMyrAdapter)
}

// NewMyrAdapter creates a Myr Adapter
// {3} - ARTIFACT CREATURE
func NewMyrAdapter(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Myr Adapter")
	card.ManaCost = "{3}"
	card.Types = []string{"ARTIFACT", "CREATURE"}
	card.Subtypes = []string{"MYR"}
	card.Power = "1"
	card.Toughness = "1"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewBoostEffect(amount, amount)).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	return card, nil
}
