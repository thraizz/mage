package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/token"
)

func init() {
	cards.Register("Ashiok Wicked Manipulator", NewAshiokWickedManipulator)
}

// NewAshiokWickedManipulator creates a Ashiok Wicked Manipulator
// {3}{B}{B} - PLANESWALKER
func NewAshiokWickedManipulator(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Ashiok Wicked Manipulator")
	card.ManaCost = "{3}{B}{B}"
	card.Types = []string{"PLANESWALKER"}
	card.Subtypes = []string{"ASHIOK"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Loyalty = "5"
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - LookLibraryAndPickControllerEffect(2, 1, PutCards.EXILED, PutCards.HAND)
	// card.AddAbility(ability0)
	token1_0, err := token.GetToken("AshiokWickedManipulatorNightmareToken")
	if err != nil {
		return nil, err
	}
	ability1, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewCreateTokenEffect(token1_0, 2)).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability1)
	return card, nil
}
