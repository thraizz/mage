package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/token"
)

func init() {
	cards.Register("Dovin Grand Arbiter", NewDovinGrandArbiter)
}

// NewDovinGrandArbiter creates a Dovin Grand Arbiter
// {1}{W}{U} - PLANESWALKER
func NewDovinGrandArbiter(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Dovin Grand Arbiter")
	card.ManaCost = "{1}{W}{U}"
	card.Types = []string{"PLANESWALKER"}
	card.Subtypes = []string{"DOVIN"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Loyalty = "3"
	card.SetCode = "M21"
	card.Rarity = "common"

	token0_0, err := token.GetToken("ThopterColorlessToken")
	if err != nil {
		return nil, err
	}
	ability0, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewCreateTokenEffect(token0_0)).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	// TODO: Implement spell ability with unmapped effects
	//   - LookLibraryAndPickControllerEffect(                 10, 3, PutCards.HAND, PutCards.BO...)
	// card.AddAbility(ability1)
	return card, nil
}