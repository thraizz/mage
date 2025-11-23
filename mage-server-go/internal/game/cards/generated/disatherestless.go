package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/token"
)

func init() {
	cards.Register("Disa The Restless", NewDisaTheRestless)
}

// NewDisaTheRestless creates a Disa The Restless
// {2}{B}{R}{G} - CREATURE
func NewDisaTheRestless(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Disa The Restless")
	card.ManaCost = "{2}{B}{R}{G}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"HUMAN", "SCOUT"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "5"
	card.Toughness = "6"
	card.SetCode = "M21"
	card.Rarity = "common"

	token0_0, err := token.GetToken("TarmogoyfToken")
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
