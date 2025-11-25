package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/token"
)

func init() {
	cards.Register("Quintorius Kand", NewQuintoriusKand)
}

// NewQuintoriusKand creates a Quintorius Kand
// {3}{R}{W} - PLANESWALKER
func NewQuintoriusKand(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Quintorius Kand")
	card.ManaCost = "{3}{R}{W}"
	card.Types = []string{"PLANESWALKER"}
	card.Subtypes = []string{"QUINTORIUS"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Loyalty = "4"
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement triggered ability: SpellCastControllerTriggeredAbility
	//   - Effect: DamagePlayersEffect()
	// card.AddAbility(ability0)
	token1_0, err := token.GetToken("Spirit32Token")
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
