package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/token"
)

func init() {
	cards.Register("Sarkhan Fireblood", NewSarkhanFireblood)
}

// NewSarkhanFireblood creates a Sarkhan Fireblood
// {1}{R}{R} - PLANESWALKER
func NewSarkhanFireblood(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Sarkhan Fireblood")
	card.ManaCost = "{1}{R}{R}"
	card.Types = []string{"PLANESWALKER"}
	card.Subtypes = []string{"SARKHAN"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Loyalty = "3"
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - DoIfCostPaid(                 new DrawCardSourceControllerEffec...)
	// card.AddAbility(ability0)
	ability1, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewCreateTokenEffect(new DragonToken2(), 4)).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability1)
	return card, nil
}