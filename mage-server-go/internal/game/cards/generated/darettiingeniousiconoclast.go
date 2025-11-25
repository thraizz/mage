package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/token"
)

func init() {
	cards.Register("Daretti Ingenious Iconoclast", NewDarettiIngeniousIconoclast)
}

// NewDarettiIngeniousIconoclast creates a Daretti Ingenious Iconoclast
// {1}{B}{R} - PLANESWALKER
func NewDarettiIngeniousIconoclast(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Daretti Ingenious Iconoclast")
	card.ManaCost = "{1}{B}{R}"
	card.Types = []string{"PLANESWALKER"}
	card.Subtypes = []string{"DARETTI"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Loyalty = "3"
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement triggered ability: LoyaltyAbility
	//   - Effect: DoIfCostPaid(                         new DestroyTargetEffect()...)
	//
	// Targets:
	//   - abilities.NewTargetRequirement(1, 1, abilities.NewPermanentTargetFilter())
	// card.AddAbility(ability0)
	token1_0, err := token.GetToken("DarettiConstructToken")
	if err != nil {
		return nil, err
	}
	ability1, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewCreateTokenEffect(token1_0)).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability1)
	return card, nil
}
