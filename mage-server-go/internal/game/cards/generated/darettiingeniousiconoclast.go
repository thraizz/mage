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

	token0_0, err := token.GetToken("DarettiConstructToken")
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
	return card, nil
}
