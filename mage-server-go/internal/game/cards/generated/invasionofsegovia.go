package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/token"
)

func init() {
	cards.Register("Invasion Of Segovia", NewInvasionOfSegovia)
}

// NewInvasionOfSegovia creates a Invasion Of Segovia
// {2}{U} - BATTLE
func NewInvasionOfSegovia(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Invasion Of Segovia")
	card.ManaCost = "{2}{U}"
	card.Types = []string{"BATTLE"}
	card.Subtypes = []string{"SIEGE"}
	card.SetCode = "M21"
	card.Rarity = "common"

	token0_0, err := token.GetToken("Kraken11Token")
	if err != nil {
		return nil, err
	}
	ability0, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewCreateTokenEffect(token0_0, 2)).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	return card, nil
}
