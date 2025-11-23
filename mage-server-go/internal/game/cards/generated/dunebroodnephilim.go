package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/token"
)

func init() {
	cards.Register("Dune Brood Nephilim", NewDuneBroodNephilim)
}

// NewDuneBroodNephilim creates a Dune Brood Nephilim
// {B}{R}{G}{W} - CREATURE
func NewDuneBroodNephilim(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Dune Brood Nephilim")
	card.ManaCost = "{B}{R}{G}{W}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"NEPHILIM"}
	card.Power = "3"
	card.Toughness = "3"
	card.SetCode = "M21"
	card.Rarity = "common"

	token0_0, err := token.GetToken("DuneBroodNephilimToken")
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
